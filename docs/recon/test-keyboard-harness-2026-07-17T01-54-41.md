# Keyboard harness results

Run: 2026-07-17T01:54:45+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 37 pass, 1 fail, 0 prepare fail; total 75.8s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.3 | no | — | — |
| home-clears-selection | pass | 1.0 | no | — | — |
| shift-right-from-home | pass | 3.5 | no | — | — |
| shift-left-from-end | pass | 3.3 | no | — | — |
| shift-right-after-home-no-stale-anchor | pass | 3.3 | no | — | — |
| shift-down-after-arrow-down | pass | 3.4 | no | — | — |
| shift-up-after-arrow-down | pass | 3.4 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.6 | no | — | — |
| down-one-logical-line | pass | 3.9 | no | — | — |
| shift-left-repeat-from-end | pass | 3.4 | no | — | — |
| alt-backspace-deletes-word | pass | 1.6 | no | — | — |
| ctrl-backspace-deletes-line | pass | 2.2 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.4 | no | — | — |
| cm-line-down-basic | pass | 3.4 | no | — | — |
| cm-line-down-last-line | pass | 1.3 | no | — | — |
| combo-alt-left | pass | 3.5 | no | — | — |
| combo-alt-right | pass | 3.5 | no | — | — |
| combo-ctrl-home | pass | 1.1 | no | — | — |
| combo-ctrl-end | pass | 1.3 | no | — | — |
| bs-plain | pass | 1.0 | no | — | — |
| wrap-down-one-visual-line | pass | 1.1 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.5 | no | — | — |
| undo-redo-len | pass | 0.9 | no | — | — |
| gap-up-at-doc-start | pass | 1.2 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.3 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.3 | no | — | — |
| gap-collapse-selection-left | pass | 1.0 | no | — | — |
| gap-collapse-selection-right | pass | 1.0 | no | — | — |
| gap-delete-forward | pass | 1.1 | no | — | — |
| gap-delete-with-selection | fail | 0.9 | no | — | step 3: selEnd want 1278 got 1277; selLen want 4 got 3; state={1277 1274 1277 1551 Writerdeck harness dummy — ikke i vanlig notatliste  Før… |
| gap-select-all | pass | 0.6 | no | — | — |
| gap-enter-new-line | pass | 1.0 | no | — | — |
| gap-type-replaces-selection | pass | 1.2 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 1.2 | no | — | — |
| gap-undo-chain | pass | 0.9 | no | — | — |
| gap-empty-doc-backspace | pass | 0.4 | no | — | — |
| gap-shift-down-mid-wrapping-paras | pass | 3.7 | no | — | — |
| gap-shift-up-mid-wrapping-paras | pass | 3.4 | no | — | — |
