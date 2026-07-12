// Writerdeck-server — see main.go for overview.

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

// watchPhysicalButtons reads gpio-keys events (Home, Power, page buttons).
// Supervisor mode (s != nil): Home relay + Power sleep/wake (see session.sleepForPower).
// Standalone mode (s == nil): Home sends quit to ec then returns.
func watchPhysicalButtons(s *session, ec *editorConn) {
	f, err := os.Open(buttonDev)
	if err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: button watcher: %v (OK on non-device machines)\n", err)
		return
	}
	defer f.Close()
	fmt.Fprintln(os.Stderr, "writerdeck-server: watching physical buttons on "+buttonDev)
	var debounce time.Time
	var leftDown, rightDown bool
	var chordDebounce time.Time
	for {
		var ev inputEvent
		if err := binary.Read(f, binary.LittleEndian, &ev); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: button read error: %v\n", err)
			return
		}
		if ev.Type != evKey {
			continue
		}

		// Page-button chord: hold left+right together to launch Writerdeck from stock UI.
		if ev.Code == keyLeft || ev.Code == keyRight {
			if ev.Code == keyLeft {
				leftDown = ev.Value == 1
			} else {
				rightDown = ev.Value == 1
			}
			if leftDown && rightDown && ev.Value == 1 && time.Since(chordDebounce) >= keyboardDebounceMs {
				chordDebounce = time.Now()
				handleIdleLaunch(s, "page buttons (left+right)")
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

		// Power or Wakeup: sleep while editing, wake after suspend.
		if s == nil {
			continue
		}
		if s.isSleeping() {
			currentNoteMu.Lock()
			note := currentNote
			currentNoteMu.Unlock()
			fmt.Fprintln(os.Stderr, "writerdeck-server: power button -- waking from sleep")
			go func() { _ = s.wakeFromSleep(note) }()
			continue
		}
		if s.isActive() {
			fmt.Fprintln(os.Stderr, "writerdeck-server: power button -- sleep")
			go s.sleepForPower()
		}
	}
}

// watchHomeButton is kept as an alias for callers that haven't been renamed yet.
func watchHomeButton(s *session, ec *editorConn) { watchPhysicalButtons(s, ec) }

// findKeyboardInputDevices returns /dev/input/event* nodes that look like USB
// keyboards (name contains "keyboard"), excluding gpio-keys on event1.
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
		if dev == buttonDev {
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
