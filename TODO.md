# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, server-side GitHub sync, and encryption round 1 are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. Physical Home — **done** (session `EVIOCGRAB` + fork `3be2de4` without `suppressNextHomeKey`; Writerdeck binary deployed). Please press physical Home once from edit, read, and Lobby to confirm. [docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md); [docs/decisions.md](docs/decisions.md) §28.
2. Keyboard editing — harness sign-off **110/110/0** @ `14-52-09` (critical **38/38/0**, includes mid-wrapping + cross-para Shift). [docs/editor-testing/todo.md](docs/editor-testing/todo.md).
3. Keywriter fork migration — **done.** Owned fork [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (`master`); helpers + C++ infra + Lobby/shell in-tree. Policy: [docs/decisions.md](docs/decisions.md) §3. Checklist: [docs/editor-migration-1-to-QML/todo-handoff-keywriter-fork.md](docs/editor-migration-1-to-QML/todo-handoff-keywriter-fork.md). Active rule: `.cursor/rules/writerdeck.mdc`.
4. Edit helpers QML → C++ **Phase A** — **done.** Pure text math + undo in fork `EditHelper` (Phases 0 → A3; full **110/110/0** @ `10-12-39`, fork `a92ad2b`). [docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md](docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md).
5. Edit helpers QML → C++ **Phase B** — **done.** Key-chord dispatcher in fork `EditHelper` (fork `57bfc21`; full **110/110/0** @ `10-29-42`). [docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md](docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md).
6. Edit helpers QML → C++ **Phase C** — **done.** Visual-line math in `EditHelper`; fork `6a15e08`; full **110/110/0** @ `14-52-09`. [docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md](docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md).
7. After migration 2 (A–C) — **done.** **Keep** hand-tuned wrap/caret gaps and custom `EditHelper` undo (do not lean on Qt undo or rewrite wrap for purity). [docs/decisions.md](docs/decisions.md) §30; handoff § After A–C.
8. Fork wrap-up hygiene — **done** (pending device verify this session). Fork owns QML assembly (`assemble-qml.sh` → committed `main.qml`); `build-keywriter.sh` is clone + assert + build only. Ship tip in fork README (`e2a8436` / tip `0bb3b70`). Verify: edit-session, critical **38/38/0**, full **110/110/0**.
9. Fork upstream ancestry (optional, separate session) — restore a real merge-base with [dps/remarkable-keywriter](https://github.com/dps/remarkable-keywriter) so future upstream pulls are ordinary merges. Do not mix with item 8. Policy already in [docs/decisions.md](docs/decisions.md) §3.

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); Writerdeck-keywriter fork → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next: owner physical Home check ([docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md)), or fork wrap-up item 8 (CI insert/concat → in-tree; fresh session). Keyboard harness **110/110/0** @ `14-52-09` (fork `6a15e08`). Migration 2 keep: [docs/decisions.md](docs/decisions.md) §30. Integrity: [integrity-audit.md](docs/integrity-audit.md).
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
