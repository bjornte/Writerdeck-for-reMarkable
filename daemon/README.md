# daemon/ — Writerdeck-server

The Go daemon that runs on the reMarkable — see [../docs/architecture.md](../docs/architecture.md). Deploys as `/home/root/Writerdeck-server`: a single static ARMv7 binary, no on-device dependencies, always-on once installed.

## What it does
> No `/dev/uinput` on this kernel (it can't load it — exports trimmed; see [../DONE.md](../DONE.md)). Writerdeck-server feeds a patched Writerdeck over a local socket (`/run/Writerdeck.sock`) instead.

- Serves the capture page + WebSocket on `:8000` — `index.html` is embedded via `go:embed` (nothing extra to deploy). The browser sends key events; the server forwards them to Writerdeck as integer Unicode codepoints (`{"t":"text","cp":N}`, escaping-proof) plus named keys/modifiers — Writerdeck takes Qt input, so there is no fd to swap; it replays decoded keys.
- Notes file-manager API on `/api/notes` (list / read / create / rename / delete over `/home/root/Writerdeck-user-documents`), gated by a per-boot PIN (`/api/pin` → HttpOnly `writerdeck_token` session cookie that also guards the WS upgrade).
- Settings API on `/api/settings` and rotate on `POST /api/rotate` — font, PIN length, and display rotation persist to `/home/root/.Writerdeck/settings.json`; rotation is pushed to Writerdeck on connect (`setrotation`) and after phone rotate.
- Supervisor / lifecycle split — owns the `xochitl ↔ Writerdeck` toggle in Go: keeps serving `:8000` even under the stock GUI and summons Writerdeck on demand (`/api/open`, `/api/launch`); boot auto-launches one editor session; the physical Home button relays through to Writerdeck (two-level Home: edit→Lobby, Lobby→quit→xochitl). **USB Escape** (evdev hotplug watch) launches Writerdeck to the Lobby from stock UI when no editor session is active — not a wake-from-sleep path (power button only). Pushes Lobby info to Writerdeck on socket connect (`IP`, PIN, sync flags, note count, formatted last sync) via `{"t":"info",…}`; `POST /api/sync/ack` (after phone reconcile) stores `lastSyncAt` and re-pushes. Session lines (`editor started`, `editor process exited`, `home button -- relaying to editor`) go to stderr → `journalctl -u writerdeck.service` when run under systemd.

## Layout
```
daemon/
  go.mod
  main.go        Writerdeck-server: flags, WebSocket + HTTP handlers, notes API, PIN/session, supervisor, editor-feed socket
  index.html     capture page + phone file-manager (keydown/keyup → WebSocket), embedded via go:embed
  app.js         browser UI logic (embedded)
```

## Build & deploy
Static ARMv7 (`CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`). Deploys originate on the Mac — the only host that reaches the tablet:
```bash
bash ../scripts/deploy-rmkbd.sh      # cross-build → /home/root/Writerdeck-server  (rmkw = binary-only Writerdeck redeploy)
```
The ThinkPad can cross-build for a compile check (`../scripts/build-rmkbd.ps1`) but can't reach the device.

## Troubleshooting

**Edit from browser → stock UI reloads in one beat** — Writerdeck started then exited; `session.end()` restarted `xochitl`. Check `journalctl -u writerdeck.service` for `QQmlApplicationEngine failed to load component` (broken QML patch) or `editor started` immediately followed by `editor process exited`. Rebuild and redeploy Writerdeck (`rmkw` after CI). Automated check: `bash ../scripts/test-edit-session.sh` (see [decisions.md](../docs/decisions.md) #21).

**Do not `pkill -f /home/root/Writerdeck`** — that pattern also matches `Writerdeck-server`. Use `pidof Writerdeck` for the editor; `pkill -f /home/root/Writerdeck-server` for the server.
