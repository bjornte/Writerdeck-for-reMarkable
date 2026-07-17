# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, server-side GitHub sync, and encryption round 1 are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md). Words: [docs/terms.md](docs/terms.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. Physical Home — **done** (session `EVIOCGRAB` + fork `3be2de4` without `suppressNextHomeKey`; Writerdeck binary deployed). Please press physical Home once from edit, read, and Lobby to confirm. [docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md); [docs/decisions.md](docs/decisions.md) §28.
2. Keyboard editing — harness sign-off **110/110/0** @ `14-52-09` (critical **38/38/0**, includes mid-wrapping + cross-para Shift). [docs/editor-testing/todo.md](docs/editor-testing/todo.md).
3. Keywriter fork migration — **done.** Owned fork [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (`master`); helpers + C++ infra + Lobby/shell in-tree. Policy: [docs/decisions.md](docs/decisions.md) §3. Checklist: [docs/editor-migration-1-to-QML/todo-handoff-keywriter-fork.md](docs/editor-migration-1-to-QML/todo-handoff-keywriter-fork.md). Active rule: `.cursor/rules/writerdeck.mdc`.
4. Edit helpers QML → C++ **Phase A** — **done.** Pure text math + undo in fork `EditHelper` (Phases 0 → A3; full **110/110/0** @ `10-12-39`, fork `a92ad2b`). [docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md](docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md).
5. Edit helpers QML → C++ **Phase B** — **done.** Key-chord dispatcher in fork `EditHelper` (fork `57bfc21`; full **110/110/0** @ `10-29-42`). [docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md](docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md).
6. Edit helpers QML → C++ **Phase C** — **done.** Visual-line math in `EditHelper`; fork `6a15e08`; full **110/110/0** @ `14-52-09`. [docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md](docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md).
7. After migration 2 (A–C) — **done.** **Keep** hand-tuned wrap/caret gaps and custom `EditHelper` undo (do not lean on Qt undo or rewrite wrap for purity). [docs/decisions.md](docs/decisions.md) §30; handoff § After A–C.
8. Fork wrap-up hygiene — **done.** Fork owns QML assembly (`assemble-qml.sh` → committed `main.qml`); `build-keywriter.sh` is clone + assert + build only. Edit-session PASS; critical **38/38/0** @ `17-22-24`; full **110/110/0** @ `17-23-47`; fork tip `0bb3b70`.
9. Fork upstream ancestry — **done.** Ours-merge `5946cae` links [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) to [dps/remarkable-keywriter](https://github.com/dps/remarkable-keywriter) (merge-base `ddc9e73`; tree unchanged, no force-push). How to pull upstream: [docs/decisions.md](docs/decisions.md) §3; fork README.

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); Writerdeck-keywriter fork → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next: owner physical Home check ([docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md)). Keyboard harness **110/110/0** @ `17-23-47` (fork `0bb3b70`). Fork ancestry linked @ `5946cae` ([docs/decisions.md](docs/decisions.md) §3). Migration 2 keep: §30. Integrity: [integrity-audit.md](integrity-audit.md).
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
