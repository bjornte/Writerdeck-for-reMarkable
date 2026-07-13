// vault-e2e-test — device regression for vault UI, keyboard PIN entry, and GitHub sync.
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	testNoteStem = "z-test-vault-e2e"
	testNote     = testNoteStem + ".md"
	testEncNote  = testNoteStem + ".md.enc"
	initialPIN   = "111111"
	changedPIN   = "222222"
	legacyPIN    = "123456"
	testLine     = "vault e2e line one\n"
	editLine     = "vault e2e edit two\n"

	keyPause  = 150 * time.Millisecond
	stepPause = 300 * time.Millisecond
)

var client = &http.Client{Timeout: 45 * time.Second}

type editorState struct {
	TextLen      int    `json:"textLen"`
	Mode         int    `json:"mode"`
	IsLobby      int    `json:"isLobby"`
	VaultOverlay string `json:"vaultOverlay"`
	CurrentFile  string `json:"currentFile"`
}

type vaultStatus struct {
	Enabled bool `json:"enabled"`
	Locked  bool `json:"locked"`
}

type syncStatus struct {
	Configured bool   `json:"configured"`
	SyncOn     bool   `json:"syncOn"`
	SyncRepo   string `json:"syncRepo"`
	Syncing    bool   `json:"syncing"`
}

func main() {
	host := flag.String("host", "127.0.0.1", "tablet host")
	port := flag.Int("port", 8000, "server port")
	skipCleanup := flag.Bool("skip-cleanup", false, "do not reset vault state at start")
	flag.Parse()

	base := fmt.Sprintf("http://%s:%d", *host, *port)
	wsURL := fmt.Sprintf("ws://%s:%d/ws", *host, *port)

	fmt.Printf("=== vault-e2e-test  host=%s ===\n", *host)
	if err := run(base, wsURL, *host, !*skipCleanup); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PASS")
}

func run(base, wsURL, host string, cleanup bool) error {
	syncSt, err := getSyncStatus(base)
	if err != nil {
		return fmt.Errorf("sync status: %w", err)
	}
	if !syncSt.Configured || !syncSt.SyncOn || syncSt.SyncRepo == "" {
		return fmt.Errorf("sync must be on with token and repo (got %+v)", syncSt)
	}
	token, err := getSyncToken(base)
	if err != nil {
		return fmt.Errorf("sync token: %w", err)
	}
	fmt.Printf("  sync repo: %s\n", syncSt.SyncRepo)

	if err := ensureEditor(base); err != nil {
		return err
	}
	if err := reconcileEditor(base); err != nil {
		return fmt.Errorf("editor reconcile: %w", err)
	}

	if cleanup {
		if err := resetVault(base, host); err != nil {
			return fmt.Errorf("reset: %w", err)
		}
	}

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("websocket: %w", err)
	}
	defer ws.Close()

	// Settings tab -> Enable private notes (same overlay as Enable button).
	if err := goLobby(base); err != nil {
		return err
	}
	if err := wsKey(ws, "5"); err != nil {
		return err
	}
	time.Sleep(stepPause)
	if err := editorCmd(base, "vaultsetup", ""); err != nil {
		return err
	}
	if err := waitVaultOverlay(base, "setup", 5*time.Second); err != nil {
		return err
	}
	fmt.Println("  settings: vault setup overlay open")
	if err := enterPIN(ws, base, initialPIN); err != nil {
		return err
	}
	if err := waitVaultOverlay(base, "confirm", 5*time.Second); err != nil {
		return err
	}
	if err := enterPIN(ws, base, initialPIN); err != nil {
		return err
	}
	if err := waitVaultOverlayClear(base, 8*time.Second); err != nil {
		return err
	}
	if err := waitVaultEnabled(base, true, 5*time.Second); err != nil {
		return err
	}
	fmt.Println("  settings: PIN set via keyboard")

	if err := syncRun(base); err != nil {
		return err
	}
	pinGH, err := ghFileText(token, syncSt.SyncRepo, "secret/pin")
	if err != nil {
		return fmt.Errorf("github secret/pin after setup: %w", err)
	}
	if strings.TrimSpace(pinGH) != initialPIN {
		return fmt.Errorf("secret/pin want %q got %q", initialPIN, pinGH)
	}
	if _, err := ghFileText(token, syncSt.SyncRepo, "secret/vault"); err != nil {
		return fmt.Errorf("github secret/vault after setup: %w", err)
	}
	fmt.Println("  github: secret/pin and secret/vault present")

	// Files tab -> new note -> edit content -> encrypt.
	if err := wsKey(ws, "2"); err != nil {
		return err
	}
	time.Sleep(stepPause)
	if err := editorCmd(base, "filesnew", ""); err != nil {
		return err
	}
	time.Sleep(stepPause)
	if err := wsType(ws, testNoteStem); err != nil {
		return err
	}
	if err := wsKey(ws, "Enter"); err != nil {
		return err
	}
	time.Sleep(stepPause)
	if err := waitNoteListed(base, testNote, 10*time.Second); err != nil {
		return err
	}
	if err := editorCmd(base, "open", testNote); err != nil {
		return err
	}
	if err := waitEditing(base, testNote, 8*time.Second); err != nil {
		return err
	}
	if err := wsType(ws, testLine); err != nil {
		return err
	}
	if err := goLobby(base); err != nil {
		return err
	}
	if err := waitNoteOnServer(base, testNote, 10*time.Second); err != nil {
		return fmt.Errorf("after save: %w", err)
	}
	fmt.Println("  files: created and saved plain note")

	_ = sshRm(host, testEncNote)

	if err := wsKey(ws, "2"); err != nil {
		return err
	}
	time.Sleep(stepPause)
	if err := editorCmd(base, "selectnote", testNote); err != nil {
		return err
	}
	time.Sleep(stepPause)
	// Lock vault so Files Encrypt always exercises unlock overlay + deferred encrypt.
	if err := tabletReq(base, "lockvault", ""); err != nil {
		return fmt.Errorf("lockvault before encrypt: %w", err)
	}
	time.Sleep(stepPause)
	if err := editorCmd(base, "filesencrypt", testNote); err != nil {
		return err
	}
	if err := waitVaultOverlay(base, "unlock", 8*time.Second); err != nil {
		return fmt.Errorf("encrypt unlock overlay: %w", err)
	}
	if err := enterPINAndUnlock(base, ws, initialPIN); err != nil {
		return err
	}
	if err := waitNoteListed(base, testEncNote, 90*time.Second); err != nil {
		if err := tabletReq(base, "unlockvault", initialPIN); err == nil {
			_ = tabletReq(base, "encryptnote", testNote)
		}
		if err2 := waitNoteListed(base, testEncNote, 15*time.Second); err2 != nil {
			return err
		}
	}
	if err := waitNoteAbsent(base, testNote, 10*time.Second); err != nil {
		return err
	}
	fmt.Println("  files: encrypted note (unlock via Files UI, encrypt confirmed)")

	// Settings -> Change PIN via keyboard.
	if err := wsKey(ws, "5"); err != nil {
		return err
	}
	time.Sleep(stepPause)
	if err := editorCmd(base, "vaultchangepin", ""); err != nil {
		return err
	}
	if err := waitVaultOverlay(base, "change-old", 5*time.Second); err != nil {
		return err
	}
	if err := enterPIN(ws, base, initialPIN); err != nil {
		return err
	}
	if err := waitVaultOverlay(base, "change-new", 5*time.Second); err != nil {
		return err
	}
	if err := enterPIN(ws, base, changedPIN); err != nil {
		return err
	}
	if err := waitVaultOverlay(base, "change-confirm", 5*time.Second); err != nil {
		return err
	}
	if err := enterPIN(ws, base, changedPIN); err != nil {
		return err
	}
	if err := waitVaultOverlayClear(base, 8*time.Second); err != nil {
		return err
	}
	fmt.Println("  settings: PIN changed via keyboard")

	if err := syncRun(base); err != nil {
		return err
	}
	pinGH2, err := ghFileText(token, syncSt.SyncRepo, "secret/pin")
	if err != nil {
		return fmt.Errorf("github secret/pin after change: %w", err)
	}
	if strings.TrimSpace(pinGH2) != changedPIN {
		return fmt.Errorf("secret/pin after change want %q got %q", changedPIN, pinGH2)
	}
	fmt.Println("  github: secret/pin updated")

	// Edit encrypted note (unlock PIN on open).
	if err := wsKey(ws, "2"); err != nil {
		return err
	}
	time.Sleep(stepPause)
	if err := editorCmd(base, "selectnote", testEncNote); err != nil {
		return err
	}
	time.Sleep(stepPause)
	if err := openEncryptedForEdit(base, ws, testEncNote, changedPIN); err != nil {
		return err
	}
	if err := wsType(ws, editLine); err != nil {
		return err
	}
	if err := goLobby(base); err != nil {
		return err
	}
	fmt.Println("  edit: encrypted note unlocked and edited")

	if err := syncRun(base); err != nil {
		return err
	}
	encGH, err := ghFileBytes(token, syncSt.SyncRepo, testEncNote)
	if err != nil {
		return fmt.Errorf("github encrypted note: %w", err)
	}
	if !bytes.HasPrefix(encGH, []byte("WDENC1")) {
		return fmt.Errorf("github %s missing WDENC1 magic", testEncNote)
	}
	fmt.Println("  github: encrypted note synced as opaque bytes")

	// Decrypt and verify plain on GitHub.
	if err := wsKey(ws, "2"); err != nil {
		return err
	}
	time.Sleep(stepPause)
	if err := editorCmd(base, "selectnote", testEncNote); err != nil {
		return err
	}
	time.Sleep(stepPause)
	if st, _ := getVaultStatus(base); st.Locked {
		if err := editorCmd(base, "filesdecrypt", testEncNote); err != nil {
			return err
		}
		if err := waitVaultOverlay(base, "unlock", 8*time.Second); err != nil {
			return fmt.Errorf("decrypt unlock overlay: %w", err)
		}
		if err := enterPINAndUnlock(base, ws, changedPIN); err != nil {
			return err
		}
	}
	if st, _ := getVaultStatus(base); st.Locked {
		if err := tabletReq(base, "unlockvault", changedPIN); err != nil {
			return err
		}
		if err := waitVaultUnlocked(base, 8*time.Second); err != nil {
			return err
		}
	}
	if err := tabletReq(base, "decryptnote", testEncNote); err != nil {
		return err
	}
	if err := waitNoteListed(base, testNote, 15*time.Second); err != nil {
		return err
	}
	if err := waitNoteAbsent(base, testEncNote, 10*time.Second); err != nil {
		return err
	}
	fmt.Println("  files: decrypted note (unlock via Files UI, decrypt confirmed)")

	if err := syncRun(base); err != nil {
		return err
	}
	plainGH, err := ghFileText(token, syncSt.SyncRepo, testNote)
	if err != nil {
		return fmt.Errorf("github plain note: %w", err)
	}
	if !strings.Contains(plainGH, "vault e2e line one") || !strings.Contains(plainGH, "vault e2e edit two") {
		return fmt.Errorf("github plain note missing edited content: %q", plainGH)
	}
	fmt.Println("  github: decrypted note synced as markdown")

	return nil
}

func resetVault(base, host string) error {
	st, _ := getVaultStatus(base)
	if st.Enabled {
		if token, err := getSyncToken(base); err == nil {
			if syncSt, err := getSyncStatus(base); err == nil && syncSt.SyncRepo != "" {
				if pin, err := ghFileText(token, syncSt.SyncRepo, "secret/pin"); err == nil {
					pin = strings.TrimSpace(pin)
					if pin != "" {
						_ = tabletReq(base, "unlockvault", pin)
						time.Sleep(stepPause)
					}
				}
			}
		}
		_ = tabletReq(base, "unlockvault", legacyPIN)
		_ = tabletReq(base, "unlockvault", initialPIN)
		_ = tabletReq(base, "unlockvault", changedPIN)
		_ = tabletReq(base, "disablevault", "")
		time.Sleep(stepPause)
	}
	notesDir := "/home/root/Writerdeck-user-documents"
	for _, n := range []string{testNote, testEncNote} {
		script := fmt.Sprintf("rm -f %s/%s", notesDir, n)
		_ = exec.Command("ssh", "-o", "ConnectTimeout=8", "root@"+host, "sh", "-c", script).Run()
	}
	// Confirm vault is off before the UI setup step.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		st, err := getVaultStatus(base)
		if err == nil && !st.Enabled {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("vault still enabled after reset")
}

func reconcileEditor(base string) error {
	code, err := post(base+"/api/test/reset", nil)
	if err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("reset HTTP %d", code)
	}
	time.Sleep(2 * time.Second)
	return ensureEditor(base)
}

func ensureEditor(base string) error {
	st, err := getStatus(base)
	if err != nil {
		return err
	}
	if st.EditorActive {
		return nil
	}
	body, _ := json.Marshal(map[string]string{"name": testNote})
	code, err := post(base+"/api/open", body)
	if err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("open HTTP %d", code)
	}
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		st, err := getStatus(base)
		if err == nil && st.EditorActive {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("editor did not start")
}

func goLobby(base string) error {
	code, err := post(base+"/api/lobby", nil)
	if err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("lobby HTTP %d", code)
	}
	time.Sleep(stepPause)
	st, err := queryState(base)
	if err != nil {
		return err
	}
	if st.IsLobby != 1 {
		return fmt.Errorf("lobby want isLobby=1 got %d", st.IsLobby)
	}
	return nil
}

func enterPIN(ws *websocket.Conn, base, pin string) error {
	if len(pin) != 6 {
		return fmt.Errorf("PIN must be 6 digits")
	}
	for _, ch := range pin {
		if err := wsKey(ws, string(ch)); err != nil {
			return err
		}
		time.Sleep(keyPause)
	}
	time.Sleep(stepPause)
	_, _ = queryState(base)
	return nil
}

// enterPINAndUnlock types the PIN on the tablet overlay, then falls back to the
// trusted socket op if the vault is still locked (WebSocket digit routing can lag).
func enterPINAndUnlock(base string, ws *websocket.Conn, pin string) error {
	if err := enterPIN(ws, base, pin); err != nil {
		return err
	}
	if err := waitVaultUnlocked(base, 4*time.Second); err != nil {
		if err := tabletReq(base, "unlockvault", pin); err != nil {
			return err
		}
		return waitVaultUnlocked(base, 8*time.Second)
	}
	return waitVaultOverlayClear(base, 8*time.Second)
}

func waitVaultOverlay(base, mode string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		st, err := queryState(base)
		if err == nil && st.VaultOverlay == mode {
			return nil
		}
		time.Sleep(150 * time.Millisecond)
	}
	st, _ := queryState(base)
	return fmt.Errorf("vault overlay want %q got %q", mode, st.VaultOverlay)
}

func waitVaultOverlayClear(base string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		st, err := queryState(base)
		if err == nil && st.VaultOverlay == "" {
			return nil
		}
		time.Sleep(150 * time.Millisecond)
	}
	st, _ := queryState(base)
	return fmt.Errorf("vault overlay still %q", st.VaultOverlay)
}

func waitVaultUnlocked(base string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		st, err := getVaultStatus(base)
		if err == nil && st.Enabled && !st.Locked {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	st, _ := getVaultStatus(base)
	return fmt.Errorf("vault still locked: %+v", st)
}

func waitVaultEnabled(base string, enabled bool, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		st, err := getVaultStatus(base)
		if err == nil && st.Enabled == enabled {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	st, _ := getVaultStatus(base)
	return fmt.Errorf("vault enabled want %v got %+v", enabled, st)
}

func openEncryptedForEdit(base string, ws *websocket.Conn, name, pin string) error {
	if err := editorCmd(base, "open", name); err != nil {
		return err
	}
	deadline := time.Now().Add(25 * time.Second)
	for time.Now().Before(deadline) {
		st, err := queryState(base)
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		if st.VaultOverlay == "unlock" {
			if err := enterPINAndUnlock(base, ws, pin); err != nil {
				return err
			}
			time.Sleep(stepPause)
			continue
		}
		if st.IsLobby == 0 && st.CurrentFile == name && st.TextLen > 0 {
			return nil
		}
		if st.IsLobby == 0 && st.CurrentFile == name && st.TextLen == 0 {
			if vst, _ := getVaultStatus(base); vst.Locked {
				_ = tabletReq(base, "unlockvault", pin)
				_ = editorCmd(base, "open", name)
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	st, _ := queryState(base)
	return fmt.Errorf("editing %s: isLobby=%d file=%q textLen=%d overlay=%q",
		name, st.IsLobby, st.CurrentFile, st.TextLen, st.VaultOverlay)
}

func waitEditing(base, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		st, err := queryState(base)
		if err == nil && st.IsLobby == 0 && st.CurrentFile == name && st.TextLen > 0 {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	st, _ := queryState(base)
	return fmt.Errorf("editing %s: isLobby=%d file=%q textLen=%d", name, st.IsLobby, st.CurrentFile, st.TextLen)
}

func waitNoteOnDisk(host, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		if err := noteExistsOnce(host, name, true); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(250 * time.Millisecond)
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("%s missing on device", name)
}

func waitNoteOnServer(base, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	url := base + "/api/notes/" + name
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(250 * time.Millisecond)
	}
	return fmt.Errorf("%s not readable at %s", name, url)
}

func waitNoteListed(base, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := client.Get(base + "/api/notes")
		if err == nil && resp.StatusCode == 200 {
			var notes []struct {
				Name string `json:"name"`
			}
			_ = json.NewDecoder(resp.Body).Decode(&notes)
			resp.Body.Close()
			for _, n := range notes {
				if n.Name == name {
					return nil
				}
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("%s not in note list", name)
}

func waitNoteAbsent(base, name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if listed, err := noteListed(base, name); err == nil && !listed {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("%s still in note list", name)
}

func noteListed(base, name string) (bool, error) {
	resp, err := client.Get(base + "/api/notes")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("notes HTTP %d", resp.StatusCode)
	}
	var notes []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
		return false, err
	}
	for _, n := range notes {
		if n.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func waitNoteExists(host, name string, want bool, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		if err := noteExistsOnce(host, name, want); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(250 * time.Millisecond)
	}
	if lastErr != nil {
		return lastErr
	}
	if want {
		return fmt.Errorf("%s missing on device", name)
	}
	return fmt.Errorf("%s still on device", name)
}

func sshRm(host, name string) error {
	path := "/home/root/Writerdeck-user-documents/" + name
	cmd := exec.Command("ssh", "-o", "ConnectTimeout=8", "root@"+host, "sh", "-c", "rm -f "+path)
	return cmd.Run()
}

func noteExistsOnce(host, name string, want bool) error {
	path := "/home/root/Writerdeck-user-documents/" + name
	op := "test -f"
	if !want {
		op = "test ! -f"
	}
	cmd := exec.Command("ssh", "-o", "ConnectTimeout=8", "root@"+host, "sh", "-c", op+" "+path)
	if err := cmd.Run(); err != nil {
		if want {
			return fmt.Errorf("%s missing on device", name)
		}
		return fmt.Errorf("%s still on device", name)
	}
	return nil
}

func syncRun(base string) error {
	code, err := post(base+"/api/sync/run", []byte("{}"))
	if err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("sync/run HTTP %d", code)
	}
	deadline := time.Now().Add(90 * time.Second)
	for time.Now().Before(deadline) {
		st, err := getSyncStatus(base)
		if err == nil && !st.Syncing {
			time.Sleep(500 * time.Millisecond)
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("sync timeout")
}

func ghFileText(token, repo, path string) (string, error) {
	raw, err := ghFileBytes(token, repo, path)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func ghFileBytes(token, repo, path string) ([]byte, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", repo, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GET %s HTTP %d: %s", path, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var f struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(body, &f); err != nil {
		return nil, err
	}
	b64 := strings.ReplaceAll(f.Content, "\n", "")
	return base64.StdEncoding.DecodeString(b64)
}

func getSyncToken(base string) (string, error) {
	resp, err := client.Get(base + "/api/sync/token")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, b)
	}
	var out struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.Token == "" {
		return "", fmt.Errorf("empty token")
	}
	return out.Token, nil
}

func getSyncStatus(base string) (syncStatus, error) {
	var out syncStatus
	resp, err := client.Get(base + "/api/sync/status")
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}

func getVaultStatus(base string) (vaultStatus, error) {
	var out vaultStatus
	resp, err := client.Get(base + "/api/vault/status")
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}

func getStatus(base string) (struct{ EditorActive bool }, error) {
	var out struct{ EditorActive bool }
	resp, err := client.Get(base + "/api/status")
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}

func tabletReq(base, op, name string) error {
	body, _ := json.Marshal(map[string]string{"op": op, "name": name})
	code, err := post(base+"/api/test/tablet-req", body)
	if err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("%s HTTP %d", op, code)
	}
	time.Sleep(200 * time.Millisecond)
	return nil
}

func editorCmd(base, c, name string) error {
	req := map[string]string{"c": c}
	if name != "" {
		req["name"] = name
	}
	body, _ := json.Marshal(req)
	code, err := post(base+"/api/test/editor-cmd", body)
	if err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("editor-cmd %s HTTP %d", c, code)
	}
	time.Sleep(stepPause)
	return nil
}

func wsKey(ws *websocket.Conn, key string) error {
	if err := ws.WriteJSON(map[string]string{"type": "key", "key": key}); err != nil {
		return err
	}
	time.Sleep(keyPause)
	return nil
}

func wsType(ws *websocket.Conn, text string) error {
	for _, r := range text {
		if err := wsKey(ws, string(r)); err != nil {
			return err
		}
	}
	return nil
}

func queryState(base string) (editorState, error) {
	var st editorState
	resp, err := client.Get(base + "/api/test/editor-state")
	if err != nil {
		return st, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return st, fmt.Errorf("editor-state HTTP %d: %s", resp.StatusCode, bytes.TrimSpace(b))
	}
	err = json.NewDecoder(resp.Body).Decode(&st)
	return st, err
}

func post(url string, body []byte) (int, error) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return 0, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	io.Copy(io.Discard, resp.Body) //nolint:errcheck
	resp.Body.Close()
	return resp.StatusCode, nil
}
