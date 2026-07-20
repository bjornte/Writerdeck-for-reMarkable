package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// testResetHandler quits the active editor session. The keyboard harness does not
// use this; it sandbox-resets in-process via harnessprepare instead.
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
		Op    string   `json:"op"`
		Name  string   `json:"name"`
		Old   string   `json:"old"`
		Notes []string `json:"notes"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Op == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	switch req.Op {
	case "setreadfont", "setpindigits", "setkeyboardlayout",
		"setvaultpin", "changevaultpin", "verifyvaultpin",
		"encryptnote", "decryptnote", "disablevault",
		"offerdownload":
		handleEditorReq(req.Op, req.Name, req.Old)
	case "vaultrewrap":
		if err := vaultRewrapFromOldSecret(req.Old, req.Name, req.Notes); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "unsupported op", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// testEditorCmdHandler forwards {"t":"cmd","c":"..."} to the editor socket (harness UI triggers).
func testEditorCmdHandler(w http.ResponseWriter, r *http.Request) {
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
	if globalEC == nil || !globalEC.ready() {
		http.Error(w, "no active editor session", http.StatusConflict)
		return
	}
	var req struct {
		C       string `json:"c"`
		Name    string `json:"name,omitempty"`
		W       int    `json:"w,omitempty"`
		Pos     int    `json:"pos,omitempty"`
		Degrees *int   `json:"degrees,omitempty"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.C == "" {
		http.Error(w, "bad request: need {c}", http.StatusBadRequest)
		return
	}
	var line string
	switch {
	case req.Name != "" && req.W > 0:
		line = fmt.Sprintf(`{"t":"cmd","c":%q,"name":%q,"w":%d}`, req.C, req.Name, req.W)
	case req.Name != "":
		line = fmt.Sprintf(`{"t":"cmd","c":%q,"name":%q}`, req.C, req.Name)
	case req.W > 0:
		line = fmt.Sprintf(`{"t":"cmd","c":%q,"w":%d}`, req.C, req.W)
	case req.Pos > 0 || req.C == "harnesssetcursor":
		line = fmt.Sprintf(`{"t":"cmd","c":%q,"pos":%d}`, req.C, req.Pos)
	case req.Degrees != nil:
		line = fmt.Sprintf(`{"t":"cmd","c":%q,"degrees":%d}`, req.C, *req.Degrees)
	default:
		line = fmt.Sprintf(`{"t":"cmd","c":%q}`, req.C)
	}
	globalEC.write([]byte(line))
	if req.C == "harnessopen" && req.Name != "" {
		currentNoteMu.Lock()
		currentNote = req.Name
		currentNoteMu.Unlock()
	}
	if req.C == "harnessprepare" {
		if _, err := queryEditorState(); err != nil {
			http.Error(w, err.Error(), http.StatusGatewayTimeout)
			return
		}
	}
	if req.C == "harnesssetcursor" {
		if _, err := queryEditorState(); err != nil {
			http.Error(w, err.Error(), http.StatusGatewayTimeout)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
