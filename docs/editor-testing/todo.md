# TODO: Keyboard editing + harness

Mac-style editing in Writerdeck (`handleMacArrow`, `handleMacBackspace`, `handleMacEditKeys` in `build-keywriter.sh`). Drive fixes through the device harness — not manual Lobby typing.

Read: [scenario-cookbook.md](scenario-cookbook.md), [scenario-catalog.md](scenario-catalog.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22. Scores: [milestone-runs.md](milestone-runs.md).

Root pointer only: [TODO.md](../../TODO.md) item 2.

## Current score (device)

| Milestone | Result | Note |
|-----------|--------|------|
| Best full suite (pre-rewrite) | **89 / 4** (+1 prep) of **94** | `00-37-27` @ `bdccee9` |
| Sign-off gate | **102/102 PASS** after rewrite lands | `--fast`, single session |

### Harness rewrite — confirmed complete (no device run yet)

All applicable scenario files under `daemon/cmd/edit-harness/scenarios_*.go` are rewritten. Suite size **102** (was 94). Do not treat the pre-rewrite score as current.

What landed:

1. **`read-overscroll-clamps`** — Esc to preview, page past EOF, ten extra downs must keep `contentY` unchanged, one up decreases it. Tag `read`. Runner capture: `CaptureContentY` / `ExpectContentYEqCaptured` / `ExpectContentYLtCaptured`.
2. **Omnidirectional + N=1/3/7** on motion, selection, shrink/extend, page, and combo cases (unary edits keep clear setup).
3. **Shared prose fixture** (`fixtureProse` → harness-only `z-test-keyboard-harness.md`): Norwegian æøå / accents, two bullet lists, long horizontal + 12 uniform vertical lines + word line. Specialized Content kept for wrap (W=320), empty doc, goal-column shapes, tall page bodies.
4. **Unicode-safe lengths** — prepare/`TextLen` expects use `editorLen` (BMP/QString), not Go byte `len()`.

### Remaining product bugs (from last device run @ `bdccee9`) — fix after device re-baseline

- Shift+Down then Shift+Up shrink (`shift-down-then-up-shrinks`)
- Arrow Down goal column across short line (`cm-line-down-goal-col`)
- Shift-select near EOF (`cm-select-down-up-doc-end`)
- Unicode word-delete prepare flake (`gap-unicode-alt-backspace`) — length fix may help; re-check on device
- Option+Backspace with selection (`gap-alt-bs-with-selection`)
- **New expected fail until QML clamp:** `read-overscroll-clamps`

Hardware page cmds were proven via harness inject; physical gpio page buttons still need exclusive grab ([handoff](../todo-handoff-physical-home-input.md)).

## Next (device verify — only after this confirmation)

1. `bash scripts/test-edit-session.sh`
2. Full `bash scripts/test-keyboard-harness.sh --fast` → update [milestone-runs.md](milestone-runs.md); sign-off **102/102**
3. Recalibrate wrap N=3/N=7 offsets if provisional values drift (`wrap_fixtures.go`)
4. Product QML fixes for remaining failures + read overscroll clamp — triage once, batch fix, one Writerdeck deploy

## Do not retry

- Treating keyboard Left/Right as page-scroll (gpio Key_Left confusion — [lessons.md](../lessons.md)).
- One-press / one-direction coverage as “done” for walk-back or repeat bugs.
- Separate WebSocket wake after prepare for modified keys.
- End-prime before modified scenarios — wiped text / jumped to EOF.
- Nested `invokeMethod` for `socketRouteKey` on the GUI thread — deadlock.
- Duplicate `Keys.onPressed` on query TextEdit — QML crash loop.
- Per-scenario deploy loops — triage once, batch fix, one deploy ([lessons.md](../lessons.md) § Harness batch workflow).

## Harness inventory (post-rewrite: 102)

| File | Block |
|------|--------|
| `scenarios.go` | Core (9) |
| `scenarios_regression.go` | Regression `\n` (7) |
| `scenarios_cm.go` | CodeMirror vertical (11) |
| `scenarios_combo.go` | Alt/Ctrl / Shift combos (25) |
| `scenarios_bs.go` | Backspace / delete (5) |
| `scenarios_wrap.go` | Wrapped paragraph (15) |
| `scenarios_undo.go` | Undo/redo (5) |
| `scenarios_gaps.go` | Gap coverage (17) |
| `scenarios_hw.go` | Hardware page buttons (2) |
| `scenarios_read.go` | Reading-mode overscroll (1) |
| `scenarios_touch.go` | Touch → visual goal-x (3) |
| `scenarios_selection.go` | Shift reverse (2) |
| `fixtures.go` | Shared prose + tall/wrap helpers |
| `main.go`, `report.go` | Runner, contentY capture, `showlobby` teardown |
| `wrap_fixtures.go` | Calibrated wrap offsets (W=320; N=3/7 provisional) |

Mode: **sandbox-prepare**. Hardware pages: cmd inject + `contentY`, not Arrow keys. Tags: `-t critical`, `-t hw`, `-t read`, `-t wrap`, `-t undo`.

## Acceptance (post-rewrite)

Full suite **102/102 PASS** with `--fast`, single session, clean `journalctl`. `test-edit-session.sh` PASS. `read-overscroll-clamps` PASS. Omnidirectional + 1/3/7 + prose fixtures as rewritten.

## Dev loop

[lessons.md](../lessons.md) § Harness batch workflow.
