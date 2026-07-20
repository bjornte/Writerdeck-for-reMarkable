# Automated typing tests

Prove Mac/Linux-style editing on the tablet with scripted tests — not by typing in the Lobby by hand.

Read: this file, [milestone-runs.md](milestone-runs.md), [lessons.md](../lessons.md), [decisions.md](../decisions.md) (**Typing-test strategy is failing**, then §13), [terms.md](../terms.md). Check names: [scenario-catalog.md](scenario-catalog.md).

## Test strategy status (do not remove yet)

Same rule as [decisions.md](../decisions.md) **Typing-test strategy is failing**: we do not have a working strategy yet — only an attempt that keeps going green while basic editing bugs still reach a person. Keep this banner until that decisions section is lifted with solid proof the misses have stopped recurring. Method theory (update when we learn a new failure mode): [methodology-shortcomings.md](methodology-shortcomings.md). Claim kill-tests: [basic-claims.md](basic-claims.md).

## Current score

All 125 checks passed at `21-43-24` (fork commit `19792fc`). Basic set **57** of 57 (`21-32-07`): wrap Home/End, End-then-Down kill-test, goal-column Up round-trip. Do not run edit-session at the same time as these typing tests. Treat this score under the status section above.

## Next

Owner Physical Home check ([user-should-test.md](../user-should-test.md)). Keep all typing checks passing on any future edit change: one push, one deploy, edit-session, full `--fast`, update the pass/fail log. When a human finds a basic bug that was green, fix the test’s ability to catch that failure and update [basic-claims.md](basic-claims.md) before calling the strategy improved. Priority raises are listed at the bottom of that inventory.

## Do not

Infer the moving end of a selection from the caret after select. Treat keyboard arrows as page-scroll. Deploy once per failing case. Auto-release Escape in key-feed helpers. Step wrapped lines by a tall caret rectangle alone.

## Done means

Basic set 57 of 57, full set 125 of 125, edit-session pass, clean journal after deploy.
