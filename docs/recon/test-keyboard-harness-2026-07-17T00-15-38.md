# Keyboard harness results

Run: 2026-07-17T00:15:38+02:00

Target: `192.168.1.8:8000`

Mode: sandbox-prepare (single session)

Timing: fast pauses

Summary: 0 pass, 0 fail, 15 prepare fail; total 22.7s

| Name | Result | Time (s) | Recovery | Cascade | Comments |
|------|--------|----------|----------|---------|----------|
| wrap-down-one-visual-line | prepare fail | 1.5 | yes | — | write: read note HTTP 403; prepare retries |
| wrap-down-not-jump-paragraph | prepare fail | 1.5 | yes | wrap-down-one-visual-line | write: read note HTTP 403; prepare retries; cascade suspect after wrap-down-one-visual-line |
| wrap-up-from-visual-line-2 | prepare fail | 1.5 | yes | wrap-down-not-jump-paragraph | write: read note HTTP 403; prepare retries; cascade suspect after wrap-down-not-jump-paragraph |
| wrap-shift-down-one-visual | prepare fail | 1.5 | yes | wrap-up-from-visual-line-2 | write: read note HTTP 403; prepare retries; cascade suspect after wrap-up-from-visual-line-2 |
| wrap-shift-down-then-up-shrinks | prepare fail | 1.5 | yes | wrap-shift-down-one-visual | write: read note HTTP 403; prepare retries; cascade suspect after wrap-shift-down-one-visual |
| wrap-down-last-visual-line | prepare fail | 1.5 | yes | wrap-shift-down-then-up-shrinks | write: read note HTTP 403; prepare retries; cascade suspect after wrap-shift-down-then-up-shrinks |
| wrap-shift-down-last-to-eof | prepare fail | 1.5 | yes | wrap-down-last-visual-line | write: read note HTTP 403; prepare retries; cascade suspect after wrap-down-last-visual-line |
| wrap-mixed-newline-and-wrap | prepare fail | 1.5 | yes | wrap-shift-down-last-to-eof | write: read note HTTP 403; prepare retries; cascade suspect after wrap-shift-down-last-to-eof |
| wrap-down-goal-column | prepare fail | 1.5 | yes | wrap-mixed-newline-and-wrap | write: read note HTTP 403; prepare retries; cascade suspect after wrap-mixed-newline-and-wrap |
| wrap-combo-alt-left-word | prepare fail | 1.5 | yes | wrap-down-goal-column | write: read note HTTP 403; prepare retries; cascade suspect after wrap-down-goal-column |
| wrap-combo-alt-right-word | prepare fail | 1.5 | yes | wrap-combo-alt-left-word | write: read note HTTP 403; prepare retries; cascade suspect after wrap-combo-alt-left-word |
| wrap-combo-ctrl-bs-line | prepare fail | 1.5 | yes | wrap-combo-alt-right-word | write: read note HTTP 403; prepare retries; cascade suspect after wrap-combo-alt-right-word |
| wrap-shift-left-across-wrap | prepare fail | 1.5 | yes | wrap-combo-ctrl-bs-line | write: read note HTTP 403; prepare retries; cascade suspect after wrap-combo-ctrl-bs-line |
| wrap-home-on-visual-line | prepare fail | 1.5 | yes | wrap-shift-left-across-wrap | write: read note HTTP 403; prepare retries; cascade suspect after wrap-shift-left-across-wrap |
| wrap-end-on-visual-line | prepare fail | 1.5 | yes | wrap-home-on-visual-line | write: read note HTTP 403; prepare retries; cascade suspect after wrap-home-on-visual-line |
