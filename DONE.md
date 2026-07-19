# What's shipped

Writerdeck adds word processing to the first generation reMarkable. Write with a physical Bluetooth or USB keyboard. Notes save as Markdown.

Open work: [TODO.md](TODO.md). How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

## Your words

Notes are plain Markdown on disk. Most importantly, your documents must be preserved at all times. Sync and file operations must not silently modify any note, including the open one. Details and leftovers: [docs/integrity-audit.md](docs/integrity-audit.md).

## Starting the Writerdeck

Power on lands in the normal reMarkable interface. The Writerdeck server starts in the background so the phone page stays reachable.

From that stock UI: press both left and right page buttons together, or Esc on a USB keyboard, to open the Lobby. From a Mac on the same Wi-Fi: `wd` (or `bash scripts/lobby.sh`). On the tablet over SSH: `~/wd`. On the phone page, **Show PIN on tablet** also opens the Lobby when it was closed.

## How you write

Hook up a physical keyboard, either USB (with an OTG cable) or Bluetooth (via your phone. The Lobby Keyboard tab shows both: Bluetooth first with the phone URL and QR, then USB layout).

## Saving & exiting

Pressing home while in a note saves and returns to Files. Home from the Lobby quits to the stock reMarkable UI.

## Syncing notes

Notes can optionally be synced to a private GitHub repo of your choosing. Ensuring that your work is preserved, and never lost, is the core success factor of the sync. Git history is an extra safety layer. Any missing note at either end is always copied over. Mass deletes are not supported. Checks are in place to never flush notes of content.

Notes syncing is done from the reMarkable, and not the phone. The repository access token stays in the phone's browser and the tablet's RAM. If the token is missing one place but available at the other, it is copied over. For security purposes, it is never saved to reMarkable's disk. Details: [server-sync-implementation.md](docs/server-sync-implementation.md).

## Settings

Except for setting the sync token, which is done on the phone, all settings are available from the Lobby

## Connection between tablet and phone

When using a Bluetooth keyboard, the phone is a bridge for typing. Secondary jobs on the same page: paste-at-cursor while a note is open for edit (not on Lobby Files), accepting a tablet Download offer, and sync token entry. Details: [docs/browser-vs-tablet.md](docs/browser-vs-tablet.md).

Phone keyboard-first: the notes list and Upload are gone from the browser; Lobby Files Download prompts open phones with “Download here?”.

## Lobby

The Lobby is keyboard-first: focus returns after touch, tabs and actions have chords (Shortcuts lists them), the file list pages instead of scrolling, and the private PIN accepts USB or phone-forwarded digits. Edit, new, rename, and similar prompts that need typing first check for a USB keyboard or an open phone/laptop page; if neither is there, a tip shows how to connect, with a QR code for the phone URL. Cursor’s built-in browser does not count as that page — only a real phone or laptop browser does. A wrong private PIN on Encrypt (and similar) stays on the pad with a clear message.

## Security and private notes

PIN each boot — six digits, four, or none. **Note!** Without a PIN, anyone on your Wi-Fi can reach your notes. Five wrong guesses lock that IP for a minute.

Optional vault: a second PIN on the tablet only. Per-note encrypt/decrypt; ciphertext as `.md.enc`. PIN every open. Recovery material can sync to GitHub under `secret/`.

## Editor

Built from our [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) fork of David Singleton’s keywriter. Typing math lives in C++ EditHelper; the screen is QML. Mac/Linux shortcuts pass all automated typing checks. Physical Home is owned by the server while you edit. Power sleeps and wakes with a save. Norwegian USB layout works on hardware.

We stop short of replacing Qt’s text box ([decisions.md](docs/decisions.md) §5–§6).

## Prove it still works

After editor changes: edit-session check, then automated typing tests if caret work moved, Lobby keyboard test if Home/Lobby moved. After vault work: vault scripts. Do not retry uinput — keys use the socket.

## First-time install

Visitors: Download ZIP, `bash scripts/install.sh --start`, answer password + Wi-Fi prompts. Binaries come from GitHub Releases (`keywriter`, `server`) — no `gh` login, no Go. Autostart is enabled only after a phone-page health check. Details and boot-risk follow-up: [docs/install-onboarding/](docs/install-onboarding/).
