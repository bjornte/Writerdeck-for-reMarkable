# daemon/ — Writerdeck-server

Go daemon on the reMarkable. Deploys as `/home/root/Writerdeck-server` — static ARMv7, always-on. Architecture: [../docs/architecture.md](../docs/architecture.md).

No `/dev/uinput` on this kernel — feeds Writerdeck over `/run/Writerdeck.sock`.

- Capture page + WebSocket on `:8000`; forwards keys as `{"t":"text","cp":N}`.
- Notes API `/api/notes` on `Writerdeck-user-documents/`; PIN auth.
- Settings `/api/settings`; rotation `POST /api/rotate`.
- xochitl ↔ Writerdeck lifecycle; USB Escape and L+R page-button launch when idle.
- GitHub sync engine (`syncengine.go`).

## Layout

```
main.go, editor.go, input.go, websocket.go, notes.go, settings.go,
auth.go, lobby.go, handlers.go, session.go,
syncengine.go, syncgithub.go, syncapi.go, syncmeta.go,
index.html, app.js, connection.js, notes-ui.js, panels.js,
deps.js, state.js, sync.js, app.css  (embedded)
```

## Build

```bash
bash ../scripts/deploy-rmkbd.sh
```

## Troubleshooting

Edit → stock UI flash: Writerdeck exited — check journalctl for QML errors. Run `test-edit-session.sh`.

Do not `pkill -f /home/root/Writerdeck` — matches server. Use `pidof Writerdeck`.
