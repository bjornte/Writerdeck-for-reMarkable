# Keyboard harness results

Run: 2026-07-16T11:33:21+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 2 pass, 5 fail, 0 prepare fail; total 13.2s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 4.0 | no | — | — |
| home-clears-selection | pass | 0.9 | no | — | — |
| shift-right-from-home | fail | 1.6 | no | — | step 11: cursor want 1282 got 1283; selStart want 1282 got 1281; selEnd want 1282 got 1283; selLen want 0 got 2; state={1283 1281 1283 1551… |
| shift-left-from-end | fail | 1.6 | no | — | step 11: cursor want 1313 got 1312; selStart want 1313 got 1312; selEnd want 1313 got 1312; state={1312 1312 1312 1551 Writerdeck harness d… |
| shift-right-after-home-no-stale-anchor | fail | 1.6 | no | — | step 11: cursor want 1282 got 1283; selStart want 1282 got 1281; selEnd want 1282 got 1283; selLen want 0 got 2; state={1283 1281 1283 1551… |
| shift-left-after-end-no-stale-anchor | fail | 1.6 | no | — | step 11: cursor want 1313 got 1312; selStart want 1313 got 1312; selEnd want 1313 got 1312; state={1312 1312 1312 1551 Writerdeck harness d… |
| shift-down-after-arrow-down | fail | 2.0 | no | — | step 14: state: HTTP 409: no active editor session; editor unreachable after fail: HTTP 409: no active editor session |
