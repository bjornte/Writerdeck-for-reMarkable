# Keyboard harness

Prove Mac/Linux-style editing on the tablet with the automated harness — not by typing in the Lobby by hand.

Read: this file, [milestone-runs.md](milestone-runs.md), [lessons.md](../lessons.md), [decisions.md](../decisions.md) §13, [terms.md](../terms.md). Scenario names: [scenario-catalog.md](scenario-catalog.md).

## Current score

Full suite **110/110/0** @ `17-23-47` (fork `0bb3b70`). Critical **38/38/0**. Wrap and undo tags green. Edit-session passed the same day. Do not run edit-session in parallel with this harness.

## Next

Owner Physical Home check ([user-should-test.md](../user-should-test.md)). Keep the harness green on any future edit change: one push, one deploy, edit-session, full `--fast`, update the scoreboard.

## Do not

Infer the moving end of a selection from the caret after select. Treat keyboard arrows as page-scroll. Deploy once per failing scenario. Auto-release Escape in inject helpers. Step wrapped lines by a tall caret rectangle alone.

## Acceptance

Critical **38/38/0**, full **110/110/0**, edit-session pass, clean journal after deploy.
