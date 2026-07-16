# What's shipped

Writerdeck for reMarkable 1 adds keyboard word processing to the first-gen reMarkable. Type from a phone (Bluetooth keyboard over Wi-Fi) or USB via an OTG cable. Notes save as Markdown.

Open work: [TODO.md](TODO.md). How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

## Core loop

Power on and the tablet boots into a distraction-free editor. Open `http://<tablet-ip>:8000/` on the phone for the keyboard bridge; open notes on the tablet Files tab. Keystrokes travel WebSocket to Writerdeck-server, which feeds Writerdeck over `/run/Writerdeck.sock`. No app to install on the phone.

Home from the editor saves and returns to the Lobby on the Files tab with the note you were editing selected. Home from the Lobby quits to the stock reMarkable UI. Writerdeck-server keeps serving port 8000 either way — relaunch from the tablet Files tab, Esc on a USB keyboard, or left and right page buttons together from stock UI.

## Lobby

Six-tab pager on e-ink: Files, Keyboard, Sync, Settings, Shortcuts, Home (digits 1–6). Boot and Home from edit open on Files. Touch the tabs or use Tab, arrows, or digits 1–6. On socket connect the server pushes IP, PIN, sync state, note count, and formatted last sync; it re-pushes when wlan0 gets an address, a reconcile finishes, or notes change.

The Files tab lists notes from the server over a trusted socket. New, Edit, Read, Rename, and Delete work by touch or USB keys (`n`, Enter, `v`, `r`, `d`). With private notes on, a second row offers Encrypt, New encrypted, or Decrypt by touch (USB `x` / `y` on Files). Edit opens the note in type mode; Read opens preview on e-ink. A second tap on an already-selected row opens Edit. Lobby new/rename uses an inline prompt with a movable cursor (arrow keys, Home/End). Failed encrypt or decrypt shows a red error on the Files tab. After Home from edit, Lobby keyboard focus stays on `lobbyFocus` so USB and WebSocket keys keep working. Home from read returns to the Lobby (not quit). Show PIN on tablet (phone button) drops back to the Lobby when a second device needs the PIN.

Launch from stock UI: Mac `wd` or `bash scripts/lobby.sh`; on tablet SSH, `~/wd`; USB Esc; L+R page buttons. Ctrl-K note switcher works from the editor on USB keyboard (device verified).

## Phone companion

Upload and download `.md` notes from the note list. Type mode when the tablet opens a note for edit (Files, Ctrl-K, or reconnect) via WebSocket `openedit` — keystrokes to e-ink with an echo footer. When the tablet is in read preview, new/rename, or delete confirm, the server broadcasts `openread` or `lobbyinput` so a Bluetooth keyboard paired to the phone still forwards keys (read view or a green banner over the list). **Paste from phone** in Type mode inserts clipboard text at the tablet cursor (keystroke replay, not a new file). Home on the tablet drops the phone back to the list. Dark type mode for OLED phones. No phone preview, Edit, or file CRUD — tablet Lobby Files tab.

## Security

PIN on the tablet each boot — 6 digits, 4 digits, or none (none warns that anyone on your Wi-Fi can connect). Five wrong guesses from one IP lock that IP for 60 seconds. Auth cookie until 04:00 local time.

Optional private-note encryption: a second 6-digit vault PIN (tablet only), independent of the pairing PIN. Per-note Encrypt/Decrypt from Lobby Files; ciphertext as `.md.enc` with `WDENC1` on disk and on GitHub. PIN required for each open, read, edit, encrypt, or decrypt — no persistent unlocked state. Failed encrypt/decrypt surfaces an error on the Files tab (corrupt ciphertext, wrong format, name clash). Recovery material syncs to `secret/pin` and `secret/vault`. Phone download decrypts after the tablet enters its PIN; the phone waits. Details: [decisions.md](docs/decisions.md) §31.

## Settings and sync

Reading font, PIN length, display rotation, and Exit Writerdeck live on the tablet Lobby Settings tab. GitHub sync lives in Notes sync setup on the phone (bar: Sync setup) — toggle, repo, token, Save, Sync. Connection indicator refreshes via status every five seconds. Token in phone browser and tablet RAM only; auto-restore after restart via WebSocket `needtoken`; `syncOn`, `syncRepo`, and `syncMeta` on disk.

## Editor

Built from owned fork [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) of remarkable-keywriter (**keywriter**: Qt 5, C++, QML), deployed as Writerdeck. Socket input, Lobby, and Mac/Linux-style edit helpers live in the fork; `build-keywriter.sh` is CI glue ([decisions.md](docs/decisions.md) §3). Full-panel via linuxfb. Norwegian and other Unicode via the browser path; USB Norwegian (`no.qmap`) verified on hardware — æ ø å, AltGr, `@`, `{` `}`. Reading fonts: Inter, Literata, EB Garamond, DejaVu. Page up and down in read and edit. Block cursor hides while typing. Ctrl-K note switcher. Mac/Linux-style navigation in edit (Ctrl/Alt chords; partial — see [editor-testing/todo.md](docs/editor-testing/todo.md)). Physical Home is one path: server `EVIOCGRAB` on gpio-keys while editing ([decisions.md](docs/decisions.md) §28). Power button saves, shows a sleep screen, suspends; press again to wake (device verified). USB Esc from stock UI launches to Lobby. Rotation in settings, pushed on connect.

## GitHub sync

Optional, off by default. Writerdeck-server reconciles with a private repo — pull missing either way, push local-only, clashes keep both copies. Marker-aware delete when the local copy is pristine with stored sha. Empty-push guard refuses zero bytes over previously-synced content. Triggers on boot, three-minute poll, token verify, Sync now, Home, power, and CRUD. Skips only the note open on the tablet. Lobby Sync tab shows TOKEN NEEDED, SYNC OFFLINE, and SYNC FAILED. Token auto-restore after restart: [server-sync-implementation.md](docs/server-sync-implementation.md).

## Document integrity

Plain Markdown on disk, no silent overwrite of live edits, durable saves. Slices 1–11 shipped. Residual risks: [docs/integrity-audit.md](docs/integrity-audit.md).

## Infrastructure

Static Go binary at `/home/root/Writerdeck-server`. Writerdeck built in CI. Cold-boot autostart via `writerdeck.service`. Keep-awake during editor sessions only. On-device layout: [docs/architecture.md](docs/architecture.md).

Regression: `bash scripts/test-edit-session.sh` — POST `/api/open` must keep Writerdeck running, xochitl down, and editorActive true for several seconds. After arrow, selection, or `handleKey` QML changes: `bash scripts/test-keyboard-harness.sh --fast` (**110/110/0** for sign-off; [milestone-runs.md](docs/editor-testing/milestone-runs.md)). After Lobby or `handleHome` QML changes: `bash scripts/test-lobby-keyboard.sh` — lobby keys after return from edit, Home-from-read must not quit Writerdeck. After vault or encryption changes: `bash scripts/test-vault.sh` and `bash scripts/test-vault-e2e.sh` (sync on required for E2E).

`/dev/uinput` is dead on this kernel. The editor is fed over a socket. Do not retry uinput.
