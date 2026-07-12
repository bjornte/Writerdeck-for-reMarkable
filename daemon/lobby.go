// Writerdeck-server — see main.go for overview.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

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
	keyboardLayout := curSettings.KeyboardLayout
	if keyboardLayout == "" {
		keyboardLayout = "us"
	}
	lastSync := formatLastSyncAgo(curSettings.LastSyncAt)
	if n := len(curSettings.PendingSync); n > 0 {
		pending := "sync pending"
		if n > 1 {
			pending = fmt.Sprintf("%d sync ops pending", n)
		}
		if lastSync == "" {
			lastSync = pending
		} else {
			lastSync = lastSync + " (" + pending + ")"
		}
	}
	settingsMu.Unlock()
	syncReady := syncEng.ready()
	syncing := syncEng.isSyncing()
	if syncOn && syncRepo != "" && !syncReady {
		lastSync = "Token needed — add in phone Setup"
	}
	infoMsg, _ := json.Marshal(struct {
		T              string `json:"t"`
		IP             string `json:"ip"`
		PIN            string `json:"pin"`
		SyncOn         bool   `json:"syncOn"`
		SyncRepo       string `json:"syncRepo"`
		NoteCount      int    `json:"noteCount"`
		LastSync       string `json:"lastSync"`
		SyncReady      bool   `json:"syncReady"`
		Syncing        bool   `json:"syncing"`
		KeyboardLayout string `json:"keyboardLayout"`
	}{"info", ip, pin, syncOn, syncRepo, countNotes(), lastSync, syncReady, syncing, keyboardLayout})
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
