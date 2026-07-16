# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, server-side GitHub sync, and encryption round 1 are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. Physical Home — **done** (session `EVIOCGRAB` + fork `3be2de4` without `suppressNextHomeKey`; Writerdeck binary deployed). Please press physical Home once from edit, read, and Lobby to confirm. [docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md); [docs/decisions.md](docs/decisions.md) §28.
2. Keyboard editing — harness sign-off **110/110/0** @ `00-29-12` (critical **38/38/0**, includes mid-wrapping + cross-para Shift). [docs/editor-testing/todo.md](docs/editor-testing/todo.md).
3. Keywriter fork migration — **done.** Owned fork [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (`master`); helpers + C++ infra + Lobby/shell in-tree. Policy: [docs/decisions.md](docs/decisions.md) §3. Checklist: [docs/editor-migration-1-to-QML/todo-handoff-keywriter-fork.md](docs/editor-migration-1-to-QML/todo-handoff-keywriter-fork.md). Active rule: `.cursor/rules/writerdeck.mdc`.
4. Edit helpers QML → C++ **Phase A** — **next.** Pure text math + undo into fork `EditHelper` (Phases 0 → A1 → A2 → A3). Visual wrap and key dispatch stay QML. [docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md](docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md). Active rule: `.cursor/rules/edit-helper-cpp-migration.mdc` (`writerdeck.mdc` paused).
5. Edit helpers QML → C++ **Phase B** — later (after A3). Move key-chord dispatcher (`handleMacArrow` / `handleMacBackspace` / `handleMacEditKeys`) into C++. Same handoff.
6. Edit helpers QML → C++ **Phase C** — optional later (after A/B pay off). Visual / wrap line math into C++ (`positionToRectangle` / layout access). Same handoff.
7. After migration 2 (A–C as pursued) — **evaluate, do not assume solved by the port:** (a) wrap/caret math still using hand-tuned “magic” gaps vs a cleaner layout-based approach; (b) whether custom undo-on-TextEdit should stay or move toward a purer undo model. Same handoff § After A–C.

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); Writerdeck-keywriter fork → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next: [docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md](docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md) Phase A1. Keyboard harness **110/110/0**. Integrity: [integrity-audit.md](docs/integrity-audit.md).
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
