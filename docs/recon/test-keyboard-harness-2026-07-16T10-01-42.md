# Keyboard harness results

Run: 2026-07-16T10:01:48+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 72 pass, 33 fail, 0 prepare fail; total 156.7s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 6.0 | no | — | — |
| home-clears-selection | pass | 1.0 | no | — | — |
| shift-right-from-home | fail | 2.0 | no | — | step 11: cursor want 1282 got 1283; selStart want 1282 got 1281; selEnd want 1282 got 1283; selLen want 0 got 2; state={1283 1281 1283 1551… |
| shift-left-from-end | fail | 1.6 | no | — | step 11: cursor want 1313 got 1312; selStart want 1313 got 1312; selEnd want 1313 got 1312; state={1312 1312 1312 1551 Writerdeck harness d… |
| shift-right-after-home-no-stale-anchor | fail | 1.5 | no | — | step 11: cursor want 1282 got 1283; selStart want 1282 got 1281; selEnd want 1282 got 1283; selLen want 0 got 2; state={1283 1281 1283 1551… |
| shift-left-after-end-no-stale-anchor | fail | 1.6 | no | — | step 11: cursor want 1313 got 1312; selStart want 1313 got 1312; selEnd want 1313 got 1312; state={1312 1312 1312 1551 Writerdeck harness d… |
| shift-down-after-arrow-down | fail | 2.4 | no | — | step 16: cursor want 1362 got 1343; selStart want 1324 got 1343; selEnd want 1362 got 1343; selLen want 38 got 0; state={1343 1343 1343 155… |
| shift-up-after-arrow-down | fail | 1.1 | no | — | step 6: cursor want 1495 got 1457; selStart want 1400 got 1438; selEnd want 1495 got 1457; selLen want 95 got 19; state={1457 1438 1457 155… |
| ctrl-shift-left-select-line | pass | 0.7 | no | — | — |
| down-one-logical-line | pass | 3.3 | no | — | — |
| up-one-logical-line | pass | 3.3 | no | — | — |
| shift-down-then-up-shrinks | fail | 2.6 | no | — | step 16: cursor want 1362 got 1343; selStart want 1324 got 1343; selEnd want 1362 got 1343; selLen want 38 got 0; state={1343 1343 1343 155… |
| shift-left-repeat-from-end | fail | 1.6 | no | — | step 11: cursor want 1323 got 1322; selStart want 1323 got 1322; selEnd want 1323 got 1322; state={1322 1322 1322 1551 Writerdeck harness d… |
| alt-backspace-deletes-word | fail | 1.1 | no | — | step 6: cursor want 1155 got 1161; state={1161 1161 1161 1518 Writerdeck harness dummy — ikke i vanlig notatliste  Første avsnitt — Naïve c… |
| ctrl-backspace-deletes-line | fail | 1.2 | no | — | step 6: cursorMin want >= 1457 got 1438; state={1438 1438 1438 1438 Writerdeck harness dummy — ikke i vanlig notatliste  Første avsnitt — N… |
| shift-left-repeat-mid-doc | fail | 1.6 | no | — | step 11: cursor want 1437 got 1436; selStart want 1437 got 1436; selEnd want 1437 got 1436; state={1436 1436 1436 1551 Writerdeck harness d… |
| cm-line-down-basic | pass | 3.2 | no | — | — |
| cm-line-up-basic | pass | 3.3 | no | — | — |
| cm-line-down-shorter | pass | 0.8 | no | — | — |
| cm-line-up-shorter | pass | 0.7 | no | — | — |
| cm-line-down-last-line | pass | 1.0 | no | — | — |
| cm-line-down-goal-col | fail | 0.7 | no | — | step 4: cursor want 6 got 7; selStart want 6 got 7; selEnd want 6 got 7; state={7 7 7 11 tre i femte 1 0 z-test-keyboard-harness.md 0} |
| cm-select-line-down | fail | 2.2 | no | — | step 11: cursor want 1362 got 1381; selEnd want 1362 got 1381; selLen want 0 got 19; state={1381 1362 1381 1551 Writerdeck harness dummy — … |
| cm-select-up-basic | fail | 1.2 | no | — | step 6: cursor want 1495 got 1457; selStart want 1400 got 1438; selEnd want 1495 got 1457; selLen want 95 got 19; state={1457 1438 1457 155… |
| cm-select-line-down-mid | fail | 0.6 | no | — | step 4: cursor want 4 got 5; selEnd want 4 got 5; selLen want 3 got 4; state={5 1 5 9 en to tre 1 0 z-test-keyboard-harness.md 0} |
| cm-select-down-up-doc-end | fail | 0.7 | no | — | step 4: selStart want 9 got 8; state={9 8 9 9 en to tre 1 0 z-test-keyboard-harness.md 0} |
| cm-select-up-mid | pass | 0.8 | no | — | — |
| combo-alt-left | pass | 3.5 | no | — | — |
| combo-alt-right | pass | 3.7 | no | — | — |
| combo-alt-up | pass | 1.4 | no | — | — |
| combo-alt-down | pass | 1.0 | no | — | — |
| combo-ctrl-left | fail | 0.7 | no | — | step 3: cursor want 0 got 53; selStart want 0 got 53; selEnd want 0 got 53; state={53 53 53 1551 Writerdeck harness dummy — ikke i vanlig n… |
| combo-ctrl-right | fail | 0.7 | no | — | step 3: cursor want 1551 got 1027; selStart want 1551 got 1027; selEnd want 1551 got 1027; state={1027 1027 1027 1551 Writerdeck harness du… |
| combo-ctrl-up | pass | 1.3 | no | — | — |
| combo-ctrl-down | pass | 1.1 | no | — | — |
| combo-shift-alt-left | fail | 0.6 | no | — | step 3: selStart want 1168 got 1181; selLen want 17 got 4; state={1185 1181 1185 1551 Writerdeck harness dummy — ikke i vanlig notatliste  … |
| combo-shift-alt-left-repeat | fail | 0.7 | no | — | step 3: selStart want 1168 got 1181; state={1185 1181 1185 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Første avsnitt — Naïve… |
| combo-shift-alt-right | pass | 1.1 | no | — | — |
| combo-shift-alt-right-repeat | pass | 1.0 | no | — | — |
| combo-shift-alt-up | pass | 0.6 | no | — | — |
| combo-shift-alt-down | pass | 0.5 | no | — | — |
| combo-shift-ctrl-left | pass | 1.5 | no | — | — |
| combo-shift-ctrl-left-multiline | pass | 0.7 | no | — | — |
| combo-shift-ctrl-right | fail | 0.9 | no | — | step 4: selStart want 0 got 11; selLen want 11 got 0; state={11 11 11 11 hello world 1 0 z-test-keyboard-harness.md 0} |
| combo-shift-ctrl-up | pass | 1.6 | no | — | — |
| combo-shift-ctrl-down | fail | 0.9 | no | — | step 4: selStart want 0 got 9; selLen want 9 got 0; state={9 9 9 9 en to tre 1 0 z-test-keyboard-harness.md 0} |
| combo-shift-home-line | pass | 0.6 | no | — | — |
| combo-shift-end-line | pass | 0.7 | no | — | — |
| combo-ctrl-home | pass | 1.0 | no | — | — |
| combo-ctrl-end | pass | 1.1 | no | — | — |
| combo-shift-ctrl-home | pass | 0.7 | no | — | — |
| combo-shift-ctrl-end | pass | 0.6 | no | — | — |
| bs-alt-word-mid | fail | 0.7 | no | — | step 4: cursor want 1143 got 1142; textLen want 1548 got 1543; state={1142 1142 1142 1543 Writerdeck harness dummy — ikke i vanlig notatlis… |
| bs-ctrl-line-start | pass | 0.8 | no | — | — |
| bs-shift-with-selection | pass | 0.9 | no | — | — |
| bs-plain | pass | 1.1 | no | — | — |
| delete-repeat-forward | pass | 1.0 | no | — | — |
| wrap-down-one-visual-line | pass | 1.1 | no | — | — |
| wrap-down-not-jump-paragraph | pass | 1.0 | no | — | — |
| wrap-up-from-visual-line-2 | fail | 1.2 | no | — | step 7: cursorMax want <= 65 got 80; state={80 80 80 199 word word word word word word word word word word word word word word word word wo… |
| wrap-shift-down-one-visual | fail | 0.9 | no | — | step 5: cursor want 60 got 199; selEnd want 60 got 199; state={199 0 199 199 word word word word word word word word word word word word wo… |
| wrap-shift-down-then-up-shrinks | fail | 1.2 | no | — | step 8: selLenMin want >= 1 got 0; state={0 0 0 199 word word word word word word word word word word word word word word word word word wo… |
| wrap-down-last-visual-line | pass | 1.4 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.8 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.8 | no | — | — |
| wrap-down-goal-column | pass | 0.8 | no | — | — |
| wrap-combo-alt-left-word | pass | 1.3 | no | — | — |
| wrap-combo-alt-right-word | pass | 1.3 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 1.0 | no | — | — |
| wrap-shift-left-across-wrap | pass | 1.5 | no | — | — |
| wrap-home-on-visual-line | pass | 0.8 | no | — | — |
| wrap-end-on-visual-line | pass | 0.8 | no | — | — |
| undo-redo-len | pass | 1.1 | no | — | — |
| undo-cursor-reposition | pass | 1.4 | no | — | — |
| undo-mid-line-delete | pass | 0.9 | no | — | — |
| redo-cleared-by-new-edit | pass | 1.0 | no | — | — |
| undo-after-select-delete | pass | 0.8 | no | — | — |
| gap-up-at-doc-start | pass | 1.3 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.4 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.4 | no | — | — |
| gap-plain-left-in-paragraph | fail | 2.1 | no | — | step 14: cursor want 175 got 176; selStart want 175 got 176; selEnd want 175 got 176; state={176 176 176 1551 Writerdeck harness dummy — ik… |
| gap-plain-right-in-paragraph | pass | 3.4 | no | — | — |
| gap-plain-left-at-doc-start | pass | 1.1 | no | — | — |
| gap-plain-right-at-doc-end | pass | 1.2 | no | — | — |
| gap-collapse-selection-left | pass | 0.9 | no | — | — |
| gap-collapse-selection-right | pass | 1.0 | no | — | — |
| gap-delete-forward | pass | 1.0 | no | — | — |
| gap-delete-with-selection | pass | 1.0 | no | — | — |
| gap-select-all | pass | 0.7 | no | — | — |
| gap-enter-new-line | pass | 0.8 | no | — | — |
| gap-type-replaces-selection | pass | 0.8 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 1.0 | no | — | — |
| gap-undo-chain | pass | 0.9 | no | — | — |
| gap-unicode-alt-backspace | pass | 0.5 | no | — | — |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
| gap-alt-bs-with-selection | fail | 1.4 | no | — | step 3: selStart want 1168 got 1181; selLen want 17 got 4; state={1185 1181 1185 1551 Writerdeck harness dummy — ikke i vanlig notatliste  … |
| hw-page-right-scrolls-edit | fail | 7.0 | no | — | step 17: cursor want 0 got 1781; state={1781 1781 1781 26400 Page scroll filler æøå line 0000 xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx Page scroll … |
| hw-page-left-scrolls-edit | pass | 10.8 | no | — | — |
| read-overscroll-clamps | fail | 4.9 | no | — | step 3: mode want 0 got 1; state={0 0 0 5280 Page scroll filler æøå line 0000 xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx Page scroll filler æøå line … |
| touch-down-goal-column | pass | 0.6 | no | — | — |
| touch-up-goal-column | pass | 0.6 | no | — | — |
| touch-down-shorter-line | pass | 0.5 | no | — | — |
| shift-left-then-right-shrinks | fail | 1.7 | no | — | step 11: cursor want 1297 got 1296; selStart want 1297 got 1296; selEnd want 1297 got 1296; state={1296 1296 1296 1551 Writerdeck harness d… |
| shift-right-then-left-shrinks | fail | 1.6 | no | — | step 11: cursor want 1297 got 1298; selStart want 1297 got 1296; selEnd want 1297 got 1298; selLen want 0 got 2; state={1298 1296 1298 1551… |
| shift-up-then-down-shrinks | fail | 1.1 | no | — | step 6: cursor want 1495 got 1457; selStart want 1400 got 1438; selEnd want 1495 got 1457; selLen want 95 got 19; state={1457 1438 1457 155… |
