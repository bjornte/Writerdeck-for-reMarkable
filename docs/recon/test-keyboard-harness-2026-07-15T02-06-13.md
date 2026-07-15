# Keyboard harness results

Run: 2026-07-15T02:06:13+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 20 pass, 4 fail, 0 prepare fail; total 13.2s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| combo-alt-left | pass | 0.5 | no | — | — |
| combo-alt-right | pass | 0.5 | no | — | — |
| combo-alt-up | pass | 0.5 | no | — | — |
| combo-alt-down | pass | 0.5 | no | — | — |
| combo-ctrl-left | pass | 0.5 | no | — | — |
| combo-ctrl-right | pass | 0.5 | no | — | — |
| combo-ctrl-up | pass | 0.5 | no | — | — |
| combo-ctrl-down | pass | 0.5 | no | — | — |
| combo-shift-alt-left | pass | 0.5 | no | — | — |
| combo-shift-alt-right | pass | 0.5 | no | — | — |
| combo-shift-alt-up | fail | 0.5 | no | — | step 3: cursor want 12 got 5; selEnd want 12 got 5; selLen want 12 got 5; state={5 0 5 12 para1  para2 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-down | pass | 0.5 | no | — | — |
| combo-shift-ctrl-left | pass | 0.5 | no | — | — |
| combo-shift-ctrl-right | pass | 0.5 | no | — | — |
| combo-shift-ctrl-up | fail | 0.6 | no | — | step 3: cursor want 13 got 3; selEnd want 13 got 3; selLen want 13 got 3; state={3 0 3 13 one two three 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-down | pass | 0.5 | no | — | — |
| combo-shift-home-line | pass | 0.6 | no | — | — |
| combo-shift-end-line | pass | 0.6 | no | — | — |
| combo-ctrl-home | pass | 0.5 | no | — | — |
| combo-ctrl-end | pass | 0.5 | no | — | — |
| combo-shift-ctrl-home | fail | 0.6 | no | — | step 3: cursor want 4 got 0; selEnd want 4 got 0; selLen want 4 got 0; state={0 0 0 3 def 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-end | fail | 0.5 | no | — | step 3: cursor want 7 got 0; selStart want 4 got 0; selEnd want 7 got 0; selLen want 3 got 0; state={0 0 0 3 def 1 0 z-test-keyboard-harnes… |
| wrap-combo-alt-left-word | pass | 0.8 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 0.8 | no | — | — |
