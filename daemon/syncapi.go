package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func syncTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		if !syncEng.tokenConfigured() {
			json.NewEncoder(w).Encode(map[string]bool{"configured": false}) //nolint:errcheck
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"configured": true,
			"token":      syncEng.getToken(),
		})
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Token == "" {
		syncEng.clearToken()
		syncEng.setLastError("")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true, "configured": false}) //nolint:errcheck
		return
	}
	settingsMu.Lock()
	repo := curSettings.SyncRepo
	settingsMu.Unlock()
	if repo == "" {
		http.Error(w, "set syncRepo first", http.StatusBadRequest)
		return
	}
	status, err := syncEng.verifyRepo(repo, req.Token)
	if err != nil {
		http.Error(w, "could not reach GitHub", http.StatusBadGateway)
		return
	}
	if status == 401 || status == 403 {
		http.Error(w, "token rejected", http.StatusUnauthorized)
		return
	}
	if status == 404 {
		http.Error(w, "repo not found", http.StatusNotFound)
		return
	}
	if status != 200 {
		http.Error(w, "github error", http.StatusBadGateway)
		return
	}
	syncEng.setToken(req.Token)
	syncEng.setLastError("")
	pushLobbyInfo()
	go func() {
		n, err := syncEng.reconcileAll("token")
		if err == nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: token verify reconcile: %d notes\n", n)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
		"ok": true, "configured": true, "verified": true,
	})
}

func syncStatusHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	settingsMu.Lock()
	syncOn := curSettings.SyncOn
	repo := curSettings.SyncRepo
	lastSync := curSettings.LastSyncAt
	settingsMu.Unlock()
	syncEng.syncingMu.Lock()
	syncing := syncEng.syncing
	syncEng.syncingMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
		"syncOn":      syncOn,
		"syncRepo":    repo,
		"configured":  syncEng.tokenConfigured(),
		"lastSyncAt":  lastSync,
		"lastError":   syncEng.getLastError(),
		"syncing":     syncing,
		"lastSyncAgo": formatLastSyncAgo(lastSync),
	})
}

func syncRunHandler(w http.ResponseWriter, r *http.Request) {
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
	if !syncEng.ready() {
		http.Error(w, "sync not configured", http.StatusBadRequest)
		return
	}
	n, err := syncEng.reconcileAll("manual")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
		"ok": true, "notes": n,
	})
}
