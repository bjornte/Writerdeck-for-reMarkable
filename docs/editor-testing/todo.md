# TODO: Keyboard editing + harness

**Fresh-agent entry point.** Mac/Linux-style editing in Writerdeck (Ctrl/Alt chords — same on USB Linux keyboards and phone path; QML helpers in fork `edit_mac_helpers.qml.inc`). Drive fixes through the device harness — not manual Lobby typing.

Read first: this file, [milestone-runs.md](milestone-runs.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22. Scenario names: [scenario-catalog.md](scenario-catalog.md). Porting sources: [scenario-cookbook.md](scenario-cookbook.md).

Root pointer: [TODO.md](../../TODO.md) item 2.

## Current score (device)

| Milestone | Result | Note |
|-----------|--------|------|
| Latest full suite | **107 / 0** of **107** | `23-12-40` @ fork `67656e1`; report `docs/recon/test-keyboard-harness-2026-07-16T23-12-40.md` |
| Prior full suite | **105 / 0** of **105** | `21-21-15` first sign-off |
| **Critical (gate)** | **36 / 36** | green @ `23-11-28` |
| Wrap tag | **15 / 15** | recalibrated after visual-line fix |
| Undo tag | **5 / 5** | Phase 2C |
| Sign-off gate | **107/107 PASS** | met @ `23-12-40` |

`test-edit-session.sh` PASS on deploy @ fork `67656e1`. Do not run it in parallel with the keyboard harness.

## Goal for next session

Keyboard harness sign-off is **done** (**107/107**), including mid-sentence Shift+vertical across wrapping paragraphs. Prefer Physical Home owner check. Keep critical **36/36** green on every behavior-moving deploy. Edit QML/C++ in the fork, not in `build-keywriter.sh`.

## What `67656e1` fixed

- Mid-sentence Shift+Down/Up jumped ~160 characters because `visualLineDownPos` stepped by a too-tall character box on wrapped lines.
- Fix: walk to the next distinct row `y` with a small minGap; wrap fixture offsets recalibrated (~10 chars/row at W=320).
- New scenarios: `gap-shift-down-mid-wrapping-paras`, `gap-shift-up-mid-wrapping-paras` (uni1/uni5 + bi1+1/bi3+5/bi7+7).

## Remaining fails

None @ `23-12-40`.

## Next (one batch)

1. Owner: press physical Home from edit, read, and Lobby ([todo-handoff-physical-home-input.md](../todo-handoff-physical-home-input.md)).
2. Keep harness green on any future edit QML change: one push → CI → fetch → deploy → `test-edit-session.sh` → full `--fast` → update [milestone-runs.md](milestone-runs.md).

Deploy budget: **one** Writerdeck binary deploy per session unless QML fails to launch.

## Do not retry

- Inferring the moving selection end from `query.cursorPosition` after `query.select(min, max)`.
- Treating keyboard Left/Right as page-scroll.
- Per-scenario deploy loops.
- Parallel `test-edit-session.sh` + full harness.
- Auto-sending Qt KeyRelease for Escape in `rmkbdInjectLine` (double-fires mode toggle).
- Stepping visual lines by full `positionToRectangle(pos).height` on wrapped mid-line carets.

## Harness inventory (107)

Mode: **sandbox-prepare**. Tags: `-t critical`, `-t hw`, `-t read`, `-t wrap`, `-t undo`. Single scenario: `-s NAME --fast`. Step flag: `Reprepare` rewrites note + `harnessprepare` after mutating edits.

## Acceptance

1. `-t critical --fast` → **36/36 PASS** (met)
2. Full `--fast` → **107/107 PASS** (met @ `23-12-40`)
3. `test-edit-session.sh` PASS
4. `journalctl -u writerdeck -n 30` clean after deploy
