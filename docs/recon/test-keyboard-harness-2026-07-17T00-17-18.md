# Keyboard harness results

Run: 2026-07-17T00:17:18+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 34 pass, 4 fail, 0 prepare fail; total 70.2s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.3 | no | — | — |
| home-clears-selection | pass | 0.8 | no | — | — |
| shift-right-from-home | pass | 3.2 | no | — | — |
| shift-left-from-end | pass | 3.2 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.2 | no | — | — |
| shift-down-after-arrow-down | pass | 3.3 | no | — | — |
| shift-up-after-arrow-down | fail | 1.0 | no | — | step 3: selStart want 1476 got 1495; selLen want 19 got 0; state={1495 1495 1495 1551 Writerdeck harness dummy — ikke i vanlig notatliste  … |
| ctrl-shift-left-select-line | fail | 1.1 | no | — | step 3: selStart want 1272 got 1323; selLen want 51 got 0; state={1323 1323 1323 1551 Writerdeck harness dummy — ikke i vanlig notatliste  … |
| down-one-logical-line | fail | 2.1 | no | — | step 6: cursor want 1457 got 1362; selStart want 1457 got 1362; selEnd want 1457 got 1362; state={1362 1362 1362 1551 Writerdeck harness du… |
| shift-left-repeat-from-end | fail | 3.1 | no | — | step 14: selStart want 1320 got 1323; selLen want 3 got 0; state={1323 1323 1323 1551 Writerdeck harness dummy — ikke i vanlig notatliste  … |
| alt-backspace-deletes-word | pass | 2.7 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.8 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.3 | no | — | — |
| cm-line-down-basic | pass | 3.3 | no | — | — |
| cm-line-down-last-line | pass | 1.1 | no | — | — |
| combo-alt-left | pass | 3.5 | no | — | — |
| combo-alt-right | pass | 3.6 | no | — | — |
| combo-ctrl-home | pass | 1.1 | no | — | — |
| combo-ctrl-end | pass | 1.0 | no | — | — |
| bs-plain | pass | 1.0 | no | — | — |
| wrap-down-one-visual-line | pass | 1.3 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.6 | no | — | — |
| undo-redo-len | pass | 1.1 | no | — | — |
| gap-up-at-doc-start | pass | 1.0 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.4 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.2 | no | — | — |
| gap-collapse-selection-left | pass | 0.8 | no | — | — |
| gap-collapse-selection-right | pass | 0.8 | no | — | — |
| gap-delete-forward | pass | 1.0 | no | — | — |
| gap-delete-with-selection | pass | 0.9 | no | — | — |
| gap-select-all | pass | 1.1 | no | — | — |
| gap-enter-new-line | pass | 0.6 | no | — | — |
| gap-type-replaces-selection | pass | 0.8 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 1.0 | no | — | — |
| gap-undo-chain | pass | 0.9 | no | — | — |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
| gap-shift-down-mid-wrapping-paras | pass | 3.2 | no | — | — |
| gap-shift-up-mid-wrapping-paras | pass | 3.4 | no | — | — |
