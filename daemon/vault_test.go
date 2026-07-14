package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func resetVaultTestState() {
	settingsMu.Lock()
	curSettings = settingsData{ReadFont: "Inter", PinDigits: "6"}
	settingsMu.Unlock()
	vaultMu.Lock()
	vaultDataKey = nil
	vaultCurrentPIN = ""
	vaultUnlockFails = map[string]*pinAttempt{}
	vaultMu.Unlock()
}

func TestVaultKDFAndWrap(t *testing.T) {
	resetVaultTestState()
	dir := t.TempDir()
	settingsFilePath = filepath.Join(dir, "settings.json")
	notesDirPath = dir

	if err := vaultSetupPIN("123456"); err != nil {
		t.Fatal(err)
	}
	if !vaultEnabled() {
		t.Fatal("expected encryption enabled")
	}
	if !vaultLocked() {
		t.Fatal("expected no session key after setup")
	}
	if err := vaultVerifyPIN("123456", true); err != nil {
		t.Fatalf("verify PIN: %v", err)
	}
	if vaultLocked() {
		t.Fatal("expected session key after verify")
	}
	vaultClearSession()
	if !vaultLocked() {
		t.Fatal("expected no session key after clear")
	}
	if err := vaultVerifyPIN("000000", true); err == nil {
		t.Fatal("wrong PIN should fail")
	}
}

func TestVaultEncryptRoundtrip(t *testing.T) {
	resetVaultTestState()
	dir := t.TempDir()
	settingsFilePath = filepath.Join(dir, "settings.json")
	notesDirPath = dir

	if err := vaultSetupPIN("654321"); err != nil {
		t.Fatal(err)
	}
	if err := vaultVerifyPIN("654321", true); err != nil {
		t.Fatal(err)
	}

	plain := []byte("# secret\n\nhello vault\n")
	enc, err := encryptNoteBytes(plain)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(enc, []byte(vaultMagic)) {
		t.Fatal("missing magic")
	}

	vaultClearSession()
	if _, err := decryptNoteBytes(enc); err == nil {
		t.Fatal("decrypt without session key should fail")
	}
	if err := vaultVerifyPIN("654321", true); err != nil {
		t.Fatal(err)
	}
	out, err := decryptNoteBytes(enc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(plain, out) {
		t.Fatalf("roundtrip mismatch: %q vs %q", out, plain)
	}
}

func TestVaultChangePINReWrap(t *testing.T) {
	resetVaultTestState()
	dir := t.TempDir()
	settingsFilePath = filepath.Join(dir, "settings.json")
	notesDirPath = dir

	if err := vaultSetupPIN("111111"); err != nil {
		t.Fatal(err)
	}
	if err := vaultVerifyPIN("111111", true); err != nil {
		t.Fatal(err)
	}
	plain := []byte("keep me")
	enc, err := encryptNoteBytes(plain)
	if err != nil {
		t.Fatal(err)
	}
	wrappedBefore := curSettings.WrappedDataKey
	vaultClearSession()

	if err := vaultChangePIN("111111", "222222"); err != nil {
		t.Fatal(err)
	}
	if curSettings.WrappedDataKey == wrappedBefore {
		t.Fatal("wrapped key should change after PIN change")
	}
	if err := vaultVerifyPIN("222222", true); err != nil {
		t.Fatal(err)
	}
	out, err := decryptNoteBytes(enc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(plain, out) {
		t.Fatal("ciphertext should still decrypt after PIN change")
	}
}

func TestNotesSafeEncrypted(t *testing.T) {
	dir := t.TempDir()
	notesDirPath = dir
	if p := notesSafe("diary.md.enc"); p == "" {
		t.Fatal("notesSafe should accept .md.enc")
	}
	if p := notesSafe("evil/foo.md.enc"); p != "" {
		t.Fatal("notesSafe should reject slashes")
	}
}

func TestVaultEncryptNoteFile(t *testing.T) {
	resetVaultTestState()
	dir := t.TempDir()
	settingsFilePath = filepath.Join(dir, "settings.json")
	notesDirPath = dir
	os.MkdirAll(dir, 0755)

	if err := vaultSetupPIN("333333"); err != nil {
		t.Fatal(err)
	}
	if err := vaultVerifyPIN("333333", true); err != nil {
		t.Fatal(err)
	}
	plainPath := filepath.Join(dir, "note.md")
	if err := os.WriteFile(plainPath, []byte("# hi"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := vaultEncryptNote("note.md"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(plainPath); !os.IsNotExist(err) {
		t.Fatal("plain file should be removed")
	}
	encPath := filepath.Join(dir, "note.md.enc")
	data, err := os.ReadFile(encPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := vaultVerifyPIN("333333", true); err != nil {
		t.Fatal(err)
	}
	dec, err := decryptNoteBytes(data)
	if err != nil {
		t.Fatal(err)
	}
	if string(dec) != "# hi" {
		t.Fatalf("got %q", dec)
	}
}

func TestVaultDisableRefusesUserNotes(t *testing.T) {
	resetVaultTestState()
	dir := t.TempDir()
	settingsFilePath = filepath.Join(dir, "settings.json")
	notesDirPath = dir

	if err := vaultSetupPIN("123456"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "til 1.md.enc"), []byte("WDENC1fake"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := vaultDisable(true); err == nil {
		t.Fatal("expected disable to refuse user encrypted note")
	}
	if err := os.Remove(filepath.Join(dir, "til 1.md.enc")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "z-test-vault.md.enc"), []byte("WDENC1fake"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := vaultDisable(true); err != nil {
		t.Fatalf("expected disable with z-test only: %v", err)
	}
}

func TestVaultOpErrMsg(t *testing.T) {
	if got := vaultOpErrMsg("decrypt", errors.New("cipher: message authentication failed")); got == "" {
		t.Fatal("expected decrypt auth failure message")
	}
	if got := vaultOpErrMsg("decrypt", errors.New("not an encrypted note")); !strings.Contains(got, "corrupted") {
		t.Fatalf("got %q", got)
	}
}
