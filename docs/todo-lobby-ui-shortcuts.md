# TODO: Lobby shortcuts in lobby-ui.json

Move the remaining configurable Lobby chords onto disk so SSH edits apply without rebuilding Writerdeck. Policy: [decisions.md](decisions.md) §36–§37. Defaults: [../config/lobby-ui.json](../config/lobby-ui.json). Parent: [../TODO.md](../TODO.md).

## Stay in code (do not put in JSON)

Enter / Return (edit, sync, confirms, submit). Tab / Shift-Tab and plain Left / Right (Lobby page cycle). Up / Down and Page Up / Page Down (Files list). Home (edit → Files; Lobby → quit). Esc (dismiss). Backspace / Left / Right / Home / End inside New / Rename fields. Digits 0–9 on the private PIN pad.

## Checklist

- [ ] Remove the old Ctrl-K note picker (omni overlay). Lobby Files already opens and switches notes; the picker saves no time. Drop the keybinding, the overlay UI, and help text that still mentions Ctrl-K / quick file picker.
- [ ] Put page digits in JSON: which keys jump to which Lobby tab, and the digit badges on the tab bar. Default remains 1–6 lined up with tab order (Files … Home).
- [ ] Put rotate in JSON: Ctrl-R and Ctrl+Left / Ctrl+Right (today hardwired). Keep the Settings letter chord (`settings.rotation`, default O) as it already is. Document phone browsers often eating Ctrl-R (§37).
- [ ] Put quit in JSON: Ctrl-Q as an optional chord (Home from Lobby and Settings → Exit already quit). Drop the parser rule that forever blocks the letters K and Q from JSON.
- [ ] Optional: add `sync.now` to JSON if Sync should have a Ctrl-letter badge; today Sync is Enter-only and QML already looks up a missing key.
- [ ] Refresh `shortcuts.body`, Settings help lines, Lobby keyboard tests, and §36 so they match (no longer “Ctrl-K / Ctrl-Q stay only in code”).
- [ ] Verify on tablet: edit `lobby-ui.json` over SSH, confirm chords and badges update without a binary rebuild; bad JSON keeps the last good load.

## Out of scope here

Button labels, page copy, colors, and radii — [todo-lobby-ui-chrome.md](todo-lobby-ui-chrome.md).
