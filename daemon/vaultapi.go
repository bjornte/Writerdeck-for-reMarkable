// Writerdeck-server — vault status API for phone download wait.

package main

import (
	"encoding/json"
	"net/http"
)

func vaultStatusHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct { //nolint:errcheck
		Enabled bool `json:"enabled"`
		Locked  bool `json:"locked"`
	}{vaultEnabled(), vaultLocked()})
}
