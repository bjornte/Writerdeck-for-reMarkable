// Writerdeck-server -- see main.go for overview.

package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Physical buttons may span one or two evdev nodes:
//   rM1: gpio-keys on event1 (Home, Left, Right, Power, Wakeup)
//   rM2: snvs-powerkey on event0 (Power only); no Home/page gpio-keys
// Paths are resolved at open time by device name, not hardcoded eventN.
var (
	buttonDevMu      sync.Mutex
	buttonDevPaths   []string // gpio-keys and/or powerkey paths we opened
	buttonDevFiles   []*os.File
	buttonDevGrabbed bool
	// powerMenuDevice is true on rM2 (powerkey only, no gpio-keys Home).
	powerMenuDevice bool

	powerMenuMu      sync.Mutex
	powerMenuWaiting bool
)

// findInputByName returns /dev/input/event* whose sysfs name matches any of the
// substrings (case-insensitive). First match wins for each name; duplicates skipped.
func findInputByName(substrings ...string) string {
	entries, err := os.ReadDir("/sys/class/input")
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "event") {
			continue
		}
		b, err := os.ReadFile(filepath.Join("/sys/class/input", e.Name(), "device", "name"))
		if err != nil {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(string(b)))
		for _, sub := range substrings {
			if strings.Contains(name, strings.ToLower(sub)) {
				return "/dev/input/" + e.Name()
			}
		}
	}
	return ""
}

// resolveButtonDevices picks gpio-keys (rM1 buttons) and/or snvs-powerkey (rM2 power).
// Falls back to /dev/input/event1 for older images if nothing matches by name.
func resolveButtonDevices() []string {
	seen := map[string]bool{}
	var out []string
	add := func(path string) {
		if path == "" || seen[path] {
			return
		}
		seen[path] = true
		out = append(out, path)
	}
	add(findInputByName("gpio-keys", "gpio_keys"))
	add(findInputByName("snvs-powerkey", "snvs_powerkey", "powerkey"))
	if len(out) == 0 {
		add("/dev/input/event1")
	}
	return out
}

// openButtonDev opens gpio-keys and/or powerkey once. Safe before the watcher
// and before the first session.start().
func openButtonDev() error {
	buttonDevMu.Lock()
	defer buttonDevMu.Unlock()
	if len(buttonDevFiles) > 0 {
		return nil
	}
	paths := resolveButtonDevices()
	var files []*os.File
	var opened []string
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: open %s: %v\n", p, err)
			continue
		}
		files = append(files, f)
		opened = append(opened, p)
	}
	if len(files) == 0 {
		return fmt.Errorf("no button input devices (tried %v)", paths)
	}
	buttonDevFiles = files
	buttonDevPaths = opened
	// rM2: snvs-powerkey only. rM1: gpio-keys (may also list powerkey).
	hasGpio := findInputByName("gpio-keys", "gpio_keys") != ""
	hasPowerKey := findInputByName("snvs-powerkey", "snvs_powerkey", "powerkey") != ""
	powerMenuDevice = !hasGpio && hasPowerKey
	fmt.Fprintf(os.Stderr, "writerdeck-server: button devices: %s (powerMenu=%v)\n",
		strings.Join(opened, ", "), powerMenuDevice)
	return nil
}

// usePowerMenu reports whether short-press power should open Sleep/Lobby/Exit
// instead of sleeping immediately (rM2 power-only hardware).
func usePowerMenu() bool {
	buttonDevMu.Lock()
	defer buttonDevMu.Unlock()
	return powerMenuDevice
}

// grabButtonDev takes exclusive EVIOCGRAB on every open button device so Qt
// evdev cannot see Home/Power/page. Call before spawning Writerdeck.
func grabButtonDev() {
	buttonDevMu.Lock()
	defer buttonDevMu.Unlock()
	if len(buttonDevFiles) == 0 {
		fmt.Fprintln(os.Stderr, "writerdeck-server: button grab skipped (device not open)")
		return
	}
	if buttonDevGrabbed {
		return
	}
	for i, f := range buttonDevFiles {
		if err := evdevGrab(f.Fd()); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: EVIOCGRAB %s failed: %v\n", buttonDevPaths[i], err)
			continue
		}
		fmt.Fprintln(os.Stderr, "writerdeck-server: exclusive grab on "+buttonDevPaths[i])
	}
	buttonDevGrabbed = true
}

// ungrabButtonDev releases EVIOCGRAB so stock xochitl can read buttons again.
func ungrabButtonDev() {
	buttonDevMu.Lock()
	defer buttonDevMu.Unlock()
	if len(buttonDevFiles) == 0 || !buttonDevGrabbed {
		return
	}
	for i, f := range buttonDevFiles {
		if err := evdevUngrab(f.Fd()); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: EVIOCGRAB release %s failed: %v\n", buttonDevPaths[i], err)
			continue
		}
		fmt.Fprintln(os.Stderr, "writerdeck-server: released grab on "+buttonDevPaths[i])
	}
	buttonDevGrabbed = false
}

func isButtonDevPath(path string) bool {
	buttonDevMu.Lock()
	defer buttonDevMu.Unlock()
	for _, p := range buttonDevPaths {
		if p == path {
			return true
		}
	}
	return false
}

// handlePowerKey wakes after suspend, opens the rM2 power menu, or sleeps (rM1).
// On rM2, wake usually happens when systemctl suspend returns; KEY_POWER may not arrive.
func handlePowerKey(s *session) {
	if s == nil {
		return
	}
	if s.isSleeping() {
		s.mu.Lock()
		note := s.sleepNote
		s.mu.Unlock()
		fmt.Fprintln(os.Stderr, "writerdeck-server: power button -- waking from sleep")
		go func() { _ = s.wakeFromSleep(note) }()
		return
	}
	if !s.isActive() {
		return
	}
	if usePowerMenu() {
		powerMenuMu.Lock()
		waiting := powerMenuWaiting
		powerMenuMu.Unlock()
		if waiting {
			// Second press while menu is open: tell QML to choose Sleep.
			fmt.Fprintln(os.Stderr, "writerdeck-server: power button -- menu re-press (sleep)")
			s.ec.write([]byte(`{"t":"cmd","c":"powermenu"}`))
			return
		}
		fmt.Fprintln(os.Stderr, "writerdeck-server: power button -- power menu")
		go s.powerMenu()
		return
	}
	fmt.Fprintln(os.Stderr, "writerdeck-server: power button -- sleep")
	go s.sleepForPower()
}

// readButtonDevice reads one evdev node and dispatches Home / page / Power.
func readButtonDevice(f *os.File, path string, s *session, ec *editorConn) {
	fmt.Fprintln(os.Stderr, "writerdeck-server: watching physical buttons on "+path)
	var debounce time.Time
	var leftDown, rightDown bool
	var chordDebounce time.Time
	for {
		var ev inputEvent
		if err := binary.Read(f, binary.LittleEndian, &ev); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: button read %s: %v\n", path, err)
			return
		}
		if ev.Type != evKey {
			continue
		}

		// Page-button chord: hold left+right to launch Writerdeck from stock UI.
		if ev.Code == keyLeft || ev.Code == keyRight {
			if ev.Code == keyLeft {
				leftDown = ev.Value == 1
			} else {
				rightDown = ev.Value == 1
			}
			if leftDown && rightDown && ev.Value == 1 && time.Since(chordDebounce) >= keyboardDebounceMs {
				chordDebounce = time.Now()
				handleIdleLaunch(s, "page buttons (left+right)")
				continue
			}
			if s != nil && s.isActive() && ev.Value == 1 && !(leftDown && rightDown) {
				settingsMu.Lock()
				rot := curSettings.Rotation
				settingsMu.Unlock()
				cmd := physicalPageCmd(ev.Code == keyLeft, rot)
				payload := []byte(`{"t":"cmd","c":"` + cmd + `"}`)
				go ec.write(payload)
			}
			continue
		}

		if ev.Value != 1 {
			continue
		}
		if time.Since(debounce) < 800*time.Millisecond {
			continue
		}
		debounce = time.Now()

		if ev.Code == keyHome {
			if s != nil {
				if s.isActive() {
					fmt.Fprintln(os.Stderr, "writerdeck-server: home button -- relaying to editor")
					go func() {
						ec.writeCmdWaitAck([]byte(`{"t":"cmd","c":"home"}`), "saved", "home", saveAckTimeout)
						currentNoteMu.Lock()
						currentNote = ""
						currentNoteMu.Unlock()
						broadcast([]byte(`{"type":"exitedit","source":"home"}`))
						if syncEng.ready() {
							syncEng.reconcileAll("home")
						}
					}()
				} else {
					fmt.Fprintln(os.Stderr, "writerdeck-server: home button -- no active session, ignoring")
				}
			} else {
				fmt.Fprintln(os.Stderr, "writerdeck-server: home button pressed -- sending quit to editor")
				ec.write([]byte(`{"t":"cmd","c":"quit"}`))
				return
			}
			continue
		}

		// Only real power/wake codes -- never treat other leftover keys as power.
		if ev.Code == keyPower || ev.Code == keyWake {
			handlePowerKey(s)
		}
	}
}

// watchPhysicalButtons reads gpio-keys and/or snvs-powerkey.
// Supervisor (s != nil): Home relay + Power sleep/wake (session.sleepForPower).
// Standalone (s == nil): Home sends quit then returns.
func watchPhysicalButtons(s *session, ec *editorConn) {
	if err := openButtonDev(); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: button watcher: %v (OK on non-device machines)\n", err)
		return
	}
	buttonDevMu.Lock()
	files := append([]*os.File(nil), buttonDevFiles...)
	paths := append([]string(nil), buttonDevPaths...)
	buttonDevMu.Unlock()
	if s == nil {
		grabButtonDev()
		defer ungrabButtonDev()
	}
	if len(files) == 1 {
		readButtonDevice(files[0], paths[0], s, ec)
		return
	}
	var wg sync.WaitGroup
	for i := range files {
		wg.Add(1)
		go func(f *os.File, p string) {
			defer wg.Done()
			readButtonDevice(f, p, s, ec)
		}(files[i], paths[i])
	}
	wg.Wait()
}

// watchHomeButton is kept as an alias for callers that haven't been renamed yet.
func watchHomeButton(s *session, ec *editorConn) { watchPhysicalButtons(s, ec) }

// findKeyboardInputDevices returns /dev/input/event* nodes that look like USB
// keyboards (name contains "keyboard"), excluding gpio-keys / powerkey.
func findKeyboardInputDevices() []string {
	entries, err := os.ReadDir("/sys/class/input")
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "event") {
			continue
		}
		dev := "/dev/input/" + e.Name()
		if isButtonDevPath(dev) {
			continue
		}
		namePath := filepath.Join("/sys/class/input", e.Name(), "device", "name")
		b, err := os.ReadFile(namePath)
		if err != nil {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(string(b)))
		if !strings.Contains(name, "keyboard") {
			continue
		}
		out = append(out, dev)
	}
	return out
}

// handleIdleLaunch starts an editor session (Lobby) from the stock UI when no session
// is active and the device is not in power sleep. Used by USB Escape and the
// physical left+right page-button chord.
func handleIdleLaunch(s *session, source string) {
	if s == nil || s.isActive() || s.isSleeping() {
		return
	}
	fmt.Fprintf(os.Stderr, "writerdeck-server: %s -- launching editor to Lobby\n", source)
	go func() {
		if err := s.start(); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: %s launch failed: %v\n", source, err)
		}
	}()
}

// handleEscapeLaunch starts an editor session (Lobby) when Escape is pressed on a
// USB keyboard while the stock UI is up. Ignored during active sessions and
// power sleep -- keywriter handles Esc while editing; power button handles wake.
func handleEscapeLaunch(s *session) {
	handleIdleLaunch(s, "Escape")
}

func readUSBKeyboardEvents(dev string, s *session, debounce *struct {
	mu sync.Mutex
	t  time.Time
}) {
	f, err := os.Open(dev)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(os.Stderr, "writerdeck-server: watching USB keyboard %s for Escape (launch)\n", dev)
	for {
		var ev inputEvent
		if err := binary.Read(f, binary.LittleEndian, &ev); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: keyboard %s: %v\n", dev, err)
			return
		}
		if ev.Type != evKey || ev.Value != 1 || ev.Code != keyEsc {
			continue
		}
		debounce.mu.Lock()
		if time.Since(debounce.t) < keyboardDebounceMs {
			debounce.mu.Unlock()
			continue
		}
		debounce.t = time.Now()
		debounce.mu.Unlock()
		handleEscapeLaunch(s)
	}
}

// watchUSBKeyboardHotplug restarts the editor when a USB keyboard appears during
// an active session so the launcher can pin the device and apply the .qmap.
func watchUSBKeyboardHotplug(s *session) {
	known := make(map[string]struct{})
	for _, dev := range findKeyboardInputDevices() {
		known[dev] = struct{}{}
	}
	for {
		time.Sleep(keyboardRescan)
		for _, dev := range findKeyboardInputDevices() {
			if _, ok := known[dev]; ok {
				continue
			}
			known[dev] = struct{}{}
			if s != nil && s.isActive() {
				restartEditorForKeymap("USB keyboard " + dev + " connected")
			}
		}
	}
}

// watchUSBKeyboardForLaunch rescans for USB keyboards and listens for Escape to
// start Writerdeck when idle (stock UI). Does not intercept Esc while editing.
func watchUSBKeyboardForLaunch(s *session) {
	fmt.Fprintln(os.Stderr, "writerdeck-server: USB keyboard launch watcher started")
	var debounce struct {
		mu sync.Mutex
		t  time.Time
	}
	running := make(map[string]struct{})
	var mu sync.Mutex
	for {
		for _, dev := range findKeyboardInputDevices() {
			mu.Lock()
			if _, ok := running[dev]; ok {
				mu.Unlock()
				continue
			}
			running[dev] = struct{}{}
			mu.Unlock()
			go func(d string) {
				readUSBKeyboardEvents(d, s, &debounce)
				mu.Lock()
				delete(running, d)
				mu.Unlock()
			}(dev)
		}
		time.Sleep(keyboardRescan)
	}
}
