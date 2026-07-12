# Architecture decisions

Why the project is built this way. How it works: [architecture.md](architecture.md). Open work: [../TODO.md](../TODO.md). What's shipped: [../DONE.md](../DONE.md). Gotchas: [lessons.md](lessons.md).

---

## Document integrity

Writerdeck is a typewriter. The owner's prose must survive editing, sync, and normal device use as plain Markdown on disk. That contract gates every feature — not a polish item for later. Full audit: [integrity-audit.md](integrity-audit.md).

Files on disk are UTF-8 Markdown, never Qt HTML or rich text. While a note is open on the tablet, reconcile, phone CRUD, and rename/delete must not silently overwrite it. Saves go through defined paths, with 45-second autosave and save-before-stop on deploy and SIGTERM; a SIGKILL or crash before the next autosave or flush can still lose recent typing. If disk changes under an open session, the user must reload or get conflict UX — the buffer must not blindly win. GitHub sync assists backup; it must not delete, empty-push, or fork paths against a live edit.

No change to `daemon/`, the sync engine, `build-keywriter.sh`, or note APIs ships without an integrity pass against those rules.

## Device verification

A green deploy script is not a pass. After any change to `daemon/`, `build-keywriter.sh`, `lobby/`, or deploy scripts: rebuild Writerdeck when QML changed, deploy, relaunch, then read `journalctl -u writerdeck -n 30` on the tablet. Fail on QML parse errors or immediate `editor process exited`. When Writerdeck QML was touched, run `bash scripts/test-edit-session.sh`. "Works in theory" does not ship.

## Move features to reMarkable lobby UI

Over time, move controls from the phone to the tablet Lobby and remove duplicates once the tablet path is verified

---



## 1. Socket input, not uinput

This kernel cannot load `/dev/uinput`. Open returns ENODEV; trimmed kernel exports (`CONFIG_TRIM_UNUSED_KSYMS`) mean an out-of-tree `uinput.ko` cannot bind. The patched editor reads keystrokes from `/run/Writerdeck.sock` instead. Do not retry uinput. If keywriter itself ever fails, fallbacks are a libremarkable editor or HWR via the Wacom pen node — that draws ink, not clean Markdown.

## 2. Synthetic QKeyEvent injection

Keywriter takes keyboard input through Qt QPA — there is no input fd to swap. A socket reader in `main.cpp` posts synthetic `QKeyEvent`s to `focusWindow()`. The browser sends full JSON to the daemon; the daemon sends integer Unicode codepoints to keywriter as `{"t":"text","cp":N}`. The keymap lives browser-side, where layout is already resolved.

## 3. Build keywriter from source

The four-year-old prebuilt binary dies at the loader (`libQt5Quick.so.5`). Qt is static-linked into xochitl on current firmware, so `LD_LIBRARY_PATH` cannot rescue it. We cross-build from source in `ghcr.io/toltec-dev/qt:v3.3` (CI) and deploy a Qt5 runtime sysroot. Writerdeck renders via linuxfb on the rM1's real `/dev/fb0` — rm2fb is not needed.

## 4. No Toltec

Toltec locks firmware to a supported version range and can soft-brick on unsupported versions. That conflicts with preserving OTA updates. Only revisit if the owner accepts the version lock.

## 5. Static Go binary

The server is `CGO_ENABLED=0`, ARMv7, with no on-device runtime dependencies. It survives firmware updates and cross-compiles on the Mac in about a second. The tablet ships no Python; installing it implies Entware/Toltec and a firmware lock.

## 6. Mac builds the server; CI builds Writerdeck

The Mac is the only host that can reach the tablet over Wi-Fi, so deploys originate there. Go cross-compiles fast enough that a local `go build` is the edit-deploy loop. CI and Docker are reserved for Writerdeck, which needs the toltec Qt sysroot.

## 7. Always-on server, on-demand editor

Writerdeck-server keeps serving `:8000` even under the stock xochitl GUI. It toggles xochitl per editor session, not per server lifetime — Home returns to the GUI while the phone can still reach the server. Boot auto-launches one editor session. Keep-awake wraps only the Writerdeck child via `systemd-inhibit`, so the tablet sleeps normally under the GUI.

Any systemd unit whose `ExecStart` lives on `/home/root` must declare `RequiresMountsFor=/home/root`. Otherwise a cold boot races the mount and you get `203/EXEC`.

## 7b. USB Escape launches from stock UI

Writerdeck-server watches USB keyboard evdev nodes with hotplug rescan. Escape while idle — xochitl up, no editor session, not sleeping — starts a session to the Lobby. This is not Esc-to-wake after power sleep; the power button handles that. While a session is active, Escape is ignored here so Writerdeck keeps normal edit behaviour.

## 7c. Left and right page buttons launch from stock UI

Physical page buttons (`KEY_LEFT` / `KEY_RIGHT` on `/dev/input/event1`) are readable alongside xochitl. Both pressed while idle follows the same launch path as USB Escape, with 800 ms debounce. xochitl still receives individual button events — acceptable on the home screen; in a document the chord may briefly page.

## 8. Companion model

The tablet is the web server; there is no phone app. Writerdeck-server does all file operations on `Writerdeck-user-documents/` natively. Writerdeck changes are a Lobby overlay and opening a note via `saveAndLoad(name)`. A random PIN is minted per boot and shown on e-ink; you must hold the device to read it. Two-level Home — edit to Lobby, Lobby to quit — means you can write again without rebooting.

## 9. Share is Download plus Copy

The native iOS share sheet needs HTTPS, which we do not have on plain LAN http. Download with `Content-Disposition: attachment` and copy-to-clipboard are the reliable paths. `navigator.clipboard` also needs a secure context, so Copy falls back to a temporary textarea and `document.execCommand('copy')` on plain http.

## 10. Two-machine dev split (retired)

July 2026. Was a Mac-on-LAN plus work-laptop-over-git workaround when corporate VPN kept the work machine off the home LAN. All dev is Mac-on-LAN now.

## 11. Wi-Fi is the dev path

The Mac's USB-ethernet gadget to the tablet is inactive. Wi-Fi SSH works when `RM_HOST_WIFI` in secrets matches the tablet's DHCP address.

## 12. Deploy transport is gzip-over-ssh

`scp` wedges at a fixed ~255 KB offset on the Mac-to-Wi-Fi-to-tablet link — an SFTP windowing deadlock, not sleep or QoS. `rm_send_file` in `_env.sh` streams gzip through SSH with a post-copy size check. When scp stalls at a fixed offset on an embedded link, switch transports.

## 13. Secrets in a gitignored env file

Plaintext in `secrets/remarkable.local.env` is acceptable here: the SSH password is shown on the device settings screen, and the device lives on a home LAN. The real risk is git leakage, which the ignore prevents.

## 14. LF and ASCII on device files

CRLF or a stray non-ASCII byte breaks shell scripts and the systemd unit on the device. `.gitattributes` normalizes line endings. Markdown prose may use Unicode; code and device files may not.

## 16. Upload reuses the safe create route

Uploading a `.md` from the phone POSTs through `/api/notes` so it inherits `notesSafe()` path checks and 409-on-exists. The authoritative body size cap is `http.MaxBytesReader` on the server; client checks are UX only.

## 17. PIN and per-IP lockout

The owner chooses 6-digit, 4-digit, or no PIN. "None" opens the notes API to anyone on the Wi-Fi — the UI warns explicitly. Per-IP lockout blocks an address for 60 seconds after five wrong guesses, not globally, so an attacker cannot lock the owner out of their own device. Compare uses `subtle.ConstantTimeCompare`. Changing PIN length re-mints under lock and re-pushes the Lobby. Store length as the string enum `"6"`, `"4"`, or `"none"` — an absent field on an older settings file would read as integer zero and silently mean "none". Gotcha: `int(Uint32) % N` overflows the device's 32-bit `int`; reduce in uint32 space first.

## 18. Show PIN on tablet

A second device that arrives mid-edit sees the note, not the PIN. `POST /api/lobby` saves the open note and drops to the Lobby so the PIN is readable on e-ink. The endpoint is pre-auth by necessity; it reveals nothing over the wire. Rate-limited to about one honored call per three seconds. `showLobby()` is idempotent and never quits — distinct from `handleHome()`.

## 19. GitHub sync is a non-authoritative reconciler

The engine runs on Writerdeck-server. The token lives in browser `localStorage` and tablet RAM via `POST /api/sync/token`; it is never written to disk. Reconcile unions tablet and repo note lists and copies anything missing from either side — it never deletes on its own. Destructive ops from the browser pair to GitHub. External deletes on GitHub propagate when the local copy is pristine and carries a stored `sha`, confirmed with a per-note 404 before acting; unpushed local edits resurrect instead of deleting. Empty-push guard refuses to push a zero-byte file over a previously-synced note. Details: [server-sync-implementation.md](server-sync-implementation.md).

## 20. Display rotation persists in settings

Global rotation (0, 90, 180, 270) is stored in `.Writerdeck/settings.json`. Phone rotate POSTs to the server, which pushes `setrotation` on connect. USB Ctrl+arrow in preview relays `rotationChanged` back via `rotation_watcher`. Both binaries must be current — server-only deploy can save to disk while an old Writerdeck ignores `setrotation`.

## 21. Edit-session regression test

Phone Edit (`POST /api/open`) is the primary companion launch path. If Writerdeck exits immediately, the session ends, xochitl restarts, and it looks like the stock UI reloading — almost always a broken QML patch, not the server. `scripts/test-edit-session.sh` POSTs `/api/open` from stock UI and asserts Writerdeck stays up about eight seconds. Run it after Writerdeck or QML deploy (`rmkw`); not after server-only deploy — restart the server and spot-check the API instead. `build-keywriter.sh` asserts brace balance in `handleKey()` before write.

## 22. Keyboard selection harness

Modifier+arrow and selection behaviour are tested on the device via `daemon/cmd/edit-harness` and `scripts/test-keyboard-harness.sh`, not by reading saved note bytes. Writerdeck publishes cursor/selection over the socket; the server exposes `GET /api/test/editor-state` and `POST /api/test/reset` (hard quit). Scenarios send keys over `/ws` like the phone UI.

Default run uses soft reset: one editor launch per full suite, `PUT` plus `POST /api/reload` between scenarios so content changes are not overwritten by `saveAndLoad`. Unit coverage for `translate()` modifier masks stays in `daemon/editor_test.go`. The harness does not exercise USB evdev or `.qmap` — keep hardware checks for Norwegian Alt+arrow, AltGr, and national characters. Run after QML arrow/selection changes and after daemon test-handler changes.

## 23. On-device Writerdeck naming

On the tablet: `Writerdeck`, `Writerdeck-server`, `Writerdeck-user-documents/`, `.Writerdeck/`, `/run/Writerdeck.sock`, `writerdeck.service`. Repo script names like `deploy-rmkbd.sh` stay historical. `migrate-device-layout.sh` renames legacy paths on deploy.

## 24. Tablet file CRUD via trusted socket

The Lobby Files tab sends `{"t":"req","op":…}` over the existing Unix socket; the server performs the same disk ops as `/api/notes`. Tablet delete, rename, and create queue sync and notify the phone via `tabletcrud`. Launch the Lobby from the Mac with `wd` or `bash scripts/lobby.sh`, or on the tablet with `~/wd`.

## 25. Display sync after Lobby clears the buffer

Returning to the Lobby assigns `query.text = ""`, which breaks any `text:` binding on the `TextEdit`. `doLoad` and `toggleMode()` call `syncQueryDisplay()` so edit shows plain `doc` and preview shows `readHtml(doc)` — never RichText extracted back into `doc`. Without this, Home save can zero files. `showLobby` also clears `currentFile`.

---



## Open risks

Firmware OTA may wipe the systemd unit and regenerate the SSH password — recovery is re-deploy and re-enable. USB keyboard locales need qmaps for national layouts; the browser path already resolves Norwegian via the phone OS. Encrypted note subset is design-only. Integrity residuals: [integrity-audit.md](integrity-audit.md). uinput is closed — see decision 1. Go must be on the Mac. Rootfs is about 96% full; everything we ship lives on `/home/root/`. Do not resize rootfs — A/B OTA scheme, brick risk.