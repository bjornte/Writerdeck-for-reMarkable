# Keyboard harness results

Run: 2026-07-14T23:24:25+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 6 pass, 16 fail, 0 prepare fail; total 10.6s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| combo-alt-left | fail | 0.7 | no | — | step 3: cursor want 6 got 5; selStart want 6 got 5; selEnd want 6 got 5; state={5 5 5 5 1 0 z-test-keyboard-harness.md} |
| combo-alt-right | fail | 0.4 | no | — | step 2: cursor want 6 got 0; selStart want 6 got 0; selEnd want 6 got 0; state={0 0 0 6 1 0 z-test-keyboard-harness.md} |
| combo-alt-up | pass | 0.5 | no | — | — |
| combo-alt-down | fail | 0.5 | no | — | step 2: cursor want 7 got 0; selStart want 7 got 0; selEnd want 7 got 0; state={0 0 0 7 1 0 z-test-keyboard-harness.md} |
| combo-ctrl-left | pass | 0.5 | no | — | — |
| combo-ctrl-right | fail | 0.4 | no | — | step 2: cursor want 11 got 0; selStart want 11 got 0; selEnd want 11 got 0; state={0 0 0 11 1 0 z-test-keyboard-harness.md} |
| combo-ctrl-up | pass | 0.5 | no | — | — |
| combo-ctrl-down | fail | 0.4 | no | — | step 2: cursor want 13 got 0; selStart want 13 got 0; selEnd want 13 got 0; state={0 0 0 13 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-left | fail | 0.5 | no | — | step 3: cursor want 11 got 5; selStart want 6 got 5; selEnd want 11 got 5; selLen want 5 got 0; state={5 5 5 5 1 0 z-test-keyboard-harness.… |
| combo-shift-alt-right | fail | 0.3 | no | — | step 2: cursor want 6 got 0; selEnd want 6 got 0; selLen want 6 got 0; state={0 0 0 6 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-up | fail | 0.5 | no | — | step 3: cursor want 12 got 0; selEnd want 12 got 0; selLen want 12 got 0; state={0 0 0 7 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-down | fail | 0.4 | no | — | step 2: cursor want 7 got 0; selEnd want 7 got 0; selLen want 7 got 0; state={0 0 0 7 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-left | fail | 0.5 | no | — | step 3: cursor want 11 got 0; selEnd want 11 got 0; selLen want 11 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-right | fail | 0.4 | no | — | step 2: cursor want 11 got 0; selEnd want 11 got 0; selLen want 11 got 0; state={0 0 0 11 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-up | fail | 0.5 | no | — | step 3: cursor want 13 got 0; selEnd want 13 got 0; selLen want 13 got 0; state={0 0 0 10 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-down | fail | 0.3 | no | — | step 2: cursor want 13 got 0; selEnd want 13 got 0; selLen want 13 got 0; state={0 0 0 13 1 0 z-test-keyboard-harness.md} |
| combo-shift-home-line | pass | 0.7 | no | — | — |
| combo-shift-end-line | pass | 0.6 | no | — | — |
| combo-ctrl-home | pass | 0.5 | no | — | — |
| combo-ctrl-end | fail | 0.4 | no | — | step 2: cursor want 7 got 0; selStart want 7 got 0; selEnd want 7 got 0; state={0 0 0 7 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-home | fail | 0.6 | no | — | step 3: cursor want 4 got 0; selEnd want 4 got 0; selLen want 4 got 0; state={0 0 0 3 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-end | fail | 0.6 | no | — | step 3: cursor want 7 got 0; selStart want 4 got 0; selEnd want 7 got 0; selLen want 3 got 0; state={0 0 0 3 1 0 z-test-keyboard-harness.md} |
