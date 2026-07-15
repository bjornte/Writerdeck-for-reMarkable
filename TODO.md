# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, server-side GitHub sync, and encryption round 1 are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. Physical Home — single input path. Handoff: [docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md).
2. Keyboard editing — **90-scenario** device harness; kernel fixes for combo/wrap/undo/visual-x (best full **83-scenario** run **64/16** @ `071f998`; no full **90/90** yet). Handoff: [docs/editor-testing/todo.md](docs/editor-testing/todo.md).

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); patched keywriter → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next unchecked: Physical Home; keyboard editing — suite **90** scenarios, best full **83** run **64/16** @ `071f998` ([editor-testing/todo.md](docs/editor-testing/todo.md), [milestone-runs.md](docs/editor-testing/milestone-runs.md)). Integrity: [integrity-audit.md](docs/integrity-audit.md). After QML edits: `test-edit-session.sh` (§21); after arrow/selection QML: `test-keyboard-harness.sh --fast` (§22, **90/90** sign-off).
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
