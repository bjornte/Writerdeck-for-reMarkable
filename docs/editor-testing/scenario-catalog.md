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

Not critical (still valuable): selection shrink (`shift-down-then-up-shrinks`, `shift-left-then-right-shrinks`, `shift-right-then-left-shrinks`, `shift-up-then-down-shrinks`), caret clamps at ends, hardware page scroll (`hw-page-*`), reading overscroll (`read-overscroll-clamps`), goal-column precision, touch tap placement, combo repeat, unicode word boundaries, most wrap/combo permutations.

### Open on device @ `10-01-42` (`f42bfbe`)

**26 / 36 critical PASS.** Ten critical scenarios failed:

| Scenario | Layman |
|----------|--------|
| `shift-right-from-home` | Shift+Right then reverse does not shrink cleanly (bi 1+1 / 3+5) |
| `shift-left-from-end` | Shift+Left then reverse — same |
| `shift-right-after-home-no-stale-anchor` | Same pattern on stale-anchor variant |
| `shift-down-after-arrow-down` | Shift+Down then Shift+Up shrink on vertical lines |
| `shift-up-after-arrow-down` | Shift+Up then Shift+Down shrink |
| `shift-left-repeat-from-end` | Shift+Left grow then reverse |
| `shift-left-repeat-mid-doc` | Shift+Left mid-doc then reverse |
| `alt-backspace-deletes-word` | Alt+Backspace word delete on prose |
| `ctrl-backspace-deletes-line` | Ctrl+Backspace line delete on prose |
| `wrap-up-from-visual-line-2` | Up after Down on wrapped paragraph |

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
| `cm-select-line-down-mid` | Shift+Down from mid-line extends to next line. |
| `cm-select-down-up-doc-end` | At EOF, Shift+Down clamps then Shift+Up yields bounded selection. |
| `cm-select-up-basic` | Pattern: Shift+Up grow, Shift+Down shrink from end of vertical block. |
| `cm-select-up-mid` | Shift+Up from mid last line selects upward. |

## Modifier combos (25)

Word jumps use pattern blocks on the prose word line. Doc-bound jumps use uni1/uni5 idempotent clamp where applicable.

| Scenario | Behavior |
|----------|----------|
| `combo-alt-left` / `combo-alt-right` | Pattern word motion both directions. |
| `combo-alt-up` / `combo-alt-down` | Paragraph motion; pattern + clamp. |
| `combo-ctrl-left` / `combo-ctrl-right` | Doc start/end from mid-prose caret; pattern + clamp. |
| `combo-ctrl-up` / `combo-ctrl-down` | Doc start/end vertical; pattern + clamp. |
| `combo-shift-alt-left` / `combo-shift-alt-left-repeat` | Word select backward; pattern. |
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
| `wrap-shift-down-then-up-shrinks` | Pattern extend then shrink on wrapped block. |
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

## Gap coverage (19)

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
| `hw-page-right-scrolls-edit` | Pattern `pageright`: raises `contentY` by ~1500px each; caret stays put. |
| `hw-page-left-scrolls-edit` | After seeding seven page-rights, pattern `pageleft` returns to `contentY` 0. |

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
