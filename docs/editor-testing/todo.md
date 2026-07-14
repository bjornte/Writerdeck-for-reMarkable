# TODO: Keyboard editing bugs + harness loop

Fix flaky Mac-style editing in Writerdeck (Qt `TextEdit` + custom `handleKey` in `build-keywriter.sh`). Drive fixes through the device keyboard harness — not manual Lobby typing.

Read first: [scenario-cookbook.md](scenario-cookbook.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22, `.cursor/rules/writerdeck.mdc`.

## Reported bugs

- Alt+Backspace should delete word before cursor, not line or wrong span
- Ctrl+Backspace should delete line
- Arrow Down in a multiline paragraph should move one visual/logical line, not jump paragraph or end-of-line
- Shift+Down then Shift+Up should shrink downward selection; currently expands upward
- Repeated Shift+Left should extend selection; currently stuck at one word/char
- Ctrl+Z undo sometimes does the wrong thing

Phone/WebSocket path is in scope for harness; USB/qmap is spot-check only after qmap changes ([lessons.md](../lessons.md)).

## Harness (extend, don't reinvent)

| Piece | Path |
|-------|------|
| Scenarios (core) | `daemon/cmd/edit-harness/scenarios.go` |
| Scenarios (regressions) | `daemon/cmd/edit-harness/scenarios_regression.go` |
| Runner | `daemon/cmd/edit-harness/main.go` |
| Shell wrapper | `scripts/test-keyboard-harness.sh` |
| Editor state API | `GET /api/test/editor-state` |
| QML key logic | `third_party/keywriter/build-keywriter.sh` |

### Dev loop (fast)

1. Pick or add scenario — [scenario-cookbook.md](scenario-cookbook.md) or `scenarios_regression.go`.
2. Lint: `bash scripts/test-keyboard-harness.sh --unit` (no tablet).
3. QML changed → deploy Writerdeck.
4. `bash scripts/test-keyboard-harness.sh -s NAME --fast` then `--no-prepare` on repeats.

See [README.md](README.md) for flags (`-m`, `--fast`, `--no-prepare`).

### Regression scenarios in repo (may FAIL until fixed)

- `down-one-logical-line`
- `shift-down-then-up-shrinks`
- `shift-left-repeat-from-end`
- `alt-backspace-deletes-word`
- `ctrl-backspace-deletes-line`

Cookbook lists additional scenarios not yet coded — port as needed.

## Implementation notes

Custom selection/cursor math in `handleMacArrow` likely fights Qt defaults. Wrapped lines need visual-line logic, not `\n`-only math. Do not swap editors in this task.

## Acceptance

- All `scenarios.go` and `scenarios_regression.go` PASS on device.
- `--unit` PASS on Mac.
- After QML deploy: `test-edit-session.sh` PASS; `journalctl -u writerdeck -n 30` clean.

## Out of scope

Editor replacement, USB evdev harness, phone browser UI.

## Resume prompt

> Read `docs/editor-testing/todo.md` and `scenario-cookbook.md`. Pick one failing scenario, run `bash scripts/test-keyboard-harness.sh -s NAME --fast`, fix `handleKey` in `build-keywriter.sh`, redeploy Writerdeck, re-run until PASS. Add cookbook scenarios before fixing new bugs.
