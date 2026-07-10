# Lessons learned

Operational gotchas from building Writerdeck — the stuff that burned time once and shouldn't burn it again. Architectural *why* lives in [decisions.md](decisions.md); shipped features in [DONE.md](../DONE.md).

---

## Deploy & staleness

**Three layers of "my change did nothing."** (1) CI keywriter binary lags the git push. (2) Browser caches the capture page — serve with `Cache-Control: no-store`. (3) A live editor session keeps the old keywriter binary; respawn after deploy (Home→reopen or reboot).

**`rmkw` is binary-only.** Fonts live in the Qt sysroot (~14 MB). After a font change: `RM_FORCE_SYSROOT=1 bash scripts/deploy-keywriter.sh -b`, then respawn the editor.

**`deploy-rmkbd.sh` kills rmkbd.** Follow with `systemctl restart rm1-writerdeck`.

**scp deadlocks** at a fixed offset on the Mac→Wi-Fi→tablet link. Use `rm_send_file` (gzip-over-ssh) in `_env.sh`.

**ETXTBSY on deploy** — kill by full path before copying; stream to `.new`, then `mv`.

**Browser rotate needs keywriter deploy** — `POST /api/rotate` is rmkbd-only, but the tablet must handle the socket `rotate` cmd. `deploy-rmkbd.sh` alone leaves an old keywriter that ignores it.

## systemd & device

**`RequiresMountsFor=/home/root`** on any unit whose `ExecStart` lives on `/home` — otherwise cold boot races the mount and you get `203/EXEC`.

**`HOME=/home/root` in launch-keywriter.sh** — under systemd, root's `$HOME` is `/`, so keywriter's save path breaks without the export.

**Keep the tablet awake** — it drops Wi-Fi on suspend.

## keywriter / QML

**Every save path must sync `query.text → doc` before `saveFile()`** in edit mode. A bare `saveFile()` writes stale `doc`. Guards: saveAndQuit, handleHome, showLobby, saveAndLoad, omni switcher.

**Lobby is a clean no-file state** — clear `currentFile` on every return; guard `saveFile()` when empty. A stale `currentFile` resurrects deleted notes.

**Ctrl+K / modifier flags** — keywriter's `ctrlPressed` bool only flips on a standalone Control key. Injected keys use the modifier *flag*; `handleKeyDown` must also read `event.modifiers & Qt.ControlModifier`.

**Qt 5.15 RichText ignores `margin-bottom` on `<p>`/`<li>`.** Use `line-height` or spacer nodes; always verify on device.

**Font IDs must match Qt family names exactly** or the editor silently falls back to DejaVu.

**QML `Text` needs explicit `width` + `wrapMode`** or long Lobby copy runs off-screen.

**Apostrophes in Python patch heredocs** — use `' + chr(39) + '`, not a literal `'`.

**No cursor blink on e-ink** — it ghosts and smears. Hide while typing won.

## Browser / capture page

**Capture must stand down when an overlay is up** — PIN screen or paste modal. Otherwise keystrokes leak to the tablet.

**`display: ''` restores the stylesheet value** — if CSS says `display:none`, setting `''` keeps it hidden. Set an explicit value.

**Inline `onclick` can't reach IIFE closures** — use `addEventListener`.

**`navigator.clipboard` needs HTTPS** — on plain http LAN, Copy falls back to `execCommand('copy')`.

**Load sync flags at page init**, not when the Settings panel opens — otherwise auto-sync silently skips.

**Async primitives must return their promise** — `reconcileAll` didn't wait on `pushNote`; concurrent GitHub PUTs lost commits.

**GitHub token is per-origin** — new DHCP IP = new browser origin = re-enter token.

## Sync

**Destructive sync ops need per-note confirmation** — `reconcileAll` maps a failed remote list to `[]`; without a 404 guard, one network blip would mass-delete the tablet.

**Boot in edit mode, don't inject Escape** — daemon, editor, and client have independent lifetimes; a synthetic toggle desyncs on reconnect.

## CI / patches

**One patch file = one target file.** Multi-file `git apply --recount` can't tell where hunks end; second-file edits go through `build-keywriter.sh` sed/python.

**Font CI: one hard-failing `RUN` per font** with `fc-list | grep` assertion. A trailing `|| true` swallows download failures.

**`int(Uint32) % N` overflows 32-bit `int` on device** — modulo in `uint32` space first.

## Recon on BusyBox

This `od` is a stub — pull raw bytes to the Mac and decode with BSD `od`. No `timeout` — use `dd & sleep & kill`.
