# Keyboard harness results

Run: 2026-07-14T23:16:21+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 23 pass, 39 fail, 0 prepare fail; total 35.3s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.3 | no | — | — |
| home-clears-selection | pass | 0.5 | no | — | — |
| shift-right-from-home | fail | 0.6 | no | — | step 4: cursor want 3 got 0; selEnd want 3 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| shift-left-from-end | pass | 0.6 | no | — | — |
| shift-right-after-home-no-stale-anchor | fail | 0.6 | no | — | step 4: cursor want 1 got 0; selEnd want 1 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| shift-down-after-arrow-down | fail | 0.6 | no | — | cursor on line3: cursor want 12 got 11; selStart want 12 got 11; selEnd want 12 got 11; textLen want 23 got 18; state={11 11 11 18 1 0 z-te… |
| shift-up-after-arrow-down | fail | 0.5 | no | — | cursor on line3: cursor want 12 got 11; selStart want 12 got 11; selEnd want 12 got 11; state={11 11 11 18 1 0 z-test-keyboard-harness.md} |
| ctrl-shift-left-select-line | pass | 0.5 | no | — | — |
| down-one-logical-line | fail | 0.4 | no | — | step 2: textLen want 5 got 3; state={0 0 0 3 1 0 z-test-keyboard-harness.md} |
| shift-down-then-up-shrinks | fail | 0.5 | no | — | cursor line3: cursor want 12 got 11; state={11 11 11 18 1 0 z-test-keyboard-harness.md} |
| shift-left-repeat-from-end | pass | 0.6 | no | — | — |
| alt-backspace-deletes-word | pass | 0.5 | no | — | — |
| ctrl-backspace-deletes-line | pass | 0.6 | no | — | — |
| shift-left-repeat-mid-doc | pass | 0.7 | no | — | — |
| cm-line-down-basic | fail | 0.6 | no | — | step 3: cursor want 4 got 3; selStart want 4 got 3; selEnd want 4 got 3; state={3 3 3 4 1 0 z-test-keyboard-harness.md} |
| cm-line-down-shorter | fail | 0.7 | no | — | step 5: cursor want 4 got 2; selStart want 4 got 2; selEnd want 4 got 2; state={2 2 2 2 1 0 z-test-keyboard-harness.md} |
| cm-line-down-last-line | fail | 0.7 | no | — | step 4: cursor want 3 got 0; selStart want 3 got 0; selEnd want 3 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| cm-line-down-goal-col | fail | 0.7 | no | — | step 4: cursor want 11 got 8; selStart want 11 got 8; selEnd want 11 got 8; state={8 8 8 8 1 0 z-test-keyboard-harness.md} |
| cm-select-line-down | fail | 0.5 | no | — | step 3: cursor want 4 got 3; selEnd want 4 got 3; selLen want 4 got 3; state={3 0 3 10 1 0 z-test-keyboard-harness.md} |
| cm-select-line-down-mid | fail | 0.7 | no | — | step 4: cursor want 7 got 8; selEnd want 7 got 8; selLen want 5 got 6; state={8 2 8 10 1 0 z-test-keyboard-harness.md} |
| cm-select-down-up-doc-end | fail | 0.4 | no | — | step 2: cursor want 13 got 3; state={3 3 3 13 1 0 z-test-keyboard-harness.md} |
| cm-select-up-basic | fail | 0.7 | no | — | step 3: cursor want 13 got 3; selStart want 4 got 0; selEnd want 13 got 3; selLen want 9 got 3; state={3 0 3 13 1 0 z-test-keyboard-harness… |
| cm-select-up-mid | fail | 0.8 | no | — | step 5: cursor want 9 got 8; selStart want 4 got 3; selEnd want 9 got 8; state={8 3 8 10 1 0 z-test-keyboard-harness.md} |
| combo-alt-left | fail | 0.5 | no | — | step 3: cursor want 6 got 5; selStart want 6 got 5; selEnd want 6 got 5; state={5 5 5 5 1 0 z-test-keyboard-harness.md} |
| combo-alt-right | fail | 0.4 | no | — | step 2: cursor want 6 got 5; selStart want 6 got 5; selEnd want 6 got 5; state={5 5 5 5 1 0 z-test-keyboard-harness.md} |
| combo-alt-up | pass | 0.5 | no | — | — |
| combo-alt-down | fail | 0.4 | no | — | step 2: cursor want 7 got 0; selStart want 7 got 0; selEnd want 7 got 0; state={0 0 0 7 1 0 z-test-keyboard-harness.md} |
| combo-ctrl-left | pass | 0.5 | no | — | — |
| combo-ctrl-right | fail | 0.4 | no | — | step 2: cursor want 11 got 0; selStart want 11 got 0; selEnd want 11 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| combo-ctrl-up | pass | 0.6 | no | — | — |
| combo-ctrl-down | fail | 0.4 | no | — | step 2: cursor want 13 got 0; selStart want 13 got 0; selEnd want 13 got 0; state={0 0 0 10 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-left | fail | 0.7 | no | — | step 3: cursor want 11 got 5; selStart want 6 got 5; selEnd want 11 got 5; selLen want 5 got 0; state={5 5 5 5 1 0 z-test-keyboard-harness.… |
| combo-shift-alt-right | fail | 0.4 | no | — | step 2: cursor want 6 got 5; selStart want 0 got 5; selEnd want 6 got 5; selLen want 6 got 0; state={5 5 5 5 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-up | fail | 0.7 | no | — | step 3: cursor want 12 got 0; selEnd want 12 got 0; selLen want 12 got 0; state={0 0 0 7 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-down | fail | 0.4 | no | — | step 2: cursor want 7 got 0; selEnd want 7 got 0; selLen want 7 got 0; state={0 0 0 7 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-left | fail | 0.6 | no | — | step 3: cursor want 11 got 0; selEnd want 11 got 0; selLen want 11 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-right | fail | 0.5 | no | — | step 2: cursor want 11 got 0; selEnd want 11 got 0; selLen want 11 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-up | fail | 0.5 | no | — | step 3: cursor want 13 got 0; selEnd want 13 got 0; selLen want 13 got 0; state={0 0 0 10 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-down | fail | 0.5 | no | — | step 2: cursor want 13 got 0; selEnd want 13 got 0; selLen want 13 got 0; state={0 0 0 10 1 0 z-test-keyboard-harness.md} |
| combo-shift-home-line | pass | 0.6 | no | — | — |
| combo-shift-end-line | pass | 0.7 | no | — | — |
| combo-ctrl-home | pass | 0.5 | no | — | — |
| combo-ctrl-end | fail | 0.4 | no | — | step 2: cursor want 7 got 0; selStart want 7 got 0; selEnd want 7 got 0; state={0 0 0 4 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-home | fail | 0.5 | no | — | step 3: cursor want 4 got 0; selEnd want 4 got 0; selLen want 4 got 0; state={0 0 0 3 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-end | fail | 0.7 | no | — | step 3: cursor want 7 got 0; selStart want 4 got 0; selEnd want 7 got 0; selLen want 3 got 0; state={0 0 0 3 1 0 z-test-keyboard-harness.md} |
| bs-alt-word-mid | pass | 0.8 | no | — | — |
| bs-ctrl-line-start | pass | 0.6 | no | — | — |
| bs-shift-with-selection | pass | 0.7 | no | — | — |
| bs-plain | pass | 0.6 | no | — | — |
| wrap-down-one-visual-line | fail | 0.5 | no | — | step 3: cursorMin want >= 1 got 0; textLen want 199 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| wrap-down-not-jump-paragraph | fail | 0.5 | no | — | step 3: cursorMin want >= 1 got 0; textLen want 199 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| wrap-up-from-visual-line-2 | fail | 0.5 | no | — | after down: cursorMin want >= 1 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| wrap-shift-down-one-visual | fail | 0.5 | no | — | step 3: cursorMin want >= 1 got 0; selLenMin want >= 1 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| wrap-shift-down-then-up-shrinks | fail | 0.8 | no | — | step 5: selLenMin want >= 1 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| wrap-down-last-visual-line | pass | 0.5 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.5 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.6 | no | — | — |
| undo-redo-len | fail | 1.0 | no | — | step 8: cursor want 5 got 0; textLen want 5 got 0; state={0 0 0 0 1 0 z-test-keyboard-harness.md} |
| undo-cursor-reposition | fail | 0.4 | no | — | step 2: textLen want 27 got 23; state={0 0 0 23 1 0 z-test-keyboard-harness.md} |
| undo-mid-line-delete | fail | 0.5 | no | — | step 3: cursor want 4 got 3; textLen want 7 got 4; state={3 3 3 4 1 0 z-test-keyboard-harness.md} |
| redo-cleared-by-new-edit | pass | 1.0 | no | — | — |
| undo-after-select-delete | pass | 0.7 | no | — | — |
