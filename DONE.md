# What's shipped

Writerdeck adds keyboard word processing to the first-gen reMarkable. Use a USB keyboard on an OTG cable, or a Bluetooth keyboard paired to your phone (keys cross Wi-Fi to the tablet). Notes save as Markdown.

Open work: [TODO.md](TODO.md). How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

## Your words

Notes are plain Markdown on disk. Sync and file ops must not silently empty or overwrite an open note. Details and leftovers: [docs/integrity-audit.md](docs/integrity-audit.md).

## How you write

Power on and you land in the Lobby Files tab. Open the tablet address on the phone, enter the PIN, pair a keyboard to the phone if you like. Open a note on the tablet — the phone enters Type mode and keys land on e-ink.

Home from a note saves and returns to Files with that note selected. Home from the Lobby quits to the stock reMarkable UI. The server keeps answering on port 8000 either way. Relaunch with USB Esc, both page buttons together, or `wd` / `~/wd`.

## Phone and tablet split

The phone is a keyboard bridge, upload/download, paste-at-cursor, and sync token entry. Day-to-day files and settings live on the tablet Lobby. Full split: [docs/browser-vs-tablet.md](docs/browser-vs-tablet.md).

## Security and private notes

PIN each boot — six digits, four, or none (none warns that anyone on your Wi-Fi can connect). Five wrong guesses lock that IP for a minute.

Optional vault: a second PIN on the tablet only. Per-note encrypt/decrypt; ciphertext as `.md.enc`. PIN every open. Recovery material can sync to GitHub under `secret/`. [decisions.md](docs/decisions.md) §12.

## Editor

Built from our [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) fork of Dave’s keywriter. Typing math lives in C++ EditHelper; the screen is QML. Mac/Linux chords are harness-green (**110/110/0**). Physical Home is owned by the server while you edit. Power sleeps and wakes with a save. Norwegian USB layout works on hardware.

We stop short of replacing Qt’s text box ([decisions.md](docs/decisions.md) §5–§6).

## Sync

Optional GitHub reconcile — copy missing notes both ways, never mass-delete, refuse empty push over a known-good note. Token stays in the browser and tablet RAM. [server-sync-implementation.md](docs/server-sync-implementation.md).

## Prove it still works

After editor changes: edit-session smoke test, then keyboard harness if caret work moved, Lobby keyboard test if Home/Lobby moved. After vault work: vault scripts. Do not retry uinput — keys use the socket.
