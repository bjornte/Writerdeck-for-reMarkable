# TODO: Keyboard editing + harness

**Fresh-agent entry point.** Mac/Linux-style editing in Writerdeck (Ctrl/Alt chords — same on USB Linux keyboards and phone path; QML helpers in fork `edit_mac_helpers.qml.inc`). Drive fixes through the device harness — not manual Lobby typing.

Read first: this file, [milestone-runs.md](milestone-runs.md), [lessons.md](../lessons.md) § Keyboard and selection, [decisions.md](../decisions.md) §22. Scenario names: [scenario-catalog.md](scenario-catalog.md). Porting sources: [scenario-cookbook.md](scenario-cookbook.md).

Root pointer: [TODO.md](../../TODO.md) item 2.

## Current score (device)

| Milestone | Result | Note |
|-----------|--------|------|
| Latest full suite | **105 / 0** (0 prepare fail) of **105** | `21-21-15` @ fork `48b5d26`; report `docs/recon/test-keyboard-harness-2026-07-16T21-21-15.md` |
| Prior full suite | **93 / 12** (0 prepare fail) of **105** | `18-57-31` @ Phase 3 Lobby/shell |
| **Critical (gate)** | **36 / 36** | green @ `21-21-15` |
| Wrap tag | **15 / 15** | Phase 2B |
| Undo tag | **5 / 5** | Phase 2C @ `17-31-41` |
| Best pre-rewrite | **89 / 4** (+1 prep) of **94** | `00-37-27` @ `bdccee9` |
| Sign-off gate | **105/105 PASS** | met @ `21-21-15` |

`test-edit-session.sh` PASS on deploy @ fork `48b5d26`. Do not run it in parallel with the keyboard harness.

## Goal for next session

Keyboard harness sign-off is **105/105**, but one real-world gap is open below. Prefer that, or Physical Home owner check. Keep critical **36/36** green on every behavior-moving deploy. Edit QML/C++ in the fork, not in `build-keywriter.sh`.

## Open product gaps

- [ ] **Shift+Up/Down mid-sentence across wrapping paragraphs** — From the middle of a sentence in a long wrapping paragraph, Shift+Down (and Up) across several visual lines / into the next paragraph does not match Mac/usual editor behaviour. Existing harness covers Shift+vertical on short `\n` lines (`cm-select-*`, `shift-*-shrinks`) and on **one** wrapped block mostly from the left edge (`wrap-shift-down-*`). Nothing asserts mid-sentence selection walking through multi-paragraph prose. Add a scenario on `fixtureProse` (mid-caret → Shift+Down×N / Shift+Up×N), then fix in the fork.

## What `48b5d26` / `6a07cc6` fixed

- Ctrl+Left/Right → document bounds (matched Ctrl+Up/Down).
- Shift+Ctrl Right/Down via `shiftAnchor`/`shiftHead` (raw `query.select` collapsed on repeat).
- `goalXTrackSuspended` so keepGoalColumn survives short-line landings.
- Vertical shift snap only when the source line wraps visually; EOF forced select only on wrapped last lines.
- Escape auto-release blocked in `main.cpp` (double `toggleMode` cancelled Esc).
- Mid-word Alt+BS no longer eats the preceding space; harness Shift+Alt+Left expects fixed to 12-word line.
- `scrollDown` clamps at `contentHeight - height`; hw fixture tall enough for page-9 overshoot under clamp.

## Remaining fails

None @ `21-21-15`.

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

## Harness inventory (105)

Mode: **sandbox-prepare**. Tags: `-t critical`, `-t hw`, `-t read`, `-t wrap`, `-t undo`. Single scenario: `-s NAME --fast`. Step flag: `Reprepare` rewrites note + `harnessprepare` after mutating edits.

## Acceptance

1. `-t critical --fast` → **36/36 PASS** (met)
2. Full `--fast` → **105/105 PASS** (met @ `21-21-15`)
3. `test-edit-session.sh` PASS
4. `journalctl -u writerdeck -n 30` clean after deploy
