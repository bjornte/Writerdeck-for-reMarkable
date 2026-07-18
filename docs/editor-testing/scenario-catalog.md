# Scenario catalog

All **122** device harness scenarios in implementation-agnostic terms. Each loads a note (usually the shared Norwegian prose fixture into harness-only `z-test-keyboard-harness.md`), performs keystrokes (Mac-style modifiers — Ctrl/Alt — over the phone/WebSocket path), then asserts caret position, selection range, and document length or content.

**Motion/selection pattern** (reset caret between blocks — never grow-to-N-then-peel):

| Block | Meaning |
|-------|---------|
| uni 1 | one press one way |
| uni 5 | five presses one way |
| bi 1+1 | grow 1, reverse 1 |
| bi 3+5 | grow 3, reverse 5 (overshoot past anchor) |
| bi 7+7 | grow 7, reverse 7 |

Both directions on applicable axes. Fixture includes ≥2 long wrapping paragraphs with æøå, two bullet lists, word line, horizontal line, and 12 equal vertical lines.

## Dialect (who wins)

Phone Meta maps to Ctrl ([decisions.md](../decisions.md) §3), so harness **Ctrl = Mac ⌘**.

| Rank | Source | Role |
|------|--------|------|
| 1 | **Apple Cocoa** (`StandardKeyBinding.dict`, TextEdit) | Prose chords: ⌘/Option arrows, deletes, paragraph unit |
| 2 | **CodeMirror** `@codemirror/commands` `standardKeymap` | Home/End caret to visual-line ends (Apple Home/End only scroll); wrap-boundary checks |
| 3 | **Qt** `QTextEdit` | Not the target keymap (Ctrl+arrow = word there). Use only to know stock widget defaults |

### Units

Writers use two views of the same note at once (Finseth, *The Craft of Text Editing*, ch. 4 — MIT-hosted). **Meaning** units live in the Markdown file; **layout** units are what you see on the wrapped screen. Mixing them is how “line” tests go green while ⌘→ still jumps a whole paragraph.

| Term | Kind | Meaning |
|------|------|---------|
| Visual line | Layout | One soft-wrapped row on screen. |
| Logical line | Meaning | Text between `\n` characters. |
| Paragraph | Meaning | Apple `paragraphRange`: each `\n`-delimited segment, **including empty lines**. Option+Up/Down. |
| Word | Meaning | Apple `moveWord*` / Option+Left/Right (Latin prose; not CM “group”). |
| Document | Meaning | Whole note. |
| Goal column | Layout | Horizontal spot Up/Down tries to keep across uneven rows. |

Layout-motion scenarios that claim “line” under soft wrap must use a wrap fixture and must not accept caret at paragraph end (`motion_test.go`).

### Chords

| Chord (harness) | Must do | Source |
|-----------------|---------|--------|
| Ctrl+Left / Ctrl+Right | Visual-line start/end; further presses stay if already there | Apple ⌘←/→ |
| Shift+Ctrl+Left / Right | Extend selection with that motion (visual line, not whole wrap-paragraph) | Apple |
| Home / End | Visual-line start/end (not scroll-only; not paragraph end) | CodeMirror |
| Shift+Home / Shift+End | Select to visual-line start/end | CodeMirror (with Home/End) |
| Alt+Left / Alt+Right | Word | Apple |
| Alt+Up / Alt+Down | Paragraph boundary; **empty paragraphs count** | Apple |
| Ctrl+Up / Ctrl+Down, Ctrl+Home / Ctrl+End | Document start/end | Apple ⌘↑/↓ and Mod+Home/End |
| Ctrl+Backspace | Delete to visual-line start (not whole logical line) | Apple ⌘⌫ / CM `deleteLineBoundaryBackward` |
| Alt+Backspace | Delete word backward | Apple |
| Shift + motion | Extend selection with the same unit | Apple / CM |
| Plain Up/Down | One visual row; keep goal column | Apple / CM |

### Crosswalk (Apple / Lexical / harness)

| Unit (Google Selection API) | Apple (Mac) | Lexical-style name | Harness |
|----------------------------|-------------|--------------------|---------|
| `lineboundary` | ⌘← / ⌘→ | `moveToLineBeginning` / `moveToLineEnd` | Ctrl+Left / Ctrl+Right → `wrap-ctrl-*` |
| `lineboundary` | Home / End (we follow CM, not Apple scroll) | — | Home / End → `wrap-home-*` / `wrap-end-*` |
| `word` | ⌥← / ⌥→ | `moveToPrevWord` / next | Alt+Left / Alt+Right |
| `paragraphboundary` | ⌥↑ / ⌥↓ | `moveToParagraphBeginning` / End | Alt+Up / Alt+Down |
| `documentboundary` | ⌘↑ / ⌘↓ | `moveToEditorBeginning` / End | Ctrl+Up / Ctrl+Down |
| (delete lineboundary) | ⌘⌫ | `deleteLineBackward` | Ctrl+Backspace → `wrap-combo-ctrl-bs-line` |

Go helpers: `daemon/cmd/edit-harness/motion.go`. Line-boundary scenarios cannot assert caret at paragraph end (`motion_test.go`).

Short hard-broken fixtures only prove the no-wrap case (visual = logical). Wrap proof for Ctrl+Left/Right is `wrap-ctrl-*`. Vertical Up/Down preserve **visual x**. Hardware pages: `pageleft`/`pageright` + `contentY`. Reading overscroll: Esc → preview then page cmds.

Filter critical: `bash scripts/test-keyboard-harness.sh -t critical --fast` (**42**). Authoritative names: `--list`. Implementation: `daemon/cmd/edit-harness/scenarios_*.go`, helpers in `pattern.go` and `motion.go`.

## Critical (42)

Must pass for basic editing. Tag: `critical`. Live scoreboard: [milestone-runs.md](milestone-runs.md).

Framed as manuscript unit-tasks (Card, Moran & Newell / GOMS — locate, then modify, then check the caret): keyboard path only; no mouse.

| Goal | Group | Scenarios |
|------|-------|-----------|
| Locate (layout) | Visual rows / line ends | `wrap-down-one-visual-line`, `wrap-up-from-visual-line-2`, `wrap-ctrl-left`, `wrap-ctrl-right`, `gap-plain-left-moves-caret`, `gap-plain-right-moves-caret` |
| Locate (meaning) | Char / hard line / word / doc | `down-one-logical-line`, `cm-line-down-basic`, `cm-line-down-last-line`, `combo-alt-left`, `combo-alt-right`, `combo-ctrl-home`, `combo-ctrl-end`, `gap-up-at-doc-start` |
| Locate + select | Shift extend / collapse | `load-cursor-at-start`, `home-clears-selection`, `shift-right-from-home`, `shift-left-from-end`, `shift-right-after-home-no-stale-anchor`, `shift-down-after-arrow-down`, `shift-up-after-arrow-down`, `shift-left-repeat-from-end`, `shift-left-repeat-mid-doc`, `ctrl-shift-left-select-line`, `gap-collapse-selection-left`, `gap-collapse-selection-right`, `gap-shift-down-mid-wrapping-paras`, `gap-shift-up-mid-wrapping-paras` |
| Modify | Type / delete / clipboard | `bs-plain`, `gap-delete-forward`, `gap-delete-with-selection`, `gap-empty-doc-backspace`, `alt-backspace-deletes-word`, `ctrl-backspace-deletes-line`, `gap-enter-new-line`, `gap-type-replaces-selection`, `gap-select-all`, `gap-copy-paste`, `gap-cut-paste` |
| Verify | Undo / redo | `undo-redo-len`, `gap-undo-chain`, `gap-redo-shift-ctrl-z` |

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
| `combo-alt-up-double-blank` / `combo-alt-down-double-blank` | Empty `\n` paragraphs are real stops (Apple); must not skip blanks to the far side in one press. |
| `combo-alt-up-prose-double-blank` | Alt+Up into trailing double-blank section; next press leaves the blank band toward earlier text. |
| `combo-ctrl-left` / `combo-ctrl-right` | Hard-`\n` only (`en\nto\ntre`): line start/end where visual equals logical. Wrap proof is `wrap-ctrl-*`. |
| `combo-ctrl-up` / `combo-ctrl-down` | Doc start/end vertical; pattern + clamp. |
| `combo-shift-alt-left` / `combo-shift-alt-left-repeat` | Word select backward; pattern. |
| `combo-shift-alt-left-after-type` | After Shift-select + type, Shift+Alt+Left re-anchors (no stale head). |
| `combo-shift-left-after-type` | Same type-then-nav guard for plain Shift+Left. |
| `combo-shift-alt-right` / `combo-shift-alt-right-repeat` | Word select forward; pattern. |
| `combo-shift-alt-up` / `combo-shift-alt-down` | Paragraph select. |
| `combo-shift-ctrl-left` / `combo-shift-ctrl-right` | Line select on one-line docs; pattern + clamp. |
| `combo-shift-ctrl-left-multiline` | Shift+Ctrl+Left on line 2 selects that line only. |
| `combo-shift-ctrl-up` / `combo-shift-ctrl-down` | Whole-doc select; pattern + clamp. |
| `combo-shift-home-line` / `combo-shift-end-line` | Line select via Shift+Home/End (CodeMirror Home/End family). |
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

## Wrapped paragraph (19)

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
| `wrap-combo-ctrl-bs-line` | Mid visual row: Ctrl+Backspace deletes to that row’s start only (not whole paragraph). |
| `wrap-shift-left-across-wrap` | Pattern Shift+Left grow, Shift+Right shrink across wrap boundary. |
| `wrap-home-on-visual-line` / `wrap-end-on-visual-line` | Home/End on second visual row; End lands at next wrap point (`wrapEndVisualRow1`), not paragraph end. |
| `wrap-ctrl-left` / `wrap-ctrl-right` | Mid visual row: Ctrl+Left/Right to that row’s ends; further presses stay. Critical. |
| `wrap-shift-ctrl-left` / `wrap-shift-ctrl-right` | Same motion with Shift: select to visual-line end only. |

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
