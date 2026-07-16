# Keyboard harness results

Run: 2026-07-16T22:27:49+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 6 pass, 3 fail, 1 prepare fail; total 33.7s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.6 | no | — | — |
| home-clears-selection | pass | 0.8 | no | — | — |
| shift-right-from-home | pass | 3.2 | no | — | — |
| shift-left-from-end | pass | 3.2 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 4.0 | no | — | — |
| shift-down-after-arrow-down | fail | 3.5 | no | suspect | step 21: cursor want 1362 got 0; selStart want 1362 got 0; selEnd want 1362 got 0; state={0 0 0 1551 Writerdeck harness dummy — ikke i vanl… |
| shift-up-after-arrow-down | fail | 6.5 | yes | shift-down-after-arrow-down | step 14: cursor want 1495 got 0; selStart want 1438 got 0; selEnd want 1495 got 0; selLen want 57 got 0; state={0 0 0 1551 Writerdeck harne… |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | fail | 3.2 | no | suspect | step 21: cursor want 1362 got 0; selStart want 1362 got 0; selEnd want 1362 got 0; state={0 0 0 199 word word word word word word word word… |
| shift-left-repeat-from-end | prepare fail | 8.2 | yes | down-one-logical-line | textLen want 1551 got 199; prepare retries; cascade suspect after down-one-logical-line |
