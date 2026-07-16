# Keyboard harness results

Run: 2026-07-16T23:12:41+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 107 pass, 0 fail, 0 prepare fail; total 193.5s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 1.8 | no | — | — |
| home-clears-selection | pass | 0.8 | no | — | — |
| shift-right-from-home | pass | 3.4 | no | — | — |
| shift-left-from-end | pass | 3.3 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.3 | no | — | — |
| shift-left-after-end-no-stale-anchor | pass | 3.4 | no | — | — |
| shift-down-after-arrow-down | pass | 3.3 | no | — | — |
| shift-up-after-arrow-down | pass | 3.3 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | pass | 3.4 | no | — | — |
| up-one-logical-line | pass | 3.3 | no | — | — |
| shift-down-then-up-shrinks | pass | 3.4 | no | — | — |
| shift-left-repeat-from-end | pass | 3.3 | no | — | — |
| alt-backspace-deletes-word | pass | 1.5 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.6 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.3 | no | — | — |
| cm-line-down-basic | pass | 3.2 | no | — | — |
| cm-line-up-basic | pass | 3.3 | no | — | — |
| cm-line-down-shorter | pass | 0.8 | no | — | — |
| cm-line-up-shorter | pass | 0.7 | no | — | — |
| cm-line-down-last-line | pass | 1.0 | no | — | — |
| cm-line-down-goal-col | pass | 0.8 | no | — | — |
| cm-select-line-down | pass | 3.2 | no | — | — |
| cm-select-up-basic | pass | 3.3 | no | — | — |
| cm-select-line-down-mid | pass | 0.7 | no | — | — |
| cm-select-down-up-doc-end | pass | 1.2 | no | — | — |
| cm-select-up-mid | pass | 0.8 | no | — | — |
| combo-alt-left | pass | 3.4 | no | — | — |
| combo-alt-right | pass | 3.4 | no | — | — |
| combo-alt-up | pass | 1.1 | no | — | — |
| combo-alt-down | pass | 1.0 | no | — | — |
| combo-ctrl-left | pass | 1.2 | no | — | — |
| combo-ctrl-right | pass | 1.2 | no | — | — |
| combo-ctrl-up | pass | 1.4 | no | — | — |
| combo-ctrl-down | pass | 1.2 | no | — | — |
| combo-shift-alt-left | pass | 1.1 | no | — | — |
| combo-shift-alt-left-repeat | pass | 1.1 | no | — | — |
| combo-shift-alt-right | pass | 1.1 | no | — | — |
| combo-shift-alt-right-repeat | pass | 1.0 | no | — | — |
| combo-shift-alt-up | pass | 0.6 | no | — | — |
| combo-shift-alt-down | pass | 0.5 | no | — | — |
| combo-shift-ctrl-left | pass | 1.5 | no | — | — |
| combo-shift-ctrl-left-multiline | pass | 0.7 | no | — | — |
| combo-shift-ctrl-right | pass | 1.5 | no | — | — |
| combo-shift-ctrl-up | pass | 1.5 | no | — | — |
| combo-shift-ctrl-down | pass | 1.4 | no | — | — |
| combo-shift-home-line | pass | 0.7 | no | — | — |
| combo-shift-end-line | pass | 0.6 | no | — | — |
| combo-ctrl-home | pass | 1.0 | no | — | — |
| combo-ctrl-end | pass | 1.1 | no | — | — |
| combo-shift-ctrl-home | pass | 0.6 | no | — | — |
| combo-shift-ctrl-end | pass | 0.7 | no | — | — |
| bs-alt-word-mid | pass | 0.7 | no | — | — |
| bs-ctrl-line-start | pass | 0.7 | no | — | — |
| bs-shift-with-selection | pass | 0.9 | no | — | — |
| bs-plain | pass | 1.0 | no | — | — |
| delete-repeat-forward | pass | 1.0 | no | — | — |
| wrap-down-one-visual-line | pass | 1.1 | no | — | — |
| wrap-down-not-jump-paragraph | pass | 1.0 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.6 | no | — | — |
| wrap-shift-down-one-visual | pass | 1.1 | no | — | — |
| wrap-shift-down-then-up-shrinks | pass | 1.3 | no | — | — |
| wrap-down-last-visual-line | pass | 1.3 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.8 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.8 | no | — | — |
| wrap-down-goal-column | pass | 0.8 | no | — | — |
| wrap-combo-alt-left-word | pass | 1.3 | no | — | — |
| wrap-combo-alt-right-word | pass | 1.3 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 0.9 | no | — | — |
| wrap-shift-left-across-wrap | pass | 1.2 | no | — | — |
| wrap-home-on-visual-line | pass | 0.9 | no | — | — |
| wrap-end-on-visual-line | pass | 0.8 | no | — | — |
| undo-redo-len | pass | 1.1 | no | — | — |
| undo-cursor-reposition | pass | 1.4 | no | — | — |
| undo-mid-line-delete | pass | 1.0 | no | — | — |
| redo-cleared-by-new-edit | pass | 1.0 | no | — | — |
| undo-after-select-delete | pass | 0.7 | no | — | — |
| gap-up-at-doc-start | pass | 1.0 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.3 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.3 | no | — | — |
| gap-plain-left-in-paragraph | pass | 3.3 | no | — | — |
| gap-plain-right-in-paragraph | pass | 3.4 | no | — | — |
| gap-plain-left-at-doc-start | pass | 1.1 | no | — | — |
| gap-plain-right-at-doc-end | pass | 1.1 | no | — | — |
| gap-collapse-selection-left | pass | 0.8 | no | — | — |
| gap-collapse-selection-right | pass | 0.8 | no | — | — |
| gap-delete-forward | pass | 1.0 | no | — | — |
| gap-delete-with-selection | pass | 0.9 | no | — | — |
| gap-select-all | pass | 0.6 | no | — | — |
| gap-enter-new-line | pass | 0.6 | no | — | — |
| gap-type-replaces-selection | pass | 0.9 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 0.9 | no | — | — |
| gap-undo-chain | pass | 0.9 | no | — | — |
| gap-unicode-alt-backspace | pass | 0.6 | no | — | — |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
| gap-alt-bs-with-selection | pass | 0.9 | no | — | — |
| gap-shift-down-mid-wrapping-paras | pass | 3.4 | no | — | — |
| gap-shift-up-mid-wrapping-paras | pass | 3.4 | no | — | — |
| hw-page-right-scrolls-edit | pass | 10.0 | no | — | — |
| hw-page-left-scrolls-edit | pass | 10.7 | no | — | — |
| read-overscroll-clamps | pass | 8.0 | no | — | — |
| touch-down-goal-column | pass | 3.0 | no | — | — |
| touch-up-goal-column | pass | 0.6 | no | — | — |
| touch-down-shorter-line | pass | 0.6 | no | — | — |
| shift-left-then-right-shrinks | pass | 3.5 | no | — | — |
| shift-right-then-left-shrinks | pass | 3.3 | no | — | — |
| shift-up-then-down-shrinks | pass | 3.3 | no | — | — |
