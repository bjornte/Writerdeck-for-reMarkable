# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, server-side GitHub sync, and encryption round 1 are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. Physical Home — single input path (exclusive gpio grab so page buttons and Home are not confused with keyboard keys). Handoff: [docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md).
2. Keyboard editing — harness done (**105** scenarios). Critical gate green (**36/36**); full suite **93/12** @ `18-57-31` (Phase 3 Lobby/shell; Patch LOC **386**). Product sign-off still **105/105** — [docs/editor-testing/todo.md](docs/editor-testing/todo.md). Do not prioritize burning down the leftover fails ahead of Phase 3.
3. Keywriter fork migration — **preferred path out of patch-script debt.** Handoff: [docs/editor-migration/todo-handoff-keywriter-fork.md](docs/editor-migration/todo-handoff-keywriter-fork.md). Rule: `.cursor/rules/keywriter-fork-migration.mdc`. Policy: [docs/decisions.md](docs/decisions.md) §3.
   - **Fork:** [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (`master`) — CI pinned; helpers + C++ infra + Lobby/shell QML in-tree (`68f6e32`).
   - **Phase 2:** done (A–D).
   - **Phase 3 (next):** document fork ownership; restore general `writerdeck.mdc` rule. Lobby/shell shrink verified (critical **36/36**, full **93/12**, Patch LOC **386**).

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); patched keywriter → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next: **keywriter fork Phase 3** — document fork ownership + restore general rules (Lobby/shell QML in fork `68f6e32`); handoff [todo-handoff-keywriter-fork.md](docs/editor-migration/todo-handoff-keywriter-fork.md). Physical Home later. Keyboard **105/105** is product sign-off only. Integrity: [integrity-audit.md](docs/integrity-audit.md).
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
