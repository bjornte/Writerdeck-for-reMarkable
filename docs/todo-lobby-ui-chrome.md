# TODO: Lobby chrome still hardwired

Labels, help copy, and visual knobs that still live in QML instead of `/home/root/.Writerdeck/lobby-ui.json`. Goal: change them over SSH without a rebuild ([decisions.md](decisions.md) §36). Parent: [../TODO.md](../TODO.md). Shortcuts are already on disk: [todo-lobby-ui-shortcuts.md](todo-lobby-ui-shortcuts.md).

Wishlist overlap: [improvements.md](improvements.md) (Lobby / phone).

## Already on disk

Borders, selected-border, margins, row/tab/action heights, spacing, label and badge point sizes, text / border / badge text colors. Dialog titles and buttons, several Settings help strings, private on/off blurbs, `home.tip`, Shortcuts title/body, Files Enter badge. Ctrl-letter *action* chords for Files / Keyboard / Settings (except the gaps in the shortcuts todo).

## Checklist — copy and labels

- [ ] Tab titles (`Files`, `Keyboard`, `Sync`, `Settings`, `Shortcuts`, `Home`).
- [ ] Files: Prev / Next / `Page N`; action labels (New, Edit, Read, Rename, Delete, Download, Encrypt, New encrypted, Decrypt); list marker `▶` and `[private]` suffix.
- [ ] Home: brand title, tagline, “N note(s) on this device”, open-source line (`home.tip` is already JSON).
- [ ] Keyboard: Bluetooth / USB headlines and connected suffixes, pairing and OTG bodies, “Layout”, “PIN: …”, layout button labels (US QWERTY, Norwegian).
- [ ] Sync: page title, status templates, TOKEN NEEDED / SYNC FAILED / SYNC OFFLINE, token body, Sync now / Syncing… / Token needed button text, auto-sync footnote.
- [ ] Settings: page title `"Settings"`; section headers (Reading font, Private notes, PIN for phone pairing, Display rotation, Service); Enable / Change PIN / Exit Writerdeck; PIN option labels and Wi-Fi warn; rotation degree labels.
- [ ] Vault pad: prompts (setup, confirm, change, open, edit, read, encrypt, decrypt, download, wrong PIN); Cancel (dialog Cancel is already JSON — vault should use it).
- [ ] Errors: name-already-exists, Operation failed, and related fallbacks shown in Lobby dialogs.

## Checklist — visual

- [ ] Button fill and selected fill (today roughly `#f0f0f0` / `#e8e8e8`; tabs `#f5f5f5` / `#e0e0e0`).
- [ ] Page / dialog backgrounds, dialog scrim (`#dddddd`), vault pad wash (`#f8f8f8`).
- [ ] Corner radii (buttons ~6, dialogs ~8, keycap badges ~3, sync banners ~4).
- [ ] Layout leftovers: Files page-strip height, list row side inset, tab-row extra height, dialog width fraction and padding.
- [ ] Type sizes beyond `labelPointSize` / `badgePointSize` (Home title, list rows, section titles, dialog titles, sync banners).
- [ ] Separate color for shortcut keycap borders (today they share `borderColor` with outer buttons).

## Verify

After each batch: edit the tablet JSON, confirm Lobby reloads (`lobby-ui: loaded …` in the journal), spot-check the changed page. Prefer one Writerdeck deploy per batch of QML wiring, not per string.
