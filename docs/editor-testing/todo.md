# TODO: Keyboard editing + harness

Mac-style editing in Writerdeck (`handleMacArrow`, `handleMacBackspace`, `handleMacEditKeys` in `build-keywriter.sh`). Drive fixes through the device harness — not manual Lobby typing.

Read: [scenario-cookbook.md](scenario-cookbook.md), [scenario-catalog.md](scenario-catalog.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22. Scores: [milestone-runs.md](milestone-runs.md).

Root pointer only: [TODO.md](../../TODO.md) item 2.

## Current score (device)

| Milestone | Result | Note |
|-----------|--------|------|
| Best full suite (pre-rewrite) | **89 / 4** (+1 prep) of **94** | `00-37-27` @ `bdccee9` |
| Rewrite baseline (102) | **85 / 17** (0 prep) of **102** | `01-54-49`; edit-session PASS |
| Sign-off gate | **105/105 PASS** | `--fast`, single session |

### Harness rewrite — confirmed complete (no device run yet)

All applicable scenario files under `daemon/cmd/edit-harness/scenarios_*.go` are rewritten. Suite size **102** (was 94). Do not treat the pre-rewrite score as current.

What landed:

1. **`read-overscroll-clamps`** — Esc to preview, page past EOF, ten extra downs must keep `contentY` unchanged, one up decreases it. Tag `read`. Runner capture: `CaptureContentY` / `ExpectContentYEqCaptured` / `ExpectContentYLtCaptured`.
2. **Omnidirectional + N=1/3/7** on motion, selection, shrink/extend, page, and combo cases (unary edits keep clear setup).
3. **Shared prose fixture** (`fixtureProse` → harness-only `z-test-keyboard-harness.md`): Norwegian æøå / accents, two bullet lists, long horizontal + 12 uniform vertical lines + word line. Specialized Content kept for wrap (W=320), empty doc, goal-column shapes, tall page bodies.
4. **Unicode-safe lengths** — prepare/`TextLen` expects use `editorLen` (BMP/QString), not Go byte `len()`.

### Pattern (motion / selection / page / word-nav)

Applicable scenarios use **SetCursor reset between blocks** so each case is pure:

| Block | Meaning |
|-------|---------|
| uni 1 | one press one way |
| uni 5 | five presses one way |
| bi 1+1 | grow 1, reverse 1 |
| bi 3+5 | grow 3, reverse 5 (intentional **overshoot** past anchor) |
| bi 7+7 | grow 7, reverse 7 |

Both directions (Left↔Right, Up↔Down, page-right↔page-left). Helpers: `pattern.go`.

Shared `fixtureProse`: at least two long wrapping paragraphs (Norwegian æøå), two bullet lists, word line, horizontal line, 12 equal vertical lines.

### Remaining product bugs (from last device run) — fix after re-baseline

- Rewrite baseline `01-54-49`: **85/17/0** of 102 (old 1/3/7 stacking). Next full run re-baselines **105**.
- Shift reverse still grows (no pure 1+1 / 3+5 yet until this pattern rewrite lands on device)
- `cm-line-down-goal-col`, read Esc→preview (`read-overscroll-clamps`), provisional wrap offsets, Ctrl+Left/Right semantics

Hardware page cmds proven via inject; physical gpio still needs exclusive grab ([handoff](../todo-handoff-physical-home-input.md)).

## Next (device verify)

1. `bash scripts/test-edit-session.sh`
2. Full `bash scripts/test-keyboard-harness.sh --fast` → update [milestone-runs.md](milestone-runs.md); sign-off **105/105**
3. Recalibrate wrap offsets if needed
4. Product QML fixes — triage once, batch fix, one Writerdeck deploy

## Do not retry

- Treating keyboard Left/Right as page-scroll (gpio Key_Left confusion — [lessons.md](../lessons.md)).
- One-press / one-direction coverage, or grow-to-N-then-peel, as “done” for walk-back bugs.
- Separate WebSocket wake after prepare for modified keys.
- End-prime before modified scenarios — wiped text / jumped to EOF.
- Nested `invokeMethod` for `socketRouteKey` on the GUI thread — deadlock.
- Duplicate `Keys.onPressed` on query TextEdit — QML crash loop.
- Per-scenario deploy loops — triage once, batch fix, one deploy ([lessons.md](../lessons.md) § Harness batch workflow).

## Harness inventory (post-pattern: 105)

| File | Block |
|------|--------|
| `scenarios.go` | Core (9) |
| `scenarios_regression.go` | Regression (7) |
| `scenarios_cm.go` | CodeMirror vertical (10) |
| `scenarios_combo.go` | Alt/Ctrl / Shift combos (25) |
| `scenarios_bs.go` | Backspace / delete (5) |
| `scenarios_wrap.go` | Wrapped paragraph (15) |
| `scenarios_undo.go` | Undo/redo (5) |
| `scenarios_gaps.go` | Gap coverage (19) |
| `scenarios_hw.go` | Hardware page buttons (2) |
| `scenarios_read.go` | Reading-mode overscroll (1) |
| `scenarios_touch.go` | Touch → visual goal-x (3) |
| `scenarios_selection.go` | Shift reverse (3) |
| `fixtures.go` / `pattern.go` | Prose + uni/bi step helpers |
| `main.go`, `report.go` | Runner, contentY capture, `showlobby` teardown |
| `wrap_fixtures.go` | Calibrated wrap offsets (W=320; N=3/7 provisional) |

Mode: **sandbox-prepare**. Tags: `-t critical`, `-t hw`, `-t read`, `-t wrap`, `-t undo`.

## Acceptance

Full suite **105/105 PASS** with `--fast`. Pure bi 1+1 / 3+5 / 7+7 on applicable axes. `read-overscroll-clamps` PASS after QML clamp.

## Dev loop

[lessons.md](../lessons.md) § Harness batch workflow.
