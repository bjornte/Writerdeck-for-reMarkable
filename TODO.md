# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, server-side GitHub sync, and encryption round 1 are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. Physical Home — single input path (exclusive gpio grab so page buttons and Home are not confused with keyboard keys). Handoff: [docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md).
2. Keyboard editing — harness done (**105** scenarios). Latest **74/31** (critical **26/36**, 10 open) @ `11-38-40` @ `7603357`. Next: critical pass **36/36** then full **105/105** — [docs/editor-testing/todo.md](docs/editor-testing/todo.md).
3. Keywriter fork migration — fork upstream keywriter and retire behavioral patching from `third_party/keywriter/build-keywriter.sh` in phases. Phase 1: pin builds to Writerdeck fork (`KEYWRITER_REPO`/`KEYWRITER_REF`) with no behavior change. Phase 2: move edit-mode behavior patches (arrow/selection/undo/wrap) from script into forked source, keep script for deterministic build glue. Phase 3: shrink script to minimal bootstrap/deploy patch set and document fork ownership/update policy in [docs/decisions.md](docs/decisions.md).

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); patched keywriter → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next unchecked: Physical Home; keyboard editing critical pass (26/36 @ `7603357`, 10 failures listed in [editor-testing/todo.md](docs/editor-testing/todo.md)); keywriter fork migration. Integrity: [integrity-audit.md](docs/integrity-audit.md). After QML deploy: `test-edit-session.sh` then `test-keyboard-harness.sh -t critical --fast` (36/36) then `--fast` full 105/105 (§22). Do not run edit-session and harness in parallel.
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
