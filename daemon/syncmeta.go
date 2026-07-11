package main

import "strconv"

// noteSyncMeta tracks the last known GitHub blob SHA and local content fingerprint
// for one note. Stored in settings.json (non-secret).
type noteSyncMeta struct {
	SHA       string `json:"sha"`
	LocalHash string `json:"localHash"`
}

// strHash is djb2 — matches sync.js strHash() (signed 32-bit overflow per step).
func strHash(s string) string {
	h := int32(5381)
	for _, r := range s {
		if r > 0xffff {
			r -= 0x10000
			h = ((h << 5) + h + int32(0xd800+(r>>10)))
			h = ((h << 5) + h + int32(0xdc00+(r&0x3ff)))
		} else {
			h = ((h << 5) + h + int32(r))
		}
	}
	return strconv.FormatUint(uint64(uint32(h)), 10)
}

func (e *syncEngine) getMeta(name string) (noteSyncMeta, bool) {
	settingsMu.Lock()
	defer settingsMu.Unlock()
	if curSettings.SyncMeta == nil {
		return noteSyncMeta{}, false
	}
	m, ok := curSettings.SyncMeta[name]
	return m, ok
}

func (e *syncEngine) setMeta(name, sha, localHash string) {
	settingsMu.Lock()
	defer settingsMu.Unlock()
	if curSettings.SyncMeta == nil {
		curSettings.SyncMeta = map[string]noteSyncMeta{}
	}
	curSettings.SyncMeta[name] = noteSyncMeta{SHA: sha, LocalHash: localHash}
	saveSettingsLocked()
}

func (e *syncEngine) removeMeta(name string) {
	settingsMu.Lock()
	defer settingsMu.Unlock()
	if curSettings.SyncMeta == nil {
		return
	}
	delete(curSettings.SyncMeta, name)
	saveSettingsLocked()
}
