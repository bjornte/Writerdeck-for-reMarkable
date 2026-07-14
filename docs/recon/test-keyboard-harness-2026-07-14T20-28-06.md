# Keyboard harness results

Run: 2026-07-14T20:28:06+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 27 pass, 35 fail, 0 prepare fail; total 38.7s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.7 | no | — | — |
| home-clears-selection | pass | 0.6 | no | — | — |
| shift-right-from-home | pass | 2.0 | no | — | — |
| shift-left-from-end | pass | 0.6 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 0.6 | no | — | — |
| shift-down-after-arrow-down | pass | 0.7 | no | — | — |
| shift-up-after-arrow-down | pass | 0.8 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.5 | no | — | — |
| down-one-logical-line | pass | 0.5 | no | — | — |
| shift-down-then-up-shrinks | pass | 0.8 | no | — | — |
| shift-left-repeat-from-end | pass | 0.6 | no | — | — |
| alt-backspace-deletes-word | pass | 0.7 | no | — | — |
| ctrl-backspace-deletes-line | pass | 0.7 | no | — | — |
| shift-left-repeat-mid-doc | pass | 0.7 | no | — | — |
| cm-line-down-basic | pass | 0.5 | no | — | — |
| cm-line-down-shorter | fail | 0.8 | no | — | step 5: cursor want 4 got 5; selStart want 4 got 5; selEnd want 4 got 5; state={5 5 5 5 1 0 z-test-keyboard-harness.md} |
| cm-line-down-last-line | pass | 0.6 | no | — | — |
| cm-line-down-goal-col | fail | 0.7 | no | — | step 4: cursor want 11 got 7; selStart want 11 got 7; selEnd want 11 got 7; state={7 7 7 11 1 0 z-test-keyboard-harness.md} |
| cm-select-line-down | pass | 0.5 | no | — | — |
| cm-select-line-down-mid | fail | 0.8 | no | — | step 4: cursor want 7 got 6; selEnd want 7 got 6; selLen want 5 got 4; state={6 2 6 13 1 0 z-test-keyboard-harness.md} |
| cm-select-down-up-doc-end | fail | 0.7 | no | — | step 2: cursor want 13 got 3; state={3 3 3 13 1 0 z-test-keyboard-harness.md} |
| cm-select-up-basic | fail | 0.8 | no | — | step 3: cursor want 13 got 3; selStart want 4 got 0; selEnd want 13 got 3; selLen want 9 got 3; state={3 0 3 13 1 0 z-test-keyboard-harness… |
| cm-select-up-mid | fail | 0.8 | no | — | step 5: selStart want 4 got 5; selLen want 5 got 4; state={9 5 9 13 1 0 z-test-keyboard-harness.md} |
| combo-alt-left | fail | 0.5 | no | — | step 3: cursor want 6 got 5; selStart want 6 got 5; selEnd want 6 got 5; state={5 5 5 5 1 0 z-test-keyboard-harness.md} |
| combo-alt-right | fail | 0.4 | no | — | step 2: cursor want 6 got 0; selStart want 6 got 0; selEnd want 6 got 0; state={0 0 0 11 1 0 z-test-keyboard-harness.md} |
| combo-alt-up | pass | 0.4 | no | — | — |
| combo-alt-down | fail | 0.5 | no | — | step 2: cursor want 7 got 0; selStart want 7 got 0; selEnd want 7 got 0; state={0 0 0 12 1 0 z-test-keyboard-harness.md} |
| combo-ctrl-left | pass | 0.7 | no | — | — |
| combo-ctrl-right | fail | 0.4 | no | — | step 2: cursor want 11 got 0; selStart want 11 got 0; selEnd want 11 got 0; state={0 0 0 11 1 0 z-test-keyboard-harness.md} |
| combo-ctrl-up | pass | 0.5 | no | — | — |
| combo-ctrl-down | fail | 0.4 | no | — | step 2: cursor want 13 got 0; selStart want 13 got 0; selEnd want 13 got 0; state={0 0 0 13 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-left | fail | 0.5 | no | — | step 3: cursor want 11 got 5; selStart want 6 got 5; selEnd want 11 got 5; selLen want 5 got 0; state={5 5 5 5 1 0 z-test-keyboard-harness.… |
| combo-shift-alt-right | fail | 0.4 | no | — | step 2: cursor want 6 got 0; selEnd want 6 got 0; selLen want 6 got 0; state={0 0 0 11 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-up | fail | 0.6 | no | — | step 3: cursor want 12 got 0; selEnd want 12 got 0; selLen want 12 got 0; state={0 0 0 7 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-down | fail | 0.4 | no | — | step 2: cursor want 7 got 0; selEnd want 7 got 0; selLen want 7 got 0; state={0 0 0 12 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-left | fail | 0.5 | no | — | step 3: cursor want 11 got 0; selEnd want 11 got 0; selLen want 11 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-right | fail | 0.4 | no | — | step 2: cursor want 11 got 0; selEnd want 11 got 0; selLen want 11 got 0; state={0 0 0 11 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-up | fail | 0.5 | no | — | step 3: cursor want 13 got 0; selEnd want 13 got 0; selLen want 13 got 0; state={0 0 0 10 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-down | fail | 0.4 | no | — | step 2: cursor want 13 got 0; selEnd want 13 got 0; selLen want 13 got 0; state={0 0 0 13 1 0 z-test-keyboard-harness.md} |
| combo-shift-home-line | fail | 0.6 | no | — | step 4: selStart want 0 got 4; selLen want 7 got 3; state={7 4 7 7 1 0 z-test-keyboard-harness.md} |
| combo-shift-end-line | fail | 0.6 | no | — | step 4: cursor want 7 got 0; selStart want 4 got 0; selEnd want 7 got 0; selLen want 3 got 0; state={0 0 0 0 1 1 }; editor in lobby after f… |
| combo-ctrl-home | pass | 0.5 | no | — | — |
| combo-ctrl-end | fail | 0.4 | no | — | step 2: cursor want 7 got 0; selStart want 7 got 0; selEnd want 7 got 0; state={0 0 0 7 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-home | fail | 0.5 | no | — | step 3: cursor want 4 got 0; selEnd want 4 got 0; selLen want 4 got 0; state={0 0 0 3 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-end | fail | 0.5 | no | — | step 3: cursor want 7 got 0; selStart want 4 got 0; selEnd want 7 got 0; selLen want 3 got 0; state={0 0 0 3 1 0 z-test-keyboard-harness.md} |
| bs-alt-word-mid | fail | 0.8 | no | — | step 4: textLen want 5 got 8; state={5 5 5 8 1 0 z-test-keyboard-harness.md} |
| bs-ctrl-line-start | fail | 0.6 | no | — | step 4: cursor want 6 got 0; textLen want 6 got 5; state={0 0 0 5 1 0 z-test-keyboard-harness.md} |
| bs-shift-with-selection | pass | 0.6 | no | — | — |
| bs-plain | pass | 0.6 | no | — | — |
| wrap-down-one-visual-line | fail | 0.5 | no | — | step 3: cursorMax want <= 198 got 199; state={199 199 199 199 1 0 z-test-keyboard-harness.md} |
| wrap-down-not-jump-paragraph | fail | 0.5 | no | — | step 3: cursorMax want <= 198 got 199; state={199 199 199 199 1 0 z-test-keyboard-harness.md} |
| wrap-up-from-visual-line-2 | fail | 0.5 | no | — | after down: cursorMax want <= 198 got 199; state={199 199 199 199 1 0 z-test-keyboard-harness.md} |
| wrap-shift-down-one-visual | pass | 0.5 | no | — | — |
| wrap-shift-down-then-up-shrinks | pass | 0.8 | no | — | — |
| wrap-down-last-visual-line | pass | 0.6 | no | — | — |
| wrap-shift-down-last-to-eof | fail | 0.5 | no | — | step 3: selLenMin want >= 1 got 0; state={199 199 199 199 1 0 z-test-keyboard-harness.md} |
| wrap-mixed-newline-and-wrap | pass | 0.5 | no | — | — |
| undo-redo-len | fail | 0.9 | no | — | step 6: cursor want 5 got 0; textLen want 5 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| undo-cursor-reposition | fail | 0.8 | no | — | step 6: cursor want 31 got 8; state={8 8 8 31 1 0 z-test-keyboard-harness.md} |
| undo-mid-line-delete | fail | 0.9 | no | — | step 5: cursor want 4 got 1; state={1 1 1 4 1 0 z-test-keyboard-harness.md} |
| redo-cleared-by-new-edit | fail | 0.8 | no | — | step 5: cursor want 3 got 0; textLen want 3 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| undo-after-select-delete | fail | 0.8 | no | — | step 7: cursor want 6 got 0; selStart want 6 got 0; selEnd want 6 got 0; textLen want 6 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.m… |
