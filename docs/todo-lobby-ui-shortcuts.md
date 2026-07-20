# TODO: Lobby shortcuts in lobby-ui.json

Move the remaining configurable Lobby shortcuts onto disk so SSH edits apply without rebuilding Writerdeck. Policy: [decisions.md](decisions.md) §36–§37. Defaults: [../config/lobby-ui.json](../config/lobby-ui.json). Parent: [../TODO.md](../TODO.md). Chord letters and special values live in that JSON — not in decisions.

## Rules (agreed)

One shortcut per action — one line in the file, one live binding.

A letter value means Ctrl (or Cmd) plus that letter. Two special values do **not** need Ctrl: `enter` (Return / Enter key) and `hardware_home` (the reMarkable 1’s physical middle Home button — not the Home key on a USB or phone keyboard).

Whenever a shortcut value is `enter`, the on-screen key badge shows the single glyph **↩** — never the word “enter”.

**Tabs** are ordinary actions too — same JSON, same rules. Always write every tab line, with an empty value for “no shortcut”:

```json
"tabs.files": "",
"tabs.keyboard": "",
"tabs.sync": "",
"tabs.settings": "",
"tabs.shortcuts": "",
"tabs.about": ""
```

Empty means no jump key and no tab badge. **Remove** the hardwired digits 1–6 page jump from code (do not leave it as a fallback). Lobby page switching without a configured tab shortcut is Tab / Shift-Tab and plain Left / Right only. Later you may fill a tab line (e.g. `"tabs.files": "f"` → Ctrl-F). A future Finder-style note jump needs bare letters **and** digits, so do not spend bare digit keys on Lobby tabs without deciding that trade-off; fix §37’s old “digits for tabs / letters for Finder” split when refreshing docs.

**Clashes** — only within the same scope. Same letter on two Files actions (or two `global.*` lines): the later line in the file wins. Same letter on Files vs Settings is fine (page-scoped). Duplicate names (two `global.quit` lines) — only the last counts. If an LLM spots a same-scope clash while working, say so.

**Home: two lines** (do not fold into one):

```json
"global.toLobby": "hardware_home",
"global.quit": "hardware_home"
```

- `global.toLobby` — from a **note** (edit or read) back to Lobby.
- `global.quit` — from the **Lobby**, quit Writerdeck to the stock reMarkable UI.
- Default both to `"hardware_home"`. Settings → Exit still works.
- A letter on either line (e.g. `"q"`) means that Ctrl chord owns **that** action only; `hardware_home` no longer does that one. The other line stays independent.
- Never treat the USB / phone **keyboard Home** key as `hardware_home`. Keyboard Home keeps its normal jobs (caret Home/End in the editor and in New / Rename fields). It must not return to Lobby or quit Writerdeck.

**Rotate** — only the existing Settings letter (`settings.rotation`, default O → Ctrl-O). Drop Ctrl-R and Ctrl+Left / Ctrl+Right rotate.

**Sync now** — `"sync.now": "enter"` (Enter key; badge ↩).

**Edit selected note** — prefer `"files.edit": "enter"` (badge ↩ from the value) if wiring it like other JSON shortcuts is a small change. If that turns into a large detour, leave Enter hardwired for Edit and keep using the existing badge string until chrome work.

Drop the parser rule that forever blocks K and Q. Remove the old Ctrl-K note picker (overlay, binding, help text).

## Stay in code (do not put in JSON)

Tab / Shift-Tab and plain Left / Right (cycle Lobby pages). Up / Down and Page Up / Page Down (Files list). Esc (dismiss). Backspace / Left / Right / Home / End inside New / Rename **text fields** and Home / End caret motion in the editor (not Lobby quit / back). Digits 0–9 on the private PIN pad. Enter for dialog confirms / submit where that is structural, not a badgeable action.

## Checklist

- [x] Remove Ctrl-K omni note picker (binding, overlay UI, help / README / home tip).
- [x] Lobby tab shortcuts in JSON (six empty defaults above). Delete hardwired 1–6 jump and hardwired digit badges; badges only when a tab line is set.
- [x] Rotate: keep `settings.rotation` only; remove hardwired Ctrl-R and Ctrl+arrow rotate.
- [x] Home path: add `global.toLobby` and `global.quit`, both default `"hardware_home"`. Letter on a line owns that action only. Keyboard Home ≠ `hardware_home`. Unblock K and Q in the parser.
- [x] Sync: `"sync.now": "enter"`; treat `enter` as a first-class shortcut value (no Ctrl); badge glyph always ↩ when value is `enter`.
- [x] Edit: `"files.edit": "enter"` (wired like other JSON shortcuts; badge via `shortcutBadge`).
- [x] Refresh `shortcuts.body`, Settings help, §16 / §22 / §36 / §37, README, [user-should-test.md](user-should-test.md) — no Ctrl-K picker; no hardwired 1–6; keyboard Home not Lobby/quit; fix §37 Finder/digits wording; no Ctrl-R / arrow rotate. Do not re-list chord tables in decisions — point at `lobby-ui.json`.
- [ ] Verify on tablet: SSH-edit `lobby-ui.json`, chords and badges update without a binary rebuild; bad JSON keeps the last good load. Lobby keyboard test after deploy.

## Out of scope here

Button labels, page copy, colors, and radii — [todo-lobby-ui-chrome.md](todo-lobby-ui-chrome.md).
