package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// productVersion is the build stamped into the binary (YYYY-MM-DD or YYYY-MM-DD.N).
// Set at link time: -ldflags "-X main.productVersion=$(tr -d '[:space:]' < VERSION)"
// Keep the repo-root VERSION file in sync -- that is what GitHub "latest" checks.
var productVersion = "unknown"

const githubVersionURL = "https://raw.githubusercontent.com/bjornte/Writerdeck-for-reMarkable/main/VERSION"

func localProductVersion() string {
	v := strings.TrimSpace(productVersion)
	if v == "" {
		return "unknown"
	}
	return v
}

func formatVersionMessage(local, latest string, fetchErr error) string {
	if fetchErr != nil || strings.TrimSpace(latest) == "" {
		return fmt.Sprintf("Writerdeck version %s (couldn't reach GitHub to check for updates)", local)
	}
	latest = strings.TrimSpace(latest)
	if latest == local {
		return fmt.Sprintf("Writerdeck version %s (latest)", local)
	}
	return fmt.Sprintf("Writerdeck version %s. Latest on GitHub is %s.", local, latest)
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

// GET /api/version -- local build stamp only (fast).
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

// GET /api/version/check -- compare to VERSION on GitHub main.
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
	local := localProductVersion()
	latest, err := fetchGitHubVersion()
	msg := formatVersionMessage(local, latest, err)
	out := map[string]interface{}{
		"version": local,
		"message": msg,
	}
	if err == nil {
		out["latest"] = latest
	} else {
		out["error"] = err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
