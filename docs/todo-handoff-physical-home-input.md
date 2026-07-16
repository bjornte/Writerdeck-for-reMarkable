# Handoff: Physical Home — single input path

Give this file to a fresh agent. Goal: one physical Home press → one handler. No timers, no pairing flags.

Read first: [architecture.md](architecture.md), [decisions.md](decisions.md) §8 (two-level Home), `.cursor/rules/writerdeck.mdc` (deploy verification).

## Problem

The reMarkable 1 middle (Home) button is wired through gpio-keys on `/dev/input/event1`. Today two stacks see the same press:

1. **Go** — `daemon/input.go` `watchPhysicalButtons()` reads `event1`, sends `{"t":"cmd","c":"home"}` on the editor socket, and in parallel broadcasts `exitedit`, clears `currentNote`, runs sync reconcile.
2. **Qt** — Writerdeck's evdev plugin also reads gpio-keys, delivering `Qt.Key_Home` into QML `handleKey()`.

One press can therefore run `handleHome()` twice: cmd path returns edit/read → Lobby, then `Key_Home` sees `isLobby` and quits. Read mode had the same bug via upstream `Qt.quit()` on `Key_Home` (fixed in edit 7p of `build-keywriter.sh`).

Interim fix (shipped in `9981f86`): `handleHome(fromPhysicalCmd)` sets `suppressNextHomeKey` when the cmd path transitions to Lobby; the duplicate `Key_Home` consumes the flag. This works but is a bridge between two pipelines — not the long-term shape.

## Target architecture

```
Physical gpio-keys (Home, Power, page buttons) on /dev/input/event1
  └── Writerdeck-server ONLY (exclusive grab — Qt must not see event1)
        └── socket cmd → QML handleHome() / prepareSleep() / idle launch chord

USB keyboard (/dev/input/event* except event1) + phone WebSocket
  └── /run/Writerdeck.sock inject → handleKey() → handleHome() for USB Home
```

After this: delete `suppressNextHomeKey`, the `fromPhysicalCmd` argument, and the C++ `Q_ARG(QVariant, true)` on cmd home in fork [`bjornte/Writerdeck-keywriter`](https://github.com/bjornte/Writerdeck-keywriter) `main.cpp` (socket cmd handler).

## Implementation steps (in order)

1. **Exclusive grab in Go** — In `watchPhysicalButtons()` (`daemon/input.go`, constants in `daemon/editor.go`: `buttonDev`, `keyHome`, etc.), grab `event1` so Qt evdev never receives physical Home/Power/page buttons. Linux: `EVIOCGRAB` / `ioctl` on the open fd (verify on device — may need `golang.org/x/sys/unix` or a tiny cgo-free ioctl). Grab must happen before Writerdeck starts and stay held for the session lifetime. USB keyboard devices are already discovered separately via `findKeyboardInputDevices()` (excludes `buttonDev`).

2. **Verify Qt no longer gets physical Home** — Deploy server only first; with Writerdeck running, physical Home should change lobby/quit only via cmd path. No duplicate quit. USB keyboard Home (`Key_Home` over socket or evdev on a USB event node) must still work from edit, read, and lobby.

3. **Remove pairing hack from QML** — In fork `main.qml` (Writerdeck-keywriter): drop `suppressNextHomeKey` property, `fromPhysicalCmd` parameter, and the consume branch at the top of `handleKey`. Restore cmd home to plain `invokeSaveCmd("handleHome", "home")` in fork `main.cpp`. Keep a single `Key_Home` entry point at the top of `handleKey` → `handleHome()`.

4. **Move server side-effects to editor ack (optional but recommended)** — Today `input.go` fires `exitedit` and sync when the button is pressed, in a goroutine parallel to save. Tie `exitedit`, `currentNote` clear, and `syncEng.reconcileAll("home")` to the successful `saved`/`home` ack from the editor (same pattern as other save-driven state). Button handler should only send cmd home.

5. **Document** — Add a short ADR in `docs/decisions.md` (why exclusive grab; USB Home unchanged). One line in `docs/lessons.md` if grab ordering matters (grab before spawn).

## Files

| Area | Path |
|------|------|
| Button watcher | `daemon/input.go`, `daemon/editor.go` |
| QML patches | fork `main.qml` (Writerdeck-keywriter) |
| Socket cmd | fork [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) `main.cpp` |
| Tests | `daemon/cmd/lobby-keyboard-test/main.go`, `scripts/test-lobby-keyboard.sh` |
| Docs | `docs/decisions.md`, optionally `docs/lessons.md` |

## Acceptance criteria

- Physical Home from edit → Lobby (save completes, no instant quit).
- Physical Home from read → Lobby (regression that crashed Writerdeck before edit 7p).
- Physical Home from Lobby → quit Writerdeck, xochitl returns; server keeps serving `:8000`.
- Second physical Home press (800 ms debounce in Go) from Lobby → quit; no double-fire on first press.
- USB keyboard Home: edit → Lobby, lobby → quit, read → Lobby — all single press.
- `grep suppressNextHomeKey` and `grep homeArrivalGrace` empty in repo.
- No new QML timers for Home semantics.

## Verify (mandatory before done)

1. `bash scripts/deploy-rmkbd.sh` → `systemctl restart writerdeck` on device.
2. QML change: `git push` → `bash scripts/fetch-keywriter-dist.sh` → `bash scripts/deploy-keywriter.sh -b` (relaunch Writerdeck — server restart alone is not enough).
3. `bash scripts/test-lobby-keyboard.sh` — must PASS (includes Home-from-read via `POST /api/test/home`).
4. `bash scripts/test-edit-session.sh`.
5. SSH: `journalctl -u writerdeck -n 30` — no QML parse errors, no instant `editor process exited`.
6. Manual on device: Home from read, Home from edit, Home×2 from Lobby; USB Home from lobby if keyboard available.

## Constraints

No jailbreak; preserve OTA. Static Go only: `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`. Do not `pkill -f /home/root/Writerdeck` (matches Writerdeck-server). Never local `docker build` for Writerdeck — CI + `fetch-keywriter-dist.sh`.

## Context links

Two-level Home semantics: [decisions.md](decisions.md) §8. Page-button chord on same `event1`: §7c. Socket inject: [decisions.md](decisions.md) §1. Prior commit with interim fix: `9981f86`.

## Resume prompt (paste into a fresh chat)

> Read `docs/handoff-physical-home-input.md`, then `.cursor/rules/writerdeck.mdc`, `docs/architecture.md`, and `docs/decisions.md` §8. Implement exclusive gpio grab on `/dev/input/event1` in the daemon, remove `suppressNextHomeKey`, verify with `test-lobby-keyboard.sh` and manual Home on device.
