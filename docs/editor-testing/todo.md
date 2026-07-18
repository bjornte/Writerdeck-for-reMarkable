# Automated typing tests

Prove Mac/Linux-style editing on the tablet with scripted tests — not by typing in the Lobby by hand.

Read: this file, [milestone-runs.md](milestone-runs.md), [lessons.md](../lessons.md), [decisions.md](../decisions.md) §13, [terms.md](../terms.md). Check names: [scenario-catalog.md](scenario-catalog.md).

## Current score

All 112 checks passed at `15-40-20` (fork commit `df1d38b`). Basic set 40 of 40, including copy/cut/paste. Edit-session and lobby-keyboard passed the same day. Do not run edit-session at the same time as these typing tests.

## Next

Owner Physical Home check ([user-should-test.md](../user-should-test.md)). Keep all typing checks passing on any future edit change: one push, one deploy, edit-session, full `--fast`, update the pass/fail log.

## Do not

Infer the moving end of a selection from the caret after select. Treat keyboard arrows as page-scroll. Deploy once per failing case. Auto-release Escape in key-feed helpers. Step wrapped lines by a tall caret rectangle alone.

## Done means

Basic set 40 of 40, full set 112 of 112, edit-session pass, clean journal after deploy.
