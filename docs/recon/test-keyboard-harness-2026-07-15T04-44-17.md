# Keyboard harness results

Run: 2026-07-15T04:44:17+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 9 pass, 5 fail, 0 prepare fail; total 21.4s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| wrap-down-one-visual-line | fail | 0.6 | no | — | step 3: cursor want 20 got 11; textLen want 199 got 11; state={11 11 11 11 hello world 1 0 z-test-keyboard-harness.md} |
| wrap-down-not-jump-paragraph | fail | 3.7 | no | — | step 3: state: HTTP 504: editor state timeout |
| wrap-up-from-visual-line-2 | pass | 4.1 | no | — | — |
| wrap-shift-down-one-visual | pass | 0.8 | no | — | — |
| wrap-shift-down-then-up-shrinks | pass | 0.9 | no | — | — |
| wrap-down-last-visual-line | fail | 0.9 | no | — | step 3: cursor want 199 got 0; selStart want 199 got 0; selEnd want 199 got 0; textLen want 199 got 7; state={0 0 0 7 abc def 1 0 z-test-ke… |
| wrap-shift-down-last-to-eof | fail | 0.9 | no | — | step 3: cursor want 199 got 7; selEnd want 199 got 7; textLen want 199 got 7; state={7 6 7 7 abc def 1 0 z-test-keyboard-harness.md} |
| wrap-mixed-newline-and-wrap | pass | 0.7 | no | — | — |
| wrap-down-goal-column | pass | 0.9 | no | — | — |
| wrap-combo-alt-left-word | pass | 0.8 | no | — | — |
| wrap-combo-ctrl-bs-line | fail | 0.7 | no | suspect | step 2: cursor want 199 got 0; textLen want 199 got 7; state={0 0 0 7 abc def 1 0 z-test-keyboard-harness.md}; may have poisoned next scena… |
| wrap-shift-left-across-wrap | pass | 4.6 | yes | wrap-combo-ctrl-bs-line | prepare retries; cascade suspect after wrap-combo-ctrl-bs-line |
| wrap-home-on-visual-line | pass | 0.8 | no | — | — |
| wrap-end-on-visual-line | pass | 0.9 | no | — | — |
