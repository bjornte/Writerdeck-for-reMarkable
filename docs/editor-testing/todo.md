# Keyboard harness — remaining verify

Mac-style editing over the phone/WebSocket path is covered by `scenarios.go` and `scenarios_regression.go` (device PASS with `--fast --hard-reset` as of 2026-07). This file tracks what is still open and how to extend the harness.

Read: [scenario-cookbook.md](scenario-cookbook.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22.

## Open

- Undo/redo — five scenarios in `scenarios_undo.go` need device PASS (`bash scripts/test-keyboard-harness.sh -m undo --fast --hard-reset`).

## Harness layout

| Piece | Path |
|-------|------|
| Scenarios (core) | `daemon/cmd/edit-harness/scenarios.go` |
| Scenarios (regressions) | `daemon/cmd/edit-harness/scenarios_regression.go` |
| Scenarios (undo) | `daemon/cmd/edit-harness/scenarios_undo.go` |
| Runner | `daemon/cmd/edit-harness/main.go` |
| Shell wrapper | `scripts/test-keyboard-harness.sh` |
| QML key logic | `third_party/keywriter/build-keywriter.sh` |

## Dev loop (batch first)

Debugging and sign-off are different jobs. See [lessons.md](../lessons.md) § Device verify and iteration and § Harness batch workflow.

Triage once: `--unit`, then `--fast --hard-reset`. Classify FAILs from `docs/recon/test-keyboard-harness-*.txt`. Confirm with `-s NAME --fast` on the same binary — no deploy between. Batch harness fixes locally; batch QML fixes; at most one Writerdeck deploy per session before the next full rerun. Sign-off: full suite `--hard-reset`, `test-edit-session.sh`, clean `journalctl`.

Port new cases from [scenario-cookbook.md](scenario-cookbook.md) into `scenarios_regression.go` or `scenarios_undo.go`, then `--unit`.

## Resume prompt

> Read this file and `scenario-cookbook.md`. Run `bash scripts/test-keyboard-harness.sh -m undo --fast --hard-reset`. Fix QML or scenarios in batch; one deploy max. Sign-off: full suite `--hard-reset`, `test-edit-session.sh`, clean journalctl.
