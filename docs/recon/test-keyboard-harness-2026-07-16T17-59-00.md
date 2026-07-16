# Keyboard harness results

Run: 2026-07-16T17:59:00+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 92 pass, 13 fail, 0 prepare fail; total 172.4s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 1.9 | no | — | — |
| home-clears-selection | pass | 0.9 | no | — | — |
| shift-right-from-home | pass | 3.3 | no | — | — |
| shift-left-from-end | pass | 3.4 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.3 | no | — | — |
| shift-left-after-end-no-stale-anchor | pass | 3.3 | no | — | — |
| shift-down-after-arrow-down | pass | 3.4 | no | — | — |
| shift-up-after-arrow-down | pass | 3.3 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.5 | no | — | — |
| down-one-logical-line | pass | 3.3 | no | — | — |
| up-one-logical-line | pass | 3.4 | no | — | — |
| shift-down-then-up-shrinks | pass | 3.3 | no | — | — |
| shift-left-repeat-from-end | pass | 3.3 | no | — | — |
| alt-backspace-deletes-word | pass | 1.6 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.6 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.5 | no | — | — |
| cm-line-down-basic | pass | 3.3 | no | — | — |
| cm-line-up-basic | pass | 3.4 | no | — | — |
| cm-line-down-shorter | pass | 0.6 | no | — | — |
| cm-line-up-shorter | pass | 0.7 | no | — | — |
| cm-line-down-last-line | pass | 1.1 | no | — | — |
| cm-line-down-goal-col | fail | 0.8 | no | — | step 4: cursor want 6 got 7; selStart want 6 got 7; selEnd want 6 got 7; state={7 7 7 11 tre i femte 1 0 z-test-keyboard-harness.md 0} |
| cm-select-line-down | pass | 3.2 | no | — | — |
| cm-select-up-basic | pass | 3.3 | no | — | — |
| cm-select-line-down-mid | fail | 0.7 | no | — | step 4: cursor want 4 got 5; selEnd want 4 got 5; selLen want 3 got 4; state={5 1 5 9 en to tre 1 0 z-test-keyboard-harness.md 0} |
| cm-select-down-up-doc-end | fail | 0.8 | no | — | step 4: selStart want 9 got 8; state={9 8 9 9 en to tre 1 0 z-test-keyboard-harness.md 0} |
| cm-select-up-mid | pass | 0.8 | no | — | — |
| combo-alt-left | pass | 3.5 | no | — | — |
| combo-alt-right | pass | 3.4 | no | — | — |
| combo-alt-up | pass | 1.1 | no | — | — |
| combo-alt-down | pass | 1.0 | no | — | — |
| combo-ctrl-left | fail | 0.7 | no | — | step 3: cursor want 0 got 53; selStart want 0 got 53; selEnd want 0 got 53; state={53 53 53 1551 Writerdeck harness dummy — ikke i vanlig n… |
| combo-ctrl-right | fail | 0.6 | no | — | step 3: cursor want 1551 got 1027; selStart want 1551 got 1027; selEnd want 1551 got 1027; state={1027 1027 1027 1551 Writerdeck harness du… |
| combo-ctrl-up | pass | 1.1 | no | — | — |
| combo-ctrl-down | pass | 1.4 | no | — | — |
| combo-shift-alt-left | fail | 0.6 | no | — | step 3: selStart want 1168 got 1181; selLen want 17 got 4; state={1185 1181 1185 1551 Writerdeck harness dummy — ikke i vanlig notatliste  … |
| combo-shift-alt-left-repeat | fail | 0.6 | no | — | step 3: selStart want 1168 got 1181; state={1185 1181 1185 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Første avsnitt — Naïve… |
| combo-shift-alt-right | pass | 1.1 | no | — | — |
| combo-shift-alt-right-repeat | pass | 1.1 | no | — | — |
| combo-shift-alt-up | pass | 0.6 | no | — | — |
| combo-shift-alt-down | pass | 0.6 | no | — | — |
| combo-shift-ctrl-left | pass | 1.4 | no | — | — |
| combo-shift-ctrl-left-multiline | pass | 0.7 | no | — | — |
| combo-shift-ctrl-right | fail | 0.9 | no | — | step 4: selStart want 0 got 11; selLen want 11 got 0; state={11 11 11 11 hello world 1 0 z-test-keyboard-harness.md 0} |
| combo-shift-ctrl-up | pass | 1.6 | no | — | — |
| combo-shift-ctrl-down | fail | 1.0 | no | — | step 4: selStart want 0 got 9; selLen want 9 got 0; state={9 9 9 9 en to tre 1 0 z-test-keyboard-harness.md 0} |
| combo-shift-home-line | pass | 0.6 | no | — | — |
| combo-shift-end-line | pass | 0.7 | no | — | — |
| combo-ctrl-home | pass | 1.0 | no | — | — |
| combo-ctrl-end | pass | 1.0 | no | — | — |
| combo-shift-ctrl-home | pass | 0.6 | no | — | — |
| combo-shift-ctrl-end | pass | 0.5 | no | — | — |
| bs-alt-word-mid | fail | 0.6 | no | — | step 4: cursor want 1143 got 1142; textLen want 1548 got 1547; state={1142 1142 1142 1547 Writerdeck harness dummy — ikke i vanlig notatlis… |
| bs-ctrl-line-start | pass | 0.7 | no | — | — |
| bs-shift-with-selection | pass | 0.8 | no | — | — |
| bs-plain | pass | 1.0 | no | — | — |
| delete-repeat-forward | pass | 1.1 | no | — | — |
| wrap-down-one-visual-line | pass | 1.2 | no | — | — |
| wrap-down-not-jump-paragraph | pass | 1.0 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.4 | no | — | — |
| wrap-shift-down-one-visual | pass | 1.2 | no | — | — |
| wrap-shift-down-then-up-shrinks | pass | 1.3 | no | — | — |
| wrap-down-last-visual-line | pass | 1.2 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.8 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.7 | no | — | — |
| wrap-down-goal-column | pass | 1.1 | no | — | — |
| wrap-combo-alt-left-word | pass | 1.4 | no | — | — |
| wrap-combo-alt-right-word | pass | 1.3 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 0.9 | no | — | — |
| wrap-shift-left-across-wrap | pass | 1.2 | no | — | — |
| wrap-home-on-visual-line | pass | 1.0 | no | — | — |
| wrap-end-on-visual-line | pass | 0.8 | no | — | — |
| undo-redo-len | pass | 0.9 | no | — | — |
| undo-cursor-reposition | pass | 1.4 | no | — | — |
| undo-mid-line-delete | pass | 1.0 | no | — | — |
| redo-cleared-by-new-edit | pass | 1.2 | no | — | — |
| undo-after-select-delete | pass | 0.8 | no | — | — |
| gap-up-at-doc-start | pass | 1.1 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.2 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.3 | no | — | — |
| gap-plain-left-in-paragraph | pass | 3.5 | no | — | — |
| gap-plain-right-in-paragraph | pass | 3.7 | no | — | — |
| gap-plain-left-at-doc-start | pass | 1.0 | no | — | — |
| gap-plain-right-at-doc-end | pass | 1.1 | no | — | — |
| gap-collapse-selection-left | pass | 0.8 | no | — | — |
| gap-collapse-selection-right | pass | 0.9 | no | — | — |
| gap-delete-forward | pass | 1.0 | no | — | — |
| gap-delete-with-selection | pass | 1.1 | no | — | — |
| gap-select-all | pass | 0.7 | no | — | — |
| gap-enter-new-line | pass | 0.7 | no | — | — |
| gap-type-replaces-selection | pass | 0.8 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 0.9 | no | — | — |
| gap-undo-chain | pass | 0.8 | no | — | — |
| gap-unicode-alt-backspace | pass | 0.7 | no | — | — |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
| gap-alt-bs-with-selection | fail | 0.7 | no | — | step 3: selStart want 1168 got 1181; selLen want 17 got 4; state={1185 1181 1185 1551 Writerdeck harness dummy — ikke i vanlig notatliste  … |
| hw-page-right-scrolls-edit | pass | 9.1 | no | — | — |
| hw-page-left-scrolls-edit | pass | 10.8 | no | — | — |
| read-overscroll-clamps | fail | 1.0 | no | — | step 3: mode want 0 got 1; state={0 0 0 5280 Page scroll filler æøå line 0000 xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx Page scroll filler æøå line … |
| touch-down-goal-column | pass | 0.7 | no | — | — |
| touch-up-goal-column | pass | 0.5 | no | — | — |
| touch-down-shorter-line | pass | 0.5 | no | — | — |
| shift-left-then-right-shrinks | pass | 3.5 | no | — | — |
| shift-right-then-left-shrinks | pass | 3.3 | no | — | — |
| shift-up-then-down-shrinks | fail | 1.3 | no | — | step 6: selStart want 1400 got 1419; selLen want 95 got 76; state={1495 1419 1495 1551 Writerdeck harness dummy — ikke i vanlig notatliste … |
