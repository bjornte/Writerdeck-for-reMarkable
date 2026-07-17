# Browser vs tablet

What each surface does. How the system fits: [architecture.md](architecture.md). Wishlist: [improvements.md](improvements.md).

The phone is a keyboard bridge, import/export helper, and place to paste the GitHub token. Day-to-day files and settings live on the tablet Lobby.

## Phone

After the PIN: upload and download notes from the list. There is no phone preview and no phone Edit.

When the tablet opens a note for edit, the phone enters Type mode and forwards keys. Paste from phone inserts clipboard text at the tablet cursor — it does not create a file. Bluetooth keyboards pair to the phone. The phone also forwards keys during read preview and Lobby prompts (new, rename, delete confirm).

Sync setup, connection status, and “Show PIN on tablet” live here too. Download uses a normal file attachment; plain LAN http has no reliable iOS Share sheet.

## Tablet

Files tab lists notes. Edit to type; Read to preview. Esc toggles. Rename and delete from Files. Settings holds reading font, PIN length, rotation, and Exit. Keyboard tab picks USB layout. Sync tab can run Sync now.

Home from edit returns to Files with that note selected. Launch from stock UI with USB Esc, both page buttons, or `wd` / `~/wd`. File ops use the trusted socket ([decisions.md](decisions.md) §19).
