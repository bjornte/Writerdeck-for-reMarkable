# Keyboard harness results

Run: 2026-07-16T17:07:24+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 34 pass, 2 fail, 0 prepare fail; total 68.5s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.4 | no | — | — |
| home-clears-selection | pass | 0.9 | no | — | — |
| shift-right-from-home | pass | 3.8 | no | — | — |
| shift-left-from-end | pass | 3.4 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.5 | no | — | — |
| shift-down-after-arrow-down | pass | 3.3 | no | — | — |
| shift-up-after-arrow-down | pass | 3.3 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | pass | 3.5 | no | — | — |
| shift-left-repeat-from-end | fail | 2.0 | no | — | step 14: selStart want 1320 got 1321; selLen want 3 got 2; state={1323 1321 1323 1551 Writerdeck harness dummy — ikke i vanlig notatliste  … |
| alt-backspace-deletes-word | pass | 1.5 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.6 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.4 | no | — | — |
| cm-line-down-basic | fail | 3.5 | no | — | step 21: cursor want 1362 got 1381; selStart want 1362 got 1381; selEnd want 1362 got 1381; state={1381 1381 1381 1551 Writerdeck harness d… |
| cm-line-down-last-line | pass | 1.0 | no | — | — |
| combo-alt-left | pass | 3.5 | no | — | — |
| combo-alt-right | pass | 3.5 | no | — | — |
| combo-ctrl-home | pass | 1.1 | no | — | — |
| combo-ctrl-end | pass | 1.8 | no | — | — |
| bs-plain | pass | 1.6 | no | — | — |
| wrap-down-one-visual-line | pass | 1.3 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.6 | no | — | — |
| undo-redo-len | pass | 1.0 | no | — | — |
| gap-up-at-doc-start | pass | 1.0 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.3 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.4 | no | — | — |
| gap-collapse-selection-left | pass | 1.0 | no | — | — |
| gap-collapse-selection-right | pass | 0.9 | no | — | — |
| gap-delete-forward | pass | 1.1 | no | — | — |
| gap-delete-with-selection | pass | 1.0 | no | — | — |
| gap-select-all | pass | 1.2 | no | — | — |
| gap-enter-new-line | pass | 1.1 | no | — | — |
| gap-type-replaces-selection | pass | 1.0 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 1.0 | no | — | — |
| gap-undo-chain | pass | 1.0 | no | — | — |
| gap-empty-doc-backspace | pass | 0.6 | no | — | — |
