# Keyboard harness results

Run: 2026-07-16T01:54:50+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 85 pass, 17 fail, 0 prepare fail; total 103.1s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 2.8 | no | — | — |
| home-clears-selection | pass | 0.8 | no | — | — |
| shift-right-from-home | pass | 1.0 | no | — | — |
| shift-left-from-end | pass | 1.1 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 1.1 | no | — | — |
| shift-left-after-end-no-stale-anchor | pass | 1.0 | no | — | — |
| shift-down-after-arrow-down | pass | 1.0 | no | — | — |
| shift-up-after-arrow-down | fail | 0.8 | no | — | step 6: cursor want 794 got 775; selStart want 737 got 756; selEnd want 794 got 775; state={775 756 775 832 Writerdeck harness dummy — ikke… |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | pass | 1.0 | no | — | — |
| up-one-logical-line | pass | 1.0 | no | — | — |
| shift-down-then-up-shrinks | pass | 1.5 | no | — | — |
| shift-left-repeat-from-end | pass | 1.1 | no | — | — |
| alt-backspace-deletes-word | pass | 1.1 | no | — | — |
| ctrl-backspace-deletes-line | fail | 0.8 | no | — | step 6: cursor want 795 got 776; state={776 776 776 776 Writerdeck harness dummy — ikke i vanlig notatliste Naïve café résumé: æøå på Færøy… |
| shift-left-repeat-mid-doc | pass | 1.2 | no | — | — |
| cm-line-down-basic | pass | 1.1 | no | — | — |
| cm-line-up-basic | pass | 1.0 | no | — | — |
| cm-line-down-shorter | pass | 0.7 | no | — | — |
| cm-line-up-shorter | pass | 0.7 | no | — | — |
| cm-line-down-last-line | pass | 1.0 | no | — | — |
| cm-line-down-goal-col | fail | 0.8 | no | — | step 4: cursor want 6 got 7; selStart want 6 got 7; selEnd want 6 got 7; state={7 7 7 11 tre i femte 1 0 z-test-keyboard-harness.md 0} |
| cm-select-line-down | pass | 1.0 | no | — | — |
| cm-select-line-down-mid | fail | 0.6 | no | — | step 4: cursor want 4 got 5; selEnd want 4 got 5; selLen want 3 got 4; state={5 1 5 9 en to tre 1 0 z-test-keyboard-harness.md 0} |
| cm-select-down-up-doc-end | fail | 0.6 | no | — | step 4: selStart want 9 got 8; state={9 8 9 9 en to tre 1 0 z-test-keyboard-harness.md 0} |
| cm-select-up-basic | fail | 0.7 | no | — | step 3: selStart want 813 got 795; state={832 795 832 832 Writerdeck harness dummy — ikke i vanlig notatliste Naïve café résumé: æøå på Fær… |
| cm-select-up-mid | pass | 0.8 | no | — | — |
| combo-alt-left | pass | 1.1 | no | — | — |
| combo-alt-right | pass | 1.1 | no | — | — |
| combo-alt-up | pass | 1.2 | no | — | — |
| combo-alt-down | pass | 1.0 | no | — | — |
| combo-ctrl-left | fail | 0.6 | no | — | step 3: cursor want 0 got 52; selStart want 0 got 52; selEnd want 0 got 52; state={52 52 52 832 Writerdeck harness dummy — ikke i vanlig no… |
| combo-ctrl-right | fail | 0.7 | no | — | step 3: cursor want 832 got 321; selStart want 832 got 321; selEnd want 832 got 321; state={321 321 321 832 Writerdeck harness dummy — ikke… |
| combo-ctrl-up | pass | 1.1 | no | — | — |
| combo-ctrl-down | pass | 1.1 | no | — | — |
| combo-shift-alt-left | pass | 1.0 | no | — | — |
| combo-shift-alt-left-repeat | pass | 1.0 | no | — | — |
| combo-shift-alt-right | pass | 1.1 | no | — | — |
| combo-shift-alt-right-repeat | pass | 1.1 | no | — | — |
| combo-shift-alt-up | pass | 0.6 | no | — | — |
| combo-shift-alt-down | pass | 0.5 | no | — | — |
| combo-shift-ctrl-left | pass | 1.5 | no | — | — |
| combo-shift-ctrl-left-multiline | pass | 0.7 | no | — | — |
| combo-shift-ctrl-right | fail | 0.9 | no | — | step 4: selStart want 0 got 11; selLen want 11 got 0; state={11 11 11 11 hello world 1 0 z-test-keyboard-harness.md 0} |
| combo-shift-ctrl-up | pass | 1.6 | no | — | — |
| combo-shift-ctrl-down | fail | 0.8 | no | — | step 4: selStart want 0 got 9; selLen want 9 got 0; state={9 9 9 9 en to tre 1 0 z-test-keyboard-harness.md 0} |
| combo-shift-home-line | pass | 0.6 | no | — | — |
| combo-shift-end-line | pass | 0.7 | no | — | — |
| combo-ctrl-home | pass | 1.0 | no | — | — |
| combo-ctrl-end | pass | 1.0 | no | — | — |
| combo-shift-ctrl-home | pass | 0.6 | no | — | — |
| combo-shift-ctrl-end | pass | 0.6 | no | — | — |
| bs-alt-word-mid | fail | 0.6 | no | — | step 4: cursor want 436 got 435; textLen want 829 got 824; state={435 435 435 824 Writerdeck harness dummy — ikke i vanlig notatliste Naïve… |
| bs-ctrl-line-start | pass | 0.7 | no | — | — |
| bs-shift-with-selection | pass | 0.9 | no | — | — |
| bs-plain | pass | 1.0 | no | — | — |
| delete-repeat-forward | pass | 1.1 | no | — | — |
| wrap-down-one-visual-line | pass | 1.2 | no | — | — |
| wrap-down-not-jump-paragraph | pass | 1.0 | no | — | — |
| wrap-up-from-visual-line-2 | fail | 1.2 | no | — | step 7: cursorMax want <= 65 got 80; state={80 80 80 199 word word word word word word word word word word word word word word word word wo… |
| wrap-shift-down-one-visual | fail | 0.9 | no | — | step 5: cursor want 60 got 199; selEnd want 60 got 199; state={199 0 199 199 word word word word word word word word word word word word wo… |
| wrap-shift-down-then-up-shrinks | fail | 1.2 | no | — | step 8: selLenMin want >= 1 got 0; state={0 0 0 199 word word word word word word word word word word word word word word word word word wo… |
| wrap-down-last-visual-line | pass | 1.3 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.9 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.8 | no | — | — |
| wrap-down-goal-column | pass | 0.9 | no | — | — |
| wrap-combo-alt-left-word | pass | 1.4 | no | — | — |
| wrap-combo-alt-right-word | pass | 1.3 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 0.8 | no | — | — |
| wrap-shift-left-across-wrap | pass | 1.4 | no | — | — |
| wrap-home-on-visual-line | pass | 0.9 | no | — | — |
| wrap-end-on-visual-line | pass | 0.9 | no | — | — |
| undo-redo-len | pass | 1.0 | no | — | — |
| undo-cursor-reposition | pass | 1.4 | no | — | — |
| undo-mid-line-delete | pass | 0.9 | no | — | — |
| redo-cleared-by-new-edit | pass | 1.2 | no | — | — |
| undo-after-select-delete | pass | 0.8 | no | — | — |
| gap-up-at-doc-start | pass | 1.0 | no | — | — |
| gap-plain-left-moves-caret | pass | 1.1 | no | — | — |
| gap-plain-right-moves-caret | pass | 1.0 | no | — | — |
| gap-plain-left-at-doc-start | pass | 1.0 | no | — | — |
| gap-plain-right-at-doc-end | pass | 1.3 | no | — | — |
| gap-collapse-selection-left | pass | 0.8 | no | — | — |
| gap-collapse-selection-right | pass | 0.9 | no | — | — |
| gap-delete-forward | pass | 1.1 | no | — | — |
| gap-delete-with-selection | pass | 1.1 | no | — | — |
| gap-select-all | pass | 0.7 | no | — | — |
| gap-enter-new-line | pass | 0.7 | no | — | — |
| gap-type-replaces-selection | pass | 0.8 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 0.9 | no | — | — |
| gap-undo-chain | pass | 0.8 | no | — | — |
| gap-unicode-alt-backspace | pass | 0.7 | no | — | — |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
| gap-alt-bs-with-selection | pass | 0.8 | no | — | — |
| hw-page-right-scrolls-edit | pass | 2.2 | no | — | — |
| hw-page-left-scrolls-edit | pass | 3.4 | no | — | — |
| read-overscroll-clamps | fail | 2.6 | no | — | step 3: mode want 0 got 1; state={0 0 0 5280 Page scroll filler æøå line 0000 xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx Page scroll filler æøå line … |
| touch-down-goal-column | pass | 0.6 | no | — | — |
| touch-up-goal-column | pass | 0.5 | no | — | — |
| touch-down-shorter-line | pass | 0.5 | no | — | — |
| shift-left-then-right-shrinks | fail | 1.4 | no | — | step 11: cursor want 575 got 579; selEnd want 575 got 579; selLen want 4 got 8; state={579 571 579 832 Writerdeck harness dummy — ikke i va… |
| shift-right-then-left-shrinks | fail | 1.2 | no | — | step 9: cursor want 584 got 585; selStart want 578 got 577; selEnd want 584 got 585; selLen want 6 got 8; state={585 577 585 832 Writerdeck… |
