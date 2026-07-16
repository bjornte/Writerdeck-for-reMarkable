# Milestone full-suite runs

Scoreboard for full keyboard tests on the tablet (`bash scripts/test-keyboard-harness.sh --fast`). Add a row after every full run. Detailed logs: `docs/recon/test-keyboard-harness-*.{md,txt}` (older runs in [harness-runs.md](../recon/harness-runs.md)).

**Critical tests** = the smaller “basic editing” gate (38 today). Sign-off needs the full suite green: **110/110/0**.

| Time of test | All tests (total/passed/failed) | Critical tests (total/passed/failed) | Comments |
|--------------|---------------------------------|--------------------------------------|----------|
| 14 Jul 2026, 23:24 | 62/37/25 | — | First baseline |
| 15 Jul 2026, 23:53 | 90/80/9 | 34/34/0 | Undo works; first basic editing gate green (shortcuts and wrap mostly fixed on the way up to 90) |
| 16 Jul 2026, 10:01 | 105/72/33 | 36/26/10 | Full 105 / more realistic set; harder cases dropped the score |
| 16 Jul 2026, 12:41 | 105/91/14 | 36/36/0 | Shift-select fixed; basic editing gate fully green |
| 16 Jul 2026, 18:57 | 105/93/12 | 36/36/0 | Editor lives in our own fork now |
| 16 Jul 2026, 21:21 | 105/105/0 | 36/36/0 | Sign-off — all tests pass |
| 16 Jul 2026, 23:12 | 107/107/0 | 36/36/0 | Mid-sentence Shift+Up/Down across wrapping paragraphs fixed |
| 17 Jul 2026, 00:29 | 110/110/0 | 38/38/0 | Hardened: mid-wrapping Shift now critical; cross-para + mid-short-line select |

Other landmarks: shortcut-chord tag **25/25**; wrap **15/15**; undo **5/5**. Fork migration finished the same day as the first sign-off.
