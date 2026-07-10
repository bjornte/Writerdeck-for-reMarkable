# Writerdeck for reMarkable 1 — TODO

> rM1-Writerdeck turns a first-gen reMarkable 1 e-paper tablet into a distraction-free Markdown typewriter: type on an iPhone (or laptop) keyboard, the keystrokes travel over Wi-Fi to the tablet, which shows the text on e-ink and saves `.md`.
>
> This file is just the open work. How it works → [docs/architecture.md](docs/architecture.md). Why → [docs/decisions.md](docs/decisions.md). What's shipped → [DONE.md](DONE.md). Gotchas → [docs/lessons.md](docs/lessons.md).
>
> Keystrokes reach the editor over a local socket, not `/dev/uinput` (this kernel can't load uinput — see [docs/decisions.md](docs/decisions.md)). Verify each item on the device before checking it off.

---

## Status

Phases 0–8 are done — the Companion appliance works end-to-end (see [DONE.md](DONE.md)). Optional polish below remains.

---

## Next up

1. **Power button** — make functional during Writerdeck sessions (see investigation below).
2. **Lobby: file-picker button** — blocked until USB keyboard (Ctrl-K omni testing).

---

## Power button — investigation (2026-07-10)

**Hardware:** `KEY_POWER` (116) on `/dev/input/event1` (`gpio-keys`), alongside Home (102), page buttons (105/106), and `KEY_WAKEUP` (143). Confirmed with `evtest` on device.

**Stock OS:** `systemd-logind` has `HandlePowerKey=ignore` (`/etc/systemd/logind.conf.d/powerkey.conf`) — logind does not act on the power key. With xochitl running, xochitl/Qt handles sleep/wake UI. Suspend is available (`/sys/power/state` → `mem`).

**Writerdeck gap:** During an editor session xochitl is **stopped**. `rmkbd` watches `event1` but only handles `KEY_HOME` (`watchHomeButton` in `daemon/main.go`). Power presses are ignored. `systemd-inhibit --what=sleep` on keywriter blocks auto-suspend but does not map the physical power key.

**Likely user-visible symptom:** Power button appears dead (no sleep, no shutdown UI) while editing — unlike stock reMarkable behavior.

**Implementation sketch (not started):**
- Extend the gpio-keys watcher to handle `KEY_POWER` (short vs long press if needed).
- Short press: save open note → `systemctl suspend` (or write `mem` to `/sys/power/state`).
- Long press: optional shutdown path (`remarkable-shutdown` / `poweroff`) — needs UX care.
- On wake: rmkbd is still running (`:8000` up); editor session state TBD (resume keywriter vs drop to Lobby).
- Device-test required; cannot fully verify over SSH alone.

Related: [improvements.md](docs/improvements.md) “Some sleep logic” — same family of work.

---

## Phase 9 — Polish / stretch (optional)

- [x] Cursor navigation niceties — Down on last line → end of line; Up on first line → line start (edit mode).
- [x] Mac-style modifier+arrow navigation — device-verified; Alt=word/paragraph, Cmd=line/doc, Shift=select, Home/End line/doc (`handleMacArrow` in QML).
- [x] Sync: marker-aware delete/rename — device-verified ([decisions.md](docs/decisions.md) #19).
- [x] Lobby: repo URL when sync on — device-verified; `pushLobbyInfo` sends `syncOn`/`syncRepo`; Lobby shows `Sync: github.com/owner/repo`.
- [x] Reading view: no auto-scroll to bottom — device-verified; `ensureVisible` only in edit mode (Esc to preview keeps scroll position).
- [x] Browser: battery/Wi-Fi status in top bar — `GET /api/status`, polled every 30s.
- [x] Browser: exit Writerdeck from Settings — `POST /api/shutdown` stops editor, restores xochitl, exits `rmkbd`.

> Dev-ergonomics polish is already done (deploy ticker, binary-only `rmkw` redeploy, SSH preflight) — see [docs/architecture.md](docs/architecture.md).

> Shipped polish pruned once verified — upload, PIN chooser, Lobby-on-demand, Lobby sync repo line, read-view scroll, fonts, browser sync UI, etc. Lessons in [docs/lessons.md](docs/lessons.md) and [docs/decisions.md](docs/decisions.md). Recover specs from git history if a regression needs them.

---

## Open questions

1. Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit ⇒ a re-deploy + re-`enable` cadence (low-risk; recovery documented in [docs/decisions.md](docs/decisions.md)).

---

## Resume prompt (paste into a fresh chat)

> Project rM1-Writerdeck — a reMarkable 1 as a Wi-Fi Markdown typewriter. A static Go daemon (`rmkbd`) on the tablet serves a WebSocket + HTML capture page and feeds a patched keywriter over a local socket (this kernel can't load `/dev/uinput`); keywriter saves `.md`. The client is the Mac in dev, the iPhone in use.
> State: Phases 0–8 and most of Phase 9 polish are done & device-verified (see [DONE.md](DONE.md)). **Next:** power button during Writerdeck sessions; Lobby file-picker blocked on USB keyboard.
> Read first: [architecture](docs/architecture.md), [decisions](docs/decisions.md), [DONE](DONE.md), [lessons](docs/lessons.md). Power-button notes in **Next up** above.
> Dev: device SSH/deploy over Wi-Fi; IP in `secrets/remarkable.local.env` (`RM_HOST_WIFI`, currently `10.0.0.20`).
> Constraints: no jailbreak; preserve OTA; no Toltec; static Go binary (`CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`). SSH password gitignored in `secrets/remarkable.local.env`. Iterate over Wi-Fi; keep the tablet awake.
> Refs: editor https://github.com/dps/remarkable-keywriter · input docs https://remarkable.guide/devel/device/input.html.
