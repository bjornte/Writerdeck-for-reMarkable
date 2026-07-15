# Keyboard harness results

Run: 2026-07-15T10:13:45+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 30 pass, 4 fail, 0 prepare fail; total 23.3s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 3.2 | no | — | — |
| home-clears-selection | pass | 0.7 | no | — | — |
| shift-right-from-home | pass | 0.5 | no | — | — |
| shift-left-from-end | pass | 0.6 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 0.6 | no | — | — |
| shift-down-after-arrow-down | pass | 0.6 | no | — | — |
| shift-up-after-arrow-down | pass | 0.8 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.5 | no | — | — |
| down-one-logical-line | pass | 0.5 | no | — | — |
| shift-left-repeat-from-end | pass | 0.6 | no | — | — |
| alt-backspace-deletes-word | pass | 0.5 | no | — | — |
| ctrl-backspace-deletes-line | pass | 0.6 | no | — | — |
| shift-left-repeat-mid-doc | pass | 0.7 | no | — | — |
| cm-line-down-basic | pass | 0.5 | no | — | — |
| cm-line-down-last-line | pass | 0.8 | no | — | — |
| combo-alt-left | pass | 0.5 | no | — | — |
| combo-alt-right | pass | 0.5 | no | — | — |
| combo-ctrl-home | pass | 0.5 | no | — | — |
| combo-ctrl-end | pass | 0.5 | no | — | — |
| bs-plain | pass | 0.5 | no | — | — |
| wrap-down-one-visual-line | pass | 0.7 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 0.8 | no | — | — |
| undo-redo-len | fail | 0.9 | no | — | step 6: cursor want 5 got 0; textLen want 5 got 0; state={0 0 0 0  1 0 z-test-keyboard-harness.md} |
| gap-up-at-doc-start | pass | 0.6 | no | — | — |
| gap-collapse-selection-left | pass | 0.8 | no | — | — |
| gap-collapse-selection-right | pass | 0.7 | no | — | — |
| gap-delete-forward | fail | 0.5 | no | — | step 3: cursor want 2 got 0; state={0 0 0 4 abcd 1 0 z-test-keyboard-harness.md} |
| gap-delete-with-selection | pass | 0.6 | no | — | — |
| gap-select-all | pass | 0.4 | no | — | — |
| gap-enter-new-line | pass | 0.5 | no | — | — |
| gap-type-replaces-selection | pass | 0.6 | no | — | — |
| gap-redo-shift-ctrl-z | fail | 0.8 | no | — | step 5: cursor want 3 got 0; textLen want 3 got 0; state={0 0 0 0  1 0 z-test-keyboard-harness.md} |
| gap-undo-chain | fail | 0.7 | no | — | step 5: cursor want 2 got 1; textLen want 2 got 1; state={1 1 1 1 a 1 0 z-test-keyboard-harness.md} |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
