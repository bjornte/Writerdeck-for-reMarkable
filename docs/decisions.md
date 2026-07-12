# Writerdeck for reMarkable 1 — Architecture Decision Record (ADR)

The why behind the project's choices — the durable record so a fresh reader (or LLM) doesn't re-litigate settled ground or retry a dead path. How it works: [architecture.md](architecture.md). Open plan: [../TODO.md](../TODO.md). Shipped: [../DONE.md](../DONE.md). Operational gotchas: [lessons.md](lessons.md).

Each entry: a decision, its status, and the reasoning. Newest/most-foundational first.

---

## Document integrity — the product contract (foundational)

Status: **active — non-negotiable.** Writerdeck is a typewriter: the owner's prose must survive editing, sync, and normal device use as plain Markdown on disk. This is not a Phase 10 nice-to-have; it gates every feature. How it works: [architecture.md](architecture.md) § Document integrity. Audit: [integrity-audit.md](integrity-audit.md).

**The contract (summary):**

1. **Content** — bytes on disk are UTF-8 Markdown; never Qt `qrichtext` / HTML (slice 2 shipped 2026-07-11).
2. **Lifecycle** — no silent overwrite of an open note by reconcile, phone CRUD, or rename/delete; server + phone track tablet editor state (slices 1+3+4 shipped).
3. **Durability** — defined save paths + 45 s autosave (slice 9) + **save-before-stop** on deploy/SIGTERM (slice 11: `/api/flush-save`, `autosavenow` ack, deploy waits for exit). Atomic note writes server + tablet loopback PUT (slice 10). Residual: SIGKILL / crash before autosave or flush hook runs.
4. **Coherence** — disk change under an open session requires reload or conflict UX; buffer must not blindly win over a pull.
5. **Sync subordination** — GitHub reconcile assists backup; it must not delete, empty-push, or fork paths against a live edit (#19, #24 guards are partial).

**Feature gate:** no change to `daemon/`, `sync.js`, `build-keywriter.sh`, or note APIs ships without an integrity pass against the five points above. Incident ADRs (#24, empty-push) are patches *under* this contract, not substitutes for it.

**Shipped under this contract:** #24, slices 1–11, empty-push guard (#19). **Residual:** SIGKILL/crash before autosave or flush; Writerdeck binary must match server for `autosavenow` + loopback save.

---

## 1. Feed the editor over a local socket, not `/dev/uinput`
Status: active (the core architecture). Keystrokes reach the editor over a local socket, not `/dev/uinput` — that driver is absent and un-addable on this kernel: `/dev/uinput` opens `ENODEV`, and the kernel's exports are trimmed (`CONFIG_TRIM_UNUSED_KSYMS` — only ~375 symbols exported vs >10,000 on a module-friendly kernel), so an out-of-tree `uinput.ko` resolves every import to `Unknown symbol (err -2)` and cannot load (only an OTA-risky kernel-image swap would change that). Don't spend time trying to make uinput work on this kernel — it was built green in CI and still wouldn't bind. A patched keywriter reads keystrokes from the socket instead. Fallbacks if keywriter itself ever fails: a self-contained `libremarkable` editor, or HWR injection into the Wacom pen node (draws ink, not clean Markdown).

## 2. Inject synthetic `QKeyEvent`s into keywriter
Status: active. keywriter takes keyboard input via Qt QPA — there is no `/dev/input` fd to swap. So a detached socket-reader thread in `main.cpp` (`socket-inject.patch`, applied with `git apply --recount`, built `-pthread`) posts synthetic `QKeyEvent`s to `focusWindow()` (the `QQuickWindow`, which replays the real key path including QML `Keys` handlers). Wire format: browser→daemon is full JSON; daemon→keywriter rides text as an integer Unicode codepoint (`{"t":"text","cp":N}` → `QString::fromUcs4`, escaping-proof). The keymap lives browser-side (forward `event.key`, layout already resolved), not in the editor.

## 3. Build keywriter from source (the 4-yr-old prebuilt is dead)
Status: resolved. The prebuilt binary dies at the loader (`libQt5Quick.so.5`) — Qt is static-linked into xochitl on 2026 firmware, so `LD_LIBRARY_PATH` can't rescue it. Cross-build from source in `ghcr.io/toltec-dev/qt:v3.3` (CI) + deploy a Qt5 runtime sysroot. Renders via linuxfb alone (the rM1 has a real `/dev/fb0`; linuxfb drives the EPDC, so rm2fb is not needed). The "do not regress" build fixes live in [build-keywriter.sh](../third_party/keywriter/build-keywriter.sh) and [lessons.md](lessons.md).

## 4. No Toltec
Status: by design. Toltec locks the firmware to a supported version range and can soft-brick on unsupported versions — it conflicts with "preserve OTA updates". Only revisit if the owner ever accepts the version lock for a Python fast-path.

## 5. Static Go binary (`CGO_ENABLED=0`, ARMv7)
Status: by design. No on-device runtime deps; survives firmware updates; cross-compiles in ~1 s. The tablet ships no Python, and installing one implies Entware/Toltec + a firmware lock.

## 6. Build Writerdeck-server on the Mac host, Writerdeck in CI
Status: active. The Mac is the only machine that can reach the tablet, so deploys originate there regardless; a static Go binary cross-compiles in ~1 s, so a host `go build` is the fastest edit→deploy loop. CI/Docker is reserved for Writerdeck (upstream keywriter tree; needs the toltec Qt sysroot).

## 7. Always-on Writerdeck-server + on-demand Writerdeck (the lifecycle split)
Status: built (Phase 8 slice 8c), device-verified. Writerdeck-server is its own always-on service that keeps serving `:8000` even under the stock GUI, and owns the xochitl↔Writerdeck toggle in Go — xochitl is stopped/started per *editor session*, not per *server lifetime*, so pressing Home returns to the GUI while `:8000` keeps serving and the phone re-launches writing with no reboot. Boot auto-launches one session (power-on = typewriter). Keep-awake is scoped to a session: `systemd-inhibit` wraps *Writerdeck* (the on-demand child), so the sleep-block lasts only during an editor session and the tablet sleeps normally under the GUI — zero Go D-Bus code.

> Gotcha banked: the unit's `ExecStart` is `/home/root/Writerdeck-server`, which lives on a late-mounting partition → a cold boot raced the mount and failed `203/EXEC`. Any systemd unit whose `ExecStart`, working dir, or `EnvironmentFile` lives on `/home/root` *must* declare `RequiresMountsFor=/home/root`. (Phase 7 accidentally dodged this when `ExecStart` was `/usr/bin/systemd-inhibit`, on rootfs.)

## 7b. USB Escape launches Writerdeck from stock UI (not wake)
Status: built. `Writerdeck-server` watches USB keyboard evdev nodes (hotplug rescan). **Escape while idle** (stock `xochitl` up, no editor session, not sleeping) → `session.start()` → Lobby — the keyboard counterpart to phone **Edit** without pre-selecting a note. **Not** Esc-to-wake after power sleep (that path didn't work on device; power button only). While a session is active, Esc is ignored here so Writerdeck keeps normal edit/preview/omni behaviour.

## 7c. Left+right page buttons launch Writerdeck from stock UI
Status: built. Physical page buttons (`KEY_LEFT` 105 / `KEY_RIGHT` 106 on `/dev/input/event1`) are readable alongside `xochitl` (no `EVIOCGRAB`). **Both pressed while idle** → same `handleIdleLaunch` path as USB Escape → Lobby. Debounced 800 ms; ignored during an active session or power sleep. Gives a tablet-only launch gesture without USB or the phone. `xochitl` still receives the individual button events in parallel — acceptable on the home screen; in a document the chord may briefly page.

## 8. Phase 8 Companion — phone owns files, tablet owns editing — with PIN-on-tablet auth
Status: built (slices 8a–8e), device-verified. The owner wanted a boot Lobby (IP + connect how-to + Home-to-exit + GitHub line) and full note management from the phone. De-risked by inspection: the tablet *is* the web server (Writerdeck-server), so no phone app; Writerdeck-server (Go) does all file ops on `/home/root/Writerdeck-user-documents/*.md` natively; the only Writerdeck changes are a Lobby overlay (mirroring the existing `isOmni` Rectangle) and "open note X" = `saveAndLoad(name)`. PIN shown on the tablet per boot (no stored secret; you must hold the device to read it). Two-level Home (edit → Lobby, Lobby → quit) is the nice consequence: it removes the "reboot to write again" limit.

## 9. Share = best-effort fallback (Download + Copy)
Status: built (Phase 8 slice 8f), device-verified. The native iOS share sheet needs a secure context (HTTPS), unavailable on plain LAN http — so we shipped the reliable fallback: a `Content-Disposition: attachment` download route + copy-to-clipboard, both in the phone's Read view. `navigator.clipboard` *also* needs a secure context, so Copy falls back to a temporary `<textarea>` + `document.execCommand('copy')` on plain http — same constraint, banked. Revisit the native sheet only if a secure context becomes available.

## 10. Two-machine split + git bridge (tunnel deferred)
Status: active. Author-specific dev-environment workaround, not part of the product: a corporate VPN keeps the dev laptop off the LAN, so device work runs on a second machine and git bridges the two. A reverse tunnel was *feasible* but deferred (infra + exposure risk) in favor of the git bridge; the clean fix is an IT split-tunnel exception.

## 11. Wi-Fi is the dev path (not USB)
Status: active. The Mac's USB-ethernet gadget is inactive (no DHCP lease / no macOS RNDIS); Wi-Fi SSH works when `RM_HOST_WIFI` in secrets matches the tablet (DHCP — changes with network).

## 12. Deploy transport = gzip-over-ssh (scp deadlocks)
Status: active. `scp` wedges at a fixed ~255 KB offset on this Mac→Wi-Fi→tablet link (SFTP app-level windowing deadlock — *not* sleep, *not* QoS; `IPQoS`+keepalives didn't clear it). `rm_send_file` ([scripts/_env.sh](../scripts/_env.sh)) streams the file gzip'd through the SSH channel (`gzip -c src | ssh … 'gzip -dc > dst'`) with a post-copy `wc -c` size check. Lesson: when scp stalls at a fixed offset on an embedded link, switch transports — don't tune scp.

## 13. Secrets in a gitignored env file (plaintext)
Status: active. Low threat model (the password is shown on the device screen, home LAN); the real risk is git leakage, which the ignore prevents. `.gitignore` was created *before* git init so credentials never entered history.

## 14. LF + ASCII enforced via `.gitattributes`
Status: active. CRLF or a stray non-ASCII byte breaks shell scripts / the systemd unit on the device (and once broke PowerShell parsing). `.md` prose may use Unicode; code/device files may not.

## 15. Model usage protocol — Opus plans/inspects, Sonnet writes code
Status: active. Match model strength to task: keep high-reasoning for design/review/diagnosis, save cost/latency on mechanical coding; the assistant flags switch points, the human drives the picker. Authoritative copy in [.github/copilot-instructions.md](../.github/copilot-instructions.md).

## 16. File upload reuses the safe create route; the size cap is server-side
Status: built (Phase 9), device-verified. Uploading a `.md` from the phone deliberately POSTs through the existing `/api/notes` create path rather than a new endpoint, so it inherits every guard already proven there: `notesSafe()` blocks `/`+`..` and forces `.md` (no traversal / arbitrary write), and the pre-write `os.Stat` 409s instead of overwriting (no clobber). Read-back is inert by construction — preview uses `textContent` (never `innerHTML`) and the GET serves `text/plain` / `text/markdown`+`attachment` — so a hostile `.md` can't run script. The client extension/size checks are UX only; the authoritative body cap is `http.MaxBytesReader` on the server, because a client check is bypassable by a direct authed request.

## 17. 6-digit PIN + per-IP brute-force lockout
Status: built (Phase 9), device-verified. The original 4-digit PIN (10,000 combos, no throttle) was the accepted weak point of the on-device-PIN / home-LAN model. Hardened to 6 digits (1,000,000 combos) plus a per-IP lockout (5 wrong guesses → that IP blocked 60 s, HTTP 429 + `Retry-After`; resets on success or expiry). Per-IP, *not* global, so an attacker can't DoS the owner out of their own device — the legitimate user connects from a different address. State is in-memory (reboot clears it, which also regenerates the PIN), pruned each attempt so the map stays small. Compare is `subtle.ConstantTimeCompare` (no timing leak; the lockout is the real defense). Gotcha banked: `int(Uint32) % N` overflows the device's 32-bit `int` — reduce in `uint32` space first. Still by-design *not* hard-enforced beyond this (no account lockout escalation, no persistence) — proportionate to a single-user home device whose PIN is physically displayed and rotates every boot.

Extended (Phase 9 P, device-verified) — the length is now owner-choosable: 6 / 4 / none. The Preferences screen lets the owner pick the PIN length or turn it off entirely. "none" is a real downgrade — it opens the notes API (and the WS upgrade) to anyone on the Wi-Fi — so the settings row carries an explicit warning and it is strictly owner-chosen; 4-digit drops back to 10,000 combos, but the per-IP 5-guesses/60 s lockout still stretches full exhaustion to hours, so it stays a reasonable convenience. The PIN is now runtime-mutable (guarded by `authMu`, alongside `authToken` and a derived `pinRequired`): a length change re-mints the PIN under lock, re-pushes the Lobby `info` so the tablet shows it at once, clears stale `pinAttempts`, and hands the changer a fresh `writerdeck_token` cookie so turning a PIN *on* doesn't instantly 401 them. Stored as a string enum (`"6"/"4"/"none"`, default `"6"`), not an int — an absent field on an older font-only `settings.json` would read `0` and silently mean "none"; the loader defaults `"" → "6"`.

## 18. `/api/lobby` (Lobby-on-demand) is pre-auth but leaks nothing over the wire
Status: built (Phase 9 L), device-verified. A second device that arrives *after* the owner is already editing finds the tablet showing the note, not the PIN — so it has nothing to authenticate with. `POST /api/lobby` (the "Show PIN on tablet" button on the PIN screen) asks the tablet to save the open note and drop to the Lobby, bringing the PIN back onto the e-ink. It is necessarily pre-auth (you can't gate "show me the PIN" behind already holding the PIN), and that's acceptable because it reveals the PIN *only* on the physical e-ink, never over the network — the HTTP response carries nothing secret — so it doesn't weaken the must-hold-the-device model (#8); it just asks the tablet to show its own screen. Rate-limited to one honored call per ~3 s (`429 + Retry-After`) so a LAN actor can't grief by spam-flipping the screen — and the open note is saved *before* each flip, so the worst case is annoyance, not data loss. The QML `showLobby()` is idempotent and never quits, deliberately distinct from `handleHome()` (which quits *from* the Lobby — reusing `home` would have been a trap). If no editor session is live, the call starts one (booting straight into the Lobby), so it works from the stock GUI too.

## 19. GitHub note-sync is a non-authoritative reconciler — delete/rename only from the browser
Status: built (Phase 9), device-verified. **Updated 2026-07-12:** reconcile engine runs on **Writerdeck-server** (not the phone browser). Token: browser `localStorage ghToken` + tablet RAM via `POST /api/sync/token` (never on tablet disk). Tablet **Sync now** + automatic triggers; phone **Setup** for token/repo only. Reconciler behaviour below is unchanged.

The optional two-way GitHub sync (off by default) deliberately reconciles by *unioning* the tablet's and the repo's note lists and copying any note missing from one side to the other — it never deletes on its own. That is the intended safety property: a reconciler that cannot delete cannot lose a note, which sidesteps the documented real-git-on-mobile instability (isomorphic-git crashes / packfile corruption on low-RAM devices — the reason we sync via GitHub's plain Contents API, not real git). Accepted cost: destructive ops must go through the *phone browser*, which pairs them — the UI's Delete also calls GitHub DELETE + stored `sha` and Rename deletes the old path then pushes the new. A delete or rename made *outside* the browser — VS Code, `git`, the GitHub web UI — used to resurrect or duplicate on the next sync; the marker-aware path in [../daemon/syncengine.go](../daemon/syncengine.go) (device-verified) now treats "on the tablet + stored `sha` + pristine + gone from GitHub" as a real delete (confirmed per-note 404 before acting), covering the common GitHub-side delete and (as a side effect) external rename as delete-old + pull-new. Two safety invariants: unpushed local edit resurrects (push) rather than deletes; delete only fires after a fresh per-note `GET` returns 404 — guarding against `reconcileAll` mapping a transient remote-list failure to `[]`. Still (correctly) leaves a purely-local unpushed tablet delete alone: GitHub stays authoritative for deletes.

**Empty-push guard (2026-07-11):** after a Lobby binding bug zeroed tablet files, sync pushed empties to GitHub. `pushNote` now refuses to push `content === ""` when `ghLocalHash` was non-empty; reconcile pulls from GitHub instead; clash handler restores from GitHub when tablet is empty and remote is not — no `(tablet copy)` junk. See #24.

## 20. Display rotation persists in `settings.json`
Status: built, device-verified. Global `root.rotation` (0/90/180/270) is owner-chosen and stored in `.Writerdeck/settings.json` as `"rotation"`. Phone **Rotate tablet 90°** (`POST /api/rotate`) increments, saves, and pushes `{"t":"cmd","c":"setrotation","degrees":N}` to Writerdeck. On every editor socket connect, Writerdeck-server restores the saved angle the same way (alongside `setfont`). USB Ctrl+←/→ in preview/read still rotates in QML; `rotation_watcher` (moc'd helper in the keywriter build) relays `rotationChanged` as `{"t":"rotation","degrees":N}` so the server persists without the phone. Both binaries must be current — server-only deploy can save to disk while an old Writerdeck ignores `setrotation` and boots at 0°.

## 21. Edit-from-browser regression test (`test-edit-session.sh`)
Status: built (2026-07-11), device-verified. Phone **Edit** (`POST /api/open`) is the primary companion launch path — if Writerdeck exits immediately, `session.end()` restarts `xochitl` and looks like the stock UI "reloading." That failure mode is almost always a broken `main.qml` patch (QML parse error), **not** Writerdeck-server, `sync.js`, or the USB Escape watcher.

`scripts/test-edit-session.sh` automates the check: from stock UI, POST `/api/open`, assert `Writerdeck` stays up ~8 s, `xochitl` stays down, and `/api/status` reports `editorActive: true`.

**When to run (device):**

| Change | Run `test-edit-session.sh`? | Instead |
|---|---|---|
| `build-keywriter.sh`, `socket-inject.patch`, `lobby_bridge.*` → Writerdeck deploy (`rmkw`) | **Yes** | — |
| `daemon/` only → `deploy-rmkbd.sh` (Go, embedded `app.js`/`sync.js`) | **No** | `systemctl start writerdeck` (deploy stops the server); spot-check the changed API (`curl /api/status`, browser smoke) |
| Docs / scripts with no binary change | **No** | — |

Build-time guard: `build-keywriter.sh` asserts `{`/`}` balance in `handleKey()` before write (patch 7p once regressed here). Logs to `docs/recon/`.

## 22. On-device Writerdeck naming (2026-07)
Status: built, device-verified. Binaries and paths on the tablet use Writerdeck branding (`Writerdeck`, `Writerdeck-server`, `Writerdeck-user-documents/`, `.Writerdeck/`, `/run/Writerdeck.sock`, `writerdeck.service`, `writerdeck_token`). Repo script names (`deploy-rmkbd.sh`, `third_party/keywriter/`) and the GitHub repo folder stay unchanged. `scripts/migrate-device-layout.sh` renames legacy paths and removes old binaries on deploy; see the on-device table in [architecture.md](architecture.md).

## 23. Tablet file CRUD via trusted socket (Lobby Files)
Status: built (Phase 10 partial), device-verified 2026-07-11. Extends #8 without exposing unauthenticated LAN HTTP: Writerdeck sends `{"t":"req","op":"noteslist|createnote|deletenote|renamenote",…}` over the existing Unix socket; Writerdeck-server performs the same disk ops as `/api/notes`. Six-tab Lobby (**Home · Files · Keyboard · Sync · Settings · Shortcuts**) with touch + Tab/arrows/1–6 navigation; Files page lists notes, supports New/Open/Rename/Delete (touch buttons + `n`/Enter/`r`/`d`). Open still uses `saveAndLoad` (phone path unchanged). Tablet delete/rename/create queues `pendingSync` and notifies the phone via `tabletcrud` WebSocket — GitHub pairs immediately when the browser is connected, or on next connect/reconcile (slice 7; was next-reconcile-only in #19). Mac/tablet launch helpers: `wd` / `bash scripts/lobby.sh` / on-device `~/wd`.

## 24. `doLoad` must re-sync `query.text` after Lobby clears the editor
Status: built, device-verified 2026-07-11. Returning to Lobby (`handleHome`, `showLobby`) assigns `query.text = ""`, which breaks the QML `text: doc` one-way binding. Without `query.text = response` in `doLoad`'s XHR callback (edit 2b), the next Home save copies empty `query.text` into `doc` and `saveFile()` zeroes the file (`save -> 0` in journal). First open after boot worked; open → Home → open another wiped notes and cascaded into GitHub via sync. `showLobby` now also clears `currentFile` for a clean no-file Lobby state. Recovery script: `scripts/restore-wiped-notes.sh`. Lesson banked in [lessons.md](lessons.md).

---

## Known gaps & open risks

- Firmware update (OTA) may break the setup *(open · low)*. An OTA can wipe the systemd unit and regenerates the SSH password. Mitigation: we ship only a static binary + user files + one unit (no Toltec), so the OTA itself stays intact; recovery = re-deploy + re-`enable`, re-record the password. This is the one genuinely open operational risk — tracked in [../TODO.md](../TODO.md) open questions.
- USB keyboard locales *(open · medium)*. Browser/Bluetooth path resolves layout in the phone OS (Norwegian works). USB path uses Qt evdev with **US QWERTY** default — Norwegian æøå and other national layouts need per-layout `.qmap` files via `QT_QPA_EVDEV_KEYBOARD_PARAMETERS` ([remarkable-keywriter#1](https://github.com/dps/remarkable-keywriter/issues/1)); `loadkeys` / `setxkbmap` do not apply. Planned: ship qmaps + `settings.json` picker — [improvements.md](improvements.md), [../TODO.md](../TODO.md) Phase 10.
- Per-note / subfolder encryption *(open · design)*. Global PIN gates the API; no subset protection yet. Encrypted subfolder with passphrase-derived keys is the leading option — design in [improvements.md](improvements.md); implementation not started.
- Tablet file management *(partial · shipped)*. Lobby **Files** tab + socket CRUD covers list/create/open/rename/delete on tablet (#23). Upload, download, copy, paste, and GitHub token entry remain browser-only.
- **Document integrity — residual risks** *(low)*. Slices 1–11 shipped; SIGKILL/crash before autosave or flush hook; binary mismatch (old Writerdeck without `autosavenow`/loopback save). Known open + unknown: [integrity-audit.md](integrity-audit.md).
- `/dev/uinput` is unavailable and unfixable on this kernel (decision 1). Closed, not a to-do — recorded so nobody retries it.
- Go toolchain must be on the Mac (`brew install go`) — the only device-reachable host.
- Disk: only the tiny rootfs is tight, and nothing we ship goes there. `/` (rootfs, ~228 MB) is 96% full; everything we deploy lives on `/home/root/` (separate multi-GB partition). Don't resize rootfs (A/B OTA scheme; brick risk).
