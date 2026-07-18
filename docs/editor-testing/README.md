# Editor testing

Automated typing and selection checks on the real tablet, over the same path the phone uses. USB layout quirks still need a human.

Start here: [todo.md](todo.md).

Calling typing work done means all 112 checks passed. The basic “editing works” set is 40 checks. Policy: [decisions.md](../decisions.md) §13. Gotchas: [lessons.md](../lessons.md).

todo.md — current scores and next steps.

milestone-runs.md — pass/fail log of full runs; update after each full run.

scenario-catalog.md — all checks by name.

scenario-cookbook.md — where checks came from (CodeMirror, Qt, Ace).

Code lives in `daemon/cmd/edit-harness/` and `scripts/test-keyboard-harness.sh`.
