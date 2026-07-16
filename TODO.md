# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, server-side GitHub sync, and encryption round 1 are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. Physical Home — single input path (exclusive gpio grab so page buttons and Home are not confused with keyboard keys). Handoff: [docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md).
2. Keyboard editing — harness done (**105** scenarios). Critical gate green (**36/36**); full suite **93/12** @ `18-57-31` (Patch LOC **386**). Product sign-off still **105/105** — [docs/editor-testing/todo.md](docs/editor-testing/todo.md). Leftover fails are product polish, not blocked on fork migration.
3. Keywriter fork migration — **done.** Owned fork [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (`master`); helpers + C++ infra + Lobby/shell in-tree (`68f6e32`). Ownership and upstream-merge policy: [docs/decisions.md](docs/decisions.md) §3. Checklist: [docs/editor-migration/todo-handoff-keywriter-fork.md](docs/editor-migration/todo-handoff-keywriter-fork.md). Active rule: `.cursor/rules/writerdeck.mdc`.

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); Writerdeck-keywriter fork → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next: **Physical Home** ([todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md)), or leftover keyboard harness fails toward **105/105**. Fork migration done (`68f6e32`; Patch LOC **386**). Integrity: [integrity-audit.md](docs/integrity-audit.md).
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
