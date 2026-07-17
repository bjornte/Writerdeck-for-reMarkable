# What's shipped

Writerdeck adds word processing to the first generation reMarkable. Write with a physical Bluetooth or USB keyboard. Notes save as Markdown.

Open work: [TODO.md](TODO.md). How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

## Your words

Notes are plain Markdown on disk. Most importantly, your documents must be preserved at all times. Sync and file operations must not silently modify any note, including the open one. Details and leftovers: [docs/integrity-audit.md](docs/integrity-audit.md).

## Starting the Writerdeck

If the reMarkable is off, turn it on, and the app launches.

If you're in the main reMarkable user interface, press both left & right pagination buttons simulateously, and you land in the Lobby Files tab. If you have connected a USB keyboard, pressing Esc will also launch the app.

If you have a dev setup for the app on your computer, `wd` will also launch the Writerdeck.

## How you write

Hook up a physical keyboard, either USB (with an OTG cable) or Bluetooth (via your phone. The lobby's keyboard page tells you how).

## Saving & exiting

Pressing home while in a note saves and returns to Files. Home from the Lobby quits to the stock reMarkable UI.

## Syncing notes

Notes can optionally be synced to a private GitHub repo of your choosing. Ensuring that your work is preserved, and never lost, is the core success factor of the sync. Git history is an extra safety layer. Any missing note at either end is always copied over. Mass deletes are not supported. Checks are in place to never flush notes of content.

Notes syncing is done from the reMarkable, and not the phone. The repository access token stays in the phone's browser and the tablet's RAM. If the token is missing one place but available at the other, it is copied over. For security purposes, it is never saved to reMarkable's disk. Details: [server-sync-implementation.md](docs/server-sync-implementation.md).

## Settings

Except for setting the sync token, which is done on the phone, all settings are available from the Lobby

## Connection between tablet and phone

When using a Bluetooth keyboard, the phone is a bridge for typing. However, the interface has several secondary functions, which are useful also from a laptop: upload, download, paste-at-cursor, and sync token entry. Details: [docs/browser-vs-tablet.md](docs/browser-vs-tablet.md).

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
