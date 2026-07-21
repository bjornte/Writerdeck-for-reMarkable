// Writerdeck-server — see main.go for overview.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		exec.Command("systemctl", "start", "xochitl").Run() //nolint:errcheck
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
		exec.Command("systemctl", "start", "xochitl").Run() //nolint:errcheck
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

// sleepForPower saves via QML, shows the sleep screen, stops keywriter (releases
// systemd-inhibit), and suspends. The e-ink frame persists until wake.
func (s *session) sleepForPower() {
	if !s.isActive() {
		return
	}
	if !s.ec.writeCmdWaitAck([]byte(`{"t":"cmd","c":"preparesleep"}`), "saved", "preparesleep", saveAckTimeout) {
		fmt.Fprintln(os.Stderr, "writerdeck-server: preparesleep save ack missed -- continuing")
	}
	if !s.ec.waitAck("ready", "preparesleep", paintAckTimeout) {
		fmt.Fprintln(os.Stderr, "writerdeck-server: sleep screen ready ack missed -- continuing")
	}
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
	waitSyncAck(syncAckTimeout)

	s.mu.Lock()
	s.sleeping = true
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
	fmt.Fprintln(os.Stderr, "writerdeck-server: suspending")
	exec.Command("systemctl", "suspend").Run() //nolint:errcheck
}

// wakeFromSleep starts a fresh editor session and reopens the note that was open.
func (s *session) wakeFromSleep(noteName string) error {
	s.mu.Lock()
	if !s.sleeping {
		s.mu.Unlock()
		return nil
	}
	s.sleeping = false
	s.mu.Unlock()

	if err := s.start(); err != nil {
		return err
	}
	if noteName != "" {
		for i := 0; i < 10; i++ {
			if s.ec.ready() {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		if !s.ec.ready() {
			return fmt.Errorf("editor socket not ready after wake")
		}
		editorName := filepath.Base(noteName)
		cmd, _ := json.Marshal(struct {
			T    string `json:"t"`
			C    string `json:"c"`
			Name string `json:"name"`
		}{"cmd", "open", editorName})
		s.ec.write(cmd)
	}
	if syncEng.ready() {
		go func() { _, _ = syncEng.reconcileAll("wake") }()
	}
	return nil
}
