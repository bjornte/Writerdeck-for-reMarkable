// Writerdeck-server — see main.go for overview.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// --- Settings API ---

// settingsFilePath is the JSON settings store on the device.
// Override with --settings-file for local dev (mirrors --notes-dir).
var settingsFilePath = "/home/root/.Writerdeck/settings.json"

// settingsData is the on-disk and in-memory settings schema.
type settingsData struct {
	ReadFont        string `json:"readFont"`
	PinDigits       string `json:"pinDigits"`       // "6", "4", or "none"; default "6"
	Rotation        int    `json:"rotation"`        // display rotation in degrees (0, 90, 180, 270)
	KeyboardLayout  string `json:"keyboardLayout"`  // USB evdev qmap id: "us", "no"; default "us"
	SyncOn     bool                       `json:"syncOn"`               // GitHub two-way sync enabled
	SyncRepo   string                     `json:"syncRepo"`          // "owner/repo" of the notes repo; token never stored here
	LastSyncAt int64                      `json:"lastSyncAt,omitempty"` // unix seconds of last reconcile
	SyncMeta   map[string]noteSyncMeta    `json:"syncMeta,omitempty"` // per-note GitHub SHA + local hash (non-secret)
	PendingSync []pendingSyncOp           `json:"pendingSync,omitempty"` // tablet CRUD awaiting sync (legacy drain)
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

// keyboardLayoutOption is the USB qmap allow-list (tablet Keyboard tab via setkeyboardlayout).
type keyboardLayoutOption struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

var keyboardLayoutRegistry = []keyboardLayoutOption{
	{ID: "us", Label: "US QWERTY"},
	{ID: "no", Label: "Norwegian"},
}

func normalizeKeyboardLayout(id string) string {
	for _, k := range keyboardLayoutRegistry {
		if k.ID == id {
			return id
		}
	}
	return "us"
}

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
		if s.KeyboardLayout == "" {
			s.KeyboardLayout = "us"
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
