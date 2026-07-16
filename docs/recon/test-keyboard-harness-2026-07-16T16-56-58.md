# Keyboard harness results

Run: 2026-07-16T16:56:58+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 35 pass, 1 fail, 0 prepare fail; total 68.8s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.4 | no | — | — |
| home-clears-selection | pass | 0.8 | no | — | — |
| shift-right-from-home | pass | 3.9 | no | — | — |
| shift-left-from-end | pass | 3.4 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 4.0 | no | — | — |
| shift-down-after-arrow-down | pass | 3.9 | no | — | — |
| shift-up-after-arrow-down | pass | 3.7 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | pass | 3.4 | no | — | — |
| shift-left-repeat-from-end | pass | 3.3 | no | — | — |
| alt-backspace-deletes-word | pass | 1.6 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.5 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.3 | no | — | — |
| cm-line-down-basic | pass | 3.5 | no | — | — |
| cm-line-down-last-line | pass | 1.1 | no | — | — |
| combo-alt-left | pass | 3.4 | no | — | — |
| combo-alt-right | pass | 3.8 | no | — | — |
| combo-ctrl-home | pass | 1.0 | no | — | — |
| combo-ctrl-end | pass | 1.6 | no | — | — |
| bs-plain | pass | 1.1 | no | — | — |
| wrap-down-one-visual-line | pass | 1.6 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.6 | no | — | — |
| undo-redo-len | pass | 0.9 | no | — | — |
| gap-up-at-doc-start | pass | 1.2 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.7 | no | — | — |
| gap-plain-right-moves-caret | fail | 2.9 | no | — | step 19: cursor want 1304 got 1303; selStart want 1304 got 1303; selEnd want 1304 got 1303; state={1303 1303 1303 1551 Writerdeck harness d… |
| gap-collapse-selection-left | pass | 0.8 | no | — | — |
| gap-collapse-selection-right | pass | 0.8 | no | — | — |
| gap-delete-forward | pass | 1.0 | no | — | — |
| gap-delete-with-selection | pass | 0.9 | no | — | — |
| gap-select-all | pass | 0.6 | no | — | — |
| gap-enter-new-line | pass | 0.7 | no | — | — |
| gap-type-replaces-selection | pass | 0.8 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 0.9 | no | — | — |
| gap-undo-chain | pass | 0.9 | no | — | — |
| gap-empty-doc-backspace | pass | 0.5 | no | — | — |
