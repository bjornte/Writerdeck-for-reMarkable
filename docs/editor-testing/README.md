# Editor testing

Automated typing and selection checks on the real tablet, over the same path the phone uses. USB layout quirks still need a human.

Start here: [todo.md](todo.md).

Sign-off is the full suite green (**110/110/0**). The critical gate is **38/38/0**. Policy: [decisions.md](../decisions.md) §13. Gotchas: [lessons.md](../lessons.md).

todo.md — current scores and next steps.

milestone-runs.md — full-suite scoreboard; update after each full run.

scenario-catalog.md — all scenarios by name.

scenario-cookbook.md — where scenarios came from (CodeMirror, Qt, Ace).

Code lives in `daemon/cmd/edit-harness/` and `scripts/test-keyboard-harness.sh`.
