# Lessons learned

Operational gotchas from building Writerdeck — the stuff that burned time once and shouldn't burn it again. Architectural *why* lives in [decisions.md](decisions.md); shipped features in [DONE.md](../DONE.md).

---

## Deploy & staleness

**Three layers of "my change did nothing."** (1) CI keywriter binary lags the git push. (2) Browser caches the capture page — serve with `Cache-Control: no-store`. (3) A live editor session keeps the old keywriter binary; respawn after deploy (Home→reopen or reboot).

**`rmkw` is binary-only.** Fonts live in the Qt sysroot (~14 MB). After a font change: `RM_FORCE_SYSROOT=1 bash scripts/deploy-keywriter.sh -b`, then respawn the editor.

**`deploy-rmkbd.sh` flushes then waits.** `POST /api/flush-save` (loopback) runs `autosavenow` on the open note; then SIGTERM + wait up to ~12 s for Writerdeck-server to exit (save + quit). Follow with `systemctl start writerdeck`. Old behaviour (`pkill` + 0.5 s) could drop unsaved buffer — documented breach, fixed slice 11.

**scp deadlocks** at a fixed offset on the Mac→Wi-Fi→tablet link. Use `rm_send_file` (gzip-over-ssh) in `_env.sh`.

**ETXTBSY on deploy** — kill by full path before copying; stream to `.new`, then `mv`.

**Browser rotate needs both binaries** — `POST /api/rotate` is handled by Writerdeck-server (saves `"rotation"` in `settings.json`), but restore and USB persistence need Writerdeck to handle `setrotation` and relay `rotationChanged` via `rotation_watcher`. `deploy-rmkbd.sh` alone leaves an old editor at 0° on relaunch; phone rotate may save without the tablet moving.

**Rotation watcher is a separate moc'd unit** — `Q_OBJECT` in patched `main.cpp` won't link on the ARM Qt build. Keep `rotation_watcher.{h,cpp}` in `edit.pro` (copied by `build-keywriter.sh`).

## systemd & device

**Supervisor logs live in journald** — with `writerdeck.service`, `Writerdeck-server` and Writerdeck Qt stderr both land in `journalctl -u writerdeck.service` (session start/stop, Home relay, QML errors, `editor process exited`). Ad-hoc `nohup … >/tmp/wd-server.log` only when running outside systemd.

**`RequiresMountsFor=/home/root`** on any unit whose `ExecStart` lives on `/home` — otherwise cold boot races the mount and you get `203/EXEC`.

**`HOME=/home/root` in Writerdeck-launcher.sh** — under systemd, root's `$HOME` is `/`, so Writerdeck's save path breaks without the export.

**No `pkill -f /home/root/Writerdeck`** — matches `Writerdeck-server` too. Kill the editor with `pidof Writerdeck`; kill the server with `pkill -f /home/root/Writerdeck-server`. Deploy scripts already do this; ad-hoc SSH restarts must too, or you stack duplicate processes.

**Keep the tablet awake** — it drops Wi-Fi on suspend.

## Writerdeck / QML

**Every save path must sync `query.text → doc` before `saveFile()`** in edit mode. A bare `saveFile()` writes stale `doc`. Guards: saveAndQuit, handleHome, showLobby, saveAndLoad, omni switcher, Ctrl-Q.

**Never clear `query.text` without re-syncing on load** — assigning `query.text = ""` (e.g. returning to Lobby) breaks the `text: doc` binding on the `TextEdit`. `doLoad` must set `query.text = response` after every file read (build-keywriter.sh edit 2b), or the next Home save runs `doc = query.text` → empty → `saveFile()` wipes the file on disk. Symptom in `journalctl -u writerdeck`: repeated `qml: Save foo.md` then `qml: save -> 0`. First open after boot can look fine; **second** open/Home cycle wipes. Device-verified fix 2026-07-11.

**Lobby Files open regression (Jul 2026)** — shipped Lobby subpages let you open notes on tablet without the phone. Testing pattern open → Home → open another triggered the binding bug above; GitHub sync then pushed empty files and clash handling created `(tablet copy).md` junk. Recovery: `bash scripts/restore-wiped-notes.sh` (tablet + `my-notes` from pre-wipe commit). Prevention: edit 2b + sync empty-push guard in `sync.js`.

**Python comments inside QML patch heredocs** — `build-keywriter.sh` embeds QML via Python string literals. A `//` JavaScript comment *outside* the string (e.g. after `'query.text = ""\n'`) is a Python `SyntaxError` and fails CI silently until you read the build log. Use `#` for Python-side comments only; QML comments must live inside the quoted string if needed.

**Socket-triggered saves ack back to Writerdeck-server** — `{"t":"saved","c":"home|open|..."}` after the QML handler finishes (BlockingQueuedConnection). Writerdeck-server waits for that before `exitedit`, GitHub push, or HTTP 200 on `/api/open`. Power sleep also gets `{"t":"ready","c":"preparesleep"}` after the e-ink sleep screen paints (~800 ms). Never guess with fixed sleeps for save timing.

**Lobby is a clean no-file state** — clear `currentFile` on every return (`handleHome`, `showLobby`); guard `saveFile()` when empty. A stale `currentFile` resurrects deleted notes.

**Ctrl+K / modifier flags** — Writerdeck's `ctrlPressed` bool only flips on a standalone Control key. Injected keys use the modifier *flag*; `handleKeyDown` must also read `event.modifiers & Qt.ControlModifier`.

**USB Escape launch** — Writerdeck-server watches USB keyboard evdev nodes (hotplug rescan every 3 s). Escape while no active session and not sleeping → `start()` (Lobby). Ignored while editing (Writerdeck owns Esc) and while sleeping (power button wakes). Not an Esc-to-wake path.

**Page-button chord launch** — same idle path when **left+right** physical page buttons are held together on `/dev/input/event1` (800 ms debounce). Tablet-only; no USB or phone needed.

**Qt 5.15 RichText ignores `margin-bottom` on `<p>`/`<li>`.** Use `line-height` or spacer nodes; always verify on device.

**Font IDs must match Qt family names exactly** or the editor silently falls back to DejaVu.

**QML `Text` needs explicit `width` + `wrapMode`** or long Lobby copy runs off-screen. The Lobby uses a `Flickable` with vertical centering when content fits (`y = max(0, (viewport − height) / 2)`); taller content scrolls from the top.

**Apostrophes in Python patch heredocs** — use `' + chr(39) + '`, not a literal `'`.

**QML patch blocks must balance braces** — patch 7p (Lobby Ctrl+arrow rotation) once opened `else if (mode == 0 || isLobby) {` without closing it, leaving `Component.onCompleted` inside `handleKey()`. Symptom: `QQmlApplicationEngine failed to load component` / `Expected token ','` at the next top-level item; Writerdeck exits immediately on `/api/open` → `session.end()` → stock UI. `build-keywriter.sh` now asserts `{`/`}` balance in `handleKey` before write; still verify with `scripts/test-edit-session.sh` on device.

**Edit → stock UI in one beat** — if logs show `editor started` then immediate `editor process exited` with a QML parse error, it's a broken `main.qml` patch, not the USB Escape watcher. Rebuild Writerdeck (`rmkw` after CI) and redeploy.

**Home can look like a crash** — two-level Home is intentional: first press saves and returns to Lobby; second press (from Lobby) quits to stock UI. Under systemd, both show as `home button -- relaying to editor` then `editor process exited` / `starting xochitl` — not a segfault. One press straight to stock UI from an open note would be a bug; check whether Lobby appeared between presses (`journalctl | grep home`).

**No cursor blink on e-ink** — it ghosts and smears. Hide while typing won.

## Browser / capture page

**Capture must stand down when an overlay is up** — PIN screen or paste modal. Otherwise keystrokes leak to the tablet.

**`display: ''` restores the stylesheet value** — if CSS says `display:none`, setting `''` keeps it hidden. Set an explicit value.

**Inline `onclick` can't reach IIFE closures** — use `addEventListener`.

**`navigator.clipboard` needs HTTPS** — on plain http LAN, Copy falls back to `execCommand('copy')`.

**Lobby last-sync on tablet** — relative time comes from `lastSyncAt` in `settings.json`. The phone must POST `/api/sync/ack` after reconcile; `reconcileAll` does this on every success (not only on power sleep). `pushLobbyInfo` formats it for the Lobby and re-pushes on ack.

**Load sync flags at page init**, not when the Preferences panel opens — otherwise auto-sync silently skips.

**Async primitives must return their promise** — `reconcileAll` didn't wait on `pushNote`; concurrent GitHub PUTs lost commits.

**GitHub token is per-origin** — new DHCP IP = new browser origin = re-enter token.

## Sync

**Destructive sync ops need per-note confirmation** — `reconcileAll` maps a failed remote list to `[]`; without a 404 guard, one network blip would mass-delete the tablet.

**Never push an empty tablet file over a previously-synced note** — `pushNote` refuses when `content === ""` and `ghLocalHash` was non-empty; reconcile pulls from GitHub instead. Clash handler skips `(tablet copy)` when tablet is empty and GitHub has content. Belt-and-suspenders after the Lobby Home wipe bug — a genuine "delete all text" edit still saves non-empty until the owner clears and re-syncs intentionally.

**Open-file tracking — slices 1+3+4 shipped, residuals remain** — tablet opens report via `notifyOpen` → server `openNote` / `openedit` WS (slice 1). Phone rename/delete of the open file notifies the editor (`noterenamed` / `notedeleted`, slice 3). Reconcile skips `openNote` from `/api/status` (slice 4), not phone `typingMode` alone. **Still open:** stale `tabletOpenNote` after phone-back; `doLoad` async races on rapid switch; clash/pull overwrites disk without auto-reload (drift banner is manual). See [integrity-audit.md](integrity-audit.md).

**Boot in edit mode, don't inject Escape** — daemon, editor, and client have independent lifetimes; a synthetic toggle desyncs on reconnect.

## CI / patches

**One patch file = one target file.** Multi-file `git apply --recount` can't tell where hunks end; second-file edits go through `build-keywriter.sh` sed/python.

**Font CI: one hard-failing `RUN` per font** with `fc-list | grep` assertion. A trailing `|| true` swallows download failures.

**`int(Uint32) % N` overflows 32-bit `int` on device** — modulo in `uint32` space first.

## Recon on BusyBox

This `od` is a stub — pull raw bytes to the Mac and decode with BSD `od`. No `timeout` — use `dd & sleep & kill`.
