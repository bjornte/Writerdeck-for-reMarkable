# Keyboard harness results

Run: 2026-07-17T01:00:25+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 14 pass, 1 fail, 0 prepare fail; total 18.2s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| wrap-down-one-visual-line | pass | 1.2 | no | — | — |
| wrap-down-not-jump-paragraph | fail | 0.6 | no | — | step 3: cursor want 10 got 0; state={0 0 0 199 word word word word word word word word word word word word word word word word word word wo… |
| wrap-up-from-visual-line-2 | pass | 1.5 | no | — | — |
| wrap-shift-down-one-visual | pass | 1.1 | no | — | — |
| wrap-shift-down-then-up-shrinks | pass | 3.5 | no | — | — |
| wrap-down-last-visual-line | pass | 1.2 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.9 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.7 | no | — | — |
| wrap-down-goal-column | pass | 0.9 | no | — | — |
| wrap-combo-alt-left-word | pass | 1.2 | no | — | — |
| wrap-combo-alt-right-word | pass | 1.3 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 0.9 | no | — | — |
| wrap-shift-left-across-wrap | pass | 1.3 | no | — | — |
| wrap-home-on-visual-line | pass | 0.8 | no | — | — |
| wrap-end-on-visual-line | pass | 0.9 | no | — | — |
