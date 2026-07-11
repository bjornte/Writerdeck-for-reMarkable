# Document integrity — audit (2026-07-11)

Living audit for the product contract in [architecture.md](architecture.md) § Document integrity and [decisions.md](decisions.md) § Document integrity. **Slices 1–11 are shipped** (`b1ce2bc`…`f72282d`, device-verified). Open items stay open until closed or accepted.

## Fixed (slices 1–11)

| Slice | What |
|---|---|
| 1 | Edit lease — `notifyOpen` + `openedit` WS; reconcile skips `tabletOpenNote` |
| 2 | Content fidelity — plain-markdown save contract, load sanitizer, `toggleMode` fix, server HTML guard |
| 3 | `notedeleted` / `noterenamed` — editor notified on phone rename/delete of open file |
| 4 | Reconcile policy — `openNote` in `/api/status`; `reconcileAll` gated on edit lease |
| 5 | OCC — GET `ETag`; PUT overwrite requires `If-Match`; sync sends revision |
| 6 | Atomic server writes — `writeNoteFile` temp+rename |
| 7 | Tablet CRUD → GitHub — `tabletcrud` WS + `pendingSync` queue |
| 8 | Disk↔buffer drift — `diskchanged` WS, phone drift banner, `POST /api/reload` |
| 9 | Autosave — 45 s `autosaveTimer` while editing |
| 10 | Tablet atomic saves — loopback `PUT /api/notes` → `writeNoteFile` |
| 11 | Save before deploy/stop — `POST /api/flush-save`, deploy graceful wait, SIGTERM flush |

Also under contract: empty-push guard (#19), Lobby Home wipe fix (#24). Shipped summary: [../DONE.md](../DONE.md) § Document integrity.

## Known open (residual risks)

Mitigated or bounded — not eliminated. Revisit before calling integrity “closed.”

**Durability**
- SIGKILL, editor segfault, or power loss between autosave ticks (up to ~45 s of typing).
- `POST /api/flush-save` / SIGTERM path fails if editor socket down, `autosavenow` ack times out, or mismatched Writerdeck binary (no loopback save).

**Save / load timing**
- `doLoad` is async — rapid note switch (Ctrl-K, Lobby Files) can interleave saves and loads.
- `/api/open` continues on save-ack timeout — previous note may not be flushed before switch.

**Concurrency / sync**
- Real edit clashes create `(tablet copy)` duplicates; clash overwrites disk without auto-reloading the editor (drift banner is manual).
- Stale `tabletOpenNote` after phone-back can skip the wrong file in reconcile until status poll refreshes.
- External GitHub edits (VS Code, web, git) — marker-aware delete helps; duplicates/resurrections still possible (#19).
- `localStorage` loss (new browser, cleared site data) — surprise push/pull/clash on next reconcile.
- Multiple browser tabs share one sync state — serialized but confusing.
- Power sleep: up to 45 s wait for phone reconcile ack; GitHub may lag if phone offline.

**Auth / ops**
- PIN `none` — anyone on LAN can mutate notes (integrity + confidentiality on untrusted Wi-Fi).
- `restore-wiped-notes.sh` — Git history only; no partial-truncation or duplicate cleanup beyond `(tablet copy)` names.

## Unknown (unbounded)

No claim that the threat surface is complete.

- Bugs in paths not stress-tested (rapid switching, multi-device, sleep/wake edge cases, first boot with sync on).
- Firmware OTA side effects on deployed binaries or unit file.
- Qt / QML regressions not caught by `test-edit-session.sh` (save paths, autosave, loopback PUT).
- Novel failure modes (network partition mid-save, partial HTTP write despite atomic rename, future feature regressions).

**Process:** before shipping note/save/sync/lifecycle changes, ask: *can this lose text, write wrong bytes, or overwrite without the user knowing?* If yes, it does not ship until mitigated or **explicitly accepted** by the owner and logged here.

## Reference patterns

| Pattern | Status | Pointer |
|---|---|---|
| Edit lease | Shipped | [Ink & Switch local-first](https://www.inkandswitch.com/essay/local-first/) |
| Plain-text save contract | Shipped | — |
| Optimistic concurrency | Shipped (HTTP) | [RFC 7232](https://httpwg.org/specs/rfc7232.html) |
| Conflict copies | Partial (`(tablet copy)`) | [Syncthing conflicts](https://docs.syncthing.net/users/syncing.html) |
| Atomic durable write | Shipped | [google/renameio](https://github.com/google/renameio) |
| CRDT / OT | Not planned | [Automerge](https://automerge.org/) |
