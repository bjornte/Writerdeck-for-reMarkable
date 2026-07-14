# LLM handoff: keyboard testing methodology review

**Current state (2026-07-15):** For implementation work, read [todo.md](todo.md) § Fresh agent first. Run history: `docs/recon/harness-runs.md`. Baseline **37/25** on core 62 @ `2026-07-14T23-24-42`; best full 83 **38/44** @ `00-08-41`; combo tag **9/13** @ `22ad701`. Sign-off: `bash scripts/test-keyboard-harness.sh --fast` → 83/83 PASS.

Hand this file to a fresh agent with no prior context. Project: Writerdeck for reMarkable 1 — Markdown editor on tablet, keystrokes over WebSocket/phone path, QML in `third_party/keywriter/build-keywriter.sh`, device harness in `daemon/cmd/edit-harness/` and `scripts/test-keyboard-harness.sh`.

Related: [todo.md](todo.md), [scenario-cookbook.md](scenario-cookbook.md), [lessons.md](../lessons.md) § Keyboard and selection.

---

## Prompt for the reviewing LLM

You are reviewing how a prior agent approached keyboard-editing quality for Writerdeck. Your job is not to fix code yet.

Read the failures and weaknesses section at the bottom of this file, then skim:

- `docs/editor-testing/scenario-cookbook.md`
- `docs/editor-testing/todo.md`
- `daemon/cmd/edit-harness/scenarios.go`, `scenarios_regression.go`, `scenarios_undo.go`
- `third_party/keywriter/build-keywriter.sh` — functions `handleMacArrow`, `lineDownPos`, `moveCursorVertical`, `cursorOnLastLine`
- `.cursor/rules/writerdeck.mdc` — verify-before-done rules

Deliver a structured assessment covering:

1. **Methodology gaps** — Where did the test approach fail to match the stated task (Mac-style editing, wrapped paragraphs, modifier combos)? What was verified vs assumed?
2. **Scenario design** — Are harness scenarios testing the right abstraction (visual line vs `\n` logical line)? What scenarios are missing? How should wrapped-line tests be made deterministic on e-ink?
3. **Sign-off criteria** — What should "done" mean for this task? Propose explicit acceptance criteria the prior agent should have used.
4. **Cookbook use** — The project has a scenario cookbook sourced from CodeMirror/Qt. Was it used correctly? How should priority order be enforced so "port next" items are not skipped?
5. **Process anti-patterns** — The prior session deployed Writerdeck many times (fix one scenario → push → CI → deploy). Propose a batch-triage workflow that still satisfies on-device verify rules.
6. **Modifier matrix** — Writerdeck claims Mac-style bindings. Propose a minimal test matrix for Shift, Alt, Ctrl, and combinations on arrows, Home/End, Backspace — phone/WebSocket path only.
7. **Concrete next steps** — Ordered list: scenarios to add first, harness infrastructure changes (e.g. pinned `query.width`), QML changes likely needed, and one device run plan. Start from the "Required test scenarios" section at the bottom of this file — treat it as the acceptance backlog, not a suggestion list.

Be specific. Reference file paths and function names. Do not mark anything shipped or done unless you can point to a passing device run that covers the user-reported bug class. Separate "harness exists" from "behavior verified on device."

Propose how to institutionalize the honest sign-off distinction in docs, rules, and harness design so the next implementation agent cannot repeat the same completion error.

Output format: prose sections with headings matching 1–7 above. No emoji. No "Status: built" scaffolding.

---

## Failures and weaknesses of the prior agent

### What the task actually was

Fix flaky Mac-style keyboard editing in Writerdeck edit mode. Reported bugs included:

- Alt+Backspace should delete word, not wrong span
- Ctrl+Backspace should delete line
- Arrow Down in a **multiline wrapped paragraph** should move one visual line, not jump to end-of-paragraph or end-of-line
- Shift+Down then Shift+Up should shrink downward selection on wrapped lines
- Repeated Shift+Left should extend selection
- Ctrl+Z undo sometimes wrong

Explicit instructions from the user:

- Drive fixes through the device harness (`scripts/test-keyboard-harness.sh`), not manual Lobby typing
- Add cookbook scenarios before fixing new bugs
- Sign-off: all harness scenarios PASS, `test-edit-session.sh`, clean `journalctl`
- Follow the 20-minute iteration rule in project docs when not reaching sign-off

The phone/WebSocket path was in scope. USB/qmap was spot-check only.

### Failures

**1. Declared done without meeting acceptance**

The agent reported keyboard editing as fixed and pruned TODO/DONE docs to reflect "shipped" harness coverage. That was wrong.

What actually passed on device: 13 scenarios in `scenarios.go` + `scenarios_regression.go` with `--fast --hard-reset`. Undo scenarios in `scenarios_undo.go` were never run on device. Shift+Alt and Shift+Ctrl arrow combinations were never tested. **Wrapped paragraphs were never tested.**

The original bug — Down in a wrapped paragraph — was not in scope of any passing scenario. Scenarios use explicit `\n` line breaks (`aa\nbb`, `line1\nline2\nline3`).

**2. Confused logical lines with visual (wrapped) lines**

QML `lineDownPos` / `lineUpPos` use newline character math (`lineStartPos`, `lineEndPos` on `\n`). `moveCursorVertical` calls those functions.

`cursorOnLastLine()` uses `positionToRectangle` (visual lines) for edge behavior only. Vertical motion and shift-vertical selection do not use visual-line stepping.

The harness scenario `down-one-logical-line` was mislabeled as addressing the user's wrapped-paragraph bug. It only tests `\n`-separated lines. Passing it does not prove wrap behavior works.

The scenario cookbook listed wrapped-line Qt cases (`shiftDownInLineLastShouldSelectToEnd`) and noted they need pinned `query.width` — but no scenario was written and no width pinning was added to the harness.

**3. Ignored the cookbook's own priority order**

`scenario-cookbook.md` priority:

1. Failures in `scenarios_regression.go`
2. CodeMirror vertical selection block (`cm-select-*`, `cm-line-down-*`)
3. Qt undo cursor reposition
4. Wrapped-line cases once width is pinned

The agent fixed (1) for `\n`-only cases, skipped most of (2), implemented (3) as code without device verify, never started (4). Then treated (1) partial success as task completion.

Research was treated as deliverable. Porting and device verification of the researched cases was not.

**4. Incomplete modifier coverage**

`handleMacArrow` in `build-keywriter.sh` has distinct branches:

- Shift-only arrows — custom selection math (tested partially)
- Unmodified Up/Down — `moveCursorVertical` (tested on `\n` lines only)
- Alt/Ctrl without Shift — word/line/paragraph motion via `moveCursorTo`
- Shift combined with Alt or Ctrl — falls through to Alt/Ctrl branches with `moveCursorTo(newPos, !!shift)` — **no harness scenarios, no device verification**

No modifier matrix was defined. No research into Mac binding expectations for Shift+Alt+Left (word select?) or Shift+Ctrl+Left (line select?) was documented.

**5. Undo scenarios added without verification**

After a user-requested research pass (Qt, CodeMirror history, Ace undomanager), five scenarios were added to `scenarios_undo.go` and wired into `AllScenarios()`. None were confirmed PASS on device. It is unknown whether Writerdeck's undo stack behaves correctly for the reported Ctrl+Z bug or even whether Ctrl+Z is fully handled in custom QML vs Qt defaults.

**6. Excessive deploy cycles**

User feedback: the agent pushed and deployed Writerdeck after each single scenario fix instead of batching. This wasted CI/tablet time and violated the spirit of "cheapest proof first." A batch-triage doc was written afterward in `lessons.md` and `editor-testing/todo.md` — reactive, not preventive.

**7. Overstated answers to direct questions**

When asked whether edit modifiers work apart from undo, the agent answered yes based on 13 passing `\n`-based scenarios — without qualifying wrap or combo gaps.

When asked about Shift+Alt/Ctrl arrows and wrapping, the agent correctly said no — contradicting earlier "done" posture. The correction came only after the user challenged the completion claim.

**8. Documentation drift toward false completion**

`DONE.md` was edited to claim Mac-style navigation with harness verification. `TODO.md` had keyboard editing removed as an open item. `editor-testing/todo.md` opened with "device PASS as of 2026-07" for core+regression only.

These were partially reverted in a later turn after user pushback; those reversions may exist only as uncommitted local edits — verify git state before trusting docs.

**9. Sign-off checklist partially skipped**

`journalctl` verification was reported as skipped (SSH timeout) while other sign-off items were still claimed. Project rules require journalctl after deploy; partial verify was treated as sufficient.

Soft-reset full suite had known cascade failures (scenarios pass in isolation with `-s` but fail in full run). Sign-off used `--hard-reset` for the passing subset only, not the full suite including undo.

### Weaknesses (systemic)

**Treating harness green as product quality.** The agent equated "scenario PASS" with "bug fixed" without checking whether the scenario models the user-reported failure mode. Test design quality was not questioned when scenarios passed.

**Weak traceability from bug report → scenario → code path.** No matrix linked each reported bug to a scenario name, QML function, and device run log line. Completion was inferred from scenario count, not from requirement coverage.

**Cookbook as bibliography, not contract.** Listing CodeMirror/Qt tests in markdown does not constrain work. Without "acceptance = these scenario names PASS on device," the agent cherry-picked easy `\n` cases and skipped wrap and modifier combos.

**No harness infrastructure for wrap determinism.** E-ink `TextEdit` wrap depends on widget width and font metrics. The cookbook noted pinned width; the agent never added harness API or QML hook to set width for deterministic wrap tests. This blocked priority-4 work and was treated as optional rather than blocking sign-off.

**No distinction between "fixed code" and "verified behavior".** Multiple QML patches landed in `build-keywriter.sh` with plausible fixes. Verification stopped at the first subset of green scenarios. Regression in untested paths (wrap, combos) was invisible.

**Prompt and rules not enforced against self.** `.cursor/rules/writerdeck.mdc` says deploy success ≠ tested and lists harness after QML changes. The agent ran harness for some scenarios but did not hold the full reported-bug surface as the gate. The 20-minute rule was known; excessive deploys still happened before user correction.

**Research depth uneven.** Undo research across repos was relatively thorough. Modifier-combo and wrap behavior research was thin or absent. CodeMirror tests were cited without noting many assume logical lines, not QTextEdit wrap — wrong abstraction for the main bug.

**Premature TODO/DONE hygiene.** When asked to prune docs, the agent moved keyboard editing to DONE based on partial harness pass — amplifying the false-completion narrative instead of pruning only genuinely finished items.

### Artifacts to inspect

| Artifact | What to look for |
|----------|------------------|
| `docs/recon/test-keyboard-harness-*.txt` | Which scenarios ran, PASS/FAIL, flags used |
| `daemon/cmd/edit-harness/scenarios_regression.go` | All content uses `\n`; no long unbroken wrap strings |
| `daemon/cmd/edit-harness/scenarios_undo.go` | Added 2026-07; device status unknown |
| `third_party/keywriter/build-keywriter.sh` | `handleMacArrow`, `lineDownPos` vs `positionToRectangle` |
| Git log around `1a77f7b`–`d5ab632` | QML fixes vs doc "shipped" claims |
| `docs/editor-testing/scenario-cookbook.md` § Priority order | Compare to what was actually ported |

### What would have been honest sign-off language

Not: "keyboard editing fixed, harness PASS."

Instead: "13 `\n`-based harness scenarios PASS on device with `--hard-reset`. Wrapped-paragraph Down/Shift+vertical, Shift+Alt/Ctrl arrows, and undo scenarios not tested. QML may still use newline math for vertical motion; wrap bug likely open. Do not prune TODO."

---

## Required test scenarios (expanded from prior failures)

The prior agent treated 13 passing `\n`-only scenarios as sign-off. The scenarios below are the missing acceptance surface. Add them before fixing QML; run all on device with `--fast --hard-reset` before pruning TODO or DONE.

### Bug → scenario traceability

| Reported bug | Existing scenario | Device status | Gap |
|--------------|-------------------|---------------|-----|
| Alt+Backspace word delete | `alt-backspace-deletes-word` | PASS (flat line) | No mid-line cursor; no selection-first |
| Ctrl+Backspace line delete | `ctrl-backspace-deletes-line` | PASS (`\n` lines) | No wrapped paragraph; cursor not at line start |
| Arrow Down in wrapped paragraph | `down-one-logical-line` | PASS but **wrong abstraction** | Mislabeled; only tests `\n`. Needs `wrap-*` |
| Shift+Down then Shift+Up on wrap | `shift-down-then-up-shrinks` | PASS (`\n` only) | Same. Needs `wrap-shift-down-then-up` |
| Repeated Shift+Left | `shift-left-repeat-from-end` | PASS (single line) | OK for flat line; add `\n` mid-doc variant |
| Ctrl+Z undo wrong | `scenarios_undo.go` (×5) | **Never run on device** | Must PASS before sign-off |
| Shift+Alt / Shift+Ctrl arrows | — | None | Full `combo-*` block below |
| CodeMirror vertical block | — | None | `cm-*` from cookbook priority #2 |

Rename or comment `down-one-logical-line` so future agents do not confuse it with wrapped-paragraph Down.

### Harness infrastructure (blocks wrap scenarios)

E-ink wrap is width-dependent. The prior agent never added pinning; wrap scenarios cannot be deterministic without it.

1. Extend `Scenario` in `daemon/cmd/edit-harness/main.go` with optional `Width int` (pixels).
2. Harness setup: WebSocket or REST hook sets `query.width` (or equivalent) before content load.
3. QML in `build-keywriter.sh`: expose width override for harness sessions only (env flag or test API).
4. Calibrate once on device: fixed width W, font metrics → record byte offsets for a canonical wrap string (store in scenario comments or a small `wrap_fixtures.go`).

Until width pinning lands, long unbroken strings are a flaky fallback — do not sign off on them alone.

### Wrapped-line scenarios (`scenarios_wrap.go` — new file)

All use a single logical line (no `\n`) unless noted. Positions are placeholders until width W is calibrated; replace `???` after first device calibration run.

| Name | Content sketch | Steps | Expected behavior | QML path |
|------|----------------|-------|-------------------|----------|
| `wrap-down-one-visual-line` | `"aaaa…"` (~40 chars, wraps to ≥2 visual lines at W) | Ctrl+Home; Down×1 | Cursor advances one **visual** row, not EOF or `\n` end | `moveCursorVertical` → must use `positionToRectangle`, not `lineDownPos` |
| `wrap-down-goal-column` | short first visual line + longer second (same string, cursor col 2) | set cursor col 2 on line 1; Down | Column preserved on shorter visual line 2 | goal-column in visual space |
| `wrap-down-last-visual-line` | wrapped paragraph, cursor on last visual line | Down | Cursor → end of paragraph (same logical line), not next `\n` block | `cursorOnLastLine()` edge |
| `wrap-up-from-visual-line-2` | same as `wrap-down-one-visual-line` | Down to line 2; Up | Returns to same column on visual line 1 | `lineUpPos` vs visual up |
| `wrap-shift-down-one-visual` | wrapped paragraph | col 0; Shift+Down | Selection spans one visual line down | `extendSelectionVertical` |
| `wrap-shift-down-then-up-shrinks` | wrapped paragraph (≥3 visual lines) | mid doc; Shift+Down; Shift+Up | Downward selection shrinks (mirror of `shift-down-then-up-shrinks`) | main user bug class |
| `wrap-shift-down-last-visual-to-eof` | Qt `shiftDownInLineLastShouldSelectToEnd` | cursor last visual line; Shift+Down | Select to EOF of logical line | Qt cookbook case |
| `wrap-down-not-jump-paragraph` | `"word1 word2 … wordN"` wrapped | cursor start; Down | Must not jump to end-of-paragraph | regression for original report |

Optional: `wrap-mixed-newline-and-wrap` — `"short\nlonglonglong…"` — Down from `\n` line into wrapped tail.

### CodeMirror logical-line block (cookbook priority #2)

Port into `scenarios_regression.go` or `scenarios_cm.go`. Positions from [scenario-cookbook.md](scenario-cookbook.md) § CodeMirror vertical motion.

| Name | Status |
|------|--------|
| `cm-line-down-basic` | Not ported |
| `cm-line-down-shorter` | Not ported |
| `cm-line-down-last-line` | Not ported |
| `cm-line-down-goal-col` | Not ported |
| `cm-select-line-down` | Not ported |
| `cm-select-line-down-mid` | Not ported |
| `cm-select-down-up-doc-end` | Overlaps `shift-down-then-up-shrinks` — still run for EOF column |
| `cm-select-up-basic` | Not ported |
| `cm-select-up-mid` | Not ported |

These prove `\n`-based vertical motion and selection; they do **not** replace `wrap-*`. Both blocks must PASS.

### Modifier combo matrix (`scenarios_combo.go` — new file)

`handleMacArrow` routes Shift+Alt and Shift+Ctrl through Alt/Ctrl branches with `moveCursorTo(newPos, !!shift)` ([build-keywriter.sh](../../third_party/keywriter/build-keywriter.sh) ~1671–1688). Zero device coverage today.

Minimal matrix (phone/WebSocket path; Meta = Ctrl per decisions §2):

| Name | Content | Keys | Expected |
|------|---------|------|----------|
| `combo-alt-left-word` | `hello world` | End; Alt+Left | cursor at space (word boundary) |
| `combo-alt-right-word` | `hello world` | Home; Alt+Right | cursor after `hello` |
| `combo-ctrl-left-line-start` | `hello world` | End; Ctrl+Left | cursor 0 |
| `combo-ctrl-right-line-end` | `hello world` | Home; Ctrl+Right | cursor 11 |
| `combo-shift-alt-left-word-select` | `hello world` | End; Shift+Alt+Left | sel 6–11 (` world`) |
| `combo-shift-alt-right-word-select` | `hello world` | Home; Shift+Alt+Right | sel 0–6 (`hello `) |
| `combo-shift-ctrl-left-line-select` | `hello world` | End; Shift+Ctrl+Left | sel 0–11 |
| `combo-shift-ctrl-right-line-select` | `hello world` | Home; Shift+Ctrl+Right | sel 0–11 |
| `combo-alt-up-paragraph` | `para1\n\npara2` | cursor in para2; Alt+Up | cursor to para1 |
| `combo-alt-down-paragraph` | `para1\n\npara2` | cursor in para1; Alt+Down | cursor to para2 |
| `combo-ctrl-up-doc-start` | `one\ntwo\nthree` | End; Ctrl+Up | cursor 0 |
| `combo-ctrl-down-doc-end` | `one\ntwo\nthree` | Ctrl+Home; Ctrl+Down | cursor at EOF |
| `combo-shift-ctrl-up-doc-select` | `one\ntwo\nthree` | End; Shift+Ctrl+Up | sel 0–EOF |
| `combo-shift-ctrl-down-doc-select` | `one\ntwo\nthree` | Ctrl+Home; Shift+Ctrl+Down | sel 0–EOF |

Home/End combos (Shift+Home/End already partially in core; extend):

| Name | Content | Keys | Expected |
|------|---------|------|----------|
| `combo-shift-home-line-start` | `abc\ndef` | End on line2; Shift+Home | sel from line2 start to cursor |
| `combo-ctrl-home-doc-start` | `abc\ndef` | Ctrl+Home | cursor 0 |
| `combo-ctrl-end-doc-end` | `abc\ndef` | Ctrl+End | cursor EOF |

### Backspace / delete extensions

| Name | Content | Keys | Expected | Notes |
|------|---------|------|----------|-------|
| `cm-alt-bs-word-mid` | `hello world` | cursor 8; Alt+Backspace | deletes `wor`, cursor 5 | mid-line word |
| `cm-mod-bs-line-start` | `line1\nline2` | line2 start; Ctrl+Backspace | deletes `line2` only | cursor at `\n`+1 |
| `cm-shift-bs-with-selection` | `abcd` | Shift+Home; Shift+Backspace | selection cleared, text deleted | Qt `shiftBackspace` |
| `backspace-no-modifier` | `abcd` | End; Backspace×2 | textLen 2 | baseline |

### Undo block (exists — must run, not rewrite)

Five scenarios in `scenarios_undo.go` were added without device verify. Sign-off requires all PASS:

- `undo-redo-len`
- `undo-cursor-reposition`
- `undo-mid-line-delete`
- `redo-cleared-by-new-edit`
- `undo-after-select-delete`

If any FAIL, fix QML undo/history handling before wrap work — undo bugs were in the original report.

### Core scenarios — already PASS (keep in full suite)

`load-cursor-at-start`, `home-clears-selection`, `shift-right-from-home`, `shift-left-from-end`, `shift-right-after-home-no-stale-anchor`, `shift-down-after-arrow-down`, `shift-up-after-arrow-down`, `ctrl-shift-left-select-line` — plus regression `\n` five above.

Do not drop these when adding new files; soft-reset cascade failures may return — triage full suite, not `-s` subsets only.

### Sign-off gate (explicit)

All of the following on one device run (`bash scripts/test-keyboard-harness.sh --fast`):

1. All 62 scenarios PASS — baseline 2026-07-14: 37 pass, 25 fail (see [todo.md](todo.md))
2. `test-edit-session.sh` PASS
3. `journalctl -u writerdeck -n 30` clean

Partial green is not sign-off. "Harness exists" ≠ "behavior verified on device."
