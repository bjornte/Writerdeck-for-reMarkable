# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, server-side GitHub sync, and encryption round 1 are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. Physical Home — single input path (exclusive gpio grab so page buttons and Home are not confused with keyboard keys). Handoff: [docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md).
2. Keyboard editing — harness done (**105** scenarios). Critical gate green (**36/36**); full suite **91/14** @ `12-41-15` @ `377a053`. Product sign-off still **105/105** — [docs/editor-testing/todo.md](docs/editor-testing/todo.md). Do not prioritize burning down the 14 leftover fails ahead of the keywriter fork.
3. Keywriter fork migration — **preferred path out of patch-script debt.** Handoff: [docs/todo-handoff-keywriter-fork.md](docs/todo-handoff-keywriter-fork.md). Rule: `.cursor/rules/keywriter-fork-migration.mdc`. Policy: [docs/decisions.md](docs/decisions.md) §3.
   - **Fork:** [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (`master`) — CI pinned (Phase 1 done; critical **36/36**).
   - **Phase 2 (next):** move edit behavior into forked C++/QML by criticality (A→D). Do not first fix leftover harness fails.
   - **Phase 3:** shrink script; restore general `writerdeck.mdc` rule.

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); patched keywriter → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next: **keywriter fork Phase 2A** — caret/selection/backspace into [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter); handoff [todo-handoff-keywriter-fork.md](docs/todo-handoff-keywriter-fork.md). Physical Home later. Keyboard **105/105** is product sign-off only. Integrity: [integrity-audit.md](docs/integrity-audit.md).
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
