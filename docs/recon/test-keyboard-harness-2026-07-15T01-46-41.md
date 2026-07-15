# Keyboard harness results

Run: 2026-07-15T01:46:41+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 41 pass, 39 fail, 1 prepare fail; total 56.3s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 2.8 | no | — | — |
| home-clears-selection | pass | 0.5 | no | — | — |
| shift-right-from-home | pass | 0.6 | no | — | — |
| shift-left-from-end | pass | 0.6 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 0.7 | no | — | — |
| shift-down-after-arrow-down | pass | 0.7 | no | — | — |
| shift-up-after-arrow-down | pass | 0.6 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | pass | 0.5 | no | — | — |
| shift-down-then-up-shrinks | pass | 0.9 | no | — | — |
| shift-left-repeat-from-end | pass | 0.6 | no | — | — |
| alt-backspace-deletes-word | pass | 0.5 | no | — | — |
| ctrl-backspace-deletes-line | pass | 0.6 | no | — | — |
| shift-left-repeat-mid-doc | pass | 0.8 | no | — | — |
| cm-line-down-basic | pass | 0.5 | no | — | — |
| cm-line-down-shorter | fail | 0.5 | no | — | step 3: cursor want 2 got 0; state={0 0 0 5 one t 1 0 z-test-keyboard-harness.md} |
| cm-line-down-last-line | pass | 0.6 | no | — | — |
| cm-line-down-goal-col | fail | 0.7 | no | — | step 4: cursor want 11 got 6; selStart want 11 got 6; selEnd want 11 got 6; state={6 6 6 11 one t three 1 0 z-test-keyboard-harness.md} |
| cm-select-line-down | pass | 0.5 | no | — | — |
| cm-select-line-down-mid | fail | 0.8 | no | — | step 4: cursor want 7 got 4; selStart want 2 got 0; selEnd want 7 got 4; selLen want 5 got 4; state={4 0 4 13 one two three 1 0 z-test-keyb… |
| cm-select-down-up-doc-end | fail | 0.8 | no | — | step 5: selStart want 8 got 7; selLen want 4 got 5; state={12 7 12 13 one two three 1 0 z-test-keyboard-harness.md} |
| cm-select-up-basic | fail | 0.7 | no | — | step 3: selStart want 4 got 7; selLen want 9 got 6; state={13 7 13 13 one two three 1 0 z-test-keyboard-harness.md} |
| cm-select-up-mid | fail | 0.9 | no | — | step 5: cursor want 9 got 8; selEnd want 9 got 8; selLen want 5 got 4; state={8 4 8 13 one two three 1 0 z-test-keyboard-harness.md} |
| combo-alt-left | fail | 0.6 | no | — | step 3: cursor want 6 got 5; selStart want 6 got 5; selEnd want 6 got 5; state={5 5 5 5 hello 1 0 z-test-keyboard-harness.md} |
| combo-alt-right | fail | 0.5 | no | — | step 2: cursor want 6 got 0; selStart want 6 got 0; selEnd want 6 got 0; state={0 0 0 6  world 1 0 z-test-keyboard-harness.md} |
| combo-alt-up | pass | 0.5 | no | — | — |
| combo-alt-down | fail | 0.5 | no | — | step 2: cursor want 7 got 0; selStart want 7 got 0; selEnd want 7 got 0; state={0 0 0 7   para2 1 0 z-test-keyboard-harness.md} |
| combo-ctrl-left | pass | 0.5 | no | — | — |
| combo-ctrl-right | pass | 0.5 | no | — | — |
| combo-ctrl-up | pass | 0.5 | no | — | — |
| combo-ctrl-down | pass | 0.5 | no | — | — |
| combo-shift-alt-left | fail | 0.5 | no | — | step 3: cursor want 11 got 5; selStart want 6 got 5; selEnd want 11 got 5; selLen want 5 got 0; state={5 5 5 5 hello 1 0 z-test-keyboard-ha… |
| combo-shift-alt-right | fail | 0.5 | no | — | step 2: cursor want 6 got 0; selEnd want 6 got 0; selLen want 6 got 0; state={0 0 0 6  world 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-up | fail | 0.5 | no | — | step 3: cursor want 12 got 0; selEnd want 12 got 0; selLen want 12 got 0; state={0 0 0 7   para2 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-down | fail | 0.5 | no | — | step 2: cursor want 7 got 0; selEnd want 7 got 0; selLen want 7 got 0; state={0 0 0 7   para2 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-left | fail | 0.6 | no | — | step 3: cursor want 11 got 0; selEnd want 11 got 0; selLen want 11 got 0; state={0 0 0 0  1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-right | fail | 0.7 | no | — | step 2: cursor want 11 got 0; selEnd want 11 got 0; selLen want 11 got 0; state={0 0 0 11 hello world 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-up | fail | 0.6 | no | — | step 3: cursor want 13 got 0; selEnd want 13 got 0; selLen want 13 got 0; state={0 0 0 10  two three 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-down | fail | 0.5 | no | — | step 2: cursor want 13 got 0; selEnd want 13 got 0; selLen want 13 got 0; state={0 0 0 13 one two three 1 0 z-test-keyboard-harness.md} |
| combo-shift-home-line | pass | 0.6 | no | — | — |
| combo-shift-end-line | pass | 0.6 | no | — | — |
| combo-ctrl-home | pass | 0.5 | no | — | — |
| combo-ctrl-end | pass | 0.5 | no | — | — |
| combo-shift-ctrl-home | fail | 0.6 | no | — | step 3: cursor want 4 got 0; selEnd want 4 got 0; selLen want 4 got 0; state={0 0 0 3 def 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-end | fail | 0.6 | no | — | step 3: cursor want 7 got 0; selStart want 4 got 0; selEnd want 7 got 0; selLen want 3 got 0; state={0 0 0 3 def 1 0 z-test-keyboard-harnes… |
| bs-alt-word-mid | fail | 0.8 | no | — | step 2: cursor want 8 got 0; state={0 0 0 11 hello world 1 0 z-test-keyboard-harness.md} |
| bs-ctrl-line-start | pass | 0.5 | no | — | — |
| bs-shift-with-selection | fail | 0.7 | no | — | step 5: cursor want 0 got 4; selEnd want 0 got 4; textLen want 0 got 4; state={4 0 4 4 abcd 1 0 z-test-keyboard-harness.md} |
| bs-plain | pass | 0.5 | no | — | — |
| wrap-down-one-visual-line | fail | 0.7 | no | — | step 3: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word wor… |
| wrap-down-not-jump-paragraph | fail | 0.8 | no | — | step 3: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word wor… |
| wrap-up-from-visual-line-2 | fail | 0.6 | no | — | after down: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word… |
| wrap-shift-down-one-visual | fail | 0.7 | no | — | step 3: cursor want 10 got 20; selEnd want 10 got 20; state={20 0 20 199 word word word word word word word word word word word word word w… |
| wrap-shift-down-then-up-shrinks | fail | 1.0 | no | — | step 5: selLenMin want >= 1 got 0; state={40 40 40 199 word word word word word word word word word word word word word word word word word… |
| wrap-down-last-visual-line | pass | 0.8 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.8 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.8 | no | — | — |
| wrap-down-goal-column | fail | 0.7 | no | — | step 3: cursor want 2 got 0; state={0 0 0 177 abword word word word word word word word word word word word word word word word word word w… |
| wrap-combo-alt-left-word | fail | 0.9 | no | — | step 4: textLen want 199 got 194; state={194 194 194 194 word word word word word word word word word word word word word word word word wo… |
| wrap-combo-ctrl-bs-line | pass | 0.8 | no | — | — |
| wrap-shift-left-across-wrap | fail | 0.8 | no | — | step 3: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word wor… |
| wrap-home-on-visual-line | fail | 0.6 | no | — | step 3: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word wor… |
| wrap-end-on-visual-line | fail | 0.7 | no | — | step 3: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word wor… |
| undo-redo-len | fail | 0.9 | no | — | step 6: cursor want 5 got 0; textLen want 5 got 0; state={0 0 0 0  1 0 z-test-keyboard-harness.md} |
| undo-cursor-reposition | fail | 0.8 | no | — | step 6: cursor want 31 got 8; state={8 8 8 31 Blahfive lines in this textedit 1 0 z-test-keyboard-harness.md} |
| undo-mid-line-delete | fail | 0.9 | no | — | step 5: cursor want 4 got 1; state={1 1 1 4 adef 1 0 z-test-keyboard-harness.md} |
| redo-cleared-by-new-edit | fail | 0.8 | no | — | step 5: cursor want 3 got 0; textLen want 3 got 0; state={0 0 0 0  1 0 z-test-keyboard-harness.md} |
| undo-after-select-delete | fail | 0.8 | no | — | step 7: cursor want 6 got 0; selStart want 6 got 0; selEnd want 6 got 0; textLen want 6 got 0; state={0 0 0 0  1 0 z-test-keyboard-harness.… |
| gap-up-at-doc-start | pass | 0.5 | no | — | — |
| gap-plain-left-no-cursor-move | pass | 0.5 | no | — | — |
| gap-plain-right-no-cursor-move | pass | 0.5 | no | — | — |
| gap-collapse-selection-left | pass | 0.8 | no | — | — |
| gap-collapse-selection-right | pass | 0.7 | no | — | — |
| gap-delete-forward | fail | 0.6 | no | — | step 3: cursor want 2 got 0; state={0 0 0 4 abcd 1 0 z-test-keyboard-harness.md} |
| gap-delete-with-selection | pass | 0.6 | no | — | — |
| gap-select-all | pass | 0.4 | no | — | — |
| gap-enter-new-line | pass | 0.5 | no | — | — |
| gap-type-replaces-selection | pass | 0.6 | no | — | — |
| gap-redo-shift-ctrl-z | fail | 0.8 | no | — | step 5: cursor want 3 got 0; textLen want 3 got 0; state={0 0 0 0  1 0 z-test-keyboard-harness.md} |
| gap-undo-chain | fail | 0.7 | no | suspect | step 5: cursor want 2 got 1; textLen want 2 got 1; state={1 1 1 1 a 1 0 z-test-keyboard-harness.md}; may have poisoned next scenario |
| gap-unicode-alt-backspace | prepare fail | 2.6 | yes | gap-undo-chain | textLen want 13 got 11; prepare retries; cascade suspect after gap-undo-chain |
