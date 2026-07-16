# Keyboard harness results

Run: 2026-07-16T16:59:09+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 32 pass, 4 fail, 0 prepare fail; total 61.3s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.6 | no | — | — |
| home-clears-selection | pass | 0.9 | no | — | — |
| shift-right-from-home | pass | 3.2 | no | — | — |
| shift-left-from-end | pass | 3.3 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.5 | no | — | — |
| shift-down-after-arrow-down | fail | 1.6 | no | — | step 11: cursor want 1362 got 1381; selEnd want 1362 got 1381; selLen want 0 got 19; state={1381 1362 1381 1551 Writerdeck harness dummy — … |
| shift-up-after-arrow-down | pass | 3.3 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | pass | 3.4 | no | — | — |
| shift-left-repeat-from-end | pass | 3.3 | no | — | — |
| alt-backspace-deletes-word | pass | 1.6 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.7 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.5 | no | — | — |
| cm-line-down-basic | pass | 3.2 | no | — | — |
| cm-line-down-last-line | pass | 1.1 | no | — | — |
| combo-alt-left | pass | 3.4 | no | — | — |
| combo-alt-right | fail | 1.1 | no | — | step 6: cursor want 1168 got 1162; selStart want 1168 got 1162; selEnd want 1168 got 1162; state={1162 1162 1162 1551 Writerdeck harness du… |
| combo-ctrl-home | pass | 1.0 | no | — | — |
| combo-ctrl-end | pass | 1.0 | no | — | — |
| bs-plain | pass | 1.0 | no | — | — |
| wrap-down-one-visual-line | pass | 1.1 | no | — | — |
| wrap-up-from-visual-line-2 | fail | 1.0 | no | — | step 3: cursor want 140 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word wo… |
| undo-redo-len | pass | 0.9 | no | — | — |
| gap-up-at-doc-start | pass | 1.2 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.3 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.6 | no | — | — |
| gap-collapse-selection-left | pass | 0.8 | no | — | — |
| gap-collapse-selection-right | pass | 0.8 | no | — | — |
| gap-delete-forward | pass | 1.2 | no | — | — |
| gap-delete-with-selection | pass | 1.1 | no | — | — |
| gap-select-all | pass | 0.7 | no | — | — |
| gap-enter-new-line | pass | 0.6 | no | — | — |
| gap-type-replaces-selection | pass | 0.9 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 1.0 | no | — | — |
| gap-undo-chain | fail | 0.5 | no | — | step 3: cursor want 5 got 6; textLen want 5 got 6; state={6 6 6 6 abc æø 1 0 z-test-keyboard-harness.md 0} |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
