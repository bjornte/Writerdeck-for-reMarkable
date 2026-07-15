# Keyboard harness results

Run: 2026-07-15T02:04:34+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 0 pass, 1 fail, 0 prepare fail; total 0.6s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| combo-shift-ctrl-right | fail | 0.6 | no | — | step 2: selStart want 0 got 11; selLen want 11 got 0; state={11 11 11 11 hello world 1 0 z-test-keyboard-harness.md} |
