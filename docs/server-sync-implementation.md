# Server-side GitHub sync

Shipped and device-verified July 2026. The reconcile engine runs in Writerdeck-server (Go). The phone handles token entry and sync settings only. The tablet Lobby Sync tab shows last sync, Sync now when configured, TOKEN NEEDED when RAM is empty, and SYNC OFFLINE / SYNC FAILED when Wi-Fi or GitHub errors apply.

Reconciler behaviour: [decisions.md](decisions.md) §19. Integrity contract: [integrity-audit.md](integrity-audit.md).

## How it works

Engine code lives in `daemon/syncengine.go`, `syncgithub.go`, and `syncmeta.go`. The GitHub token sits in browser `localStorage` (`ghToken`) and tablet RAM via `POST /api/sync/token` — never on disk. Per-note sha and localHash (djb2, same as the old phone engine) persist in `settings.json` as `syncMeta`.

### Token restore after restart

Tablet RAM clears on every `writerdeck` restart. A browser that already saved a token restores sync without manual Save:

1. **WebSocket `needtoken`** — when sync is on, repo is set, and tablet RAM is empty, the server sends `{"type":"needtoken"}` to each browser on connect (and debounced broadcast when lobby info is pushed). `connection.js` calls `respondToNeedToken()` in `sync.js`, which fetches `/api/sync/status`, then `POST /api/sync/token`.
2. **Page load / reconnect** — `refreshSyncStatus()` on WebSocket open also pushes when `syncOn && !configured && ghToken`.
3. **Tablet → browser** — when the tablet has a token but the browser does not, `GET /api/sync/token` pulls it into `localStorage`.

Token POST calls `pushLobbyInfo()` immediately so the Lobby clears TOKEN NEEDED before reconcile finishes.

### Phone UI

Top bar: **Sync setup** opens the GitHub panel. **Save** verifies repo + token; **Sync** runs reconcile (`POST /api/sync/run`). Primary button is Save when no token is saved locally, Sync when one is. Offline sync status appears on the Notes sync setup panel when `/api/sync/status` is unreachable.

Reconcile runs on boot, a three-minute ticker, Home, power sleep, CRUD, manual run, and token verify. Power sleep runs reconcile on the server before suspend — the browser no longer waits. Phone `sync.js` keeps drift UX only; no GitHub fetch calls remain.

Endpoints: `POST /api/sync/token` set or clear token and verify against repo; `GET /api/sync/token` read token from tablet RAM (browser pull); `GET /api/sync/status` for configured, lastSyncAt, syncOn, syncRepo, lastError, syncing; `POST /api/sync/run` trigger reconcile with optional wait; `POST /api/sync/ack` updates lastSyncAt and Lobby after reconcile. Settings also hold syncOn, syncRepo, lastSyncAt, syncMeta, and a legacy pendingSync queue the server drains directly.

## Verify in logs

There is no `needtoken` line in the journal — the message is WebSocket-only. A successful auto-restore after restart looks like:

```
writerdeck-server: client connected 192.168.x.x:…
writerdeck-server: sync reconcile (token)
writerdeck-server: token verify reconcile: 11 notes
```

`sync reconcile (token)` means `POST /api/sync/token` succeeded. Poll `GET /api/sync/status` — `configured` should flip to `true` within a second or two of the browser reconnecting.

```bash
source scripts/_env.sh
rm_ssh 'journalctl -u writerdeck --since "1 hour ago" --no-pager' | rg 'client connected|sync reconcile \(token\)|token verify'
```

## Device sign-off matrix

Token POST then service restart — token cleared from RAM, browser auto-repost via `needtoken` (journal: `sync reconcile (token)`). Token survives Wi-Fi IP change (re-enter once per browser origin if IP changes). Edit on tablet, exit Home — GitHub updated within one reconcile. Edit on GitHub, sync run — tablet updated if file not open in editor. Both edit same note — `(tablet copy).md` plus primary. Empty tablet, non-empty GitHub — pull restore, no junk copy. Tablet delete, rename, create — GitHub reflects without phone. Open note during reconcile — that file skipped, others sync. Power sleep — reconcile completes, no browser hang. Sync off — no GitHub traffic. Bad or expired token — error banner, local saves continue. Phone closed overnight — ticker reconcile backs up tablet edits.

Rollback: redeploy the previous Writerdeck-server binary from git.
