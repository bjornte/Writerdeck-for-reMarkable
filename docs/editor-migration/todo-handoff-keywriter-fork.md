# Handoff: Keywriter fork migration

**Migration complete.** Active rule: `.cursor/rules/writerdeck.mdc`. Policy: [../decisions.md](../decisions.md) §3. Root queue: [../../TODO.md](../../TODO.md) item 3 (shipped). Live keyboard scores: [../editor-testing/todo.md](../editor-testing/todo.md) (**110/110/0**, critical **38/38/0**). Phase checkmarks below keep the scores from the day each slice landed (critical was still 36 then).

## What we left

`third_party/keywriter/build-keywriter.sh` used to rewrite upstream C++/QML with huge string patches every CI build. That emergency layer is retired for edit behavior. **keywriter** (Qt 5 / C++ / QML) is the editor engine; **Writerdeck** is our on-device binary. **QML** = screen and typing behavior; **C++** = startup, display, socket keys — see [../architecture.md](../architecture.md) § On the tablet.

CI pins to owned fork [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (`master`) via `KEYWRITER_REPO` / `KEYWRITER_REF`. Edit helpers live in fork file `edit_mac_helpers.qml.inc` (includes cursor/autosave Timers and text-change Connections); socket/`lobby_bridge`/`rotation_watcher` are in-tree C++; Lobby/shell QML and `lobby/*.inc` are in-tree too (`68f6e32`). `build-keywriter.sh` inserts helpers before `showLobby()`, concatenates Lobby subpages + sleep screen, and otherwise asserts + builds.

## Slices (all checked)

### Phase 1 — pin, no behavior change

- [x] Create Writerdeck-owned fork of [remarkable-keywriter](https://github.com/dps/remarkable-keywriter) — **done:** [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (default branch `master`).
- [x] Wire CI / Dockerfile / `build-keywriter.sh` so builds clone that fork (`KEYWRITER_REPO` / `KEYWRITER_REF`), still applying today’s patch script unchanged.
- [x] One CI build + `fetch-keywriter-dist.sh` + `deploy-keywriter.sh -b` + `test-edit-session.sh` + `-t critical --fast` → **36/36**. Documented in [../decisions.md](../decisions.md) §3.

### Phase 2 — move critical groups into forked source

- [x] **A — Caret, shift selection, backspace/delete** — helpers in fork `edit_mac_helpers.qml.inc` (`568ee3f`). Critical **36/36**; full suite **92/13** @ `14-29-52`.
- [x] **B — Wrap / visual line** — Shift+Down EOF fix (`904ec77`). Wrap tag **15/15**; critical **36/36**; full suite **91/14** @ `17-14-44`.
- [x] **C — Undo / redo** — undo props in helpers (`6676614`). Undo tag **5/5**; critical **36/36**; full suite **90/15** @ `17-31-53`.
- [x] **D — Combos / polish** — Keys wiring (`b0f17a5`). Critical **36/36**; full suite **93/12** @ `17-47-29`.

### Phase 3 — shrink script + ownership

- [x] **Connections / Timers** — fork `db0781e`; Patch LOC **1802**; full suite **92/13** @ `17-59-00`.
- [x] **C++ infra** — fork `f7c84e9`; Patch LOC **1778**; full suite **93/12** @ `18-10-10`.
- [x] Lobby/shell QML + `lobby/*.inc` in fork (`68f6e32`); Patch LOC **386**; critical **36/36** @ `18-56-12`; full suite **93/12** @ `18-57-31`.
- [x] Document fork ownership, default branch, and how to merge upstream keywriter in [../decisions.md](../decisions.md) §3.
- [x] Restore `.cursor/rules/writerdeck.mdc` (`alwaysApply: true`); retire `keywriter-fork-migration.mdc` (`alwaysApply: false`).

## Do not (still true)

- Grow new editor behavior in `build-keywriter.sh`.
- Local `docker build` for Writerdeck — CI + `fetch-keywriter-dist.sh` only.
- Parallel `test-edit-session.sh` + full keyboard harness.
- `pkill -f /home/root/Writerdeck` (matches Writerdeck-server).
