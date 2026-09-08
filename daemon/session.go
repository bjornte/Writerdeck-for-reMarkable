// Writerdeck-server — see main.go for overview.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// --- Session manager ---
// An editor session is a sub-lifecycle: xochitl stopped, keywriter running
// (with systemd-inhibit in launch-keywriter.sh holding the sleep lock).
// rmkbd itself is always-on; sessions are started/stopped on demand.

// session holds the state of one editor sub-lifecycle.
type session struct {
	mu         sync.Mutex
	active     bool
	sleeping   bool // power-button sleep: editor stopped, xochitl stays down
	sleepNote  string // note to reopen after suspend (currentNote is cleared for the phone)
	cmd        *exec.Cmd
	doneCh     chan struct{}
	editorPath string
	ec         *editorConn
}

// isSleeping returns true after a power-button sleep until wake completes.
func (s *session) isSleeping() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sleeping
}

// isActive returns true if an editor session is currently running.
func (s *session) isActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active
}

// start stops xochitl, spawns the editor, and marks the session active.
// Holds the mutex for the duration so concurrent start calls are serialized.
// Returns an error if a session is already active.
func (s *session) start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		return fmt.Errorf("session already active")
	}
	fmt.Fprintln(os.Stderr, "writerdeck-server: session: stopping xochitl")
	if err := exec.Command("systemctl", "stop", "xochitl").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: warning: stop xochitl: %v\n", err)
	}
	time.Sleep(time.Second)
	// Grab gpio-keys before Writerdeck starts so Qt evdev never sees physical
	// Home/Power/page buttons (avoids duplicate handleHome). Idle xochitl keeps
	// the buttons because we only grab while a session is active.
	grabButtonDev()
	cmd := exec.Command(s.editorPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Setpgid gives the editor+inhibit wrapper their own process group so a
	// Kill(-pgid, SIGTERM) SIGTERM fallback reaches all child processes.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		ungrabButtonDev()
		startXochitl()
		return fmt.Errorf("start editor: %w", err)
	}
	fmt.Fprintf(os.Stderr, "writerdeck-server: session: editor started (pid %d)\n", cmd.Process.Pid)
	doneCh := make(chan struct{})
	s.cmd = cmd
	s.doneCh = doneCh
	s.active = true
	go func() {
		cmd.Wait() //nolint:errcheck
		fmt.Fprintln(os.Stderr, "writerdeck-server: session: editor process exited")
		s.end()
	}()
	if syncEng.ready() {
		go func() { _, _ = syncEng.reconcileAll("app") }()
	}
	return nil
}

// end marks the session inactive and restarts xochitl (unless sleeping for power).
func (s *session) end() {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	wasSleeping := s.sleeping
	s.active = false
	ch := s.doneCh
	s.cmd = nil
	s.doneCh = nil
	s.mu.Unlock()
	ungrabButtonDev()
	if wasSleeping {
		fmt.Fprintln(os.Stderr, "writerdeck-server: session: editor stopped for sleep (xochitl stays down)")
	} else {
		currentNoteMu.Lock()
		currentNote = ""
		currentNoteMu.Unlock()
		broadcast([]byte(`{"type":"exitedit"}`))
		fmt.Fprintln(os.Stderr, "writerdeck-server: session: starting xochitl")
		startXochitl()
	}
	if ch != nil {
		close(ch)
	}
}

// quit sends a graceful quit to the editor, waits for it to exit,
// and falls back to SIGTERM on the process group after 3 s.
func (s *session) quit() {
	s.mu.Lock()
	active := s.active
	doneCh := s.doneCh
	cmd := s.cmd
	s.mu.Unlock()
	if !active || doneCh == nil {
		return
	}
	fmt.Fprintln(os.Stderr, "writerdeck-server: session: sending quit to editor")
	if s.ec != nil && s.ec.ready() {
		flushEditorSave()
	}
	if s.ec != nil {
		s.ec.write([]byte(`{"t":"cmd","c":"quit"}`))
	}
	select {
	case <-doneCh:
		fmt.Fprintln(os.Stderr, "writerdeck-server: session: editor exited cleanly")
	case <-time.After(12 * time.Second):
		fmt.Fprintf(os.Stderr, "writerdeck-server: session: 12s timeout -- SIGTERM to process group")
		if cmd != nil && cmd.Process != nil {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM) //nolint:errcheck
		}
		<-doneCh
	}
}

const powerMenuAckTimeout = 120 * time.Second

// powerMenu shows the rM2 Sleep/Lobby/Exit overlay and acts on the choice.
// Sleep runs suspend after QML has already painted the sleep screen.
func (s *session) powerMenu() {
	if !s.isActive() || s.isSleeping() {
		return
	}
	powerMenuMu.Lock()
	if powerMenuWaiting {
		powerMenuMu.Unlock()
		return
	}
	powerMenuWaiting = true
	powerMenuMu.Unlock()
	defer func() {
		powerMenuMu.Lock()
		powerMenuWaiting = false
		powerMenuMu.Unlock()
	}()

	s.ec.drainPowerMenuChoice()
	s.ec.write([]byte(`{"t":"cmd","c":"powermenu"}`))
	choice := s.ec.waitPowerMenuChoice(powerMenuAckTimeout)
	fmt.Fprintf(os.Stderr, "writerdeck-server: power menu choice=%q\n", choice)
	switch choice {
	case "sleep":
		// Re-read at choice time: Lobby may have cleared the open note since the menu opened.
		currentNoteMu.Lock()
		noteToReopen := currentNote
		currentNoteMu.Unlock()
		s.suspendAfterSleepPaint(noteToReopen)
	case "sleeplobby":
		// Sleep started from Lobby — wake back to Lobby, not the last document.
		s.suspendAfterSleepPaint("")
	case "lobby":
		// Same as Home from edit: no note is open for sync / next wake.
		currentNoteMu.Lock()
		currentNote = ""
		currentNoteMu.Unlock()
		broadcast([]byte(`{"type":"exitedit","source":"lobby"}`))
	case "exit", "cancel", "":
		// QML handled Exit locally; cancel/timeout leaves UI as-is.
	default:
		fmt.Fprintf(os.Stderr, "writerdeck-server: unknown power menu choice %q\n", choice)
	}
}

// sleepForPower saves via QML, shows the sleep screen, stops keywriter (releases
// systemd-inhibit), and suspends. The e-ink frame persists until wake.
//
// On reMarkable 2, systemctl suspend often returns at Sleep target -- before the
// device has actually suspended -- so we wait for systemd-suspend.service to
// finish. Wake runs after that cycle; rM2 often resumes without delivering
// KEY_POWER to userspace.
func (s *session) sleepForPower() {
	if !s.isActive() || s.isSleeping() {
		return
	}
	currentNoteMu.Lock()
	noteToReopen := currentNote
	currentNoteMu.Unlock()

	if !s.ec.writeCmdWaitAck([]byte(`{"t":"cmd","c":"preparesleep"}`), "saved", "preparesleep", saveAckTimeout) {
		fmt.Fprintln(os.Stderr, "writerdeck-server: preparesleep save ack missed -- continuing")
	}
	if !s.ec.waitAck("ready", "preparesleep", paintAckTimeout) {
		fmt.Fprintln(os.Stderr, "writerdeck-server: sleep screen ready ack missed -- continuing")
	}
	s.suspendAfterSleepPaint(noteToReopen)
}

// suspendAfterSleepPaint finishes power-sleep after the sleep screen is on
// e-ink: sync, stop editor, suspend, wake and reopen noteToReopen.
func (s *session) suspendAfterSleepPaint(noteToReopen string) {
	currentNoteMu.Lock()
	currentNote = ""
	currentNoteMu.Unlock()
	beginSyncWait()
	go func() {
		if syncEng.ready() {
			syncEng.reconcileAllBlocking("power", syncAckTimeout)
		} else {
			signalSyncAck()
		}
	}()
	broadcast([]byte(`{"type":"exitedit","source":"power"}`))
	// rM2: leave USB-C as gadget so a charger or laptop works after sleep.
	restoreUSBGadget()
	waitSyncAck(syncAckTimeout)

	s.mu.Lock()
	s.sleeping = true
	s.sleepNote = noteToReopen
	cmd := s.cmd
	doneCh := s.doneCh
	s.mu.Unlock()

	if cmd != nil && cmd.Process != nil {
		fmt.Fprintln(os.Stderr, "writerdeck-server: stopping editor before suspend")
		syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM) //nolint:errcheck
		if doneCh != nil {
			select {
			case <-doneCh:
			case <-time.After(3 * time.Second):
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) //nolint:errcheck
				<-doneCh
			}
		}
	}
	suspendAndWait()
	fmt.Fprintln(os.Stderr, "writerdeck-server: resumed from suspend -- waking editor")
	if err := s.wakeFromSleep(""); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: wake after suspend failed: %v\n", err)
	}
}

// systemdUnitState returns systemctl show ActiveState for unit ("" on error).
func systemdUnitState(unit string) string {
	out, err := exec.Command("systemctl", "show", unit, "--property=ActiveState", "--value").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// suspendAndWait requests suspend and blocks until systemd-suspend has run and
// finished (device resumed). Needed because systemctl suspend on rM2 can return
// before the actual sleep/wake cycle completes.
func suspendAndWait() {
	fmt.Fprintln(os.Stderr, "writerdeck-server: suspending")
	before := time.Now()
	errCh := make(chan error, 1)
	go func() {
		errCh <- exec.Command("systemctl", "suspend").Run()
	}()

	deadline := time.Now().Add(3 * time.Minute)
	sawSuspend := false
	systemctlDone := false
	for time.Now().Before(deadline) {
		state := systemdUnitState("systemd-suspend.service")
		if state == "activating" || state == "active" {
			sawSuspend = true
		}
		if sawSuspend && state != "activating" && state != "active" {
			fmt.Fprintf(os.Stderr, "writerdeck-server: suspend cycle done after %s\n", time.Since(before).Round(time.Millisecond))
			time.Sleep(500 * time.Millisecond) // display / epaper settle
			return
		}
		select {
		case err := <-errCh:
			systemctlDone = true
			if err != nil {
				fmt.Fprintf(os.Stderr, "writerdeck-server: systemctl suspend: %v\n", err)
			}
		default:
		}
		// Suspend never started (failed request) -- do not spin forever.
		if systemctlDone && !sawSuspend && time.Since(before) > 5*time.Second {
			fmt.Fprintln(os.Stderr, "writerdeck-server: suspend did not start -- continuing")
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Fprintln(os.Stderr, "writerdeck-server: suspend wait timed out -- continuing wake")
}

// wakeFromSleep starts a fresh editor session and reopens the note that was open.
// noteName may be empty; then sleepNote from the session is used. sleeping stays
// true until the editor is up so a failed mid-wake start does not bring xochitl
// back over a blank screen.
func (s *session) wakeFromSleep(noteName string) error {
	s.mu.Lock()
	if !s.sleeping {
		s.mu.Unlock()
		return nil
	}
	if noteName == "" {
		noteName = s.sleepNote
	}
	s.mu.Unlock()

	if err := s.start(); err != nil {
		if s.isActive() {
			return nil // concurrent wake already started the editor
		}
		return err
	}
	for i := 0; i < 20; i++ {
		if s.ec.ready() {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if !s.ec.ready() {
		fmt.Fprintln(os.Stderr, "writerdeck-server: wake: socket not ready -- stopping editor for retry")
		s.quit()
		return fmt.Errorf("editor socket not ready after wake")
	}
	if noteName != "" {
		editorName := filepath.Base(noteName)
		fmt.Fprintf(os.Stderr, "writerdeck-server: wake: reopening %q\n", editorName)
		cmd, _ := json.Marshal(struct {
			T    string `json:"t"`
			C    string `json:"c"`
			Name string `json:"name"`
		}{"cmd", "open", editorName})
		// Restore currentNote before open so a second sleep still knows what to reopen.
		currentNoteMu.Lock()
		currentNote = editorName
		currentNoteMu.Unlock()
		broadcastOpenEdit(editorName)
		if !s.ec.writeCmdWaitAck(cmd, "saved", "open", saveAckTimeout) {
			fmt.Fprintln(os.Stderr, "writerdeck-server: wake open save ack missed -- continuing")
		}
	}

	s.mu.Lock()
	s.sleeping = false
	s.sleepNote = ""
	s.mu.Unlock()

	if syncEng.ready() {
		go func() { _, _ = syncEng.reconcileAll("wake") }()
	}
	return nil
}
