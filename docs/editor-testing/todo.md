# TODO: Keyboard editing + harness

Mac-style editing in Writerdeck (`handleMacArrow`, `handleMacBackspace`, `handleMacEditKeys` in `build-keywriter.sh` patch 7n/7o). Drive fixes through the device harness — not manual Lobby typing.

Read: [scenario-cookbook.md](scenario-cookbook.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22.

## Fresh agent — start here

**Scores:** [milestone-runs.md](milestone-runs.md) (update after each full `--fast` run). Detail: [harness-runs.md](../recon/harness-runs.md).

| Milestone | Pass / fail | Commit / note |
|-----------|-------------|---------------|
| Core 62 baseline | 37 / 25 | `23-24-42` |
| Best full 83 | 64 / 16 | `04-45-43` (`071f998`) — stopped early |
| **Suite size** | **90** | +touch, +selection, +shift-alt repeat, +shift-ctrl multiline (not full-run yet) |

**Suite:** **90 scenarios**. Sign-off: **90/90 PASS** with `--fast`.

**Deploy:** `git push` → `fetch-keywriter-dist.sh` → `deploy-keywriter.sh -b`. Relaunch Writerdeck after binary deploy (`kill` + `POST /api/open`). Harness-only: `deploy-rmkbd.sh` when `/api/test/*` or harness runner changed. `systemctl restart writerdeck` if deploy-rmkbd stops the server.

**One Writerdeck deploy per debugging session** unless QML fails to launch. Batch kernel fixes, one push/CI/fetch/deploy, then full harness.

### Open failure clusters (local / unverified on device)

| Cluster | Scenario | Fix assumption |
|---------|----------|----------------|
| Shift+Ctrl+Left multiline | `combo-shift-ctrl-left-multiline` | Mac expects sel **0–caret**, not line start–caret. `socketRouteKey` already sets `lo=0` for Shift+Ctrl+Left; if harness still fails, check whether the key reaches `handleMacArrow` instead and uses `lineStartPos` via `extendSelectionHorizontal` without the Shift+Ctrl early block — ensure phone path hits `socketRouteKey` Shift+Ctrl branch; USB path must use `newPos=0` + `extendSelectionHorizontal`, not plain Ctrl+Left. |
| Shift+Alt arrow repeat | `combo-shift-alt-*-repeat` | Second press used anchor `pos` for word boundary. Fixed in tree: `selectionExtendFrom(key)` + `extendSelectionHorizontal` in `socketRouteKey` and `handleMacArrow`. Deploy + verify. |
| Visual goal-x / touch | `-t touch`, `-t cm`, `cm-line-down-goal-col` | Replaced character `goalColumn` with `goalX` from `positionToRectangle`; `lineDownPos`/`lineUpPos` always call `visualLine*Pos(pos, goalX)`. Tap/`harnesssetcursor` → `rememberGoalX`. `cm-line-down-goal-col` now expects cursor **6** (visual), not 11 (char column). Deploy + verify. |
| Shift+Left then Right | `shift-left-then-right-shrinks` | When head is at min end and user Shift+Right, shrink selection — head/anchor math in `handleMacArrow` shift+Right/Left blocks. Partial fix in tree; verify on device. |
| Undo cursor/text | `gap-undo-chain`, undo tag | `handleMacUndo` / `pendingRedoCursor` restore wrong caret after undo chain; triage undo tag before full suite (poisons late scenarios). |
| Wrap / CM selection | `wrap-*`, `cm-select-*`, `shift-down-then-up-shrinks` | Mostly visual-x + `extendSelectionVertical` anchor math; see [harness-runs.md](../recon/harness-runs.md). Wrap offsets calibrated @ W=320 in `wrap_fixtures.go` (20/40/24). |

Plain **Ctrl+nav from cursor 0** fixed @ `22ad701`. Combo tag **25/25** and wrap **17/17** @ `071f998` on the 83-scenario suite (pre–goal-x / pre–multiline scenario).

### Do not retry

- Separate WebSocket wake after prepare — does not unblock modified keys.
- End-prime before modified scenarios — wiped `query.text` / jumped to EOF.
- Nested `invokeMethod` for `socketRouteKey` on the GUI thread — deadlock (504 on editor-state).
- Duplicate `Keys.onPressed` on query TextEdit — QML crash loop.
- Routing only modified nav to focus item without `socketRouteKey` — Qt release poisoned buffer.
- Redundant `Ctrl+Home` in scenario setup when prime already positions.
- Per-scenario deploy loops — triage once, batch fix, one deploy ([lessons.md](../lessons.md) § Harness batch workflow).

### Harness inventory (90 scenarios)

| File | Block |
|------|--------|
| `scenarios.go` | Core (8) |
| `scenarios_regression.go` | Regression `\n` (6) |
| `scenarios_cm.go` | CodeMirror vertical (9) |
| `scenarios_combo.go` | Alt/Ctrl / Shift combos (25) |
| `scenarios_bs.go` | Backspace extensions (4) |
| `scenarios_wrap.go` | Wrapped paragraph (14) |
| `scenarios_undo.go` | Undo/redo (5) |
| `scenarios_gaps.go` | Gap coverage (15) |
| `scenarios_touch.go` | Touch → visual goal-x (3) |
| `scenarios_selection.go` | Shift reverse (1) |
| `main.go`, `report.go` | Runner, markdown reports, `showlobby` teardown |
| `wrap_fixtures.go` | Calibrated wrap offsets (W=320) |

Mode: **sandbox-prepare** — no editor quit between scenarios. Harness returns to Lobby via `showlobby` when done.

## Acceptance

Full suite **90/90 PASS** with `--fast`, single session, clean `journalctl`. `test-edit-session.sh` PASS.

## Dev loop

[lessons.md](../lessons.md) § Harness batch workflow. Triage once, batch fixes, one Writerdeck deploy, rerun full suite.

Per-scenario: `bash scripts/test-keyboard-harness.sh -s NAME --fast`. Match: `-m combo`, `-t wrap`, `-t touch`.
