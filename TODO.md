# Writerdeck for reMarkable 1 — TODO

> Writerdeck for reMarkable 1 turns a first-gen reMarkable 1 e-paper tablet into a distraction-free Markdown typewriter: type on an iPhone (or laptop) keyboard, the keystrokes travel over Wi-Fi to the tablet, which shows the text on e-ink and saves `.md`.
>
> This file is just the open work. How it works → [docs/architecture.md](docs/architecture.md). Why → [docs/decisions.md](docs/decisions.md). What's shipped → [DONE.md](DONE.md). Gotchas → [docs/lessons.md](docs/lessons.md).
>
> Keystrokes reach the editor over a local socket, not `/dev/uinput` (this kernel can't load uinput — see [docs/decisions.md](docs/decisions.md)). Verify each item on the device before checking it off.

---

## Status

Phases 0–8 and document integrity slices 1–11 are done — the Companion appliance works end-to-end (see [DONE.md](DONE.md)). Optional polish below remains.

---

## Next up

1. **Power button** — needs device test (implementation in [DONE.md](DONE.md) § Editor).
2. **Lobby Ctrl-K on USB keyboard** — needs device verify.

---

## Phase 10 — locales & protection

Design notes: [docs/improvements.md](docs/improvements.md).

### USB keyboard locales (Norwegian first)

- [ ] Generate `no.qmap` (and `us.qmap` baseline) via `ckbcomp` + `kmap2qmap`; ship in `keymaps/`, deploy to `/home/root/keymaps/`.
- [ ] Extend `Writerdeck-launcher.sh` to set `QT_QPA_EVDEV_KEYBOARD_PARAMETERS` from `settings.json` → `keyboardLayout` (default `us`).
- [ ] Hotplug-safe device path — rescan or match Writerdeck-server’s keyboard discovery; document event-node variance.
- [ ] Lobby / Preferences: layout picker (browser + optional tablet Keyboard subpage).
- [ ] Device test: æ ø å Æ Ø Å, AltGr, `@`, `{` `}` on Norwegian USB keyboard.

Ref: [remarkable-keywriter#1](https://github.com/dps/remarkable-keywriter/issues/1) — `loadkeys` / `setxkbmap` do **not** work for Qt apps on rM.

### Encrypted / password-protected note subset

- [ ] Design ADR: encrypted subfolder (e.g. `private/`), passphrase-derived key, session unlock, sync exclusion.
- [ ] Implement only after design sign-off — Go `crypto/*`, rate-limited unlock, locked entries in list API.

---

## Open questions

1. Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit ⇒ a re-deploy + re-`enable` cadence (low-risk; recovery documented in [docs/decisions.md](docs/decisions.md)).

---

## Resume prompt (paste into a fresh chat)

> Project Writerdeck for reMarkable 1 — a reMarkable 1 as a Wi-Fi Markdown typewriter. Writerdeck-server (`/home/root/Writerdeck-server`, built from `daemon/`) serves a WebSocket + HTML capture page and feeds a patched Writerdeck editor over `/run/Writerdeck.sock` (this kernel can't load `/dev/uinput`); Writerdeck saves `.md` to `Writerdeck-user-documents/`. The client is the Mac in dev, the iPhone in use.
> State: Phases 0–8, integrity slices 1–11, and Phase 9 polish are done (see [DONE.md](DONE.md)). **Next:** power button device test; Lobby Ctrl-K USB verify. **Phase 10:** USB locale qmaps and encrypted subfolder — [improvements.md](docs/improvements.md), [TODO.md](TODO.md). Integrity audit: [integrity-audit.md](docs/integrity-audit.md). After Writerdeck/QML edits: `bash scripts/test-edit-session.sh` ([decisions](docs/decisions.md) #21).
> Read first: [architecture](docs/architecture.md), [decisions](docs/decisions.md), [DONE](DONE.md), [lessons](docs/lessons.md), [browser-vs-tablet](docs/browser-vs-tablet.md), [integrity-audit](docs/integrity-audit.md), [improvements](docs/improvements.md). Power-button notes in **Next up** above.
> Dev: device SSH/deploy over Wi-Fi; IP in `secrets/remarkable.local.env` (`RM_HOST_WIFI`, currently `192.168.1.8`).
> Constraints: no jailbreak; preserve OTA; no Toltec; static Go binary (`CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`). SSH password gitignored in `secrets/remarkable.local.env`. Iterate over Wi-Fi; keep the tablet awake.
> Refs: editor https://github.com/dps/remarkable-keywriter · keyboard layouts https://github.com/dps/remarkable-keywriter/issues/1 · input docs https://remarkable.guide/devel/device/input.html.
