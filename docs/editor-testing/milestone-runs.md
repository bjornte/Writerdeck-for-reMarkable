# Milestone full-suite runs

Scoreboard for `bash scripts/test-keyboard-harness.sh --fast`. Add a row after every full run. Critical is the smaller basic-editing gate. Sign-off needs the full suite green.

| Time of test | All tests (total/passed/failed) | Critical tests (total/passed/failed) | Comments |
|--------------|---------------------------------|--------------------------------------|----------|
| 16 Jul 2026, 21:21 | 105/105/0 | 36/36/0 | First full-suite sign-off |
| 16 Jul 2026, 23:12 | 107/107/0 | 36/36/0 | Mid-sentence wrap Shift fixed |
| 17 Jul 2026, 00:29 | 110/110/0 | 38/38/0 | Mid-wrapping Shift now critical |
| 17 Jul 2026, 10:12 | 110/110/0 | 38/38/0 | Phase A EditHelper math + undo; fork `a92ad2b` |
| 17 Jul 2026, 10:29 | 110/110/0 | 38/38/0 | Phase B chords; fork `57bfc21` |
| 17 Jul 2026, 14:52 | 110/110/0 | 38/38/0 | Phase C wrap walk; fork `6a15e08`; keep wrap gaps + custom undo (§5) |
| 17 Jul 2026, 17:23 | 110/110/0 | 38/38/0 | QML assembly in fork; tip `0bb3b70` |

Wrap tag **15/15**; undo tag **5/5**.
