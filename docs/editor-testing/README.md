# Editor testing

Keyboard and selection regression for Writerdeck edit mode. Phone/WebSocket path via `edit-harness`; USB is spot-check only.

| File | Purpose |
|------|---------|
| [todo.md](todo.md) | **Fresh agent handoff** — baseline 37/25, failure clusters, do-not-retry, next steps |
| [scenario-cookbook.md](scenario-cookbook.md) | Scenario specs; most blocks now ported into `scenarios_*.go` |
| [llm-handoff-test-failures.md](llm-handoff-test-failures.md) | Historical methodology review; see todo.md for current state |

Code: `daemon/cmd/edit-harness/`, `scripts/test-keyboard-harness.sh`. Policy: [decisions.md](../decisions.md) §22. Gotchas and batch workflow: [lessons.md](../lessons.md) § Keyboard and selection.
