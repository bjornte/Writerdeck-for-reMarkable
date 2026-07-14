# Keyboard harness run history

Consolidated from per-run recon files on 2026-07-15. Individual `test-keyboard-harness-*.md` reports pruned.

Sign-off gate: **83/83 PASS** (`bash scripts/test-keyboard-harness.sh --fast`).

## Run log (all sessions)

| Run (local) | Kind | N | Pass | Fail | Prep | Notes |
|-------------|------|---|------|------|------|-------|
| 2026-07-14T19-49-27 | partial | 18 | 13 | 5 | 0 | legacy soft-reset, 18 scenarios |
| 2026-07-14T20-15-40 | full-62 | 62 | 26 | 36 | 0 |  |
| 2026-07-14T20-28-06 | full-62 | 62 | 27 | 35 | 0 |  |
| 2026-07-14T20-42-51 | infra | 62 | 0 | 0 | 62 | infra: 62/62 prepare fail |
| 2026-07-14T22-06-06 | full-62 | 62 | 32 | 30 | 0 |  |
| 2026-07-14T23-06-59 | full-62 | 62 | 36 | 26 | 0 |  |
| 2026-07-14T23-14-30 | full-62 | 62 | 37 | 25 | 0 | duplicate baseline 37/25 |
| 2026-07-14T23-15-20 | isolated | 1 | 0 | 1 | 0 | combo-ctrl-right |
| 2026-07-14T23-15-34 | isolated | 1 | 0 | 1 | 0 | combo-ctrl-right |
| 2026-07-14T23-15-56 | isolated | 1 | 0 | 1 | 0 | combo-ctrl-right |
| 2026-07-14T23-16-04 | isolated | 1 | 0 | 1 | 0 | combo-ctrl-right |
| 2026-07-14T23-16-12 | isolated | 1 | 1 | 0 | 0 | combo-ctrl-left |
| 2026-07-14T23-16-21 | full-62 | 62 | 23 | 39 | 0 | mid-session regression 23/39 |
| 2026-07-14T23-24-24 | combo-tag | 22 | 6 | 16 | 0 |  |
| 2026-07-14T23-24-42 | full-62 | 62 | 37 | 25 | 0 |  |
| 2026-07-15T00-06-54 | full-83 | 83 | 26 | 56 | 1 | first 83-scenario run |
| 2026-07-15T00-08-41 | full-83 | 83 | 38 | 44 | 1 |  |
| 2026-07-15T00-17-48 | full-83 | 83 | 34 | 48 | 1 |  |
| 2026-07-15T00-43-13 | full-83 | 81 | 27 | 53 | 1 |  |
| 2026-07-15T00-48-45 | isolated | 1 | 0 | 1 | 0 | shift-right-from-home |
| 2026-07-15T00-48-54 | isolated | 1 | 0 | 1 | 0 | bs-plain |
| 2026-07-15T00-50-12 | isolated | 1 | 1 | 0 | 0 | shift-right-from-home |
| 2026-07-15T00-52-59 | infra | 83 | 0 | 0 | 83 | infra: server not restarted |
| 2026-07-15T00-56-17 | full-83 | 81 | 35 | 45 | 1 |  |
| 2026-07-15T01-00-00 | isolated | 1 | 0 | 1 | 0 | combo-ctrl-right |
| 2026-07-15T01-00-21 | isolated | 1 | 0 | 1 | 0 | combo-alt-right |
| 2026-07-15T01-01-00 | isolated | 1 | 0 | 1 | 0 | combo-ctrl-right |
| 2026-07-15T01-01-32 | isolated | 1 | 0 | 1 | 0 | combo-ctrl-right |
| 2026-07-15T01-02-35 | isolated | 1 | 1 | 0 | 0 | combo-ctrl-left |
| 2026-07-15T01-02-50 | isolated | 1 | 0 | 1 | 0 | combo-ctrl-right |
| 2026-07-15T01-03-04 | isolated | 1 | 1 | 0 | 0 | combo-ctrl-left |
| 2026-07-15T01-03-12 | isolated | 1 | 1 | 0 | 0 | combo-ctrl-left |
| 2026-07-15T01-05-40 | isolated | 1 | 0 | 1 | 0 | combo-ctrl-right |
| 2026-07-15T01-05-57 | isolated | 1 | 0 | 1 | 0 | combo-alt-right |
| 2026-07-15T01-09-02 | isolated | 1 | 0 | 1 | 0 | deadlock: socketRouteKey on GUI thread |
| 2026-07-15T01-09-09 | isolated | 1 | 0 | 0 | 1 | combo-ctrl-left |
| 2026-07-15T01-11-19 | isolated | 1 | 0 | 0 | 1 | combo-alt-right |
| 2026-07-15T01-13-29 | isolated | 1 | 0 | 1 | 0 | combo-alt-left |
| 2026-07-15T01-15-18 | isolated | 1 | 0 | 1 | 0 | combo-shift-ctrl-right |
| 2026-07-15T01-18-17 | combo-tag | 22 | 6 | 16 | 0 |  |
| 2026-07-15T01-18-46 | isolated | 1 | 0 | 1 | 0 | combo-ctrl-right |
| 2026-07-15T01-22-07 | isolated | 1 | 0 | 1 | 0 | combo-ctrl-right |
| 2026-07-15T01-25-33 | isolated | 1 | 1 | 0 | 0 | combo-ctrl-right |
| 2026-07-15T01-25-34 | isolated | 1 | 1 | 0 | 0 | combo-ctrl-down |
| 2026-07-15T01-25-35 | isolated | 1 | 1 | 0 | 0 | combo-ctrl-end |
| 2026-07-15T01-25-36 | isolated | 1 | 1 | 0 | 0 | combo-ctrl-left |
| 2026-07-15T01-25-41 | combo-tag | 22 | 9 | 13 | 0 | combo tag @22ad701 — best combo (9/13) |

## Milestone full suites

Canonical table: [docs/editor-testing/milestone-runs.md](../editor-testing/milestone-runs.md) (update after each full `--fast` run).

| Run | Suite | Pass | Fail | Prep | vs prior | Context |
|-----|-------|------|------|------|----------|---------|
| 2026-07-14T20-15-40 | 62 | 26 | 36 | 0 | +6 | early session, pre-harness |
| 2026-07-14T22-06-06 | 62 | 32 | 30 | 0 | +6 | harness hardening |
| 2026-07-14T23-06-59 | 62 | 36 | 26 | 0 | +4 | pre-baseline |
| 2026-07-14T23-24-42 | 62 | 37 | 25 | 0 | +1 | **baseline anchor 37/25** |
| 2026-07-15T00-08-41 | 83 | 38 | 44 | 1 | new gaps | **best 83: 38/44** pre-QML |
| 2026-07-15T00-17-48 | 83 | 34 | 48 | 1 | −4 | df2f850 QML deploy regression |
| 2026-07-15T00-43-13 | 83 | 27 | 53 | 1 | −7 | 4c4d816 worst 83 (27/53) |
| 2026-07-15T00-56-17 | 83 | 35 | 45 | 1 | +8 | 0a339c9 partial recovery (35/45) |

## Combo tag (22 scenarios)

| Run | Pass | Fail | Context |
|-----|------|------|---------|
| 2026-07-14T23-24-24 | 6 | 16 | at 62 baseline |
| 2026-07-15T01-18-17 | 6 | 16 | 9bbf282 socketRouteKey |
| 2026-07-15T01-25-41 | 9 | 13 | 22ad701 Ctrl fast-path — ctrl-right/down/end pass |

## Per-scenario matrix (core 62)

P=pass F=fail. Columns: baseline → best 83 → latest 83.

| Scenario | base 37/25 | best83 38/44 | latest83 35/45 | @22ad701 |
|----------|------------|--------------|----------------|----------|
| load-cursor-at-start | P | P | P | — |
| home-clears-selection | P | P | P | — |
| shift-right-from-home | P | P | P | — |
| shift-left-from-end | P | P | P | — |
| shift-right-after-home-no-stale-anchor | P | P | P | — |
| shift-down-after-arrow-down | P | P | P | — |
| shift-up-after-arrow-down | P | P | P | — |
| ctrl-shift-left-select-line | P | P | P | — |
| down-one-logical-line | P | P | P | — |
| shift-down-then-up-shrinks | P | P | P | — |
| shift-left-repeat-from-end | P | P | P | — |
| alt-backspace-deletes-word | P | P | P | — |
| ctrl-backspace-deletes-line | P | P | P | — |
| shift-left-repeat-mid-doc | P | P | P | — |
| cm-line-down-basic | P | P | P | — |
| cm-line-down-shorter | P | P | P | — |
| cm-line-down-last-line | P | P | P | — |
| cm-line-down-goal-col | F | F | F | — |
| cm-select-line-down | P | P | P | — |
| cm-select-line-down-mid | F | F | F | — |
| cm-select-down-up-doc-end | F | F | F | — |
| cm-select-up-basic | F | F | F | — |
| cm-select-up-mid | F | F | F | — |
| combo-alt-left | F | F | F | — |
| combo-alt-right | F | F | F | — |
| combo-alt-up | P | P | P | — |
| combo-alt-down | F | F | F | — |
| combo-ctrl-left | P | P | P | P |
| combo-ctrl-right | F | F | F | P |
| combo-ctrl-up | P | P | P | P |
| combo-ctrl-down | F | F | F | P |
| combo-shift-alt-left | F | F | F | — |
| combo-shift-alt-right | F | F | F | — |
| combo-shift-alt-up | F | F | F | — |
| combo-shift-alt-down | F | F | F | — |
| combo-shift-ctrl-left | F | F | F | — |
| combo-shift-ctrl-right | F | F | F | — |
| combo-shift-ctrl-up | F | F | F | — |
| combo-shift-ctrl-down | F | F | F | — |
| combo-shift-home-line | P | P | P | — |
| combo-shift-end-line | P | P | P | — |
| combo-ctrl-home | P | P | P | — |
| combo-ctrl-end | F | F | F | P |
| combo-shift-ctrl-home | F | F | F | — |
| combo-shift-ctrl-end | F | F | F | — |
| bs-alt-word-mid | P | P | P | — |
| bs-ctrl-line-start | P | P | P | — |
| bs-shift-with-selection | P | P | P | — |
| bs-plain | P | F | P | — |
| wrap-down-one-visual-line | P | F | F | — |
| wrap-down-not-jump-paragraph | P | F | F | — |
| wrap-up-from-visual-line-2 | F | F | F | — |
| wrap-shift-down-one-visual | P | F | F | — |
| wrap-shift-down-then-up-shrinks | P | P | F | — |
| wrap-down-last-visual-line | P | P | F | — |
| wrap-shift-down-last-to-eof | P | P | F | — |
| wrap-mixed-newline-and-wrap | P | P | P | — |
| undo-redo-len | F | F | F | — |
| undo-cursor-reposition | F | F | F | — |
| undo-mid-line-delete | F | F | F | — |
| redo-cleared-by-new-edit | P | P | F | — |
| undo-after-select-delete | P | P | F | — |

## Regressions

Pass @ baseline → fail @ best 83: 4 scenarios
- `bs-plain`
- `wrap-down-one-visual-line`
- `wrap-down-not-jump-paragraph`
- `wrap-shift-down-one-visual`

Pass @ best 83 → fail @ post-QML (00-17-48): 4 scenarios
- `combo-shift-home-line`
- `wrap-shift-down-then-up-shrinks`
- `wrap-down-last-visual-line`
- `wrap-shift-down-last-to-eof`

Recovered @ 0a339c9 (00-56-17): 5 scenarios
- `combo-shift-home-line`
- `bs-plain`
- `gap-delete-with-selection`
- `gap-enter-new-line`
- `gap-type-replaces-selection`

## Backtrack guidance

Best full-suite score so far: **38/44** @ `00-08-41` (harness fixes, pre-`df2f850` QML). QML deploys `df2f850`/`4c4d816` regressed to 27–34 pass; `0a339c9` recovered to 35 but not past 38.

Best combo tag: **9/13** @ `01-25-41` (`22ad701` socketRouteKey + Ctrl fast-path). No full 83 run on that build yet.

Do not wholesale revert editor — socket routing on current head (`22ad701`) fixes a real delivery bug. Re-run full `--fast` on that build and add a row to [milestone-runs.md](../editor-testing/milestone-runs.md).
