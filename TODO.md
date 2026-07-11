# Writerdeck for reMarkable 1 — TODO

> Writerdeck for reMarkable 1 turns a first-gen reMarkable 1 e-paper tablet into a distraction-free Markdown typewriter: type on an iPhone (or laptop) keyboard, the keystrokes travel over Wi-Fi to the tablet, which shows the text on e-ink and saves `.md`.
>
> This file is just the open work. How it works → [docs/architecture.md](docs/architecture.md). Why → [docs/decisions.md](docs/decisions.md). What's shipped → [DONE.md](DONE.md). Gotchas → [docs/lessons.md](docs/lessons.md).
>
> Keystrokes reach the editor over a local socket, not `/dev/uinput` (this kernel can't load uinput — see [docs/decisions.md](docs/decisions.md)). Verify each item on the device before checking it off.

---

## Status

Phases 0–8 are done — the Companion appliance works end-to-end (see [DONE.md](DONE.md)). Optional polish below remains.

---

## Next up

1. **Power button** — implemented (Option A): save → sleep screen → suspend; power again to wake. **Needs device test.**

---

## Power button — investigation (2026-07-10)

**Shipped (pending device test):** Short power press while editing → `prepareSleep()` (save + on-screen message) → stop Writerdeck → `systemctl suspend`. Message:

> Writerdeck is sleeping.
>
> Press power to wake.

Power again after wake → restart editor session and reopen the note that was open.

**Hardware:** `KEY_POWER` (116) on `/dev/input/event1` (`gpio-keys`), alongside Home (102), page buttons (105/106), and `KEY_WAKEUP` (143).

**Stock OS:** `systemd-logind` has `HandlePowerKey=ignore`. With xochitl running, xochitl handles sleep/wake. During Writerdeck sessions xochitl is stopped; Writerdeck-server watches Power/Wakeup via `watchPhysicalButtons`.

---

## Phase 9 — Polish / stretch (optional)

- [x] Cursor navigation niceties — Down on last line → end of line; Up on first line → line start (edit mode).
- [x] Mac-style modifier+arrow navigation — device-verified; Alt=word/paragraph, Cmd=line/doc, Shift=select, Home/End line/doc (`handleMacArrow` in QML).
- [x] Sync: marker-aware delete/rename — device-verified ([decisions.md](docs/decisions.md) #19).
- [x] Lobby: repo URL when sync on — device-verified; `pushLobbyInfo` sends `syncOn`/`syncRepo`; Lobby Syncing section shows repo and last-sync time.
- [x] Lobby: structured sections — Notes, Syncing, Keyboard connection, Shortcuts, footer; vertically centered when content fits.
- [x] Lobby: last sync on tablet — `lastSyncAt` in settings; `/api/sync/ack` after every successful reconcile; shown in Syncing section.
- [x] Reading view: no auto-scroll to bottom — device-verified; `ensureVisible` only in edit mode (Esc to preview keeps scroll position).
- [x] Browser: connection status in top bar — single indicator: **Tablet offline** / **Connecting…** / **Connected · 85%** (`GET /api/status`, polled every 5 s; HTTP is source of truth).
- [x] Browser: exit Writerdeck from Preferences — `POST /api/shutdown` stops editor, restores xochitl, exits Writerdeck-server.
- [x] USB Escape launch — from stock UI (no session), Esc on USB keyboard → Lobby (`watchUSBKeyboardForLaunch` in `daemon/main.go`). Not Esc-to-wake (power button only).
- [x] Edit from browser — regression fixed (2026-07-11): patch 7p missing `}` in `handleKey()` broke QML load; `scripts/test-edit-session.sh` guards it.
- [ ] Lobby: Ctrl-K note picker from Lobby — shipped; **verify on USB keyboard**.

> Dev-ergonomics polish is already done (deploy ticker, binary-only `rmkw` redeploy, SSH preflight) — see [docs/architecture.md](docs/architecture.md).

> Shipped polish pruned once verified — upload, PIN chooser, Lobby-on-demand, Lobby sync repo line, read-view scroll, fonts, browser sync UI, etc. Lessons in [docs/lessons.md](docs/lessons.md) and [docs/decisions.md](docs/decisions.md) (#21 Edit-session test). Recover specs from git history if a regression needs them.

---

## Open questions

1. Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit ⇒ a re-deploy + re-`enable` cadence (low-risk; recovery documented in [docs/decisions.md](docs/decisions.md)).

---

## Resume prompt (paste into a fresh chat)

> Project Writerdeck for reMarkable 1 — a reMarkable 1 as a Wi-Fi Markdown typewriter. Writerdeck-server (`/home/root/Writerdeck-server`, built from `daemon/`) serves a WebSocket + HTML capture page and feeds a patched Writerdeck editor over `/run/Writerdeck.sock` (this kernel can't load `/dev/uinput`); Writerdeck saves `.md` to `Writerdeck-user-documents/`. The client is the Mac in dev, the iPhone in use.
> State: Phases 0–8 and most of Phase 9 polish are done & device-verified (see [DONE.md](DONE.md)). **Next:** power button device test; Lobby Ctrl-K verify on USB keyboard. After Writerdeck/QML edits, run `bash scripts/test-edit-session.sh` ([decisions](docs/decisions.md) #21).
> Read first: [architecture](docs/architecture.md), [decisions](docs/decisions.md), [DONE](DONE.md), [lessons](docs/lessons.md). Power-button notes in **Next up** above.
> Dev: device SSH/deploy over Wi-Fi; IP in `secrets/remarkable.local.env` (`RM_HOST_WIFI`, currently `192.168.1.8`).
> Constraints: no jailbreak; preserve OTA; no Toltec; static Go binary (`CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`). SSH password gitignored in `secrets/remarkable.local.env`. Iterate over Wi-Fi; keep the tablet awake.
> Refs: editor https://github.com/dps/remarkable-keywriter · input docs https://remarkable.guide/devel/device/input.html.
