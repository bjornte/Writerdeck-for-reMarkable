# TODO: Keyboard editing + harness

Mac-style editing in Writerdeck (`handleMacArrow`, `handleMacBackspace`, `handleMacEditKeys` in `build-keywriter.sh`). Drive fixes through the device harness — not manual Lobby typing.

Read: [scenario-cookbook.md](scenario-cookbook.md), [scenario-catalog.md](scenario-catalog.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22. Scores: [milestone-runs.md](milestone-runs.md).

Root pointer only: [TODO.md](../../TODO.md) item 2.

## Current score (device)

| Milestone | Result | Note |
|-----------|--------|------|
| Best full suite (pre-rewrite) | **89 / 4** (+1 prep) of **94** | `00-37-27` @ `bdccee9` |
| Rewrite baseline (102) | **85 / 17** (0 prep) of **102** | `01-54-49` |
| Pattern baseline (105) | **72 / 33** (0 prep) of **105** | `10-01-42` @ `f42bfbe`; edit-session PASS |
| Critical @ pattern baseline | **26 / 10** of **36** | same run; see failures below |
| Sign-off gate | **105/105 PASS** | `--fast`, single session |

Report: `docs/recon/test-keyboard-harness-2026-07-16T10-01-42.md`.

## Harness rewrite — complete

All applicable scenario files under `daemon/cmd/edit-harness/scenarios_*.go` use the pattern below. Suite size **105** (was 94 pre-rewrite). Helpers: `pattern.go`. Shared `fixtureProse` loads into harness-only `z-test-keyboard-harness.md` — Norwegian æøå, two long wrapping paragraphs, two bullet lists, word line, horizontal line, twelve uniform vertical lines.

### Motion pattern (reset caret between blocks)

| Block | Meaning |
|-------|---------|
| uni 1 | one press one way |
| uni 5 | five presses one way |
| bi 1+1 | grow 1, reverse 1 (must shrink / return) |
| bi 3+5 | grow 3, reverse 5 (intentional overshoot past anchor) |
| bi 7+7 | grow 7, reverse 7 |

Both directions on applicable axes (Left↔Right, Up↔Down, page-right↔page-left). No grow-to-N-then-peel stacking.

### Remaining failures @ `10-01-42` (layman)

**Critical (10 of 36 failed)** — basic editing bar not met:

- Shift+arrow selection: grow right/left from line start/end, then reverse to shrink (`shift-right-from-home`, `shift-left-from-end`, `shift-right-after-home-no-stale-anchor`, `shift-left-repeat-from-end`, `shift-left-repeat-mid-doc`). Reverse direction grows the highlight instead of shrinking — including pure bi 1+1 and bi 3+5.
- Shift+Up/Down after moving vertically (`shift-down-after-arrow-down`, `shift-up-after-arrow-down`). Same shrink-on-reverse bug on logical lines.
- Alt+Backspace and Ctrl+Backspace on the long prose note (`alt-backspace-deletes-word`, `ctrl-backspace-deletes-line`).
- Wrapped paragraph: Up after Down does not return to the right place (`wrap-up-from-visual-line-2`).

**Non-critical clusters** (still open): `shift-down-then-up-shrinks` and dedicated selection-shrink scenarios; `cm-line-down-goal-col`; combo Ctrl+Left/Right on prose; wrap shift reverse; `read-overscroll-clamps` (Esc does not enter preview); `hw-page-right-scrolls-edit` cursor drift; provisional wrap N=3/N=7 offsets.

Hardware page cmds proven via harness inject; physical gpio still needs exclusive grab ([handoff](../todo-handoff-physical-home-input.md)).

## Next (product fixes)

1. Triage critical cluster: Shift reverse-walk shrink (horizontal + vertical) — one QML batch.
2. Alt/Ctrl+Backspace word/line delete on prose fixture.
3. Wrap Up-after-Down calibration; read-mode Esc→preview; Ctrl+Left/Right doc jump semantics if still failing on `-t critical` subset.
4. One Writerdeck deploy, then `bash scripts/test-keyboard-harness.sh --fast` → update [milestone-runs.md](milestone-runs.md); sign-off **105/105**.

## Do not retry

- Treating keyboard Left/Right as page-scroll (gpio Key_Left confusion — [lessons.md](../lessons.md)).
- One-press / one-direction coverage, or grow-to-N-then-peel, as “done” for walk-back bugs.
- Separate WebSocket wake after prepare for modified keys.
- End-prime before modified scenarios — wiped text / jumped to EOF.
- Nested `invokeMethod` for `socketRouteKey` on the GUI thread — deadlock.
- Duplicate `Keys.onPressed` on query TextEdit — QML crash loop.
- Per-scenario deploy loops — triage once, batch fix, one deploy ([lessons.md](../lessons.md) § Harness batch workflow).

## Harness inventory (105)

| File | Block |
|------|--------|
| `scenarios.go` | Core (9) |
| `scenarios_regression.go` | Regression (7) |
| `scenarios_cm.go` | CodeMirror vertical (11) |
| `scenarios_combo.go` | Alt/Ctrl / Shift combos (25) |
| `scenarios_bs.go` | Backspace / delete (5) |
| `scenarios_wrap.go` | Wrapped paragraph (15) |
| `scenarios_undo.go` | Undo/redo (5) |
| `scenarios_gaps.go` | Gap coverage (19) |
| `scenarios_hw.go` | Hardware page buttons (2) |
| `scenarios_read.go` | Reading-mode overscroll (1) |
| `scenarios_touch.go` | Touch → visual goal-x (3) |
| `scenarios_selection.go` | Shift reverse (3) |
| `fixtures.go` / `pattern.go` | Prose + pattern helpers |
| `main.go`, `report.go` | Runner, contentY capture, `showlobby` teardown |
| `wrap_fixtures.go` | Calibrated wrap offsets (W=320; provisional multi-step) |

Mode: **sandbox-prepare**. Tags: `-t critical`, `-t hw`, `-t read`, `-t wrap`, `-t undo`.

## Acceptance

Full suite **105/105 PASS** with `--fast`. Critical **36/36 PASS**. Pure bi 1+1 / 3+5 / 7+7 on applicable axes. `read-overscroll-clamps` PASS after QML clamp.

## Dev loop

[lessons.md](../lessons.md) § Harness batch workflow.
