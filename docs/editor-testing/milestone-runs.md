# Milestone full-suite runs

Scoreboard for full keyboard tests on the tablet (`bash scripts/test-keyboard-harness.sh --fast`). Add a row after every full run. Detailed logs: `docs/recon/test-keyboard-harness-*.{md,txt}` (older runs in [harness-runs.md](../recon/harness-runs.md)).

**Critical tests** = the smaller “basic editing” gate (38 today). Sign-off needs the full suite green: **110/110/0**.

| Time of test | All tests (total/passed/failed) | Critical tests (total/passed/failed) | Comments |
|--------------|---------------------------------|--------------------------------------|----------|
| 16 Jul 2026, 21:21 | 105/105/0 | 36/36/0 | First full-suite sign-off (fork migration day; climbed from early baselines) |
| 16 Jul 2026, 23:12 | 107/107/0 | 36/36/0 | Mid-sentence Shift+Up/Down across wrapping paragraphs fixed |
| 17 Jul 2026, 00:29 | 110/110/0 | 38/38/0 | Hardened: mid-wrapping Shift now critical; cross-para + mid-short-line select |
| 17 Jul 2026, 10:12 | 110/110/0 | 38/38/0 | Phase A3 sign-off; fork `a92ad2b`; Phases 0–A2 critical gates along the way (stub / pure math / undo); first run 109/110 flake (`10-09-05`), retry green |
| 17 Jul 2026, 10:29 | 110/110/0 | 38/38/0 | Phase B key-chord dispatcher → EditHelper; fork `57bfc21`; Patch LOC **279**; edit-session PASS @ `10-29-19`; full green first run @ `10-29-42` |
| 17 Jul 2026, 14:52 | 110/110/0 | 38/38/0 | Phase C sign-off; fork `6a15e08` after `b3e2fe0` device miss (invokeMethod arg types; caret width may be 0 — do not gate on `QRectF::isValid()`); first device run 71/110 @ `14-46-12`; edit-session PASS @ `14-51-45`; critical **38/38/0** @ `14-55-48`; Patch LOC **177**. After A–C evaluation (docs only): keep wrap gaps + custom EditHelper undo — [decisions.md](../decisions.md) §30 |

Other landmarks: shortcut-chord tag **25/25**; wrap **15/15**; undo **5/5**.
