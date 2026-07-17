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
| 17 Jul 2026, 01:32 | — | 38/38/0 | Phase 0 EditHelper stub (no behavior change); fork `9320231`; Patch LOC **36** |
| 17 Jul 2026, 01:57 | — | 38/38/0 | Phase A1 pure text math → EditHelper; wrap **15/15/0** @ `01-59-37`; fork `aa9912b`; Patch LOC **175** |
| 17 Jul 2026, 03:50 | — | 38/38/0 | Phase A2 undo/redo → EditHelper; undo **5/5/0** @ `03-47-34`; fork `a92ad2b`; Patch LOC **168** |
| 17 Jul 2026, 10:12 | 110/110/0 | 38/38/0 | Phase A3 sign-off (migration 2 Phase A); fork `a92ad2b`; edit-session PASS @ `10-08-06`; first run 109/110 (shift-select flake @ `10-09-05`), retry green |
| 17 Jul 2026, 10:29 | 110/110/0 | 38/38/0 | Phase B key-chord dispatcher → EditHelper; fork `57bfc21`; Patch LOC **279**; edit-session PASS @ `10-29-19`; full green first run @ `10-29-42` |
| 17 Jul 2026, 11:27 | — | — | Phase C visual-line math → EditHelper; fork `b3e2fe0`; Patch LOC **181**; CI build green (run `29569954349`); deploy pending until tablet online |
| 17 Jul 2026, 14:46 | 110/71/39 | — | Phase C first device run @ `192.168.86.30` (fork `b3e2fe0`); Up/Down jumped to doc bounds — `queryRectAt` invokeMethod used wrong arg types; fix pushed as `6a15e08` |
| 17 Jul 2026, 14:52 | 110/110/0 | 38/38/0 | Phase C sign-off; fork `6a15e08`; layout access fix; edit-session PASS @ `14-51-45`; full green @ `14-52-09`; critical **38/38/0** @ `14-55-48`; Patch LOC **177** |

Other landmarks: shortcut-chord tag **25/25**; wrap **15/15**; undo **5/5**. Fork migration finished the same day as the first sign-off.
