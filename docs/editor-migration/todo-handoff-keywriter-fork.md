# Handoff: Keywriter fork migration

Give this file to a fresh agent. One slice per session. Active rule: `.cursor/rules/keywriter-fork-migration.mdc`. Archived general rules: `.cursor/rules/writerdeck-backup.mdc`.

Policy: [../decisions.md](../decisions.md) §3. Root queue: [../../TODO.md](../../TODO.md) item 3. Keyboard scores: [../editor-testing/todo.md](../editor-testing/todo.md).

## What we are leaving

`third_party/keywriter/build-keywriter.sh` rewrites upstream C++/QML with huge string patches every CI build. That is emergency architecture. **keywriter** (Qt 5 / C++ / QML) is the editor engine; **Writerdeck** is our on-device binary. **QML** = screen and typing behavior; **C++** = startup, display, socket keys — see [../architecture.md](../architecture.md) § On the tablet.

CI pins to owned fork [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (`master`) via `KEYWRITER_REPO` / `KEYWRITER_REF`. Edit helpers live in fork file `edit_mac_helpers.qml.inc` (includes cursor/autosave Timers and text-change Connections); socket/`lobby_bridge`/`rotation_watcher` are in-tree C++. `build-keywriter.sh` inserts the helpers file before `showLobby()`, injects a one-liner `handleMacKeysOnPressed` call, and still applies Lobby/shell QML Python patches.

## Priority

Migrate **critical feature groups** into the fork first, in bulk. Do **not** first fix the 14 leftover non-critical harness fails. Only touch a fail when its feature group is being migrated. Harness **105/105** is product sign-off, not the migration order. Keep critical **36/36** green after every behavior-moving deploy.

Quality is the paramount driver for this migration. Check patterns from similar editors before inventing behavior. For the project as a whole, **document integrity** is absolute — [../decisions.md](../decisions.md) § Document integrity, [../architecture.md](../architecture.md), [../integrity-audit.md](../integrity-audit.md).

## Slices (check off in order)

### Phase 1 — pin, no behavior change

- [x] Create Writerdeck-owned fork of [remarkable-keywriter](https://github.com/dps/remarkable-keywriter) — **done:** [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (default branch `master`).
- [x] Wire CI / Dockerfile / `build-keywriter.sh` so builds clone that fork (`KEYWRITER_REPO` / `KEYWRITER_REF`), still applying today’s patch script unchanged.
- [x] One CI build + `fetch-keywriter-dist.sh` + `deploy-keywriter.sh -b` + `test-edit-session.sh` + `-t critical --fast` → **36/36**. Documented in [../decisions.md](../decisions.md) §3.

### Phase 2 — move critical groups into forked source

Do one lettered group per session (or per deploy cycle). After each group: remove the corresponding patches from `build-keywriter.sh`, rebuild, deploy, critical harness green.

- [x] **A — Caret, shift selection, backspace/delete**  
  Helpers moved to fork [`edit_mac_helpers.qml.inc`](https://github.com/bjornte/Writerdeck-keywriter/blob/master/edit_mac_helpers.qml.inc) (`568ee3f`); script no longer embeds the string. Wrap/undo/combo **bodies** rode along (B–D finished Keys/props later). Critical **36/36**; full suite **92/13** @ `14-29-52`.

- [x] **B — Wrap / visual line**  
  Bodies already in fork. Fixed Shift+Down EOF jump on wrapped paragraphs (`904ec77` — snap only when crossing a newline). Wrap tag **15/15**; critical **36/36** @ `17-13-30`. Harness expect for full reverse shrink updated. No wrap-only scraps left in `build-keywriter.sh`. Full suite **91/14** @ `17-14-44`.

- [x] **C — Undo / redo**  
  Bodies already in fork. Undo property decls moved into [`edit_mac_helpers.qml.inc`](https://github.com/bjornte/Writerdeck-keywriter/blob/master/edit_mac_helpers.qml.inc) (`6676614`); script asserts presence. Undo tag **5/5**; critical **36/36** @ `17-34-55`; full suite **90/15** @ `17-31-53`.

- [x] **D — Combos / polish**  
  Keys wiring + remaining edit/cursor/harness props moved into [`edit_mac_helpers.qml.inc`](https://github.com/bjornte/Writerdeck-keywriter/blob/master/edit_mac_helpers.qml.inc) (`b0f17a5` — `handleMacKeysOnPressed`); script injects one call and asserts. Critical **36/36** @ `17-46-14`; full suite **93/12** @ `17-47-29`.

### Phase 3 — shrink script + ownership

- [x] **Connections / Timers** — cursor Timer, autosave Timer, and text-change Connections moved into [`edit_mac_helpers.qml.inc`](https://github.com/bjornte/Writerdeck-keywriter/blob/master/edit_mac_helpers.qml.inc) (`db0781e`); script asserts. Critical **36/36** @ `17-57-45`; full suite **92/13** @ `17-59-00` (Patch LOC **1802**). Critical edit paths no longer live only as script string patches.
- [x] **C++ infra** — socket reader, `lobby_bridge`, `rotation_watcher`, and toltec `edit.pro` live in the fork (`f7c84e9`); script asserts instead of `git apply` / COPY. Critical **36/36** + full suite — fill after deploy. Lobby/shell QML Python patches remain.
- [ ] `build-keywriter.sh` further shrink toward build glue only (Lobby/shell QML Python patches still dominate; tiny deterministic patches if any).
- [ ] Document fork ownership, default branch, and how to merge upstream keywriter in [../decisions.md](../decisions.md) §3.
- [ ] Restore `.cursor/rules/writerdeck-backup.mdc` → `writerdeck.mdc` (`alwaysApply: true`); retire or set `alwaysApply: false` on `keywriter-fork-migration.mdc`.

## Do not

- Grow new editor behavior in `build-keywriter.sh`.
- Prioritize the leftover 14 harness fails ahead of groups A–C.
- Local `docker build` for Writerdeck — CI + `fetch-keywriter-dist.sh` only.
- Parallel `test-edit-session.sh` + full keyboard harness.
- Declare the migration done while critical edit paths still live only as script string patches.

## Fresh session prompt

> Read `docs/editor-migration/todo-handoff-keywriter-fork.md` and follow `.cursor/rules/keywriter-fork-migration.mdc`. Do the next unchecked slice only. Criticality-first bulk migration into the keywriter fork — not leftover harness fail cleanup. After deploy: edit-session + `-t critical --fast` (36/36). Update this handoff checklist when the slice ships.
