# Editor testing

Keyboard and selection regression for Writerdeck edit mode. Phone/WebSocket path via `edit-harness`; USB is spot-check only.

| File | Purpose |
|------|---------|
| [todo.md](todo.md) | Open bugs, harness loop, acceptance, resume prompt |
| [scenario-cookbook.md](scenario-cookbook.md) | Scenarios borrowed from CodeMirror and Qt; port into `scenarios_regression.go` |

Code: `daemon/cmd/edit-harness/`, `scripts/test-keyboard-harness.sh`. Policy: [decisions.md](../decisions.md) §22. Gotchas: [lessons.md](../lessons.md) § Keyboard and selection.
