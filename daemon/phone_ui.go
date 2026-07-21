// Writerdeck-server -- phone UI language packs (follow Lobby language).
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed phone-ui-i18n/en.json
var phoneUIEN []byte

//go:embed phone-ui-i18n/no.json
var phoneUINO []byte

//go:embed phone-ui-i18n/es.json
var phoneUIES []byte

//go:embed phone-ui-i18n/de.json
var phoneUIDE []byte

//go:embed phone-ui-i18n/fr.json
var phoneUIFR []byte

var phoneUIPacks = map[string][]byte{
	"en": phoneUIEN,
	"no": phoneUINO,
	"es": phoneUIES,
	"de": phoneUIDE,
	"fr": phoneUIFR,
}

func lobbyUIJSONPath() string {
	return filepath.Join(filepath.Dir(settingsFilePath), "lobby-ui.json")
}

func readLobbyLanguage() string {
	raw, err := os.ReadFile(lobbyUIJSONPath())
	if err != nil {
		return "en"
	}
	var doc struct {
		Language string `json:"language"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return "en"
	}
	lang := strings.TrimSpace(strings.ToLower(doc.Language))
	if _, ok := phoneUIPacks[lang]; !ok {
		return "en"
	}
	return lang
}

func loadPhoneUIStrings(lang string) (map[string]string, error) {
	raw, ok := phoneUIPacks[lang]
	if !ok {
		raw = phoneUIPacks["en"]
		lang = "en"
	}
	var stringsMap map[string]string
	if err := json.Unmarshal(raw, &stringsMap); err != nil {
		return nil, err
	}
	if lang != "en" {
		var en map[string]string
		if err := json.Unmarshal(phoneUIPacks["en"], &en); err == nil {
			for k, v := range en {
				if stringsMap[k] == "" {
					stringsMap[k] = v
				}
			}
		}
	}
	return stringsMap, nil
}

// phoneUIHandler serves GET /api/phone-ui: { language, strings }.
// Language follows lobby-ui.json (same as Lobby); packs ship in the binary.
func phoneUIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	lang := readLobbyLanguage()
	stringsMap, err := loadPhoneUIStrings(lang)
	if err != nil {
		http.Error(w, fmt.Sprintf("phone-ui pack: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
		"language": lang,
		"strings":  stringsMap,
	})
}
