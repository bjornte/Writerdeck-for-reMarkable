# Improvements

Wishlist — not a backlog. Shipped: [../DONE.md](../DONE.md). Open verify: [../TODO.md](../TODO.md). Phone vs tablet: [browser-vs-tablet.md](browser-vs-tablet.md).

- Consider having a presentation mode where the screen is sent to the browser. Should be triggered from the tablet and accepted in the browser, similar to the download function. Should support having two browsers connected: One (typically a laptop) to present the material, another (typically a mobile) to act as the Bluetooth keyboard bridge.
  - Consider a presenter inside View mode that is like the presenter view in Confluence, where it shows only one title's content at a time.
- A smarter scroll. reinvent the scroll bar of windowing, but optimized for e-ink 
- Editing and modifier keys
  - Multiple Shift+down arrow also sometimes selects wrong (too much). Tests are either incorrect or too soft.
- Editing, other issues
  - Park for now; seems to be working: When rotated, full screen width is not used in editor or reader
  - Sometimes, touching the screen to move the cursor stops working, then starts working again. Hard to reproduce.
- General
  - Parked, evaluate later: When necessary functionality is migrated away from the phone and over to the reMarkable, hopefully, the PIN to connect from the phone can be removed. BUT this needs consideration. Requires all "security sensitive features", e.g. download, upload, copy & paste, to be acknowledged on the device. When done (if at all): Update all processes, including install scripts.
  - Windows installer still open. Mac/Linux installer remembers Wi-Fi / password / GitHub repo + token and can open a prefilled token page — [TODO.md](../TODO.md) Open.
- Lobby / phone
  - Bring Upload/import back on the Lobby or a thin phone control when needed (API route still exists — [decisions.md](decisions.md) §30).
  - Lobby chrome / Latin i18n on disk (shipped) — [todo-lobby-ui-chrome.md](todo-lobby-ui-chrome.md). Change language via `"language"` in `lobby-ui.json` for now; add a Settings-tab language picker later. Do not leave old English keys in the tablet `strings` block — they override the language pack. Phone UI follows that language (`daemon/phone-ui-i18n/`).
  - Korean, Chinese, and Japanese Lobby strings need CJK fonts on the tablet (today only Latin Noto/DejaVu). Do not ship those packs until fonts land.
  - Images in an `/img` folder (`![text](img/…)` in edit, show in read). Hard today: sync and APIs are flat Markdown-only, no tablet upload path, and e-ink RichText `<img>` is unproven. A local preview spike is possible; shipping with GitHub sync is weeks of work.



## Possible later

- More USB keyboard layouts beyond US / Norwegian / Spanish / German / French — add sources under `keymaps/`, generate, check on hardware (the phone path does not exercise USB maps).
- Richer edit chrome: indentation helpers, headline jumps, a modest status bar. Not WYSIWYG-in-edit — that stays out ([decisions.md](decisions.md) §1).
- Smoother e-ink redraw inside the current editor if full-frame updates ever feel too heavy. Stay on keywriter; do not switch the product to a terminal stack.
- A read-only Markdown mirror on the phone while you type on e-ink, maybe with tap-to-move-cursor.
- Niceties: bulk delete, search, HTTPS for a native Share sheet.
- reMarkable 2 — only if people ask for it. The hard part is drawing to the screen without Toltec while keeping the Qt editor; launch can lean on the phone, `wd`, or USB Esc. Effort is in the same ballpark as making typing trustworthy, not a rewrite of the product. Stance: [decisions.md](decisions.md) §33.

