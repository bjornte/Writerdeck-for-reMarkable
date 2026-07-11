// Writerdeck-server -- reMarkable Wi-Fi typewriter daemon.
//
// Serves a WebSocket on 0.0.0.0:8000/ws and forwards received key events
// to the patched Writerdeck editor over a local Unix socket (/run/Writerdeck.sock).
//
// Architecture (two layers, two parsers):
//   Browser --WebSocket--> Writerdeck-server --Unix socket--> Writerdeck (patched keywriter)
//
// WebSocket message (JSON, from browser keydown):
//   {"type":"key","key":"<KeyboardEvent.key>"}
//
// Editor-feed wire format (NDJSON to keywriter's naive C++ parser):
//   {"t":"text","cp":<unicode-codepoint-int>}   -- single printable char
//   {"t":"key","k":"Escape|Return|Backspace|Tab|ArrowUp|ArrowDown|ArrowLeft|ArrowRight"}
//   {"t":"cmd","c":"home|open|notedeleted|noterenamed|..."}  -- editor commands
// Browser <- Writerdeck-server: {"type":"openedit","name":"<file>.md"} on tablet open/rename
//   {"type":"tabletcrud","op":"createnote|deletenote|renamenote","name":"...","oldName":"..."}
//     on tablet Lobby Files CRUD — phone pairs ghDelete/pushNote; queue drains on connect.
//   {"type":"diskchanged","name":"<file>.md"} when disk is written under an open editor (slice 8).
//
// Integer codepoints are escaping-proof: JSON special chars in typed text
// can never corrupt the naive C++ substring parser (see socket-inject.patch).
//
// Usage on the device:
//   /home/root/Writerdeck-server               # serve on :8000
//   /home/root/Writerdeck-server -v            # verbose key logging
//   /home/root/Writerdeck-server --selftest    # one-shot hello world+Return
//   /home/root/Writerdeck-server --port 9000   # custom port
package main

import (
	"bufio"
	_ "embed"
	crand "crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
)

// Static assets embedded at compile time; all served with Cache-Control: no-store.
// app.css, app.js, state.js, and sync.js are split out so index.html stays markup-only.
//
//go:embed index.html
var indexHTML []byte

//go:embed app.css
var appCSS []byte

//go:embed app.js
var appJS []byte

//go:embed state.js
var stateJS []byte

//go:embed sync.js
var syncJS []byte

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML) //nolint:errcheck
}

const (
	sockPath    = "/run/Writerdeck.sock"
	defaultPort = 8000
)

// inputEvent is the Linux input_event layout on ARM32 (16 bytes, little-endian).
// Used to read physical button events from /dev/input/event1 (gpio-keys node).
type inputEvent struct {
	Sec   uint32
	Usec  uint32
	Type  uint16
	Code  uint16
	Value int32
}

const (
	evKey     = 1   // EV_KEY
	keyEsc    = 1   // KEY_ESC -- USB keyboard launch from stock UI
	keyHome   = 102 // KEY_HOME -- middle (home) button, confirmed on /dev/input/event1
	keyLeft   = 105 // KEY_LEFT -- physical page-left button
	keyRight  = 106 // KEY_RIGHT -- physical page-right button
	keyPower  = 116 // KEY_POWER -- top power button
	keyWake   = 143 // KEY_WAKEUP -- fires on some wake paths
	buttonDev = "/dev/input/event1"

	saveAckTimeout     = 10 * time.Second // wait for keywriter {"t":"saved",...}
	paintAckTimeout    = 3 * time.Second  // e-ink sleep screen {"t":"ready",...}
	syncAckTimeout     = 45 * time.Second // power sleep waits for GitHub reconcile
	keyboardRescan     = 3 * time.Second  // hotplug rescan for USB keyboards
	keyboardDebounceMs = 800 * time.Millisecond
)

// wsMsg is the JSON message received from the browser on keydown.
type wsMsg struct {
	Type  string `json:"type"`  // always "key"
	Key   string `json:"key"`   // KeyboardEvent.key value
	Shift bool   `json:"shift"` // event.shiftKey
	Ctrl  bool   `json:"ctrl"`  // event.ctrlKey
	Alt   bool   `json:"alt"`   // event.altKey
	Meta  bool   `json:"meta"`  // event.metaKey (Cmd on Mac/iPhone)
}

// namedKeys maps browser KeyboardEvent.key values to keywriter named keys.
// Only keys whose value is NOT a single printable rune need an entry here.
var namedKeys = map[string]string{
	"Enter":      "Return",
	"Backspace":  "Backspace",
	"Tab":        "Tab",
	"Escape":     "Escape",
	"Home":       "Home",
	"End":        "End",
	"ArrowUp":    "ArrowUp",
	"ArrowDown":  "ArrowDown",
	"ArrowLeft":  "ArrowLeft",
	"ArrowRight": "ArrowRight",
}

// editorConn holds the live connection to keywriter's socket.
// rmkbd dials and redials; keywriter is the server.
// keywriter writes back {"t":"saved"|"ready","c":"<cmd>"} ack lines.
type editorConn struct {
	mu   sync.Mutex
	conn net.Conn

	ackMu   sync.Mutex
	ackWait []*ackWait
}

type ackWait struct {
	typ string
	cmd string
	ch  chan struct{}
}

func (e *editorConn) registerAckWait(typ, cmd string) chan struct{} {
	ch := make(chan struct{}, 1)
	e.ackMu.Lock()
	e.ackWait = append(e.ackWait, &ackWait{typ: typ, cmd: cmd, ch: ch})
	e.ackMu.Unlock()
	return ch
}

func (e *editorConn) cancelAckWait(typ, cmd string) {
	e.ackMu.Lock()
	defer e.ackMu.Unlock()
	for i, w := range e.ackWait {
		if w.typ == typ && w.cmd == cmd {
			e.ackWait = append(e.ackWait[:i], e.ackWait[i+1:]...)
			return
		}
	}
}

func (e *editorConn) signalAck(typ, cmd string) {
	e.ackMu.Lock()
	defer e.ackMu.Unlock()
	for i, w := range e.ackWait {
		if w.typ == typ && w.cmd == cmd {
			select {
			case w.ch <- struct{}{}:
			default:
			}
			e.ackWait = append(e.ackWait[:i], e.ackWait[i+1:]...)
			return
		}
	}
}

func (e *editorConn) clearAckWaits() {
	e.ackMu.Lock()
	defer e.ackMu.Unlock()
	e.ackWait = nil
}

func (e *editorConn) waitAck(typ, cmd string, timeout time.Duration) bool {
	ch := e.registerAckWait(typ, cmd)
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		e.cancelAckWait(typ, cmd)
		return false
	}
}

func (e *editorConn) writeCmdWaitAck(cmd []byte, typ, cmdName string, timeout time.Duration) bool {
	ch := e.registerAckWait(typ, cmdName)
	e.write(cmd)
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		e.cancelAckWait(typ, cmdName)
		fmt.Fprintf(os.Stderr, "writerdeck-server: ack timeout %s/%s\n", typ, cmdName)
		return false
	}
}

func (e *editorConn) handleEditorLine(line []byte) {
	var msg struct {
		T       string `json:"t"`
		C       string `json:"c"`
		Op      string `json:"op"`
		Name    string `json:"name"`
		Old     string `json:"old"`
		Degrees int    `json:"degrees"`
	}
	if json.Unmarshal(line, &msg) != nil {
		return
	}
	switch msg.T {
	case "req":
		handleEditorReq(msg.Op, msg.Name, msg.Old)
	case "open":
		name := filepath.Base(msg.Name)
		if notesSafe(name) == "" {
			return
		}
		currentNoteMu.Lock()
		currentNote = name
		currentNoteMu.Unlock()
		broadcastOpenEdit(name)
	case "saved", "ready":
		e.signalAck(msg.T, msg.C)
	case "rotation":
		settingsMu.Lock()
		curSettings.Rotation = normalizeRotation(msg.Degrees)
		saveSettingsLocked()
		settingsMu.Unlock()
	}
}

func (e *editorConn) write(line []byte) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.conn == nil {
		return
	}
	_, err := e.conn.Write(append(line, '\n'))
	if err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: editor socket write error: %v -- will redial\n", err)
		e.conn.Close()
		e.conn = nil
	}
}

func (e *editorConn) set(c net.Conn) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.conn != nil {
		e.conn.Close()
	}
	e.conn = c
}

func (e *editorConn) ready() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.conn != nil
}

// ipv4OnInterface returns the first IPv4 on name, or "" if the interface is
// down or has no suitable address.
func ipv4OnInterface(name string) string {
	iface, err := net.InterfaceByName(name)
	if err != nil || iface.Flags&net.FlagUp == 0 {
		return ""
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			continue
		}
		if ip4 := ip.To4(); ip4 != nil {
			return ip4.String()
		}
	}
	return ""
}

// getLocalIP returns the device's Wi-Fi IPv4, preferring wlan0, or "?" if none
// is up yet. Used to populate the tablet Lobby screen.
func getLocalIP() string {
	if ip := ipv4OnInterface("wlan0"); ip != "" {
		return ip
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return "?"
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		if ip := ipv4OnInterface(iface.Name); ip != "" {
			return ip
		}
	}
	return "?"
}

// dialLoop keeps a live connection to keywriter's socket, redialling on loss.
func dialLoop(ec *editorConn) {
	logged := false
	for {
		c, err := net.Dial("unix", sockPath)
		if err != nil {
			if !logged {
				fmt.Fprintf(os.Stderr, "writerdeck-server: waiting for editor socket (retrying silently until connected)...\n")
				logged = true
			}
			time.Sleep(time.Second)
			continue
		}
		logged = false
		fmt.Fprintln(os.Stderr, "writerdeck-server: connected to editor socket")
		ec.set(c)
		// Push lobby info so keywriter displays the current IP and PIN (or the
		// no-PIN text when pin is ""). Always send even in no-PIN mode so the
		// QML conditional can render the appropriate Lobby line.
		pushLobbyInfo()
		pushNotesList()
		// Push persisted font so a freshly-spawned editor reflects the saved choice.
		settingsMu.Lock()
		fontFamily := curSettings.ReadFont
		settingsMu.Unlock()
		if fontFamily != "" {
			fontMsg, _ := json.Marshal(struct {
				T      string `json:"t"`
				C      string `json:"c"`
				Family string `json:"family"`
			}{"cmd", "setfont", fontFamily})
			ec.write(fontMsg)
		}
		settingsMu.Lock()
		rotation := curSettings.Rotation
		settingsMu.Unlock()
		rotMsg, _ := json.Marshal(struct {
			T       string `json:"t"`
			C       string `json:"c"`
			Degrees int    `json:"degrees"`
		}{"cmd", "setrotation", rotation})
		ec.write(rotMsg)
		// Read ack lines until the connection dies.
		sc := bufio.NewScanner(c)
		for sc.Scan() {
			ec.handleEditorLine(sc.Bytes())
		}
		fmt.Fprintln(os.Stderr, "writerdeck-server: editor socket closed -- redialling")
		ec.clearAckWaits()
		ec.set(nil)
	}
}

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

// translate converts a browser key event to an editor-feed NDJSON line.
// Returns nil if the key should be ignored (e.g. lone modifier keys).
//
// Three-way classification:
//   Action -- Ctrl or Meta(Cmd) held: forward key name + modifier mask.
//             Letters are uppercased so C++ maps A-Z to Qt::Key_A..Z.
//   Named  -- arrow/Home/End/Enter/etc.: forward name + mask.
//             Shift is forwarded so Shift+Arrow selects text in Qt TextEdit.
//   Text   -- single printable rune, no Ctrl/Meta: codepoint only.
//             event.key already has Shift/Alt baked in; adding a modifier
//             would double-apply (e.g. 'A' -> Shift+'a' -> garbled).
func translate(ev wsMsg) []byte {
	key := ev.Key
	isAction := ev.Ctrl || ev.Meta

	// Modifier bitmask for the editor-feed (Shift=1, Ctrl=2, Alt=4, Meta=8).
	// C++ injector maps both Ctrl(2) and Meta(8) to Qt::ControlModifier.
	mask := 0
	if ev.Shift {
		mask |= 1
	}
	if ev.Ctrl {
		mask |= 2
	}
	if ev.Alt {
		mask |= 4
	}
	if ev.Meta {
		mask |= 8
	}

	// Named control key (arrow, Enter, Backspace, Tab, Escape, Home, End)?
	if kwKey, ok := namedKeys[key]; ok {
		if mask != 0 {
			return []byte(fmt.Sprintf(`{"t":"key","k":%q,"m":%d}`, kwKey, mask))
		}
		return []byte(fmt.Sprintf(`{"t":"key","k":%q}`, kwKey))
	}

	// Action: Ctrl or Meta(Cmd) held on a printable key.
	// Uppercase so C++ maps the single letter to Qt::Key_A + (c - 'A').
	if isAction && utf8.RuneCountInString(key) == 1 {
		r, _ := utf8.DecodeRuneInString(key)
		if r != utf8.RuneError {
			return []byte(fmt.Sprintf(`{"t":"key","k":%q,"m":%d}`, strings.ToUpper(key), mask))
		}
	}

	// Text: single printable codepoint, no Ctrl/Meta.
	// Modifiers already baked into the rune by the browser.
	if !isAction && utf8.RuneCountInString(key) == 1 {
		r, _ := utf8.DecodeRuneInString(key)
		if r != utf8.RuneError {
			return []byte(fmt.Sprintf(`{"t":"text","cp":%d}`, r))
		}
	}

	// Ignore everything else (modifier keys, dead keys, multi-char specials, etc.)
	return nil
}

var upgrader = websocket.Upgrader{
	// LAN use; no auth in Phase 3.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// logEvery controls how often the terse (non-verbose) log prints a running
// key count. Per-key translation detail is gated behind -v (for keymap
// debugging); by default the log stays quiet -- a periodic count plus a
// per-session total is enough to confirm keys are flowing, without flooding
// the device log with one line per keystroke.
const logEvery = 25

// --- WS broadcast hub ---
//
// Every connected browser is registered as a wsClient. A dedicated writer
// goroutine per client owns all conn.Write calls (gorilla WS forbids
// concurrent writers). broadcast() fans out a server-push message to all
// clients; sends are non-blocking so a slow/dead client cannot stall the
// caller.
type wsClient struct {
	conn *websocket.Conn
	send chan []byte
}

var (
	wsClientsMu sync.Mutex
	wsClients   = make(map[*wsClient]bool)

	syncAckMu sync.Mutex
	syncAckCh chan struct{} // set while power-sleep waits for browser reconcile
)

// signalSyncAck unblocks sleepForPower after the phone browser finishes GitHub sync.
func signalSyncAck() {
	syncAckMu.Lock()
	ch := syncAckCh
	syncAckCh = nil
	syncAckMu.Unlock()
	if ch != nil {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func beginSyncWait() {
	syncAckMu.Lock()
	syncAckCh = make(chan struct{}, 1)
	syncAckMu.Unlock()
}

func waitSyncAck(timeout time.Duration) {
	syncAckMu.Lock()
	ch := syncAckCh
	syncAckMu.Unlock()
	if ch == nil {
		return
	}
	select {
	case <-ch:
		fmt.Fprintln(os.Stderr, "writerdeck-server: sync ack received")
	case <-time.After(timeout):
		fmt.Fprintln(os.Stderr, "writerdeck-server: sync ack timeout -- proceeding")
	}
	syncAckMu.Lock()
	if syncAckCh == ch {
		syncAckCh = nil
	}
	syncAckMu.Unlock()
}

// broadcastOpenEdit tells phone clients which note the tablet editor holds open.
func broadcastOpenEdit(name string) {
	if name == "" {
		return
	}
	msg, _ := json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{"openedit", name})
	broadcast(msg)
}

// broadcast pushes msg to every registered browser client.
func broadcast(msg []byte) {
	wsClientsMu.Lock()
	defer wsClientsMu.Unlock()
	for c := range wsClients {
		select {
		case c.send <- msg:
		default:
		}
	}
}

// enqueuePendingSync records a tablet file op for the phone to pair on GitHub.
func enqueuePendingSync(op, name, oldName string) {
	settingsMu.Lock()
	curSettings.PendingSync = append(curSettings.PendingSync, pendingSyncOp{Op: op, Name: name, OldName: oldName})
	saveSettingsLocked()
	settingsMu.Unlock()
}

// clearPendingSync removes all queued tablet sync ops (after phone has paired them).
func clearPendingSync() {
	settingsMu.Lock()
	if len(curSettings.PendingSync) == 0 {
		settingsMu.Unlock()
		return
	}
	curSettings.PendingSync = nil
	saveSettingsLocked()
	settingsMu.Unlock()
	pushLobbyInfo()
}

// notifyTabletCrud queues the op and tells connected phone browsers to pair it on GitHub.
func notifyTabletCrud(op, name, oldName string) {
	if name != "" {
		if p := notesSafe(name); p != "" {
			name = filepath.Base(p)
		}
	}
	if oldName != "" {
		if p := notesSafe(oldName); p != "" {
			oldName = filepath.Base(p)
		}
	}
	enqueuePendingSync(op, name, oldName)
	msg, _ := json.Marshal(struct {
		Type    string `json:"type"`
		Op      string `json:"op"`
		Name    string `json:"name"`
		OldName string `json:"oldName,omitempty"`
	}{"tabletcrud", op, name, oldName})
	broadcast(msg)
	pushLobbyInfo()
}

// maybeBroadcastDiskChanged notifies phone browsers when disk was written for the open note.
func maybeBroadcastDiskChanged(name string) {
	base := filepath.Base(name)
	currentNoteMu.Lock()
	open := currentNote != "" && currentNote == base
	currentNoteMu.Unlock()
	if !open {
		return
	}
	msg, _ := json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{"diskchanged", base})
	broadcast(msg)
}

// currentNote is the basename (.md) of the note the editor currently has open.
// Protected by currentNoteMu. Set by openHandler and editor {"t":"open"} reports;
// cleared by watchHomeButton, session.end(), and the DELETE handler on a match.
var (
	currentNoteMu sync.Mutex
	currentNote   string
)

func wsHandler(ec *editorConn, verbose bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authMu.Lock()
		required := pinRequired
		tok := authToken
		authMu.Unlock()
		if required {
			cookie, err := r.Cookie("writerdeck_token")
			if err != nil || cookie.Value != tok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: WS upgrade error: %v\n", err)
			return
		}
		// Register in the broadcast hub; writer goroutine owns all conn writes.
		client := &wsClient{conn: conn, send: make(chan []byte, 8)}
		wsClientsMu.Lock()
		wsClients[client] = true
		wsClientsMu.Unlock()
		defer func() {
			wsClientsMu.Lock()
			delete(wsClients, client)
			wsClientsMu.Unlock()
			close(client.send) // signals writer goroutine to drain and exit
			conn.Close()
		}()
		go func() {
			for msg := range client.send {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					conn.Close() // read loop will see error and exit
					return
				}
			}
		}()
		// Tell new clients which note the tablet editor holds open (edit lease).
		currentNoteMu.Lock()
		openNote := currentNote
		currentNoteMu.Unlock()
		if openNote != "" {
			msg, _ := json.Marshal(struct {
				Type string `json:"type"`
				Name string `json:"name"`
			}{"openedit", openNote})
			select {
			case client.send <- msg:
			default:
			}
		}
		remote := r.RemoteAddr
		fmt.Fprintf(os.Stderr, "writerdeck-server: client connected %s\n", remote)
		var keys int
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Fprintf(os.Stderr, "writerdeck-server: client disconnected %s: %v (%d keys forwarded)\n", remote, err, keys)
				return
			}
			var ev wsMsg
			if err := json.Unmarshal(msg, &ev); err != nil || ev.Type != "key" {
				continue
			}
			line := translate(ev)
			if line == nil {
				continue
			}
			keys++
			if verbose {
				fmt.Fprintf(os.Stderr, "writerdeck-server: key=%q -> %s\n", ev.Key, line)
			} else if keys%logEvery == 0 {
				fmt.Fprintf(os.Stderr, "writerdeck-server: forwarded %d keys\n", keys)
			}
			ec.write(line)
		}
	}
}

// --- Notes API ---

// notesDirPath is where .md notes are stored.
// Override with --notes-dir for local testing (default: /home/root/Writerdeck-user-documents).
var notesDirPath = "/home/root/Writerdeck-user-documents"

// noteInfo is the JSON shape returned by GET /api/notes.
type noteInfo struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

// notesSafe validates a filename and returns its full path, or "".
// Rejects empty names, slashes, "..", and appends ".md" if absent.
func notesSafe(name string) string {
	if name == "" || strings.Contains(name, "/") || strings.Contains(name, "..") {
		return ""
	}
	if !strings.HasSuffix(name, ".md") {
		name += ".md"
	}
	return filepath.Join(notesDirPath, name)
}

// rejectsHtmlNoteContent reports Qt qrichtext / HTML accidentally saved as .md.
func rejectsHtmlNoteContent(content string) bool {
	if len(content) < 15 {
		return false
	}
	head := strings.ToLower(content)
	if len(head) > 512 {
		head = head[:512]
	}
	return strings.HasPrefix(head, "<!doctype html") ||
		strings.HasPrefix(head, "<html") ||
		strings.Contains(content, `name="qrichtext"`)
}

// noteETag is a content-hash revision token for optimistic concurrency (RFC 7232).
func noteETag(content []byte) string {
	sum := sha256.Sum256(content)
	return `"` + hex.EncodeToString(sum[:8]) + `"`
}

// ifMatchOK reports whether the If-Match header allows writing over etag.
func ifMatchOK(ifMatch, etag string) bool {
	if ifMatch == "" {
		return false
	}
	if ifMatch == "*" {
		return true
	}
	for _, part := range strings.Split(ifMatch, ",") {
		part = strings.TrimSpace(part)
		if part == etag || part == "*" {
			return true
		}
	}
	return false
}

// writeNoteFile atomically writes note bytes (write-temp-then-rename, like settings.json).
func writeNoteFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, content, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// notesListHandler serves GET /api/notes (list) and POST /api/notes (create).
func notesListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if r.URL.Path != "/api/notes" {
		http.NotFound(w, r)
		return
	}
	if !checkAuth(w, r) {
		return
	}
	switch r.Method {
	case http.MethodGet:
		entries, err := os.ReadDir(notesDirPath)
		if err != nil {
			if os.IsNotExist(err) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("[]\n")) //nolint:errcheck
				return
			}
			http.Error(w, "cannot read notes dir", http.StatusInternalServerError)
			return
		}
		var notes []noteInfo
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			info, err := e.Info()
			if err != nil {
				continue
			}
			notes = append(notes, noteInfo{
				Name:     e.Name(),
				Size:     info.Size(),
				Modified: info.ModTime().Format(time.RFC3339),
			})
		}
		if notes == nil {
			notes = []noteInfo{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes) //nolint:errcheck

	case http.MethodPost:
		// Cap the request body. The client also limits uploads to 1 MB, but a
		// client-side check is bypassable (e.g. a direct curl by an authed
		// caller), so this is the authoritative limit for ALL create paths
		// (New / New w/ paste / Upload). 2 MiB leaves headroom over a 1 MB text
		// file once it is wrapped + escaped in the JSON {name, content} envelope.
		// When exceeded, the read fails and the decode below returns 400.
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
		var req struct {
			Name    string `json:"name"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			http.Error(w, "bad request: need {name}", http.StatusBadRequest)
			return
		}
		p := notesSafe(req.Name)
		if p == "" {
			http.Error(w, "invalid name", http.StatusBadRequest)
			return
		}
		if _, err := os.Stat(p); err == nil {
			http.Error(w, "already exists", http.StatusConflict)
			return
		}
		content := req.Content
		if content == "" {
			content = "# " + strings.TrimSuffix(req.Name, ".md") + "\n"
		}
		if rejectsHtmlNoteContent(content) {
			http.Error(w, "refusing HTML/qrichtext payload", http.StatusUnsupportedMediaType)
			return
		}
		if err := writeNoteFile(p, []byte(content)); err != nil {
			http.Error(w, "write failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		pushLobbyInfo()
		pushNotesList()

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// notesItemHandler serves GET /api/notes/{name} (read),
// DELETE /api/notes/{name} (delete), and PATCH /api/notes/{name} (rename).
func notesItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, If-Match")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/api/notes/")
	// Download variant: GET /api/notes/{name}/download
	// Strip the suffix BEFORE notesSafe() -- notesSafe rejects names containing '/'.
	download := strings.HasSuffix(name, "/download")
	if download {
		name = strings.TrimSuffix(name, "/download")
	}
	p := notesSafe(name)
	if p == "" {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		data, err := os.ReadFile(p)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if download {
			base := filepath.Base(p)
			w.Header().Set("Content-Disposition", `attachment; filename="`+base+`"`)
			w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		} else {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}
		w.Header().Set("ETag", noteETag(data))
		w.Write(data) //nolint:errcheck

	case http.MethodDelete:
		if err := deleteNoteFile(filepath.Base(p)); err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, "delete failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodPut:
		// Upsert: write or overwrite content. Used by the sync engine to apply a
		// version pulled from GitHub. Overwrite requires If-Match (content ETag).
		// 2 MiB limit matches POST /api/notes.
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
		var putReq struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&putReq); err != nil {
			http.Error(w, "bad request: need {content}", http.StatusBadRequest)
			return
		}
		if rejectsHtmlNoteContent(putReq.Content) {
			http.Error(w, "refusing HTML/qrichtext payload", http.StatusUnsupportedMediaType)
			return
		}
		existing, err := os.ReadFile(p)
		editorLocal := isLoopback(r)
		if err == nil {
			if !editorLocal {
				etag := noteETag(existing)
				ifMatch := r.Header.Get("If-Match")
				if !ifMatchOK(ifMatch, etag) {
					w.Header().Set("ETag", etag)
					http.Error(w, "If-Match required or revision mismatch", http.StatusPreconditionFailed)
					return
				}
			}
		} else if !os.IsNotExist(err) {
			http.Error(w, "read failed", http.StatusInternalServerError)
			return
		}
		if err := writeNoteFile(p, []byte(putReq.Content)); err != nil {
			http.Error(w, "write failed", http.StatusInternalServerError)
			return
		}
		if !editorLocal {
			maybeBroadcastDiskChanged(filepath.Base(p))
		}
		w.Header().Set("ETag", noteETag([]byte(putReq.Content)))
		w.WriteHeader(http.StatusOK)

	case http.MethodPatch:
		// Rename: body {"name":"new-name.md"}
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			http.Error(w, "bad request: need {name}", http.StatusBadRequest)
			return
		}
		if err := renameNoteFile(filepath.Base(p), req.Name); err != nil {
			if strings.Contains(err.Error(), "already taken") {
				http.Error(w, "name already taken", http.StatusConflict)
				return
			}
			if strings.Contains(err.Error(), "invalid") {
				http.Error(w, "invalid name", http.StatusBadRequest)
				return
			}
			http.Error(w, "rename failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// --- Settings API ---

// settingsFilePath is the JSON settings store on the device.
// Override with --settings-file for local dev (mirrors --notes-dir).
var settingsFilePath = "/home/root/.Writerdeck/settings.json"

// settingsData is the on-disk and in-memory settings schema.
type settingsData struct {
	ReadFont  string `json:"readFont"`
	PinDigits string `json:"pinDigits"` // "6", "4", or "none"; default "6"
	Rotation  int    `json:"rotation"`  // display rotation in degrees (0, 90, 180, 270)
	SyncOn     bool  `json:"syncOn"`               // GitHub two-way sync enabled
	SyncRepo   string `json:"syncRepo"`          // "owner/repo" of the notes repo; token never stored here
	LastSyncAt int64  `json:"lastSyncAt,omitempty"` // unix seconds of last browser reconcile ack
	PendingSync []pendingSyncOp `json:"pendingSync,omitempty"` // tablet CRUD awaiting phone GitHub pairing
}

// pendingSyncOp is one queued tablet file op for the phone browser to mirror on GitHub.
type pendingSyncOp struct {
	Op      string `json:"op"`                // createnote, deletenote, renamenote
	Name    string `json:"name"`              // target .md basename
	OldName string `json:"oldName,omitempty"` // renamenote source basename
}

// normalizeRotation maps any integer to a 0-359 degree value.
func normalizeRotation(deg int) int {
	deg %= 360
	if deg < 0 {
		deg += 360
	}
	return deg
}

// isValidGitHubRepo returns true iff repo is a non-empty "owner/repo" string
// where both parts contain only characters valid in GitHub owner/repo names.
func isValidGitHubRepo(repo string) bool {
	parts := strings.SplitN(repo, "/", 3)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	valid := func(s string) bool {
		for _, c := range s {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
				(c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.') {
				return false
			}
		}
		return true
	}
	return valid(parts[0]) && valid(parts[1])
}

// fontOption is one entry in the reading-view font picker.
type fontOption struct {
	ID    string `json:"id"`    // exact Qt internal family name (must match TTF)
	Label string `json:"label"` // human-readable label shown in the phone UI
}

// fontRegistry is the canonical allow-list. IDs must exactly match the Qt
// internal family names as reported by fc-query; a wrong name silently falls
// back to DejaVu on the device with no error from Qt.
var fontRegistry = []fontOption{
	{ID: "Inter", Label: "Inter"},
	{ID: "Literata", Label: "Literata"},
	{ID: "EB Garamond", Label: "EB Garamond"},
	{ID: "DejaVu Sans", Label: "DejaVu Sans"},
}

var (
	settingsMu  sync.Mutex
	curSettings = settingsData{ReadFont: "Inter", PinDigits: "6"}
)

// globalEC is set in main() so settingsHandler can push font changes to the
// editor without requiring supervisor mode. Always the same *editorConn as
// activeSess.ec in supervisor mode.
var globalEC *editorConn

// loadSettings reads the persisted settings file and populates curSettings.
// Missing file is silently ignored (first run). Invalid JSON uses the default.
func loadSettings() {
	settingsMu.Lock()
	defer settingsMu.Unlock()
	data, err := os.ReadFile(settingsFilePath)
	if err != nil {
		return // first run or unreadable; keep default
	}
	var s settingsData
	if json.Unmarshal(data, &s) == nil {
		if s.ReadFont == "" {
			s.ReadFont = "Inter" // upgrade: missing field keeps default
		}
		if s.PinDigits == "" {
			s.PinDigits = "6" // upgrade: existing font-only settings.json had no pin field
		}
		curSettings = s
	}
}

// saveSettingsLocked writes curSettings to disk atomically.
// Caller must hold settingsMu.
func saveSettingsLocked() {
	dir := filepath.Dir(settingsFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: settings mkdir %s: %v\n", dir, err)
		return
	}
	data, err := json.Marshal(curSettings)
	if err != nil {
		return
	}
	tmp := settingsFilePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: settings write: %v\n", err)
		return
	}
	os.Rename(tmp, settingsFilePath) //nolint:errcheck
}

// settingsHandler serves GET /api/settings (read) and POST /api/settings (write).
func settingsHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		settingsMu.Lock()
		resp := struct {
			ReadFont  string       `json:"readFont"`
			Fonts     []fontOption `json:"fonts"`
			PinDigits string       `json:"pinDigits"`
			PinOpts   []string     `json:"pinOpts"`
			SyncOn    bool         `json:"syncOn"`
			SyncRepo  string       `json:"syncRepo"`
		}{curSettings.ReadFont, fontRegistry, curSettings.PinDigits, []string{"6", "4", "none"},
			curSettings.SyncOn, curSettings.SyncRepo}
		settingsMu.Unlock()
		json.NewEncoder(w).Encode(resp) //nolint:errcheck

	case http.MethodPost:
		var req struct {
			ReadFont  string  `json:"readFont"`
			PinDigits string  `json:"pinDigits"`
			SyncOn    *bool   `json:"syncOn"`
			SyncRepo  *string `json:"syncRepo"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if req.ReadFont != "" {
			// Validate against registry -- prevents arbitrary family injection.
			valid := false
			for _, f := range fontRegistry {
				if f.ID == req.ReadFont {
					valid = true
					break
				}
			}
			if !valid {
				http.Error(w, "unknown font", http.StatusBadRequest)
				return
			}
			settingsMu.Lock()
			curSettings.ReadFont = req.ReadFont
			saveSettingsLocked()
			font := curSettings.ReadFont
			settingsMu.Unlock()
			// Push to editor if a connection is alive.
			if globalEC != nil {
				cmd, _ := json.Marshal(struct {
					T      string `json:"t"`
					C      string `json:"c"`
					Family string `json:"family"`
				}{"cmd", "setfont", font})
				globalEC.write(cmd)
			}
		}
		if req.PinDigits != "" {
			// Validate the enum; 400 on an unknown value.
			pinLen := 0
			switch req.PinDigits {
			case "6":
				pinLen = 6
			case "4":
				pinLen = 4
			case "none":
				pinLen = 0
			default:
				http.Error(w, `pinDigits must be "6", "4", or "none"`, http.StatusBadRequest)
				return
			}
			newPIN := generatePIN(pinLen)
			newToken := generateToken()
			authMu.Lock()
			authPIN = newPIN
			authToken = newToken
			pinRequired = pinLen > 0
			authMu.Unlock()
			// Clear stale lockouts so the fresh PIN starts clean.
			pinMu.Lock()
			pinAttempts = map[string]*pinAttempt{}
			pinMu.Unlock()
			// Persist.
			settingsMu.Lock()
			curSettings.PinDigits = req.PinDigits
			saveSettingsLocked()
			settingsMu.Unlock()
			// Push new Lobby info to editor so the tablet shows the updated PIN at once.
			// (No-PIN mode sends pin="" so the QML conditional renders the friendly line.)
			pushLobbyInfo()
			pushNotesList()
			// Issue a fresh cookie so the changer stays authed after the change.
			// Without this, switching from no-PIN to 6-digit would instantly 401 the changer.
			exp := nextMorningCutoff(time.Now())
			http.SetCookie(w, &http.Cookie{
				Name:     "writerdeck_token",
				Value:    newToken,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				Expires:  exp,
				MaxAge:   int(time.Until(exp).Seconds()),
			})
		}
		if req.SyncRepo != nil {
			repo := *req.SyncRepo
			if repo != "" && !isValidGitHubRepo(repo) {
				http.Error(w, `syncRepo must be "owner/repo"`, http.StatusBadRequest)
				return
			}
			settingsMu.Lock()
			curSettings.SyncRepo = repo
			saveSettingsLocked()
			settingsMu.Unlock()
			pushLobbyInfo()
			pushNotesList()
		}
		if req.SyncOn != nil {
			settingsMu.Lock()
			curSettings.SyncOn = *req.SyncOn
			saveSettingsLocked()
			settingsMu.Unlock()
			pushLobbyInfo()
			pushNotesList()
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// --- Auth ---

var (
	authMu      sync.Mutex // guards authPIN, authToken, pinRequired
	authPIN     string     // PIN generated at startup; shown on the tablet Lobby
	authToken   string     // session token issued when the PIN is verified
	pinRequired bool       // false in no-PIN mode (checkAuth always passes)
	activeSess  *session   // non-nil only in supervisor (--editor) mode
)

// generatePIN mints a cryptographically random decimal PIN of the given length.
// length 0 returns "" (no-PIN mode). Length 4 produces a 4-digit PIN,
// length 6 a 6-digit PIN. Reduces in uint32 space before converting to int:
// int is 32-bit on the ARMv7 device and a raw Uint32 can exceed int32 max.
func generatePIN(length int) string {
	if length == 0 {
		return ""
	}
	var b [4]byte
	if _, err := crand.Read(b[:]); err != nil {
		if length == 4 {
			return "0000"
		}
		return "000000"
	}
	v := binary.BigEndian.Uint32(b[:])
	if length == 4 {
		return fmt.Sprintf("%04d", v%10000)
	}
	return fmt.Sprintf("%06d", v%1000000)
}

func generateToken() string {
	var b [16]byte
	if _, err := crand.Read(b[:]); err != nil {
		return "insecure-fallback"
	}
	return hex.EncodeToString(b[:])
}

var (
	lobbyIPMu         sync.Mutex
	lastPushedLobbyIP string
)

// readAllNotes returns metadata for every .md file in the notes directory.
func readAllNotes() []noteInfo {
	entries, err := os.ReadDir(notesDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []noteInfo{}
		}
		return nil
	}
	var notes []noteInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		notes = append(notes, noteInfo{
			Name:     e.Name(),
			Size:     info.Size(),
			Modified: info.ModTime().Format(time.RFC3339),
		})
	}
	if notes == nil {
		notes = []noteInfo{}
	}
	return notes
}

// pushNotesList sends the full note list to the editor for the Lobby Files page.
func pushNotesList() {
	notes := readAllNotes()
	msg, _ := json.Marshal(struct {
		T     string     `json:"t"`
		Items []noteInfo `json:"items"`
	}{"notes", notes})
	if globalEC != nil {
		globalEC.write(msg)
	}
}

// createNoteFile writes a new note; name is validated via notesSafe.
func createNoteFile(name, content string) error {
	p := notesSafe(name)
	if p == "" {
		return fmt.Errorf("invalid name")
	}
	if _, err := os.Stat(p); err == nil {
		return fmt.Errorf("already exists")
	}
	if content == "" {
		content = "# " + strings.TrimSuffix(filepath.Base(p), ".md") + "\n"
	}
	if rejectsHtmlNoteContent(content) {
		return fmt.Errorf("refusing HTML/qrichtext payload")
	}
	if err := writeNoteFile(p, []byte(content)); err != nil {
		return err
	}
	pushLobbyInfo()
	pushNotesList()
	return nil
}

// deleteNoteFile removes a note and notifies the editor if it was open.
func deleteNoteFile(name string) error {
	p := notesSafe(name)
	if p == "" {
		return fmt.Errorf("invalid name")
	}
	if err := os.Remove(p); err != nil {
		return err
	}
	currentNoteMu.Lock()
	wasOpen := currentNote != "" && filepath.Base(p) == currentNote
	if wasOpen {
		currentNote = ""
	}
	currentNoteMu.Unlock()
	if wasOpen && activeSess != nil && activeSess.isActive() {
		if globalEC != nil {
			globalEC.write([]byte(`{"t":"cmd","c":"notedeleted"}`))
		}
		broadcast([]byte(`{"type":"exitedit"}`))
	}
	pushLobbyInfo()
	pushNotesList()
	return nil
}

// renameNoteFile renames a note on disk and notifies the editor if it was open.
func renameNoteFile(oldName, newName string) error {
	oldP := notesSafe(oldName)
	newP := notesSafe(newName)
	if oldP == "" || newP == "" {
		return fmt.Errorf("invalid name")
	}
	if _, err := os.Stat(newP); err == nil {
		return fmt.Errorf("name already taken")
	}
	if err := os.Rename(oldP, newP); err != nil {
		return err
	}
	newBase := filepath.Base(newP)
	currentNoteMu.Lock()
	wasOpen := currentNote != "" && filepath.Base(oldP) == currentNote
	if wasOpen {
		currentNote = newBase
	}
	currentNoteMu.Unlock()
	if wasOpen && activeSess != nil && activeSess.isActive() {
		if globalEC != nil {
			cmd, _ := json.Marshal(struct {
				T    string `json:"t"`
				C    string `json:"c"`
				Name string `json:"name"`
			}{"cmd", "noterenamed", newBase})
			globalEC.write(cmd)
		}
		broadcastOpenEdit(newBase)
	}
	pushLobbyInfo()
	pushNotesList()
	return nil
}

// handleEditorReq serves trusted file ops from the local editor over the socket.
func handleEditorReq(op, name, oldName string) {
	switch op {
	case "noteslist":
		pushNotesList()
	case "createnote":
		if err := createNoteFile(name, ""); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: editor createnote: %v\n", err)
		} else {
			notifyTabletCrud("createnote", name, "")
		}
	case "deletenote":
		if err := deleteNoteFile(name); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: editor deletenote: %v\n", err)
		} else {
			notifyTabletCrud("deletenote", name, "")
		}
	case "renamenote":
		if err := renameNoteFile(oldName, name); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: editor renamenote: %v\n", err)
		} else {
			notifyTabletCrud("renamenote", name, oldName)
		}
	default:
		fmt.Fprintf(os.Stderr, "writerdeck-server: unknown editor req op %q\n", op)
	}
}

// countNotes returns the number of .md files in the notes directory.
func countNotes() int {
	entries, err := os.ReadDir(notesDirPath)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			n++
		}
	}
	return n
}

// formatLastSyncAgo turns a unix timestamp into a short relative time string.
func formatLastSyncAgo(unix int64) string {
	if unix <= 0 {
		return ""
	}
	d := time.Since(time.Unix(unix, 0))
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	}
	if d < 24*time.Hour {
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	}
	days := int(d.Hours() / 24)
	if days == 1 {
		return "1 day ago"
	}
	return fmt.Sprintf("%d days ago", days)
}

// pushLobbyInfo sends {"t":"info",...} to the editor socket so the Lobby
// reflects current connect, sync, and notes state.
func pushLobbyInfo() {
	ip := getLocalIP()
	authMu.Lock()
	pin := authPIN
	authMu.Unlock()
	settingsMu.Lock()
	syncOn := curSettings.SyncOn
	syncRepo := curSettings.SyncRepo
	lastSync := formatLastSyncAgo(curSettings.LastSyncAt)
	if n := len(curSettings.PendingSync); n > 0 {
		pending := "sync pending"
		if n > 1 {
			pending = fmt.Sprintf("%d sync ops pending", n)
		}
		if lastSync == "" {
			lastSync = pending + " — open phone browser"
		} else {
			lastSync = lastSync + " (" + pending + ")"
		}
	}
	settingsMu.Unlock()
	infoMsg, _ := json.Marshal(struct {
		T         string `json:"t"`
		IP        string `json:"ip"`
		PIN       string `json:"pin"`
		SyncOn    bool   `json:"syncOn"`
		SyncRepo  string `json:"syncRepo"`
		NoteCount int    `json:"noteCount"`
		LastSync  string `json:"lastSync"`
	}{"info", ip, pin, syncOn, syncRepo, countNotes(), lastSync})
	if globalEC != nil {
		globalEC.write(infoMsg)
	}
	lobbyIPMu.Lock()
	lastPushedLobbyIP = ip
	lobbyIPMu.Unlock()
}

// watchLobbyIP re-pushes lobby info when wlan0 gets an address after boot or
// when the DHCP lease changes. pushLobbyInfo on socket connect alone is too
// early on cold boot: DHCP often completes after keywriter has already shown
// http://?:8000.
func watchLobbyIP() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		ip := getLocalIP()
		lobbyIPMu.Lock()
		changed := ip != lastPushedLobbyIP
		lobbyIPMu.Unlock()
		if !changed || globalEC == nil || !globalEC.ready() {
			continue
		}
		pushLobbyInfo()
		pushNotesList()
	}
}

// checkAuth returns true if the request is authorized.
// Always returns true for OPTIONS (preflight) or when PIN auth is disabled.
// Loopback requests are trusted (tablet editor HTTP saves — slice 10).
// When PIN auth is enabled, checks the writerdeck_token session cookie.
func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodOptions {
		return true
	}
	if isLoopback(r) {
		return true
	}
	authMu.Lock()
	required := pinRequired
	tok := authToken
	authMu.Unlock()
	if !required {
		return true
	}
	cookie, err := r.Cookie("writerdeck_token")
	if err != nil || cookie.Value != tok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

// isLoopback reports whether r came from the tablet editor on localhost.
func isLoopback(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	return host == "127.0.0.1" || host == "::1"
}

// nextMorningCutoff returns the next local 04:00. If it is currently before
// today's 04:00 it returns today's 04:00; otherwise tomorrow's. Used as the
// auth-cookie expiry so a full day's writing never re-prompts for the PIN,
// while a fresh morning re-gates once. (A reboot also re-gates independently:
// authToken is regenerated per boot, so a stale cookie value stops matching --
// and a mid-day rmkbd restart likewise asks once, an accepted cost of not
// persisting the token to disk.)
func nextMorningCutoff(now time.Time) time.Time {
	cutoff := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
	if !now.Before(cutoff) {
		cutoff = cutoff.AddDate(0, 0, 1)
	}
	return cutoff
}

// --- PIN brute-force throttle ---
// A 6-digit PIN has 1,000,000 combinations; without throttling, someone on the
// LAN could exhaust it. We lock an IP out for pinLockout after pinMaxFails
// consecutive wrong guesses (HTTP 429 + Retry-After). State is in-memory and
// per-IP: a reboot clears it (and regenerates the PIN), and tracking by IP
// means a malicious actor cannot lock out the owner, who connects from a
// different address.
const (
	pinMaxFails = 5
	pinLockout  = 60 * time.Second
)

type pinAttempt struct {
	fails       int
	lockedUntil time.Time
}

var (
	pinMu       sync.Mutex
	pinAttempts = map[string]*pinAttempt{}
)

// clientIP returns the host portion of r.RemoteAddr (no port), or the raw
// value if it cannot be split.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// pinHandler handles POST /api/pin: validates the PIN and sets a session cookie.
func pinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		PIN string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Read auth state once under lock; use local copies for the rest of the handler.
	authMu.Lock()
	pin := authPIN
	required := pinRequired
	authMu.Unlock()

	// No PIN required: just issue a cookie (handles a client stuck on the PIN
	// screen when the owner switches to no-PIN mode) and return OK.
	if !required {
		authMu.Lock()
		tok := authToken
		authMu.Unlock()
		exp := nextMorningCutoff(time.Now())
		http.SetCookie(w, &http.Cookie{
			Name:     "writerdeck_token",
			Value:    tok,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Expires:  exp,
			MaxAge:   int(time.Until(exp).Seconds()),
		})
		w.WriteHeader(http.StatusOK)
		return
	}

	ip := clientIP(r)
	now := time.Now()
	pinMu.Lock()
	// Prune served-out lockouts so the map stays small and an expired lockout
	// resets that IP's counter automatically.
	for k, a := range pinAttempts {
		if a.fails >= pinMaxFails && now.After(a.lockedUntil) {
			delete(pinAttempts, k)
		}
	}
	if a := pinAttempts[ip]; a != nil && a.fails >= pinMaxFails && now.Before(a.lockedUntil) {
		retry := int(a.lockedUntil.Sub(now).Seconds()) + 1
		pinMu.Unlock()
		w.Header().Set("Retry-After", fmt.Sprintf("%d", retry))
		http.Error(w, "too many attempts", http.StatusTooManyRequests)
		return
	}
	// Constant-time compare so the response time does not leak how many leading
	// digits matched (the per-IP lockout above is the primary defense).
	if subtle.ConstantTimeCompare([]byte(req.PIN), []byte(pin)) != 1 {
		a := pinAttempts[ip]
		if a == nil {
			a = &pinAttempt{}
			pinAttempts[ip] = a
		}
		a.fails++
		locked := a.fails >= pinMaxFails
		if locked {
			a.lockedUntil = now.Add(pinLockout)
		}
		pinMu.Unlock()
		if locked {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(pinLockout.Seconds())))
			http.Error(w, "too many attempts", http.StatusTooManyRequests)
			return
		}
		http.Error(w, "wrong PIN", http.StatusUnauthorized)
		return
	}
	// Success: clear this IP's failure record.
	delete(pinAttempts, ip)
	pinMu.Unlock()

	// Re-read token under lock in case a concurrent PIN-length change regenerated
	// it during this request (so the issued cookie matches the current authToken).
	authMu.Lock()
	currentTok := authToken
	authMu.Unlock()

	// Expire the cookie at the next local 04:00 so a full day's writing never
	// re-prompts for the PIN, but a fresh morning (and any reboot) re-gates once.
	// Set both Expires and MaxAge: MaxAge wins in modern browsers, Expires is the
	// fallback for older ones -- both point at the same wall-clock moment.
	exp := nextMorningCutoff(time.Now())
	http.SetCookie(w, &http.Cookie{
		Name:     "writerdeck_token",
		Value:    currentTok,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  exp,
		MaxAge:   int(time.Until(exp).Seconds()),
	})
	w.WriteHeader(http.StatusOK)
}

// --- Lobby-on-demand ---
// A rate-limited pre-auth endpoint that tells the editor to show the Lobby.
// Pre-auth by design: reveals the PIN only on the physical e-ink screen, never
// over the network, so it does not weaken the "must hold the device" model.
var (
	lobbyMu      sync.Mutex
	lobbyLastReq time.Time
	lobbyMinGap  = 3 * time.Second
)

// lobbyHandler handles POST /api/lobby: tells the editor to show the Lobby.
// Registered outside checkAuth (pre-auth) -- reveals PIN only on the e-ink,
// not over the network. Rate-limited to ~3 s so a LAN actor cannot spam
// Lobby-flips (notes are saved before the view switch, so it is annoying
// but not data-destructive).
func lobbyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	now := time.Now()
	lobbyMu.Lock()
	if !lobbyLastReq.IsZero() && now.Sub(lobbyLastReq) < lobbyMinGap {
		remaining := lobbyMinGap - now.Sub(lobbyLastReq)
		lobbyMu.Unlock()
		w.Header().Set("Retry-After", fmt.Sprintf("%d", int(remaining.Seconds())+1))
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}
	lobbyLastReq = now
	lobbyMu.Unlock()

	if activeSess == nil {
		http.Error(w, "not in supervisor mode", http.StatusNotImplemented)
		return
	}
	if activeSess.isActive() {
		// Editor is running: save current note then show Lobby; wait for save ack.
		activeSess.ec.writeCmdWaitAck([]byte(`{"t":"cmd","c":"showlobby"}`), "saved", "showlobby", saveAckTimeout)
	} else {
		// No active session: start one -- it boots directly into the Lobby.
		if err := activeSess.start(); err != nil {
			http.Error(w, "could not start session: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	pushLobbyInfo()
	w.WriteHeader(http.StatusOK)
}

// statusHandler serves GET /api/status: tablet battery, Wi-Fi, editor state, open note.
func statusHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	battery := -1
	charging := false
	if data, err := os.ReadFile("/sys/class/power_supply/bq27441-0/capacity"); err == nil {
		if p, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil && p >= 0 && p <= 100 {
			battery = p
		}
	}
	if st, err := os.ReadFile("/sys/class/power_supply/bq27441-0/status"); err == nil {
		s := strings.TrimSpace(string(st))
		charging = s == "Charging" || s == "Full"
	}
	wifi := false
	if st, err := os.ReadFile("/sys/class/net/wlan0/operstate"); err == nil {
		wifi = strings.TrimSpace(string(st)) == "up"
	}
	editorActive := false
	if activeSess != nil {
		editorActive = activeSess.isActive()
	}
	currentNoteMu.Lock()
	openNote := currentNote
	currentNoteMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct { //nolint:errcheck
		Battery      int    `json:"battery"`
		Charging     bool   `json:"charging"`
		Wifi         bool   `json:"wifi"`
		IP           string `json:"ip"`
		EditorActive bool   `json:"editorActive"`
		OpenNote     string `json:"openNote"`
	}{battery, charging, wifi, getLocalIP(), editorActive, openNote})
}

// shutdownHandler handles POST /api/shutdown: end the editor session, restore
// xochitl, and exit rmkbd. The phone companion disconnects until rmkbd is
// started again (SSH, reboot hook, or systemd).
func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if activeSess == nil {
		http.Error(w, "not in supervisor mode", http.StatusNotImplemented)
		return
	}
	w.WriteHeader(http.StatusOK)
	go func() {
		time.Sleep(200 * time.Millisecond)
		if activeSess.isActive() {
			activeSess.quit()
		} else {
			exec.Command("systemctl", "start", "xochitl").Run() //nolint:errcheck
		}
		fmt.Fprintln(os.Stderr, "writerdeck-server: shutdown requested from phone -- exiting")
		os.Exit(0)
	}()
}

// syncAckHandler handles POST /api/sync/ack: the phone browser calls this after
// reconcileAll completes so power-button sleep can suspend without cutting Wi-Fi
// mid-upload.
func syncAckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signalSyncAck()
	settingsMu.Lock()
	curSettings.LastSyncAt = time.Now().Unix()
	saveSettingsLocked()
	settingsMu.Unlock()
	pushLobbyInfo()
	w.WriteHeader(http.StatusOK)
}

// pendingSyncHandler serves GET /api/sync/pending (queued tablet CRUD ops).
func pendingSyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	settingsMu.Lock()
	ops := curSettings.PendingSync
	if ops == nil {
		ops = []pendingSyncOp{}
	}
	settingsMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ops) //nolint:errcheck
}

// pendingClearHandler handles POST /api/sync/pending/clear after the phone pairs
// queued tablet ops on GitHub.
func pendingClearHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	clearPendingSync()
	w.WriteHeader(http.StatusOK)
}

// reloadHandler handles POST /api/reload: tell the tablet editor to reload the
// open note from disk (slice 8 — after a pull/clash changed disk under the buffer).
func reloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if activeSess == nil || !activeSess.isActive() || globalEC == nil {
		http.Error(w, "editor not active", http.StatusServiceUnavailable)
		return
	}
	currentNoteMu.Lock()
	name := currentNote
	currentNoteMu.Unlock()
	if name == "" {
		http.Error(w, "no open note", http.StatusConflict)
		return
	}
	globalEC.write([]byte(`{"t":"cmd","c":"reloadnote"}`))
	w.WriteHeader(http.StatusOK)
}

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
	cmd := exec.Command(s.editorPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Setpgid gives the editor+inhibit wrapper their own process group so a
	// Kill(-pgid, SIGTERM) SIGTERM fallback reaches all child processes.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
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
	s.ec.write([]byte(`{"t":"cmd","c":"quit"}`))
	select {
	case <-doneCh:
		fmt.Fprintln(os.Stderr, "writerdeck-server: session: editor exited cleanly")
	case <-time.After(3 * time.Second):
		fmt.Fprintln(os.Stderr, "writerdeck-server: session: 3s timeout -- SIGTERM to process group")
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
	beginSyncWait()
	broadcast([]byte(`{"type":"exitedit","source":"power","awaitSync":true}`))
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
	if noteName == "" {
		return nil
	}
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
	return nil
}

// launchHandler handles POST /api/launch: starts an editor session if idle.
// Returns 409 if a session is already active; 501 if not in supervisor mode.
func launchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if activeSess == nil {
		http.Error(w, "not in supervisor mode", http.StatusNotImplemented)
		return
	}
	if err := activeSess.start(); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// openHandler handles POST /api/open: opens a specific note in the editor.
// If no session is active, starts one first (same as /api/launch).
// Sends {"t":"cmd","c":"open","name":"<file>"} to the editor, which calls
// saveAndLoad(name) in QML -- saves the current note then loads the new one.
func openHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "bad request: need {name}", http.StatusBadRequest)
		return
	}
	p := notesSafe(req.Name)
	if p == "" {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}
	if activeSess == nil {
		http.Error(w, "not in supervisor mode", http.StatusNotImplemented)
		return
	}
	// Ensure an editor session is running; start one if idle.
	if !activeSess.isActive() {
		if err := activeSess.start(); err != nil {
			http.Error(w, "could not start editor: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Wait up to 5 s for the editor socket to connect. Needed when we just
	// started a fresh session above; already-active sessions are instant.
	for i := 0; i < 10; i++ {
		if activeSess.ec.ready() {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	// Send saveAndLoad command -- QML saves the current note then calls doLoad(name).
	editorName := filepath.Base(p) // e.g. "scratch.md"
	cmd, _ := json.Marshal(struct {
		T    string `json:"t"`
		C    string `json:"c"`
		Name string `json:"name"`
	}{"cmd", "open", editorName})
	currentNoteMu.Lock()
	currentNote = editorName
	currentNoteMu.Unlock()
	broadcastOpenEdit(editorName)
	if !activeSess.ec.writeCmdWaitAck(cmd, "saved", "open", saveAckTimeout) {
		fmt.Fprintln(os.Stderr, "writerdeck-server: open save ack missed -- continuing")
	}
	w.WriteHeader(http.StatusOK)
}

// rotateHandler handles POST /api/rotate: rotates the editor display 90 degrees
// clockwise, persists the new angle to settings.json, and pushes it to Writerdeck.
func rotateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if activeSess == nil {
		http.Error(w, "not in supervisor mode", http.StatusNotImplemented)
		return
	}
	if !activeSess.isActive() {
		http.Error(w, "no active editor session", http.StatusConflict)
		return
	}
	settingsMu.Lock()
	curSettings.Rotation = normalizeRotation(curSettings.Rotation + 90)
	rotation := curSettings.Rotation
	saveSettingsLocked()
	settingsMu.Unlock()
	rotMsg, _ := json.Marshal(struct {
		T       string `json:"t"`
		C       string `json:"c"`
		Degrees int    `json:"degrees"`
	}{"cmd", "setrotation", rotation})
	activeSess.ec.write(rotMsg)
	w.WriteHeader(http.StatusOK)
}

// selftest sends "hello world" + Return over the editor socket,
// replicating the Phase 2 smoke test without needing a browser.
// No leading Escape: keywriter now boots in edit mode.
func selftest(ec *editorConn) {
	fmt.Fprintln(os.Stderr, "writerdeck-server: --selftest: waiting for editor socket...")
	for !ec.ready() {
		time.Sleep(500 * time.Millisecond)
	}
	time.Sleep(3 * time.Second) // let keywriter finish QML init

	send := func(line []byte) {
		fmt.Fprintf(os.Stderr, "writerdeck-server: selftest send %s\n", line)
		ec.write(line)
		time.Sleep(100 * time.Millisecond)
	}

	// No leading Escape: keywriter now boots in edit mode (see socket-inject.patch).
	// Sending Escape here would toggle *out* of edit mode.
	for _, r := range "hello world" {
		send([]byte(fmt.Sprintf(`{"t":"text","cp":%d}`, r)))
	}
	send([]byte(`{"t":"key","k":"Return"}`))
	fmt.Fprintln(os.Stderr, "writerdeck-server: selftest done")
}

func main() {
	port := flag.Int("port", defaultPort, "WebSocket server port")
	doSelftest := flag.Bool("selftest", false, "send hardcoded hello world+Return and exit")
	verbose := flag.Bool("v", false, "verbose: log every translated key (default logs only a periodic count)")
	editorPath := flag.String("editor", "", "path to editor launch script; rmkbd spawns it as a child and owns its lifecycle (supervisor mode, used by systemd unit)")
	flag.StringVar(&notesDirPath, "notes-dir", notesDirPath, "directory for .md notes (default: /home/root/Writerdeck-user-documents; override for local dev)")
	flag.StringVar(&settingsFilePath, "settings-file", settingsFilePath, "path to settings JSON (default /home/root/.Writerdeck/settings.json; override for local dev)")
	flag.Parse()

	loadSettings()

	// Determine PIN length from persisted settings (loaded above).
	var bootPinLen int
	switch curSettings.PinDigits {
	case "4":
		bootPinLen = 4
	case "none":
		bootPinLen = 0
	default:
		bootPinLen = 6 // "6" or any unrecognised value
	}
	authMu.Lock()
	authPIN = generatePIN(bootPinLen)
	authToken = generateToken()
	pinRequired = bootPinLen > 0
	authMu.Unlock()
	if authPIN != "" {
		fmt.Fprintf(os.Stderr, "writerdeck-server: PIN is %s (will be shown on tablet Lobby; for now read from this log)\n", authPIN)
	} else {
		fmt.Fprintln(os.Stderr, "writerdeck-server: no PIN required (pinDigits=none)")
	}

	ec := &editorConn{}
	globalEC = ec
	go dialLoop(ec)
	go watchLobbyIP()

	if *doSelftest {
		selftest(ec)
		return
	}

	addr := fmt.Sprintf(":%d", *port)
	fmt.Fprintf(os.Stderr, "writerdeck-server: serving capture page on http://<device-ip>%s/\n", addr)
	fmt.Fprintf(os.Stderr, "writerdeck-server: serving WebSocket on %s/ws\n", addr)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/app.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(appCSS) //nolint:errcheck
	})
	http.HandleFunc("/app.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Write(appJS) //nolint:errcheck
	})
	http.HandleFunc("/state.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Write(stateJS) //nolint:errcheck
	})
	http.HandleFunc("/sync.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Write(syncJS) //nolint:errcheck
	})
	http.HandleFunc("/ws", wsHandler(ec, *verbose))
	http.HandleFunc("/api/pin", pinHandler)
	http.HandleFunc("/api/launch", launchHandler)
	http.HandleFunc("/api/open", openHandler)
	http.HandleFunc("/api/rotate", rotateHandler)
	http.HandleFunc("/api/notes", notesListHandler)
	http.HandleFunc("/api/notes/", notesItemHandler)
	http.HandleFunc("/api/settings", settingsHandler)
	http.HandleFunc("/api/lobby", lobbyHandler) // pre-auth: reveals PIN on e-ink only
	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/api/shutdown", shutdownHandler)
	http.HandleFunc("/api/sync/ack", syncAckHandler)
	http.HandleFunc("/api/sync/pending", pendingSyncHandler)
	http.HandleFunc("/api/sync/pending/clear", pendingClearHandler)
	http.HandleFunc("/api/reload", reloadHandler)

	if *editorPath != "" {
		// Supervisor mode: rmkbd is always-on; editor sessions are on-demand.
		// xochitl stop/start happens per session in Go (session.start/end).
		// ExecStopPost=start xochitl in the unit stays as a safety net.
		activeSess = &session{editorPath: *editorPath, ec: ec}

		// Always-on Home watcher: loops for rmkbd's lifetime.
		go watchHomeButton(activeSess, ec)
		go watchUSBKeyboardForLaunch(activeSess)

		// HTTP server always-on in the background.
		go func() {
			if err := http.ListenAndServe(addr, nil); err != nil {
				fmt.Fprintf(os.Stderr, "writerdeck-server: HTTP server: %v\n", err)
			}
		}()

		// Reconcile: kill any stray keywriter from a previous crash so two
		// editors don't fight for the framebuffer on startup.
		fmt.Fprintln(os.Stderr, "writerdeck-server: reconcile: killing any stray Writerdeck editor")
		exec.Command("sh", "-c", "for p in $(pidof Writerdeck 2>/dev/null); do kill $p 2>/dev/null; done").Run() //nolint:errcheck
		time.Sleep(500 * time.Millisecond)

		// Auto-launch first session: power-on = typewriter (unchanged behaviour).
		fmt.Fprintln(os.Stderr, "writerdeck-server: auto-launching editor session on boot")
		if err := activeSess.start(); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: auto-launch failed: %v\n", err)
		}

		// Block until SIGTERM/SIGINT; gracefully end any active session first.
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		sig := <-sigCh
		fmt.Fprintf(os.Stderr, "writerdeck-server: signal %v received\n", sig)
		if activeSess.isActive() {
			fmt.Fprintln(os.Stderr, "writerdeck-server: ending active session before exit")
			activeSess.quit()
		}
		fmt.Fprintln(os.Stderr, "writerdeck-server: exiting (ExecStopPost safety net restarts xochitl if needed)")
		os.Exit(0)
	}

	// Stand-alone mode: dev/test scripts manage keywriter separately.
	// Still watch the home button: sends a single quit to ec (one-shot).
	go watchHomeButton(nil, ec)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: server error: %v\n", err)
		os.Exit(1)
	}
}
