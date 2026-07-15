# Editor testing

Keyboard and selection regression for Writerdeck edit mode. Phone/WebSocket path via `edit-harness`; USB is spot-check only. Sign-off: **87/87 PASS** (`--fast`).

| File | Purpose |
|------|---------|
| [scenario-catalog.md](scenario-catalog.md) | **All 87 scenarios** — business-logic inventory |
| [todo.md](todo.md) | **Fresh agent handoff** — scores, failure clusters, do-not-retry |
| [milestone-runs.md](milestone-runs.md) | **Full-suite scoreboard** — update after each `--fast` full run (no `-s`/`--tag`) |
| [harness-runs.md](../recon/harness-runs.md) | Consolidated run log and per-scenario matrix |
| [scenario-cookbook.md](scenario-cookbook.md) | Source catalogs (CodeMirror/Qt) and porting notation |
| [llm-handoff-test-failures.md](llm-handoff-test-failures.md) | Historical methodology review; see todo.md for current state |

Code: `daemon/cmd/edit-harness/`, `scripts/test-keyboard-harness.sh`. Policy: [decisions.md](../decisions.md) §22. Gotchas and batch workflow: [lessons.md](../lessons.md) § Keyboard and selection.
