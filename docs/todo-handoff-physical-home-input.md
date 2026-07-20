# Handoff: Physical Home — single input path

**Done.** One physical Home press → one handler. Exclusive `EVIOCGRAB` on `/dev/input/event1` while a Writerdeck session is active; `suppressNextHomeKey` removed from the fork.

Read first: [architecture.md](architecture.md), [decisions.md](decisions.md) §16, `.cursor/rules/writerdeck.mdc`.

## Architecture (shipped)

```
Physical gpio-keys (Home, Power, page buttons) on /dev/input/event1
  └── Writerdeck-server ONLY while session active (EVIOCGRAB before spawn)
        └── socket cmd → QML handleHome() / prepareSleep() / idle launch chord

USB keyboard (/dev/input/event* except event1) + phone WebSocket
  └── /run/Writerdeck.sock inject → handleKey() → handleHome() for USB Home
```

## What shipped

1. **Exclusive grab in Go** — `daemon/input.go` + `daemon/evdev_grab_linux.go`; `grabButtonDev()` in `session.start()` before editor spawn; `ungrabButtonDev()` in `session.end()`. Idle xochitl keeps buttons (no grab).
2. **Fork cleanup** — [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) `3be2de4`: dropped `suppressNextHomeKey` / `fromPhysicalCmd`; cmd home uses `invokeSaveCmd("handleHome", "home")`.
3. **Docs** — [decisions.md](decisions.md) §16; lessons note on grab-before-spawn.

## Verify

- `deploy-rmkbd.sh` — journal shows `exclusive grab on /dev/input/event1`
- Writerdeck binary from fork `3be2de4` deployed (`fetch-keywriter-dist` + `deploy-keywriter.sh -b`)
- `test-edit-session.sh` PASS
- `test-lobby-keyboard.sh` PASS (Home-from-read via API)
- Manual: press physical Home from edit / read / Lobby once each (owner spot-check)

## Constraints

No jailbreak; preserve OTA. Static Go only: `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`. Do not `pkill -f /home/root/Writerdeck`. Never local `docker build` for Writerdeck.
