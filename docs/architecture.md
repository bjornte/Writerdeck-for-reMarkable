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

Writerdeck-server is a static Go binary — Wi-Fi, APIs, sync, PIN, launching the editor. Source in `daemon/`. Product version is one auto date stamp for the whole product (`YYYY-MM-DD`, or `.N` only when you force a second ship the same day with `scripts/product-version.sh --bump`). Server and editor each carry it; Lobby About shows one number (the older if they differ) and compares to repo-root `VERSION` on GitHub `main`. CI and `deploy-rmkbd.sh` keep `VERSION` current via `scripts/product-version.sh --write` — do not hand-edit that file for routine builds. Why: [decisions.md](decisions.md) §38.

Writerdeck is the full-screen editor, built from our fork of [keywriter](https://github.com/dps/remarkable-keywriter). QML draws the screen and applies edits. C++ starts the app, talks to the display, feeds keys from the socket, and runs EditHelper (math, shortcuts, wrap, undo). Keep hand-tuned wrap gaps and custom undo; do not replace Qt’s text box ([decisions.md](decisions.md) §5–§6).

New editor behavior belongs in [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter). CI (GitHub Actions) clones, checks, and compiles — it does not stitch QML.

Under `/home/root/`: the two binaries, `Writerdeck-user-documents/` for notes, `.Writerdeck/settings.json`, `.Writerdeck/lobby-ui.json`, a Qt runtime, and the launcher script. They meet on `/run/Writerdeck.sock`. The phone page is embedded in the server.

Lobby look, wording, language, and shortcut chords are owned by `lobby-ui.json` and `lobby-ui-i18n/<lang>.json` ([decisions.md](decisions.md) §36). Edit those files on the tablet; Writerdeck reloads them (watch plus a short mtime poll). Repo source of truth for first install is `config/lobby-ui.json` and `config/lobby-ui-i18n/`; `deploy-keywriter.sh` seeds missing tablet files only. Journal shows `lobby-ui: loaded … lang=… (rev N)` on each successful load.

## Phone and Lobby

No phone app — open Safari to the tablet. The phone is a keyboard bridge: paste, sync token, and accepting a Lobby Download offer. Files and settings live in the Lobby ([browser-vs-tablet.md](browser-vs-tablet.md)). Lobby Files paginates on e-ink — fixed pages, no flick; when notes spill a page, Prev / Page N/M / Next sits above the action buttons ([decisions.md](decisions.md) §35). Create and rename go through the trusted editor socket; a name that already exists is refused and shown inside the New / Rename dialog ([decisions.md](decisions.md) §19). Uniqueness ignores letter case and treats plain and encrypted titles as one name; the Files list is sorted the same way.

The server stays up under the stock UI. Boot leaves xochitl on screen; open Writerdeck with page buttons, USB Esc, or `wd`. Home from Lobby brings xochitl back; the phone can still reach port 8000. The PIN is shown in the Lobby. Lobby → Keyboard lists Bluetooth (phone URL, PIN, QR) then USB layout, with live `(connected)` / `(not connected)` on each headline while that tab is open. Sync is optional and change-driven (boot, app open, document open, wake, Wi-Fi up, Home, power sleep — not a timer); the token never hits disk ([server-sync-implementation.md](server-sync-implementation.md)).

Touch Edit / New / Rename (and similar) without a USB keyboard or an open phone page shows a connect tip with the phone URL and a QR image (same QR as Keyboard). An open page counts once it sends WebSocket `hello`; leftover sockets without hello do not. Cursor’s embedded browser (User-Agent contains `Cursor/` or `Electron/`) does not send or count as hello, so agent tabs do not suppress the tip ([decisions.md](decisions.md) §34).

## Constraints

No jailbreak; keep OTA (over-the-air updates) — so no Toltec. One static Go binary on the tablet. Markdown on disk; HTML there is a bug. Device scripts are ASCII and LF. The tablet drops Wi-Fi when it sleeps — keep it awake while developing.

## Device facts

This product is for the reMarkable 1. reMarkable 2 needs a different display path; we are open to exploring that if the community wants it ([decisions.md](decisions.md) §33).

SSH as `root` over Wi-Fi. Password and `RM_HOST_WIFI` live in `secrets/remarkable.local.env` (install prompts if empty). After an OTA (over-the-air update) the password changes. On iPhone hotspot the tablet is often `172.20.10.5`. Visitors fetch prebuilt editor (`keywriter` Release) and server (`server` Release); developers with Go build the server locally via `deploy-rmkbd.sh`.

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

Then relaunch the editor and read `journalctl -u writerdeck` (Writerdeck only — stock UI and system noise live in the rest of `/var/log`). After QML changes: edit-session check. After caret work: automated typing tests. After Lobby/Home: Lobby keyboard test. Deploy uses gzip over SSH, not scp. The same deploy seeds `lobby-ui.json` only if that file is absent on the tablet.
