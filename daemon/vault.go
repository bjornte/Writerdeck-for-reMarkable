// Writerdeck-server — optional at-rest note encryption (ADR decisions.md §31).

package main

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/scrypt"
)

const (
	vaultMagic     = "WDENC1"
	vaultScryptN   = 32768
	vaultScryptR   = 8
	vaultScryptP   = 1
	vaultKeyLen    = 32
	vaultSaltLen   = 32
	vaultPINLen    = 6
	vaultVerifyMsg = "WDV1"

	vaultMaxFails = 5
	vaultLockout  = 60 * time.Second

	secretPinPath   = "secret/pin"
	secretVaultPath = "secret/vault"
)

var (
	vaultMu          sync.Mutex
	vaultDataKey     []byte // nil unless a note-editing session or one-shot crypto op is active
	vaultCurrentPIN  string // RAM only — for secret/pin sync after setup/change PIN
	vaultUnlockFails = map[string]*pinAttempt{}
)

type secretVaultJSON struct {
	Salt           string `json:"salt"`
	WrappedDataKey string `json:"wrappedDataKey"`
	Verifier       string `json:"verifier"`
}

func isEncryptedNoteName(name string) bool {
	return strings.HasSuffix(name, ".md.enc")
}

func validVaultPIN(pin string) bool {
	if len(pin) != vaultPINLen {
		return false
	}
	for _, c := range pin {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func vaultEnabled() bool {
	settingsMu.Lock()
	on := curSettings.EncryptionEnabled
	settingsMu.Unlock()
	return on
}

func vaultLocked() bool {
	vaultMu.Lock()
	defer vaultMu.Unlock()
	return vaultDataKey == nil
}

func deriveKEK(pin string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(pin), salt, vaultScryptN, vaultScryptR, vaultScryptP, vaultKeyLen)
}

func aesGCMSeal(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := crand.Read(nonce); err != nil {
		return nil, err
	}
	return append(nonce, gcm.Seal(nil, nonce, plaintext, nil)...), nil
}

func aesGCMOpen(key, sealed []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize()
	if len(sealed) < ns {
		return nil, errors.New("ciphertext too short")
	}
	return gcm.Open(nil, sealed[:ns], sealed[ns:], nil)
}

func loadVaultSalt() ([]byte, error) {
	settingsMu.Lock()
	b64 := curSettings.VaultSalt
	settingsMu.Unlock()
	if b64 == "" {
		return nil, errors.New("vault salt missing")
	}
	return base64.StdEncoding.DecodeString(b64)
}

func unwrapDataKeyFromSecret(pin, saltB64, wrappedB64, verifierB64 string) ([]byte, error) {
	if saltB64 == "" || wrappedB64 == "" || verifierB64 == "" {
		return nil, errors.New("incomplete vault secret")
	}
	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return nil, err
	}
	kek, err := deriveKEK(pin, salt)
	if err != nil {
		return nil, err
	}
	wrapped, err := base64.StdEncoding.DecodeString(wrappedB64)
	if err != nil {
		return nil, err
	}
	verifier, err := base64.StdEncoding.DecodeString(verifierB64)
	if err != nil {
		return nil, err
	}
	plainVerify, err := aesGCMOpen(kek, verifier)
	if err != nil || string(plainVerify) != vaultVerifyMsg {
		return nil, errors.New("wrong PIN")
	}
	dataKey, err := aesGCMOpen(kek, wrapped)
	if err != nil {
		return nil, errors.New("wrong PIN")
	}
	if len(dataKey) != vaultKeyLen {
		return nil, errors.New("invalid data key")
	}
	return dataKey, nil
}

func unwrapDataKey(pin string) ([]byte, error) {
	settingsMu.Lock()
	saltB64 := curSettings.VaultSalt
	wrappedB64 := curSettings.WrappedDataKey
	verifierB64 := curSettings.VaultVerifier
	settingsMu.Unlock()
	return unwrapDataKeyFromSecret(pin, saltB64, wrappedB64, verifierB64)
}

// listEncryptedNoteBasenames returns flat .md.enc names in notesDirPath.
func listEncryptedNoteBasenames() ([]string, error) {
	entries, err := os.ReadDir(notesDirPath)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if isEncryptedNoteName(e.Name()) {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

func hasUserEncryptedNotes() (bool, error) {
	names, err := listEncryptedNoteBasenames()
	if err != nil {
		return false, err
	}
	for _, n := range names {
		if !strings.HasPrefix(n, "z-test-") {
			return true, nil
		}
	}
	return false, nil
}

func vaultRewrapNote(name string, oldKey, newKey []byte) error {
	p := notesSafe(name)
	if p == "" || !isEncryptedNoteName(filepath.Base(p)) {
		return fmt.Errorf("invalid encrypted name")
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	plain, err := decryptNoteBytesWithKey(data, oldKey)
	if err != nil {
		return err
	}
	enc, err := encryptNoteBytesWithKey(plain, newKey)
	if err != nil {
		return err
	}
	return writeNoteFile(p, enc)
}

func vaultCheckUnlockThrottle() error {
	now := time.Now()
	vaultMu.Lock()
	defer vaultMu.Unlock()
	a := vaultUnlockFails["tablet"]
	if a != nil && a.fails >= vaultMaxFails && now.Before(a.lockedUntil) {
		return fmt.Errorf("too many attempts")
	}
	return nil
}

func vaultRecordUnlockFail() bool {
	now := time.Now()
	vaultMu.Lock()
	defer vaultMu.Unlock()
	a := vaultUnlockFails["tablet"]
	if a == nil {
		a = &pinAttempt{}
		vaultUnlockFails["tablet"] = a
	}
	a.fails++
	if a.fails >= vaultMaxFails {
		a.lockedUntil = now.Add(vaultLockout)
		return true
	}
	return false
}

func vaultClearUnlockFails() {
	vaultMu.Lock()
	delete(vaultUnlockFails, "tablet")
	vaultMu.Unlock()
}

// vaultVerifyPIN checks the vault PIN and stores the data key in RAM.
// keepSession=true keeps the key for note editing (cleared on Lobby return);
// keepSession=false keeps it only until the caller or a one-shot op clears it.
func vaultVerifyPIN(pin string, keepSession bool) error {
	if !vaultEnabled() {
		return errors.New("encryption not enabled")
	}
	if !validVaultPIN(pin) {
		return errors.New("invalid PIN format")
	}
	if err := vaultCheckUnlockThrottle(); err != nil {
		return err
	}
	dataKey, err := unwrapDataKey(pin)
	if err != nil {
		if vaultRecordUnlockFail() {
			return fmt.Errorf("too many attempts")
		}
		return err
	}
	vaultMu.Lock()
	vaultDataKey = dataKey
	vaultMu.Unlock()
	vaultClearUnlockFails()
	pushVaultPINOK()
	vaultNotifyPINGranted()
	if !keepSession {
		// One-shot ops (encrypt/decrypt/download) clear the key themselves.
	}
	return nil
}

func vaultClearSession() {
	vaultMu.Lock()
	vaultDataKey = nil
	vaultMu.Unlock()
}

func vaultClearSessionOnLobby() {
	vaultClearSession()
}

// vaultClearSessionIfIdle drops the data key when no encrypted note is open for editing.
func vaultClearSessionIfIdle() {
	currentNoteMu.Lock()
	open := currentNote
	currentNoteMu.Unlock()
	if open != "" && isEncryptedNoteName(open) {
		return
	}
	vaultClearSession()
}

func vaultSetupPIN(pin string) error {
	if !validVaultPIN(pin) {
		return errors.New("PIN must be 6 digits")
	}
	settingsMu.Lock()
	if curSettings.EncryptionEnabled {
		settingsMu.Unlock()
		return errors.New("encryption already enabled")
	}
	settingsMu.Unlock()

	salt := make([]byte, vaultSaltLen)
	if _, err := crand.Read(salt); err != nil {
		return err
	}
	dataKey := make([]byte, vaultKeyLen)
	if _, err := crand.Read(dataKey); err != nil {
		return err
	}
	kek, err := deriveKEK(pin, salt)
	if err != nil {
		return err
	}
	wrapped, err := aesGCMSeal(kek, dataKey)
	if err != nil {
		return err
	}
	verifier, err := aesGCMSeal(kek, []byte(vaultVerifyMsg))
	if err != nil {
		return err
	}

	settingsMu.Lock()
	curSettings.EncryptionEnabled = true
	curSettings.VaultSalt = base64.StdEncoding.EncodeToString(salt)
	curSettings.WrappedDataKey = base64.StdEncoding.EncodeToString(wrapped)
	curSettings.VaultVerifier = base64.StdEncoding.EncodeToString(verifier)
	saveSettingsLocked()
	settingsMu.Unlock()

	vaultMu.Lock()
	vaultCurrentPIN = pin
	vaultUnlockFails = map[string]*pinAttempt{}
	vaultMu.Unlock()

	pushLobbyInfo()
	syncEng.tryPushVaultSecrets()
	return nil
}

func vaultChangePIN(oldPIN, newPIN string) error {
	if !validVaultPIN(oldPIN) || !validVaultPIN(newPIN) {
		return errors.New("PIN must be 6 digits")
	}
	if !vaultEnabled() {
		return errors.New("encryption not enabled")
	}
	if err := vaultCheckUnlockThrottle(); err != nil {
		return err
	}
	dataKey, err := unwrapDataKey(oldPIN)
	if err != nil {
		if vaultRecordUnlockFail() {
			return fmt.Errorf("too many attempts")
		}
		return errors.New("wrong PIN")
	}
	vaultClearUnlockFails()

	salt, err := loadVaultSalt()
	if err != nil {
		return err
	}
	kek, err := deriveKEK(newPIN, salt)
	if err != nil {
		return err
	}
	wrapped, err := aesGCMSeal(kek, dataKey)
	if err != nil {
		return err
	}
	verifier, err := aesGCMSeal(kek, []byte(vaultVerifyMsg))
	if err != nil {
		return err
	}

	settingsMu.Lock()
	curSettings.WrappedDataKey = base64.StdEncoding.EncodeToString(wrapped)
	curSettings.VaultVerifier = base64.StdEncoding.EncodeToString(verifier)
	saveSettingsLocked()
	settingsMu.Unlock()

	vaultMu.Lock()
	vaultCurrentPIN = newPIN
	vaultMu.Unlock()

	pushLobbyInfo()
	syncEng.tryPushVaultSecrets()
	return nil
}

func vaultDisable(testHarness bool) error {
	if !vaultEnabled() {
		return nil
	}
	enc, err := listEncryptedNoteBasenames()
	if err != nil {
		return err
	}
	if len(enc) > 0 {
		if !testHarness {
			return fmt.Errorf("cannot disable vault: %d encrypted note(s) on disk", len(enc))
		}
		for _, n := range enc {
			if !strings.HasPrefix(n, "z-test-") {
				return fmt.Errorf("cannot disable vault: user encrypted note %s", n)
			}
		}
	}
	settingsMu.Lock()
	curSettings.EncryptionEnabled = false
	curSettings.VaultSalt = ""
	curSettings.VaultVerifier = ""
	curSettings.WrappedDataKey = ""
	saveSettingsLocked()
	settingsMu.Unlock()
	vaultMu.Lock()
	vaultDataKey = nil
	vaultCurrentPIN = ""
	vaultMu.Unlock()
	pushLobbyInfo()
	return nil
}

func encryptNoteBytesWithKey(plaintext, dataKey []byte) ([]byte, error) {
	if dataKey == nil {
		return nil, errors.New("vault PIN required")
	}
	sealed, err := aesGCMSeal(dataKey, plaintext)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 0, len(vaultMagic)+len(sealed))
	out = append(out, vaultMagic...)
	out = append(out, sealed...)
	return out, nil
}

func decryptNoteBytesWithKey(ciphertext, dataKey []byte) ([]byte, error) {
	if len(ciphertext) < len(vaultMagic) || string(ciphertext[:len(vaultMagic)]) != vaultMagic {
		return nil, errors.New("not an encrypted note")
	}
	if dataKey == nil {
		return nil, errors.New("vault PIN required")
	}
	return aesGCMOpen(dataKey, ciphertext[len(vaultMagic):])
}

func encryptNoteBytes(plaintext []byte) ([]byte, error) {
	vaultMu.Lock()
	dk := vaultDataKey
	vaultMu.Unlock()
	return encryptNoteBytesWithKey(plaintext, dk)
}

func decryptNoteBytes(ciphertext []byte) ([]byte, error) {
	vaultMu.Lock()
	dk := vaultDataKey
	vaultMu.Unlock()
	return decryptNoteBytesWithKey(ciphertext, dk)
}

func vaultEncryptNote(name string) error {
	if vaultLocked() {
		return errors.New("vault PIN required")
	}
	p := notesSafe(name)
	base := filepath.Base(p)
	if p == "" || isEncryptedNoteName(base) {
		return fmt.Errorf("invalid name")
	}
	if !strings.HasSuffix(base, ".md") {
		return fmt.Errorf("not a plain note")
	}
	plain, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	enc, err := encryptNoteBytes(plain)
	if err != nil {
		return err
	}
	encName := strings.TrimSuffix(base, ".md") + ".md.enc"
	encPath := filepath.Join(notesDirPath, encName)
	if noteNameConflict(encName, base) {
		return fmt.Errorf("already exists")
	}
	if err := writeNoteFile(encPath, enc); err != nil {
		return err
	}
	if err := os.Remove(p); err != nil {
		return err
	}
	currentNoteMu.Lock()
	wasOpen := currentNote != "" && filepath.Base(p) == currentNote
	if wasOpen {
		currentNote = encName
	}
	currentNoteMu.Unlock()
	if wasOpen && activeSess != nil && activeSess.isActive() && globalEC != nil {
		cmd, _ := json.Marshal(struct {
			T    string `json:"t"`
			C    string `json:"c"`
			Name string `json:"name"`
		}{"cmd", "noterenamed", encName})
		globalEC.write(cmd)
		broadcastOpenEdit(encName)
	}
	pushLobbyInfo()
	pushNotesList()
	notifyTabletCrud("renamenote", encName, base)
	vaultClearSessionIfIdle()
	return nil
}

func vaultDecryptNote(name string) error {
	if vaultLocked() {
		return errors.New("vault PIN required")
	}
	p := notesSafe(name)
	base := filepath.Base(p)
	if p == "" || !isEncryptedNoteName(base) {
		return fmt.Errorf("invalid encrypted name")
	}
	ciphertext, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	plain, err := decryptNoteBytes(ciphertext)
	if err != nil {
		return err
	}
	plainName := strings.TrimSuffix(base, ".md.enc") + ".md"
	plainPath := filepath.Join(notesDirPath, plainName)
	if noteNameConflict(plainName, base) {
		return fmt.Errorf("already exists")
	}
	if err := writeNoteFile(plainPath, plain); err != nil {
		return err
	}
	if err := os.Remove(p); err != nil {
		return err
	}
	currentNoteMu.Lock()
	wasOpen := currentNote != "" && filepath.Base(p) == currentNote
	if wasOpen {
		currentNote = plainName
	}
	currentNoteMu.Unlock()
	if wasOpen && activeSess != nil && activeSess.isActive() && globalEC != nil {
		cmd, _ := json.Marshal(struct {
			T    string `json:"t"`
			C    string `json:"c"`
			Name string `json:"name"`
		}{"cmd", "noterenamed", plainName})
		globalEC.write(cmd)
		broadcastOpenEdit(plainName)
	}
	pushLobbyInfo()
	pushNotesList()
	notifyTabletCrud("renamenote", plainName, base)
	vaultClearSessionIfIdle()
	return nil
}

// vaultSecretPinBytes returns secret/pin payload when encryption is on.
func vaultSecretPinBytes() ([]byte, bool) {
	if !vaultEnabled() {
		return nil, false
	}
	vaultMu.Lock()
	pin := vaultCurrentPIN
	vaultMu.Unlock()
	if pin == "" {
		return nil, false
	}
	return []byte(pin), true
}

func vaultSecretVaultJSON() ([]byte, bool) {
	settingsMu.Lock()
	if !curSettings.EncryptionEnabled {
		settingsMu.Unlock()
		return nil, false
	}
	payload := secretVaultJSON{
		Salt:           curSettings.VaultSalt,
		WrappedDataKey: curSettings.WrappedDataKey,
		Verifier:       curSettings.VaultVerifier,
	}
	settingsMu.Unlock()
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, false
	}
	return b, true
}

func vaultApplySecretVault(data []byte) error {
	var v secretVaultJSON
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v.Salt == "" || v.WrappedDataKey == "" || v.Verifier == "" {
		return errors.New("incomplete vault secret")
	}
	settingsMu.Lock()
	oldWrap := curSettings.WrappedDataKey
	settingsMu.Unlock()
	if oldWrap != "" && v.WrappedDataKey != oldWrap {
		userEnc, err := hasUserEncryptedNotes()
		if err != nil {
			return err
		}
		if userEnc {
			return errors.New("refusing vault sync: would orphan encrypted notes")
		}
	}
	settingsMu.Lock()
	curSettings.EncryptionEnabled = true
	curSettings.VaultSalt = v.Salt
	curSettings.WrappedDataKey = v.WrappedDataKey
	curSettings.VaultVerifier = v.Verifier
	saveSettingsLocked()
	settingsMu.Unlock()
	pushLobbyInfo()
	return nil
}

func vaultApplySecretPin(pin string) {
	if !validVaultPIN(pin) {
		return
	}
	vaultMu.Lock()
	if vaultCurrentPIN == "" {
		vaultCurrentPIN = pin
	}
	vaultMu.Unlock()
}

// --- vault PIN wait (phone download) ---

var (
	vaultWaitMu sync.Mutex
	vaultWaitCh []chan struct{}
)

func vaultRegisterWaiter() chan struct{} {
	ch := make(chan struct{}, 1)
	vaultWaitMu.Lock()
	vaultWaitCh = append(vaultWaitCh, ch)
	vaultWaitMu.Unlock()
	return ch
}

func vaultUnregisterWaiter(ch chan struct{}) {
	vaultWaitMu.Lock()
	defer vaultWaitMu.Unlock()
	for i, w := range vaultWaitCh {
		if w == ch {
			vaultWaitCh = append(vaultWaitCh[:i], vaultWaitCh[i+1:]...)
			return
		}
	}
}

func vaultNotifyPINGranted() {
	vaultWaitMu.Lock()
	waiters := vaultWaitCh
	vaultWaitCh = nil
	vaultWaitMu.Unlock()
	msg, _ := json.Marshal(struct {
		Type string `json:"type"`
	}{"vaultpingranted"})
	broadcast(msg)
	for _, ch := range waiters {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func waitVaultPINGranted(timeout time.Duration) bool {
	if !vaultLocked() {
		return true
	}
	ch := vaultRegisterWaiter()
	defer vaultUnregisterWaiter(ch)
	select {
	case <-ch:
		return !vaultLocked()
	case <-time.After(timeout):
		return false
	}
}

func pushVaultPINOK() {
	if globalEC == nil {
		return
	}
	globalEC.write([]byte(`{"t":"cmd","c":"vaultpinok"}` + "\n"))
}

// vaultOpErrMsg turns vault encrypt/decrypt errors into short tablet copy.
func vaultOpErrMsg(op string, err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "message authentication failed"),
		strings.Contains(msg, "ciphertext too short"),
		strings.Contains(msg, "not an encrypted note"):
		if op == "decrypt" {
			return "Cannot decrypt: wrong vault key or corrupted file. If vault was reset, recover from GitHub secret/vault history."
		}
		return "Encrypt failed: invalid document data."
	case strings.Contains(msg, "already exists"):
		return "A document with that name already exists."
	case strings.Contains(msg, "vault PIN required"):
		return "Enter PIN first."
	default:
		if op == "decrypt" {
			return "Decrypt failed: " + msg
		}
		return "Encrypt failed: " + msg
	}
}

func pushVaultOpFailed(msg string) {
	if globalEC == nil || msg == "" {
		return
	}
	cmd, _ := json.Marshal(struct {
		T   string `json:"t"`
		C   string `json:"c"`
		Msg string `json:"msg"`
	}{"cmd", "vaultopfailed", msg})
	globalEC.write(append(cmd, '\n'))
}

// vaultRewrapFromOldSecret re-encrypts notes that were sealed with a prior vault key.
func vaultRewrapFromOldSecret(oldVaultJSON, pin string, notes []string) error {
	if !validVaultPIN(pin) {
		return errors.New("PIN must be 6 digits")
	}
	if !vaultEnabled() {
		return errors.New("encryption not enabled")
	}
	var v secretVaultJSON
	if err := json.Unmarshal([]byte(oldVaultJSON), &v); err != nil {
		return err
	}
	oldKey, err := unwrapDataKeyFromSecret(pin, v.Salt, v.WrappedDataKey, v.Verifier)
	if err != nil {
		return fmt.Errorf("old vault: %w", err)
	}
	newKey, err := unwrapDataKey(pin)
	if err != nil {
		return fmt.Errorf("current vault: %w", err)
	}
	for _, name := range notes {
		if err := vaultRewrapNote(name, oldKey, newKey); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		syncEng.tryPushNote(name)
	}
	pushNotesList()
	return nil
}

func pushRequestVaultPIN(reason, name string) {
	if globalEC == nil {
		return
	}
	cmd, _ := json.Marshal(struct {
		T      string `json:"t"`
		C      string `json:"c"`
		Reason string `json:"reason"`
		Name   string `json:"name"`
	}{"cmd", "requestvaultpin", reason, name})
	globalEC.write(cmd)
}

// vaultPINMatches uses constant-time compare for test helpers.
func vaultPINMatches(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
