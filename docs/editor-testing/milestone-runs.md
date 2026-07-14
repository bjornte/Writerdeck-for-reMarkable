# Milestone full-suite runs

Hand-maintained scoreboard for full `bash scripts/test-keyboard-harness.sh --fast` sessions (no `-s`, no `--tag`). Update this table after every full-suite run — add a row with timestamp, pass/fail/prep, delta vs prior milestone, and one-line context (commit or what changed). Per-run reports: `docs/recon/test-keyboard-harness-*.{md,txt}` (new runs only; older files consolidated in [harness-runs.md](../recon/harness-runs.md)).

Sign-off gate: **83/83 PASS**.

| Run | Suite | Pass | Fail | Prep | vs prior | Context |
|-----|-------|------|------|------|----------|---------|
| 2026-07-14T20-15-40 | 62 | 26 | 36 | 0 | — | early session |
| 2026-07-14T22-06-06 | 62 | 32 | 30 | 0 | +6 | harness hardening |
| 2026-07-14T23-06-59 | 62 | 36 | 26 | 0 | +4 | pre-baseline |
| 2026-07-14T23-24-42 | 62 | 37 | 25 | 0 | +1 | **baseline anchor** |
| 2026-07-15T00-08-41 | 83 | 38 | 44 | 1 | new gaps | **best 83** (pre-QML) |
| 2026-07-15T00-17-48 | 83 | 34 | 48 | 1 | −4 | df2f850 QML deploy regression |
| 2026-07-15T00-43-13 | 83 | 27 | 53 | 1 | −7 | 4c4d816 worst 83 |
| 2026-07-15T00-56-17 | 83 | 35 | 45 | 1 | +8 | 0a339c9 partial recovery |

Combo-tag milestones (22 scenarios, `--tag combo`): 6/16 @ baseline → 9/13 @ `22ad701` (`01-25-41`). No full 83 run on `22ad701` yet.
