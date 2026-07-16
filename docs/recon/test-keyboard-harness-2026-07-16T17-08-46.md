# Keyboard harness results

Run: 2026-07-16T17:08:46+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 33 pass, 3 fail, 0 prepare fail; total 64.7s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.4 | no | — | — |
| home-clears-selection | pass | 0.9 | no | — | — |
| shift-right-from-home | pass | 3.4 | no | — | — |
| shift-left-from-end | pass | 3.4 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.4 | no | — | — |
| shift-down-after-arrow-down | pass | 3.3 | no | — | — |
| shift-up-after-arrow-down | pass | 3.3 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | fail | 3.8 | no | — | step 21: cursor want 1362 got 1438; selStart want 1362 got 1438; selEnd want 1362 got 1438; state={1438 1438 1438 1551 Writerdeck harness d… |
| shift-left-repeat-from-end | pass | 3.8 | no | — | — |
| alt-backspace-deletes-word | pass | 1.6 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.6 | no | — | — |
| shift-left-repeat-mid-doc | fail | 2.1 | no | — | step 14: selStart want 1434 got 1435; selLen want 3 got 2; state={1437 1435 1437 1551 Writerdeck harness dummy — ikke i vanlig notatliste  … |
| cm-line-down-basic | pass | 3.3 | no | — | — |
| cm-line-down-last-line | pass | 1.0 | no | — | — |
| combo-alt-left | fail | 1.6 | no | — | step 6: cursor want 1124 got 1138; selStart want 1124 got 1138; selEnd want 1124 got 1138; state={1138 1138 1138 1551 Writerdeck harness du… |
| combo-alt-right | pass | 3.8 | no | — | — |
| combo-ctrl-home | pass | 1.1 | no | — | — |
| combo-ctrl-end | pass | 1.1 | no | — | — |
| bs-plain | pass | 1.0 | no | — | — |
| wrap-down-one-visual-line | pass | 1.3 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.5 | no | — | — |
| undo-redo-len | pass | 0.9 | no | — | — |
| gap-up-at-doc-start | pass | 1.2 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.2 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.3 | no | — | — |
| gap-collapse-selection-left | pass | 0.9 | no | — | — |
| gap-collapse-selection-right | pass | 0.8 | no | — | — |
| gap-delete-forward | pass | 1.7 | no | — | — |
| gap-delete-with-selection | pass | 1.3 | no | — | — |
| gap-select-all | pass | 0.6 | no | — | — |
| gap-enter-new-line | pass | 0.7 | no | — | — |
| gap-type-replaces-selection | pass | 0.9 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 0.9 | no | — | — |
| gap-undo-chain | pass | 0.8 | no | — | — |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
