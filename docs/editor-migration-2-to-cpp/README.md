# Editor migration 2 — helpers into C++

Pure edit math, chords, wrap walk, and undo moved into C++ EditHelper. QML still draws the text box and applies results. Done.

Keep calibrated wrap gaps and custom undo. Do not fork Qt’s text box. [../decisions.md](../decisions.md) §5–§6.

Checklist: [todo-handoff-edit-helper-cpp.md](todo-handoff-edit-helper-cpp.md). Prior migration: [../editor-migration-1-to-QML/](../editor-migration-1-to-QML/).
