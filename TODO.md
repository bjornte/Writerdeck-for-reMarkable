# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, server-side GitHub sync, and encryption round 1 are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. Physical Home — single input path. Handoff: [docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md).
2. Edit scroll/cursor — patch is in `build-keywriter.sh`; run `test-keyboard-harness.sh` and spot-check last-line typing on e-ink.

Encryption regression (not open work): [docs/todo-handoff-encryption.md](docs/todo-handoff-encryption.md) — `test-vault.sh`, `test-vault-e2e.sh` (sync on).

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); patched keywriter → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next unchecked: Physical Home ([todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md)); edit scroll/cursor device verify. Integrity: [integrity-audit.md](docs/integrity-audit.md). After QML edits: `test-edit-session.sh` (§21); after arrow/selection QML: `test-keyboard-harness.sh` (§22).
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
