# Scenario catalog

All **105** device harness scenarios in implementation-agnostic terms. Each loads a note (usually the shared Norwegian prose fixture into harness-only `z-test-keyboard-harness.md`), performs keystrokes (Mac-style modifiers over the phone/WebSocket path), then asserts caret position, selection range, and document length or content.

**Motion/selection pattern** (reset caret between blocks — never grow-to-N-then-peel):

| Block | Meaning |
|-------|---------|
| uni 1 | one press one way |
| uni 5 | five presses one way |
| bi 1+1 | grow 1, reverse 1 |
| bi 3+5 | grow 3, reverse 5 (overshoot past anchor) |
| bi 7+7 | grow 7, reverse 7 |

Both directions on applicable axes. Fixture includes ≥2 long wrapping paragraphs with æøå, two bullet lists, word line, horizontal line, and 12 equal vertical lines.

**Conventions:** Doc start/end = beginning/end of document. Line start/end = start/end of current logical line (between newlines). Visual line = displayed row; a single logical line may wrap. Vertical Up/Down preserve **visual x** (`positionToRectangle`). Shift+arrow extends selection; plain arrow moves the caret. Alt = word/paragraph; Ctrl = document/line boundaries (Mac). Hardware pages: `pageleft`/`pageright` + `contentY`. Reading overscroll: Esc → preview then page cmds.

Filter critical: `bash scripts/test-keyboard-harness.sh -t critical --fast` (**36**). Authoritative names: `--list`. Implementation: `daemon/cmd/edit-harness/scenarios_*.go`, helpers in `pattern.go`.

## Critical (36)

Must pass for basic editing. Tag: `critical`. Grouped by function; each row is one scenario.

| Group | Scenarios |
|-------|-----------|
| Caret / line home | `load-cursor-at-start`, `home-clears-selection`, `gap-up-at-doc-start` |
| Plain arrows | `down-one-logical-line`, `cm-line-down-basic`, `cm-line-down-last-line`, `wrap-down-one-visual-line`, `wrap-up-from-visual-line-2`, `gap-plain-left-moves-caret`, `gap-plain-right-moves-caret` |
| Doc bounds | `combo-ctrl-home`, `combo-ctrl-end` |
| Shift+arrow select | `shift-right-from-home`, `shift-left-from-end`, `shift-right-after-home-no-stale-anchor`, `shift-down-after-arrow-down`, `shift-up-after-arrow-down`, `shift-left-repeat-from-end`, `shift-left-repeat-mid-doc`, `ctrl-shift-left-select-line`, `gap-collapse-selection-left`, `gap-collapse-selection-right` |
| Backspace / Delete | `bs-plain`, `gap-delete-forward`, `gap-delete-with-selection`, `gap-empty-doc-backspace`, `alt-backspace-deletes-word`, `ctrl-backspace-deletes-line` |
| Insert / replace | `gap-enter-new-line`, `gap-type-replaces-selection`, `gap-select-all` |
| Undo / redo | `undo-redo-len`, `gap-undo-chain`, `gap-redo-shift-ctrl-z` |
| Word nav | `combo-alt-left`, `combo-alt-right` |

Not critical (still valuable): selection shrink on reverse (`shift-down-then-up-shrinks`, `shift-left-then-right-shrinks`, `shift-right-then-left-shrinks`), caret clamps at ends, hardware page scroll (`hw-page-*`), reading overscroll clamp (`read-overscroll-clamps`), goal-column precision, touch tap placement, combo repeat, unicode word boundaries, most wrap/combo permutations.

## Core (9)

| Scenario | Behavior |
|----------|----------|
| `load-cursor-at-start` | After open, caret at document start, no selection, full prose present, edit mode. |
| `home-clears-selection` | Mid horizontal line, End then Shift+Home then Home: selection cleared to line start. |
| `shift-right-from-home` | From horizontal line start, Shift+Right at N=1/3/7 grows selection from fixed anchor. |
| `shift-left-from-end` | From horizontal line end, Shift+Left at N=1/3/7 grows selection backward. |
| `shift-right-after-home-no-stale-anchor` | Repeated Shift+Right keeps anchor at line start through N=7. |
| `shift-left-after-end-no-stale-anchor` | Repeated Shift+Left keeps anchor at line end through N=7. |
| `shift-down-after-arrow-down` | Mid vertical block, Shift+Down at N=1/3/7 extends selection downward. |
| `shift-up-after-arrow-down` | Near bottom of vertical block, Shift+Up at N=1/3/7 extends selection upward. |
| `ctrl-shift-left-select-line` | Shift+Home from horizontal line end selects the whole line. |

## Regression — logical newlines (7)

| Scenario | Behavior |
|----------|----------|
| `down-one-logical-line` | Down at N=1/3/7 across equal-width vertical lines. |
| `up-one-logical-line` | Up at N=1/3/7 (reverse of down). |
| `shift-down-then-up-shrinks` | Shift+Down extends; Shift+Up shrinks at N=1/3/7. |
| `shift-left-repeat-from-end` | Shift+Left at N=1/3/7 from horizontal line end. |
| `alt-backspace-deletes-word` | Alt+Backspace at N=1/3/7 from end of word line. |
| `ctrl-backspace-deletes-line` | Ctrl+Backspace at N=1/3/7 from last vertical line. |
| `shift-left-repeat-mid-doc` | Shift+Left at N=1/3/7 from mid vertical line end. |

## CodeMirror vertical (11)

| Scenario | Behavior |
|----------|----------|
| `cm-line-down-basic` | Down at N=1/3/7 on prose vertical block. |
| `cm-line-up-basic` | Up at N=1/3/7 (reverse). |
| `cm-line-down-shorter` | Down from mid longer line clamps onto shorter next line. |
| `cm-line-up-shorter` | Up from longer line clamps onto shorter previous line. |
| `cm-line-down-last-line` | Down on last line clamps at EOF through N=7. |
| `cm-line-down-goal-col` | Down twice preserves visual x across a short middle line. |
| `cm-select-line-down` | Shift+Down at N=1/3/7 from top of vertical block. |
| `cm-select-line-down-mid` | Shift+Down from mid-line extends to next line. |
| `cm-select-down-up-doc-end` | At EOF, Shift+Down clamps then Shift+Up yields bounded selection. |
| `cm-select-up-basic` | Shift+Up at N=1/3/7 from end of vertical block. |
| `cm-select-up-mid` | Shift+Up from mid last line selects upward. |

## Modifier combos (25)

Word jumps use the shared prose word line (`alfa…juliett`). Doc-bound jumps prove idempotent clamp at N=1/3/7.

| Scenario | Behavior |
|----------|----------|
| `combo-alt-left` / `combo-alt-right` | Word motion at N=1/3/7 both directions. |
| `combo-alt-up` / `combo-alt-down` | Paragraph motion + clamp at N=1/3/7. |
| `combo-ctrl-left` / `combo-ctrl-right` | Doc start/end from mid-prose caret; clamp N=1/3/7. |
| `combo-ctrl-up` / `combo-ctrl-down` | Doc start/end vertical; clamp N=1/3/7. |
| `combo-shift-alt-left` / `combo-shift-alt-left-repeat` | Word select backward at N=1/3/7. |
| `combo-shift-alt-right` / `combo-shift-alt-right-repeat` | Word select forward at N=1/3/7. |
| `combo-shift-alt-up` / `combo-shift-alt-down` | Paragraph select. |
| `combo-shift-ctrl-left` / `combo-shift-ctrl-right` | Line/doc select; clamp N=1/3/7. |
| `combo-shift-ctrl-left-multiline` | Shift+Ctrl+Left on line 2 selects that line only. |
| `combo-shift-ctrl-up` / `combo-shift-ctrl-down` | Whole-doc select; clamp N=1/3/7. |
| `combo-shift-home-line` / `combo-shift-end-line` | Line select via Shift+Home/End. |
| `combo-ctrl-home` / `combo-ctrl-end` | Doc Home/End from mid prose; clamp N=1/3/7. |
| `combo-shift-ctrl-home` / `combo-shift-ctrl-end` | Shift+Ctrl+Home/End from mid two-line doc. |

## Backspace extensions (5)

| Scenario | Behavior |
|----------|----------|
| `bs-alt-word-mid` | Alt+Backspace mid-word on prose word line. |
| `bs-ctrl-line-start` | Ctrl+Backspace at start of a vertical line merges upward. |
| `bs-shift-with-selection` | Shift+Backspace clears a full-line selection. |
| `bs-plain` | Backspace at N=1/3/7 from horizontal line end. |
| `delete-repeat-forward` | Delete at N=1/3/7 from horizontal line start (reverse of bs-plain). |

## Wrapped paragraph (15)

Fixed editor width (320px). Default fixture: `word ` × 40 (specialized geometry). N=3/N=7 visual-line offsets are provisional until re-calibrated on device.

| Scenario | Behavior |
|----------|----------|
| `wrap-down-one-visual-line` | Down at N=1/3/7 through visual rows. |
| `wrap-down-not-jump-paragraph` | Down stays inside the wrapped block. |
| `wrap-up-from-visual-line-2` | Up at N=1/3/7 from deep visual row back to start. |
| `wrap-shift-down-one-visual` | Shift+Down at N=1/3/7. |
| `wrap-shift-down-then-up-shrinks` | Extend then shrink with Shift+Up at N=1/3. |
| `wrap-down-last-visual-line` | Down at EOF clamps through N=7. |
| `wrap-shift-down-last-to-eof` | Shift+Down on last visual row selects through EOF. |
| `wrap-mixed-newline-and-wrap` | Down from short first line into wrapped second line. |
| `wrap-down-goal-column` | Down preserves visual x across a wrap break. |
| `wrap-combo-alt-left-word` / `wrap-combo-alt-right-word` | Alt word motion at N=1/3/7 both ways. |
| `wrap-combo-ctrl-bs-line` | Ctrl+Backspace clears the wrapped logical line. |
| `wrap-shift-left-across-wrap` | Shift+Left at N=1/3/7 across a wrap boundary. |
| `wrap-home-on-visual-line` / `wrap-end-on-visual-line` | Home/End on second visual row. |

## Undo and redo (5)

| Scenario | Behavior |
|----------|----------|
| `undo-redo-len` | Select all, delete, Undo restores, Redo re-applies delete. |
| `undo-cursor-reposition` | Insert at start, Undo from EOF restores caret; Redo restores insert caret. |
| `undo-mid-line-delete` | Delete second line, Undo restores text and caret. |
| `redo-cleared-by-new-edit` | After Undo, a new edit clears the redo stack. |
| `undo-after-select-delete` | Shift+Home select, delete, Undo restores collapsed caret at end. |

## Gap coverage (17)

| Scenario | Behavior |
|----------|----------|
| `gap-up-at-doc-start` | Up at doc start clamps through N=7. |
| `gap-plain-left-moves-caret` / `gap-plain-right-moves-caret` | Plain Left/Right at N=1/3/7 from mid horizontal line. |
| `gap-plain-left-at-doc-start` / `gap-plain-right-at-doc-end` | Clamp at ends through N=7. |
| `gap-collapse-selection-left` / `gap-collapse-selection-right` | Plain arrow collapses a shift selection. |
| `gap-delete-forward` | Delete at N=1/3/7 from mid horizontal line. |
| `gap-delete-with-selection` | Delete clears a mid-line selection. |
| `gap-select-all` | Ctrl+A selects entire prose document. |
| `gap-enter-new-line` | Enter at horizontal line end inserts newline. |
| `gap-type-replaces-selection` | Typing replaces a mid-line selection. |
| `gap-redo-shift-ctrl-z` | Shift+Ctrl+Z redoes after Undo. |
| `gap-undo-chain` | Two Undos restore successive deletions. |
| `gap-unicode-alt-backspace` | Alt+Backspace on `test résumé æøå` leaves `test résumé`. |
| `gap-empty-doc-backspace` | Backspace on empty document is a no-op. |
| `gap-alt-bs-with-selection` | Alt+Backspace with a word selection deletes the selection. |

## Hardware page buttons (2)

| Scenario | Behavior |
|----------|----------|
| `hw-page-right-scrolls-edit` | `pageright` at N=1/3/7 raises `contentY` by ~1500px each; caret stays put. |
| `hw-page-left-scrolls-edit` | After seeding seven page-rights, `pageleft` at N=1/3/7 returns to `contentY` 0. |

## Reading mode (1)

| Scenario | Behavior |
|----------|----------|
| `read-overscroll-clamps` | Esc to preview; after paging to document end, ten extra `pageright` leave `contentY` unchanged; one `pageleft` decreases `contentY`. Tag: `read`. |

## Touch (3)

| Scenario | Behavior |
|----------|----------|
| `touch-down-goal-column` | Tap mid line 1, Down lands at same visual x on line 2. |
| `touch-up-goal-column` | Tap mid line 2, Up lands at same visual x on line 1. |
| `touch-down-shorter-line` | Tap mid longer line, Down clamps to closest x on shorter next line. |

## Selection (shift reverse) (3)

| Scenario | Behavior |
|----------|----------|
| `shift-left-then-right-shrinks` | Pattern left→right: uni1/5, bi1+1, bi3+5 overshoot, bi7+7. |
| `shift-right-then-left-shrinks` | Pattern right→left (partner). |
| `shift-up-then-down-shrinks` | Pattern up→down on vertical prose lines. |

## Sources and notation

CodeMirror/Qt porting notes, marker notation, and “what not to port”: [scenario-cookbook.md](scenario-cookbook.md).
