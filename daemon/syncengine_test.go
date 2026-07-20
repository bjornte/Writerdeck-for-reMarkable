package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

// startFakeGitHub serves GitHub Contents API responses for one file name.
// Content is returned base64-encoded, matching api.github.com.
func startFakeGitHub(t *testing.T, name, sha string, raw []byte) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || !strings.Contains(r.URL.Path, "/contents/"+name) {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(ghContentFile{
			Content: base64.StdEncoding.EncodeToString(raw),
			SHA:     sha,
			Name:    name,
		})
	}))
	t.Cleanup(srv.Close)
	return srv
}

func setupPullNoteTest(t *testing.T) (e *syncEngine, notesDir string) {
	t.Helper()
	dir := t.TempDir()
	oldNotes := notesDirPath
	oldAPI := ghAPIBase
	notesDirPath = dir
	settingsMu.Lock()
	oldSettings := curSettings
	oldSettingsPath := settingsFilePath
	curSettings = settingsData{
		SyncOn:   true,
		SyncRepo: "owner/notes",
		SyncMeta: map[string]noteSyncMeta{},
	}
	settingsFilePath = filepath.Join(dir, "settings.json")
	settingsMu.Unlock()
	currentNoteMu.Lock()
	oldOpen := currentNote
	currentNote = ""
	currentNoteMu.Unlock()
	e = &syncEngine{}
	e.setToken("test-token")
	t.Cleanup(func() {
		notesDirPath = oldNotes
		ghAPIBase = oldAPI
		settingsMu.Lock()
		curSettings = oldSettings
		settingsFilePath = oldSettingsPath
		settingsMu.Unlock()
		currentNoteMu.Lock()
		currentNote = oldOpen
		currentNoteMu.Unlock()
		e.clearToken()
	})
	return e, dir
}

func TestPullNoteGhostRestore(t *testing.T) {
	const sha = "blob-sha-unchanged"
	cases := []struct {
		name    string
		payload []byte
	}{
		{"ghost.md", []byte("# restored from github\n")},
		{"ghost.md.enc", []byte("WDENC1-fake-ciphertext")},
	}
	for _, tc := range cases {
		t.Run(tc.name+" missing locally", func(t *testing.T) {
			e, dir := setupPullNoteTest(t)
			srv := startFakeGitHub(t, tc.name, sha, tc.payload)
			ghAPIBase = srv.URL
			e.setMeta(tc.name, sha, strHash(string(tc.payload)))

			changed, err := e.pullNote(tc.name)
			if err != nil {
				t.Fatalf("pullNote: %v", err)
			}
			if !changed {
				t.Fatal("ghost note: want changed=true (restore)")
			}
			got, err := os.ReadFile(filepath.Join(dir, tc.name))
			if err != nil {
				t.Fatalf("restored file missing: %v", err)
			}
			if string(got) != string(tc.payload) {
				t.Fatalf("restored content = %q, want %q", got, tc.payload)
			}
		})

		t.Run(tc.name+" present locally", func(t *testing.T) {
			e, dir := setupPullNoteTest(t)
			local := append([]byte(nil), tc.payload...)
			local = append(local, []byte("-local")...)
			if err := os.WriteFile(filepath.Join(dir, tc.name), local, 0644); err != nil {
				t.Fatal(err)
			}
			srv := startFakeGitHub(t, tc.name, sha, tc.payload)
			ghAPIBase = srv.URL
			e.setMeta(tc.name, sha, strHash(string(local)))

			changed, err := e.pullNote(tc.name)
			if err != nil {
				t.Fatalf("pullNote: %v", err)
			}
			if changed {
				t.Fatal("present note with matching SHA: want changed=false")
			}
			got, err := os.ReadFile(filepath.Join(dir, tc.name))
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != string(local) {
				t.Fatalf("local file rewritten: got %q, want %q", got, local)
			}
		})
	}
}
