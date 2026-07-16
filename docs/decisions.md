# Architecture decisions

Why the project is built this way. How it works: [architecture.md](architecture.md). Open work: [../TODO.md](../TODO.md). What's shipped: [../DONE.md](../DONE.md). Gotchas: [lessons.md](lessons.md).

---

## Document integrity

Writerdeck is a typewriter. The owner's prose must survive editing, sync, and normal device use as plain Markdown on disk. That contract gates every feature — not a polish item for later. Full audit: [integrity-audit.md](integrity-audit.md).

Files on disk are UTF-8 Markdown, never Qt HTML or rich text. While a note is open on the tablet, reconcile and tablet-side rename/delete must not silently overwrite it. Saves go through defined paths, with 45-second autosave and save-before-stop on deploy and SIGTERM; a SIGKILL or crash before the next autosave or flush can still lose recent typing. If disk changes under an open session, the user must reload or get conflict UX — the buffer must not blindly win. GitHub sync assists backup; it must not delete, empty-push, or fork paths against a live edit.

No change to `daemon/`, the sync engine, `build-keywriter.sh`, or note APIs ships without an integrity pass against those rules.

## Device verification

A green deploy script is not a pass. After any change to `daemon/`, `build-keywriter.sh`, `lobby/`, or deploy scripts: rebuild Writerdeck when QML changed, deploy, relaunch, then read `journalctl -u writerdeck -n 30` on the tablet. Fail on QML parse errors or immediate `editor process exited`. When Writerdeck QML was touched, run `bash scripts/test-edit-session.sh`. "Works in theory" does not ship.

## Move features to reMarkable lobby UI

Shipped July 2026. Phone file-manager dedup: create, read preview, Edit, rename, and delete removed from the browser; tablet Files tab owns day-to-day file ops. Phone keeps upload, download, paste-at-cursor (Type mode), and sync token entry — [browser-vs-tablet.md](browser-vs-tablet.md). Preferences dedup (July 2026): reading font, PIN length, and Exit Writerdeck moved to the Lobby Settings tab; phone Preferences removed.

---



## 1. Socket input, not uinput

This kernel cannot load `/dev/uinput`. Open returns ENODEV; trimmed kernel exports (`CONFIG_TRIM_UNUSED_KSYMS`) mean an out-of-tree `uinput.ko` cannot bind. The patched editor reads keystrokes from `/run/Writerdeck.sock` instead. Do not retry uinput. If keywriter itself ever fails, fallbacks are a libremarkable editor or HWR via the Wacom pen node — that draws ink, not clean Markdown.

## 2. Synthetic QKeyEvent injection

Keywriter takes keyboard input through Qt QPA — there is no input fd to swap. A socket reader in `main.cpp` posts synthetic `QKeyEvent`s to `focusWindow()`. The browser sends full JSON to the daemon; the daemon sends integer Unicode codepoints to keywriter as `{"t":"text","cp":N}`. The keymap lives browser-side, where layout is already resolved. Edit bindings are Mac/Linux-style: Control and Alt chords (Meta from a Mac-like phone keyboard maps to Ctrl). The tablet being Linux does not require a separate “Linux-only” shortcut set — USB Linux keyboards already use Ctrl/Alt.

## 3. Build keywriter from source

**keywriter** (*remarkable-keywriter*) is the editor engine: a **Qt 5** app in **C++** and **QML**. **Writerdeck** is our on-device binary — that engine plus our patches.

The four-year-old prebuilt binary dies at the loader (`libQt5Quick.so.5`). Qt is static-linked into xochitl on current firmware, so `LD_LIBRARY_PATH` cannot rescue it. We cross-build from source in `ghcr.io/toltec-dev/qt:v3.3` (CI) and deploy a Qt5 runtime sysroot. Writerdeck renders via linuxfb on the rM1's real `/dev/fb0` — rm2fb is not needed.

Using `build-keywriter.sh` as the patching layer is a constraint choice, not an ideal end state. It keeps builds reproducible from a clean upstream checkout, but it is brittle at this size: patch order and context coupling make regressions easier, and reviewability drops as more editor logic lives in generated string patches.

Prefer moving substantial QML/C++ edits into a maintained Writerdeck fork of keywriter and reducing `build-keywriter.sh` to build glue plus minimal deterministic patches. The long patch script is an emergency maintenance model, not a destination — do not wait for harness **105/105**, and do not spend the migration queue on leftover non-critical harness fails first.

**Fork (owned):** [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) — fork of `dps/remarkable-keywriter`, default branch `master`. CI clones that URL via `KEYWRITER_REPO` / `KEYWRITER_REF` (defaults in `build-keywriter.sh`, Dockerfile `ENV`, and `build-keywriter.yml`). Phase 1 pin verified: edit-session PASS, critical **36/36**. Patch script still applies unchanged. Handoff: [todo-handoff-keywriter-fork.md](todo-handoff-keywriter-fork.md).

Phasing ([TODO.md](../TODO.md) item 3):

1. Pin CI to the fork with **no behavior change** — **done.**
2. Move behavior from the patch script into forked C++/QML **by criticality**, in bulk groups that belong together (for example: caret + shift selection + backspace/delete; then wrap/visual line; then undo; then combos/gap polish). Critical editing paths first; remaining harness fails only when their feature group is the one being migrated.
3. Shrink the script to build glue; document fork ownership and upstream-merge policy here.

Critical **36/36** means basic editing is gated green. Full **105/105** remains product sign-off. Neither blocks starting the fork; neither should reorder Phase 2 away from criticality-first migration.

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

## 9. Share is Download

The native iOS share sheet needs HTTPS, which we do not have on plain LAN http. **Download** with `Content-Disposition: attachment` is the reliable export path on the phone (per-row button on the note list). A read-view **Copy** button existed earlier; it left with phone preview dedup. Paste-at-cursor in Type mode is insert-while-editing, not export — see [browser-vs-tablet.md](browser-vs-tablet.md).

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

The engine runs on Writerdeck-server. The token lives in browser `localStorage` and tablet RAM via `POST /api/sync/token`; it is never written to disk. After a restart clears tablet RAM, the server sends WebSocket `needtoken` to connected browsers so a saved token can be reposted automatically; the browser may also push on reconnect via `refreshSyncStatus()`. Reconcile unions tablet and repo note lists and copies anything missing from either side — it never deletes on its own. Tablet Files CRUD pairs to GitHub via the trusted socket. External deletes on GitHub propagate when the local copy is pristine and carries a stored `sha`, confirmed with a per-note 404 before acting; unpushed local edits resurrect instead of deleting. Empty-push guard refuses to push a zero-byte file over a previously-synced note. Details: [server-sync-implementation.md](server-sync-implementation.md).

## 20. Display rotation persists in settings

Global rotation (0, 90, 180, 270) is stored in `.Writerdeck/settings.json`. On editor connect the server pushes `setrotation` from saved settings. On the tablet, Ctrl-R, Ctrl-arrows, and the Lobby Settings button call `rotateScreen()`; USB rotation relays `rotationChanged` back via `rotation_watcher`. Phone Preferences no longer expose rotate — tablet only.

## 21. Edit-session regression test

`POST /api/open` launches a note from outside the tablet UI (tests and scripts). If Writerdeck exits immediately, the session ends, xochitl restarts, and it looks like the stock UI reloading — almost always a broken QML patch, not the server. `scripts/test-edit-session.sh` POSTs `/api/open` from stock UI and asserts Writerdeck stays up about eight seconds. Run it after Writerdeck or QML deploy (`rmkw`); not after server-only deploy — restart the server and spot-check the API instead. `build-keywriter.sh` asserts brace balance in `handleKey()` before write.

## 22. Keyboard selection harness

Modifier+arrow and selection behaviour are tested on the device via `daemon/cmd/edit-harness` and `scripts/test-keyboard-harness.sh`, not by reading saved note bytes. Writerdeck publishes cursor/selection/`contentY` over the socket; the server exposes `GET /api/test/editor-state` and `POST /api/test/reset` (hard quit). Scenarios send keys over `/ws` like the phone UI; hardware page turns use editor cmds `pageleft`/`pageright` (not Arrow keys). **105 scenarios**; sign-off **105/105 PASS** with `--fast`. **Critical subset: 36 scenarios** (`-t critical --fast`); must pass before claiming basic editing works. Motion/selection pattern: uni 1, uni 5, bi 1+1, bi 3+5 (overshoot), bi 7+7 — both directions. Scoreboard: [editor-testing/milestone-runs.md](editor-testing/milestone-runs.md). Handoff: [editor-testing/todo.md](editor-testing/todo.md).

Default run uses **sandbox-prepare**: one editor launch per full suite; between scenarios the harness `PUT`s note content and sends `harnessopen` + `harnessprepare` (QML `harnessSandboxReset` — reload text, cursor 0, clear undo, optional width) without quitting Writerdeck. When the run finishes (pass or fail), the harness sends `showlobby` so the tablet returns to the Lobby. Edit-mode keys are dispatched via QML `socketRouteKey()` from the socket inject thread (not raw `QKeyEvent` to focus). Tap placement and vertical Up/Down share a **visual goal-x** (`positionToRectangle`); harness simulates tap via `harnesssetcursor`. `POST /api/test/reset` remains for manual hard quit only; `--hard-reset` was removed from the harness script. Unit coverage for `translate()` modifier masks stays in `daemon/editor_test.go`. The harness does not exercise USB evdev or `.qmap` — re-check Norwegian Alt+arrow, AltGr, and national characters on hardware after qmap changes. Run after QML arrow/selection changes and after daemon test-handler changes. Harness notes use the `z-test-` prefix (decision 32).

## 23. On-device Writerdeck naming

On the tablet: `Writerdeck`, `Writerdeck-server`, `Writerdeck-user-documents/`, `.Writerdeck/`, `/run/Writerdeck.sock`, `writerdeck.service`. Repo script names like `deploy-rmkbd.sh` stay historical. `migrate-device-layout.sh` renames legacy paths on deploy.

## 24. Tablet file CRUD via trusted socket

The Lobby Files tab sends `{"t":"req","op":…}` over the existing Unix socket; the server performs the same disk ops as `/api/notes`. Tablet delete, rename, and create queue sync and notify the phone via `tabletcrud`. Launch the Lobby from the Mac with `wd` or `bash scripts/lobby.sh`, or on the tablet with `~/wd`.

## 25. Display sync after Lobby clears the buffer

Returning to the Lobby assigns `query.text = ""`, which breaks any `text:` binding on the `TextEdit`. `doLoad` and `toggleMode()` call `syncQueryDisplay()` so edit shows plain `doc` and preview shows `readHtml(doc)` — never RichText extracted back into `doc`. Without this, Home save can zero files. `showLobby` also clears `currentFile`.

## 26. Plain-text edit mode

Edit mode stays `TextEdit.PlainText` — raw Markdown bytes in `doc`, monospace on screen. Read mode (Esc) renders via sundown into RichText. Formatted headings, bold, and italic belong in read mode only.

RichText in edit mode was tried upstream and in early patches; reading formatted `query.text` back into `doc` caused empty saves, HTML on disk, and corrupted previews. The integrity contract gates any WYSIWYG-in-edit proposal: a display-only overlay or C++ syntax highlighter could work in theory, but true hide-the-markers editing is out of scope unless someone writes an ADR and re-tests slices 1–11.

## 27. Lobby Files: Edit and Read

The Files tab offers Edit and Read instead of a single Open. Edit runs `saveAndLoad()` — type mode; the server broadcasts `openedit` and the phone enters Type mode. Read runs `doLoad()` in preview (`mode=0`) and broadcasts `openread` so a Bluetooth keyboard on the phone can still send Esc to toggle edit. USB: Enter edits, `v` reads. Touch: double-tap or a second tap on the already-selected row opens Edit.

## 33. Bluetooth remote key capture on the phone

Bluetooth keyboards pair to the phone; keystrokes reach the tablet over the same WebSocket inject path as Type mode. The phone only captures when the tablet asks: `openedit` (full Type mode), `openread` (read preview — Esc to edit), or `lobbyinput` (Files new, rename, new-encrypted, or delete confirm). The tablet signals via `notifyReadOpen` / `notifyLobbyInput` on the editor socket; the server fans out to browsers. Do not capture in plain Browse mode — browser shortcuts must still work.

## 28. Physical Home duplicate delivery (interim)

gpio-keys on `/dev/input/event1` reaches both Writerdeck-server (`cmd home`) and Qt evdev (`Key_Home`). Without pairing, one press from edit/read could land in the Lobby then immediately quit. Interim fix: `handleHome(fromPhysicalCmd)` sets `suppressNextHomeKey` on the cmd path; the duplicate `Key_Home` consumes it. Durable fix: exclusive grab on `event1` in Go so Qt never sees physical Home — [handoff-physical-home-input.md](handoff-physical-home-input.md).

## 29. Lobby keyboard regression test

After Home from edit, `lobbyFocus` must keep USB and WebSocket keys working — a bare `FocusScope` without `Keys` handlers stole focus. `scripts/test-lobby-keyboard.sh` opens a note, drops to Lobby via `POST /api/lobby`, sends Enter over WebSocket (Files is the default tab), reopens the note, and asserts Home-from-read does not quit Writerdeck (`POST /api/test/home`). Run after Lobby or `handleHome` QML changes alongside `test-edit-session.sh`.

## 33. Lobby tab order (Files first)

Files is tab 1 so boot and Home from edit land on the note list, not the welcome Home screen (tab 6). Returning from edit saves the open filename in `lobbyLastEditedFile` and selects that row after the notes list refreshes. Ctrl-K from the Lobby still opens the note picker on Files. Vault e2e harnesses use digit `1` for Files and `4` for Settings.

## 31. Optional at-rest encryption (private notes)

Two independent PINs: the LAN pairing PIN (`pinDigits`, phone browser) and the vault encryption PIN (tablet only). `pinDigits: none` does not disable encryption.

Encrypted notes use suffix `.md.enc` beside plain `.md` in `Writerdeck-user-documents/`. Crypto: AES-256-GCM per file with a random 32-byte data key; user PIN derives a KEK via scrypt (N=32768); settings store `encryptionEnabled`, `vaultSalt`, `vaultVerifier`, `wrappedDataKey` only. On-disk format: magic `WDENC1` + nonce + ciphertext+tag. Stdlib AES-GCM plus `golang.org/x/crypto/scrypt`; `CGO_ENABLED=0`.

Unlocked state is gone: the data key lives in server RAM only during an active note-editing session (one encrypted note) or briefly for a one-shot Files encrypt/decrypt or phone download. PIN every time you open, read, edit, encrypt, or decrypt — including note switches. Edit/read toggle on the same note does not re-prompt; save on Lobby exit does not re-prompt while that session key is held. Returning to the Lobby clears the session key. PIN entry is tablet-only: touch numpad or USB digits + Enter. Failed attempts rate-limit like pairing PIN auth. Files Encrypt/Decrypt and New encrypted show the PIN overlay first; after a correct PIN the deferred op runs once via `vaultpinok`.

Per-note encrypt/decrypt from Lobby Files — no bulk encrypt on enable. With private notes on, Files shows a second touch row (Encrypt and New encrypted on plain notes, Decrypt on `.md.enc`); Settings is enable and change PIN only. Failed encrypt/decrypt pushes `vaultopfailed` to the tablet — a red message on the Files tab (corrupt ciphertext, bad format, name clash). Opening an encrypted note with the wrong vault key shows the same error path, not a blank editor. Disabling vault or applying a new `secret/vault` from sync is refused while non-test encrypted notes exist — `vaultChangePIN` re-wraps the same data key and is safe; `disablevault` plus fresh setup mints a new key and orphans old ciphertext unless notes are re-wrapped. GitHub sync treats `.md.enc` as opaque bytes and mirrors recovery material under `secret/pin` (plaintext PIN) and `secret/vault` (JSON wrap metadata). `secret/` is excluded from phone note APIs.

Phone download decrypts server-side after the tablet enters its PIN; if no session key is present, the server returns 423, pushes `requestvaultpin` to the tablet, and the phone waits. No encryption PIN UI on the phone. Forgotten PIN: recover from `secret/pin` on GitHub after re-deploy.

UTF-8 Markdown integrity applies to `.md` only; `.md.enc` are opaque on disk. See [integrity-audit.md](integrity-audit.md).

## 32. Device test harness note names

Notes created or opened by device regression scripts (`test-edit-session.sh`, `test-keyboard-harness.sh`, `test-vault.sh`, `test-vault-e2e.sh`, and future harnesses) use the `z-test-` filename prefix so they sort last in the Lobby Files list and are easy to tell from user notes. Example: `z-test-keyboard-harness.md`, `z-test-vault-e2e.md`.

---

## Open risks

Firmware OTA may wipe the systemd unit and regenerate the SSH password — recovery is re-deploy and re-enable. USB `us`/`no` qmaps and Lobby layout picker are shipped and device-verified. Encrypted notes subset is implemented (decisions.md §31). Integrity residuals: [integrity-audit.md](integrity-audit.md). uinput is closed — see decision 1. Go must be on the Mac. Rootfs is about 96% full; everything we ship lives on `/home/root/`. Do not resize rootfs — A/B OTA scheme, brick risk.

The editor patch stack is an active architectural risk: too much behavior lives as generated string patches in one script. Migrate to a Writerdeck fork of keywriter (decision 3, [TODO.md](../TODO.md) item 3) by moving **critical feature groups first**, not by clearing the leftover harness fail list first. Harness **105/105** remains product sign-off; it does not set migration order.