# Server-side GitHub sync

Shipped and device-verified July 2026. The reconcile engine runs in Writerdeck-server (Go). The phone handles token entry and sync settings only. The tablet Lobby has Sync now and a TOKEN NEEDED warning when RAM is empty.

Reconciler behaviour: [decisions.md](decisions.md) §19. Integrity contract: [integrity-audit.md](integrity-audit.md).

## How it works

Engine code lives in `daemon/syncengine.go`, `syncgithub.go`, and `syncmeta.go`. The GitHub token sits in browser localStorage and tablet RAM via `POST /api/sync/token` — never on disk. The phone re-posts the saved token when tablet RAM is empty after a service restart. Per-note sha and localHash (djb2, same as the old phone engine) persist in `settings.json` as `syncMeta`.

Reconcile runs on boot, a three-minute ticker, Home, power sleep, CRUD, manual run, and token verify. Power sleep runs reconcile on the server before suspend — the browser no longer waits. Phone `sync.js` keeps drift UX only; no GitHub fetch calls remain.

Endpoints: `POST /api/sync/token` set or clear token and verify against repo; `GET /api/sync/status` for configured, lastSyncAt, syncOn, syncRepo, lastError, syncing; `POST /api/sync/run` trigger reconcile with optional wait; `POST /api/sync/ack` updates lastSyncAt and Lobby after reconcile. Settings also hold syncOn, syncRepo, lastSyncAt, syncMeta, and a legacy pendingSync queue the server drains directly.

## Device sign-off matrix

Token POST then service restart — token cleared from RAM, re-POST works. Token survives Wi-Fi IP change. Edit on tablet, exit Home — GitHub updated within one reconcile. Edit on GitHub, sync run — tablet updated if file not open in editor. Both edit same note — `(tablet copy).md` plus primary. Empty tablet, non-empty GitHub — pull restore, no junk copy. Tablet delete, rename, create — GitHub reflects without phone. Phone delete and rename — GitHub reflects. Open note during reconcile — that file skipped, others sync. Power sleep — reconcile completes, no browser hang. Sync off — no GitHub traffic. Bad or expired token — error banner, local saves continue. Phone closed overnight — ticker reconcile backs up tablet edits.

Rollback: redeploy the previous Writerdeck-server binary from git.
