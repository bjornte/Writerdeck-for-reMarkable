# daemon/ — Writerdeck-server

Go daemon on the reMarkable. Deploys as `/home/root/Writerdeck-server` — static ARMv7, always-on. Architecture: [../docs/architecture.md](../docs/architecture.md).

No `/dev/uinput` on this kernel — feeds Writerdeck over `/run/Writerdeck.sock`.

- Capture page + WebSocket on `:8000`; forwards keys as `{"t":"text","cp":N}`.
- WebSocket `openedit` / `exitedit` — phone Type mode when tablet opens/closes a note (`followTabletOpen` in `notes-ui.js`).
- Notes API `/api/notes` on `Writerdeck-user-documents/` (list, upload create, download); PIN auth. Phone DELETE/PATCH removed — tablet CRUD via socket.
- Settings `/api/settings`.
- xochitl ↔ Writerdeck lifecycle; USB Escape and L+R page-button launch when idle.
- GitHub sync engine (`syncengine.go`); WebSocket `needtoken` when tablet RAM lacks a token; test endpoints `/api/test/*`; `cmd/edit-harness/` and `cmd/lobby-keyboard-test/` for device regression tests.

Embedded phone UI (`index.html`, `notes-ui.js`, …): note list with Upload/Download; Type mode + paste modal; no preview or phone Edit. See [../docs/browser-vs-tablet.md](../docs/browser-vs-tablet.md).

## Layout

```
main.go, editor.go, editorstate.go, input.go, websocket.go, notes.go, settings.go,
auth.go, lobby.go, handlers.go, session.go, testhandlers.go,
syncengine.go, syncgithub.go, syncapi.go, syncmeta.go,
cmd/edit-harness/,
index.html, app.js, connection.js, notes-ui.js, panels.js,
deps.js, state.js, sync.js, app.css  (embedded)
```

## Build

```bash
bash ../scripts/deploy-rmkbd.sh
```

## Troubleshooting

Stock UI flash after open: Writerdeck exited — check journalctl for QML errors. Run `test-edit-session.sh`; after Lobby/`handleHome` changes also `test-lobby-keyboard.sh`.

Do not `pkill -f /home/root/Writerdeck` — matches server. Use `pidof Writerdeck`.
