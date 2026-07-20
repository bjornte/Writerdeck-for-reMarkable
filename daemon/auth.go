// Writerdeck-server — see main.go for overview.

package main

import (
	crand "crypto/rand"
	"crypto/subtle"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// --- Auth ---

var (
	authMu      sync.Mutex // guards authPIN, authToken, pinRequired
	authPIN     string     // PIN generated at startup; shown on the tablet Lobby
	authToken   string     // session token issued when the PIN is verified
	pinRequired bool       // false in no-PIN mode (checkAuth always passes)
	activeSess  *session   // non-nil only in supervisor (--editor) mode
)

// generatePIN mints a cryptographically random decimal PIN of the given length.
// length 0 returns "" (no-PIN mode). Length 4 produces a 4-digit PIN,
// length 6 a 6-digit PIN. Reduces in uint32 space before converting to int:
// int is 32-bit on the ARMv7 device and a raw Uint32 can exceed int32 max.
func generatePIN(length int) string {
	if length == 0 {
		return ""
	}
	var b [4]byte
	if _, err := crand.Read(b[:]); err != nil {
		if length == 4 {
			return "0000"
		}
		return "000000"
	}
	v := binary.BigEndian.Uint32(b[:])
	if length == 4 {
		return fmt.Sprintf("%04d", v%10000)
	}
	return fmt.Sprintf("%06d", v%1000000)
}

func generateToken() string {
	var b [16]byte
	if _, err := crand.Read(b[:]); err != nil {
		return "insecure-fallback"
	}
	return hex.EncodeToString(b[:])
}

var (
	lobbyIPMu         sync.Mutex
	lastPushedLobbyIP string
)
// checkAuth returns true if the request is authorized.
// Always returns true for OPTIONS (preflight) or when PIN auth is disabled.
// Loopback requests are trusted (tablet editor HTTP saves — slice 10).
// When PIN auth is enabled, checks the writerdeck_token session cookie.
func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodOptions {
		return true
	}
	if isLoopback(r) {
		return true
	}
	authMu.Lock()
	required := pinRequired
	tok := authToken
	authMu.Unlock()
	if !required {
		return true
	}
	cookie, err := r.Cookie("writerdeck_token")
	if err != nil || cookie.Value != tok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

// isLoopback reports whether r came from the tablet editor on localhost.
func isLoopback(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	return host == "127.0.0.1" || host == "::1"
}

// nextMorningCutoff returns the next local 04:00. If it is currently before
// today's 04:00 it returns today's 04:00; otherwise tomorrow's. Used as the
// auth-cookie expiry so a full day's writing never re-prompts for the PIN,
// while a fresh morning re-gates once. (A reboot also re-gates independently:
// authToken is regenerated per boot, so a stale cookie value stops matching --
// and a mid-day rmkbd restart likewise asks once, an accepted cost of not
// persisting the token to disk.)
func nextMorningCutoff(now time.Time) time.Time {
	cutoff := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
	if !now.Before(cutoff) {
		cutoff = cutoff.AddDate(0, 0, 1)
	}
	return cutoff
}

// --- PIN brute-force throttle ---
// A 6-digit PIN has 1,000,000 combinations; without throttling, someone on the
// LAN could exhaust it. We lock an IP out for pinLockout after pinMaxFails
// consecutive wrong guesses (HTTP 429 + Retry-After). State is in-memory and
// per-IP: a reboot clears it (and regenerates the PIN), and tracking by IP
// means a malicious actor cannot lock out the owner, who connects from a
// different address.
const (
	pinMaxFails = 5
	pinLockout  = 60 * time.Second
)

type pinAttempt struct {
	fails       int
	lockedUntil time.Time
}

var (
	pinMu       sync.Mutex
	pinAttempts = map[string]*pinAttempt{}
)

// clientIP returns the host portion of r.RemoteAddr (no port), or the raw
// value if it cannot be split.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// pinHandler handles POST /api/pin: validates the PIN and sets a session cookie.
func pinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		PIN string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Read auth state once under lock; use local copies for the rest of the handler.
	authMu.Lock()
	pin := authPIN
	required := pinRequired
	authMu.Unlock()

	// No PIN required: just issue a cookie (handles a client stuck on the PIN
	// screen when the owner switches to no-PIN mode) and return OK.
	if !required {
		authMu.Lock()
		tok := authToken
		authMu.Unlock()
		exp := nextMorningCutoff(time.Now())
		http.SetCookie(w, &http.Cookie{
			Name:     "writerdeck_token",
			Value:    tok,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Expires:  exp,
			MaxAge:   int(time.Until(exp).Seconds()),
		})
		w.WriteHeader(http.StatusOK)
		return
	}

	ip := clientIP(r)
	now := time.Now()
	pinMu.Lock()
	// Prune served-out lockouts so the map stays small and an expired lockout
	// resets that IP's counter automatically.
	for k, a := range pinAttempts {
		if a.fails >= pinMaxFails && now.After(a.lockedUntil) {
			delete(pinAttempts, k)
		}
	}
	if a := pinAttempts[ip]; a != nil && a.fails >= pinMaxFails && now.Before(a.lockedUntil) {
		retry := int(a.lockedUntil.Sub(now).Seconds()) + 1
		pinMu.Unlock()
		w.Header().Set("Retry-After", fmt.Sprintf("%d", retry))
		http.Error(w, "too many attempts", http.StatusTooManyRequests)
		return
	}
	// Constant-time compare so the response time does not leak how many leading
	// digits matched (the per-IP lockout above is the primary defense).
	if subtle.ConstantTimeCompare([]byte(req.PIN), []byte(pin)) != 1 {
		a := pinAttempts[ip]
		if a == nil {
			a = &pinAttempt{}
			pinAttempts[ip] = a
		}
		a.fails++
		locked := a.fails >= pinMaxFails
		if locked {
			a.lockedUntil = now.Add(pinLockout)
		}
		pinMu.Unlock()
		if locked {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(pinLockout.Seconds())))
			http.Error(w, "too many attempts", http.StatusTooManyRequests)
			return
		}
		http.Error(w, "wrong PIN", http.StatusUnauthorized)
		return
	}
	// Success: clear this IP's failure record.
	delete(pinAttempts, ip)
	pinMu.Unlock()

	// Re-read token under lock in case a concurrent PIN-length change regenerated
	// it during this request (so the issued cookie matches the current authToken).
	authMu.Lock()
	currentTok := authToken
	authMu.Unlock()

	// Expire the cookie at the next local 04:00 so a full day's writing never
	// re-prompts for the PIN, but a fresh morning (and any reboot) re-gates once.
	// Set both Expires and MaxAge: MaxAge wins in modern browsers, Expires is the
	// fallback for older ones -- both point at the same wall-clock moment.
	exp := nextMorningCutoff(time.Now())
	http.SetCookie(w, &http.Cookie{
		Name:     "writerdeck_token",
		Value:    currentTok,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  exp,
		MaxAge:   int(time.Until(exp).Seconds()),
	})
	w.WriteHeader(http.StatusOK)
}
