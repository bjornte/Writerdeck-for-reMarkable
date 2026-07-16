# Editor testing

Keyboard and selection regression for Writerdeck edit mode. Phone/WebSocket path via `edit-harness`; USB is spot-check only.

**Start here for a new agent:** [todo.md](todo.md) (scores, critical failures, next batch, do-not-retry).

Sign-off: **110/110/0** (`bash scripts/test-keyboard-harness.sh --fast`). Critical gate: **38/38/0** (`-t critical --fast`). Scoreboard columns use total/passed/failed.

Current baseline: **110/110/0** @ `00-29-12` @ fork `67656e1`; critical **38/38/0**. Report: `docs/recon/test-keyboard-harness-2026-07-17T00-29-12.md`.

| File | Purpose |
|------|---------|
| [todo.md](todo.md) | **Fresh agent handoff** — scores, critical failures, next batch, do-not-retry |
| [milestone-runs.md](milestone-runs.md) | **Full-suite scoreboard** — update after each `--fast` full run |
| [scenario-catalog.md](scenario-catalog.md) | All 110 scenarios — business-logic inventory |
| [harness-runs.md](../recon/harness-runs.md) | Consolidated run log |
| [scenario-cookbook.md](scenario-cookbook.md) | Source catalogs (CodeMirror/Qt) |
| [llm-handoff-test-failures.md](llm-handoff-test-failures.md) | Historical methodology review |

Code: `daemon/cmd/edit-harness/`, `scripts/test-keyboard-harness.sh`. Policy: [decisions.md](../decisions.md) §22. Gotchas: [lessons.md](../lessons.md) § Keyboard and selection.
