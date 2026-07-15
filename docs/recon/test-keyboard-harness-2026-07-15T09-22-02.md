# Keyboard harness results

Run: 2026-07-15T09:22:02+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 0 pass, 1 fail, 0 prepare fail; total 0.8s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| bs-plain | fail | 0.8 | no | — | step 3: cursor want 2 got 4; textLen want 2 got 4; text want "ab" got "abcd"; state={4 3 4 4 abcd 1 0 z-test-keyboard-harness.md} |
