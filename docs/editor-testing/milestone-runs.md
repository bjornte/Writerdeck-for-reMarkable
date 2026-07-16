# Milestone full-suite runs

Hand-maintained scoreboard for full `bash scripts/test-keyboard-harness.sh --fast` sessions (no `-s`, no `--tag`). Update this table after every full-suite run — add a row with timestamp, total pass, critical pass/total, fail/prepare-fail, delta vs prior milestone, **Patch LOC** (`wc -l third_party/keywriter/build-keywriter.sh` on the binary under test), and one-line context (commit or what changed). Per-run reports: `docs/recon/test-keyboard-harness-*.{md,txt}` (new runs only; older files consolidated in [harness-runs.md](../recon/harness-runs.md)).

**Prepare fail** = scenarios that never ran because sandbox prepare failed (rewrite note / `harnessprepare` / verify), not keyboard assertion fails.

**Critical pass** = pass/total of the `-t critical` gate (36 scenarios today).

**Patch LOC** = line count of the emergency patch script `third_party/keywriter/build-keywriter.sh` at the commit under test (or nearest prior commit when the run had no SHA). Tracks patch-script debt shrinking as work moves into the keywriter fork.

Sign-off gate: **105/105 PASS** (pattern rewrite uni1/uni5/bi1+1/bi3+5/bi7+7; prior 102 baseline was 85/17). Critical gate: **36/36**.

| Run | Suite | Total pass | Critical pass | Fail | Prepare fail | vs prior | Patch LOC | Context |
|-----|-------|------------|----------------|------|--------------|----------|-----------|---------|
| 2026-07-14T20-15-40 | 62 | 26 | — | 36 | 0 | — | 2265 | early session |
| 2026-07-14T22-06-06 | 62 | 32 | — | 30 | 0 | +6 | 2318 | harness hardening |
| 2026-07-14T23-06-59 | 62 | 36 | — | 26 | 0 | +4 | 2349 | pre-baseline |
| 2026-07-14T23-24-42 | 62 | 37 | — | 25 | 0 | +1 | 2356 | **baseline anchor** |
| 2026-07-15T00-08-41 | 83 | 38 | — | 44 | 1 | new gaps | 2356 | **best 83** (pre-QML) |
| 2026-07-15T00-17-48 | 83 | 34 | — | 48 | 1 | −4 | 2428 | df2f850 QML deploy regression |
| 2026-07-15T00-43-13 | 83 | 27 | — | 53 | 1 | −7 | 2519 | 4c4d816 worst 83 |
| 2026-07-15T00-56-17 | 83 | 35 | — | 45 | 1 | +8 | 2515 | 0a339c9 partial recovery |
| 2026-07-15T02-03-16 | 83 | 46 | — | 34 | 1 | +11 | 2650 | 7d00156 selection-collapse fix; stopped 81/83 |
| 2026-07-15T02-07-09 | 83 | 52 | — | 28 | 1 | +6 | 2657 | 1e62aff direct query.select for shift+combo; stopped 81/83 |
| 2026-07-15T04-45-43 | 83 | 64 | — | 16 | 1 | +12 | 2658 | 071f998 combo positioning + wrap cal 20/40/24 + backspace guard; stopped 81/83 |
| 2026-07-15T09-47-15 | 90 | 68 | — | 21 | 1 | +4 | 2667 | f1ceaaa first full 90; combo/wrap 42/42; undo 0/7; touch 0/3; delete forward broken |
| 2026-07-15T10-07-59 | 90 | 73 | 31/34 | 16 | 1 | +5 | 2662 | 2ee3a92 goalX + delete cursor + shift+right repeat; undo only critical fail |
| 2026-07-15T23-53-00 | 90 | 80 | 34/34 | 9 | 1 | +7 | 2752 | 11625d7 custom edit undo stack; remaining: touch 0/3, plain L/R scroll, shift vertical shrink, goal-col |
| 2026-07-16T00-37-27 | 94 | 89 | — | 4 | 1 | +9 | 2783 | bdccee9 Left/Right caret + hw page cmds (contentY 0→1500→3000); touch 3/3; remaining: shift vertical shrink, goal-col, unicode prepare, alt-bs+sel |
| 2026-07-16T01-54-49 | 102 | 85 | — | 17 | 0 | first 102 | 2783 | harness rewrite baseline (prose + N=1/3/7 + read-overscroll); edit-session PASS; report `docs/recon/test-keyboard-harness-2026-07-16T01-54-49.md` |
| 2026-07-16T10-01-42 | 105 | 72 | 26/36 | 33 | 0 | first 105 | 2783 | pattern uni1/uni5/bi1+1/bi3+5/bi7+7 + multi-paragraph prose; edit-session PASS; report `docs/recon/test-keyboard-harness-2026-07-16T10-01-42.md` |
| 2026-07-16T11-38-40 | 105 | 74 | 26/36 | 31 | 0 | +2 | 2769 | `7603357` shift shrink/backspace batch; edit-session PASS; report `docs/recon/test-keyboard-harness-2026-07-16T11-38-40.md` |
| 2026-07-16T12-41-15 | 105 | 91 | 36/36 | 14 | 0 | +17 | 2766 | `377a053` persistent shiftAnchor/shiftHead + backspace Reprepare + wrap-up expects; edit-session PASS; report `docs/recon/test-keyboard-harness-2026-07-16T12-41-15.md` |
| 2026-07-16T14-29-52 | 105 | 92 | 36/36 | 13 | 0 | +1 | 1957 | Phase 2A: helpers in fork `edit_mac_helpers.qml.inc` (`568ee3f`); script loads file; critical 36/36; edit-session PASS; report `docs/recon/test-keyboard-harness-2026-07-16T14-29-52.md` |
| 2026-07-16T17-14-44 | 105 | 91 | 36/36 | 14 | 0 | −1 | 1957 | Phase 2B: wrap Shift+Down visual (`904ec77`); wrap tag 15/15; critical gate 36/36 @ `17-13-30`; edit-session PASS; report `docs/recon/test-keyboard-harness-2026-07-16T17-14-44.md` |

Combo-tag milestones (25 combo scenarios, `--tag combo`): 6/16 @ baseline → 9/13 @ `22ad701` → **25/25** @ `071f998` (`04-45-43`, 83-scenario suite). Wrap tag **15/15** @ Phase 2B (`17-14-44`). Full 94: **89/94** @ `bdccee9` (`00-37-27`). First full **102**: **85/17/0** @ `01-54-49`. First full **105**: **72/33/0** @ `10-01-42` @ `f42bfbe` (critical 26/36). Latest full **105**: **91/14/0** @ `17-14-44` @ Phase 2B (fork `904ec77`) (**critical 36/36** @ `17-13-30`; Patch LOC **1957**).
