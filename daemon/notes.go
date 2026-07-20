// Writerdeck-server — see main.go for overview.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// --- Notes API ---

// notesDirPath is where .md notes are stored.
// Override with --notes-dir for local testing (default: /home/root/Writerdeck-user-documents).
var notesDirPath = "/home/root/Writerdeck-user-documents"

// noteInfo is the JSON shape returned by GET /api/notes.
type noteInfo struct {
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Modified  string `json:"modified"`
	Encrypted bool   `json:"encrypted,omitempty"`
	Locked    bool   `json:"locked,omitempty"`
}

// notesSafe validates a filename and returns its full path, or "".
// Rejects empty names, slashes, "..". Appends ".md" if no suffix given.
// Accepts flat basenames ending in .md or .md.enc.
func notesSafe(name string) string {
	if name == "" || strings.Contains(name, "/") || strings.Contains(name, "..") {
		return ""
	}
	if strings.HasSuffix(name, ".md.enc") {
		return filepath.Join(notesDirPath, name)
	}
	if strings.HasSuffix(name, ".md") {
		return filepath.Join(notesDirPath, name)
	}
	return filepath.Join(notesDirPath, name+".md")
}

func isNoteListName(name string) bool {
	if strings.HasSuffix(name, ".md.enc") {
		return true
	}
	if strings.HasSuffix(name, ".md") && !strings.HasSuffix(name, ".md.enc") {
		return true
	}
	return false
}

// noteTitleKey is the case-insensitive stem used for uniqueness and list order.
// "Doc.md", "doc.md.enc", and "DOC" all share the key "doc".
func noteTitleKey(name string) string {
	base := filepath.Base(name)
	base = strings.TrimSuffix(base, ".md.enc")
	base = strings.TrimSuffix(base, ".md")
	return strings.ToLower(base)
}

// noteNameConflict reports whether another note already uses the same title key
// (case-insensitive; plain and encrypted share a key). ignore is a full path or
// basename to skip (the file being renamed/encrypted).
func noteNameConflict(candidate, ignore string) bool {
	candKey := noteTitleKey(candidate)
	if candKey == "" {
		return false
	}
	ignoreBase := ""
	if ignore != "" {
		ignoreBase = filepath.Base(ignore)
	}
	entries, err := os.ReadDir(notesDirPath)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() || !isNoteListName(e.Name()) {
			continue
		}
		if ignoreBase != "" && e.Name() == ignoreBase {
			continue
		}
		if noteTitleKey(e.Name()) == candKey {
			return true
		}
	}
	return false
}

// sortNotesByTitle orders notes by case-insensitive title, then by name for ties.
func sortNotesByTitle(notes []noteInfo) {
	sort.Slice(notes, func(i, j int) bool {
		ki, kj := noteTitleKey(notes[i].Name), noteTitleKey(notes[j].Name)
		if ki != kj {
			return ki < kj
		}
		return strings.ToLower(notes[i].Name) < strings.ToLower(notes[j].Name)
	})
}

// rejectsHtmlNoteContent reports Qt qrichtext / HTML accidentally saved as .md.
func rejectsHtmlNoteContent(content string) bool {
	if len(content) < 15 {
		return false
	}
	head := strings.ToLower(content)
	if len(head) > 512 {
		head = head[:512]
	}
	return strings.HasPrefix(head, "<!doctype html") ||
		strings.HasPrefix(head, "<html") ||
		strings.Contains(content, `name="qrichtext"`)
}

// noteETag is a content-hash revision token for optimistic concurrency (RFC 7232).
func noteETag(content []byte) string {
	sum := sha256.Sum256(content)
	return `"` + hex.EncodeToString(sum[:8]) + `"`
}

// ifMatchOK reports whether the If-Match header allows writing over etag.
func ifMatchOK(ifMatch, etag string) bool {
	if ifMatch == "" {
		return false
	}
	if ifMatch == "*" {
		return true
	}
	for _, part := range strings.Split(ifMatch, ",") {
		part = strings.TrimSpace(part)
		if part == etag || part == "*" {
			return true
		}
	}
	return false
}

// writeNoteFile atomically writes note bytes (write-temp-then-rename, like settings.json).
func writeNoteFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, content, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// notesListHandler serves GET /api/notes (list) and POST /api/notes (create).
func notesListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	if r.URL.Path != "/api/notes" {
		http.NotFound(w, r)
		return
	}
	if !checkAuth(w, r) {
		return
	}
	switch r.Method {
	case http.MethodGet:
		notes := readAllNotes()
		if notes == nil {
			notes = []noteInfo{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes) //nolint:errcheck

	case http.MethodPost:
		// Cap the request body. The client also limits uploads to 1 MB, but a
		// client-side check is bypassable (e.g. a direct curl by an authed
		// caller), so this is the authoritative limit for ALL create paths
		// (New / New w/ paste / Upload). 2 MiB leaves headroom over a 1 MB text
		// file once it is wrapped + escaped in the JSON {name, content} envelope.
		// When exceeded, the read fails and the decode below returns 400.
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
		var req struct {
			Name    string `json:"name"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			http.Error(w, "bad request: need {name}", http.StatusBadRequest)
			return
		}
		p := notesSafe(req.Name)
		if p == "" {
			http.Error(w, "invalid name", http.StatusBadRequest)
			return
		}
		if noteNameConflict(p, "") {
			http.Error(w, "already exists", http.StatusConflict)
			return
		}
		if isEncryptedNoteName(filepath.Base(p)) {
			http.Error(w, "cannot create encrypted note from phone", http.StatusForbidden)
			return
		}
		content := req.Content
		if content == "" {
			content = "# " + strings.TrimSuffix(req.Name, ".md") + "\n"
		}
		if rejectsHtmlNoteContent(content) {
			http.Error(w, "refusing HTML/qrichtext payload", http.StatusUnsupportedMediaType)
			return
		}
		if err := writeNoteFile(p, []byte(content)); err != nil {
			http.Error(w, "write failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		pushLobbyInfo()
		pushNotesList()
		mirrorPhoneCreate(filepath.Base(p))

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// notesItemHandler serves GET /api/notes/{name} (read or download) and
// PUT /api/notes/{name} (content upsert for sync and tablet loopback).
func notesItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, If-Match")
		return
	}
	if !checkAuth(w, r) {
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/api/notes/")
	// Download variant: GET /api/notes/{name}/download
	// Strip the suffix BEFORE notesSafe() -- notesSafe rejects names containing '/'.
	download := strings.HasSuffix(name, "/download")
	if download {
		name = strings.TrimSuffix(name, "/download")
	}
	p := notesSafe(name)
	if p == "" {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		data, err := os.ReadFile(p)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		editorLocal := isLoopback(r)
		enc := isEncryptedNoteName(filepath.Base(p))
		if enc {
			if !editorLocal {
				if download {
					if vaultLocked() {
						pushRequestVaultPIN("download", filepath.Base(p))
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusLocked)
						w.Write([]byte(`{"error":"vault pin required","message":"Enter private PIN on tablet"}` + "\n")) //nolint:errcheck
						return
					}
					plain, err := decryptNoteBytes(data)
					if err != nil {
						http.Error(w, "cannot decrypt", http.StatusInternalServerError)
						return
					}
					data = plain
					vaultClearSessionIfIdle()
				} else {
					http.Error(w, "encrypted note", http.StatusForbidden)
					return
				}
			} else if vaultLocked() {
				http.Error(w, "vault PIN required", http.StatusLocked)
				return
			} else {
				plain, err := decryptNoteBytes(data)
				if err != nil {
					if editorLocal {
						pushVaultOpFailed(vaultOpErrMsg("decrypt", err))
					}
					http.Error(w, "cannot decrypt", http.StatusInternalServerError)
					return
				}
				data = plain
			}
		}
		if download {
			base := filepath.Base(p)
			if enc {
				base = strings.TrimSuffix(base, ".enc")
			}
			w.Header().Set("Content-Disposition", `attachment; filename="`+base+`"`)
			w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		} else {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}
		w.Header().Set("ETag", noteETag(data))
		w.Write(data) //nolint:errcheck

	case http.MethodPut:
		// Upsert: write or overwrite content. Used by the sync engine to apply a
		// version pulled from GitHub. Overwrite requires If-Match (content ETag).
		// 2 MiB limit matches POST /api/notes.
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
		var putReq struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&putReq); err != nil {
			http.Error(w, "bad request: need {content}", http.StatusBadRequest)
			return
		}
		if rejectsHtmlNoteContent(putReq.Content) {
			http.Error(w, "refusing HTML/qrichtext payload", http.StatusUnsupportedMediaType)
			return
		}
		editorLocal := isLoopback(r)
		enc := isEncryptedNoteName(filepath.Base(p))
		if enc && !editorLocal {
			http.Error(w, "encrypted notes are tablet-only", http.StatusForbidden)
			return
		}
		if enc && vaultLocked() {
			http.Error(w, "vault PIN required", http.StatusLocked)
			return
		}
		writeBytes := []byte(putReq.Content)
		if enc {
			var err error
			writeBytes, err = encryptNoteBytes(writeBytes)
			if err != nil {
				http.Error(w, "encrypt failed", http.StatusInternalServerError)
				return
			}
		}
		existing, err := os.ReadFile(p)
		if err == nil {
			if !editorLocal {
				etag := noteETag(existing)
				ifMatch := r.Header.Get("If-Match")
				if !ifMatchOK(ifMatch, etag) {
					w.Header().Set("ETag", etag)
					http.Error(w, "If-Match required or revision mismatch", http.StatusPreconditionFailed)
					return
				}
			}
		} else if !os.IsNotExist(err) {
			http.Error(w, "read failed", http.StatusInternalServerError)
			return
		}
		if err := writeNoteFile(p, writeBytes); err != nil {
			http.Error(w, "write failed", http.StatusInternalServerError)
			return
		}
		if !editorLocal {
			maybeBroadcastDiskChanged(filepath.Base(p))
		}
		etagBytes := writeBytes
		if enc && editorLocal {
			etagBytes = []byte(putReq.Content)
		}
		w.Header().Set("ETag", noteETag(etagBytes))
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
// readAllNotes returns metadata for every note file in the notes directory,
// ordered by case-insensitive title.
func readAllNotes() []noteInfo {
	entries, err := os.ReadDir(notesDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []noteInfo{}
		}
		return nil
	}
	var notes []noteInfo
	for _, e := range entries {
		if e.IsDir() || !isNoteListName(e.Name()) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		enc := isEncryptedNoteName(e.Name())
		notes = append(notes, noteInfo{
			Name:      e.Name(),
			Size:      info.Size(),
			Modified:  info.ModTime().Format(time.RFC3339),
			Encrypted: enc,
		})
	}
	if notes == nil {
		notes = []noteInfo{}
	}
	sortNotesByTitle(notes)
	return notes
}

// pushNotesList sends the full note list to the editor for the Lobby Files page.
func pushNotesList() {
	notes := readAllNotes()
	msg, _ := json.Marshal(struct {
		T                 string     `json:"t"`
		Items             []noteInfo `json:"items"`
		EncryptionEnabled bool       `json:"encryptionEnabled"`
	}{"notes", notes, vaultEnabled()})
	if globalEC != nil {
		globalEC.write(msg)
	}
}

// createNoteFile writes a new note; name is validated via notesSafe.
func createNoteFile(name, content string) error {
	p := notesSafe(name)
	if p == "" {
		return fmt.Errorf("invalid name")
	}
	if noteNameConflict(p, "") {
		return fmt.Errorf("already exists")
	}
	if content == "" {
		base := filepath.Base(p)
		title := strings.TrimSuffix(base, ".md.enc")
		title = strings.TrimSuffix(title, ".md")
		content = "# " + title + "\n"
	}
	if rejectsHtmlNoteContent(content) {
		return fmt.Errorf("refusing HTML/qrichtext payload")
	}
	if err := writeNoteFile(p, []byte(content)); err != nil {
		return err
	}
	pushLobbyInfo()
	pushNotesList()
	return nil
}

// deleteNoteFile removes a note and notifies the editor if it was open.
func deleteNoteFile(name string) error {
	p := notesSafe(name)
	if p == "" {
		return fmt.Errorf("invalid name")
	}
	if err := os.Remove(p); err != nil {
		return err
	}
	currentNoteMu.Lock()
	wasOpen := currentNote != "" && filepath.Base(p) == currentNote
	if wasOpen {
		currentNote = ""
	}
	currentNoteMu.Unlock()
	if wasOpen && activeSess != nil && activeSess.isActive() {
		if globalEC != nil {
			globalEC.write([]byte(`{"t":"cmd","c":"notedeleted"}`))
		}
		broadcast([]byte(`{"type":"exitedit"}`))
	}
	pushLobbyInfo()
	pushNotesList()
	return nil
}

// noteOpErrMsg maps create/rename disk errors to a short Lobby sentence.
func noteOpErrMsg(op string, err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if strings.Contains(msg, "already exists") || strings.Contains(msg, "name already taken") {
		return "A note with that name already exists."
	}
	if op == "rename" {
		return "Could not rename note."
	}
	return "Could not create note."
}

// renameNoteFile renames a note on disk and notifies the editor if it was open.
func renameNoteFile(oldName, newName string) error {
	oldP := notesSafe(oldName)
	newP := notesSafe(newName)
	if oldP == "" || newP == "" {
		return fmt.Errorf("invalid name")
	}
	if oldP == newP {
		return nil // same path — no-op (avoids false "already exists")
	}
	if noteNameConflict(newP, oldP) {
		return fmt.Errorf("name already taken")
	}
	if err := os.Rename(oldP, newP); err != nil {
		return err
	}
	newBase := filepath.Base(newP)
	currentNoteMu.Lock()
	wasOpen := currentNote != "" && filepath.Base(oldP) == currentNote
	if wasOpen {
		currentNote = newBase
	}
	currentNoteMu.Unlock()
	if wasOpen && activeSess != nil && activeSess.isActive() {
		if globalEC != nil {
			cmd, _ := json.Marshal(struct {
				T    string `json:"t"`
				C    string `json:"c"`
				Name string `json:"name"`
			}{"cmd", "noterenamed", newBase})
			globalEC.write(cmd)
		}
		broadcastOpenEdit(newBase)
	}
	pushLobbyInfo()
	pushNotesList()
	return nil
}
