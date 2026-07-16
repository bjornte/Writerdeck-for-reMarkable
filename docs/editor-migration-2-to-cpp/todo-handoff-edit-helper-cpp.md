# Handoff: Edit helpers QML → C++ (Phase A)

**Active.** Do the next unchecked item below. When a slice lands, check it off, record fork SHA + harness scores, and update [../editor-testing/milestone-runs.md](../editor-testing/milestone-runs.md). Active rule: `.cursor/rules/edit-helper-cpp-migration.mdc`. Paused: `.cursor/rules/writerdeck.mdc`. Policy: [../decisions.md](../decisions.md) §3. Root queue: [../../TODO.md](../../TODO.md) items 4–6 (A next; B/C later). Live keyboard scores: [../editor-testing/todo.md](../editor-testing/todo.md) (**110/110/0**, critical **38/38/0** @ `00-29-12`, fork `67656e1`).

Prior migration (done): [../editor-migration-1-to-QML/todo-handoff-keywriter-fork.md](../editor-migration-1-to-QML/todo-handoff-keywriter-fork.md). Folder overview: [README.md](README.md).

## Goal

Same typing behavior, clearer brain. Port **pure text math** and **undo/redo** from fork `edit_mac_helpers.qml.inc` into C++ `EditHelper` (QObject, `Q_INVOKABLE`), exposed to QML like `lobby_bridge`. QML still owns the on-screen `TextEdit`, Timers, Connections, harness hooks, visual wrap motion, and key dispatch (`handleMacKeysOnPressed` / `handleMacArrow` / …).

**Behavior lock:** after each slice, typing must match pre-slice harness results. Prefer bit-identical helper outputs for the same `(text, pos)` inputs.

## Where to edit

| Layer | Path |
|-------|------|
| Fork (source of truth) | Local clone or sibling `Writerdeck-keywriter` — [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) `master` |
| Helpers today | fork `edit_mac_helpers.qml.inc` |
| Existing C++ pattern | fork `lobby_bridge.{h,cpp}`, `edit_utils.h`, `main.cpp` context property, `edit.pro` |
| This repo build glue | `third_party/keywriter/build-keywriter.sh` — **assert only**, no new behavior |
| Verify | `test-edit-session.sh` then `test-keyboard-harness.sh` (never parallel) |

## Inventory (from `edit_mac_helpers.qml.inc`)

**Move in Phase A (pure string / undo — no layout):**
`isSpaceChar`, `lineStartPos`, `lineEndPos`, `lineCharCount`, `wordLeftPos`, `wordRightPos`, `deleteWordLeftPos`, `deleteLineLeftPos`, `paragraphUpPos`, `paragraphDownPos`, `insertTextDelta`, `isOneCharInsert`, undo stack ops (`captureEditState` data shape, `pushEditUndoWithMerge`, `editUndo` / `editRedo` merge rules, `clearEditUndoStacks`). Logical `lineDownPos` / `lineUpPos` only if they do **not** call `positionToRectangle`.

**Stay QML for Phase A:**
Anything using `query.positionToRectangle` / `goalX` / `visualLine*` / `lineWrapsVisually` / `onWrappedLine` / `macLineStartPos` / `macLineEndPos` when wrap-aware; `moveCursorTo`, `applyShiftSelection`, `extendSelection*`, `handleMacArrow`, `handleMacBackspace`, `handleMacEditKeys`, `handleMacKeysOnPressed`, `socketRouteKey`, harness helpers, Timers, Connections, `publishEditorState`.

**Phase B / C later (do not start until Phase A done):** key-chord dispatcher into C++; visual-line math into C++.

## Suggested C++ shape

- `edit_helper.h` / `edit_helper.cpp` — `class EditHelper : public QObject`
- Register in `main.cpp` as context property e.g. `editHelper` (alongside `writerdeck` / `EditUtils`)
- Add sources to `edit.pro`
- QML calls e.g. `editHelper.wordLeft(text, pos)` and keeps applying results to `query`
- Undo: either C++ owns stacks via methods, or C++ pure merge helpers while QML briefly still holds arrays — prefer C++-owned stacks once Phase A2 starts, with QML only calling `beginEdit` / `undo` / `redo` / `clear` and applying returned text+cursor+selection onto `query`

Match existing fork style: ASCII, Qt 5, `Q_INVOKABLE`, no new third-party deps.

## Verify (every behavior-moving deploy)

1. Push fork → this repo CI (or wait) → `bash scripts/fetch-keywriter-dist.sh` → `bash scripts/deploy-keywriter.sh -b`
2. Relaunch Writerdeck; `journalctl -u writerdeck -n 30` — fail on QML parse / instant exit
3. `bash scripts/test-edit-session.sh`
4. Then harness (not parallel): Phase 0/A1 → at least `-t critical --fast` (**38/38/0**); Phase A2 → also `-t undo --fast` (**5/5/0**); end of Phase A → full `--fast` (**110/110/0**)
5. Update [../editor-testing/milestone-runs.md](../editor-testing/milestone-runs.md) (include fork SHA + Patch LOC)

Deploy budget: **one** Writerdeck binary deploy per agent session unless the binary fails to launch.

## Slices

### Phase 0 — pin skeleton, no behavior change

- [ ] Inventory confirmed against current fork `edit_mac_helpers.qml.inc` (list above still accurate; note any drift in this file).
- [ ] Add `EditHelper` stub in the fork (`edit_helper.h` / `.cpp`), wire `edit.pro` + `main.cpp` context property; QML can see it but **no call sites** yet (or only a unused `Q_INVOKABLE` ping).
- [ ] CI build + `fetch-keywriter-dist.sh` + `deploy-keywriter.sh -b` + `test-edit-session.sh` + `-t critical --fast` → **38/38/0**. Record fork SHA here and in milestone-runs.

### Phase A1 — pure text math behind QML wrappers

- [ ] Port string helpers listed under Inventory → `EditHelper` (`Q_INVOKABLE` or statics called from invokables).
- [ ] Change QML functions of the same names to thin wrappers that call `editHelper.*` (keep names so `handleMacArrow` / backspace paths need minimal churn).
- [ ] Do **not** move visual-line or `positionToRectangle` callers.
- [ ] `test-edit-session.sh` + `-t critical --fast` → **38/38/0**. Prefer also `-t wrap --fast` (**15/15/0**) if any logical-line helper is shared with wrap paths. Record SHA + scores.

### Phase A2 — undo / redo into C++

- [ ] Port undo merge rules + stacks into `EditHelper` (or tightly coupled C++ types). QML Timers/Connections still call into it; `restore` still assigns `query.text` / cursor / selection in QML (or via a single apply helper in QML).
- [ ] Preserve one-char insert merge behavior (today’s `isOneCharInsert` / `pushEditUndoWithMerge`).
- [ ] `-t undo --fast` → **5/5/0**; `-t critical --fast` → **38/38/0**. Record SHA + scores.

### Phase A3 — sign-off + docs

- [ ] Full `bash scripts/test-keyboard-harness.sh --fast` → **110/110/0**.
- [ ] Shrink comments in `edit_mac_helpers.qml.inc` that still say “Phase 2C living only in QML” if misleading; note in [../architecture.md](../architecture.md) that pure math/undo live in fork `EditHelper`.
- [ ] Brief note in [../decisions.md](../decisions.md) §3 (or a short subsection) that Phase A of migration 2 shipped; link this handoff.
- [ ] Restore `.cursor/rules/writerdeck.mdc` (`alwaysApply: true`); set `edit-helper-cpp-migration.mdc` to `alwaysApply: false` with an archive note. Do this when parking after A3 **or** when B/C finish — keep the migration rule active while B/C are in progress.
- [ ] Update [../editor-testing/todo.md](../editor-testing/todo.md) / [../../TODO.md](../../TODO.md) items 4–6 when each phase completes (A → check off 4 and point Next at B; etc.).

### Phase B — key-chord dispatcher (later; do not start in the same session as A)

- [ ] Move chord → action mapping from `handleMacArrow` / `handleMacBackspace` / `handleMacEditKeys` into C++ **after** A3 is green. Separate handoff pass.

### Phase C — visual line (optional later)

- [ ] Only if A/B paid off. Needs careful layout access (`positionToRectangle` or equivalent). Not required for Phase A success.

## Do not

- Grow new editor behavior in `build-keywriter.sh`.
- Local `docker build` for Writerdeck — CI + `fetch-keywriter-dist.sh` only.
- Parallel `test-edit-session.sh` + full keyboard harness.
- `pkill -f /home/root/Writerdeck` (matches Writerdeck-server).
- Move visual-line / wrap math in Phase 0–A3.
- Replace `TextEdit` or invent a new buffer format.
- Change Markdown-on-disk / plain-text edit integrity ([../decisions.md](../decisions.md) § Document integrity, §26).
- One-scenario deploy loops; batch fixes, one binary deploy per session.
- Start Phase B/C before A3 sign-off.

## Resume prompt (copy for a fresh agent)

> Re. docs/editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md, do the next unchecked item. When done, update docs/editor-testing/milestone-runs.md.
