# Keyboard harness results

Run: 2026-07-15T02:07:10+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 52 pass, 28 fail, 1 prepare fail; total 53.3s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.2 | no | — | — |
| home-clears-selection | pass | 0.5 | no | — | — |
| shift-right-from-home | pass | 0.6 | no | — | — |
| shift-left-from-end | pass | 0.6 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 0.6 | no | — | — |
| shift-down-after-arrow-down | pass | 0.6 | no | — | — |
| shift-up-after-arrow-down | pass | 0.8 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.5 | no | — | — |
| down-one-logical-line | pass | 0.5 | no | — | — |
| shift-down-then-up-shrinks | fail | 0.8 | no | — | step 7: cursor want 17 got 12; selEnd want 17 got 12; selLen want 5 got 0; state={12 12 12 23 line1 line2 line3 line4 1 0 z-test-keyboard-h… |
| shift-left-repeat-from-end | pass | 0.6 | no | — | — |
| alt-backspace-deletes-word | pass | 0.7 | no | — | — |
| ctrl-backspace-deletes-line | pass | 0.7 | no | — | — |
| shift-left-repeat-mid-doc | pass | 0.7 | no | — | — |
| cm-line-down-basic | pass | 0.5 | no | — | — |
| cm-line-down-shorter | fail | 0.8 | no | — | step 5: cursor want 4 got 5; selStart want 4 got 5; selEnd want 4 got 5; state={5 5 5 5 one t 1 0 z-test-keyboard-harness.md} |
| cm-line-down-last-line | pass | 0.6 | no | — | — |
| cm-line-down-goal-col | pass | 0.7 | no | — | — |
| cm-select-line-down | pass | 0.6 | no | — | — |
| cm-select-line-down-mid | pass | 0.8 | no | — | — |
| cm-select-down-up-doc-end | fail | 0.8 | no | — | step 5: selStart want 8 got 7; selLen want 4 got 5; state={12 7 12 13 one two three 1 0 z-test-keyboard-harness.md} |
| cm-select-up-basic | pass | 0.7 | no | — | — |
| cm-select-up-mid | pass | 0.9 | no | — | — |
| combo-alt-left | pass | 0.6 | no | — | — |
| combo-alt-right | pass | 0.5 | no | — | — |
| combo-alt-up | pass | 0.5 | no | — | — |
| combo-alt-down | pass | 0.5 | no | — | — |
| combo-ctrl-left | pass | 0.5 | no | — | — |
| combo-ctrl-right | pass | 0.5 | no | — | — |
| combo-ctrl-up | pass | 0.5 | no | — | — |
| combo-ctrl-down | pass | 0.5 | no | — | — |
| combo-shift-alt-left | pass | 0.5 | no | — | — |
| combo-shift-alt-right | pass | 0.5 | no | — | — |
| combo-shift-alt-up | fail | 0.6 | no | — | step 3: cursor want 12 got 5; selEnd want 12 got 5; selLen want 12 got 5; state={5 0 5 12 para1  para2 1 0 z-test-keyboard-harness.md} |
| combo-shift-alt-down | pass | 0.5 | no | — | — |
| combo-shift-ctrl-left | pass | 0.6 | no | — | — |
| combo-shift-ctrl-right | pass | 0.7 | no | — | — |
| combo-shift-ctrl-up | fail | 0.6 | no | — | step 3: cursor want 13 got 3; selEnd want 13 got 3; selLen want 13 got 3; state={3 0 3 13 one two three 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-down | pass | 0.5 | no | — | — |
| combo-shift-home-line | pass | 0.6 | no | — | — |
| combo-shift-end-line | pass | 0.6 | no | — | — |
| combo-ctrl-home | pass | 0.5 | no | — | — |
| combo-ctrl-end | pass | 0.5 | no | — | — |
| combo-shift-ctrl-home | fail | 0.5 | no | — | step 3: cursor want 4 got 0; selEnd want 4 got 0; selLen want 4 got 0; state={0 0 0 3 def 1 0 z-test-keyboard-harness.md} |
| combo-shift-ctrl-end | fail | 0.6 | no | — | step 3: cursor want 7 got 0; selStart want 4 got 0; selEnd want 7 got 0; selLen want 3 got 0; state={0 0 0 3 def 1 0 z-test-keyboard-harnes… |
| bs-alt-word-mid | pass | 0.9 | no | — | — |
| bs-ctrl-line-start | pass | 0.5 | no | — | — |
| bs-shift-with-selection | fail | 0.7 | no | — | step 5: cursor want 0 got 4; selEnd want 0 got 4; textLen want 0 got 4; state={4 0 4 4 abcd 1 0 z-test-keyboard-harness.md} |
| bs-plain | fail | 0.5 | no | — | step 3: cursor want 2 got 4; textLen want 2 got 4; text want "ab" got "abcd"; state={4 3 4 4 abcd 1 0 z-test-keyboard-harness.md} |
| wrap-down-one-visual-line | fail | 0.7 | no | — | step 3: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word wor… |
| wrap-down-not-jump-paragraph | fail | 0.8 | no | — | step 3: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word wor… |
| wrap-up-from-visual-line-2 | fail | 0.6 | no | — | after down: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word… |
| wrap-shift-down-one-visual | fail | 0.7 | no | — | step 3: cursor want 10 got 20; selEnd want 10 got 20; state={20 0 20 199 word word word word word word word word word word word word word w… |
| wrap-shift-down-then-up-shrinks | pass | 1.0 | no | — | — |
| wrap-down-last-visual-line | pass | 0.8 | no | — | — |
| wrap-shift-down-last-to-eof | pass | 0.8 | no | — | — |
| wrap-mixed-newline-and-wrap | pass | 0.8 | no | — | — |
| wrap-down-goal-column | fail | 0.8 | no | — | step 5: cursor want 10 got 24; state={24 24 24 177 abword word word word word word word word word word word word word word word word word w… |
| wrap-combo-alt-left-word | pass | 0.9 | no | — | — |
| wrap-combo-ctrl-bs-line | pass | 0.8 | no | — | — |
| wrap-shift-left-across-wrap | fail | 0.7 | no | — | step 3: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word wor… |
| wrap-home-on-visual-line | fail | 0.7 | no | — | step 3: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word wor… |
| wrap-end-on-visual-line | fail | 0.8 | no | — | step 3: cursor want 10 got 20; state={20 20 20 199 word word word word word word word word word word word word word word word word word wor… |
| undo-redo-len | fail | 0.6 | no | — | step 4: cursor want 0 got 5; textLen want 0 got 5; state={5 4 5 5 abc d 1 0 z-test-keyboard-harness.md} |
| undo-cursor-reposition | fail | 0.8 | no | — | step 6: cursor want 31 got 8; state={8 8 8 31 Blahfive lines in this textedit 1 0 z-test-keyboard-harness.md} |
| undo-mid-line-delete | fail | 0.9 | no | — | step 5: textLen want 4 got 7; state={4 3 4 7 abc def 1 0 z-test-keyboard-harness.md} |
| redo-cleared-by-new-edit | fail | 0.7 | no | — | step 3: textLen want 0 got 3; state={3 2 3 3 abc 1 0 z-test-keyboard-harness.md} |
| undo-after-select-delete | fail | 0.7 | no | — | step 5: cursor want 0 got 6; textLen want 0 got 6; state={6 0 6 6 abcdef 1 0 z-test-keyboard-harness.md} |
| gap-up-at-doc-start | pass | 0.5 | no | — | — |
| gap-plain-left-no-cursor-move | fail | 0.5 | no | — | step 4: cursor want 11 got 10; selStart want 11 got 10; selEnd want 11 got 10; state={10 10 10 11 hello world 1 0 z-test-keyboard-harness.m… |
| gap-plain-right-no-cursor-move | fail | 0.5 | no | — | step 4: cursor want 0 got 1; selStart want 0 got 1; selEnd want 0 got 1; state={1 1 1 5 hello 1 0 z-test-keyboard-harness.md} |
| gap-collapse-selection-left | pass | 0.7 | no | — | — |
| gap-collapse-selection-right | pass | 0.7 | no | — | — |
| gap-delete-forward | fail | 0.8 | no | — | step 5: cursor want 2 got 3; textLen want 3 got 4; text want "abd" got "abcd"; state={3 2 3 4 abcd 1 0 z-test-keyboard-harness.md} |
| gap-delete-with-selection | fail | 0.7 | no | — | step 5: cursor want 0 got 6; selEnd want 0 got 6; textLen want 0 got 6; text want "" got "abcdef"; state={6 0 6 6 abcdef 1 0 z-test-keyboar… |
| gap-select-all | pass | 0.4 | no | — | — |
| gap-enter-new-line | pass | 0.5 | no | — | — |
| gap-type-replaces-selection | pass | 0.7 | no | — | — |
| gap-redo-shift-ctrl-z | fail | 0.6 | no | — | step 3: cursor want 0 got 3; textLen want 0 got 3; state={3 2 3 3 abc 1 0 z-test-keyboard-harness.md} |
| gap-undo-chain | fail | 0.5 | no | suspect | step 3: cursor want 1 got 3; textLen want 1 got 3; state={3 2 3 3 abc 1 0 z-test-keyboard-harness.md}; may have poisoned next scenario |
| gap-unicode-alt-backspace | prepare fail | 2.6 | yes | gap-undo-chain | textLen want 13 got 11; prepare retries; cascade suspect after gap-undo-chain |
