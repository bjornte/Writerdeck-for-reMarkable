# Keyboard harness results

Run: 2026-07-15T04:44:15+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 21 pass, 3 fail, 0 prepare fail; total 25.6s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| combo-alt-left | pass | 0.5 | no | — | — |
| combo-alt-right | pass | 0.5 | no | — | — |
| combo-alt-up | pass | 0.5 | no | — | — |
| combo-alt-down | pass | 0.5 | no | — | — |
| combo-ctrl-left | pass | 0.6 | no | — | — |
| combo-ctrl-right | pass | 1.0 | no | — | — |
| combo-ctrl-up | pass | 0.5 | no | — | — |
| combo-ctrl-down | pass | 0.5 | no | — | — |
| combo-shift-alt-left | pass | 0.5 | no | — | — |
| combo-shift-alt-right | pass | 0.5 | no | — | — |
| combo-shift-alt-up | pass | 0.6 | no | — | — |
| combo-shift-alt-down | fail | 0.5 | no | — | step 2: cursor want 7 got 0; selEnd want 7 got 0; selLen want 7 got 0; state={0 0 0 199 word word word word word word word word word word w… |
| combo-shift-ctrl-left | pass | 0.5 | no | — | — |
| combo-shift-ctrl-right | pass | 0.5 | no | — | — |
| combo-shift-ctrl-up | pass | 0.7 | no | — | — |
| combo-shift-ctrl-down | pass | 0.6 | no | — | — |
| combo-shift-home-line | pass | 0.7 | no | — | — |
| combo-shift-end-line | pass | 8.8 | yes | — | prepare retries |
| combo-ctrl-home | pass | 0.6 | no | — | — |
| combo-ctrl-end | pass | 0.5 | no | — | — |
| combo-shift-ctrl-home | fail | 0.6 | no | — | step 4: cursor want 4 got 0; selEnd want 4 got 0; selLen want 4 got 0; state={0 0 0 3 def 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-end | fail | 3.4 | no | — | step 2: state: HTTP 504: editor state timeout |
| wrap-combo-alt-left-word | pass | 0.9 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 0.8 | no | — | — |
