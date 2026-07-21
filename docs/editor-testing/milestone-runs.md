# Pass/fail log of full typing-test runs

Log for `bash scripts/test-keyboard-harness.sh --fast`. Add a row after every full run. “Basic” is the smaller editing-works set. Calling typing work done needs the full set all passing.

Keep every run where the **number of tests changed**, plus later full-green checks. Older intermediate fails before the first full pass are omitted.

| Time of test | All tests (total/passed/failed) | Basic tests (total/passed/failed) | Comments |
|--------------|---------------------------------|--------------------------------------|----------|
| 15 Jul 2026, 23:53 | 90/80/9 | 34/34/0 | Basic set introduced |
| 16 Jul 2026, 21:21 | 105/105/0 | 36/36/0 | First full pass |
| 17 Jul 2026, 00:29 | 110/110/0 | 38/38/0 | |
| 17 Jul 2026, 17:23 | 110/110/0 | 38/38/0 | After EditHelper C++ migration |
| 18 Jul 2026, 15:40 | 112/112/0 | 40/40/0 | In-editor Ctrl+C/X/V (fork `df1d38b`) |
| 18 Jul 2026, 16:53 | 115/115/0 | 40/40/0 | Mac Cmd+Left/Right; Option+Up/Down (fork `fa205c2`) |
| 18 Jul 2026, 17:17 | 116/116/0 | 40/40/0 | Shift+Alt stale head (fork `7f1cf36`) |
| 18 Jul 2026, 17:46 | 118/118/0 | 40/40/0 | Viewport page step; Shift+Left after type (fork `84e6bf0`) |
| 18 Jul 2026, 19:39 | 122/122/0 | 42/42/0 | Soft-wrap End/Cmd line; Option blank stops (fork `2ca3e2e`) |
| 18 Jul 2026, 21:01 | 124/124/0 | 44/44/0 | Soft-wrap End affinity (fork `5d543f7`) |
| 18 Jul 2026, 21:32 | (partial) | 57/57/0 | Critical + wrap End/Home; End+Down EOF (fork `19792fc`) |
| 18 Jul 2026, 21:43 | 125/125/0 | 57/57/0 | Full green (fork `19792fc`) |

Wrap group **21/21**; undo group **5/5**.
