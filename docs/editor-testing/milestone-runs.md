# Milestone full-suite runs

Scoreboard for full keyboard tests on the tablet (`bash scripts/test-keyboard-harness.sh --fast`). Add a row after every full run. Detailed logs: `docs/recon/test-keyboard-harness-*.{md,txt}` (older runs in [harness-runs.md](../recon/harness-runs.md)).

**Critical pass** = the smaller “basic editing” gate (36 tests today). Sign-off needs the full suite green: **107/107**.

| Time of test | Num of tests | Total pass | Critical pass | Fail | Comments |
|--------------|--------------|------------|---------------|------|----------|
| 14 Jul 2026, 23:24 | 62 | 37 | — | 25 | First baseline |
| 15 Jul 2026, 04:45 | 83 | 64 | — | 16 | Shortcut chords and wrapped lines mostly fixed |
| 15 Jul 2026, 09:47 | 90 | 68 | — | 21 | Suite grew to 90; undo and touch still broken |
| 15 Jul 2026, 23:53 | 90 | 80 | 34/34 | 9 | Undo works; basic editing gate green |
| 16 Jul 2026, 00:37 | 94 | 89 | — | 4 | Arrow keys and page buttons behave; almost clean |
| 16 Jul 2026, 01:54 | 102 | 85 | — | 17 | Bigger, more realistic test set |
| 16 Jul 2026, 10:01 | 105 | 72 | 26/36 | 33 | Full 105 tests; harder cases dropped the score |
| 16 Jul 2026, 12:41 | 105 | 91 | 36/36 | 14 | Shift-select fixed; basic editing gate fully green |
| 16 Jul 2026, 18:57 | 105 | 93 | 36/36 | 12 | Editor lives in our own fork now |
| 16 Jul 2026, 21:21 | 105 | 105 | 36/36 | 0 | Sign-off — all tests pass |
| 16 Jul 2026, 23:12 | 107 | 107 | 36/36 | 0 | Mid-sentence Shift+Up/Down across wrapping paragraphs fixed |

Other landmarks: shortcut-chord tag **25/25**; wrap **15/15**; undo **5/5**. Fork migration finished the same day as the first sign-off.
