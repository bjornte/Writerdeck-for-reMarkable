# Keyboard harness results

Run: 2026-07-17T01:52:45+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 35 pass, 3 fail, 0 prepare fail; total 74.3s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 1.9 | no | — | — |
| home-clears-selection | pass | 0.8 | no | — | — |
| shift-right-from-home | pass | 3.5 | no | — | — |
| shift-left-from-end | pass | 3.3 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.4 | no | — | — |
| shift-down-after-arrow-down | pass | 3.5 | no | — | — |
| shift-up-after-arrow-down | pass | 3.4 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.5 | no | — | — |
| down-one-logical-line | fail | 1.1 | no | — | step 6: cursor want 1457 got 1438; selStart want 1457 got 1438; selEnd want 1457 got 1438; state={1438 1438 1438 1551 Writerdeck harness du… |
| shift-left-repeat-from-end | pass | 3.2 | no | — | — |
| alt-backspace-deletes-word | pass | 1.6 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.6 | no | — | — |
| shift-left-repeat-mid-doc | fail | 3.1 | no | — | step 19: selStart want 1430 got 1431; selLen want 7 got 6; state={1437 1431 1437 1551 Writerdeck harness dummy — ikke i vanlig notatliste  … |
| cm-line-down-basic | pass | 3.4 | no | — | — |
| cm-line-down-last-line | pass | 1.6 | no | — | — |
| combo-alt-left | pass | 3.5 | no | — | — |
| combo-alt-right | pass | 3.5 | no | — | — |
| combo-ctrl-home | pass | 1.3 | no | — | — |
| combo-ctrl-end | pass | 1.1 | no | — | — |
| bs-plain | pass | 1.0 | no | — | — |
| wrap-down-one-visual-line | pass | 1.3 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.6 | no | — | — |
| undo-redo-len | pass | 1.0 | no | — | — |
| gap-up-at-doc-start | pass | 1.1 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.4 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.3 | no | — | — |
| gap-collapse-selection-left | pass | 0.8 | no | — | — |
| gap-collapse-selection-right | pass | 0.8 | no | — | — |
| gap-delete-forward | pass | 1.0 | no | — | — |
| gap-delete-with-selection | pass | 1.2 | no | — | — |
| gap-select-all | pass | 0.7 | no | — | — |
| gap-enter-new-line | pass | 0.7 | no | — | — |
| gap-type-replaces-selection | pass | 1.3 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 1.0 | no | — | — |
| gap-undo-chain | pass | 1.0 | no | — | — |
| gap-empty-doc-backspace | pass | 0.6 | no | — | — |
| gap-shift-down-mid-wrapping-paras | fail | 3.5 | no | — | step 21: cursorMax want <= 268 got 290; selLenMax want <= 90 got 112; state={290 178 290 1551 Writerdeck harness dummy — ikke i vanlig nota… |
| gap-shift-up-mid-wrapping-paras | pass | 3.5 | no | — | — |
