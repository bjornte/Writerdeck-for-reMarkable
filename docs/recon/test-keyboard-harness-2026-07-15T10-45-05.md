# Keyboard harness results

Run: 2026-07-15T10:45:06+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 1 pass, 4 fail, 0 prepare fail; total 8.9s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| undo-redo-len | fail | 5.2 | no | — | step 8: cursor want 5 got 0; state={0 0 0 5 abc d 1 0 z-test-keyboard-harness.md} |
| undo-cursor-reposition | fail | 1.1 | no | — | step 8: cursor want 0 got 2; textLen want 27 got 29; state={2 2 2 29 Blfive lines in this textedit 1 0 z-test-keyboard-harness.md} |
| undo-mid-line-delete | pass | 0.9 | no | — | — |
| redo-cleared-by-new-edit | fail | 0.7 | no | — | step 6: cursor want 3 got 0; state={0 0 0 3 abc 1 0 z-test-keyboard-harness.md} |
| undo-after-select-delete | fail | 1.0 | no | — | step 7: cursor want 6 got 0; selStart want 6 got 0; selEnd want 6 got 0; state={0 0 0 6 abcdef 1 0 z-test-keyboard-harness.md} |
