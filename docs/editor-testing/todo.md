# TODO: Keyboard editing + harness

Mac-style editing in Writerdeck (`handleMacArrow`, `handleMacBackspace`, `handleMacEditKeys` in `build-keywriter.sh` patch 7n/7o). Drive fixes through the device harness — not manual Lobby typing.

Read: [scenario-cookbook.md](scenario-cookbook.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22.

## Fresh agent — start here

**Baseline (62 core):** `docs/recon/harness-runs.md` → **37 pass / 25 fail** @ `2026-07-14T23-24-42`, single session, no restarts. Latest combo tag: **9/13** @ `22ad701`.

**Suite size:** 83 scenarios (62 core + 15 gap + 6 extra wrap). Sign-off: **83/83 PASS** with `--fast`.

**Deploy:** `git push` → `fetch-keywriter-dist.sh` → `deploy-keywriter.sh -b`. Harness-only: local `go test`, no Writerdeck deploy unless `/api/test/*` changed (`deploy-rmkbd.sh`).

**One Writerdeck deploy per debugging session** unless QML fails to launch. Batch kernel fixes, one push/CI/fetch/deploy, then full harness.

### Failure clusters (baseline 25 on core 62)

| Cluster | Count | Likely fix |
|---------|------:|------------|
| Modified key from cursor 0 | ~15 | Harness End-prime poll + QML sandbox nav prime; socket route nav to focus item |
| CM selection / goal-col | 5 | Goal column + `extendSelectionVertical` |
| Undo | 3 | `handleMacUndo` cursor restore |
| Word nav off-by-one | 1 | `wordLeftPos` |
| Wrap up | 1 | `visualLineUpPos` / `lineUpPos` |

### Do not retry

- Separate WebSocket wake after prepare — does not unblock modified keys.
- End-only prime before every modified scenario — wiped docs.
- Duplicate `Keys.onPressed` on query TextEdit — QML crash loop.
- Routing only modified nav to focus item — plain Backspace/Delete dropped on window.
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
