# Keyboard harness results

Run: 2026-07-15T23:49:24+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 0 pass, 1 fail, 0 prepare fail; total 1.2s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| gap-redo-shift-ctrl-z | fail | 1.2 | no | — | step 8: cursor want 0 got 3; textLen want 0 got 3; state={3 3 3 3 abc 1 0 z-test-keyboard-harness.md} |
