# TODO: Keyboard editing + harness

**Fresh-agent entry point.** Mac/Linux-style editing in Writerdeck (Ctrl/Alt chords — same on USB Linux keyboards and phone path; QML helpers still named `handleMacArrow`, `handleMacBackspace`, `handleMacEditKeys` in `third_party/keywriter/build-keywriter.sh`). Drive fixes through the device harness — not manual Lobby typing.

Read first: this file, [milestone-runs.md](milestone-runs.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22. Scenario names: [scenario-catalog.md](scenario-catalog.md). Porting sources: [scenario-cookbook.md](scenario-cookbook.md).

Root pointer: [TODO.md](../../TODO.md) item 2.

## Current score (device)

| Milestone | Result | Note |
|-----------|--------|------|
| Latest full suite | **91 / 14** (0 prepare fail) of **105** | `12-41-15` @ `377a053`; report `docs/recon/test-keyboard-harness-2026-07-16T12-41-15.md` |
| Prior full suite | **74 / 31** (0 prepare fail) of **105** | `11-38-40` @ `7603357` |
| **Critical (gate)** | **36 / 36** | green @ `377a053` |
| Best pre-rewrite | **89 / 4** (+1 prep) of **94** | `00-37-27` @ `bdccee9` |
| Sign-off gate | **105/105 PASS** | `bash scripts/test-keyboard-harness.sh --fast`, single session |

`test-edit-session.sh` PASS on deploy @ `377a053`. Do not run it in parallel with the keyboard harness.

## Goal for next session

Prefer the **keywriter fork** migration — [todo-handoff-keywriter-fork.md](../editor-migration/todo-handoff-keywriter-fork.md), rule `keywriter-fork-migration.mdc`. Order: critical feature groups in bulk (A→D). Do **not** first burn down the 14 leftover fails. Keep critical **36/36** green on every behavior-moving deploy.


## What `377a053` fixed

- Persistent `shiftAnchor` / `shiftHead` (Qt parks caret at `selectionEnd`, so reverse Shift was moving the wrong end).
- Mid-scenario `Reprepare` after mutating alt/ctrl-backspace uni1 (stale absolute `SetCursor` on a shrunken buffer).
- Wrap-up expects matched Down×7 then Up×3 geometry (~80, not ≤65).

## Remaining fails @ `12-41-15` (14)

| Scenario | Likely area |
|----------|-------------|
| `cm-line-down-goal-col` | goalX / shorter-line landing |
| `cm-select-line-down-mid` | vertical shift snap mid-line |
| `cm-select-down-up-doc-end` | EOF vertical selection |
| `combo-ctrl-left` / `combo-ctrl-right` | Ctrl line vs doc motion on prose |
| `combo-shift-alt-left` (+ repeat) | word-select head vs `shiftHead` |
| `combo-shift-ctrl-right` / `combo-shift-ctrl-down` | shift+ctrl extend |
| `bs-alt-word-mid` | mid-word Alt+BS (off by 1) |
| `wrap-shift-down-one-visual` | Shift+Down N=3 jumps to EOF |
| `gap-alt-bs-with-selection` | same cluster as shift-alt-left |
| `read-overscroll-clamps` | reading mode / Esc |
| `shift-left-then-right-shrinks` | selection alias; uni5 off-by-one |

## Next (one batch)

1. Triage the 14 with `-s NAME --fast` on the current binary (no deploy between).
2. Prefer harness/expect fixes where geometry is wrong; QML only where behavior is wrong.
3. One push → CI → fetch → deploy → `test-edit-session.sh` → full `--fast` → update [milestone-runs.md](milestone-runs.md).

Deploy budget: **one** Writerdeck binary deploy per session unless QML fails to launch.

## Do not retry

- Inferring the moving selection end from `query.cursorPosition` after `query.select(min, max)`.
- Treating keyboard Left/Right as page-scroll.
- Per-scenario deploy loops.
- Parallel `test-edit-session.sh` + full harness.
- Declaring done while full suite is under 105/105.

## Harness inventory (105)

Mode: **sandbox-prepare**. Tags: `-t critical`, `-t hw`, `-t read`, `-t wrap`, `-t undo`. Single scenario: `-s NAME --fast`. Step flag: `Reprepare` rewrites note + `harnessprepare` after mutating edits.

## Acceptance

1. `-t critical --fast` → **36/36 PASS** (met @ `377a053`)
2. Full `--fast` → **105/105 PASS**
3. `test-edit-session.sh` PASS
4. `journalctl -u writerdeck -n 30` clean after deploy
