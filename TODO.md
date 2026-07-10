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

## Next up (lowest effort first)

1. **Lobby: file-picker button** — visible affordance for Ctrl-K omni (`build-keywriter.sh` + deploy).

---

## Phase 9 — Polish / stretch (optional)

- [ ] Cursor navigation niceties (QML patch): ArrowDown on the last line → move cursor to end of line; ArrowUp on the first line → move cursor to start. (Intercept in `Keys.onPressed` when `cursorPosition` is already on the boundary line.)
- [ ] Mac-style modifier+arrow navigation: Alt+Arrow = word jump, Cmd+Arrow = line/doc start/end, Shift+Arrow = select, Shift+Alt/Cmd+Arrow = select by word/line. Match macOS TextEdit behavior.
- [ ] Word/character count, simple status line.
- [ ] Multiple notes / quick-switch UX review.
- [ ] Paragraph spacing in Read view (postponed): Qt 5.15 RichText ignores `margin-bottom` on `<p>`/`<li>`. Next: `line-height`, spacer nodes, or Markdown pre-process — see [docs/lessons.md](docs/lessons.md).
- [x] Sync: marker-aware delete/rename — device-verified ([decisions.md](docs/decisions.md) #19).
- [x] Lobby: repo URL when sync on — device-verified; `pushLobbyInfo` sends `syncOn`/`syncRepo`; Lobby shows `Sync: github.com/owner/repo`.
- [x] Reading view: no auto-scroll to bottom — device-verified; `ensureVisible` only in edit mode (Esc to preview keeps scroll position).
- [x] Browser: battery/Wi-Fi status in top bar — `GET /api/status`, polled every 30s.
- [x] Browser: exit Writerdeck from Settings — `POST /api/shutdown` stops editor, restores xochitl, exits `rmkbd`.
- [ ] Fallback spike (only if keywriter becomes a blocker): a self-contained on-device editor using `libremarkable` (Rust framebuffer) that takes text over the socket — removes the keywriter-compat risk at the cost of building an editor. Documented fallback (see [docs/decisions.md](docs/decisions.md)), not the default.

> Dev-ergonomics polish is already done (deploy ticker, binary-only `rmkw` redeploy, SSH preflight) — see [docs/architecture.md](docs/architecture.md).

> Shipped polish pruned once verified — upload, PIN chooser, Lobby-on-demand, Lobby sync repo line, read-view scroll, fonts, browser sync UI, etc. Lessons in [docs/lessons.md](docs/lessons.md) and [docs/decisions.md](docs/decisions.md). Recover specs from git history if a regression needs them.

---

## Open questions

1. Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit ⇒ a re-deploy + re-`enable` cadence (low-risk; recovery documented in [docs/decisions.md](docs/decisions.md)).

---

## Resume prompt (paste into a fresh chat)

> Project rM1-Writerdeck — a reMarkable 1 as a Wi-Fi Markdown typewriter. A static Go daemon (`rmkbd`) on the tablet serves a WebSocket + HTML capture page and feeds a patched keywriter over a local socket (this kernel can't load `/dev/uinput`); keywriter saves `.md`. The client is the Mac in dev, the iPhone in use.
> State: Phases 0–8 and most of Phase 9 polish are done & device-verified (see Status above + [DONE.md](DONE.md)); only the optional Phase 9 stretch backlog above remains. Do not redo finished phases, retry uinput, or rebuild keywriter from scratch.
> Read first: [architecture](docs/architecture.md), [decisions](docs/decisions.md), [DONE](DONE.md), [lessons](docs/lessons.md). Pick from **Next up** or Phase 9 below.
> Dev: device SSH/deploy over Wi-Fi; IP in `secrets/remarkable.local.env` (`RM_HOST_WIFI`, currently `10.0.0.20`).
> Constraints: no jailbreak; preserve OTA; no Toltec; static Go binary (`CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`). SSH password gitignored in `secrets/remarkable.local.env`. Iterate over Wi-Fi; keep the tablet awake.
> Refs: editor https://github.com/dps/remarkable-keywriter · input docs https://remarkable.guide/devel/device/input.html.
