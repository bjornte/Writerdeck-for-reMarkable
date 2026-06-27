# Writerdeck for reMarkable 1 — TODO

> rM1-Writerdeck turns a first-gen reMarkable 1 e-paper tablet into a distraction-free Markdown typewriter: type on an iPhone (or laptop) keyboard, the keystrokes travel over Wi-Fi to the tablet, which shows the text on e-ink and saves `.md`.
>
> This file is just the open work. How it works (architecture, environment, dev loop) → [docs/architecture.md](docs/architecture.md). Why each choice was made (the ADR) → [docs/decisions.md](docs/decisions.md). What's already done (the dated progress log) → [DONE.md](DONE.md).
>
> Keystrokes reach the editor over a local socket, not `/dev/uinput` (this kernel can't load uinput — see [docs/decisions.md](docs/decisions.md)). Verify each item on the device before checking it off.

---

## Status

Phases 0–8 are done and device-verified — the Companion appliance works end-to-end (the full flow and phase table are in [DONE.md](DONE.md)). All planned build work is done; only the optional Phase 9 polish below remains.

---

## Phase 9 — Polish / stretch (optional)

- [ ] Cursor navigation niceties (QML patch): ArrowDown on the last line → move cursor to end of line; ArrowUp on the first line → move cursor to start. (Intercept in `Keys.onPressed` when `cursorPosition` is already on the boundary line.)
- [ ] Mac-style modifier+arrow navigation: Alt+Arrow = word jump, Cmd+Arrow = line/doc start/end, Shift+Arrow = select, Shift+Alt/Cmd+Arrow = select by word/line. Match macOS TextEdit behavior.
- [ ] Word/character count, simple status line.
- [ ] Multiple notes / quick-switch UX review.
- [ ] Battery/Wi-Fi indicators on the capture page.
- [ ] Paragraph spacing in Read view (postponed, was spec item C): Qt 5.15 RichText ignores `margin-bottom` on `<p>`/`<li>`, so the inter-paragraph gap didn't change on e-ink (the `readHtml()` helper stays wired, so resuming only swaps the injected CSS). Next approaches: `line-height` on `<p>`, an injected spacer paragraph (`<p>&nbsp;</p>`), or pre-process the Markdown for vertical rhythm — see [DONE.md](DONE.md) 2026-06-27.
- [ ] Fallback spike (only if keywriter becomes a blocker): a self-contained on-device editor using `libremarkable` (Rust framebuffer) that takes text over the socket — removes the keywriter-compat risk at the cost of building an editor. Documented fallback (see [docs/decisions.md](docs/decisions.md)), not the default.

> Dev-ergonomics polish is already done (deploy ticker, binary-only `rmkw` redeploy, SSH preflight) — see [docs/architecture.md](docs/architecture.md).

> Shipped Phase 9 polish was pruned once device-verified — upload, the 6/4/none PIN chooser and Lobby-on-demand (P/L), the reading-view font picker, and the edit-view ↔ browser sync (S0/S1/S2) all used to live here as checklist items + specs; their durable lessons are now in [DONE.md](DONE.md) (the dated log) and [docs/decisions.md](docs/decisions.md) (#16–#18). Recover full specs from git history if a regression needs them.

---

## Open questions

1. Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit ⇒ a re-deploy + re-`enable` cadence (low-risk; recovery documented in [docs/decisions.md](docs/decisions.md)).

---

## Resume prompt (paste into a fresh chat)

> Project rM1-Writerdeck — a reMarkable 1 as a Wi-Fi Markdown typewriter. A static Go daemon (`rmkbd`) on the tablet serves a WebSocket + HTML capture page and feeds a patched keywriter over a local socket (this kernel can't load `/dev/uinput`); keywriter saves `.md`. The client is the Mac in dev, the iPhone in use.
> State: Phases 0–8 and most of Phase 9 polish are done & device-verified (see Status above + [DONE.md](DONE.md)); only the optional Phase 9 stretch backlog above remains. Do not redo finished phases, retry uinput, or rebuild keywriter from scratch.
> Read first: [docs/architecture.md](docs/architecture.md) (how it works), [docs/decisions.md](docs/decisions.md) (the why / ADR), [DONE.md](DONE.md) (what's done), then pick an item from the Phase 9 — Polish / stretch list above.
> Dev: the assistant commits on one machine; device SSH/deploy runs on the dev machine over Wi-Fi `192.168.1.8`; git bridges them.
> Constraints: no jailbreak; preserve OTA; no Toltec; static Go binary (`CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`). SSH password gitignored in `secrets/remarkable.local.env`. Iterate over Wi-Fi; keep the tablet awake.
> Refs: editor https://github.com/dps/remarkable-keywriter · input docs https://remarkable.guide/devel/device/input.html.
