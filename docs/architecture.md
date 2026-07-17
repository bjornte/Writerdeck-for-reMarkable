# Architecture

How the system works, the device facts, and the dev loop. Open work: [../TODO.md](../TODO.md). Why: [decisions.md](decisions.md). Shipped: [../DONE.md](../DONE.md). Gotchas: [lessons.md](lessons.md).

---

## What this is

The reMarkable 1 has a large e-ink screen and a distraction-free OS, but no word processor and no native keyboard support. Writerdeck adds both. You type on an iPhone keyboard (Bluetooth, bridged over Wi-Fi) or a USB keyboard (OTG cable). The tablet runs the editor and saves Markdown.

```
 CLIENT (iPhone Safari + BT keyboard; Mac in dev)
     |  WebSocket over LAN
     v
 reMarkable 1
   Writerdeck-server — capture page, WebSocket, file API
   Writerdeck — patched keywriter (Qt 5 / C++ / QML), reads /run/Writerdeck.sock, saves .md
```

Keystrokes reach the editor through a local socket, not `/dev/uinput`. This kernel cannot load uinput; see [decisions.md](decisions.md) §1.

**Who does what:** Two programs on the tablet, plus the phone page. **Writerdeck-server** is Go compiled ahead of time into one static ARM binary — always-on, no screen: Wi‑Fi, files, sync, PIN, launching the editor. **Writerdeck** is the full-screen Qt app (QML for what you see and typing feel; C++ for startup, display, socket keys). **Phone** is HTML/JS served by the daemon. They meet over `/run/Writerdeck.sock`. “On the reMarkable” does not mean “must be QML.”

## Document integrity

Writerdeck exists to produce durable Markdown notes. Before shipping anything that touches notes, saves, opens, sync, CRUD, or editor lifecycle, ask: can this lose text, write wrong bytes, or overwrite without the user knowing?

Files on disk are UTF-8 Markdown — reject qrichtext and HTML saves. While a note is open for edit, reconcile and remote CRUD must not silently overwrite it. Saves use defined paths, 45-second autosave while editing, and save-before-stop on deploy and SIGTERM. Note writes use temp-then-rename. If disk changes under an open session, reload or surface conflict — do not save a stale buffer over fresh disk. GitHub sync assists backup; it must not trump a live edit or push empty or HTML over good content.

Slices 1–11 are shipped. Deploy both server and Writerdeck for loopback save and `autosavenow`. Policy: [decisions.md](decisions.md) § Document integrity. Residuals: [integrity-audit.md](integrity-audit.md).

## On the tablet

Under `/home/root/`:

`Writerdeck-server` is the always-on Go daemon — WebSocket, HTTP API, session lifecycle, GitHub sync engine. Source in `daemon/`.

`Writerdeck` is the full-screen Markdown editor — our patched build of upstream [keywriter](https://github.com/dps/remarkable-keywriter) (*remarkable-keywriter*). Keywriter is the editor engine: a **Qt 5** app written in **C++** and **QML**. Built in CI from `third_party/keywriter/`.

**C++ vs QML (this project):** Both live in the editor kernel (keywriter / our fork). **QML** is the screen and caret application — layout, Lobby, selection, goal-column state, applying dispatch results (`main.qml`, `edit_mac_helpers.qml.inc`). **C++** is the engine under that — starting the app, talking to the tablet display, feeding keystrokes from our socket into Qt (`main.cpp`), and pure text math / undo / key-chord dispatch / visual-line math in `EditHelper` (migration 2 Phase A @ `a92ad2b`, Phase B @ `57bfc21`, Phase C @ `6a15e08`). Day-to-day wrap fixes still land in QML apply paths; string-math, undo, chord mapping, and visual-line walk live in C++ per [editor-migration-2-to-cpp/](editor-migration-2-to-cpp/). Keep the hand-tuned wrap gaps and custom `EditHelper` undo stacks ([decisions.md](decisions.md) §30).

Example: a careful, long-lived **undo** improvement belongs in the owned fork as **C++ `EditHelper`** — not in the emergency patch script. Reach for a deeper C++ text engine only if undo had to be rebuilt inside Qt’s document model itself.

This patch pipeline is intentionally reproducible. Most edit-mode behavior lives in the Writerdeck fork [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter); `third_party/keywriter/build-keywriter.sh` is clone + assert + build only (committed fork `main.qml` already includes helpers and Lobby). Keep new editor behavior in the fork ([decisions.md](decisions.md) §3, [todo-handoff-keywriter-fork.md](editor-migration-1-to-QML/todo-handoff-keywriter-fork.md)). Pure text math and undo live in fork C++ `EditHelper` ([todo-handoff-edit-helper-cpp.md](editor-migration-2-to-cpp/todo-handoff-edit-helper-cpp.md)). If e-ink full-frame redraw becomes a problem later, see [improvements.md](improvements.md) § E-ink redraw (dirty-region ideas from yaft — not a terminal editor swap).

`Writerdeck-launcher.sh` sets Qt and e-ink launch environment; the server spawns Writerdeck with `--editor`.

`Writerdeck-user-documents/` holds your `.md` notes.

`.Writerdeck/settings.json` persists font, PIN mode, sync flags, rotation, USB keyboard layout.

`qt5/` is the Qt runtime (~14 MB).

`/run/Writerdeck.sock` is the Unix socket between server and editor.

`writerdeck.service` is the systemd unit in `/etc/systemd/system/`.

The client is embedded HTML and JavaScript inside the server binary — no separate deploy.

## The companion appliance

The tablet is the web server. Open Safari on the phone; no app to install. Phone file-manager and preferences dedup is done — upload, download, paste-at-cursor in Type mode, and sync token entry stay on the phone. USB keyboard layout and device settings (font, PIN, rotation, exit) live on the tablet Lobby. Split: [browser-vs-tablet.md](browser-vs-tablet.md).

Writerdeck-server keeps `:8000` up even under stock xochitl. It summons Writerdeck on demand: stop xochitl, spawn Writerdeck, and on Home from the Lobby restart xochitl while the server keeps serving. Boot auto-launches one editor session. Notes live in `Writerdeck-user-documents/`; the Go API rejects unsafe paths via `notesSafe()`.

A random PIN is minted per boot and shown in the Lobby. The phone POSTs it for an HttpOnly session cookie that gates the notes API and WebSocket. Length is owner-choosable: 6, 4, or none. Per-IP lockout backs the PIN modes.

Notes sync setup on the phone (bar: Sync setup) covers GitHub sync — toggle, repo, token, Save, and Sync. The Lobby is a six-tab pager — Files, Keyboard, Sync, Settings, Shortcuts, Home (1–6) — fed by socket info messages with IP, PIN, sync state, note count, and last sync. Boot and Home from edit land on Files; Home (tab 6) is the welcome screen. Settings on the tablet covers reading font, PIN length, display rotation, and Exit Writerdeck. Files CRUD goes through a trusted socket, not unauthenticated LAN HTTP. Launch from stock UI: USB Escape, left and right page buttons together, Mac `wd`, tablet `~/wd`.

Browse mode on the phone shows a note list for upload and download only — no key capture by default, no preview, no phone-initiated Edit. **Type mode** forwards keystrokes when the tablet opens a note for edit; **remote key mode** forwards when the tablet is in read preview (`openread`) or a Lobby Files prompt (`lobbyinput` — new, rename, delete confirm). Paste from phone replays clipboard text at the cursor (modal in `notes-ui.js`, not upload). The server broadcasts `openedit`, `openread`, or `lobbyinput`; the phone mirrors unless an overlay has focus. Home on the tablet broadcasts `exitedit` so the phone drops back to the list.

GitHub sync is optional and off by default. The server is the engine; the token stays in RAM (browser `localStorage` + tablet RAM, never disk). After a service restart the server asks connected browsers for the token via WebSocket `needtoken`; see [server-sync-implementation.md](server-sync-implementation.md).

IP is detected from wlan0 and re-pushed to the Lobby when it changes.

## Constraints

No jailbreak; preserve OTA — so no Toltec. No on-device runtime beyond one static Go binary (`CGO_ENABLED=0`, ARMv7). Markdown on disk; HTML there is a bug. Device files are ASCII and LF — see `.gitattributes`. Keep the tablet awake; it drops Wi-Fi on suspend. Latency is e-ink refresh, not the LAN. Micro-USB OTG accepts a plain USB keyboard if Wi-Fi stalls.

---

## Device facts

reMarkable 1, codename zero-gravitas. OS `20260506100933`, kernel `5.4.70-v1.6.3-rm10x`.

SSH over Wi-Fi: `ssh root@<tablet-ip>`. Set `RM_HOST_WIFI` in `secrets/remarkable.local.env`; reserve the tablet's MAC on the router so DHCP stays stable. If deploy/SSH fails on home Wi-Fi, check whether the Mac and tablet are on an iPhone Personal Hotspot — tablet is often `172.20.10.5` (`export RM_HOST=172.20.10.5`; phone UI `http://172.20.10.5:8000/`). USB at `10.11.99.1` is dead on the Mac. SSH password is on the device settings screen and regenerates after every OTA — gitignored in secrets.

Notes: `/home/root/Writerdeck-user-documents/`. Deploy the binary to `/home/root/Writerdeck`, not into the notes directory.

Buttons on `/dev/input/event1` (value 1 = press): Home 102, left 105, right 106, power 116. Server watches them always; while a Writerdeck session is active it takes exclusive `EVIOCGRAB` so Qt never sees gpio Home ([decisions.md](decisions.md) §28). Idle: no grab — L+R chord launches, xochitl still gets buttons. Power sleep/wake and USB Escape idle launch as before.

Rootfs is about 96% full (~228 MB), but nothing we ship goes there. Binary, Qt sysroot, and notes live on `/home/root/`, a separate multi-GB partition. Do not resize rootfs.

References: [remarkable-keywriter](https://github.com/dps/remarkable-keywriter) · build with `ghcr.io/toltec-dev/qt:v3.3` (not `latest` — that is Qt6) · [input docs](https://remarkable.guide/devel/device/input.html) · [crazy-cow HWR fallback](https://github.com/machinelevel/sp425-crazy-cow) · [libremarkable](https://github.com/canselcik/libremarkable)

---

## Build and deploy

Writerdeck-server from the Mac:

```bash
bash scripts/deploy-rmkbd.sh
```

Writerdeck from CI — QML is baked into the binary; rebuild before deploy:

```bash
git push && bash scripts/fetch-keywriter-dist.sh && bash scripts/deploy-keywriter.sh -b
```

`deploy-keywriter.sh` only pushes `dist/Writerdeck`; it does not rebuild. After QML edits, fetch or Docker-build first, relaunch the editor, check `journalctl -u writerdeck`. Requires Go on the Mac (`brew install go`).

`deploy-rmkbd.sh` calls `POST /api/flush-save`, then SIGTERM-waits up to about twelve seconds before replacing the binary. Follow with `systemctl start writerdeck`.

Aliases via `bash scripts/install-alias.sh`: `rmkw` for binary-only Writerdeck deploy; `test-edit-session.sh` after QML changes ([decisions.md](decisions.md) §21); `test-keyboard-harness.sh` after arrow/selection QML (§22); `test-e2e.sh -s` for the full pipeline without rebuilding the server; `rmpush` to commit and push. Deploy uses gzip-over-ssh, not scp ([decisions.md](decisions.md) §12).

After Writerdeck or QML changes, run `test-edit-session.sh`; add `test-keyboard-harness.sh --fast` when `handleKey` or selection logic changed (110/110/0 sign-off, [editor-testing/milestone-runs.md](editor-testing/milestone-runs.md)); add `test-lobby-keyboard.sh` when Lobby navigation, `handleHome`, or `lobbyFocus` changed ([decisions.md](decisions.md) §29). After server or embedded JS only, restart the server and smoke-test the API or browser. After both, do both. Iterate over Wi-Fi; logs in `journalctl -u writerdeck.service`.
