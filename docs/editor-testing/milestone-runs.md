# Pass/fail log of full typing-test runs

Log for `bash scripts/test-keyboard-harness.sh --fast`. Add a row after every full run. “Basic” is the smaller editing-works set. Calling typing work done needs the full set all passing.

Rows below keep every run where the **number of tests changed** (full set or basic set), plus later full-green checks at 110+.

| Time of test | All tests (total/passed/failed) | Basic tests (total/passed/failed) | Comments |
|--------------|---------------------------------|--------------------------------------|----------|
| 14 Jul 2026, 23:24 | 62/37/25 | — | First test |
| 15 Jul 2026, 04:45 | 83/64/16 | — | |
| 15 Jul 2026, 09:47 | 90/68/21 | — | |
| 15 Jul 2026, 23:53 | 90/80/9 | 34/34/0 | Basic set introduced. First pass for all basic tests |
| 16 Jul 2026, 00:37 | 94/89/4 | — | |
| 16 Jul 2026, 01:54 | 102/85/17 | — | |
| 16 Jul 2026, 10:01 | 105/72/33 | 36/26/10 | |
| 16 Jul 2026, 21:21 | 105/105/0 | 36/36/0 | First full pass |
| 16 Jul 2026, 23:12 | 107/107/0 | 36/36/0 | |
| 17 Jul 2026, 00:29 | 110/110/0 | 38/38/0 | |
| 17 Jul 2026, 17:23 | 110/110/0 | 38/38/0 | After QML -> c++ migration |
| 18 Jul 2026, 15:40 | 112/112/0 | 40/40/0 | In-editor Ctrl+C/X/V (fork `df1d38b`) |
| 18 Jul 2026, 16:53 | 115/115/0 | 40/40/0 | Mac Cmd+Left/Right line; Option+Up/Down paragraph; double-blank fixture; page buttons + rotation (fork `fa205c2`) |
| 18 Jul 2026, 17:17 | 116/116/0 | 40/40/0 | Shift+Alt stale head after typing (fork `7f1cf36`) |
| 18 Jul 2026, 17:46 | 118/118/0 | 40/40/0 | Viewport page step (landscape shorter); Shift+Left after type (fork `84e6bf0`) |
| 18 Jul 2026, 19:39 | 122/122/0 | 42/42/0 | Soft-wrap End/Cmd+Left/Right stay; Cmd+Backspace visual row; Option+Up/Down blank stops (fork `2ca3e2e`) |

Wrap group **15/15**; undo group **5/5**.
