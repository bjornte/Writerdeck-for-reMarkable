# Editor migration 2 — QML edit helpers → C++

Move pure text math and undo out of fork `edit_mac_helpers.qml.inc` into a small C++ `EditHelper` (same TextEdit on screen; QML stays the hands). **Phase A done** (fork `a92ad2b`; full **110/110/0** @ `10-12-39`). **Phase B done** (key-chord dispatcher; fork `57bfc21`; full **110/110/0** @ `10-29-42`). **Phase C done** (visual-line math; fork `6a15e08`; full **110/110/0** @ `14-52-09`). Post-port evaluation remains in the handoff § After A–C.

- [todo-handoff-edit-helper-cpp.md](todo-handoff-edit-helper-cpp.md) — checklist (Phases A–C done; § After A–C next)
- Standing Cursor rule: `.cursor/rules/writerdeck.mdc` (migration rule archived)
- Prior migration (done): [../editor-migration-1-to-QML/](../editor-migration-1-to-QML/)
- Policy: [../decisions.md](../decisions.md) §3 · How: [../architecture.md](../architecture.md) · Scores: [../editor-testing/](../editor-testing/)
- Fork source of truth: [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (`master`)

## Why

Migration 1 put edit behavior in owned QML. That worked and harness-signed. A large QML helper pile is harder to keep tidy than typed C++ with clear function boundaries. Phase A ports **behavior-identical** pure string math + undo; it must not change typing feel.

## Out of scope for Phase A

- Visual / wrap line motion — **Phase C done** (`6a15e08`; full **110/110/0** @ `14-52-09`); QML keeps `goalX` + apply
- Full key-chord dispatcher rewrite (`handleMacArrow` etc.) — Phase B done @ `57bfc21`
- Replacing Qt `TextEdit` or rewriting the text engine
- Lobby / save / sync / phone UI
- Cleaning wrap “magic” thresholds or redesigning custom undo — evaluate after A–C ([todo-handoff-edit-helper-cpp.md](todo-handoff-edit-helper-cpp.md) § After A–C; [../../TODO.md](../../TODO.md) item 7)
