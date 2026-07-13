# Browser vs tablet

What each surface can do today. Contract: [architecture.md](architecture.md). Future work: [improvements.md](improvements.md).

The phone is a keyboard bridge, import/export helper, and GitHub token entry surface. Day-to-day file management and device settings live on the tablet Lobby.

## Phone

Connect with the PIN, then use the note list for **Upload** and **Download** only — not for reading or opening notes. **Upload** imports an external `.md`, `.markdown`, or `.txt` as a new note on disk. **Download** saves one note to the phone via `Content-Disposition: attachment`. There is no plaintext preview and no phone-initiated Edit.

Open a note on the tablet for **edit** (Files → Edit, Enter, or second tap on a row). The server sends WebSocket `openedit` and the phone enters **Type mode**: keystrokes forward to e-ink with an echo footer. In Type mode, **Paste from phone** (or "Paste from here" on iPad) opens a modal, reads the clipboard, and replays the text at the current tablet cursor through the existing keystroke path — it does not create a note and is not Upload.

Bluetooth keyboards pair to the phone, not the tablet. Besides Type mode, the phone also captures keys when the tablet is in **read preview** (`openread` — Esc toggles to edit), **Files new/rename** (`lobbyinput`), or **delete confirm** (Enter/Esc). Read mode shows a minimal typing panel; lobby prompts show a banner over the note list.

Also on the phone: GitHub sync in Notes sync setup (bar: Sync setup); connection status; Show PIN on tablet. Keystroke layout resolves in the phone OS (Norwegian works today).

## Tablet

List notes in the Files tab; create via New or Ctrl-K; open for edit via Edit, Enter, double-tap, or second tap on the selected row; open for read via Read or `v` (preview on e-ink); toggle edit and preview with Esc; rename and delete with `r` and `d`; reading font and PIN length in Settings; rotate with Ctrl-R or Ctrl-arrows; Exit Writerdeck in Settings (confirm with Enter); sync status and Sync now on the Sync tab; USB keyboard layout in the Keyboard tab; type from USB (Qt evdev qmap) or Bluetooth (same WebSocket path as the phone). Launch Lobby: Esc, L+R page buttons, `wd`, `~/wd`.

Upload, download, paste-at-cursor, and GitHub token entry stay browser-only. Sync engine runs on the tablet; Sync now lives on the Lobby Sync tab. Tablet Files CRUD uses the trusted socket — [decisions.md](decisions.md) §24.
