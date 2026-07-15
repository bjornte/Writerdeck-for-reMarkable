# Keyboard harness results

Run: 2026-07-16T00:37:28+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 89 pass, 4 fail, 1 prepare fail; total 67.2s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 2.6 | no | — | — |
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
| alt-backspace-deletes-word | pass | 0.6 | no | — | — |
| ctrl-backspace-deletes-line | pass | 0.8 | no | — | — |
| shift-left-repeat-mid-doc | pass | 0.7 | no | — | — |
| cm-line-down-basic | pass | 0.5 | no | — | — |
| cm-line-down-shorter | pass | 0.8 | no | — | — |
| cm-line-down-last-line | pass | 0.6 | no | — | — |
| cm-line-down-goal-col | fail | 0.7 | no | — | step 4: cursor want 6 got 7; selStart want 6 got 7; selEnd want 6 got 7; state={7 7 7 11 one t three 1 0 z-test-keyboard-harness.md 0} |
| cm-select-line-down | pass | 0.5 | no | — | — |
| cm-select-line-down-mid | pass | 0.8 | no | — | — |
| cm-select-down-up-doc-end | fail | 0.8 | no | — | step 5: selStart want 8 got 7; selLen want 4 got 5; state={12 7 12 13 one two three 1 0 z-test-keyboard-harness.md 0} |
| cm-select-up-basic | pass | 0.7 | no | — | — |
| cm-select-up-mid | pass | 0.8 | no | — | — |
| combo-alt-left | pass | 0.5 | no | — | — |
| combo-alt-right | pass | 0.5 | no | — | — |
| combo-alt-up | pass | 0.6 | no | — | — |
| combo-alt-down | pass | 0.5 | no | — | — |
| combo-ctrl-left | pass | 0.5 | no | — | — |
| combo-ctrl-right | pass | 0.5 | no | — | — |
| combo-ctrl-up | pass | 0.5 | no | — | — |
| combo-ctrl-down | pass | 0.5 | no | — | — |
| combo-shift-alt-left | pass | 0.5 | no | — | — |
| combo-shift-alt-left-repeat | pass | 0.5 | no | — | — |
| combo-shift-alt-right | pass | 0.6 | no | — | — |
| combo-shift-alt-right-repeat | pass | 0.5 | no | — | — |
| combo-shift-alt-up | pass | 0.6 | no | — | — |
| combo-shift-alt-down | pass | 0.5 | no | — | — |
| combo-shift-ctrl-left | pass | 0.5 | no | — | — |
| combo-shift-ctrl-left-multiline | pass | 0.7 | no | — | — |
| combo-shift-ctrl-right | pass | 0.5 | no | — | — |
| combo-shift-ctrl-up | pass | 0.7 | no | — | — |
| combo-shift-ctrl-down | pass | 0.5 | no | — | — |
| combo-shift-home-line | pass | 0.6 | no | — | — |
| combo-shift-end-line | pass | 0.7 | no | — | — |
| combo-ctrl-home | pass | 0.5 | no | — | — |
| combo-ctrl-end | pass | 0.5 | no | — | — |
| combo-shift-ctrl-home | pass | 0.5 | no | — | — |
| combo-shift-ctrl-end | pass | 0.6 | no | — | — |
| bs-alt-word-mid | pass | 0.8 | no | — | — |
| bs-ctrl-line-start | pass | 0.5 | no | — | — |
| bs-shift-with-selection | pass | 0.6 | no | — | — |
| bs-plain | pass | 0.7 | no | — | — |
| wrap-down-one-visual-line | pass | 0.6 | no | — | — |
| wrap-down-not-jump-paragraph | pass | 0.7 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 0.9 | no | — | — |
| wrap-shift-down-one-visual | pass | 0.7 | no | — | — |
| wrap-shift-down-then-up-shrinks | pass | 1.0 | no | — | — |
| wrap-down-last-visual-line | pass | 0.8 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.7 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.6 | no | — | — |
| wrap-down-goal-column | pass | 0.8 | no | — | — |
| wrap-combo-alt-left-word | pass | 0.9 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 0.8 | no | — | — |
| wrap-shift-left-across-wrap | pass | 1.0 | no | — | — |
| wrap-home-on-visual-line | pass | 0.8 | no | — | — |
| wrap-end-on-visual-line | pass | 0.8 | no | — | — |
| undo-redo-len | pass | 1.0 | no | — | — |
| undo-cursor-reposition | pass | 1.3 | no | — | — |
| undo-mid-line-delete | pass | 1.0 | no | — | — |
| redo-cleared-by-new-edit | pass | 1.2 | no | — | — |
| undo-after-select-delete | pass | 0.8 | no | — | — |
| gap-up-at-doc-start | pass | 0.5 | no | — | — |
| gap-plain-left-moves-caret | pass | 0.5 | no | — | — |
| gap-plain-right-moves-caret | pass | 0.5 | no | — | — |
| gap-plain-left-at-doc-start | pass | 0.5 | no | — | — |
| gap-plain-right-at-doc-end | pass | 0.5 | no | — | — |
| gap-collapse-selection-left | pass | 0.8 | no | — | — |
| gap-collapse-selection-right | pass | 0.7 | no | — | — |
| gap-delete-forward | pass | 0.7 | no | — | — |
| gap-delete-with-selection | pass | 0.7 | no | — | — |
| gap-select-all | pass | 0.4 | no | — | — |
| gap-enter-new-line | pass | 0.5 | no | — | — |
| gap-type-replaces-selection | pass | 0.6 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 1.0 | no | — | — |
| gap-undo-chain | pass | 0.9 | no | — | — |
| gap-unicode-alt-backspace | prepare fail | 2.6 | yes | — | textLen want 13 got 11; prepare retries |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
| gap-alt-bs-with-selection | fail | 0.7 | no | — | step 5: cursor want 5 got 6; selStart want 5 got 6; selEnd want 5 got 6; textLen want 5 got 6; text want "hello" got "hello "; state={6 6 6… |
| hw-page-right-scrolls-edit | pass | 1.1 | no | — | — |
| hw-page-left-scrolls-edit | pass | 1.3 | no | — | — |
| touch-down-goal-column | pass | 0.7 | no | — | — |
| touch-up-goal-column | pass | 0.6 | no | — | — |
| touch-down-shorter-line | pass | 0.6 | no | — | — |
| shift-left-then-right-shrinks | pass | 0.7 | no | — | — |
