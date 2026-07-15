# Keyboard harness results

Run: 2026-07-15T10:36:26+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 3 pass, 2 fail, 0 prepare fail; total 7.9s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| undo-redo-len | pass | 4.6 | no | — | — |
| undo-cursor-reposition | fail | 0.8 | no | — | step 6: cursor want 31 got 8; state={8 8 8 31 Blahfive lines in this textedit 1 0 z-test-keyboard-harness.md} |
| undo-mid-line-delete | pass | 0.9 | no | — | — |
| redo-cleared-by-new-edit | fail | 0.7 | no | — | step 5: cursor want 3 got 0; textLen want 3 got 1; state={0 0 0 1 a 1 0 z-test-keyboard-harness.md} |
| undo-after-select-delete | pass | 0.9 | no | — | — |
