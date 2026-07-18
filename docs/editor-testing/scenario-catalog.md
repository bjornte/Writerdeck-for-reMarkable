# Scenario catalog

All **118** device harness scenarios in implementation-agnostic terms. Each loads a note (usually the shared Norwegian prose fixture into harness-only `z-test-keyboard-harness.md`), performs keystrokes (Mac/Linux-style modifiers — Ctrl/Alt — over the phone/WebSocket path), then asserts caret position, selection range, and document length or content.

**Motion/selection pattern** (reset caret between blocks — never grow-to-N-then-peel):

| Block | Meaning |
|-------|---------|
| uni 1 | one press one way |
| uni 5 | five presses one way |
| bi 1+1 | grow 1, reverse 1 |
| bi 3+5 | grow 3, reverse 5 (overshoot past anchor) |
| bi 7+7 | grow 7, reverse 7 |

Both directions on applicable axes. Fixture includes ≥2 long wrapping paragraphs with æøå, two bullet lists, word line, horizontal line, and 12 equal vertical lines.

**Conventions:** Doc start/end = beginning/end of document. Line start/end = start/end of current logical line (between newlines). Visual line = displayed row; a single logical line may wrap. Vertical Up/Down preserve **visual x** (`positionToRectangle`). Shift+arrow extends selection; plain arrow moves the caret. Alt = word/paragraph; Ctrl/Cmd+Left/Right = line; Ctrl/Cmd+Up/Down and Ctrl+Home/End = document (Mac). Hardware pages: `pageleft`/`pageright` + `contentY`. Reading overscroll: Esc → preview then page cmds.

Filter critical: `bash scripts/test-keyboard-harness.sh -t critical --fast` (**40**). Authoritative names: `--list`. Implementation: `daemon/cmd/edit-harness/scenarios_*.go`, helpers in `pattern.go`.

## Critical (40)

Must pass for basic editing. Tag: `critical`. Grouped by function; each row is one scenario. Live scoreboard: [milestone-runs.md](milestone-runs.md).

| Group | Scenarios |
|-------|-----------|
| Caret / line home | `load-cursor-at-start`, `home-clears-selection`, `gap-up-at-doc-start` |
| Plain arrows | `down-one-logical-line`, `cm-line-down-basic`, `cm-line-down-last-line`, `wrap-down-one-visual-line`, `wrap-up-from-visual-line-2`, `gap-plain-left-moves-caret`, `gap-plain-right-moves-caret` |
| Doc bounds | `combo-ctrl-home`, `combo-ctrl-end` |
| Shift+arrow select | `shift-right-from-home`, `shift-left-from-end`, `shift-right-after-home-no-stale-anchor`, `shift-down-after-arrow-down`, `shift-up-after-arrow-down`, `shift-left-repeat-from-end`, `shift-left-repeat-mid-doc`, `ctrl-shift-left-select-line`, `gap-collapse-selection-left`, `gap-collapse-selection-right`, `gap-shift-down-mid-wrapping-paras`, `gap-shift-up-mid-wrapping-paras` |
| Backspace / Delete | `bs-plain`, `gap-delete-forward`, `gap-delete-with-selection`, `gap-empty-doc-backspace`, `alt-backspace-deletes-word`, `ctrl-backspace-deletes-line` |
| Insert / replace | `gap-enter-new-line`, `gap-type-replaces-selection`, `gap-select-all` |
| Clipboard | `gap-copy-paste`, `gap-cut-paste` |
| Undo / redo | `undo-redo-len`, `gap-undo-chain`, `gap-redo-shift-ctrl-z` |
| Word nav | `combo-alt-left`, `combo-alt-right` |

Not critical (still valuable): selection shrink on short lines (`shift-*-shrinks`), cross-paragraph Shift (`gap-shift-*-across-para-break`), mid-column short-line Shift (`gap-shift-down-mid-short-lines`), caret clamps at ends, hardware page scroll (`hw-page-*`), reading overscroll (`read-overscroll-clamps`), goal-column precision, touch tap placement, combo repeat, unicode word boundaries, most wrap/combo permutations.

## Core (9)

| Scenario | Behavior |
|----------|----------|
| `load-cursor-at-start` | After open, caret at document start, no selection, full prose present, edit mode. |
| `home-clears-selection` | Mid horizontal line, End then Shift+Home then Home: selection cleared to line start. |
| `shift-right-from-home` | Pattern uni1/5, bi1+1, bi3+5, bi7+7: Shift+Right grow then Shift+Left shrink from mid horizontal line. |
| `shift-left-from-end` | Same pattern from near horizontal line end. |
| `shift-right-after-home-no-stale-anchor` | Same as shift-right-from-home (duplicate coverage). |
| `shift-left-after-end-no-stale-anchor` | Same as shift-left-from-end (duplicate coverage). |
| `shift-down-after-arrow-down` | Pattern on vertical block: Shift+Down grow, Shift+Up shrink. |
| `shift-up-after-arrow-down` | Pattern near bottom of vertical block: Shift+Up grow, Shift+Down shrink. |
| `ctrl-shift-left-select-line` | Shift+Home from horizontal line end selects the whole line. |

## Regression — logical newlines (7)

| Scenario | Behavior |
|----------|----------|
| `down-one-logical-line` | Pattern uni1/5, bi1+1, bi3+5, bi7+7: Down on vertical block. |
| `up-one-logical-line` | Reverse: Up on vertical block. |
| `shift-down-then-up-shrinks` | Pattern: Shift+Down grow, Shift+Up shrink (vertical). |
| `shift-left-repeat-from-end` | Pattern: Shift+Left grow, Shift+Right shrink from near line end. |
| `alt-backspace-deletes-word` | uni1 and uni5 Alt+Backspace from word-line seed. |
| `ctrl-backspace-deletes-line` | uni1 and uni5 Ctrl+Backspace from last vertical line. |
| `shift-left-repeat-mid-doc` | Pattern: Shift+Left grow, Shift+Right shrink mid vertical line. |

## CodeMirror vertical (11)

| Scenario | Behavior |
|----------|----------|
| `cm-line-down-basic` | Pattern on vertical block (Down). |
| `cm-line-up-basic` | Pattern on vertical block (Up). |
| `cm-line-down-shorter` | Down from mid longer line clamps onto shorter next line. |
| `cm-line-up-shorter` | Up from longer line clamps onto shorter previous line. |
| `cm-line-down-last-line` | Down on last line clamps at EOF through uni5. |
| `cm-line-down-goal-col` | Down twice preserves visual x across a short middle line. |
| `cm-select-line-down` | Pattern: Shift+Down grow, Shift+Up shrink from top of vertical block. |
| `cm-select-line-down-mid` | Shift+Down from mid-sentence wrapping prose (one visual row band). |
| `cm-select-up-mid` | Shift+Up from mid-sentence wrapping prose (partner). |
| `cm-select-down-up-doc-end` | At EOF, Shift+Down clamps then Shift+Up yields bounded selection. |
| `cm-select-up-basic` | Pattern: Shift+Up grow, Shift+Down shrink from end of vertical block. |

## Modifier combos (29)

Word jumps use pattern blocks on the prose word line. Doc-bound jumps use uni1/uni5 idempotent clamp where applicable.

| Scenario | Behavior |
|----------|----------|
| `combo-alt-left` / `combo-alt-right` | Pattern word motion both directions. |
| `combo-alt-up` / `combo-alt-down` | Paragraph motion (current boundary, then prev/next); pattern + clamp. |
| `combo-alt-up-double-blank` / `combo-alt-down-double-blank` | Same across two consecutive blank lines. |
| `combo-alt-up-prose-double-blank` | Alt+Up across trailing double-blank section in shared prose fixture. |
| `combo-ctrl-left` / `combo-ctrl-right` | Line start/end (Mac Cmd); pattern + clamp. |
| `combo-ctrl-up` / `combo-ctrl-down` | Doc start/end vertical; pattern + clamp. |
| `combo-shift-alt-left` / `combo-shift-alt-left-repeat` | Word select backward; pattern. |
| `combo-shift-alt-left-after-type` | After Shift-select + type, Shift+Alt+Left re-anchors (no stale head). |
| `combo-shift-left-after-type` | Same type-then-nav guard for plain Shift+Left. |
| `combo-shift-alt-right` / `combo-shift-alt-right-repeat` | Word select forward; pattern. |
| `combo-shift-alt-up` / `combo-shift-alt-down` | Paragraph select. |
| `combo-shift-ctrl-left` / `combo-shift-ctrl-right` | Line/doc select; pattern + clamp. |
| `combo-shift-ctrl-left-multiline` | Shift+Ctrl+Left on line 2 selects that line only. |
| `combo-shift-ctrl-up` / `combo-shift-ctrl-down` | Whole-doc select; pattern + clamp. |
| `combo-shift-home-line` / `combo-shift-end-line` | Line select via Shift+Home/End. |
| `combo-ctrl-home` / `combo-ctrl-end` | Doc Home/End from mid prose; pattern + clamp. |
| `combo-shift-ctrl-home` / `combo-shift-ctrl-end` | Shift+Ctrl+Home/End from mid two-line doc. |

## Backspace extensions (5)

| Scenario | Behavior |
|----------|----------|
| `bs-alt-word-mid` | Alt+Backspace mid-word on prose word line. |
| `bs-ctrl-line-start` | Ctrl+Backspace at start of a vertical line merges upward. |
| `bs-shift-with-selection` | Shift+Backspace clears a full-line selection. |
| `bs-plain` | Pattern: Backspace from horizontal line end. |
| `delete-repeat-forward` | Pattern: Delete from horizontal line start (reverse of bs-plain). |

## Wrapped paragraph (15)

Fixed editor width (320px). Default fixture: `word ` × 40 (specialized geometry). Multi-step visual-line offsets in `wrap_fixtures.go` are provisional until re-calibrated on device.

| Scenario | Behavior |
|----------|----------|
| `wrap-down-one-visual-line` | Pattern Down through visual rows. |
| `wrap-down-not-jump-paragraph` | Down stays inside the wrapped block. |
| `wrap-up-from-visual-line-2` | Pattern Up from deep visual row back to start. |
| `wrap-shift-down-one-visual` | Pattern Shift+Down grow, Shift+Up shrink. |
| `wrap-shift-down-then-up-shrinks` | Mid-sentence on wrapped block: full uni/bi Shift+Down grow / Shift+Up reverse. |
| `wrap-down-last-visual-line` | Down at EOF clamps through uni5. |
| `wrap-shift-down-last-to-eof` | Shift+Down on last visual row selects through EOF. |
| `wrap-mixed-newline-and-wrap` | Down from short first line into wrapped second line. |
| `wrap-down-goal-column` | Down preserves visual x across a wrap break. |
| `wrap-combo-alt-left-word` / `wrap-combo-alt-right-word` | Pattern Alt word motion both ways. |
| `wrap-combo-ctrl-bs-line` | Ctrl+Backspace clears the wrapped logical line. |
| `wrap-shift-left-across-wrap` | Pattern Shift+Left grow, Shift+Right shrink across wrap boundary. |
| `wrap-home-on-visual-line` / `wrap-end-on-visual-line` | Home/End on second visual row. |

## Undo and redo (5)

| Scenario | Behavior |
|----------|----------|
| `undo-redo-len` | Select all, delete, Undo restores, Redo re-applies delete. |
| `undo-cursor-reposition` | Insert at start, Undo from EOF restores caret; Redo restores insert caret. |
| `undo-mid-line-delete` | Delete second line, Undo restores text and caret. |
| `redo-cleared-by-new-edit` | After Undo, a new edit clears the redo stack. |
| `undo-after-select-delete` | Shift+Home select, delete, Undo restores collapsed caret at end. |

## Gap coverage (26)

| Scenario | Behavior |
|----------|----------|
| `gap-up-at-doc-start` | Up at doc start clamps through uni5. |
| `gap-plain-left-moves-caret` / `gap-plain-right-moves-caret` | Pattern Left/Right from mid horizontal line. |
| `gap-plain-left-in-paragraph` / `gap-plain-right-in-paragraph` | Pattern Left/Right inside long wrapping paragraph. |
| `gap-plain-left-at-doc-start` / `gap-plain-right-at-doc-end` | Clamp at ends through uni5. |
| `gap-collapse-selection-left` / `gap-collapse-selection-right` | Plain arrow collapses a shift selection. |
| `gap-delete-forward` | Pattern Delete from mid horizontal line. |
| `gap-delete-with-selection` | Delete clears a mid-line selection. |
| `gap-select-all` | Ctrl+A selects entire prose document. |
| `gap-copy-paste` | Ctrl+C then Ctrl+V duplicates selected text at the caret (in-editor clipboard). Critical. |
| `gap-cut-paste` | Ctrl+X removes selection; Ctrl+V inserts it elsewhere. Critical. |
| `gap-enter-new-line` | Enter at horizontal line end inserts newline. |
| `gap-type-replaces-selection` | Typing replaces a mid-line selection. |
| `gap-redo-shift-ctrl-z` | Shift+Ctrl+Z redoes after Undo. |
| `gap-undo-chain` | Two Undos restore successive deletions. |
| `gap-unicode-alt-backspace` | Alt+Backspace on `test résumé æøå` leaves `test résumé`. |
| `gap-empty-doc-backspace` | Backspace on empty document is a no-op. |
| `gap-alt-bs-with-selection` | Alt+Backspace with a word selection deletes the selection. |
| `gap-shift-down-mid-wrapping-paras` | Mid-sentence Shift+Down grow / Shift+Up reverse across wrapping paragraphs (uni1/5, bi1+1/3+5/7+7). Critical. |
| `gap-shift-up-mid-wrapping-paras` | Mid-sentence Shift+Up grow / Shift+Down reverse (partner). Critical. |
| `gap-shift-down-across-para-break` | Near end of wrapping para1, Shift+Down enters para2; reverse shrinks. |
| `gap-shift-up-across-para-break` | Early in wrapping para2, Shift+Up reaches para1; reverse shrinks. |
| `gap-shift-down-mid-short-lines` | Mid-column Shift+Down/Up on equal-width hard-newline lines (full pattern). |

## Hardware page buttons (2)

| Scenario | Behavior |
|----------|----------|
| `hw-page-right-scrolls-edit` | Pattern `pageright`: raises `contentY` by ~85% of viewport height each; caret stays put. |
| `hw-page-left-scrolls-edit` | After seeding seven page-rights, pattern `pageleft` returns to `contentY` 0. |
| `hw-page-step-shrinks-in-landscape` | Landscape page step is smaller than portrait (viewport-relative, not fixed 1500px). |

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
