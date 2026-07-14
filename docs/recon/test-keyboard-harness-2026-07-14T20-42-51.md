# Keyboard harness results

Run: 2026-07-14T20:42:56+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 0 pass, 0 fail, 62 prepare fail; total 681.3s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| load-cursor-at-start | prepare fail | 209.3 | yes | — | ensure editor: editor not ready after open; prepare retries |
| home-clears-selection | prepare fail | 63.0 | yes | load-cursor-at-start | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: host is down; prepare retrie… |
| shift-right-from-home | prepare fail | 1.5 | yes | home-clears-selection | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: host is down; prepare retrie… |
| shift-left-from-end | prepare fail | 1.5 | yes | shift-right-from-home | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: host is down; prepare retrie… |
| shift-right-after-home-no-stale-anchor | prepare fail | 1.5 | yes | shift-left-from-end | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: host is down; prepare retrie… |
| shift-down-after-arrow-down | prepare fail | 1.5 | yes | shift-right-after-home-no-stale-anchor | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: host is down; prepare retrie… |
| shift-up-after-arrow-down | prepare fail | 1.5 | yes | shift-down-after-arrow-down | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: host is down; prepare retrie… |
| ctrl-shift-left-select-line | prepare fail | 58.5 | yes | shift-up-after-arrow-down | ensure editor: editor not ready after open; prepare retries; cascade suspect after shift-up-after-arrow-down |
| down-one-logical-line | prepare fail | 187.8 | yes | ctrl-shift-left-select-line | ensure editor: Post "http://192.168.1.8:8000/api/open": dial tcp 192.168.1.8:8000: connect: connection refused; prepare retries; cascade su… |
| shift-down-then-up-shrinks | prepare fail | 32.1 | yes | down-one-logical-line | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| shift-left-repeat-from-end | prepare fail | 2.3 | yes | shift-down-then-up-shrinks | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| alt-backspace-deletes-word | prepare fail | 2.5 | yes | shift-left-repeat-from-end | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| ctrl-backspace-deletes-line | prepare fail | 2.3 | yes | alt-backspace-deletes-word | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| shift-left-repeat-mid-doc | prepare fail | 2.4 | yes | ctrl-backspace-deletes-line | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| cm-line-down-basic | prepare fail | 2.4 | yes | shift-left-repeat-mid-doc | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| cm-line-down-shorter | prepare fail | 2.6 | yes | cm-line-down-basic | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| cm-line-down-last-line | prepare fail | 2.4 | yes | cm-line-down-shorter | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| cm-line-down-goal-col | prepare fail | 2.3 | yes | cm-line-down-last-line | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| cm-select-line-down | prepare fail | 2.4 | yes | cm-line-down-goal-col | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| cm-select-line-down-mid | prepare fail | 2.4 | yes | cm-select-line-down | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| cm-select-down-up-doc-end | prepare fail | 2.4 | yes | cm-select-line-down-mid | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| cm-select-up-basic | prepare fail | 2.3 | yes | cm-select-down-up-doc-end | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| cm-select-up-mid | prepare fail | 2.4 | yes | cm-select-up-basic | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-alt-left | prepare fail | 2.3 | yes | cm-select-up-mid | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-alt-right | prepare fail | 2.4 | yes | combo-alt-left | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-alt-up | prepare fail | 2.4 | yes | combo-alt-right | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-alt-down | prepare fail | 2.3 | yes | combo-alt-up | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-ctrl-left | prepare fail | 2.3 | yes | combo-alt-down | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-ctrl-right | prepare fail | 2.4 | yes | combo-ctrl-left | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-ctrl-up | prepare fail | 2.3 | yes | combo-ctrl-right | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-ctrl-down | prepare fail | 2.4 | yes | combo-ctrl-up | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-alt-left | prepare fail | 2.4 | yes | combo-ctrl-down | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-alt-right | prepare fail | 2.3 | yes | combo-shift-alt-left | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-alt-up | prepare fail | 2.3 | yes | combo-shift-alt-right | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-alt-down | prepare fail | 2.4 | yes | combo-shift-alt-up | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-ctrl-left | prepare fail | 2.4 | yes | combo-shift-alt-down | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-ctrl-right | prepare fail | 2.3 | yes | combo-shift-ctrl-left | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-ctrl-up | prepare fail | 2.3 | yes | combo-shift-ctrl-right | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-ctrl-down | prepare fail | 2.4 | yes | combo-shift-ctrl-up | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-home-line | prepare fail | 2.3 | yes | combo-shift-ctrl-down | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-end-line | prepare fail | 2.4 | yes | combo-shift-home-line | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-ctrl-home | prepare fail | 2.3 | yes | combo-shift-end-line | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-ctrl-end | prepare fail | 2.4 | yes | combo-ctrl-home | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-ctrl-home | prepare fail | 2.3 | yes | combo-ctrl-end | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| combo-shift-ctrl-end | prepare fail | 2.5 | yes | combo-shift-ctrl-home | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| bs-alt-word-mid | prepare fail | 2.3 | yes | combo-shift-ctrl-end | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| bs-ctrl-line-start | prepare fail | 2.4 | yes | bs-alt-word-mid | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| bs-shift-with-selection | prepare fail | 2.3 | yes | bs-ctrl-line-start | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| bs-plain | prepare fail | 2.4 | yes | bs-shift-with-selection | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| wrap-down-one-visual-line | prepare fail | 2.3 | yes | bs-plain | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| wrap-down-not-jump-paragraph | prepare fail | 2.4 | yes | wrap-down-one-visual-line | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| wrap-up-from-visual-line-2 | prepare fail | 2.3 | yes | wrap-down-not-jump-paragraph | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| wrap-shift-down-one-visual | prepare fail | 2.4 | yes | wrap-up-from-visual-line-2 | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| wrap-shift-down-then-up-shrinks | prepare fail | 2.3 | yes | wrap-shift-down-one-visual | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| wrap-down-last-visual-line | prepare fail | 2.4 | yes | wrap-shift-down-then-up-shrinks | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| wrap-shift-down-last-to-eof | prepare fail | 2.3 | yes | wrap-down-last-visual-line | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| wrap-mixed-newline-and-wrap | prepare fail | 2.5 | yes | wrap-shift-down-last-to-eof | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| undo-redo-len | prepare fail | 2.3 | yes | wrap-mixed-newline-and-wrap | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| undo-cursor-reposition | prepare fail | 2.4 | yes | undo-redo-len | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| undo-mid-line-delete | prepare fail | 2.3 | yes | undo-cursor-reposition | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| redo-cleared-by-new-edit | prepare fail | 2.4 | yes | undo-mid-line-delete | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
| undo-after-select-delete | prepare fail | 2.3 | yes | redo-cleared-by-new-edit | write: Get "http://192.168.1.8:8000/api/notes/z-test-keyboard-harness.md": dial tcp 192.168.1.8:8000: connect: connection refused; prepare … |
