# Browser vs tablet — capability matrix

What each surface can do today. Contract: [architecture.md](architecture.md). Shipped features: [../DONE.md](../DONE.md). Future work: [improvements.md](improvements.md).

The phone/Mac companion (`daemon/index.html` + `app.js` + `sync.js`) is the full control surface. The e-ink Lobby (six-tab pager) and Ctrl-K picker cover most day-to-day file ops on tablet.

| Capability | Browser (authed) | Tablet (USB / Lobby) |
|---|---|---|
| List notes | Yes — scrollable list with size/date | **Files** tab — `ListView` from socket `noteslist`; note count on Home |
| Create note | **New**, **Upload**, paste-on-create | **Files → New** or **Ctrl-K** → type name → Enter |
| Open / edit | **Edit** → Type mode + `POST /api/open` | **Files → Open** / Enter / double-tap; **Ctrl-K** picker |
| Read / preview | Read view (`textContent`, safe) | Esc toggles edit/preview in Writerdeck |
| Rename | Read view → **Rename** (`PATCH /api/notes/{name}`) | **Files → Rename** or `r` (socket `renamenote`) |
| Delete | Read view → **Delete** (`DELETE`); pairs with GitHub if sync on | **Files → Delete** or `d` + confirm (socket `deletenote`) |
| Download / copy | **Download**, **Copy** (http fallback) | No |
| Paste at cursor | **Paste from here** (Type mode) | No |
| Font (read view) | Preferences → pick Inter / Literata / EB Garamond / DejaVu | Phone pushes `setfont`; Settings tab shows current font |
| PIN length | Preferences → 6 / 4 / none | Display only (Lobby shows PIN; no change on device) |
| Display rotation | Preferences → **Rotate tablet 90°** | **Ctrl-R** / Ctrl+←/→; Settings tab **Rotate 90** button |
| Sync config | **Setup** panel — on/off, repo, token (`POST /api/sync/token` → tablet RAM) | Read-only status on **Sync** tab |
| Sync run | — (use tablet Lobby **Sync now**) | **Sync now** button + automatic triggers |
| Connection status | Top bar — offline / connecting / connected + battery | Not shown on e-ink |
| Show PIN on tablet | **Show PIN on tablet** (`POST /api/lobby`, pre-auth) | N/A — you are looking at the PIN |
| Exit Writerdeck | Preferences → **Exit Writerdeck** (`POST /api/shutdown`) | **Ctrl-Q** or Home from Lobby |
| Launch Lobby | **Edit** without note (starts session) | **Esc** / L+R page buttons; Mac `wd`; tablet `~/wd` |
| Keystrokes | WebSocket — layout resolved by phone OS | USB — Qt evdev, **US QWERTY default**; BT — same as browser path |

**Takeaway:** upload, download, copy, paste, and GitHub **setup** (toggle, repo, token) remain browser-only. **Sync now** lives on the tablet Lobby Sync tab; the engine runs on the tablet. Tablet **Files** CRUD via trusted socket — shipped and device-verified 2026-07-11 ([decisions.md](decisions.md) #23).
