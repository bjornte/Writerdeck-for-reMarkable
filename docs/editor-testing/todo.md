# TODO: Keyboard editing + harness

**Fresh-agent entry point.** Mac-style editing in Writerdeck (`handleMacArrow`, `handleMacBackspace`, `handleMacEditKeys` in `third_party/keywriter/build-keywriter.sh`). Drive fixes through the device harness — not manual Lobby typing.

Read first: this file, [milestone-runs.md](milestone-runs.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22. Scenario names: [scenario-catalog.md](scenario-catalog.md). Porting sources: [scenario-cookbook.md](scenario-cookbook.md).

Root pointer: [TODO.md](../../TODO.md) item 2.

## Current score (device)

| Milestone | Result | Note |
|-----------|--------|------|
| Latest full suite | **74 / 31** (0 prep) of **105** | `11-38-40` @ `7603357`; report `docs/recon/test-keyboard-harness-2026-07-16T11-38-40.md` |
| Prior full suite | **72 / 33** (0 prep) of **105** | `10-01-42` @ `f42bfbe` |
| **Critical (gate)** | **26 / 10** of **36** | unchanged after `7603357`; tag filter `-t critical` |
| Best pre-rewrite | **89 / 4** (+1 prep) of **94** | `00-37-27` @ `bdccee9` |
| Sign-off gate | **105/105 PASS** | `bash scripts/test-keyboard-harness.sh --fast`, single session |

`test-edit-session.sh` PASS on deploy @ `7603357`. Do not run it in parallel with the keyboard harness — concurrent `/api/open` can kill the editor mid-suite.

## Goal for next session

**Critical pass first:** `bash scripts/test-keyboard-harness.sh -t critical --fast` → **36/36 PASS**. Then one full `--fast` run → update [milestone-runs.md](milestone-runs.md). Full **105/105** is sign-off; do not prune TODO until both pass.

## Critical failures @ `11-38-40` (10 open)

All are QML in `build-keywriter.sh` unless noted. Harness path: WebSocket → `socketRouteKey()` → `handleMacArrow` / `handleMacBackspace`.

| Scenario | Symptom (layman) | QML to inspect |
|----------|------------------|----------------|
| `shift-right-from-home` | Shift+Right then Shift+Left should shrink selection to nothing; reverse grows or leaves 2-char highlight | `handleMacArrow` shift+Left/Right; `query.select` + `cursorPosition` |
| `shift-left-from-end` | Same shrink bug growing left from line end | same |
| `shift-right-after-home-no-stale-anchor` | Same as `shift-right-from-home` | same |
| `shift-down-after-arrow-down` | Shift+Down then Shift+Up should shrink vertical selection; selection collapses or wrong line span | `extendSelectionVertical`, `lineUpForSelection`, `lineDownPos` |
| `shift-up-after-arrow-down` | Shift+Up selection lands short (95 chars want, 19 got) | same |
| `shift-left-repeat-from-end` | Repeated Shift+Left shrink-on-reverse off by 1 | same horizontal shift path |
| `shift-left-repeat-mid-doc` | Same off-by-one on mid-doc line | same |
| `alt-backspace-deletes-word` | Alt+Backspace on prose word line: cursor 6 chars too far right (wrong word boundary) | `deleteWordLeftPos`, `wordLeftPos`, `handleMacBackspace` |
| `ctrl-backspace-deletes-line` | Ctrl+Backspace at line start: stops ~19 chars short of expected line boundary on prose | `deleteLineLeftPos`, `handleMacBackspace` |
| `wrap-up-from-visual-line-2` | Up after Down on wrapped line jumps too far (cursor 80 vs max 65) | `visualLineUpPos`, `lineUpPos`, `goalX` |

Pattern note: horizontal shift failures cluster on **bi 1+1 reverse** (step 11) — grow one char, reverse one char, expect collapsed selection at anchor. Vertical failures cluster on **bi 1+1 / bi 3+5 reverse** after multi-line grow.

## What `7603357` tried (did not fix critical)

Commit `7603357` batched:

- Replaced shift+Left/Right anchor math with “opposite selection end = anchor, move active caret one char” via `query.select(anchorH, newHead)`.
- Removed `lineUpForSelection` special-case `head - 1` at line start.
- Removed Alt+Backspace forward word extension; added Ctrl+Backspace newline swallow.

Result: overall **+2** pass (74 vs 72) but **critical still 26/36**. Do not declare done on partial full-suite movement. Next agent should treat shift shrink as unsolved and try a different model (explicit anchor property, or `select(min,max)` then `cursorPosition = newHead`, or CodeMirror-style anchor/head separate from Qt `selectionStart`/`selectionEnd` ordering).

## Non-critical still open (after critical)

31 full-suite fails total; 21 are non-critical. Largest clusters: `cm-*` goal-col/vertical select, `combo-ctrl-left/right`, wrap shift reverse, `read-overscroll-clamps`, `hw-page-right-scrolls-edit`, dedicated `scenarios_selection.go` shrink aliases. Defer until critical is green.

## Next (one QML batch)

1. Triage: `bash scripts/test-keyboard-harness.sh -t critical --fast` — confirm all 10 failures on current binary (no deploy between `-s` checks).
2. Fix horizontal + vertical shift shrink in one `build-keywriter.sh` diff (`handleMacArrow`, `extendSelectionVertical`, helpers). Read [lessons.md](../lessons.md): `query.cursorPosition` after `query.select()` can collapse selection; `moveCursorSelection` direction enums break vertical.
3. Fix `deleteWordLeftPos` / `deleteLineLeftPos` for prose fixture (Unicode words, logical lines with `\n`).
4. Fix `visualLineUpPos` / wrap-up calibration (`wrap_fixtures.go` offsets if needed).
5. One push → CI → `fetch-keywriter-dist.sh` → `deploy-keywriter.sh -b` → `test-edit-session.sh` → `-t critical --fast` → if 36/36, full `--fast` → update [milestone-runs.md](milestone-runs.md).

Deploy budget: **one** Writerdeck binary deploy per session unless QML fails to launch ([lessons.md](../lessons.md) § Harness batch workflow).

## Do not retry

- Treating keyboard Left/Right as page-scroll (gpio Key_Left confusion).
- One-press / one-direction only, or grow-to-N-then-peel, as done for shrink bugs.
- Separate WebSocket wake after prepare for modified keys.
- End-prime before modified scenarios.
- Nested `invokeMethod` for `socketRouteKey` on GUI thread — deadlock.
- Duplicate `Keys.onPressed` on query TextEdit — crash loop.
- Per-scenario deploy loops.
- Parallel `test-edit-session.sh` + full harness.
- Assuming `7603357` shift logic is correct — device proved otherwise.

## Harness rewrite — complete

105 scenarios; pattern uni1 / uni5 / bi1+1 / bi3+5 / bi7+7. Shared `fixtureProse` in `fixtures.go`. Helpers: `pattern.go`. Critical list: `daemon/cmd/edit-harness/keys.go` → `criticalScenarios`.

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
| `wrap_fixtures.go` | Calibrated wrap offsets (W=320) |

Mode: **sandbox-prepare**. Tags: `-t critical`, `-t hw`, `-t read`, `-t wrap`, `-t undo`. Single scenario: `-s NAME --fast`.

## Acceptance

1. `-t critical --fast` → **36/36 PASS**
2. Full `--fast` → **105/105 PASS**
3. `test-edit-session.sh` PASS
4. `journalctl -u writerdeck -n 30` clean after deploy

## Dev loop

[lessons.md](../lessons.md) § Harness batch workflow. Editor patches live in `third_party/keywriter/build-keywriter.sh` (brittle — see [decisions.md](../decisions.md) §3, [architecture.md](../architecture.md)); long-term: fork migration in [TODO.md](../../TODO.md) item 3.
