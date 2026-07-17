# Pass/fail log of full typing-test runs

Log for `bash scripts/test-keyboard-harness.sh --fast`. Add a row after every full run. “Basic” is the smaller editing-works set. Calling typing work done needs the full set all passing.

| Time of test | All tests (total/passed/failed) | Basic tests (total/passed/failed) | Comments |
|--------------|---------------------------------|--------------------------------------|----------|
| 16 Jul 2026, 21:21 | 105/105/0 | 36/36/0 | First full pass |
| 16 Jul 2026, 23:12 | 107/107/0 | 36/36/0 | Mid-sentence wrap Shift fixed |
| 17 Jul 2026, 00:29 | 110/110/0 | 38/38/0 | Mid-wrapping Shift now in the basic set |
| 17 Jul 2026, 10:12 | 110/110/0 | 38/38/0 | Phase A EditHelper math + undo; fork `a92ad2b` |
| 17 Jul 2026, 10:29 | 110/110/0 | 38/38/0 | Phase B shortcuts; fork `57bfc21` |
| 17 Jul 2026, 14:52 | 110/110/0 | 38/38/0 | Phase C wrap motion; fork `6a15e08`; keep wrap gaps + custom undo (§5) |
| 17 Jul 2026, 17:23 | 110/110/0 | 38/38/0 | Screen-file assembly in fork; commit `0bb3b70` |

Wrap group **15/15**; undo group **5/5**.
