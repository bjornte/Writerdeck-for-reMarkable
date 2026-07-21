package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// productVersion is the build stamped into the server binary (YYYY-MM-DD or YYYY-MM-DD.N).
// Set at link time from scripts/product-version.sh via -ldflags.
var productVersion = "unknown"

const githubVersionURL = "https://raw.githubusercontent.com/bjornte/Writerdeck-for-reMarkable/main/VERSION"

var productVersionRe = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})(?:\.(\d+))?$`)

func localProductVersion() string {
	v := strings.TrimSpace(productVersion)
	if v == "" {
		return "unknown"
	}
	return v
}

// parseProductVersion returns sort key (date, buildN). Bare dates use buildN 0.
// Unknown / empty / garbage sorts as older than any real stamp.
func parseProductVersion(s string) (date string, buildN int, ok bool) {
	s = strings.TrimSpace(s)
	if s == "" || s == "unknown" {
		return "", 0, false
	}
	m := productVersionRe.FindStringSubmatch(s)
	if m == nil {
		return "", 0, false
	}
	n := 0
	if m[2] != "" {
		n, _ = strconv.Atoi(m[2])
	}
	return m[1], n, true
}

// compareProductVersions returns -1 if a older, 0 if equal, 1 if a newer.
func compareProductVersions(a, b string) int {
	da, na, oka := parseProductVersion(a)
	db, nb, okb := parseProductVersion(b)
	if !oka && !okb {
		if a == b {
			return 0
		}
		return strings.Compare(a, b)
	}
	if !oka {
		return -1
	}
	if !okb {
		return 1
	}
	if da != db {
		if da < db {
			return -1
		}
		return 1
	}
	if na < nb {
		return -1
	}
	if na > nb {
		return 1
	}
	return 0
}

// combineProductVersion picks one stamp for About: both must agree.
// If they differ, returns the older stamp and mismatched=true.
func combineProductVersion(server, editor string) (product string, mismatched bool) {
	server = strings.TrimSpace(server)
	editor = strings.TrimSpace(editor)
	if server == "" {
		server = "unknown"
	}
	if editor == "" {
		editor = "unknown"
	}
	if server == editor {
		return server, false
	}
	if compareProductVersions(server, editor) <= 0 {
		return server, true
	}
	return editor, true
}

func formatVersionMessage(product, latest string, fetchErr error, mismatched bool) string {
	if mismatched {
		base := fmt.Sprintf("Writerdeck version %s (server and editor differ — update both)", product)
		if fetchErr != nil || strings.TrimSpace(latest) == "" {
			return base
		}
		latest = strings.TrimSpace(latest)
		if compareProductVersions(product, latest) < 0 {
			return fmt.Sprintf("%s. Latest on GitHub is %s.", base, latest)
		}
		return base
	}
	if fetchErr != nil || strings.TrimSpace(latest) == "" {
		return fmt.Sprintf("Writerdeck version %s (couldn't reach GitHub to check for updates)", product)
	}
	latest = strings.TrimSpace(latest)
	cmp := compareProductVersions(product, latest)
	if cmp == 0 {
		return fmt.Sprintf("Writerdeck version %s (latest)", product)
	}
	if cmp < 0 {
		return fmt.Sprintf("Writerdeck version %s. Latest on GitHub is %s.", product, latest)
	}
	return fmt.Sprintf("Writerdeck version %s (newer than GitHub %s)", product, latest)
}

// versionCheckStatus is a stable code for Lobby i18n (message stays English for phone/logs).
func versionCheckStatus(product, latest string, fetchErr error, mismatched bool) string {
	if mismatched {
		if fetchErr == nil && strings.TrimSpace(latest) != "" &&
			compareProductVersions(product, strings.TrimSpace(latest)) < 0 {
			return "mismatchLatest"
		}
		return "mismatch"
	}
	if fetchErr != nil || strings.TrimSpace(latest) == "" {
		return "offline"
	}
	cmp := compareProductVersions(product, strings.TrimSpace(latest))
	if cmp == 0 {
		return "latest"
	}
	if cmp < 0 {
		return "outdated"
	}
	return "ahead"
}

func fetchGitHubVersion() (string, error) {
	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest(http.MethodGet, githubVersionURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Writerdeck-server")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub VERSION HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64))
	if err != nil {
		return "", err
	}
	v := strings.TrimSpace(string(body))
	if v == "" {
		return "", fmt.Errorf("empty VERSION on GitHub")
	}
	// One line only; ignore trailing junk.
	if i := strings.IndexAny(v, "\r\n"); i >= 0 {
		v = strings.TrimSpace(v[:i])
	}
	return v, nil
}

// GET /api/version -- local server build stamp only (fast).
func versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !checkAuth(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"version": localProductVersion(),
	})
}

// GET /api/version/check -- one product stamp from server + editor, vs GitHub VERSION.
// Query: editor=YYYY-MM-DD (required for an honest combined result; omit => editor unknown).
func versionCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !checkAuth(w, r) {
		return
	}
	server := localProductVersion()
	editor := strings.TrimSpace(r.URL.Query().Get("editor"))
	if editor == "" {
		editor = "unknown"
	}
	product, mismatched := combineProductVersion(server, editor)
	latest, err := fetchGitHubVersion()
	msg := formatVersionMessage(product, latest, err, mismatched)
	status := versionCheckStatus(product, latest, err, mismatched)
	out := map[string]interface{}{
		"version":     product,
		"server":      server,
		"editor":      editor,
		"mismatched":  mismatched,
		"status":      status,
		"message":     msg,
	}
	if err == nil {
		out["latest"] = latest
	} else {
		out["error"] = err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
