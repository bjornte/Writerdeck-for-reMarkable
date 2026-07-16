# Keyboard harness results

Run: 2026-07-16T17:20:44+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 34 pass, 2 fail, 0 prepare fail; total 66.5s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | pass | 0.4 | no | — | — |
| home-clears-selection | pass | 0.9 | no | — | — |
| shift-right-from-home | pass | 3.3 | no | — | — |
| shift-left-from-end | pass | 3.4 | no | — | — |
| shift-right-after-home-no-stale-anchor | fail | 2.3 | no | — | step 16: selStart want 1280 got 1281; selLen want 2 got 1; state={1282 1281 1282 1551 Writerdeck harness dummy — ikke i vanlig notatliste  … |
| shift-down-after-arrow-down | pass | 3.7 | no | — | — |
| shift-up-after-arrow-down | pass | 3.9 | no | — | — |
| ctrl-shift-left-select-line | pass | 0.5 | no | — | — |
| down-one-logical-line | pass | 3.4 | no | — | — |
| shift-left-repeat-from-end | pass | 3.4 | no | — | — |
| alt-backspace-deletes-word | fail | 1.7 | no | — | step 6: cursor want 1155 got 1161; state={1161 1161 1161 1527 Writerdeck harness dummy — ikke i vanlig notatliste  Første avsnitt — Naïve c… |
| ctrl-backspace-deletes-line | pass | 1.6 | no | — | — |
| shift-left-repeat-mid-doc | pass | 3.7 | no | — | — |
| cm-line-down-basic | pass | 3.3 | no | — | — |
| cm-line-down-last-line | pass | 1.1 | no | — | — |
| combo-alt-left | pass | 3.4 | no | — | — |
| combo-alt-right | pass | 3.4 | no | — | — |
| combo-ctrl-home | pass | 1.1 | no | — | — |
| combo-ctrl-end | pass | 1.3 | no | — | — |
| bs-plain | pass | 1.0 | no | — | — |
| wrap-down-one-visual-line | pass | 1.7 | no | — | — |
| wrap-up-from-visual-line-2 | pass | 1.5 | no | — | — |
| undo-redo-len | pass | 1.0 | no | — | — |
| gap-up-at-doc-start | pass | 1.0 | no | — | — |
| gap-plain-left-moves-caret | pass | 3.3 | no | — | — |
| gap-plain-right-moves-caret | pass | 3.3 | no | — | — |
| gap-collapse-selection-left | pass | 0.9 | no | — | — |
| gap-collapse-selection-right | pass | 0.8 | no | — | — |
| gap-delete-forward | pass | 1.2 | no | — | — |
| gap-delete-with-selection | pass | 0.8 | no | — | — |
| gap-select-all | pass | 0.7 | no | — | — |
| gap-enter-new-line | pass | 0.6 | no | — | — |
| gap-type-replaces-selection | pass | 0.8 | no | — | — |
| gap-redo-shift-ctrl-z | pass | 0.9 | no | — | — |
| gap-undo-chain | pass | 0.8 | no | — | — |
| gap-empty-doc-backspace | pass | 0.5 | no | — | — |
