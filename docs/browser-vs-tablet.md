# Browser vs tablet

What the phone is for vs what the tablet is for. How the system fits: [architecture.md](architecture.md). Wishlist: [improvements.md](improvements.md).

The phone is a keyboard bridge, paste helper, and place to paste the GitHub token. Day-to-day files and settings live on the tablet Lobby. Download starts on the tablet and asks open phone browsers to save the file.

## Phone

After the PIN (or with PIN set to none): the page lands on the keyboard shell — gray Writerdeck mark centered. Keys forward to the tablet whenever the page is connected. There is no phone note list and no phone Edit. When Writerdeck has quit (stock UI) a Launch Writerdeck button appears under the logo; on the Files tab that button is hidden.

When the tablet opens a note for edit, the phone shows that note’s name and keeps capturing keys. Paste from phone inserts clipboard text at the tablet cursor — only while that note is open for edit, not on Lobby Files. Bluetooth keyboards pair to the phone. The phone also forwards keys during read preview, Lobby name prompts, delete confirm, and private PIN entry.

When Lobby taps Download, open phone pages get a “Download here?” prompt. Accepting uses the normal attachment download (encrypted notes still ask for the private PIN on the tablet first).

Observe is off by default. Enable over the LAN (keeps the phone Observe button hidden for everyone else):

```bash
source scripts/_env.sh
curl -s -X POST "http://$RM_HOST:8000/api/settings" -H 'Content-Type: application/json' -d '{"observe":true}'
```

Or edit `/home/root/.Writerdeck/settings.json` (`"observe": true`) over SSH. After Stop, the trace sits at `GET /api/observe/export` — say in Cursor chat that you found a bug and the agent pulls it. USB keys typed on the tablet are not recorded.

Sync setup, connection status, and “Show PIN on tablet” live here too. Plain LAN http has no reliable iOS Share sheet.

## Tablet

Files tab lists notes on fixed pages (Up/Down, PgUp/PgDn) — no flick scroll, no sliding one-row window ([decisions.md](decisions.md) §35). When notes spill one screen, Prev / Page N/M / Next sits above the action buttons, separated by a line. Edit to type; Read to preview; Download offers the selected note to open phone browsers. Esc toggles edit/preview. Rename and delete from Files. Settings holds reading font, PIN length, rotation, and Exit. Keyboard tab: two boxes — Bluetooth first (phone URL, PIN, same QR as the connect tip), then USB layout. Each headline shows `(connected)` or `(not connected)` and refreshes about every two seconds while you stay on that tab. Sync tab can run Sync now.

Lobby is driven by the keyboard: focus returns after touch, and every main control has a Ctrl-chord (Shortcuts tab). Switch Lobby pages with Tab, Shift-Tab, or Left/Right — not digits. Borders, help text, dialog copy, and which letter each Ctrl-chord uses come from `/home/root/.Writerdeck/lobby-ui.json` on the tablet ([decisions.md](decisions.md) §36–§37) — not from the phone. Chords avoid browser-reserved Cmd/Ctrl letters and numbers so a phone keyboard can reach the tablet. Bare letters and digits are kept free for a later Finder-style document jump. Private PIN on the tablet accepts USB digits and digits forwarded from the phone. Without a USB keyboard or an open phone/laptop page, Edit / New / Rename show a short connect tip (QR for the phone URL). An open page counts after WebSocket `hello`; Cursor’s embedded browser is excluded (User-Agent `Cursor/` or `Electron/`), so agent tabs do not hide the tip.

Home from edit returns to Files with that note selected. Launch from stock UI with USB Esc, both page buttons, phone **Show PIN on tablet**, or `wd` / `~/wd`. File ops use the trusted socket ([decisions.md](decisions.md) §19).
