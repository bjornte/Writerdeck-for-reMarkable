# Server-side GitHub sync — implementation plan

Move the GitHub reconcile engine from the phone browser (`daemon/sync.js`) to Writerdeck-server (Go). The phone keeps **token entry and sync settings UI only** until device-verified; then **remove all reconcile logic from the browser** (no dual engine).

**Status (2026-07-12):** shipped and device-verified. Engine runs on tablet; phone **Setup** panel saves token + repo; **Sync now** on tablet Lobby. Browser keeps `localStorage ghToken` as a convenience cache and auto-reposts to tablet RAM when the service restarts. Tablet **Sync** tab shows a high-contrast **TOKEN NEEDED** box when `syncReady` is false (e-ink: black border, not colour alone).

Supersedes the phone-as-engine model in [research-github-sync.md](research-github-sync.md) § Implementation and [decisions.md](decisions.md) #19 (reconciler behaviour unchanged; execution location moves).

**Contract:** same integrity rules as [integrity-audit.md](integrity-audit.md) — edit lease, empty-push guard, copy-on-clash, no delete-without-confirm, sequential pushes.

---

## Design summary

| Concern | Today | Target |
|---|---|---|
| Sync engine | `sync.js` in phone browser | Go in Writerdeck-server |
| GitHub token | `localStorage ghToken` (phone-only) | Browser `localStorage ghToken` **and** tablet RAM via `POST /api/sync/token`. Phone auto-reposts saved token when tablet RAM is empty after service restart. Never on tablet disk. |
| Per-note SHA/hash | `localStorage ghSha_*`, `ghLocalHash_*` | `settings.json` → `syncMeta` map (non-secret) |
| Triggers | WS connect, poll, exitedit, tabletcrud, manual | Server hooks + optional phone “Sync now” → `POST /api/sync/run` |
| Power sleep | Browser reconcile + `/api/sync/ack` | Server reconcile inline in `sleepForPower` |
| `pendingSync` queue | Browser drains via `/api/sync/pending` | Server drains on enqueue or at reconcile start |

Token on tablet survives: editor quit/reopen, Lobby cycles, Wi‑Fi/IP change. Cleared from tablet RAM on: service restart, reboot, explicit clear, or auth failure (401/403). Browser `ghToken` survives service restarts and re-fills tablet RAM on next page load (or **Save & verify**).

---

## New / changed server surface

### In-memory (never persisted)

```go
var (
    syncTokenMu sync.RWMutex
    syncToken   string // GitHub PAT; empty = not configured
)
```

### Persisted in `settings.json`

```json
{
  "syncOn": true,
  "syncRepo": "owner/repo",
  "lastSyncAt": 1720713600,
  "syncMeta": {
    "foo.md": { "sha": "abc…", "localHash": "1234567890" },
    "bar.md": { "sha": "def…", "localHash": "9876543210" }
  },
  "pendingSync": []
}
```

- `localHash`: same djb2 as `sync.js` `strHash()` — required for reconcile classification and migration from phone `localStorage` (one-time import optional, not required for greenfield).
- `syncMeta` is **not secret** (SHAs and content fingerprints only).

### HTTP endpoints

| Method | Path | Auth | Purpose |
|---|---|---|---|
| POST | `/api/sync/token` | PIN cookie | Body `{ "token": "ghp_…" }` or `{ "token": "" }` to clear. Verify with `GET /repos/{repo}`. Never write token to disk. |
| GET | `/api/sync/status` | PIN cookie | `{ "configured": bool, "lastSyncAt", "syncOn", "syncRepo", "lastError", "syncing": bool }` — drives phone UI. |
| POST | `/api/sync/run` | PIN cookie | `{ "wait": true }` optional — trigger reconcile; returns `{ "ok", "notes": N }`. |
| POST | `/api/sync/import-meta` | PIN cookie | **Optional slice:** one-shot import of `{ "meta": { "foo.md": { "sha", "localHash" } } }` from phone `localStorage` before cutover. |

**Keep (behaviour changes internally):**

- `POST /api/sync/ack` — still updates `lastSyncAt` + Lobby; called by server after its own reconcile (not only browser).
- `GET/POST /api/sync/pending`, `POST /api/sync/pending/clear` — **remove after cutover** (server processes queue directly); keep during parallel-run slice.

**Remove from phone after cutover:**

- All GitHub `fetch()` calls in `sync.js`.
- `startSyncPoll`, `reconcileAll`, `pushNote`, `pullNoteAndUpdate`, `ghDelete`, `handleClash`, `drainPendingSync`, `applyTabletCrud`, `verifyGitHubRepo` (logic moves to Go; verify becomes token POST response).

---

## Go module layout

Extract from `main.go` (single file today) into:

```
daemon/
  main.go          — HTTP routes, hooks, lifecycle
  syncengine.go    — reconcile, push, pull, clash, delete, strHash
  syncgithub.go    — Contents API client (GET list, GET/PUT/DELETE file, GET repo)
  syncmeta.go      — syncMeta load/save under settingsMu
```

No new module dependencies — stdlib `net/http`, `encoding/json`, `encoding/base64` only.

---

## Port map (`sync.js` → Go)

| JS | Go |
|---|---|
| `strHash(s)` | `strHash(s string) string` — identical djb2 |
| `pushNote` | `(*syncEngine).pushNote(name string) error` |
| `pullNoteAndUpdate` | `pullNote(name string) error` — reads/writes disk via existing `readNoteFile` / `writeNoteFile` |
| `handleClash` | `handleClash(name, tabletContent string) error` |
| `ghDelete` | `ghDelete(name string) error` |
| `reconcileOne` | `reconcileOne(name, remoteSha string) error` |
| `reconcileAll` | `reconcileAll(reason string) (int, error)` — mutex, skip `currentNote` |
| `applyRemoteDelete` | `applyRemoteDelete(name string) error` — confirm 404 then `deleteNoteFile` |
| `applyTabletCrud` | inline in `notifyTabletCrud` or drain loop |
| `syncReady()` | `syncEngine.ready()` — `syncOn && syncRepo && token != ""` |

**Disk I/O:** engine calls internal note helpers directly — no loopback HTTP. Respects existing atomic `writeNoteFile` and edit lease via `currentNoteMu`.

**Open note:** skip push/pull/reconcile/delete for `currentNote` only — `reconcileAll` still runs for all other notes (edit lease is per-note, not whole-run abort).

**Phone Setup — Save & verify:** POSTs token → async `reconcileAll("token")`. While a note is open on the tablet, reconcile skips that file and syncs the rest. UI polls `/api/sync/status` until idle; if the open note was skipped, shows *"Token saved — synced other notes; «name» skipped while open."*

---

## Trigger migration

| Trigger | Current | Server hook |
|---|---|---|
| Service start | — | If `syncReady()`, `go reconcileAll("boot")` after settings load |
| 3 min poll | `startSyncPoll` in JS | `time.Ticker` goroutine in `main` |
| Editor save ack (`saved`) | phone `exitedit` → reconcile/push | After non-open note save path: `go pushNote(name)` if not `currentNote` |
| Home / power exit | WS `exitedit` → browser reconcile | `session` path: flush save → `reconcileAll` (blocking) → `signalSyncAck` — **drop browser wait** |
| Tablet CRUD | `notifyTabletCrud` → pending + WS | Process op immediately if `syncReady()`, else leave in `pendingSync` |
| Phone open note (switch) | `pullNote` prev + `pushNote` prev | Server: on `openHandler` after save ack, `pushNote(previous)` if closed; pull target before open if not `currentNote` |
| Phone rename/delete | JS `ghDelete` + `pushNote` | `notesHandler` PATCH/DELETE → mirror to GitHub after disk op |
| Manual “Sync now” | `reconcileAll('manual')` | `POST /api/sync/run` |
| Token save | verify + reconcile in JS | `POST /api/sync/token` → verify → `reconcileAll("token")` |
| WS connect | reconcile on connect | **Remove** — server does not need browser connected |

### Power sleep (slice critical)

Replace:

```
beginSyncWait → broadcast exitedit awaitSync → waitSyncAck(45s) [browser]
```

With:

```
flush save → reconcileAll("power") [server, same 45s ctx] → signalSyncAck internally → broadcast exitedit (informational only)
```

Phone may still receive `exitedit` for UI (drop typing view); it must **not** run reconcile.

---

## Phone UI (interim → final)

### Slice A — parallel run (server engine + JS still present)

- Token field POSTs to `/api/sync/token` **and** keeps working with old path temporarily behind a flag — **skip**; go straight to server-only in slice B to avoid dual reconcile.
- Sync panel reads `GET /api/sync/status` for “configured / last sync / error”.
- “Save & verify” → POST token + repo; show server verify result.
- Yellow banner: `syncOn && !configured` → “Add token in Sync”.
- **Disable** JS reconcile triggers (comment/guard `sync.js` exports to no-op) while testing server — or delete JS engine in slice C only after verify.

### Slice C — phone removal (after device verify)

Delete from `sync.js` everything except:

- `recordEditorDiskBaseline`, `checkDiskDrift`, `notifyDiskChanged` (slice 8 drift UX — disk-only, no GitHub)
- `updateSyncBannerFromState` (rewritten to use `/api/sync/status`)

Delete from `app.js`:

- All `reconcileAll` / `pushNote` / `ghDelete` / `pullNoteAndUpdate` calls
- `startSyncPoll`, `loadSyncConfig` reconcile on startup
- WS handlers that call sync on `exitedit`, `tabletcrud`, connect
- `verifyGitHubRepo` import; wire Save & verify to `/api/sync/token`

Remove legacy `localStorage` keys from the old phone engine: `ghSha_*`, `ghLocalHash_*`, `ghPushFailed_*`, `ghLastSync`. **Keep** `ghToken` — browser cache for re-posting to tablet RAM.

Remove endpoints only used by browser engine: `/api/sync/pending`, `/api/sync/pending/clear` (or keep as no-op stubs one release).

Update [browser-vs-tablet.md](browser-vs-tablet.md): sync engine server-side; phone = token + toggle + manual run.

---

## Implementation slices

Ship and **device-verify each slice** before the next. Integrity gate applies to every slice.

### Slice 1 — Go GitHub client + metadata

- [ ] `syncgithub.go`: repo verify, list `.md`, get/put/delete content (base64).
- [ ] `syncmeta.go`: `syncMeta` in settings; load/save; djb2 `strHash`.
- [ ] Unit-test `strHash` against known JS vectors (empty, ASCII, UTF-8 æøå).
- [ ] No triggers yet.

### Slice 2 — Token API + in-memory store

- [ ] `POST /api/sync/token`, `GET /api/sync/status`.
- [ ] Verify token against `syncRepo`; clear on empty body.
- [ ] Phone Sync panel wired to new endpoints (still no auto reconcile).
- [ ] Device: token survives reconnect; cleared on `systemctl restart writerdeck`.

### Slice 3 — push / pull / delete / clash (single note)

- [ ] Port `pushNote`, `pullNote`, `ghDelete`, `handleClash` with empty-push guard.
- [ ] `POST /api/sync/run` runs full reconcile (port `reconcileAll` + `reconcileOne`).
- [ ] Edit lease: skip `currentNote`.
- [ ] Device: create note → run → appears on GitHub; edit on GitHub → run → tablet updated; clash → `(tablet copy).md`.

### Slice 4 — Lifecycle hooks

- [ ] Boot reconcile goroutine + 3 min ticker.
- [ ] `notifyTabletCrud` → immediate GitHub mirror; drop `pendingSync` persistence when server handles all ops (or drain at reconcile start until then).
- [ ] `openHandler` push previous / pull new; phone PATCH/DELETE mirror.
- [ ] After editor `saved` ack for closed notes → async push.
- [ ] Device: tablet Files CRUD syncs without phone open; note switch push/pull.

### Slice 5 — Power sleep + ack path

- [ ] `sleepForPower`: server `reconcileAll` before suspend; internal ack.
- [ ] WS `exitedit` no longer sets `awaitSync` dependency on browser (field ignored or removed).
- [ ] Device: power button sleep → journal shows server reconcile; GitHub current after wake.

### Slice 6 — Remove phone engine

- [ ] Strip `sync.js` GitHub logic; trim `app.js` triggers.
- [ ] Remove `/api/sync/pending*` if unused.
- [ ] Update architecture, decisions (#19 addendum), integrity-audit, lessons, DONE.
- [ ] Full regression: `scripts/test-edit-session.sh`, manual sync matrix below.

---

## Device test matrix (sign-off before slice 6)

| # | Scenario | Pass criteria |
|---|---|---|
| 1 | Token POST, reboot service | Token cleared; re-POST works |
| 2 | Token POST, Wi‑Fi change | Token still in RAM; sync works on new IP |
| 3 | Edit tablet, exit Home | GitHub updated within one reconcile |
| 4 | Edit GitHub, sync run | Tablet file updated; not open in editor |
| 5 | Both edit same note | `(tablet copy).md` + GitHub version in primary |
| 6 | Empty tablet / non-empty GitHub | Pull restore, no junk copy (#24) |
| 7 | Tablet delete/rename/create | GitHub reflects without phone |
| 8 | Phone delete/rename | GitHub reflects |
| 9 | Open note while reconcile | Open file skipped; others sync |
| 10 | Power sleep | Reconcile completes; no 45s hang waiting for browser |
| 11 | Sync off | No GitHub traffic |
| 12 | Bad/expired token | 401 banner; local saves continue |
| 13 | Phone closed overnight | Ticker reconcile; tablet edits backed up |

---

## Docs to update at cutover

- [architecture.md](architecture.md) § GitHub note-sync
- [decisions.md](decisions.md) — new #25 or #19 addendum: server-side engine, ephemeral token
- [integrity-audit.md](integrity-audit.md) — close “phone offline / localStorage / multi-tab” residuals
- [browser-vs-tablet.md](browser-vs-tablet.md)
- [lessons.md](lessons.md) — deploy Writerdeck: `fetch-keywriter-dist.sh` before `deploy-keywriter.sh`; restart server alone does not reload running Writerdeck GUI
- [research-github-sync.md](research-github-sync.md) — strikethrough phone-engine paragraph, link here

---

## Rollback

Until slice 6 ships, keep `sync.js` reconcile code on a git tag. Slice 6 is a one-way delete — rollback = redeploy previous `Writerdeck-server` + restore embedded `sync.js` from tag.

---

## Resume prompt (paste for implementation chat)

> Implement server-side GitHub sync per [docs/server-sync-implementation.md](server-sync-implementation.md). Start at slice N. Port behaviour from [daemon/sync.js](../daemon/sync.js); hooks in [daemon/main.go](../daemon/main.go). Token in RAM only. After device-verified slice 6, remove phone reconcile code. Constraints: `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`; integrity contract in [integrity-audit.md](integrity-audit.md).
