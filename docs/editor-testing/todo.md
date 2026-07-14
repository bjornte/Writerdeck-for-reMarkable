# TODO: Keyboard editing + harness

Mac-style editing in Writerdeck (`handleMacArrow`, `handleMacBackspace`, `handleMacUndo` in `build-keywriter.sh` patch 7n/7o). Drive fixes through the device harness — not manual Lobby typing.

Read: [scenario-cookbook.md](scenario-cookbook.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22.

## Fresh agent — start here (2026-07-14)

**Baseline:** `bash scripts/test-keyboard-harness.sh --fast` → **37 pass, 25 fail**, single editor session, no restarts. Latest report: `docs/recon/test-keyboard-harness-2026-07-14T23-24-42.md`.

**Session arc:** 26/36 → 32/30 → 36/26 → **37/25** (flat on last two runs). Do not re-try approaches that already failed (see § Do not retry).

**Deploy on tablet:** `git pull && bash scripts/fetch-keywriter-dist.sh && bash scripts/deploy-keywriter.sh -b && systemctl restart writerdeck`. Harness-only changes: `go run` locally, no Writerdeck deploy unless `/api/test/*` changed (`deploy-rmkbd.sh`).

**One Writerdeck deploy per debugging session** unless QML fails to launch. Batch QML fixes, then one push/CI/fetch/deploy, then full harness.

### Failure clusters (25)

| Cluster | Count | Examples | Likely fix |
|---------|------:|----------|------------|
| Modified key from cursor 0 | ~15 | `combo-*-right`, `combo-*-down`, `combo-ctrl-end` | Plain `End` on **same scenario WebSocket** before first modified step; or QML prime after `harnessSandboxReset` |
| CM selection / goal-col | 5 | `cm-line-down-goal-col`, `cm-select-*` | Goal column at EOF; `extendSelectionVertical` anchor math |
| Undo | 3 | `undo-redo-len`, `undo-cursor-reposition`, `undo-mid-line-delete` | `handleMacUndo` + Qt undo stack cursor restore over WebSocket |
| Word nav off-by-one | 1 | `combo-alt-left` (want 6 got 5) | `wordLeftPos` lands on space not word start |
| Wrap up | 1 | `wrap-up-from-visual-line-2` | `visualLineUpPos` / `lineUpPos` on wrapped paragraph |

Everything else in the suite passes, including core/regression, most wrap scenarios, backspace extensions, `combo-shift-end-line`, `combo-shift-home-line`.

### Do not retry

- Separate WebSocket wake after prepare (`ArrowLeft`, `ArrowUp`, `End`+`Home`) — does not unblock modified keys.
- `End`-only prime before every modified scenario — wiped docs (23/39 run).
- Duplicate `Keys.onPressed` on query TextEdit (patch 6c) — QML crash loop (`Property value set multiple times`). Mac keys already hook via patch **7o** on query's existing handler.
- Routing all modified nav to `activeFocusItem` in socket inject — no score change vs prior run (37/25 unchanged).
- Redundant `Ctrl+Home` before modified keys in scenarios — poisons next key when press-only; already removed from combo setup.

### Next work (ordered)

1. Harness: if first step key has modifiers, send plain `End` on the **same** `RunScenario` WebSocket, then `Ctrl+Home` or scenario-specific reset if cursor must stay 0 — prove on `combo-ctrl-right` before full deploy.
2. QML: `wordLeftPos` — return word start (6) not preceding space (5) for `hello world`.
3. QML: CM goal column + vertical selection (`cm-line-down-goal-col`, five `cm-select-*`).
4. QML: `visualLineUpPos` for `wrap-up-from-visual-line-2`.
5. QML: undo cursor/text restore for three undo scenarios.

## Harness inventory (62 scenarios)

| File | Block |
|------|--------|
| `scenarios.go` | Core (8) |
| `scenarios_regression.go` | Regression `\n` (5) |
| `scenarios_cm.go` | CodeMirror vertical (9) |
| `scenarios_combo.go` | Alt/Ctrl / Shift combos (22) |
| `scenarios_bs.go` | Backspace extensions (4) |
| `scenarios_wrap.go` | Wrapped paragraph, width 320 (8) |
| `scenarios_undo.go` | Undo/redo (5) |
| `main.go`, `report.go` | Runner, markdown reports |
| `scripts/test-keyboard-harness.sh` | Shell wrapper |

Mode: **sandbox-prepare** — `harnessopen` + `harnessprepare` (QML `harnessSandboxReset`), no editor quit between scenarios. `--hard-reset` removed.

Wrapped width: `harnessWrapWidth` in `wrap_fixtures.go`; QML `harnessSetWidth` / `harnessprepare` `w` field.

## Shipped this session (do not re-litigate)

- Sandbox prepare without editor restart; markdown reports in `docs/recon/`.
- Goal column, visual line down/up, Shift+Home/End routing, Alt+Backspace whole-word, EOF shift+down selection, modifier release (Shift-only not Ctrl+Shift).
- Home release in edit mode no longer calls `handleHome` → lobby (`combo-shift-end-line` fixed).
- `forceActiveFocus` + `ctrlPressed` reset in `harnessSandboxReset`.

Commits on `main`: `e84e1ad`, `d228218`, `089dfdb` (socket routing, no harness score gain).

## Acceptance

Full suite **62/62 PASS** with `--fast`, single session, clean `journalctl` (no QML parse errors, no restart loop). `test-edit-session.sh` PASS. Do not sign off on partial pass.

## Dev loop

[lessons.md](../lessons.md) § Harness batch workflow. Triage once (`--unit`, full `--fast`), batch fixes, one Writerdeck deploy, rerun full suite.

Per-scenario: `bash scripts/test-keyboard-harness.sh -s NAME --fast`. Match: `-m combo`.

## Resume prompt

> Keyboard harness: 62 scenarios, baseline **37 pass / 25 fail** (`docs/recon/test-keyboard-harness-2026-07-14T23-24-42.md`). Read this file § Fresh agent and § Do not retry. First fix: modified-key-from-zero cluster (~15 fails) via same-WS `End` prime or QML sandbox prime. One deploy per session. `bash scripts/test-keyboard-harness.sh --fast` for sign-off check.
