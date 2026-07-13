# TODO: Encrypted notes

Fresh-agent handoff. Read [architecture.md](architecture.md), [decisions.md](decisions.md) (document integrity), [browser-vs-tablet.md](browser-vs-tablet.md), [improvements.md](improvements.md). Device verify per `.cursor/rules/writerdeck.mdc`.

Sign-off is done. Implement in order below. Write ADR as **decisions.md §31** before shipping.

---

## Problem

LAN pairing PIN gates the phone API; it does not protect note bodies on disk or in a private GitHub repo. Owner wants optional at-rest encryption with a **user-chosen 6-digit PIN** (not the boot-random pairing PIN), unlock UI on the **tablet only**, GitHub backup of ciphertext, and PIN recovery via a **plaintext file in the repo** (`secret/`) — acceptable because the notes repo is private.

## Decisions

**Two PINs.** Pairing PIN (`pinDigits`, phone browser) and encryption PIN (vault) are independent. `pinDigits: none` does not disable encryption.

**File shape.** Encrypted notes use suffix `.md.enc` beside plain `.md` in `Writerdeck-user-documents/`. No `private/` subfolder. Plain notes unchanged.

**Crypto.** AES-256-GCM per file; random 32-byte data key; PIN → scrypt (high N) → KEK wraps data key. Persist in `.Writerdeck/settings.json`: `encryptionEnabled`, `vaultSalt`, `vaultVerifier`, `wrappedDataKey` only — never the PIN. On disk: magic `WDENC1` + nonce + ciphertext+tag. Stdlib + `golang.org/x/crypto/scrypt`; `CGO_ENABLED=0`.

**Vault session.** Unlocked = `dataKey` in server RAM. **Lock on every return to Lobby** (`handleHome` / `isLobby = true`). Unlock via tablet numpad (touch) or USB/BT keyboard digits + Enter. Rate-limit failed unlocks (reuse `auth.go` pattern).

**Integrity carve-out.** UTF-8 Markdown contract applies to `.md` only. `.md.enc` are opaque binary; plaintext exists in editor buffer only while vault unlocked. Update [integrity-audit.md](integrity-audit.md).

**Per-note lifecycle.** No bulk encrypt on enable. Files tab: **Encrypt** and **Decrypt** actions (second row of file actions or equivalent). Encrypt: read plaintext → write `.md.enc` → delete `.md`. Decrypt: reverse when vault unlocked.

**GitHub sync.** When sync is on, treat `.md.enc` as opaque bytes (extend `listLocalNoteNames`, remote list filter, push/pull — skip HTML guard). Also sync recovery material under repo path `secret/`:
- `secret/pin` — current encryption PIN as plain text (rewrite on PIN change)
- `secret/vault` — JSON: `salt`, `wrappedDataKey`, `verifier` (for restore without local settings)

Exclude `secret/` from phone note list and LAN metadata APIs. Sync engine may read/write `secret/`; phone UI must not surface `secret/` contents.

**Forgot PIN.** Recover from `secret/pin` on GitHub after re-deploy. No email relay. If repo and tablet are both gone and PIN forgotten, data is lost — acceptable.

**Phone download.** Decrypt server-side only when vault unlocked. Note need not be open on tablet. If vault locked: reject download, push unlock request to editor socket, show root-level PIN overlay on tablet (any app state — Lobby, edit, read). Phone shows “Enter private PIN on tablet”; poll or WebSocket until unlocked, then retry. No encryption PIN UI on phone.

**Phone otherwise.** List may show `.md.enc` names as locked; no body, no upload into `.md.enc`, no decrypt API for remote clients.

---

## Touch points

| Area | Files |
|------|-------|
| Crypto + vault | new `daemon/vault.go`; unit tests |
| Notes I/O | `daemon/notes.go` — `notesSafe` for `.md.enc`; loopback GET decrypt / PUT encrypt; remote GET 403 |
| Sync | `daemon/syncengine.go`, `daemon/syncgithub.go` — `.md.enc` + `secret/*` |
| Socket ops | `daemon/editor.go` `handleEditorReq`; `daemon/lobby.go` `pushLobbyInfo` vault flags |
| Tablet UI | `third_party/keywriter/lobby/*.inc`, `build-keywriter.sh`, `lobby_bridge.cpp/.h` — numpad overlay, Settings enable/change PIN, Files encrypt/decrypt |
| Phone | `daemon/notes-ui.js`, `daemon/handlers.go` download path, WebSocket unlock-wait |
| Settings schema | `daemon/settings.go` |
| ADR | `docs/decisions.md` §31 |
| Prune stale | `docs/improvements.md` email-recovery bullet; `TODO.md` Phase 10 pointer |

`notesSafe` today rejects `/`. Encrypted names are flat basenames like `diary.md.enc`.

Editor save path is loopback `PUT /api/notes/{name}` from QML `saveFile()` — encrypt at server on PUT for `.md.enc`; decrypt on GET. No keywriter crypto in C++ unless unavoidable.

---

## Tablet UI sketch

**Settings → Private notes:** Off | Setup (email not required; recovery is GitHub `secret/pin`) | On (change PIN, status locked/unlocked).

**Numpad overlay:** Root `Item`, `z` above Lobby and editor. Digits 0–9, Backspace, Done; keyboard 0–9 + Enter mirrors pad. Used for setup, confirm, unlock-on-demand (download prompt), change PIN.

**Files:** Row marker for `.md.enc`. Buttons: Encrypt (plain selected) / Decrypt (`.md.enc` selected). New plain note unchanged; optional “New encrypted” creates `name.md.enc` when vault unlocked.

---

## Socket / API (add)

Tablet → server (`req` op): `setvaultpin`, `unlockvault`, `lockvault`, `encryptnote`, `decryptnote`.

Server → tablet (`cmd`): `requestvaultunlock` `{reason, name}` for phone-download path.

`pushLobbyInfo` / `notes` list: `encryptionEnabled`, `vaultLocked`; per item `encrypted: true`.

---

## Implementation order

1. ADR §31 + `vault.go` + tests (KDF, wrap, roundtrip, wrong PIN, PIN change re-wrap only).
2. `notes.go` read/write/list; lock on Lobby in editor/home path.
3. Tablet numpad + Settings setup/change + vault lock/unlock.
4. Files encrypt/decrypt + open/save `.md.enc` through loopback.
5. Sync `.md.enc` + `secret/pin` + `secret/vault`.
6. Phone download + tablet unlock-on-demand + WS/poll retry.
7. `integrity-audit.md` residuals; `test-edit-session.sh` or new `scripts/test-vault.sh` for loopback encrypt/save/lock/unlock.

---

## Verify (device)

1. Enable encryption, set PIN, create encrypted note, save, return to Lobby → locked.
2. Unlock, edit, save, lock again — disk file is not valid UTF-8 markdown.
3. Decrypt note → plain `.md` returns.
4. Sync on: GitHub has `foo.md.enc`, `secret/pin`, `secret/vault`; pairing PIN unrelated.
5. Phone download encrypted note with vault locked → tablet PIN overlay → download succeeds without opening note on e-ink.
6. Phone download with vault already unlocked → immediate.
7. Change PIN → `secret/pin` updated on next sync; old ciphertext still readable.

Do not declare done without tablet + phone browser checks per project rules.
