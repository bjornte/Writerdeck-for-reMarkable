# Improvements

Wishlist — not a backlog. Shipped: [../DONE.md](../DONE.md). Open verify: [../TODO.md](../TODO.md). Phone vs tablet: [browser-vs-tablet.md](browser-vs-tablet.md).

* Editing and modifier keys
    * Multiple Shift+down arrow also sometimes selects wrong (too much). Tests are either incorrect or too soft.
* Editing, other issues
    * Park for now; seems to be working: When rotated, full screen width is not used in editor or reader
    * Sometimes, touching the screen to move the cursor stops working, then starts working again. Hard to reproduce.
* General
    * Parked, evaluate later: When necessary functionality is migrated away from the phone and over to the reMarkable, hopefully, the PIN to connect from the phone can be removed. BUT this needs consideration. Requires all "security sensitive features", e.g. download, upload, copy & paste, to be acknowledged on the device. When done (if at all): Update all processes, including install scripts.
    * Windows installer still open. Mac/Linux installer remembers Wi-Fi / password / GitHub repo + token and can open a prefilled token page — [TODO.md](../TODO.md) Open.
* Lobby / phone
    * Bring Upload/import back on the Lobby or a thin phone control when needed (API route still exists — [decisions.md](decisions.md) §30).
    * Separate `lobby-ui.json` color for shortcut keycap borders (today they share `borderColor` with outer buttons — cannot gray only the squares).
    * `lobby-ui.json` control for every Lobby type size (Files list, titles, help, dialogs, tabs — not only `labelPointSize` / `badgePointSize`).
    * More `lobby-ui.json` knobs for Lobby chrome still hard-coded in QML:
        * Colors: button fill (`#f0f0f0` / selected `#e8e8e8`), page/dialog backgrounds, dialog scrim (`#dddddd`), vault pad wash (`#f8f8f8`).
        * Corner radius (buttons mostly `6`, dialogs `8`).
        * Files page strip height (`48`), list row side inset (`8`), tab-row extra height (`+8`), dialog box width fraction / padding.
        * Remaining copy: tab titles, Files Prev/Next/Page, Home blurb + open-source line, Keyboard/Sync/Settings section titles and bodies, vault pad prompts, Sync “TOKEN NEEDED” / auto-sync footnote, dialog “PIN: …”.
        * Selection marker in the Files list (`▶`) and “[private]” suffix.

## Possible later
* More USB keyboard layouts beyond US and Norwegian — add sources under `keymaps/`, generate, check on hardware (the phone path does not exercise USB maps).
* Richer edit chrome: indentation helpers, headline jumps, a modest status bar. Not WYSIWYG-in-edit — that stays out ([decisions.md](decisions.md) §1).
* Smoother e-ink redraw inside the current editor if full-frame updates ever feel too heavy. Stay on keywriter; do not switch the product to a terminal stack.
* A read-only Markdown mirror on the phone while you type on e-ink, maybe with tap-to-move-cursor.
* Niceties: bulk delete, search, HTTPS for a native Share sheet.
* reMarkable 2 — only if people ask for it. The hard part is drawing to the screen without Toltec while keeping the Qt editor; launch can lean on the phone, `wd`, or USB Esc. Effort is in the same ballpark as making typing trustworthy, not a rewrite of the product. Stance: [decisions.md](decisions.md) §33.
