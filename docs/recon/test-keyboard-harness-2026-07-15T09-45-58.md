# Keyboard harness results

Run: 2026-07-15T09:45:59+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 67 pass, 16 fail, 1 prepare fail; total 59.6s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.9 | no | — | — |
| home-clears-selection | pass | 0.6 | no | — | — |
| shift-right-from-home | pass | 0.6 | no | — | — |
| shift-left-from-end | pass | 0.7 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 0.7 | no | — | — |
| shift-down-after-arrow-down | pass | 0.7 | no | — | — |
| shift-up-after-arrow-down | pass | 0.8 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.5 | no | — | — |
| down-one-logical-line | pass | 0.5 | no | — | — |
| shift-down-then-up-shrinks | fail | 0.8 | no | — | step 7: cursor want 17 got 12; selEnd want 17 got 12; selLen want 5 got 0; state={12 12 12 23 line1 line2 line3 line4 1 0 z-test-keyboard-h… |
| shift-left-repeat-from-end | pass | 0.6 | no | — | — |
| alt-backspace-deletes-word | pass | 0.7 | no | — | — |
| ctrl-backspace-deletes-line | pass | 0.7 | no | — | — |
| shift-left-repeat-mid-doc | pass | 0.7 | no | — | — |
| cm-line-down-basic | pass | 0.5 | no | — | — |
| cm-line-down-shorter | fail | 0.8 | no | — | step 5: cursor want 4 got 5; selStart want 4 got 5; selEnd want 4 got 5; state={5 5 5 5 one t 1 0 z-test-keyboard-harness.md} |
| cm-line-down-last-line | pass | 0.6 | no | — | — |
| cm-line-down-goal-col | fail | 0.7 | no | — | step 4: cursor want 6 got 11; selStart want 6 got 11; selEnd want 6 got 11; state={11 11 11 11 one t three 1 0 z-test-keyboard-harness.md} |
| cm-select-line-down | pass | 0.5 | no | — | — |
| cm-select-line-down-mid | pass | 0.8 | no | — | — |
| cm-select-down-up-doc-end | fail | 0.8 | no | — | step 5: selStart want 8 got 7; selLen want 4 got 5; state={12 7 12 13 one two three 1 0 z-test-keyboard-harness.md} |
| cm-select-up-basic | pass | 1.0 | no | — | — |
| cm-select-up-mid | pass | 0.7 | no | — | — |
| combo-alt-left | pass | 0.7 | no | — | — |
| combo-alt-right | pass | 0.5 | no | — | — |
| combo-alt-up | pass | 0.5 | no | — | — |
| combo-alt-down | pass | 0.5 | no | — | — |
| combo-ctrl-left | pass | 0.5 | no | — | — |
| combo-ctrl-right | pass | 0.5 | no | — | — |
| combo-ctrl-up | pass | 0.5 | no | — | — |
| combo-ctrl-down | pass | 0.5 | no | — | — |
| combo-shift-alt-left | pass | 0.5 | no | — | — |
| combo-shift-alt-left-repeat | fail | 0.5 | no | — | step 3: selStart want 0 got 6; selLen want 11 got 5; state={11 6 11 11 hello world 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-right | pass | 0.5 | no | — | — |
| combo-shift-alt-right-repeat | fail | 0.5 | no | — | step 2: selStart want 0 got 6; selLen want 11 got 5; state={11 6 11 11 hello world 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-up | pass | 0.6 | no | — | — |
| combo-shift-alt-down | pass | 0.5 | no | — | — |
| combo-shift-ctrl-left | pass | 0.5 | no | — | — |
| combo-shift-ctrl-left-multiline | pass | 0.7 | no | — | — |
| combo-shift-ctrl-right | pass | 0.5 | no | — | — |
| combo-shift-ctrl-up | pass | 0.7 | no | — | — |
| combo-shift-ctrl-down | pass | 0.5 | no | — | — |
| combo-shift-home-line | pass | 0.7 | no | — | — |
| combo-shift-end-line | pass | 0.6 | no | — | — |
| combo-ctrl-home | pass | 0.5 | no | — | — |
| combo-ctrl-end | pass | 0.5 | no | — | — |
| combo-shift-ctrl-home | pass | 0.5 | no | — | — |
| combo-shift-ctrl-end | pass | 0.6 | no | — | — |
| bs-alt-word-mid | pass | 0.8 | no | — | — |
| bs-ctrl-line-start | pass | 0.5 | no | — | — |
| bs-shift-with-selection | pass | 0.6 | no | — | — |
| bs-plain | pass | 0.6 | no | — | — |
| wrap-down-one-visual-line | pass | 0.6 | no | — | — |
| wrap-down-not-jump-paragraph | pass | 0.7 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 0.9 | no | — | — |
| wrap-shift-down-one-visual | pass | 0.6 | no | — | — |
| wrap-shift-down-then-up-shrinks | pass | 1.0 | no | — | — |
| wrap-down-last-visual-line | pass | 0.8 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.9 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.7 | no | — | — |
| wrap-down-goal-column | pass | 1.0 | no | — | — |
| wrap-combo-alt-left-word | pass | 1.0 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 1.1 | no | — | — |
| wrap-shift-left-across-wrap | pass | 0.9 | no | — | — |
| wrap-home-on-visual-line | pass | 0.9 | no | — | — |
| wrap-end-on-visual-line | pass | 0.8 | no | — | — |
| undo-redo-len | fail | 0.9 | no | — | step 6: cursor want 5 got 0; textLen want 5 got 0; state={0 0 0 0  1 0 z-test-keyboard-harness.md} |
| undo-cursor-reposition | fail | 0.9 | no | — | step 6: cursor want 31 got 8; state={8 8 8 31 Blahfive lines in this textedit 1 0 z-test-keyboard-harness.md} |
| undo-mid-line-delete | fail | 0.7 | no | — | step 5: cursor want 4 got 1; state={1 1 1 4 adef 1 0 z-test-keyboard-harness.md} |
| redo-cleared-by-new-edit | fail | 0.8 | no | — | step 5: cursor want 3 got 0; textLen want 3 got 0; state={0 0 0 0  1 0 z-test-keyboard-harness.md} |
| undo-after-select-delete | fail | 0.8 | no | — | step 7: cursor want 6 got 0; selStart want 6 got 0; selEnd want 6 got 0; textLen want 6 got 0; state={0 0 0 0  1 0 z-test-keyboard-harness.… |
| gap-up-at-doc-start | pass | 0.6 | no | — | — |
| gap-plain-left-no-cursor-move | fail | 0.8 | no | — | step 4: cursor want 11 got 10; selStart want 11 got 10; selEnd want 11 got 10; state={10 10 10 11 hello world 1 0 z-test-keyboard-harness.m… |
| gap-plain-right-no-cursor-move | fail | 0.5 | no | — | step 4: cursor want 0 got 1; selStart want 0 got 1; selEnd want 0 got 1; state={1 1 1 5 hello 1 0 z-test-keyboard-harness.md} |
| gap-collapse-selection-left | pass | 0.7 | no | — | — |
| gap-collapse-selection-right | pass | 1.0 | no | — | — |
| gap-delete-forward | fail | 0.7 | no | — | step 5: cursor want 2 got 0; state={0 0 0 3 abd 1 0 z-test-keyboard-harness.md} |
| gap-delete-with-selection | pass | 0.8 | no | — | — |
| gap-select-all | pass | 0.4 | no | — | — |
| gap-enter-new-line | pass | 1.2 | no | — | — |
| gap-type-replaces-selection | pass | 0.7 | no | — | — |
| gap-redo-shift-ctrl-z | fail | 0.7 | no | — | step 5: cursor want 3 got 0; textLen want 3 got 0; state={0 0 0 0  1 0 z-test-keyboard-harness.md} |
| gap-undo-chain | fail | 0.8 | no | suspect | step 5: cursor want 2 got 1; textLen want 2 got 1; state={1 1 1 1 a 1 0 z-test-keyboard-harness.md}; may have poisoned next scenario |
| gap-unicode-alt-backspace | prepare fail | 2.7 | yes | gap-undo-chain | textLen want 13 got 11; prepare retries; cascade suspect after gap-undo-chain |
