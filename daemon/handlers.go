// Writerdeck-server — see main.go for overview.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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

// requestShutdown ends the editor session, restores xochitl, and exits the server.
func requestShutdown(source string) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		if activeSess != nil && activeSess.isActive() {
			activeSess.quit()
		} else {
			exec.Command("systemctl", "start", "xochitl").Run() //nolint:errcheck
		}
		fmt.Fprintf(os.Stderr, "writerdeck-server: shutdown requested from %s -- exiting\n", source)
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

// flushEditorSave asks the tablet editor to flush the open note buffer to disk.
// Used before deploy/shutdown; returns false if the save ack times out.
func flushEditorSave() bool {
	if activeSess == nil || !activeSess.isActive() || globalEC == nil || !globalEC.ready() {
		return true
	}
	currentNoteMu.Lock()
	open := currentNote != ""
	currentNoteMu.Unlock()
	if !open {
		return true
	}
	return globalEC.writeCmdWaitAck([]byte(`{"t":"cmd","c":"autosavenow"}`), "saved", "autosavenow", saveAckTimeout)
}

// flushSaveHandler handles POST /api/flush-save: save the open editor buffer
// before deploy or server stop. Loopback-trusted (tablet editor path).
func flushSaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if !isLoopback(r) && !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !flushEditorSave() {
		http.Error(w, "editor save ack missed", http.StatusGatewayTimeout)
		return
	}
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
	currentNoteMu.Lock()
	prevNote := currentNote
	currentNoteMu.Unlock()
	if syncEng.ready() && editorName != prevNote {
		_, _ = syncEng.pullNote(editorName)
	}
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
	if syncEng.ready() && prevNote != "" && prevNote != editorName {
		syncEng.tryPushNote(prevNote)
	}
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
