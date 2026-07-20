package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// syncEngine runs GitHub reconcile on the tablet. Token lives in RAM only.
// Sync is change-driven (Home, power sleep, CRUD, manual, token, boot) — not polled.
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

func (e *syncEngine) needsBrowserToken() bool {
	settingsMu.Lock()
	on := curSettings.SyncOn
	repo := curSettings.SyncRepo
	settingsMu.Unlock()
	return on && repo != "" && !e.tokenConfigured()
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
	data, ok := readNoteBytes(name)
	if !ok {
		return "", false
	}
	if isEncryptedNoteName(name) {
		return string(data), true
	}
	return string(data), true
}

func writeNoteContentSync(name, content string) error {
	p := notesSafe(name)
	if p == "" {
		return fmt.Errorf("invalid name")
	}
	var data []byte
	if isEncryptedNoteName(name) {
		data = []byte(content)
	} else {
		if rejectsHtmlNoteContent(content) {
			return fmt.Errorf("refusing HTML payload")
		}
		data = []byte(content)
	}
	if err := writeNoteFile(p, data); err != nil {
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
		if !ent.IsDir() && isNoteListName(ent.Name()) {
			names = append(names, ent.Name())
		}
	}
	return names
}

func readNoteBytes(name string) ([]byte, bool) {
	p := notesSafe(name)
	if p == "" {
		return nil, false
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, false
	}
	return data, true
}

func noteContentHash(name string, data []byte) string {
	return strHash(string(data))
}

// syncReasonForcesRemote is true for explicit Sync (phone or Lobby). Those
// always list GitHub so remote-only edits can still be pulled when the tablet is clean.
func syncReasonForcesRemote(reason string) bool {
	return reason == "manual" || reason == "tablet"
}

// contentNeedsPush reports whether local bytes differ from the last synced fingerprint.
// No SHA or no meta means never pushed — treat as dirty.
func (e *syncEngine) contentNeedsPush(name string, data []byte) bool {
	meta, hasMeta := e.getMeta(name)
	if !hasMeta || meta.SHA == "" {
		return true
	}
	return meta.LocalHash != strHash(string(data))
}

// localDirtyCount counts notes and vault secrets that need a push, plus pending CRUD.
// When excludeOpen is true, the note open in the editor is omitted (reconcile skips it).
func (e *syncEngine) localDirtyCount(excludeOpen bool) int {
	open := ""
	if excludeOpen {
		open = e.openNote()
	}
	n := 0
	for _, name := range listLocalNoteNames() {
		if name == open {
			continue
		}
		data, ok := readNoteBytes(name)
		if !ok {
			continue
		}
		if e.contentNeedsPush(name, data) {
			n++
		}
	}
	if vaultEnabled() {
		if pin, ok := vaultSecretPinBytes(); ok && e.contentNeedsPush(secretPinPath, pin) {
			n++
		}
		if vaultJSON, ok := vaultSecretVaultJSON(); ok && e.contentNeedsPush(secretVaultPath, vaultJSON) {
			n++
		}
	}
	settingsMu.Lock()
	pending := len(curSettings.PendingSync)
	settingsMu.Unlock()
	return n + pending
}

func (e *syncEngine) handleClashBytes(name string, tabletData []byte) error {
	gh, _, err := e.ghGetFile(name)
	if err != nil || gh == nil {
		return err
	}
	ghData, err := ghDecodeBytes(gh.Content)
	if err != nil {
		return err
	}
	if bytes.Equal(ghData, tabletData) {
		e.setMeta(name, gh.SHA, strHash(string(tabletData)))
		return nil
	}
	if len(tabletData) == 0 && len(ghData) > 0 {
		p := notesSafe(name)
		if p == "" {
			return fmt.Errorf("invalid name")
		}
		if err := writeNoteFile(p, ghData); err != nil {
			return err
		}
		maybeBroadcastDiskChanged(name)
		pushLobbyInfo()
		pushNotesList()
		e.setMeta(name, gh.SHA, strHash(string(ghData)))
		return nil
	}
	copyBase := strings.TrimSuffix(name, ".md.enc") + " (tablet copy).md.enc"
	if !isEncryptedNoteName(name) {
		copyBase = strings.TrimSuffix(name, ".md") + " (tablet copy).md"
	}
	_ = writeNoteFile(notesSafe(copyBase), tabletData)
	if err := writeNoteFile(notesSafe(name), ghData); err != nil {
		return err
	}
	e.setMeta(name, gh.SHA, strHash(string(ghData)))
	fmt.Fprintf(os.Stderr, "writerdeck-server: sync clash on %s — tablet copy saved as %s\n", name, copyBase)
	broadcastClash(name, copyBase)
	return nil
}

func (e *syncEngine) pushNote(name string) (bool, error) {
	if !e.ready() {
		return false, nil
	}
	if name == e.openNote() {
		return false, nil
	}
	if isEncryptedNoteName(name) {
		data, ok := readNoteBytes(name)
		if !ok {
			return false, nil
		}
		// Encrypted notes are opaque bytes and must carry the WDENC1 header.
		// If an older buggy run wrote raw bytes without the header, re-wrap them
		// before pushing so GitHub always stores valid WDENC1 blobs.
		if !bytes.HasPrefix(data, []byte(vaultMagic)) {
			if vaultLocked() {
				return false, fmt.Errorf("cannot rewrap %s without vault PIN", name)
			}
			rewrapped, err := encryptNoteBytes(data)
			if err != nil {
				return false, err
			}
			if p := notesSafe(name); p != "" {
				_ = writeNoteFile(p, rewrapped)
			}
			data = rewrapped
		}
		meta, hasMeta := e.getMeta(name)
		emptyHash := strHash("")
		h := strHash(string(data))
		if len(data) == 0 && hasMeta && meta.LocalHash != "" && meta.LocalHash != emptyHash {
			return false, nil
		}
		if hasMeta && meta.SHA != "" && meta.LocalHash == h {
			return false, nil
		}
		sha := meta.SHA
		newSHA, status, err := e.ghPutBytes(name, data, sha)
		if err != nil {
			if status == 409 || status == 422 {
				return true, e.handleClashBytes(name, data)
			}
			return false, err
		}
		e.setMeta(name, newSHA, h)
		e.setLastError("")
		return true, nil
	}
	content, ok := readNoteContent(name)
	if !ok {
		return false, nil
	}
	meta, hasMeta := e.getMeta(name)
	emptyHash := strHash("")
	if content == "" && hasMeta && meta.LocalHash != "" && meta.LocalHash != emptyHash {
		return false, nil // empty-push guard
	}
	h := strHash(content)
	if hasMeta && meta.SHA != "" && meta.LocalHash == h {
		return false, nil
	}
	sha := meta.SHA
	newSHA, status, err := e.ghPutFile(name, content, sha)
	if err != nil {
		if status == 409 || status == 422 {
			return true, e.handleClash(name, content)
		}
		return false, err
	}
	e.setMeta(name, newSHA, h)
	e.setLastError("")
	return true, nil
}

func (e *syncEngine) pullNote(name string) (bool, error) {
	if !e.ready() || name == e.openNote() {
		return false, nil
	}
	gh, status, err := e.ghGetFile(name)
	if err != nil || gh == nil {
		if status == 401 || status == 403 {
			e.setLastError("GitHub token rejected")
		}
		return false, err
	}
	meta, _ := e.getMeta(name)
	if meta.SHA == gh.SHA {
		// Fast-path only when the tablet still actually has the file. If the
		// note is missing locally yet meta.SHA matches (a "ghost note" left by a
		// filesystem-level deletion that bypassed the app's delete path), the
		// stale SHA would otherwise suppress the re-pull forever, since the
		// GitHub blob never changes. Fall through and rewrite it so remote-only
		// notes are restored to the tablet.
		if _, ok := readNoteBytes(name); ok {
			return false, nil
		}
	}
	if isEncryptedNoteName(name) {
		raw, err := ghDecodeBytes(gh.Content)
		if err != nil {
			return false, err
		}
		p := notesSafe(name)
		if p == "" {
			return false, fmt.Errorf("invalid name")
		}
		if err := writeNoteFile(p, raw); err != nil {
			return false, err
		}
		maybeBroadcastDiskChanged(name)
		pushLobbyInfo()
		pushNotesList()
		e.setMeta(name, gh.SHA, strHash(string(raw)))
		return true, nil
	}
	ghContent, err := ghDecodeContent(gh.Content)
	if err != nil {
		return false, err
	}
	if err := writeNoteContentSync(name, ghContent); err != nil {
		return false, err
	}
	e.setMeta(name, gh.SHA, strHash(ghContent))
	return true, nil
}

func (e *syncEngine) ghDeleteNote(name string) (bool, error) {
	if !e.ready() {
		return false, nil
	}
	meta, ok := e.getMeta(name)
	if !ok || meta.SHA == "" {
		return false, nil
	}
	if err := e.ghDeleteFile(name, meta.SHA); err != nil {
		return false, err
	}
	e.removeMeta(name)
	return true, nil
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

func (e *syncEngine) applyRemoteDelete(name string) (bool, error) {
	if !e.ready() || name == e.openNote() {
		return false, nil
	}
	_, status, err := e.ghGetFile(name)
	if err != nil {
		return false, err
	}
	if status != 404 {
		return false, nil
	}
	if err := deleteNoteFile(name); err != nil && !os.IsNotExist(err) {
		return false, err
	}
	e.removeMeta(name)
	return true, nil
}

func (e *syncEngine) reconcileOne(name, remoteSHA string) (string, error) {
	hasRemote := remoteSHA != ""
	content, hasLocal := readNoteContent(name)
	if hasLocal && !hasRemote {
		meta, hasMeta := e.getMeta(name)
		if !hasMeta || meta.SHA == "" {
			ok, err := e.pushNote(name)
			if ok {
				return "pushed " + name, err
			}
			return "", err
		}
		if meta.LocalHash != strHash(content) {
			ok, err := e.pushNote(name)
			if ok {
				return "pushed " + name, err
			}
			return "", err
		}
		ok, err := e.applyRemoteDelete(name)
		if ok {
			return "removed local " + name + " (gone on GitHub)", err
		}
		return "", err
	}
	if !hasLocal && hasRemote {
		ok, err := e.pullNote(name)
		if ok {
			return "pulled " + name, err
		}
		return "", err
	}
	if !hasLocal && !hasRemote {
		return "", nil
	}
	meta, _ := e.getMeta(name)
	remoteChanged := meta.SHA != remoteSHA
	localChanged := meta.LocalHash != strHash(content)
	if remoteChanged && localChanged {
		if err := e.handleClash(name, content); err != nil {
			return "", err
		}
		return "clash on " + name, nil
	}
	if localChanged {
		emptyHash := strHash("")
		if content == "" && meta.LocalHash != "" && meta.LocalHash != emptyHash && hasRemote {
			ok, err := e.pullNote(name)
			if ok {
				return "pulled " + name, err
			}
			return "", err
		}
		ok, err := e.pushNote(name)
		if ok {
			return "pushed " + name, err
		}
		return "", err
	}
	if remoteChanged {
		ok, err := e.pullNote(name)
		if ok {
			return "pulled " + name, err
		}
		return "", err
	}
	return "", nil
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
		_, err := e.pushNote(name)
		return err
	case "deletenote":
		_, err := e.ghDeleteNote(name)
		return err
	case "renamenote":
		if _, err := e.ghDeleteNote(oldName); err != nil {
			return err
		}
		_, err := e.pushNote(name)
		return err
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

	// Auto triggers skip GitHub when nothing local needs a push. Manual/Lobby Sync
	// always lists remote so phone-side edits still land.
	if !syncReasonForcesRemote(reason) {
		dirty := e.localDirtyCount(true)
		if dirty == 0 {
			logSyncIdle(reason)
			signalSyncAck()
			return 0, nil
		}
	}

	pushLobbyInfo()

	e.drainPendingSyncOps()

	open := e.openNote()
	if open != "" {
		fmt.Fprintf(os.Stderr, "writerdeck-server: sync: leaving open note %s alone\n", open)
	}

	entries, status, err := e.ghListNotes()
	if err != nil {
		if status == 401 || status == 403 {
			e.setLastError("GitHub token rejected")
		} else if !wifiUp() {
			e.setLastError("No Wi-Fi - cannot reach GitHub")
		} else {
			e.setLastError("Could not reach GitHub")
		}
		return 0, err
	}

	remoteMap := map[string]string{}
	for _, ent := range entries {
		if ent.Type != "file" {
			continue
		}
		if strings.HasSuffix(ent.Name, ".md.enc") || (strings.HasSuffix(ent.Name, ".md") && !strings.HasSuffix(ent.Name, ".md.enc")) {
			remoteMap[ent.Name] = ent.SHA
		}
	}

	var details []string
	details = append(details, e.syncVaultSecrets(remoteMap)...)

	names := map[string]bool{}
	for n := range remoteMap {
		names[n] = true
	}
	for _, n := range listLocalNoteNames() {
		names[n] = true
	}

	for name := range names {
		if name == open {
			continue
		}
		detail, err := e.reconcileOne(name, remoteMap[name])
		if err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: sync: %s failed: %v\n", name, err)
			continue
		}
		if detail != "" {
			details = append(details, detail)
		}
	}

	e.setLastError("")
	e.markSyncComplete()
	logSyncChanged(reason, details)
	return len(details), nil
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
		if _, err := e.pushNote(name); err != nil {
			fmt.Fprintf(os.Stderr, "writerdeck-server: sync push %s: %v\n", name, err)
		}
	}()
}

func startSyncBackground() {
	// Boot reconcile only — no periodic poll. Event triggers: Home, power sleep,
	// CRUD, token verify, and explicit Sync (phone / Lobby).
	go func() {
		time.Sleep(3 * time.Second)
		if syncEng.ready() {
			syncEng.reconcileAll("boot")
		}
	}()
}

// mirrorPhoneCreate pushes a newly uploaded note to GitHub.
func mirrorPhoneCreate(name string) { syncEng.tryPushNote(name) }

func (e *syncEngine) tryPushVaultSecrets() {
	if !e.ready() || !vaultEnabled() {
		return
	}
	go func() { e.pushVaultSecrets() }()
}

func (e *syncEngine) pushVaultSecrets() []string {
	if !e.ready() || !vaultEnabled() {
		return nil
	}
	var details []string
	if pin, ok := vaultSecretPinBytes(); ok {
		if e.contentNeedsPush(secretPinPath, pin) {
			meta, _ := e.getMeta(secretPinPath)
			sha, _, err := e.ghPutBytes(secretPinPath, pin, meta.SHA)
			if err == nil {
				e.setMeta(secretPinPath, sha, strHash(string(pin)))
				details = append(details, "pushed "+secretPinPath)
			} else {
				fmt.Fprintf(os.Stderr, "writerdeck-server: sync: %s failed: %v\n", secretPinPath, err)
			}
		}
	}
	if vaultJSON, ok := vaultSecretVaultJSON(); ok {
		if e.contentNeedsPush(secretVaultPath, vaultJSON) {
			meta, _ := e.getMeta(secretVaultPath)
			sha, _, err := e.ghPutBytes(secretVaultPath, vaultJSON, meta.SHA)
			if err == nil {
				e.setMeta(secretVaultPath, sha, strHash(string(vaultJSON)))
				details = append(details, "pushed "+secretVaultPath)
			} else {
				fmt.Fprintf(os.Stderr, "writerdeck-server: sync: %s failed: %v\n", secretVaultPath, err)
			}
		}
	}
	return details
}

func (e *syncEngine) syncVaultSecrets(remoteMap map[string]string) []string {
	if !e.ready() {
		return nil
	}
	_ = remoteMap
	var details []string
	if vaultEnabled() {
		details = append(details, e.pushVaultSecrets()...)
	}
	if gh, _, err := e.ghGetFile(secretVaultPath); err == nil && gh != nil {
		meta, _ := e.getMeta(secretVaultPath)
		if meta.SHA != gh.SHA {
			raw, err := ghDecodeBytes(gh.Content)
			if err == nil {
				_ = vaultApplySecretVault(raw)
				e.setMeta(secretVaultPath, gh.SHA, strHash(string(raw)))
				details = append(details, "pulled "+secretVaultPath)
			}
		}
	}
	if gh, _, err := e.ghGetFile(secretPinPath); err == nil && gh != nil {
		meta, _ := e.getMeta(secretPinPath)
		if meta.SHA != gh.SHA {
			raw, err := ghDecodeBytes(gh.Content)
			if err == nil {
				vaultApplySecretPin(strings.TrimSpace(string(raw)))
				e.setMeta(secretPinPath, gh.SHA, strHash(string(raw)))
				details = append(details, "pulled "+secretPinPath)
			}
		}
	}
	return details
}
