# Editor testing

Automated typing and selection checks on the real tablet, over the same path the phone uses. USB layout quirks still need a human. To capture a phone-typed bug for Cursor: phone Observe ([browser-vs-tablet.md](../browser-vs-tablet.md)); USB-on-tablet keys are not recorded.

Start here: [todo.md](todo.md).

Green checks are a score, not proof basic editing works — see [decisions.md](../decisions.md) **Typing-test strategy is failing**. Method theory: [methodology-shortcomings.md](methodology-shortcomings.md) (keep updated). How to run: §13 there. Gotchas: [lessons.md](../lessons.md).

todo.md — current scores and next steps.

methodology-shortcomings.md — why the typing-test attempt keeps failing; revise when we learn new miss patterns.

basic-claims.md — writer claims × kill-tests × guarded/partial/unguarded; update with every green-suite miss.

milestone-runs.md — pass/fail log of full runs; update after each full run.

scenario-catalog.md — all checks by name.

scenario-cookbook.md — where checks came from (CodeMirror, Qt, Ace).

Code lives in `daemon/cmd/edit-harness/` and `scripts/test-keyboard-harness.sh`.
