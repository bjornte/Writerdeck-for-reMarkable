# TODO: Keyboard editing + harness

**Fresh-agent entry point.** Mac/Linux-style editing in Writerdeck (Ctrl/Alt chords — same on USB Linux keyboards and phone path; QML helpers in fork `edit_mac_helpers.qml.inc`). Drive fixes through the device harness — not manual Lobby typing.

Read first: this file, [milestone-runs.md](milestone-runs.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22. Scenario names: [scenario-catalog.md](scenario-catalog.md). Porting sources: [scenario-cookbook.md](scenario-cookbook.md).

Root pointer: [TODO.md](../../TODO.md) item 2.

## Current score (device)

Scores are **total/passed/failed**.

| Milestone | Result | Note |
|-----------|--------|------|
| Latest full suite | **110/110/0** | `00-29-12` @ fork `67656e1`; report `docs/recon/test-keyboard-harness-2026-07-17T00-29-12.md` |
| Prior full suite | **107/107/0** | `23-12-40` mid-wrapping Shift fix |
| **Critical (gate)** | **38/38/0** | green @ `00-24-22` (includes mid-wrapping Shift) |
| Wrap tag | **15/15/0** | mid-sentence wrap-shift shrinks strengthened |
| Undo tag | **5/5/0** | Phase 2C |
| Sign-off gate | **110/110/0** | met @ `00-29-12` |

`test-edit-session.sh` PASS on restart @ `00-16-42`. Do not run it in parallel with the keyboard harness.

## Goal for next session

Keyboard harness sign-off is **done** (**110/110/0**). Prefer Physical Home owner check. Keep critical **38/38/0** green on every behavior-moving deploy. Edit QML/C++ in the fork, not in `build-keywriter.sh`.

## Hardening @ `00-29-12` (harness only)

- Mid-wrapping Shift scenarios promoted to **critical**.
- `wrap-shift-down-then-up-shrinks` and `cm-select-*-mid` now seed mid-sentence wrapping text.
- New: `gap-shift-down/up-across-para-break`, `gap-shift-down-mid-short-lines`.

## Remaining fails

None @ `00-29-12`.

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
- Trusting short hard-newline Shift alone as coverage for mid-sentence wrapping selection.

## Harness inventory (110)

Mode: **sandbox-prepare**. Tags: `-t critical`, `-t hw`, `-t read`, `-t wrap`, `-t undo`. Single scenario: `-s NAME --fast`. Step flag: `Reprepare` rewrites note + `harnessprepare` after mutating edits.

## Acceptance

1. `-t critical --fast` → **38/38/0** (met)
2. Full `--fast` → **110/110/0** (met @ `00-29-12`)
3. `test-edit-session.sh` PASS
4. `journalctl -u writerdeck -n 30` clean after deploy
