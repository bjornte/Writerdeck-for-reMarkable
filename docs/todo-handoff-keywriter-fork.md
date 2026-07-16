# Handoff: Keywriter fork migration

Give this file to a fresh agent. One slice per session. Active rule: `.cursor/rules/keywriter-fork-migration.mdc`. Archived general rules: `.cursor/rules/writerdeck-backup.mdc`.

Policy: [decisions.md](decisions.md) §3. Root queue: [TODO.md](../TODO.md) item 3. Keyboard scores: [editor-testing/todo.md](editor-testing/todo.md).

## What we are leaving

`third_party/keywriter/build-keywriter.sh` rewrites upstream C++/QML with huge string patches every CI build. That is emergency architecture. **keywriter** (Qt 5 / C++ / QML) is the editor engine; **Writerdeck** is our on-device binary.

Build already supports a fork pin: `KEYWRITER_REPO` / `KEYWRITER_REF` in `build-keywriter.sh` (defaults: `dps/remarkable-keywriter` / `master`).

## Priority

Migrate **critical feature groups** into the fork first, in bulk. Do **not** first fix the 14 leftover non-critical harness fails. Only touch a fail when its feature group is being migrated. Harness **105/105** is product sign-off, not the migration order. Keep critical **36/36** green after every behavior-moving deploy.

Quality is the paramount driver for this migration. Check patterns from similar editors before inventing behavior. For the project as a whole, **document integrity** is absolute — [decisions.md](decisions.md) § Document integrity, [architecture.md](architecture.md), [integrity-audit.md](integrity-audit.md).

## Slices (check off in order)

### Phase 1 — pin, no behavior change

- [ ] Create Writerdeck-owned fork of [remarkable-keywriter](https://github.com/dps/remarkable-keywriter) (org/user repo the project controls).
- [ ] Wire CI / Dockerfile / `build-keywriter.sh` so builds clone that fork (`KEYWRITER_REPO` / `KEYWRITER_REF`), still applying today’s patch script unchanged.
- [ ] One CI build + `fetch-keywriter-dist.sh` + `deploy-keywriter.sh -b` + `test-edit-session.sh` + `-t critical --fast` → **36/36**. Document repo URL and default ref in [decisions.md](decisions.md) §3.

### Phase 2 — move critical groups into forked source

Do one lettered group per session (or per deploy cycle). After each group: remove the corresponding patches from `build-keywriter.sh`, rebuild, deploy, critical harness green.

- [ ] **A — Caret, shift selection, backspace/delete**  
  From script into fork: `handleMacArrow` horizontal/vertical shift (`shiftAnchor` / `shiftHead`), plain arrows, `handleMacBackspace` / word+line delete helpers. Critical scenarios that cover this group must stay green.

- [ ] **B — Wrap / visual line**  
  `visualLineUpPos` / `visualLineDownPos`, wrap Home/End, related goal-x. Move as one bulk; wrap harness tag as proof for this slice.

- [ ] **C — Undo / redo**  
  Custom edit undo stack and Ctrl+Z / Shift+Ctrl+Z routing. Undo scenarios as proof.

- [ ] **D — Combos / polish**  
  Alt/Ctrl motion, shift+alt/ctrl, remaining non-critical fails that belong here — only now.

### Phase 3 — shrink script + ownership

- [ ] `build-keywriter.sh` is build glue only (clone, qmake, install, tiny deterministic patches if any).
- [ ] Document fork ownership, default branch, and how to merge upstream keywriter in [decisions.md](decisions.md) §3.
- [ ] Restore `.cursor/rules/writerdeck-backup.mdc` → `writerdeck.mdc` (`alwaysApply: true`); retire or set `alwaysApply: false` on `keywriter-fork-migration.mdc`.

## Do not

- Grow new editor behavior in `build-keywriter.sh`.
- Prioritize the leftover 14 harness fails ahead of groups A–C.
- Local `docker build` for Writerdeck — CI + `fetch-keywriter-dist.sh` only.
- Parallel `test-edit-session.sh` + full keyboard harness.
- Declare the migration done while critical edit paths still live only as script string patches.

## Fresh session prompt

> Read `docs/todo-handoff-keywriter-fork.md` and follow `.cursor/rules/keywriter-fork-migration.mdc`. Do the next unchecked slice only. Criticality-first bulk migration into the keywriter fork — not leftover harness fail cleanup. After deploy: edit-session + `-t critical --fast` (36/36). Update this handoff checklist when the slice ships.
