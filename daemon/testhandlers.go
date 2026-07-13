package main

import (
	"encoding/json"
	"net/http"
)

// testResetHandler quits the active editor session (hard reset). The harness
// normally uses soft reset (PUT + reload + Home) and calls this only once per run.
func testResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if activeSess != nil && activeSess.isActive() {
		activeSess.quit()
	}
	currentNoteMu.Lock()
	currentNote = ""
	currentNoteMu.Unlock()
	w.WriteHeader(http.StatusOK)
}

// testHomeHandler sends the same home cmd as the physical Home button while editing.
func testHomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if activeSess == nil || !activeSess.isActive() {
		http.Error(w, "no active editor session", http.StatusConflict)
		return
	}
	activeSess.ec.writeCmdWaitAck([]byte(`{"t":"cmd","c":"home"}`), "saved", "home", saveAckTimeout)
	currentNoteMu.Lock()
	currentNote = ""
	currentNoteMu.Unlock()
	broadcast([]byte(`{"type":"exitedit","source":"test"}`))
	w.WriteHeader(http.StatusOK)
}

func testEditorStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if activeSess == nil || !activeSess.isActive() {
		http.Error(w, "no active editor session", http.StatusConflict)
		return
	}
	st, err := queryEditorState()
	if err != nil {
		http.Error(w, err.Error(), http.StatusGatewayTimeout)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(st) //nolint:errcheck
}

// testTabletReqHandler exercises trusted tablet socket ops (setreadfont, setpindigits).
func testTabletReqHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
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
		Op   string `json:"op"`
		Name string `json:"name"`
		Old  string `json:"old"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Op == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	switch req.Op {
	case "setreadfont", "setpindigits", "setkeyboardlayout",
		"setvaultpin", "changevaultpin", "unlockvault", "lockvault",
		"encryptnote", "decryptnote", "disablevault":
		handleEditorReq(req.Op, req.Name, req.Old)
	default:
		http.Error(w, "unsupported op", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
