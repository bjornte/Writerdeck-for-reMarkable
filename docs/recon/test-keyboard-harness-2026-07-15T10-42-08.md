# Keyboard harness results

Run: 2026-07-15T10:42:09+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 3 pass, 2 fail, 0 prepare fail; total 7.8s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| undo-redo-len | pass | 4.4 | no | — | — |
| undo-cursor-reposition | fail | 0.8 | no | — | step 6: cursor want 31 got 8; state={8 8 8 31 Blahfive lines in this textedit 1 0 z-test-keyboard-harness.md} |
| undo-mid-line-delete | fail | 1.0 | no | — | step 9: cursor want 4 got 7; state={7 7 7 7 abc def 1 0 z-test-keyboard-harness.md} |
| redo-cleared-by-new-edit | pass | 1.0 | no | — | — |
| undo-after-select-delete | pass | 0.7 | no | — | — |
