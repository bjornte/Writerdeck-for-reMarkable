package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const syncPollInterval = 3 * time.Minute

// syncEngine runs GitHub reconcile on the tablet. Token lives in RAM only.
type syncEngine struct {
	tokenMu   sync.RWMutex
	token     string

	syncingMu sync.Mutex
	syncing   bool

	errMu     sync.RWMutex
	lastError string
}

var syncEng = &syncEngine{}

func (e *syncEngine) getToken() string {
	e.tokenMu.RLock()
	defer e.tokenMu.RUnlock()
	return e.token
}

func (e *syncEngine) setToken(tok string) {
	e.tokenMu.Lock()
	e.token = tok
	e.tokenMu.Unlock()
}

func (e *syncEngine) clearToken() {
	e.setToken("")
}

func (e *syncEngine) tokenConfigured() bool {
	return e.getToken() != ""
}

func (e *syncEngine) setLastError(msg string) {
	e.errMu.Lock()
	e.lastError = msg
	e.errMu.Unlock()
}

func (e *syncEngine) getLastError() string {
	e.errMu.RLock()
	defer e.errMu.RUnlock()
	return e.lastError
}

func (e *syncEngine) syncOn() bool {
	settingsMu.Lock()
	on := curSettings.SyncOn
	settingsMu.Unlock()
	return on
}

func (e *syncEngine) ready() bool {
	settingsMu.Lock()
	repo := curSettings.SyncRepo
	on := curSettings.SyncOn
	settingsMu.Unlock()
	return on && repo != "" && e.tokenConfigured()
}

func (e *syncEngine) isSyncing() bool {
	e.syncingMu.Lock()
	defer e.syncingMu.Unlock()
	return e.syncing
}

func (e *syncEngine) openNote() string {
	currentNoteMu.Lock()
	n := currentNote
	currentNoteMu.Unlock()
	return n
}

func readNoteContent(name string) (string, bool) {
	p := notesSafe(name)
	if p == "" {
		return "", false
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return "", false
	}
	return string(data), true
}

func writeNoteContentSync(name, content string) error {
	p := notesSafe(name)
	if p == "" {
		return fmt.Errorf("invalid name")
	}
	if rejectsHtmlNoteContent(content) {
		return fmt.Errorf("refusing HTML payload")
	}
	if err := writeNoteFile(p, []byte(content)); err != nil {
		return err
	}
	maybeBroadcastDiskChanged(filepath.Base(p))
	pushLobbyInfo()
	pushNotesList()
	return nil
}

func listLocalNoteNames() []string {
	entries, err := os.ReadDir(notesDirPath)
	if err != nil {
		return nil
	}
	var names []string
	for _, ent := range entries {
		if !ent.IsDir() && strings.HasSuffix(ent.Name(), ".md") {
			names = append(names, ent.Name())
		}
	}
	return names
}

func (e *syncEngine) pushNote(name string) error {
	if !e.ready() {
		return nil
	}
	if name == e.openNote() {
		return nil
	}
	content, ok := readNoteContent(name)
	if !ok {
		return nil
	}
	meta, hasMeta := e.getMeta(name)
	emptyHash := strHash("")
	if content == "" && hasMeta && meta.LocalHash != "" && meta.LocalHash != emptyHash {
		return nil // empty-push guard
	}
	sha := meta.SHA
	newSHA, status, err := e.ghPutFile(name, content, sha)
	if err != nil {
		if status == 409 || status == 422 {
			return e.handleClash(name, content)
		}
		return err
	}
	e.setMeta(name, newSHA, strHash(content))
	e.setLastError("")
	return nil
}

func (e *syncEngine) pullNote(name string) error {
	if !e.ready() || name == e.openNote() {
		return nil
	}
	gh, status, err := e.ghGetFile(name)
	if err != nil || gh == nil {
		if status == 401 || status == 403 {
			e.setLastError("GitHub token rejected")
		}
		return err
	}
	meta, _ := e.getMeta(name)
	if meta.SHA == gh.SHA {
		return nil
	}
	ghContent, err := ghDecodeContent(gh.Content)
	if err != nil {
		return err
	}
	if err := writeNoteContentSync(name, ghContent); err != nil {
		return err
	}
	e.setMeta(name, gh.SHA, strHash(ghContent))
	return nil
}

func (e *syncEngine) ghDeleteNote(name string) error {
	if !e.ready() {
		return nil
	}
	meta, ok := e.getMeta(name)
	if !ok || meta.SHA == "" {
		return nil
	}
	if err := e.ghDeleteFile(name, meta.SHA); err != nil {
		return err
	}
	e.removeMeta(name)
	return nil
}

func (e *syncEngine) handleClash(name, tabletContent string) error {
	gh, _, err := e.ghGetFile(name)
	if err != nil || gh == nil {
		return err
	}
	ghContent, err := ghDecodeContent(gh.Content)
	if err != nil {
		return err
	}
	if ghContent == tabletContent {
		e.setMeta(name, gh.SHA, strHash(tabletContent))
		return nil
	}
	if tabletContent == "" && ghContent != "" {
		if err := writeNoteContentSync(name, ghContent); err != nil {
			return err
		}
		e.setMeta(name, gh.SHA, strHash(ghContent))
		return nil
	}
	copyBase := strings.TrimSuffix(name, ".md") + " (tablet copy).md"
	_ = createNoteFile(copyBase, tabletContent)
	if err := writeNoteContentSync(name, ghContent); err != nil {
		return err
	}
	e.setMeta(name, gh.SHA, strHash(ghContent))
	fmt.Fprintf(os.Stderr, "writerdeck-server: sync clash on %s — tablet copy saved as %s\n", name, copyBase)
	broadcastClash(name, copyBase)
	return nil
}

func broadcastClash(noteName, copyName string) {
	msg, _ := json.Marshal(struct {
		Type     string `json:"type"`
		Note     string `json:"note"`
		CopyName string `json:"copyName"`
	}{"syncclash", noteName, copyName})
	broadcast(msg)
}

func (e *syncEngine) applyRemoteDelete(name string) error {
	if !e.ready() || name == e.openNote() {
		return nil
	}
	_, status, err := e.ghGetFile(name)
	if err != nil {
		return err
	}
	if status != 404 {
		return nil
	}
	if err := deleteNoteFile(name); err != nil && !os.IsNotExist(err) {
		return err
	}
	e.removeMeta(name)
	return nil
}

func (e *syncEngine) reconcileOne(name, remoteSHA string) error {
	hasRemote := remoteSHA != ""
	content, hasLocal := readNoteContent(name)
	if hasLocal && !hasRemote {
		meta, hasMeta := e.getMeta(name)
		if !hasMeta || meta.SHA == "" {
			return e.pushNote(name)
		}
		if meta.LocalHash != strHash(content) {
			return e.pushNote(name)
		}
		return e.applyRemoteDelete(name)
	}
	if !hasLocal && hasRemote {
		return e.pullNote(name)
	}
	if !hasLocal && !hasRemote {
		return nil
	}
	meta, _ := e.getMeta(name)
	remoteChanged := meta.SHA != remoteSHA
	localChanged := meta.LocalHash != strHash(content)
	if remoteChanged && localChanged {
		return e.handleClash(name, content)
	}
	if localChanged {
		emptyHash := strHash("")
		if content == "" && meta.LocalHash != "" && meta.LocalHash != emptyHash && hasRemote {
			return e.pullNote(name)
		}
		return e.pushNote(name)
	}
	if remoteChanged {
		return e.pullNote(name)
	}
	return nil
}

func (e *syncEngine) drainPendingSyncOps() {
	if !e.ready() {
		return
	}
	settingsMu.Lock()
	ops := append([]pendingSyncOp(nil), curSettings.PendingSync...)
	settingsMu.Unlock()
	for _, op := range ops {
		_ = e.applyTabletCrud(op.Op, op.Name, op.OldName)
	}
	if len(ops) > 0 {
		clearPendingSync()
	}
}

func (e *syncEngine) applyTabletCrud(op, name, oldName string) error {
	switch op {
	case "createnote":
		return e.pushNote(name)
	case "deletenote":
		return e.ghDeleteNote(name)
	case "renamenote":
		if err := e.ghDeleteNote(oldName); err != nil {
			return err
		}
		return e.pushNote(name)
	}
	return nil
}

func (e *syncEngine) markSyncComplete() {
	settingsMu.Lock()
	curSettings.LastSyncAt = time.Now().Unix()
	saveSettingsLocked()
	settingsMu.Unlock()
	signalSyncAck()
}

func (e *syncEngine) reconcileAll(reason string) (int, error) {
	if !e.ready() {
		return 0, nil
	}
	e.syncingMu.Lock()
	if e.syncing {
		e.syncingMu.Unlock()
		return 0, nil
	}
	e.syncing = true
	e.syncingMu.Unlock()

	defer func() {
		e.syncingMu.Lock()
		e.syncing = false
		e.syncingMu.Unlock()
		pushLobbyInfo()
	}()

	pushLobbyInfo()

	fmt.Fprintf(os.Stderr, "writerdeck-server: sync reconcile (%s)\n", reason)

	e.drainPendingSyncOps()

	open := e.openNote()
	if open != "" {
		fmt.Fprintf(os.Stderr, "writerdeck-server: sync skipped — edit lease on %s\n", open)
		return 0, nil
	}

	entries, status, err := e.ghListNotes()
	if err != nil {
		if status == 401 || status == 403 {
			e.setLastError("GitHub token rejected")
		}
		return 0, err
	}

	remoteMap := map[string]string{}
	for _, ent := range entries {
		if ent.Type == "file" && strings.HasSuffix(ent.Name, ".md") {
			remoteMap[ent.Name] = ent.SHA
		}
	}

	names := map[string]bool{}
	for n := range remoteMap {
		names[n] = true
	}
	for _, n := range listLocalNoteNames() {
		names[n] = true
	}

	count := 0
	for name := range names {
		if name == open {
			continue
		}
		if err := e.reconcileOne(name, remoteMap[name]); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: sync %s: %v\n", name, err)
		}
		count++
	}

	e.setLastError("")
	e.markSyncComplete()
	return count, nil
}

func (e *syncEngine) reconcileAllBlocking(reason string, timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		_, _ = e.reconcileAll(reason)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		fmt.Fprintln(os.Stderr, "writerdeck-server: sync reconcile timeout")
		signalSyncAck()
	}
}

func (e *syncEngine) trySyncAfterCrud(op, name, oldName string) {
	if !e.ready() {
		return
	}
	go func() {
		if err := e.applyTabletCrud(op, name, oldName); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: sync crud %s: %v\n", op, err)
			return
		}
		clearPendingSync()
	}()
}

func (e *syncEngine) tryPushNote(name string) {
	if !e.ready() || name == "" || name == e.openNote() {
		return
	}
	go func() {
		if err := e.pushNote(name); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: sync push %s: %v\n", name, err)
		}
	}()
}

func startSyncBackground() {
	go func() {
		time.Sleep(3 * time.Second)
		if syncEng.ready() {
			syncEng.reconcileAll("boot")
		}
	}()
	go func() {
		ticker := time.NewTicker(syncPollInterval)
		defer ticker.Stop()
		for range ticker.C {
			if syncEng.ready() {
				syncEng.reconcileAll("poll")
			}
		}
	}()
}

// mirrorPhoneDelete pairs a browser-initiated delete with GitHub.
func mirrorPhoneDelete(name string) { syncEng.trySyncAfterCrud("deletenote", name, "") }

// mirrorPhoneRename pairs a browser-initiated rename with GitHub.
func mirrorPhoneRename(oldName, newName string) {
	if !syncEng.ready() {
		return
	}
	go func() {
		_ = syncEng.ghDeleteNote(oldName)
		_ = syncEng.pushNote(newName)
	}()
}

// mirrorPhoneCreate pushes a newly created note to GitHub.
func mirrorPhoneCreate(name string) { syncEng.tryPushNote(name) }
