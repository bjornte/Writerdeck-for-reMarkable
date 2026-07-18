# Improvements

Wishlist — not a backlog. Shipped: [../DONE.md](../DONE.md). Open verify: [../TODO.md](../TODO.md). Phone vs tablet: [browser-vs-tablet.md](browser-vs-tablet.md).

* Editing and modifier keys
    * Multiple Shift+down arrow also sometimes selects wrong (too much). Tests are either incorrect or too soft.
    * Observation mode (settings.observe, off by default): phone Observe records keys for Cursor to fetch from /api/observe/export. USB-on-tablet typing is not captured.
* Editing, other issues
    * Park for now; seems to be working: When rotated, full screen width is not used in editor or reader
    * Regression: Focus should be on the last edited file when returning to the lobby from a file.
    * Sometimes, touching the screen to move the cursor stops working, then starts working again. Hard to reproduce.
* General
    * Parked, evaluate later: When necessary functionality is migrated away from the phone and over to the reMarkable, hopefully, the PIN to connect from the phone can be removed. BUT this needs consideration. Requires all "security sensitive features", e.g. download, upload, copy & paste, to be ackowledged on the device. If done: Update all processes, including install scripts.
    * Windows installer still open. Mac/Linux installer remembers Wi-Fi / password / GitHub repo + token and can open a prefilled token page — [TODO.md](../TODO.md) Open.
* Lobby
    * The download button must move to the reMarkable, and be removed from the phone browser. It should prompt a "Download here?" type message in the open phone browser(s)
    * In the file list in the lobby, don't highlight the entire line. Rather, have a ▶︎ -like character before the name. This should lessen need for redraw. Make the ▶︎ rather big so it is not easily overlooked.
    * Make four discrete buttons for rotation orientations.
    * Keyboard focus should be in the lobby the entire time to allow keyboard navigation.
    * All affordances in the lobby should be examined to see if they can be navigated by keyboard
    * In the files interface, there are some shortcuts below. They are overlapping with the bottommost button row. Looks unwieldly. They can potentially be moved to the shortcuts tab, or into the buttons they address.
    * PIN entry should also support USB & Bluetooth keyboards.
    * Does the PIN "Done" button actually do anything?
    * Slow e-ink means file list should paginate, not scroll
    * No-keyboard tip on the tablet when edit / new / rename (and similar) would need typing: how to connect Bluetooth or USB, with a QR code for the current phone URL — [TODO.md](../TODO.md) Open.
* Phone
    * The browser is showing a "disk changed" warning in the browser seemingly whenever I'm editing a document. Why?
    * Unless something critical dictaties it from being so, focus should _always_ be on the keyboard. So inverse of today. Keyboard focus is the _basic_ state, not an exception.
    * The file list must presumably be removed to give way to keyboard focus, and so that we can remove the PIN
    * Maybe only "Paste from phone" and any connection debugging should remain.

## Possible later
* More USB keyboard layouts beyond US and Norwegian — add sources under `keymaps/`, generate, check on hardware (the phone path does not exercise USB maps).
* Richer edit chrome: indentation helpers, headline jumps, a modest status bar. Not WYSIWYG-in-edit — that stays out ([decisions.md](decisions.md) §1).
* Smoother e-ink redraw inside the current editor if full-frame updates ever feel too heavy. Stay on keywriter; do not switch the product to a terminal stack.
* A read-only Markdown mirror on the phone while you type on e-ink, maybe with tap-to-move-cursor.
* Niceties: bulk delete, search, HTTPS for a native Share sheet.
* reMarkable 2 — only if people ask for it. The hard part is drawing to the screen without Toltec while keeping the Qt editor; launch can lean on the phone, `wd`, or USB Esc. Effort is in the same ballpark as making typing trustworthy, not a rewrite of the product. Stance: [decisions.md](decisions.md) §33.
