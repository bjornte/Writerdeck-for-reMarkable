# Keyboard harness results

Run: 2026-07-17T14:46:12+02:00

Target: `192.168.86.30:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 71 pass, 39 fail, 0 prepare fail; total 177.3s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 3.5 | no | — | — |
| home-clears-selection | pass | 1.0 | no | — | — |
| shift-right-from-home | pass | 3.5 | no | — | — |
| shift-left-from-end | pass | 3.4 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.5 | no | — | — |
| shift-left-after-end-no-stale-anchor | pass | 3.3 | no | — | — |
| shift-down-after-arrow-down | fail | 0.8 | no | — | step 3: cursor want 1381 got 1551; selEnd want 1381 got 1551; selLen want 19 got 189; state={1551 1362 1551 1551 Writerdeck harness dummy —… |
| shift-up-after-arrow-down | fail | 1.0 | no | — | step 3: selStart want 1476 got 0; selLen want 19 got 1495; state={1495 0 1495 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Før… |
| ctrl-shift-left-select-line | pass | 0.7 | no | — | — |
| down-one-logical-line | fail | 0.6 | no | — | step 3: cursor want 1381 got 1551; selStart want 1381 got 1551; selEnd want 1381 got 1551; state={1551 1551 1551 1551 Writerdeck harness du… |
| up-one-logical-line | fail | 0.9 | no | — | step 3: cursor want 1476 got 0; selStart want 1476 got 0; selEnd want 1476 got 0; state={0 0 0 1551 Writerdeck harness dummy — ikke i vanli… |
| shift-down-then-up-shrinks | fail | 0.7 | no | — | step 3: cursor want 1381 got 1551; selEnd want 1381 got 1551; selLen want 19 got 189; state={1551 1362 1551 1551 Writerdeck harness dummy —… |
| shift-left-repeat-from-end | pass | 3.5 | no | — | — |
| alt-backspace-deletes-word | pass | 1.7 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.8 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.5 | no | — | — |
| cm-line-down-basic | fail | 0.7 | no | — | step 3: cursor want 1381 got 1551; selStart want 1381 got 1551; selEnd want 1381 got 1551; state={1551 1551 1551 1551 Writerdeck harness du… |
| cm-line-up-basic | fail | 0.9 | no | — | step 3: cursor want 1476 got 0; selStart want 1476 got 0; selEnd want 1476 got 0; state={0 0 0 1551 Writerdeck harness dummy — ikke i vanli… |
| cm-line-down-shorter | fail | 0.8 | no | — | step 5: cursor want 4 got 5; selStart want 4 got 5; selEnd want 4 got 5; state={5 5 5 5 tre i 1 0 z-test-keyboard-harness.md 0} |
| cm-line-up-shorter | fail | 0.4 | no | — | step 2: cursor want 2 got 5; state={5 5 5 5 i tre 1 0 z-test-keyboard-harness.md 0} |
| cm-line-down-last-line | pass | 1.2 | no | — | — |
| cm-line-down-goal-col | fail | 0.8 | no | — | step 4: cursorMax want <= 8 got 11; state={11 11 11 11 tre i femte 1 0 z-test-keyboard-harness.md 0} |
| cm-select-line-down | fail | 0.8 | no | — | step 3: cursor want 1381 got 1551; selEnd want 1381 got 1551; selLen want 19 got 189; state={1551 1362 1551 1551 Writerdeck harness dummy —… |
| cm-select-up-basic | fail | 1.0 | no | — | step 3: selStart want 1476 got 0; selLen want 19 got 1495; state={1495 0 1495 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Før… |
| cm-select-line-down-mid | fail | 0.9 | no | — | step 3: cursorMax want <= 268 got 1551; selLenMax want <= 90 got 1373; state={1551 178 1551 1551 Writerdeck harness dummy — ikke i vanlig n… |
| cm-select-down-up-doc-end | pass | 1.2 | no | — | — |
| cm-select-up-mid | fail | 0.8 | no | — | step 3: selLenMax want <= 90 got 471; state={471 0 471 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Første avsnitt — Naïve caf… |
| combo-alt-left | pass | 3.5 | no | — | — |
| combo-alt-right | pass | 3.6 | no | — | — |
| combo-alt-up | pass | 1.1 | no | — | — |
| combo-alt-down | pass | 1.0 | no | — | — |
| combo-ctrl-left | pass | 1.3 | no | — | — |
| combo-ctrl-right | pass | 1.1 | no | — | — |
| combo-ctrl-up | pass | 1.3 | no | — | — |
| combo-ctrl-down | pass | 1.1 | no | — | — |
| combo-shift-alt-left | pass | 1.1 | no | — | — |
| combo-shift-alt-left-repeat | pass | 1.2 | no | — | — |
| combo-shift-alt-right | pass | 1.2 | no | — | — |
| combo-shift-alt-right-repeat | pass | 1.1 | no | — | — |
| combo-shift-alt-up | pass | 0.6 | no | — | — |
| combo-shift-alt-down | pass | 0.6 | no | — | — |
| combo-shift-ctrl-left | pass | 1.4 | no | — | — |
| combo-shift-ctrl-left-multiline | fail | 0.9 | no | — | step 4: cursor want 5 got 9; selStart want 3 got 0; selEnd want 5 got 9; selLen want 2 got 9; state={9 0 9 9 en to tre 1 0 z-test-keyboard-… |
| combo-shift-ctrl-right | pass | 1.4 | no | — | — |
| combo-shift-ctrl-up | pass | 1.7 | no | — | — |
| combo-shift-ctrl-down | pass | 1.4 | no | — | — |
| combo-shift-home-line | fail | 0.8 | no | — | step 4: selStart want 4 got 0; selLen want 3 got 7; state={7 0 7 7 ost ost 1 0 z-test-keyboard-harness.md 0} |
| combo-shift-end-line | fail | 0.7 | no | — | step 4: cursor want 7 got 3; selStart want 4 got 0; selEnd want 7 got 3; state={3 0 3 7 ost ost 1 0 z-test-keyboard-harness.md 0} |
| combo-ctrl-home | pass | 1.3 | no | — | — |
| combo-ctrl-end | pass | 1.2 | no | — | — |
| combo-shift-ctrl-home | fail | 0.5 | no | — | step 2: cursor want 4 got 7; state={7 7 7 7 ost ost 1 0 z-test-keyboard-harness.md 0} |
| combo-shift-ctrl-end | fail | 0.5 | no | — | step 2: cursor want 4 got 7; state={7 7 7 7 ost ost 1 0 z-test-keyboard-harness.md 0} |
| bs-alt-word-mid | pass | 0.9 | no | — | — |
| bs-ctrl-line-start | pass | 0.9 | no | — | — |
| bs-shift-with-selection | pass | 0.9 | no | — | — |
| bs-plain | pass | 1.1 | no | — | — |
| delete-repeat-forward | pass | 1.3 | no | — | — |
| wrap-down-one-visual-line | fail | 0.8 | no | — | step 3: cursor want 10 got 199; state={199 199 199 199 word word word word word word word word word word word word word word word word word… |
| wrap-down-not-jump-paragraph | fail | 0.8 | no | — | step 3: cursor want 10 got 199; state={199 199 199 199 word word word word word word word word word word word word word word word word word… |
| wrap-up-from-visual-line-2 | fail | 0.9 | no | — | step 3: cursor want 70 got 199; state={199 199 199 199 word word word word word word word word word word word word word word word word word… |
| wrap-shift-down-one-visual | fail | 0.8 | no | — | step 3: cursor want 10 got 199; selEnd want 10 got 199; state={199 0 199 199 word word word word word word word word word word word word wo… |
| wrap-shift-down-then-up-shrinks | fail | 0.7 | no | — | step 3: cursorMax want <= 45 got 199; selLenMax want <= 20 got 174; state={199 25 199 199 word word word word word word word word word word… |
| wrap-down-last-visual-line | pass | 1.3 | no | — | — |
| wrap-shift-down-last-to-eof | fail | 0.9 | no | — | step 3: selLenMin want >= 1 got 0; state={199 199 199 199 word word word word word word word word word word word word word word word word w… |
| wrap-mixed-newline-and-wrap | pass | 0.7 | no | — | — |
| wrap-down-goal-column | fail | 1.0 | no | — | step 5: cursor want 14 got 177; state={177 177 177 177 abword word word word word word word word word word word word word word word word wo… |
| wrap-combo-alt-left-word | pass | 1.4 | no | — | — |
| wrap-combo-alt-right-word | pass | 1.3 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 0.9 | no | — | — |
| wrap-shift-left-across-wrap | fail | 0.7 | no | — | step 3: cursor want 10 got 199; state={199 199 199 199 word word word word word word word word word word word word word word word word word… |
| wrap-home-on-visual-line | fail | 0.8 | no | — | step 3: cursor want 10 got 199; state={199 199 199 199 word word word word word word word word word word word word word word word word word… |
| wrap-end-on-visual-line | fail | 0.8 | no | — | step 3: cursor want 10 got 199; state={199 199 199 199 word word word word word word word word word word word word word word word word word… |
| undo-redo-len | pass | 1.0 | no | — | — |
| undo-cursor-reposition | pass | 1.5 | no | — | — |
| undo-mid-line-delete | fail | 0.7 | no | — | step 3: cursor want 4 got 9; state={9 9 9 9 ost rømme 1 0 z-test-keyboard-harness.md 0} |
| redo-cleared-by-new-edit | pass | 1.2 | no | — | — |
| undo-after-select-delete | pass | 0.8 | no | — | — |
| gap-up-at-doc-start | pass | 1.3 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.4 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.4 | no | — | — |
| gap-plain-left-in-paragraph | pass | 3.4 | no | — | — |
| gap-plain-right-in-paragraph | pass | 3.4 | no | — | — |
| gap-plain-left-at-doc-start | pass | 1.1 | no | — | — |
| gap-plain-right-at-doc-end | pass | 1.1 | no | — | — |
| gap-collapse-selection-left | pass | 0.8 | no | — | — |
| gap-collapse-selection-right | pass | 0.9 | no | — | — |
| gap-delete-forward | pass | 1.3 | no | — | — |
| gap-delete-with-selection | pass | 0.9 | no | — | — |
| gap-select-all | pass | 0.7 | no | — | — |
| gap-enter-new-line | pass | 0.6 | no | — | — |
| gap-type-replaces-selection | pass | 1.0 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 1.0 | no | — | — |
| gap-undo-chain | pass | 0.9 | no | — | — |
| gap-unicode-alt-backspace | pass | 0.6 | no | — | — |
| gap-empty-doc-backspace | pass | 0.5 | no | — | — |
| gap-alt-bs-with-selection | pass | 0.8 | no | — | — |
| gap-shift-down-mid-wrapping-paras | fail | 1.1 | no | — | step 3: cursorMax want <= 268 got 1551; selLenMax want <= 90 got 1373; state={1551 178 1551 1551 Writerdeck harness dummy — ikke i vanlig n… |
| gap-shift-up-mid-wrapping-paras | fail | 0.8 | no | — | step 3: selLenMax want <= 90 got 471; state={471 0 471 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Første avsnitt — Naïve caf… |
| gap-shift-down-across-para-break | fail | 2.2 | no | — | step 5: selLenMax want <= 90 got 377; state={377 0 377 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Første avsnitt — Naïve caf… |
| gap-shift-up-across-para-break | fail | 2.2 | no | — | step 5: cursorMax want <= 531 got 1551; selLenMax want <= 90 got 1110; state={1551 441 1551 1551 Writerdeck harness dummy — ikke i vanlig n… |
| gap-shift-down-mid-short-lines | fail | 0.7 | no | — | step 3: cursor want 1428 got 1551; selEnd want 1428 got 1551; selLen want 19 got 142; state={1551 1409 1551 1551 Writerdeck harness dummy —… |
| hw-page-right-scrolls-edit | pass | 11.6 | no | — | — |
| hw-page-left-scrolls-edit | pass | 11.0 | no | — | — |
| read-overscroll-clamps | pass | 9.3 | no | — | — |
| touch-down-goal-column | fail | 2.2 | no | — | step 4: cursor want 4 got 9; selStart want 4 got 9; selEnd want 4 got 9; state={9 9 9 9 en to tre 1 0 z-test-keyboard-harness.md 0} |
| touch-up-goal-column | fail | 1.2 | no | — | step 4: cursor want 1 got 0; selStart want 1 got 0; selEnd want 1 got 0; state={0 0 0 9 en to tre 1 0 z-test-keyboard-harness.md 0} |
| touch-down-shorter-line | fail | 1.0 | no | — | step 4: cursor want 4 got 5; selStart want 4 got 5; selEnd want 4 got 5; state={5 5 5 5 tre i 1 0 z-test-keyboard-harness.md 0} |
| shift-left-then-right-shrinks | pass | 3.3 | no | — | — |
| shift-right-then-left-shrinks | pass | 3.3 | no | — | — |
| shift-up-then-down-shrinks | fail | 1.0 | no | — | step 3: selStart want 1476 got 0; selLen want 19 got 1495; state={1495 0 1495 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Før… |
