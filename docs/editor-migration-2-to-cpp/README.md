# Editor migration 2 — QML edit helpers → C++

Move pure text math and undo out of fork `edit_mac_helpers.qml.inc` into a small C++ `EditHelper` (same TextEdit on screen; QML stays the hands). **In progress — Phase 0 done (fork `9320231`); next is Phase A1.**

- [todo-handoff-edit-helper-cpp.md](todo-handoff-edit-helper-cpp.md) — checklist (do the next unchecked item)
- Active Cursor rule: `.cursor/rules/edit-helper-cpp-migration.mdc` (general `writerdeck.mdc` paused)
- Prior migration (done): [../editor-migration-1-to-QML/](../editor-migration-1-to-QML/)
- Policy: [../decisions.md](../decisions.md) §3 · How: [../architecture.md](../architecture.md) · Scores: [../editor-testing/](../editor-testing/)
- Fork source of truth: [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (`master`)

## Why

Migration 1 put edit behavior in owned QML. That worked and harness-signed. A large QML helper pile is harder to keep tidy than typed C++ with clear function boundaries. Phase A ports **behavior-identical** pure string math + undo; it must not change typing feel.

## Out of scope for Phase A

- Visual / wrap line motion (`positionToRectangle`, `goalX`, `visualLine*`) — stays QML
- Full key-chord dispatcher rewrite (`handleMacArrow` etc.) — Phase B later
- Replacing Qt `TextEdit` or rewriting the text engine
- Lobby / save / sync / phone UI
- Cleaning wrap “magic” thresholds or redesigning custom undo — evaluate after A–C ([todo-handoff-edit-helper-cpp.md](todo-handoff-edit-helper-cpp.md) § After A–C; [../../TODO.md](../../TODO.md) item 7)
