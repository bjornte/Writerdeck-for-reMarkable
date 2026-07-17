# Architecture

How Writerdeck fits together. Why: [decisions.md](decisions.md). Shipped: [../DONE.md](../DONE.md). Open: [../TODO.md](../TODO.md). Gotchas: [lessons.md](lessons.md). Words: [terms.md](terms.md).

## What this is

The reMarkable 1 has a fine e-ink screen and a quiet OS, but no word processor and no keyboard support. Writerdeck adds both. Use a USB keyboard on an OTG cable, or a Bluetooth keyboard paired to your phone — those keys reach the tablet over Wi-Fi. The tablet shows the page and keeps Markdown on disk.

```
Phone (Safari + keyboard)
    |  WebSocket on the LAN
    v
reMarkable 1
  Writerdeck-server — always on: files, sync, PIN, key relay
  Writerdeck — full-screen editor; reads /run/Writerdeck.sock; saves .md
```

Keys reach the editor over that socket. This kernel cannot load a uinput (fake keyboard device) ([decisions.md](decisions.md) §3).

## Document integrity

Before anything that touches notes, ask: can this lose text, write the wrong bytes, or overwrite without the user knowing?

Files are UTF-8 Markdown. An open note is protected from silent sync overwrite. Saves use defined paths, autosave, and save-before-stop. GitHub sync backs up; it must not empty-push or delete against a live edit. Policy: [decisions.md](decisions.md) § Document integrity. Leftovers: [integrity-audit.md](integrity-audit.md).

## Two programs on the tablet

Writerdeck-server is a static Go binary — Wi-Fi, APIs, sync, PIN, launching the editor. Source in `daemon/`.

Writerdeck is the full-screen editor, built from our fork of [keywriter](https://github.com/dps/remarkable-keywriter). QML draws the screen and applies edits. C++ starts the app, talks to the display, feeds keys from the socket, and runs EditHelper (math, shortcuts, wrap, undo). Keep hand-tuned wrap gaps and custom undo; do not replace Qt’s text box ([decisions.md](decisions.md) §5–§6).

New editor behavior belongs in [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter). CI (GitHub Actions) clones, checks, and compiles — it does not stitch QML.

Under `/home/root/`: the two binaries, `Writerdeck-user-documents/` for notes, `.Writerdeck/settings.json`, a Qt runtime, and the launcher script. They meet on `/run/Writerdeck.sock`. The phone page is embedded in the server.

## Phone and Lobby

No phone app — open Safari to the tablet. Upload, download, paste, and the sync token stay on the phone. Files and settings live in the Lobby ([browser-vs-tablet.md](browser-vs-tablet.md)).

The server stays up under the stock UI. Home from Lobby brings xochitl back; the phone can still reach port 8000. A PIN appears on e-ink each boot. Sync is optional; the token never hits disk ([server-sync-implementation.md](server-sync-implementation.md)).

## Constraints

No jailbreak; keep OTA (over-the-air updates) — so no Toltec. One static Go binary on the tablet. Markdown on disk; HTML there is a bug. Device scripts are ASCII and LF. The tablet drops Wi-Fi when it sleeps — keep it awake while developing.

## Device facts

SSH as `root` over Wi-Fi. Put the password and `RM_HOST_WIFI` in `secrets/remarkable.local.env`. After an OTA (over-the-air update) the password changes. On iPhone hotspot the tablet is often `172.20.10.5`.

While Writerdeck is open, the server grabs the physical Home button so Qt does not see it twice ([decisions.md](decisions.md) §16). Rootfs is nearly full; everything we ship lives on `/home/root`. Do not resize rootfs.

## Build and deploy

Server from the Mac:

```bash
bash scripts/deploy-rmkbd.sh
```

Editor from CI (GitHub Actions); QML is inside the binary:

```bash
git push && bash scripts/fetch-keywriter-dist.sh && bash scripts/deploy-keywriter.sh -b
```

Then relaunch the editor and read `journalctl -u writerdeck`. After QML changes: edit-session check. After caret work: automated typing tests. After Lobby/Home: Lobby keyboard test. Deploy uses gzip over SSH, not scp.
