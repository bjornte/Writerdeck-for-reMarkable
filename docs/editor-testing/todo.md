# TODO: Keyboard editing + harness

**Fresh-agent entry point.** Mac/Linux-style editing in Writerdeck (Ctrl/Alt chords — same on USB Linux keyboards and phone path; QML helpers in fork `edit_mac_helpers.qml.inc`). Drive fixes through the device harness — not manual Lobby typing.

Read first: this file, [milestone-runs.md](milestone-runs.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22. Scenario names: [scenario-catalog.md](scenario-catalog.md). Porting sources: [scenario-cookbook.md](scenario-cookbook.md).

Root pointer: [TODO.md](../../TODO.md) item 2.

## Current score (device)

Scores are **total/passed/failed**.

| Milestone | Result | Note |
|-----------|--------|------|
| Latest full suite | **110/110/0** | `14-52-09` @ fork `6a15e08` (Phase C); report `docs/recon/test-keyboard-harness-2026-07-17T14-52-09.md` |
| Prior full suite | **110/110/0** | `10-29-42` @ fork `57bfc21` (Phase B) |
| **Critical (gate)** | **38/38/0** | green @ `14-55-48` (Phase C; fork `6a15e08`) |
| Wrap tag | **15/15/0** | mid-sentence wrap-shift shrinks strengthened |
| Undo tag | **5/5/0** | green @ `03-47-34` (Phase A2; stacks in `EditHelper`) |
| Sign-off gate | **110/110/0** | met @ `14-52-09` (Phase C) |

`test-edit-session.sh` PASS @ `14-51-45`. Do not run it in parallel with the keyboard harness.

## Goal for next session

Migration 2 Phase C is **done** (visual-line math in `EditHelper`; fork `6a15e08`; full **110/110/0** @ `14-52-09`). Next: [editor-migration-2-to-cpp](../editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md) § After A–C (evaluate wrap gaps / undo model). Keep critical green on every behavior-moving deploy. Edit QML/C++ in the fork, not in `build-keywriter.sh`.

## Hardening @ `00-29-12` (harness only)

- Mid-wrapping Shift scenarios promoted to **critical**.
- `wrap-shift-down-then-up-shrinks` and `cm-select-*-mid` now seed mid-sentence wrapping text.
- New: `gap-shift-down/up-across-para-break`, `gap-shift-down-mid-short-lines`.

## Remaining fails

None @ `10-12-39`.

## Next (one batch)

1. [editor-migration-2-to-cpp](../editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md) — Phase C (optional) when resuming migration 2; or owner physical Home check ([todo-handoff-physical-home-input.md](../todo-handoff-physical-home-input.md)).
2. Keep harness green on any future edit QML/C++ change: one push → CI → fetch → deploy → `test-edit-session.sh` → full `--fast` → update [milestone-runs.md](milestone-runs.md).

Deploy budget: **one** Writerdeck binary deploy per session unless QML fails to launch.

## Do not retry

- Inferring the moving selection end from `query.cursorPosition` after `query.select(min, max)`.
- Treating keyboard Left/Right as page-scroll.
- Per-scenario deploy loops.
- Parallel `test-edit-session.sh` + full harness.
- Auto-sending Qt KeyRelease for Escape in `rmkbdInjectLine` (double-fires mode toggle).
- Stepping visual lines by full `positionToRectangle(pos).height` on wrapped mid-line carets.
- Trusting short hard-newline Shift alone as coverage for mid-sentence wrapping selection.

## Harness inventory (110)

Mode: **sandbox-prepare**. Tags: `-t critical`, `-t hw`, `-t read`, `-t wrap`, `-t undo`. Single scenario: `-s NAME --fast`. Step flag: `Reprepare` rewrites note + `harnessprepare` after mutating edits.

## Acceptance

1. `-t critical --fast` → **38/38/0** (met)
2. Full `--fast` → **110/110/0** (met @ `00-29-12`)
3. `test-edit-session.sh` PASS
4. `journalctl -u writerdeck -n 30` clean after deploy
