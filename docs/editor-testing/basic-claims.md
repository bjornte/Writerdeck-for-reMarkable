# Basic editing claims — kill-test inventory

Standing status: [decisions.md](../decisions.md) **Typing-test strategy is failing**. Method: [methodology-shortcomings.md](methodology-shortcomings.md). Critical scenarios: [scenario-catalog.md](scenario-catalog.md).

This is the strategy surface for “basic editing works.” Scenario names are evidence, not the claim. **Guarded** means we have named a dumb broken behavior and today’s harness would go red for it. **Partial** means some wrong forms fail but a known user-facing miss can stay green. **Unguarded** means critical green is possible while the claim fails for a person.

Update this file whenever a human finds a green-suite basic bug, or when a kill-test is added/proven. Same change should touch [methodology-shortcomings.md](methodology-shortcomings.md) if a new miss pattern appears.

How to raise a row: invent the broken behavior → make a red check (fixture + assert + observation that separates good from bad) → then fix the editor.

Statuses below are honest as of 18 Jul 2026 — not aspirational.

---

## Locate

| Claim | Kill-test (broken behavior that must go red) | Critical evidence | Status |
|-------|-----------------------------------------------|-------------------|--------|
| Caret moves one character left/right | Left/Right no-ops or jumps by word | `gap-plain-left-moves-caret`, `gap-plain-right-moves-caret` | Guarded |
| Up at document start does not escape or wipe | Up at offset 0 changes mode, clears text, or moves caret | `gap-up-at-doc-start` | Guarded |
| Down moves one hard-broken line | Down jumps to doc end or stays put on a multi-line block | `down-one-logical-line`, `cm-line-down-basic`, `cm-line-down-last-line` | Guarded (hard `\n` only) |
| Up/Down on soft-wrapped prose move one **visual** row | Up/Down jump a whole wrap-paragraph, or step by a tall caret rect | `wrap-down-one-visual-line`, `wrap-up-from-visual-line-2` | Guarded for big jumps |
| Up/Down keep **goal column** (visual x) across uneven/wrap rows | Down from mid-row snaps to column 0 of the next row; Up does not restore col | `wrap-down-goal-column`, `cm-line-down-goal-col` | Guarded — Down exact cursor + Up round-trip to col 2 |
| Cmd+Left/Right stop at **this visual row’s** start/end under soft wrap | Cmd+Right jumps to paragraph end; or lands at the wrap index but caret is drawn on the **next** row | `wrap-ctrl-left`, `wrap-ctrl-right`, `wrap-ctrl-right-then-left`, `wrap-end-then-up`, `wrap-end-then-down`, `wrap-home-on-visual-line`, `wrap-end-on-visual-line` (+ `assoc` / `caretY`) | Guarded — index alone was a false green; round-trip Left/Up/Down + painted caretY catch “next row” |
| Option+Left/Right move by word | Option+Left acts like plain Left or Cmd+Left | `combo-alt-left`, `combo-alt-right` | Guarded (simple word fixtures) |
| Option+Up/Down move by paragraph; blank lines are stops | Option+Up skips empty paragraphs to doc start in one press | `combo-alt-up`, `combo-alt-down`, `combo-alt-up-double-blank`, `combo-alt-down-double-blank`, `combo-alt-up-prose-double-blank` | Guarded |
| Cmd+Home / Cmd+End go to document start/end | Cmd+Home only goes to line start; Cmd+End stays mid-doc | `combo-ctrl-home`, `combo-ctrl-end` | Guarded |
| After open, caret at start, edit mode, prose present | Empty buffer, wrong mode, or caret at EOF | `load-cursor-at-start` | Guarded |

Still basic for writers but not critical: short-line `combo-ctrl-left/right` (visual=logical only).

## Select

| Claim | Kill-test | Critical evidence | Status |
|-------|-----------|-------------------|--------|
| Shift+Left/Right grow and shrink selection on a line | Shift extends then reverse does not shrink; stale anchor after Home | `shift-right-from-home`, `shift-left-from-end`, `shift-right-after-home-no-stale-anchor`, `shift-left-repeat-*` | Guarded for those patterns |
| Home clears an active selection to line start | Home leaves selection or jumps elsewhere | `home-clears-selection` | Guarded (short horizontal line) |
| Shift+Home selects to line start | Selects whole paragraph or nothing | `ctrl-shift-left-select-line` | Partial — hard-line Shift+Home; wrap visual-line select is separate |
| Shift+Cmd+Left/Right select to **visual** line ends under wrap | Selects whole wrap-paragraph | `wrap-shift-ctrl-left`, `wrap-shift-ctrl-right` | Guarded — SelLen exact to visual row slice |
| Shift+Up/Down grow/shrink on hard-broken lines | Shift+Down selects to EOF; reverse does not shrink | `shift-down-after-arrow-down`, `shift-up-after-arrow-down` | Guarded (hard `\n`) |
| Shift+Up/Down on wrapping paragraphs keep a sane selection | Selection jumps to wrong paragraph or collapses oddly mid-wrap | `gap-shift-down-mid-wrapping-paras`, `gap-shift-up-mid-wrapping-paras` | Partial — geometry-sensitive; weak expects possible |
| Plain Left/Right collapse selection then move | Collapse moves the wrong end or clears text | `gap-collapse-selection-left`, `gap-collapse-selection-right` | Guarded |

## Modify

| Claim | Kill-test | Critical evidence | Status |
|-------|-----------|-------------------|--------|
| Backspace deletes one character | No-op or deletes a word | `bs-plain` | Guarded |
| Delete forward deletes one character | No-op or deletes selection wrongly | `gap-delete-forward` | Guarded |
| Backspace/Delete with selection removes the selection | Leaves selection or deletes neighbors | `gap-delete-with-selection` | Guarded |
| Backspace on empty doc is safe | Crash, mode flip, or invents text | `gap-empty-doc-backspace` | Guarded |
| Option+Backspace deletes previous word | Deletes a character or whole line | `alt-backspace-deletes-word` | Guarded (hard-line seed) |
| Cmd+Backspace deletes to **line** start | Deletes whole paragraph on wrap; or only one character | `ctrl-backspace-deletes-line`, `wrap-combo-ctrl-bs-line` | Guarded — hard-line + wrap fixture (delete mid visual row → only those chars; wipe-paragraph fails textLen) |
| Return inserts a newline | Inserts nothing or breaks undo | `gap-enter-new-line` | Guarded |
| Typing replaces an active selection | Inserts beside selection | `gap-type-replaces-selection` | Guarded |
| Select all / copy / cut / paste round-trip | Clipboard ops no-op or wipe note | `gap-select-all`, `gap-copy-paste`, `gap-cut-paste` | Guarded (phone path) |

## Verify (undo)

| Claim | Kill-test | Critical evidence | Status |
|-------|-----------|-------------------|--------|
| Undo/redo restore length and basic caret | Undo no-op; redo wrong; length drift | `undo-redo-len`, `gap-undo-chain`, `gap-redo-shift-ctrl-z` | Guarded for length/chain; caret nuance weaker |

---

## Summary (why strategy still fails)

- Several wrap Locate/Modify claims and paragraph Option+Up/Down are now Guarded with discriminating checks. That does **not** lift the strategy-failing banner — green still is not proof basic editing works, and some Select claims stay Partial.
- Rows marked Guarded still assume the fixture creates the case; do not weaken fixtures when editing scenarios.

Count of critical scenarios: 57. Do not equate that number with “all writer claims guarded.”

## Next raise (priority)

1. When the next human green-suite miss appears: encode it here first, then kill-test, then fix — do not only promote existing scenarios.
2. Soft-wrap Shift+Up/Down expects remain Partial (geometry-sensitive).
