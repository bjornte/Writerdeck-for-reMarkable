# Keyboard harness results

Run: 2026-07-16T22:24:28+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 5 pass, 3 fail, 0 prepare fail; total 53.1s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.6 | no | — | — |
| home-clears-selection | pass | 0.8 | no | — | — |
| shift-right-from-home | pass | 3.2 | no | — | — |
| shift-left-from-end | pass | 3.2 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.2 | no | — | — |
| shift-down-after-arrow-down | fail | 3.2 | no | — | step 21: cursor want 1362 got 0; selStart want 1362 got 0; selEnd want 1362 got 0; state={0 0 0 1551 Writerdeck harness dummy — ikke i vanl… |
| shift-up-after-arrow-down | fail | 6.9 | no | — | reset uni1: setCursor: editor-cmd harnesssetcursor HTTP 504 |
| ctrl-shift-left-select-line | fail | 32.1 | no | — | step 3: selStart want 1272 got 0; selEnd want 1323 got 0; selLen want 51 got 0; state={0 0 0 0  1 1  0}; editor in lobby after fail |
