# TODO: Keyboard editing + harness

Mac-style editing in Writerdeck (`handleMacArrow`, `handleMacBackspace`, `handleMacEditKeys` in `build-keywriter.sh` patch 7n/7o). Drive fixes through the device harness — not manual Lobby typing.

Read: [scenario-cookbook.md](scenario-cookbook.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22.

## Fresh agent — start here

**Scores:** [milestone-runs.md](milestone-runs.md) (update after each full `--fast` run). Detail: [harness-runs.md](../recon/harness-runs.md).

| Milestone | Pass / fail | Commit / note |
|-----------|-------------|---------------|
| Core 62 baseline | 37 / 25 | `23-24-42` |
| Best full 83 | 38 / 44 | `00-08-41` pre-QML |
| Latest full 83 | 35 / 45 | `00-56-17` (`0a339c9`) |
| Combo tag | 9 / 13 | `01-25-41` (`22ad701`) — no full 83 on that build yet |

**Suite:** 83 scenarios (62 core + 15 gap + 6 wrap). Sign-off: **83/83 PASS** with `--fast`.

**Deploy:** `git push` → `fetch-keywriter-dist.sh` → `deploy-keywriter.sh -b`. Harness-only: local `go test`, no Writerdeck deploy unless `/api/test/*` changed (`deploy-rmkbd.sh`). `systemctl restart writerdeck` if deploy-rmkbd stops the server.

**One Writerdeck deploy per debugging session** unless QML fails to launch. Batch kernel fixes, one push/CI/fetch/deploy, then full harness.

### Open failure clusters (post-`22ad701`)

| Cluster | Count | Likely fix |
|---------|------:|------------|
| Shift+Ctrl/Alt combos | ~10 | `socketRouteKey` + `extendSelectionHorizontal` for modified shift paths |
| Wrap visual-line down | ~11 | `visualLineDownPos` / harness width 320 calibration |
| CM goal-col / selection | 5 | Goal column + `extendSelectionVertical` |
| Undo cursor/text | 3–5 | `handleMacUndo` restore |
| Alt-only from cursor 0 | 2–3 | Alt fast-path in `socketRouteKey` (Ctrl path done @ `22ad701`) |
| Word nav off-by-one | 1 | `wordLeftPos` |

Plain **Ctrl+nav from cursor 0** fixed @ `22ad701` (`socketRouteKey` from inject thread, block Ctrl/Alt nav releases).

### Do not retry

- Separate WebSocket wake after prepare — does not unblock modified keys.
- End-prime before modified scenarios — wiped `query.text` / jumped to EOF.
- Nested `invokeMethod` for `socketRouteKey` on the GUI thread — deadlock (504 on editor-state).
- Duplicate `Keys.onPressed` on query TextEdit — QML crash loop.
- Routing only modified nav to focus item without `socketRouteKey` — Qt release poisoned buffer.
- Redundant `Ctrl+Home` in scenario setup when prime already positions.

### Harness inventory (83 scenarios)

| File | Block |
|------|--------|
| `scenarios.go` | Core (8) |
| `scenarios_regression.go` | Regression `\n` (5) |
| `scenarios_cm.go` | CodeMirror vertical (9) |
| `scenarios_combo.go` | Alt/Ctrl / Shift combos (22) |
| `scenarios_bs.go` | Backspace extensions (4) |
| `scenarios_wrap.go` | Wrapped paragraph (14) |
| `scenarios_undo.go` | Undo/redo (5) |
| `scenarios_gaps.go` | Gap coverage (15) |
| `main.go`, `report.go` | Runner, markdown reports |
| `wrap_fixtures.go` | Calibrated wrap offsets (W=320) |

Mode: **sandbox-prepare** — no editor quit between scenarios. `--hard-reset` removed.

## Acceptance

Full suite **83/83 PASS** with `--fast`, single session, clean `journalctl`. `test-edit-session.sh` PASS.

## Dev loop

[lessons.md](../lessons.md) § Harness batch workflow. Triage once, batch fixes, one Writerdeck deploy, rerun full suite.

Per-scenario: `bash scripts/test-keyboard-harness.sh -s NAME --fast`. Match: `-m combo`, `-t wrap`.
