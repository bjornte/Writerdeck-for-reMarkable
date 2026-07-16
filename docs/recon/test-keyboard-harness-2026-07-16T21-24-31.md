# Keyboard harness results

Run: 2026-07-16T21:24:31+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 35 pass, 1 fail, 0 prepare fail; total 64.5s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.4 | no | — | — |
| home-clears-selection | pass | 1.0 | no | — | — |
| shift-right-from-home | pass | 3.4 | no | — | — |
| shift-left-from-end | pass | 3.3 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.3 | no | — | — |
| shift-down-after-arrow-down | fail | 2.0 | no | — | step 14: cursor want 1419 got 1362; selEnd want 1419 got 1362; selLen want 57 got 0; state={1362 1362 1362 1551 Writerdeck harness dummy — … |
| shift-up-after-arrow-down | pass | 3.2 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | pass | 3.4 | no | — | — |
| shift-left-repeat-from-end | pass | 3.3 | no | — | — |
| alt-backspace-deletes-word | pass | 1.6 | no | — | — |
| ctrl-backspace-deletes-line | pass | 1.7 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.5 | no | — | — |
| cm-line-down-basic | pass | 3.2 | no | — | — |
| cm-line-down-last-line | pass | 1.1 | no | — | — |
| combo-alt-left | pass | 3.4 | no | — | — |
| combo-alt-right | pass | 3.5 | no | — | — |
| combo-ctrl-home | pass | 1.2 | no | — | — |
| combo-ctrl-end | pass | 1.0 | no | — | — |
| bs-plain | pass | 1.1 | no | — | — |
| wrap-down-one-visual-line | pass | 1.3 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.6 | no | — | — |
| undo-redo-len | pass | 1.0 | no | — | — |
| gap-up-at-doc-start | pass | 1.0 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.3 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.4 | no | — | — |
| gap-collapse-selection-left | pass | 0.9 | no | — | — |
| gap-collapse-selection-right | pass | 0.9 | no | — | — |
| gap-delete-forward | pass | 1.0 | no | — | — |
| gap-delete-with-selection | pass | 0.9 | no | — | — |
| gap-select-all | pass | 0.6 | no | — | — |
| gap-enter-new-line | pass | 0.6 | no | — | — |
| gap-type-replaces-selection | pass | 0.8 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 1.0 | no | — | — |
| gap-undo-chain | pass | 0.8 | no | — | — |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
