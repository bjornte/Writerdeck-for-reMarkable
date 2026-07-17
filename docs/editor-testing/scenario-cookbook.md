# Scenario cookbook

Portable keyboard/selection cases; full inventory of running tests: [scenario-catalog.md](scenario-catalog.md). Sources and notation below. Implementation: `daemon/cmd/edit-harness/scenarios_*.go`.

Scintilla `test/unit/testSelection.cxx` tests internal selection structs only — not useful here.

## Notation

CodeMirror uses inline markers in one string:

- `|` — cursor (no selection)
- `<` … `>` — selection anchor to head (exclusive of markers in output)

Harness uses byte offsets in plain UTF-8 (ASCII in these cases). `\n` is one character. After each key step, assert `cursor`, `selStart`, `selEnd`, `textLen`, optional `selLen`.

Writerdeck maps Meta to Ctrl for Mac/Linux-style shortcuts over the phone path ([decisions.md](../decisions.md) §3). Chords are Control/Alt (USB Linux keyboards work as-is; Mac ⌘ is not required).

## Port template

```go
{
    Name:    "cookbook-short-name",
    Content: "one\ntwo\nthree",
    Steps: []Step{
        {Keys: []Key{{Name: "Home", Ctrl: true}}},
        {Keys: []Key{{Name: "ArrowDown"}}},
        {Expect: &StateExpect{Cursor: intp(4), SelStart: intp(4), SelEnd: intp(4)}},
    },
},
```

Run: `bash scripts/test-keyboard-harness.sh -s cookbook-short-name --fast`

## Already in `scenarios_regression.go`

| Harness name | Cookbook source | Notes |
|--------------|-----------------|-------|
| `down-one-logical-line` | CodeMirror `cursorLineDown` | `aa\nbb`, down once |
| `shift-down-then-up-shrinks` | CodeMirror `selectLineDown` + `selectLineUp` | visual-x shrink |
| `shift-left-repeat-from-end` | — | Writerdeck bug report |
| `alt-backspace-deletes-word` | CodeMirror `deleteGroupBackward` | Mac Alt+Backspace |
| `ctrl-backspace-deletes-line` | CodeMirror `deleteLineBoundaryBackward` | Mac Mod+Backspace |

## Undo / redo (in `scenarios_undo.go`)

Sources: Qt `tst_qplaintextedit` (`undoRedo`, `undoRedoShouldRepositionTextEditCursor`), CodeMirror `test-history.ts`, Ace redo-stack clearing.

| Harness name | Source | Notes |
|--------------|--------|-------|
| `undo-redo-len` | Qt `undoRedo` | delete all, Ctrl+Z, Ctrl+Y |
| `undo-cursor-reposition` | Qt `undoRedoShouldRepositionTextEditCursor` | type at start, undo/redo cursor |
| `undo-mid-line-delete` | — | delete line2, undo restores cursor |
| `redo-cleared-by-new-edit` | Ace #7024 | edit after undo kills redo |
| `undo-after-select-delete` | CodeMirror selection restore | shift+home, delete, undo |

Single-character steps use the WebSocket text path (`Key{Name:"a"}`). Run: `bash scripts/test-keyboard-harness.sh -s undo-redo-len --fast`.

## CodeMirror — vertical motion (port next)

Content uses `\n` for line breaks. Positions counted from 0.

| Suggested name | Content | Keys | Expected (cursor, selStart, selEnd) | CodeMirror test |
|----------------|---------|------|--------------------------------------|-----------------|
| `cm-line-down-basic` | `one\ntwo` | Ctrl+Home, Down | (4, 4, 4) | start of `two` |
| `cm-line-down-shorter` | `one\nt` | cursor 2, Down | (4, 4, 4) | closest x on shorter line |
| `cm-line-down-last-line` | `one` | cursor 2, Down | (3, 3, 3) | end of doc on last line |
| `cm-line-down-goal-col` | `tre\ni\nfemte` | cursor 2, Down×2 | (7–8) | visual x sticky across short middle line |
| `cm-select-line-down` | `one\ntwo\nthree` | pos 0, Shift+Down | (4, 0, 4) | sel 0–4 |
| `cm-select-line-down-mid` | `one\ntwo\nthree` | pos 2, Shift+Down | (7, 2, 7) | mid-line extend down |
| `cm-select-down-up-doc-end` | `one\ntwo\nthree` | pos 11, Shift+Down, Shift+Up | see `shift-down-then-up-shrinks` | visual x at EOF |
| `cm-select-up-basic` | `one\ntwo\nthree` | End, Shift+Up | (11, 4, 11) | select upward one line |
| `cm-select-up-mid` | `one\ntwo\nthree` | pos 9, Shift+Up | (9, 4, 9) | mid-line extend up |

## CodeMirror — backspace / delete (Mac bindings)

From `commands.ts` keymap: Alt+Backspace = word backward; Mod+Backspace = line backward.

| Suggested name | Content | Keys | Expect textLen / cursor |
|----------------|---------|------|-------------------------|
| `cm-alt-bs-word` | `hello world` | End, Alt+Backspace | len 5, cursor 5 (`hello`) |
| `cm-mod-bs-line` | `line1\nline2` | End, Ctrl+Backspace | len 6 (`line1\n`) |
| `cm-shift-bs` | `abcd` | End, Shift+Backspace | Qt: deletes selection only — define after select step |

## Qt tst_qplaintextedit — worth porting

| Qt test | Scenario idea | Harness sketch |
|---------|---------------|----------------|
| `createSelection` | Ctrl+Home, Shift+End on `Hello World` | sel 0–11 |
| `backspace` | type then Backspace×3 | textLen decreases |
| `shiftBackspace` | select then Shift+Backspace | clears selection |
| `undoRedo` | type, Ctrl+Z, Ctrl+Y | textLen round-trip — needs buffer snapshot or len only |
| `undoRedoShouldRepositionTextEditCursor` | edit mid-doc, undo | cursor returns — high value for Ctrl+Z bug |
| `shiftDownInLineLastShouldSelectToEnd` | wrapped `Foo\nBar` | Shift+Down on last visual line selects to EOF — needs fixed widget width or long unbroken line |

Wrapped-line cases use harness `Width` (320px) and calibrated byte offsets in `daemon/cmd/edit-harness/wrap_fixtures.go` via `harnessprepare`.

## Priority order

Fix open clusters in [todo.md](todo.md) first: Shift+Ctrl+Left multiline, shift+Alt repeat, visual goal-x on device, undo cursor restore, remaining wrap/CM selection. Most cookbook blocks are already in `scenarios_*.go`; add a scenario only when a reported bug has no matching name in `bash scripts/test-keyboard-harness.sh --list`.

## What not to port

- Scintilla unit tests (wrong abstraction)
- Vim/emacs internal tests
- CodeMirror folded-line / decoration tests (no fold UI in Writerdeck)
- USB-only qmap cases
