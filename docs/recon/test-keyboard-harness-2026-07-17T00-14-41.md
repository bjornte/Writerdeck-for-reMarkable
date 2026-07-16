# Keyboard harness results

Run: 2026-07-17T00:14:41+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 0 pass, 0 fail, 38 prepare fail; total 57.3s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | prepare fail | 1.5 | yes | — | write: read note HTTP 403; prepare retries |
| home-clears-selection | prepare fail | 1.5 | yes | load-cursor-at-start | write: read note HTTP 403; prepare retries; cascade suspect after load-cursor-at-start |
| shift-right-from-home | prepare fail | 1.5 | yes | home-clears-selection | write: read note HTTP 403; prepare retries; cascade suspect after home-clears-selection |
| shift-left-from-end | prepare fail | 1.5 | yes | shift-right-from-home | write: read note HTTP 403; prepare retries; cascade suspect after shift-right-from-home |
| shift-right-after-home-no-stale-anchor | prepare fail | 1.5 | yes | shift-left-from-end | write: read note HTTP 403; prepare retries; cascade suspect after shift-left-from-end |
| shift-down-after-arrow-down | prepare fail | 1.5 | yes | shift-right-after-home-no-stale-anchor | write: read note HTTP 403; prepare retries; cascade suspect after shift-right-after-home-no-stale-anchor |
| shift-up-after-arrow-down | prepare fail | 1.5 | yes | shift-down-after-arrow-down | write: read note HTTP 403; prepare retries; cascade suspect after shift-down-after-arrow-down |
| ctrl-shift-left-select-line | prepare fail | 1.5 | yes | shift-up-after-arrow-down | write: read note HTTP 403; prepare retries; cascade suspect after shift-up-after-arrow-down |
| down-one-logical-line | prepare fail | 1.5 | yes | ctrl-shift-left-select-line | write: read note HTTP 403; prepare retries; cascade suspect after ctrl-shift-left-select-line |
| shift-left-repeat-from-end | prepare fail | 1.5 | yes | down-one-logical-line | write: read note HTTP 403; prepare retries; cascade suspect after down-one-logical-line |
| alt-backspace-deletes-word | prepare fail | 1.5 | yes | shift-left-repeat-from-end | write: read note HTTP 403; prepare retries; cascade suspect after shift-left-repeat-from-end |
| ctrl-backspace-deletes-line | prepare fail | 1.5 | yes | alt-backspace-deletes-word | write: read note HTTP 403; prepare retries; cascade suspect after alt-backspace-deletes-word |
| shift-left-repeat-mid-doc | prepare fail | 1.5 | yes | ctrl-backspace-deletes-line | write: read note HTTP 403; prepare retries; cascade suspect after ctrl-backspace-deletes-line |
| cm-line-down-basic | prepare fail | 1.5 | yes | shift-left-repeat-mid-doc | write: read note HTTP 403; prepare retries; cascade suspect after shift-left-repeat-mid-doc |
| cm-line-down-last-line | prepare fail | 1.5 | yes | cm-line-down-basic | write: read note HTTP 403; prepare retries; cascade suspect after cm-line-down-basic |
| combo-alt-left | prepare fail | 1.5 | yes | cm-line-down-last-line | write: read note HTTP 403; prepare retries; cascade suspect after cm-line-down-last-line |
| combo-alt-right | prepare fail | 1.5 | yes | combo-alt-left | write: read note HTTP 403; prepare retries; cascade suspect after combo-alt-left |
| combo-ctrl-home | prepare fail | 1.5 | yes | combo-alt-right | write: read note HTTP 403; prepare retries; cascade suspect after combo-alt-right |
| combo-ctrl-end | prepare fail | 1.5 | yes | combo-ctrl-home | write: read note HTTP 403; prepare retries; cascade suspect after combo-ctrl-home |
| bs-plain | prepare fail | 1.5 | yes | combo-ctrl-end | write: read note HTTP 403; prepare retries; cascade suspect after combo-ctrl-end |
| wrap-down-one-visual-line | prepare fail | 1.5 | yes | bs-plain | write: read note HTTP 403; prepare retries; cascade suspect after bs-plain |
| wrap-up-from-visual-line-2 | prepare fail | 1.5 | yes | wrap-down-one-visual-line | write: read note HTTP 403; prepare retries; cascade suspect after wrap-down-one-visual-line |
| undo-redo-len | prepare fail | 1.5 | yes | wrap-up-from-visual-line-2 | write: read note HTTP 403; prepare retries; cascade suspect after wrap-up-from-visual-line-2 |
| gap-up-at-doc-start | prepare fail | 1.5 | yes | undo-redo-len | write: read note HTTP 403; prepare retries; cascade suspect after undo-redo-len |
| gap-plain-left-moves-caret | prepare fail | 1.5 | yes | gap-up-at-doc-start | write: read note HTTP 403; prepare retries; cascade suspect after gap-up-at-doc-start |
| gap-plain-right-moves-caret | prepare fail | 1.5 | yes | gap-plain-left-moves-caret | write: read note HTTP 403; prepare retries; cascade suspect after gap-plain-left-moves-caret |
| gap-collapse-selection-left | prepare fail | 1.5 | yes | gap-plain-right-moves-caret | write: read note HTTP 403; prepare retries; cascade suspect after gap-plain-right-moves-caret |
| gap-collapse-selection-right | prepare fail | 1.5 | yes | gap-collapse-selection-left | write: read note HTTP 403; prepare retries; cascade suspect after gap-collapse-selection-left |
| gap-delete-forward | prepare fail | 1.5 | yes | gap-collapse-selection-right | write: read note HTTP 403; prepare retries; cascade suspect after gap-collapse-selection-right |
| gap-delete-with-selection | prepare fail | 1.5 | yes | gap-delete-forward | write: read note HTTP 403; prepare retries; cascade suspect after gap-delete-forward |
| gap-select-all | prepare fail | 1.5 | yes | gap-delete-with-selection | write: read note HTTP 403; prepare retries; cascade suspect after gap-delete-with-selection |
| gap-enter-new-line | prepare fail | 1.5 | yes | gap-select-all | write: read note HTTP 403; prepare retries; cascade suspect after gap-select-all |
| gap-type-replaces-selection | prepare fail | 1.5 | yes | gap-enter-new-line | write: read note HTTP 403; prepare retries; cascade suspect after gap-enter-new-line |
| gap-redo-shift-ctrl-z | prepare fail | 1.5 | yes | gap-type-replaces-selection | write: read note HTTP 403; prepare retries; cascade suspect after gap-type-replaces-selection |
| gap-undo-chain | prepare fail | 1.5 | yes | gap-redo-shift-ctrl-z | write: read note HTTP 403; prepare retries; cascade suspect after gap-redo-shift-ctrl-z |
| gap-empty-doc-backspace | prepare fail | 1.5 | yes | gap-undo-chain | write: read note HTTP 403; prepare retries; cascade suspect after gap-undo-chain |
| gap-shift-down-mid-wrapping-paras | prepare fail | 1.5 | yes | gap-empty-doc-backspace | write: read note HTTP 403; prepare retries; cascade suspect after gap-empty-doc-backspace |
| gap-shift-up-mid-wrapping-paras | prepare fail | 1.5 | yes | gap-shift-down-mid-wrapping-paras | write: read note HTTP 403; prepare retries; cascade suspect after gap-shift-down-mid-wrapping-paras |
