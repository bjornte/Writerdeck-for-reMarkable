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

// countNotes returns the number of .md and .md.enc files in the notes directory.
func countNotes() int {
	entries, err := os.ReadDir(notesDirPath)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() && isNoteListName(e.Name()) {
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

// listenPort is the HTTP/WebSocket port (set from main after flag parse).
var listenPort = defaultPort

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
	pinDigits := curSettings.PinDigits
	keyboardLayout := curSettings.KeyboardLayout
	if keyboardLayout == "" {
		keyboardLayout = "us"
	}
	lastSyncAt := curSettings.LastSyncAt
	pendingCount := len(curSettings.PendingSync)
	lastSync := formatLastSyncAgo(lastSyncAt)
	if pendingCount > 0 {
		pending := "sync pending"
		if pendingCount > 1 {
			pending = fmt.Sprintf("%d sync ops pending", pendingCount)
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
	syncErr := syncEng.getLastError()
	syncErrKey := ""
	wifi := wifiUp()
	if syncOn && syncRepo != "" && !syncReady {
		lastSync = "Token needed — add in phone Sync setup"
	}
	if syncOn && syncRepo != "" && syncReady && !wifi {
		syncErr = "No Wi-Fi - cannot reach GitHub"
		syncErrKey = "noWifi"
	}
	phoneURL := ""
	qrPath := ""
	if ip != "" {
		phoneURL = fmt.Sprintf("http://%s:%d", ip, listenPort)
		qrPath = ensurePhoneQR(phoneURL)
	}
	infoMsg, _ := json.Marshal(struct {
		T                 string `json:"t"`
		IP                string `json:"ip"`
		PIN               string `json:"pin"`
		SyncOn            bool   `json:"syncOn"`
		SyncRepo          string `json:"syncRepo"`
		NoteCount         int    `json:"noteCount"`
		LastSync          string `json:"lastSync"`
		LastSyncAt        int64  `json:"lastSyncAt"`
		SyncPending       int    `json:"syncPending"`
		SyncReady         bool   `json:"syncReady"`
		Syncing           bool   `json:"syncing"`
		SyncError         string `json:"syncError"`
		SyncErrorKey      string `json:"syncErrorKey,omitempty"`
		Wifi              bool   `json:"wifi"`
		KeyboardLayout    string `json:"keyboardLayout"`
		PinDigits         string `json:"pinDigits"`
		EncryptionEnabled bool   `json:"encryptionEnabled"`
		PhoneConnected    bool   `json:"phoneConnected"`
		UsbKeyboard       bool   `json:"usbKeyboard"`
		Port              int    `json:"port"`
		PhoneURL          string `json:"phoneUrl"`
		QrPath            string `json:"qrPath"`
	}{"info", ip, pin, syncOn, syncRepo, countNotes(), lastSync, lastSyncAt, pendingCount, syncReady, syncing, syncErr, syncErrKey, wifi, keyboardLayout, pinDigits, vaultEnabled(),
		phoneConnected(), usbKeyboardPresent(), listenPort, phoneURL, qrPath})
	if globalEC != nil {
		globalEC.write(infoMsg)
	}
	lobbyIPMu.Lock()
	lastPushedLobbyIP = ip
	lobbyIPMu.Unlock()
	maybeBroadcastNeedToken()
}

// watchLobbyIP re-pushes lobby info when wlan0 gets an address after boot or
// when the DHCP lease changes, when Wi-Fi goes up/down, or when phone/USB
// keyboard presence changes (Lobby no-keyboard tip).
func watchLobbyIP() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	var lastWifi bool
	var haveWifi bool
	var lastPhone, lastUSB bool
	var havePresence bool
	for range ticker.C {
		ip := getLocalIP()
		w := wifiUp()
		phone := phoneConnected()
		usb := usbKeyboardPresent()
		lobbyIPMu.Lock()
		ipChanged := ip != lastPushedLobbyIP
		lobbyIPMu.Unlock()
		wifiChanged := !haveWifi || w != lastWifi
		presenceChanged := !havePresence || phone != lastPhone || usb != lastUSB
		if (!ipChanged && !wifiChanged && !presenceChanged) || globalEC == nil || !globalEC.ready() {
			continue
		}
		if presenceChanged {
			wsClientsMu.Lock()
			n, ready := 0, 0
			for c := range wsClients {
				n++
				if c.hello && !ideBrowserUA(c.ua) {
					ready++
				}
			}
			wsClientsMu.Unlock()
			fmt.Fprintf(os.Stderr, "writerdeck-server: lobby presence phone=%v usb=%v ws=%d hello=%d\n", phone, usb, n, ready)
		}
		wifiCameUp := haveWifi && !lastWifi && w
		lastWifi = w
		haveWifi = true
		lastPhone = phone
		lastUSB = usb
		havePresence = true
		pushLobbyInfo()
		if ipChanged {
			pushNotesList()
		}
		if wifiCameUp && syncEng.ready() {
			go func() { _, _ = syncEng.reconcileAll("wifi") }()
		}
	}
}

// wifiUp reports whether wlan0 is up (Lobby sync offline detection).
func wifiUp() bool {
	st, err := os.ReadFile("/sys/class/net/wlan0/operstate")
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(st)) == "up"
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
