# Improvements

Wishlist — not a backlog. Shipped: [../DONE.md](../DONE.md). Open verify: [../TODO.md](../TODO.md). Phone vs tablet: [browser-vs-tablet.md](browser-vs-tablet.md).

## Possible later

remove PIN on phone, require all "security sensitive features", e.g. download & upload, to be ackowledged on the device. Update all processes, including install scripts.

More USB keyboard layouts beyond US and Norwegian — add sources under `keymaps/`, generate, check on hardware (the phone path does not exercise USB maps).

Richer edit chrome: indentation helpers, headline jumps, a modest status bar. Not WYSIWYG-in-edit — that stays out ([decisions.md](decisions.md) §1).

Smoother e-ink redraw inside the current editor if full-frame updates ever feel too heavy. Stay on keywriter; do not switch the product to a terminal stack.

A read-only Markdown mirror on the phone while you type on e-ink, maybe with tap-to-move-cursor.

Phone niceties: bulk delete, search, HTTPS for a native Share sheet, vault PIN on the phone (tablet-only entry ships today).

No-keyboard tip on the tablet when edit / new / rename (and similar) would need typing: how to connect Bluetooth or USB, with a QR code for the current phone URL — [TODO.md](../TODO.md) Open.

Windows installer still open. Mac/Linux installer remembers password / Wi-Fi / GitHub repo + token and can open a prefilled token page — [TODO.md](../TODO.md) Open.

reMarkable 2 — only if people ask for it. The hard part is drawing to the screen without Toltec while keeping the Qt editor; launch can lean on the phone, `wd`, or USB Esc. Effort is in the same ballpark as making typing trustworthy, not a rewrite of the product. Stance: [decisions.md](decisions.md) §33.
