// Writerdeck-server — see main.go for overview.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

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

// wsMsg is the JSON message received from the browser on keydown / paste.
type wsMsg struct {
	Type   string `json:"type"`             // "key", "hello", or "paste"
	Key    string `json:"key"`              // KeyboardEvent.key value (key messages)
	Text   string `json:"text,omitempty"`   // paste body (paste messages)
	Shift  bool   `json:"shift"`            // event.shiftKey
	Ctrl   bool   `json:"ctrl"`             // event.ctrlKey
	Alt    bool   `json:"alt"`              // event.altKey
	Meta   bool   `json:"meta"`             // event.metaKey (Cmd on Mac/iPhone)
	Action string `json:"action,omitempty"` // "release" for key-up replay
}

// maxPasteRunes caps a single phone paste (Lobby Files never accepts paste).
const maxPasteRunes = 200000

// forwardPaste inserts clipboard text at the editor cursor. Refused when no
// note is open (Lobby Files / idle) so paste cannot dump into name prompts.
func forwardPaste(ec *editorConn, text string) {
	if ec == nil || text == "" {
		return
	}
	currentNoteMu.Lock()
	note := currentNote
	currentNoteMu.Unlock()
	if note == "" {
		fmt.Fprintln(os.Stderr, "writerdeck-server: paste ignored (no note open)")
		return
	}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	n := 0
	for _, r := range text {
		if n >= maxPasteRunes {
			fmt.Fprintf(os.Stderr, "writerdeck-server: paste truncated at %d runes\n", maxPasteRunes)
			break
		}
		n++
		var line []byte
		switch r {
		case '\n':
			line = []byte(`{"t":"key","k":"Return"}`)
		case '\t':
			line = []byte(`{"t":"key","k":"Tab"}`)
		default:
			line = []byte(fmt.Sprintf(`{"t":"text","cp":%d}`, r))
		}
		ec.write(line)
	}
}

// namedKeys maps browser KeyboardEvent.key values to keywriter named keys.
// Only keys whose value is NOT a single printable rune need an entry here.
var namedKeys = map[string]string{
	"Enter":      "Return",
	"Backspace":  "Backspace",
	"Delete":     "Delete",
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

	stateMu   sync.Mutex
	stateWait chan EditorState
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
	if st, ok := parseEditorState(line); ok {
		e.deliverState(st)
		return
	}
	var msg struct {
		T       string `json:"t"`
		C       string `json:"c"`
		Op      string `json:"op"`
		Name    string `json:"name"`
		Old     string `json:"old"`
		Mode    string `json:"mode"`
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
	case "openread":
		name := filepath.Base(msg.Name)
		if notesSafe(name) == "" {
			return
		}
		currentNoteMu.Lock()
		currentNote = name
		currentNoteMu.Unlock()
		broadcastOpenRead(name)
	case "lobbyinput":
		broadcastLobbyInput(msg.Mode)
	case "saved", "ready":
		if msg.T == "saved" && (msg.C == "home" || msg.C == "showlobby") {
			// Returning to the Lobby means no note is "open" for sync purposes.
			// Otherwise, sync would keep skipping the last edited note.
			currentNoteMu.Lock()
			currentNote = ""
			currentNoteMu.Unlock()
			vaultClearSessionOnLobby()
		}
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

	// Named control key (arrow, Enter, Backspace, Tab, Escape, Home, End, Delete)?
	if kwKey, ok := namedKeys[key]; ok {
		up := ev.Action == "release"
		if mask != 0 {
			if up {
				return []byte(fmt.Sprintf(`{"t":"key","k":%q,"m":%d,"u":1}`, kwKey, mask))
			}
			return []byte(fmt.Sprintf(`{"t":"key","k":%q,"m":%d}`, kwKey, mask))
		}
		if up {
			return []byte(fmt.Sprintf(`{"t":"key","k":%q,"u":1}`, kwKey))
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
// restartEditorForKeymap relaunches Writerdeck so the launcher re-reads
// settings.json and applies the .qmap on a USB keyboard (Qt reads keymap at process start).
func restartEditorForKeymap(reason string) {
	if activeSess == nil || !activeSess.isActive() {
		return
	}
	go func() {
		fmt.Fprintf(os.Stderr, "writerdeck-server: %s -- restarting editor for keymap\n", reason)
		activeSess.quit()
		if err := activeSess.start(); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: restart after keymap change: %v\n", err)
		}
	}()
}

// handleEditorReq serves trusted file ops from the local editor over the socket.
func handleEditorReq(op, name, oldName string) {
	switch op {
	case "noteslist":
		pushNotesList()
	case "lobbyinfo":
		// Lobby Keyboard tab polls presence (phone page / USB).
		pushLobbyInfo()
	case "createnote":
		if err := createNoteFile(name, ""); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: editor createnote: %v\n", err)
			pushVaultOpFailed(noteOpErrMsg("create", err))
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
			pushVaultOpFailed(noteOpErrMsg("rename", err))
		} else if notesSafe(oldName) != notesSafe(name) {
			notifyTabletCrud("renamenote", name, oldName)
		}
	case "syncnow":
		if !syncEng.ready() {
			fmt.Fprintln(os.Stderr, "writerdeck-server: syncnow ignored — not configured")
			pushLobbyInfo()
			return
		}
		go func() { _, _ = syncEng.reconcileAll("tablet") }()
	case "setkeyboardlayout":
		layout := normalizeKeyboardLayout(name)
		settingsMu.Lock()
		curSettings.KeyboardLayout = layout
		saveSettingsLocked()
		settingsMu.Unlock()
		pushLobbyInfo()
		restartEditorForKeymap("keyboard layout changed")
	case "setreadfont":
		if !applyReadFont(name) {
			fmt.Fprintf(os.Stderr, "writerdeck-server: editor setreadfont: unknown font %q\n", name)
		}
	case "setpindigits":
		if !applyPinDigits(name) {
			fmt.Fprintf(os.Stderr, "writerdeck-server: editor setpindigits: invalid %q\n", name)
		}
	case "languagechanged":
		// Lobby wrote lobby-ui.json; tell open phone pages to reload strings.
		broadcast([]byte(`{"type":"language"}`))
		fmt.Fprintln(os.Stderr, "writerdeck-server: lobby language changed — phone UI notified")
	case "setvaultpin":
		if err := vaultSetupPIN(name); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: setvaultpin: %v\n", err)
		}
	case "changevaultpin":
		if err := vaultChangePIN(oldName, name); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: changevaultpin: %v\n", err)
		}
	case "verifyvaultpin":
		keepSession := oldName == "session"
		if err := vaultVerifyPIN(name, keepSession); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: verifyvaultpin: %v\n", err)
			msg := "Wrong PIN. Try again."
			if strings.Contains(err.Error(), "too many") {
				msg = "Too many attempts. Wait and try again."
			} else if strings.Contains(err.Error(), "invalid PIN") {
				msg = "PIN must be 6 digits."
			} else if strings.Contains(err.Error(), "not enabled") {
				msg = "Encryption is not enabled."
			}
			pushVaultOpFailed(msg)
		}
	case "encryptnote":
		if err := vaultEncryptNote(name); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: encryptnote: %v\n", err)
			pushVaultOpFailed(vaultOpErrMsg("encrypt", err))
		}
	case "decryptnote":
		if err := vaultDecryptNote(name); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: decryptnote: %v\n", err)
			pushVaultOpFailed(vaultOpErrMsg("decrypt", err))
		}
	case "disablevault":
		if err := vaultDisable(true); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: disablevault: %v\n", err)
		}
	case "shutdown":
		requestShutdown("tablet Settings")
	case "offerdownload":
		p := notesSafe(name)
		if p == "" {
			fmt.Fprintf(os.Stderr, "writerdeck-server: offerdownload: bad name %q\n", name)
			return
		}
		if _, err := os.Stat(p); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: offerdownload: %v\n", err)
			return
		}
		if !phoneConnected() {
			fmt.Fprintln(os.Stderr, "writerdeck-server: offerdownload: no phone page connected")
		}
		broadcastDownloadOffer(name)
	default:
		fmt.Fprintf(os.Stderr, "writerdeck-server: unknown editor req op %q\n", op)
	}
}
