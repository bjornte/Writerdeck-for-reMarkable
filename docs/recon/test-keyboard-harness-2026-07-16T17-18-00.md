# Keyboard harness results

Run: 2026-07-16T17:18:00+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 34 pass, 2 fail, 0 prepare fail; total 67.2s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 1.9 | no | — | — |
| home-clears-selection | pass | 0.8 | no | — | — |
| shift-right-from-home | pass | 3.3 | no | — | — |
| shift-left-from-end | pass | 3.2 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.3 | no | — | — |
| shift-down-after-arrow-down | fail | 1.5 | no | — | step 6: cursor want 1457 got 1362; selEnd want 1457 got 1362; selLen want 95 got 0; state={1362 1362 1362 1551 Writerdeck harness dummy — i… |
| shift-up-after-arrow-down | pass | 3.3 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | pass | 3.4 | no | — | — |
| shift-left-repeat-from-end | pass | 3.7 | no | — | — |
| alt-backspace-deletes-word | pass | 1.7 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.6 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.4 | no | — | — |
| cm-line-down-basic | pass | 3.3 | no | — | — |
| cm-line-down-last-line | pass | 1.2 | no | — | — |
| combo-alt-left | pass | 3.8 | no | — | — |
| combo-alt-right | pass | 3.6 | no | — | — |
| combo-ctrl-home | pass | 1.2 | no | — | — |
| combo-ctrl-end | pass | 1.0 | no | — | — |
| bs-plain | pass | 1.0 | no | — | — |
| wrap-down-one-visual-line | pass | 1.1 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.5 | no | — | — |
| undo-redo-len | pass | 1.6 | no | — | — |
| gap-up-at-doc-start | pass | 1.3 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.4 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.3 | no | — | — |
| gap-collapse-selection-left | pass | 1.1 | no | — | — |
| gap-collapse-selection-right | pass | 0.8 | no | — | — |
| gap-delete-forward | pass | 1.2 | no | — | — |
| gap-delete-with-selection | pass | 0.9 | no | — | — |
| gap-select-all | pass | 0.7 | no | — | — |
| gap-enter-new-line | pass | 0.6 | no | — | — |
| gap-type-replaces-selection | pass | 0.9 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 1.1 | no | — | — |
| gap-undo-chain | fail | 0.6 | no | — | step 3: cursor want 5 got 6; textLen want 5 got 6; state={6 6 6 6 abc æø 1 0 z-test-keyboard-harness.md 0} |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
