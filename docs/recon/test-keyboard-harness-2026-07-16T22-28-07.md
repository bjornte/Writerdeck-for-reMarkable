# Keyboard harness results

Run: 2026-07-16T22:28:07+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 5 pass, 10 fail, 0 prepare fail; total 23.2s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| wrap-down-one-visual-line | fail | 7.9 | yes | — | step 3: cursor want 20 got 52; textLen want 199 got 1551; state={52 52 52 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Første … |
| wrap-down-not-jump-paragraph | fail | 1.1 | no | — | step 3: cursor want 20 got 10; state={10 10 10 199 word word word word word word word word word word word word word word word word word wor… |
| wrap-up-from-visual-line-2 | fail | 0.9 | no | — | step 3: cursor want 140 got 70; state={70 70 70 199 word word word word word word word word word word word word word word word word word wo… |
| wrap-shift-down-one-visual | fail | 0.6 | no | — | step 3: cursor want 20 got 10; selEnd want 20 got 10; state={10 0 10 199 word word word word word word word word word word word word word w… |
| wrap-shift-down-then-up-shrinks | fail | 1.1 | no | — | step 6: textLen want 199 got 1551; selLenMin want >= 1 got 0; state={0 0 0 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Første… |
| wrap-down-last-visual-line | pass | 1.2 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.7 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.6 | no | — | — |
| wrap-down-goal-column | fail | 0.8 | no | — | step 5: cursor want 24 got 14; state={14 14 14 177 abword word word word word word word word word word word word word word word word word w… |
| wrap-combo-alt-left-word | fail | 0.8 | no | — | step 4: cursorMin want >= 1 got 0; textLen want 199 got 1551; state={0 0 0 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Første… |
| wrap-combo-alt-right-word | pass | 4.8 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 0.8 | no | — | — |
| wrap-shift-left-across-wrap | fail | 0.6 | no | — | step 3: cursor want 20 got 10; state={10 10 10 199 word word word word word word word word word word word word word word word word word wor… |
| wrap-home-on-visual-line | fail | 0.6 | no | — | step 3: cursor want 20 got 10; state={10 10 10 199 word word word word word word word word word word word word word word word word word wor… |
| wrap-end-on-visual-line | fail | 0.6 | no | — | step 3: cursor want 20 got 10; state={10 10 10 199 word word word word word word word word word word word word word word word word word wor… |
