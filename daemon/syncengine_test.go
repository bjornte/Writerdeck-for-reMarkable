package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSyncReasonForcesRemote(t *testing.T) {
	if !syncReasonForcesRemote("manual") || !syncReasonForcesRemote("tablet") {
		t.Fatal("manual and tablet must force remote list")
	}
	for _, r := range []string{"boot", "home", "power", "token", "poll"} {
		if syncReasonForcesRemote(r) {
			t.Fatalf("%s must not force remote", r)
		}
	}
}

func TestContentNeedsPush(t *testing.T) {
	dir := t.TempDir()
	settingsMu.Lock()
	oldSettings := curSettings
	oldSettingsPath := settingsFilePath
	curSettings = settingsData{SyncMeta: map[string]noteSyncMeta{}}
	settingsFilePath = filepath.Join(dir, "settings.json")
	settingsMu.Unlock()
	t.Cleanup(func() {
		settingsMu.Lock()
		curSettings = oldSettings
		settingsFilePath = oldSettingsPath
		settingsMu.Unlock()
	})

	e := &syncEngine{}
	data := []byte("hello")
	if !e.contentNeedsPush("a.md", data) {
		t.Fatal("no meta: should need push")
	}
	e.setMeta("a.md", "sha1", strHash(string(data)))
	if e.contentNeedsPush("a.md", data) {
		t.Fatal("matching hash+sha: should skip")
	}
	if !e.contentNeedsPush("a.md", []byte("hello!")) {
		t.Fatal("changed content: should need push")
	}
	e.setMeta("b.md", "", strHash("x"))
	if !e.contentNeedsPush("b.md", []byte("x")) {
		t.Fatal("empty SHA: should need push")
	}
}

func TestLocalDirtyCount(t *testing.T) {
	dir := t.TempDir()
	oldNotes := notesDirPath
	notesDirPath = dir
	settingsMu.Lock()
	oldSettings := curSettings
	oldSettingsPath := settingsFilePath
	curSettings = settingsData{
		SyncOn:   true,
		SyncRepo: "o/r",
		SyncMeta: map[string]noteSyncMeta{},
	}
	settingsFilePath = filepath.Join(dir, "settings.json")
	settingsMu.Unlock()
	currentNoteMu.Lock()
	oldOpen := currentNote
	currentNote = ""
	currentNoteMu.Unlock()
	t.Cleanup(func() {
		notesDirPath = oldNotes
		settingsMu.Lock()
		curSettings = oldSettings
		settingsFilePath = oldSettingsPath
		settingsMu.Unlock()
		currentNoteMu.Lock()
		currentNote = oldOpen
		currentNoteMu.Unlock()
	})

	if err := os.WriteFile(filepath.Join(dir, "clean.md"), []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "dirty.md"), []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}

	e := &syncEngine{}
	e.setMeta("clean.md", "sha-clean", strHash("ok"))
	e.setMeta("dirty.md", "sha-dirty", strHash("old"))

	if n := e.localDirtyCount(false); n != 1 {
		t.Fatalf("dirty count = %d, want 1", n)
	}

	currentNoteMu.Lock()
	currentNote = "dirty.md"
	currentNoteMu.Unlock()
	if n := e.localDirtyCount(true); n != 0 {
		t.Fatalf("excluding open dirty note: got %d, want 0", n)
	}
	if n := e.localDirtyCount(false); n != 1 {
		t.Fatalf("including open: got %d, want 1", n)
	}
}
