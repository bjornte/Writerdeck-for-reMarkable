# Keyboard harness results

Run: 2026-07-15T23:44:20+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 4 pass, 1 fail, 0 prepare fail; total 14.6s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| undo-redo-len | pass | 7.3 | no | — | — |
| undo-cursor-reposition | fail | 2.7 | no | — | step 8: cursor want 0 got 2; textLen want 27 got 29; state={2 2 2 29 Blfive lines in this textedit 1 0 z-test-keyboard-harness.md} |
| undo-mid-line-delete | pass | 1.6 | no | — | — |
| redo-cleared-by-new-edit | pass | 1.8 | no | — | — |
| undo-after-select-delete | pass | 1.2 | no | — | — |
