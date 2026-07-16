# Keyboard harness results

Run: 2026-07-17T00:20:41+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 109 pass, 1 fail, 0 prepare fail; total 203.0s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 1.9 | no | — | — |
| home-clears-selection | pass | 1.0 | no | — | — |
| shift-right-from-home | pass | 3.3 | no | — | — |
| shift-left-from-end | pass | 3.5 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.2 | no | — | — |
| shift-left-after-end-no-stale-anchor | pass | 3.3 | no | — | — |
| shift-down-after-arrow-down | pass | 3.4 | no | — | — |
| shift-up-after-arrow-down | pass | 3.4 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | pass | 3.4 | no | — | — |
| up-one-logical-line | pass | 3.3 | no | — | — |
| shift-down-then-up-shrinks | pass | 3.4 | no | — | — |
| shift-left-repeat-from-end | pass | 3.4 | no | — | — |
| alt-backspace-deletes-word | pass | 1.6 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.9 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.3 | no | — | — |
| cm-line-down-basic | pass | 3.4 | no | — | — |
| cm-line-up-basic | pass | 3.4 | no | — | — |
| cm-line-down-shorter | pass | 1.0 | no | — | — |
| cm-line-up-shorter | pass | 0.7 | no | — | — |
| cm-line-down-last-line | pass | 1.0 | no | — | — |
| cm-line-down-goal-col | pass | 0.8 | no | — | — |
| cm-select-line-down | pass | 3.4 | no | — | — |
| cm-select-up-basic | pass | 3.4 | no | — | — |
| cm-select-line-down-mid | pass | 0.5 | no | — | — |
| cm-select-down-up-doc-end | pass | 1.1 | no | — | — |
| cm-select-up-mid | pass | 0.6 | no | — | — |
| combo-alt-left | fail | 1.5 | no | — | step 9: cursor want 1151 got 1156; selStart want 1151 got 1156; selEnd want 1151 got 1156; state={1156 1156 1156 1551 Writerdeck harness du… |
| combo-alt-right | pass | 3.6 | no | — | — |
| combo-alt-up | pass | 1.1 | no | — | — |
| combo-alt-down | pass | 1.0 | no | — | — |
| combo-ctrl-left | pass | 1.3 | no | — | — |
| combo-ctrl-right | pass | 1.1 | no | — | — |
| combo-ctrl-up | pass | 1.2 | no | — | — |
| combo-ctrl-down | pass | 1.2 | no | — | — |
| combo-shift-alt-left | pass | 1.0 | no | — | — |
| combo-shift-alt-left-repeat | pass | 1.1 | no | — | — |
| combo-shift-alt-right | pass | 1.1 | no | — | — |
| combo-shift-alt-right-repeat | pass | 1.0 | no | — | — |
| combo-shift-alt-up | pass | 0.6 | no | — | — |
| combo-shift-alt-down | pass | 0.6 | no | — | — |
| combo-shift-ctrl-left | pass | 1.4 | no | — | — |
| combo-shift-ctrl-left-multiline | pass | 0.7 | no | — | — |
| combo-shift-ctrl-right | pass | 1.6 | no | — | — |
| combo-shift-ctrl-up | pass | 1.6 | no | — | — |
| combo-shift-ctrl-down | pass | 1.4 | no | — | — |
| combo-shift-home-line | pass | 0.6 | no | — | — |
| combo-shift-end-line | pass | 0.6 | no | — | — |
| combo-ctrl-home | pass | 1.1 | no | — | — |
| combo-ctrl-end | pass | 1.1 | no | — | — |
| combo-shift-ctrl-home | pass | 0.6 | no | — | — |
| combo-shift-ctrl-end | pass | 0.7 | no | — | — |
| bs-alt-word-mid | pass | 0.7 | no | — | — |
| bs-ctrl-line-start | pass | 0.7 | no | — | — |
| bs-shift-with-selection | pass | 0.8 | no | — | — |
| bs-plain | pass | 1.1 | no | — | — |
| delete-repeat-forward | pass | 1.1 | no | — | — |
| wrap-down-one-visual-line | pass | 1.3 | no | — | — |
| wrap-down-not-jump-paragraph | pass | 0.8 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.5 | no | — | — |
| wrap-shift-down-one-visual | pass | 1.3 | no | — | — |
| wrap-shift-down-then-up-shrinks | pass | 3.4 | no | — | — |
| wrap-down-last-visual-line | pass | 1.3 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.8 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.8 | no | — | — |
| wrap-down-goal-column | pass | 1.0 | no | — | — |
| wrap-combo-alt-left-word | pass | 1.4 | no | — | — |
| wrap-combo-alt-right-word | pass | 1.3 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 0.8 | no | — | — |
| wrap-shift-left-across-wrap | pass | 1.4 | no | — | — |
| wrap-home-on-visual-line | pass | 0.8 | no | — | — |
| wrap-end-on-visual-line | pass | 1.2 | no | — | — |
| undo-redo-len | pass | 1.2 | no | — | — |
| undo-cursor-reposition | pass | 2.5 | no | — | — |
| undo-mid-line-delete | pass | 0.9 | no | — | — |
| redo-cleared-by-new-edit | pass | 1.2 | no | — | — |
| undo-after-select-delete | pass | 0.8 | no | — | — |
| gap-up-at-doc-start | pass | 1.0 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.4 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.8 | no | — | — |
| gap-plain-left-in-paragraph | pass | 3.6 | no | — | — |
| gap-plain-right-in-paragraph | pass | 3.2 | no | — | — |
| gap-plain-left-at-doc-start | pass | 1.0 | no | — | — |
| gap-plain-right-at-doc-end | pass | 1.1 | no | — | — |
| gap-collapse-selection-left | pass | 1.0 | no | — | — |
| gap-collapse-selection-right | pass | 0.8 | no | — | — |
| gap-delete-forward | pass | 1.1 | no | — | — |
| gap-delete-with-selection | pass | 0.9 | no | — | — |
| gap-select-all | pass | 0.7 | no | — | — |
| gap-enter-new-line | pass | 0.7 | no | — | — |
| gap-type-replaces-selection | pass | 0.8 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 0.9 | no | — | — |
| gap-undo-chain | pass | 0.8 | no | — | — |
| gap-unicode-alt-backspace | pass | 0.7 | no | — | — |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
| gap-alt-bs-with-selection | pass | 1.0 | no | — | — |
| gap-shift-down-mid-wrapping-paras | pass | 3.4 | no | — | — |
| gap-shift-up-mid-wrapping-paras | pass | 3.4 | no | — | — |
| gap-shift-down-across-para-break | pass | 1.1 | no | — | — |
| gap-shift-up-across-para-break | pass | 1.0 | no | — | — |
| gap-shift-down-mid-short-lines | pass | 3.5 | no | — | — |
| hw-page-right-scrolls-edit | pass | 9.9 | no | — | — |
| hw-page-left-scrolls-edit | pass | 11.5 | no | — | — |
| read-overscroll-clamps | pass | 7.6 | no | — | — |
| touch-down-goal-column | pass | 1.3 | no | — | — |
| touch-up-goal-column | pass | 0.5 | no | — | — |
| touch-down-shorter-line | pass | 0.6 | no | — | — |
| shift-left-then-right-shrinks | pass | 3.4 | no | — | — |
| shift-right-then-left-shrinks | pass | 3.4 | no | — | — |
| shift-up-then-down-shrinks | pass | 3.6 | no | — | — |
