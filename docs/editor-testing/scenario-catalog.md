# Scenario catalog

All **83** device harness scenarios in implementation-agnostic terms. Each loads a note, performs keystrokes (Mac-style modifiers over the phone/WebSocket path), then asserts caret position, selection range, and document length or content.

**Conventions:** Doc start/end = beginning/end of document. Line start/end = start/end of current logical line (between newlines). Visual line = displayed row; a single logical line may wrap. Shift+arrow extends selection; plain arrow without Shift moves the caret. Alt = word/paragraph motion; Ctrl = document or line boundaries (Mac editing model).

Authoritative names: `bash scripts/test-keyboard-harness.sh --list`. Implementation: `daemon/cmd/edit-harness/scenarios_*.go`.

## Core (8)

| Scenario | Behavior |
|----------|----------|
| `load-cursor-at-start` | After open, caret at document start, no selection, full text present, edit mode. |
| `home-clears-selection` | End then Home: caret at line start, selection cleared. |
| `shift-right-from-home` | From doc start, Shift+Right ×3 selects characters 0–3 (anchor at start). |
| `shift-left-from-end` | From doc end, Shift+Left ×3 selects last three characters (caret stays at end). |
| `shift-right-after-home-no-stale-anchor` | Two Shift+Right steps from start grow one continuous selection from anchor 0. |
| `shift-down-after-arrow-down` | Down ×2 to line 3, then Shift+Down: selection extends down one logical line. |
| `shift-up-after-arrow-down` | Down ×2 to line 3, then Shift+Up: selection extends up to line 2. |
| `ctrl-shift-left-select-line` | Shift+Home from end of line selects the whole line. |

## Regression — logical newlines (6)

| Scenario | Behavior |
|----------|----------|
| `down-one-logical-line` | Down once from first line moves to start of second line (`aa` / `bb`). |
| `shift-down-then-up-shrinks` | Shift+Down extends selection; Shift+Up shrinks it (goal-column shrink). |
| `shift-left-repeat-from-end` | Shift+Left ×3 from end of line selects the last three characters. |
| `alt-backspace-deletes-word` | Alt+Backspace at end of `hello world` leaves `hello`. |
| `ctrl-backspace-deletes-line` | Ctrl+Backspace on second line removes that line back to the newline. |
| `shift-left-repeat-mid-doc` | Shift+Left ×3 from end of second line selects three characters correctly. |

## CodeMirror vertical — explicit newlines (9)

| Scenario | Behavior |
|----------|----------|
| `cm-line-down-basic` | Down from doc start moves to start of second logical line. |
| `cm-line-down-shorter` | Down from column 2 on a longer first line lands at column 2 on a shorter second line (goal column). |
| `cm-line-down-last-line` | Down on the last line moves caret to document end. |
| `cm-line-down-goal-col` | Down twice preserves goal column across a short middle line. |
| `cm-select-line-down` | Shift+Down from doc start selects through end of first line. |
| `cm-select-line-down-mid` | Shift+Down from mid-line extends selection to target on next line. |
| `cm-select-down-up-doc-end` | At doc end, Shift+Down then Shift+Up yields a bounded upward selection. |
| `cm-select-up-basic` | Shift+Up from doc end selects upward one logical line. |
| `cm-select-up-mid` | Shift+Up from mid line 3 selects upward to line 2 boundary. |

## Modifier combos — Alt, Ctrl, Shift (22)

| Scenario | Behavior |
|----------|----------|
| `combo-alt-left` | Alt+Left from end moves back one word. |
| `combo-alt-right` | Alt+Right from start moves forward one word. |
| `combo-alt-up` | Alt+Up from end of second paragraph moves to doc start. |
| `combo-alt-down` | Alt+Down from start moves to start of next paragraph. |
| `combo-ctrl-left` | Ctrl+Left from end moves to doc start. |
| `combo-ctrl-right` | Ctrl+Right from start moves to doc end. |
| `combo-ctrl-up` | Ctrl+Up from doc end moves to doc start. |
| `combo-ctrl-down` | Ctrl+Down from start moves to doc end. |
| `combo-shift-alt-left` | Shift+Alt+Left from end selects back one word. |
| `combo-shift-alt-right` | Shift+Alt+Right from start selects forward one word. |
| `combo-shift-alt-up` | Shift+Alt+Up from doc end selects to previous paragraph start. |
| `combo-shift-alt-down` | Shift+Alt+Down from start selects through next paragraph. |
| `combo-shift-ctrl-left` | Shift+Ctrl+Left from end selects to doc start. |
| `combo-shift-ctrl-right` | Shift+Ctrl+Right from start selects to doc end. |
| `combo-shift-ctrl-up` | Shift+Ctrl+Up from doc end selects to doc start. |
| `combo-shift-ctrl-down` | Shift+Ctrl+Down from start selects to doc end. |
| `combo-shift-home-line` | Shift+Home on line 2 selects from line start to caret. |
| `combo-shift-end-line` | Shift+End on line 2 selects from caret to line end. |
| `combo-ctrl-home` | Ctrl+Home on line 2 moves to doc start. |
| `combo-ctrl-end` | Ctrl+End moves to doc end. |
| `combo-shift-ctrl-home` | Shift+Ctrl+Home on line 2 selects from doc start through line start. |
| `combo-shift-ctrl-end` | Shift+Ctrl+End on line 2 selects from line start through line end. |

## Backspace extensions (4)

| Scenario | Behavior |
|----------|----------|
| `bs-alt-word-mid` | Alt+Backspace mid-word deletes the word before the caret. |
| `bs-ctrl-line-start` | Ctrl+Backspace at start of line 2 deletes line 1 and newline. |
| `bs-shift-with-selection` | Shift+Backspace with a full-line selection clears document. |
| `bs-plain` | Backspace ×2 from end deletes last two characters. |

## Wrapped paragraph (14)

Fixed editor width (320px). Default fixture: one long unbroken paragraph (`word ` × 40) unless noted.

| Scenario | Behavior |
|----------|----------|
| `wrap-down-one-visual-line` | Down once moves to second visual row, not next paragraph. |
| `wrap-down-not-jump-paragraph` | Down does not skip past the wrapped block in one step. |
| `wrap-up-from-visual-line-2` | Up from second visual row returns to doc start. |
| `wrap-shift-down-one-visual` | Shift+Down selects through first visual-line break. |
| `wrap-shift-down-then-up-shrinks` | After moving down two visual rows, Shift+Down then Shift+Up shrinks selection. |
| `wrap-down-last-visual-line` | Down on last visual row stays at document end. |
| `wrap-shift-down-last-to-eof` | Shift+Down at last visual row selects through document end. |
| `wrap-mixed-newline-and-wrap` | Down from short first line enters wrapped tail of second logical line. |
| `wrap-down-goal-column` | Down from column 2 preserves goal column across wrap. |
| `wrap-combo-alt-left-word` | Alt+Left on wrapped paragraph moves back within text (word nav). |
| `wrap-combo-ctrl-bs-line` | Ctrl+Backspace on wrapped paragraph clears entire logical line. |
| `wrap-shift-left-across-wrap` | Shift+Left ×3 from second visual row selects backward across wrap. |
| `wrap-home-on-visual-line` | Home on second visual row moves to that visual row’s start. |
| `wrap-end-on-visual-line` | End on second visual row moves to that visual row’s end. |

## Undo and redo (5)

| Scenario | Behavior |
|----------|----------|
| `undo-redo-len` | Delete all, Undo restores text, Redo restores deletion state. |
| `undo-cursor-reposition` | Insert at doc start, Undo from end restores text and caret to insert point; Redo restores insert caret. |
| `undo-mid-line-delete` | Delete line 2, Undo restores line and caret position. |
| `redo-cleared-by-new-edit` | After Undo, a new edit clears the redo stack (Redo has no effect). |
| `undo-after-select-delete` | Select-all via Shift+Home, delete, Undo restores text and collapsed caret at end. |

## Gap coverage (15)

| Scenario | Behavior |
|----------|----------|
| `gap-up-at-doc-start` | Up at doc start leaves caret at start. |
| `gap-plain-left-no-cursor-move` | Plain Left at doc end does not move caret (page scroll only). |
| `gap-plain-right-no-cursor-move` | Plain Right at doc start does not move caret (page scroll only). |
| `gap-collapse-selection-left` | Left arrow collapses a backward selection to the near end. |
| `gap-collapse-selection-right` | Right arrow collapses a forward selection to the far end. |
| `gap-delete-forward` | Delete removes character after caret. |
| `gap-delete-with-selection` | Delete with selection clears selected text. |
| `gap-select-all` | Ctrl+A selects entire document. |
| `gap-enter-new-line` | Enter at end inserts newline. |
| `gap-type-replaces-selection` | Typing with selection replaces selection with typed character. |
| `gap-redo-shift-ctrl-z` | Shift+Ctrl+Z redoes after Undo. |
| `gap-undo-chain` | Two Undos restore successive deletions. |
| `gap-unicode-alt-backspace` | Alt+Backspace respects Unicode word boundaries (`résumé` → `test`). |
| `gap-empty-doc-backspace` | Backspace on empty document is a no-op. |
| `gap-alt-bs-with-selection` | Alt+Backspace with word selection deletes selection, leaves prior word. |

## Touch (goal column after tap)

| Scenario | Behavior |
|----------|----------|
| `touch-down-goal-column` | Tap mid line 1 (harnessSetCursor), Down lands same column on line 2. |
| `touch-up-goal-column` | Tap mid line 2, Up lands same column on line 1. |
| `touch-down-shorter-line` | Tap col 2 on longer line, Down clamps column on shorter next line. |

## Selection (shift reverse)

| Scenario | Behavior |
|----------|----------|
| `shift-left-then-right-shrinks` | Shift+Left extends selection, Shift+Right shrinks toward anchor (does not grow). |

## Sources and notation

CodeMirror/Qt porting notes, marker notation, and “what not to port”: [scenario-cookbook.md](scenario-cookbook.md).
