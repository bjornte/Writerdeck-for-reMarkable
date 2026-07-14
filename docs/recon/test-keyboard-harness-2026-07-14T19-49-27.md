# Keyboard harness results

Run: 2026-07-14T19:49:27+02:00

Target: `192.168.1.8:8000`

Mode: soft-reset (single launch)

Timing: fast pauses

Suite setup: 0.6s (one cold start, included in first scenario time)

Summary: 13 pass, 5 fail, 0 prepare fail; total 29.5s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 4.2 | no | — | — |
| home-clears-selection | pass | 0.5 | no | — | — |
| shift-right-from-home | pass | 2.0 | no | — | — |
| shift-left-from-end | pass | 4.6 | yes | — | auto hard-reset before run |
| shift-right-after-home-no-stale-anchor | pass | 4.6 | yes | — | auto hard-reset before run |
| shift-down-after-arrow-down | pass | 0.8 | no | — | — |
| shift-up-after-arrow-down | pass | 4.8 | yes | — | auto hard-reset before run |
| ctrl-shift-left-select-line | pass | 0.5 | no | — | — |
| down-one-logical-line | pass | 0.5 | no | — | — |
| shift-down-then-up-shrinks | pass | 0.9 | no | — | — |
| shift-left-repeat-from-end | pass | 0.6 | no | — | — |
| alt-backspace-deletes-word | pass | 0.6 | no | — | — |
| ctrl-backspace-deletes-line | pass | 0.6 | no | — | — |
| undo-redo-len | fail | 0.9 | no | — | step 6: cursor want 5 got 0; textLen want 5 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| undo-cursor-reposition | fail | 0.8 | no | — | step 6: cursor want 31 got 8; state={8 8 8 31 1 0 z-test-keyboard-harness.md} |
| undo-mid-line-delete | fail | 0.9 | no | — | step 5: cursor want 4 got 1; state={1 1 1 4 1 0 z-test-keyboard-harness.md} |
| redo-cleared-by-new-edit | fail | 0.8 | no | — | step 5: cursor want 3 got 0; textLen want 3 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| undo-after-select-delete | fail | 0.8 | no | — | step 7: cursor want 6 got 0; selStart want 6 got 0; selEnd want 6 got 0; textLen want 6 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
