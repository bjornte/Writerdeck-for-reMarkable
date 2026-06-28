# daemon/ — `rmkbd`

The Go daemon that runs on the reMarkable — Component A of the architecture (see [../docs/architecture.md](../docs/architecture.md)). A single static ARMv7 binary, no on-device dependencies, always-on once installed.

## What it does
> No `/dev/uinput` on this kernel (it can't load it — exports trimmed; see [../DONE.md](../DONE.md)). The daemon feeds a patched keywriter over a local socket (`/run/rmkbd.sock`) instead.

- Serves the capture page + WebSocket on `:8000` — `index.html` is embedded via `go:embed` (nothing extra to deploy). The browser sends key events; the daemon forwards them to keywriter as integer Unicode codepoints (`{"t":"text","cp":N}`, escaping-proof) plus named keys/modifiers — keywriter takes Qt input, so there is no fd to swap; it replays decoded keys.
- Notes file-manager API on `/api/notes` (list / read / create / rename / delete over `/home/root/edit`), gated by a per-boot PIN (`/api/pin` → HttpOnly session cookie that also guards the WS upgrade).
- Supervisor / lifecycle split — owns the `xochitl ↔ keywriter` toggle in Go: keeps serving `:8000` even under the stock GUI and summons keywriter on demand (`/api/open`, `/api/launch`); boot auto-launches one editor session; the physical Home button relays through to keywriter (two-level Home). Pushes the Lobby's IP + PIN to keywriter on socket connect.

## Layout
```
daemon/
  go.mod
  main.go        the daemon: flags, WebSocket + HTTP handlers, notes API, PIN/session, supervisor, editor-feed socket
  index.html     capture page + phone file-manager (keydown/keyup → WebSocket), embedded via go:embed
```

## Build & deploy
Static ARMv7 (`CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`). Deploys originate on the Mac — the only host that reaches the tablet:
```bash
bash ../scripts/deploy-rmkbd.sh      # cross-build + ship to /home/root/  (rmkw = binary-only redeploy)
```
The ThinkPad can cross-build for a compile check (`../scripts/build-rmkbd.ps1`) but can't reach the device.
