# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, server-side GitHub sync, and encryption round 1 are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. Physical Home — single input path (exclusive gpio grab so page buttons and Home are not confused with keyboard keys). Handoff: [docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md).
2. Keyboard editing — harness rewrite done (**105** scenarios, pattern uni1/uni5/bi1+1/bi3+5/bi7+7). Device baseline **72/33** (critical **26/36**) @ `10-01-42`. Product QML fixes next — [docs/editor-testing/todo.md](docs/editor-testing/todo.md).
3. Keywriter fork migration — fork upstream keywriter and retire behavioral patching from `third_party/keywriter/build-keywriter.sh` in phases. Phase 1: pin builds to Writerdeck fork (`KEYWRITER_REPO`/`KEYWRITER_REF`) with no behavior change. Phase 2: move edit-mode behavior patches (arrow/selection/undo/wrap) from script into forked source, keep script for deterministic build glue. Phase 3: shrink script to minimal bootstrap/deploy patch set and document fork ownership/update policy in [docs/decisions.md](docs/decisions.md).

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); patched keywriter → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next unchecked: Physical Home; keyboard editing (harness done, QML fixes) — [editor-testing/todo.md](docs/editor-testing/todo.md); keywriter fork migration to retire behavior patching in `build-keywriter.sh`. Integrity: [integrity-audit.md](docs/integrity-audit.md). After QML: `test-edit-session.sh` (§21); after arrow/selection QML: `test-keyboard-harness.sh --fast` 105/105 (§22).
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
