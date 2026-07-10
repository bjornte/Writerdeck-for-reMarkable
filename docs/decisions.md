# Writerdeck for reMarkable 1 — Architecture Decision Record (ADR)

The why behind the project's choices — the durable record so a fresh reader (or LLM) doesn't re-litigate settled ground or retry a dead path. How it works: [architecture.md](architecture.md). Open plan: [../TODO.md](../TODO.md). Shipped: [../DONE.md](../DONE.md). Operational gotchas: [lessons.md](lessons.md).

Each entry: a decision, its status, and the reasoning. Newest/most-foundational first.

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

## 6. Build `rmkbd` on the Mac host, keywriter in CI
Status: active. The Mac is the only machine that can reach the tablet, so deploys originate there regardless; a static Go binary cross-compiles in ~1 s, so a host `go build` is the fastest edit→deploy loop. CI/Docker is reserved for keywriter (it needs the toltec Qt sysroot).

## 7. Always-on `rmkbd` + on-demand editor (the lifecycle split)
Status: built (Phase 8 slice 8c), device-verified. `rmkbd` is its own always-on service that keeps serving `:8000` even under the stock GUI, and owns the xochitl↔keywriter toggle in Go — xochitl is stopped/started per *editor session*, not per *rmkbd lifetime*, so pressing Home returns to the GUI while `:8000` keeps serving and the phone re-launches writing with no reboot. Boot auto-launches one session (power-on = typewriter). Keep-awake is scoped to a session: `systemd-inhibit` wraps *keywriter* (the on-demand child), so the sleep-block lasts only during an editor session and the tablet sleeps normally under the GUI — zero Go D-Bus code.

> Gotcha banked: the unit's `ExecStart` is now `/home/root/rmkbd` directly, which lives on a late-mounting partition → a cold boot raced the mount and failed `203/EXEC`. Any systemd unit whose `ExecStart`, working dir, or `EnvironmentFile` lives on `/home/root` *must* declare `RequiresMountsFor=/home/root`. (Phase 7 accidentally dodged this when `ExecStart` was `/usr/bin/systemd-inhibit`, on rootfs.)

## 8. Phase 8 Companion — phone owns files, tablet owns editing — with PIN-on-tablet auth
Status: built (slices 8a–8e), device-verified. The owner wanted a boot Lobby (IP + connect how-to + Home-to-exit + GitHub line) and full note management from the phone. De-risked by inspection: the tablet *is* the web server (`rmkbd`), so no phone app; `rmkbd` (Go) does all file ops on `/home/root/edit/*.md` natively; the only keywriter changes are a Lobby overlay (mirroring the existing `isOmni` Rectangle) and "open note X" = the existing `doLoad(name)`. PIN shown on the tablet per boot (no stored secret; you must hold the device to read it). Two-level Home (edit → Lobby, Lobby → quit) is the nice consequence: it removes the "reboot to write again" limit.

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

Extended (Phase 9 P, device-verified) — the length is now owner-choosable: 6 / 4 / none. A phone settings screen lets the owner pick the PIN length or turn it off entirely. "none" is a real downgrade — it opens the notes API (and the WS upgrade) to anyone on the Wi-Fi — so the settings row carries an explicit warning and it is strictly owner-chosen; 4-digit drops back to 10,000 combos, but the per-IP 5-guesses/60 s lockout still stretches full exhaustion to hours, so it stays a reasonable convenience. The PIN is now runtime-mutable (guarded by `authMu`, alongside `authToken` and a derived `pinRequired`): a length change re-mints the PIN under lock, re-pushes the Lobby `info` so the tablet shows it at once, clears stale `pinAttempts`, and hands the changer a fresh `rmkbd_token` cookie so turning a PIN *on* doesn't instantly 401 them. Stored as a string enum (`"6"/"4"/"none"`, default `"6"`), not an int — an absent field on an older font-only `settings.json` would read `0` and silently mean "none"; the loader defaults `"" → "6"`.

## 18. `/api/lobby` (Lobby-on-demand) is pre-auth but leaks nothing over the wire
Status: built (Phase 9 L), device-verified. A second device that arrives *after* the owner is already editing finds the tablet showing the note, not the PIN — so it has nothing to authenticate with. `POST /api/lobby` (the "Show PIN on tablet" button on the PIN screen) asks the tablet to save the open note and drop to the Lobby, bringing the PIN back onto the e-ink. It is necessarily pre-auth (you can't gate "show me the PIN" behind already holding the PIN), and that's acceptable because it reveals the PIN *only* on the physical e-ink, never over the network — the HTTP response carries nothing secret — so it doesn't weaken the must-hold-the-device model (#8); it just asks the tablet to show its own screen. Rate-limited to one honored call per ~3 s (`429 + Retry-After`) so a LAN actor can't grief by spam-flipping the screen — and the open note is saved *before* each flip, so the worst case is annoyance, not data loss. The QML `showLobby()` is idempotent and never quits, deliberately distinct from `handleHome()` (which quits *from* the Lobby — reusing `home` would have been a trap). If no editor session is live, the call starts one (booting straight into the Lobby), so it works from the stock GUI too.

## 19. GitHub note-sync is a non-authoritative reconciler — delete/rename only from the browser
Status: built (Phase 9), device-verified. The optional two-way GitHub sync (off by default; the PAT lives only in the phone browser's `localStorage`, never on the tablet, which holds just non-secret `syncOn`/`syncRepo`) deliberately reconciles by *unioning* the tablet's and the repo's note lists and copying any note missing from one side to the other — it never deletes on its own. That is the intended safety property: a reconciler that cannot delete cannot lose a note, which sidesteps the documented real-git-on-mobile instability (isomorphic-git crashes / packfile corruption on low-RAM devices — the reason we sync via GitHub's plain Contents API, not real git). Accepted cost: destructive ops must go through the *phone browser*, which pairs them — the UI's Delete also calls `ghDelete` (Contents API DELETE + stored `sha`) and Rename deletes the old path then pushes the new. A delete or rename made *outside* the browser — VS Code, `git`, the GitHub web UI — used to resurrect or duplicate on the next sync; the marker-aware path in [../daemon/sync.js](../daemon/sync.js) (device-verified) now treats "on the tablet + stored `sha` + pristine + gone from GitHub" as a real delete (confirmed per-note 404 before acting), covering the common GitHub-side delete and (as a side effect) external rename as delete-old + pull-new. Two safety invariants: unpushed local edit resurrects (push) rather than deletes; delete only fires after a fresh per-note `GET` returns 404 — guarding against `reconcileAll` mapping a transient remote-list failure to `[]`. Still (correctly) leaves a purely-local unpushed tablet delete alone: GitHub stays authoritative for deletes.

---

## Known gaps & open risks

- Firmware update (OTA) may break the setup *(open · low)*. An OTA can wipe the systemd unit and regenerates the SSH password. Mitigation: we ship only a static binary + user files + one unit (no Toltec), so the OTA itself stays intact; recovery = re-deploy + re-`enable`, re-record the password. This is the one genuinely open operational risk — tracked in [../TODO.md](../TODO.md) open questions.
- `/dev/uinput` is unavailable and unfixable on this kernel (decision 1). Closed, not a to-do — recorded so nobody retries it.
- Go toolchain must be on the Mac (`brew install go`) — the only device-reachable host.
- Disk: only the tiny rootfs is tight, and nothing we ship goes there. `/` (rootfs, ~228 MB) is 96% full; everything we deploy lives on `/home/root/` (separate multi-GB partition). Don't resize rootfs (A/B OTA scheme; brick risk).
