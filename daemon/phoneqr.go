// Writerdeck-server — see main.go for overview.

package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	qrcode "github.com/skip2/go-qrcode"
)

// ideBrowserUA is true for Cursor's embedded browser (and Electron shells).
// Those clients must not count as a keyboard path for the Lobby tip.
func ideBrowserUA(ua string) bool {
	return strings.Contains(ua, "Cursor/") || strings.Contains(ua, "Electron/")
}

const phoneQRPath = "/tmp/writerdeck-phone-qr.png"

var (
	phoneQRMu  sync.Mutex
	phoneQRURL string
)

// ensurePhoneQR writes a high-contrast PNG for the phone URL (e-ink Lobby tip).
// Returns the file path, or "" if generation failed / URL empty.
func ensurePhoneQR(url string) string {
	if url == "" {
		return ""
	}
	phoneQRMu.Lock()
	defer phoneQRMu.Unlock()
	if url == phoneQRURL {
		if st, err := os.Stat(phoneQRPath); err == nil && st.Size() > 0 {
			return phoneQRPath
		}
	}
	// Medium recovery, 256px — readable on e-ink when shown ~half screen width.
	if err := qrcode.WriteFile(url, qrcode.Medium, 256, phoneQRPath); err != nil {
		fmt.Fprintf(os.Stderr, "writerdeck-server: phone QR: %v\n", err)
		return ""
	}
	phoneQRURL = url
	return phoneQRPath
}

func phoneConnected() bool {
	wsClientsMu.Lock()
	defer wsClientsMu.Unlock()
	for c := range wsClients {
		if c.hello && !ideBrowserUA(c.ua) {
			return true
		}
	}
	return false
}

func usbKeyboardPresent() bool {
	return len(findKeyboardInputDevices()) > 0
}
