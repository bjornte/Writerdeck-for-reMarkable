# TODO: Keyboard editing + harness

Mac-style editing in Writerdeck (`handleKey` / `handleMacArrow` in `build-keywriter.sh`). Drive fixes through the device keyboard harness — not manual Lobby typing.

Read: [scenario-cookbook.md](scenario-cookbook.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22.

## Not done — do not sign off without these

Thirteen harness scenarios pass on device (`scenarios.go` + `scenarios_regression.go`), but they use explicit `\n` line breaks, not wrapped paragraphs. That is not the same as the reported bugs.

Still open:

- **Wrapped paragraphs** — Down/Up and Shift+Down/Up must move/select by visual line (`positionToRectangle`), not `lineDownPos` newline math. No harness scenario yet; cookbook priority #4.
- **Shift+Alt / Shift+Ctrl arrows** — no scenarios, not device-verified (QML has code paths; unknown if correct).
- **Undo/redo** — five scenarios in `scenarios_undo.go`; not device-verified.
- **Auto-scroll on e-ink** — last wrapped line typing feel; manual spot-check after wrap logic lands.

Original reported bugs (still the acceptance bar):

- Alt+Backspace word delete — harness PASS on flat line only
- Ctrl+Backspace line delete — harness PASS on `\n` lines only
- Arrow Down in a **wrapped** multiline paragraph — not tested
- Shift+Down then Shift+Up shrink on wrapped lines — not tested
- Repeated Shift+Left — harness PASS on single line
- Ctrl+Z undo — scenarios added, not verified

## Harness

| Piece | Path |
|-------|------|
| Core | `daemon/cmd/edit-harness/scenarios.go` |
| Regressions | `daemon/cmd/edit-harness/scenarios_regression.go` |
| Undo | `daemon/cmd/edit-harness/scenarios_undo.go` |
| Runner | `daemon/cmd/edit-harness/main.go` |
| Shell | `scripts/test-keyboard-harness.sh` |
| QML | `third_party/keywriter/build-keywriter.sh` |

Wrapped-line scenarios need pinned `query.width` in harness setup or long unbroken content — see cookbook § Wrapped-line Qt cases.

## Dev loop

[lessons.md](../lessons.md) § Harness batch workflow. Triage once (`--unit`, `--fast --hard-reset`), batch fixes, one deploy max. Partial harness green is not sign-off.

## Acceptance

Full suite PASS with `--fast --hard-reset`, including wrapped-line and modifier-combo scenarios not yet written. `test-edit-session.sh` PASS. Clean `journalctl`.

## Resume prompt

> Read this file. Add wrapped-line and Shift+Alt/Ctrl scenarios first, then run full harness. Do not mark keyboard editing done until wrap and modifier combos pass on device.
