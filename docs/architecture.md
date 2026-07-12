# Writerdeck for reMarkable 1 — architecture & reference

How the system works, the device facts, and the dev/deploy loop. Open work: [../TODO.md](../TODO.md). Why: [decisions.md](decisions.md). Shipped: [../DONE.md](../DONE.md). Gotchas: [lessons.md](lessons.md).

---

## Overview

reMarkable 1: e-paper, Linux, micro-USB, no Bluetooth, no native keyboard support. Goal: sit down, type, get distraction-free prose, save Markdown — but the keyboard lives on an iPhone and keystrokes are forwarded over Wi-Fi to the tablet, which runs the editor and shows text on e-ink.

```
 CLIENT (iPhone Safari page + BT/Lightning kbd; Mac during dev)
     |   WebSocket key events  (Wi-Fi / LAN)
     v
 reMarkable 1
   - Writerdeck-server (Go): capture page + WebSocket + file API
   - Writerdeck (patched keywriter): reads /run/Writerdeck.sock, saves .md
```

Keystrokes reach the editor through a local socket rather than `/dev/uinput`: this tablet's kernel can't load uinput, so Writerdeck-server feeds the patched editor instead (the reasoning is in [decisions.md](decisions.md)).

## Document integrity (non-negotiable)

**The product contract.** Writerdeck exists to produce durable Markdown notes. Every feature, patch, and deploy must satisfy this — not as optional hardening afterward. Why it is foundational: [decisions.md](decisions.md) § Document integrity. Audit: [integrity-audit.md](integrity-audit.md).

| Obligation | Meaning |
|---|---|
| **Plain markdown on disk** | `.md` files are UTF-8 Markdown, never Qt RichText / HTML. Saves reject `qrichtext` payloads; loads sanitize corruption. |
| **Explicit durability** | Disk updates through defined save paths, **45 s autosave** while editing, and **save-before-stop** on deploy/SIGTERM (`POST /api/flush-save` + graceful shutdown wait). Residual: text typed after the last autosave if the process is SIGKILL'd or crashes with no stop hook. |
| **Single-writer awareness** | While a note is open for edit, reconcile, phone ops, and rename/delete must not silently overwrite disk or fork paths. Server and phone must know which file the tablet is editing. |
| **Buffer ↔ disk coherence** | If disk changes under an open session (pull, clash, external edit), the editor must reload or surface conflict — not save stale buffer over fresh disk. |
| **Atomic durable writes** | Note writes use write-temp-rename (settings already do); no in-place truncate on power loss. |
| **Sync is subordinate** | GitHub reconcile assists backup — it must not trump a live edit or push presentation/HTML/empty over good content. |

**Feature gate:** before shipping any change that touches notes, saves, opens, sync, CRUD, or editor lifecycle, ask: *can this lose text, write wrong bytes, or overwrite without the user knowing?* If yes, it does not ship until mitigated or explicitly accepted by the owner.

**As built (2026-07-11):** slices 1–11 shipped. Deploy and SIGTERM flush the open buffer before stopping (`/api/flush-save`); 45 s autosave covers crash/SIGKILL gaps. Deploy both server and Writerdeck for loopback save + `autosavenow` on device.

## On-device layout

What you see over SSH under `/home/root/` (deploy scripts migrate legacy names and remove old binaries):

| Path | Role |
|---|---|
| `Writerdeck-server` | Always-on Go daemon — WebSocket, HTTP API, session lifecycle |
| `Writerdeck` | Full-screen Markdown editor (patched [remarkable-keywriter](https://github.com/dps/remarkable-keywriter)) |
| `Writerdeck-launcher.sh` | Qt/e-ink launch env; spawned by Writerdeck-server with `--editor` |
| `Writerdeck-user-documents/` | Your `.md` notes |
| `.Writerdeck/settings.json` | Font, PIN mode, sync flags, rotation, USB keyboard layout |
| `qt5/` | Qt runtime for Writerdeck (~14 MB; internal support files) |
| `/run/Writerdeck.sock` | Unix socket between Writerdeck-server and Writerdeck |
| `writerdeck.service` | systemd unit (`/etc/systemd/system/`) |

Repo-side source still lives in `daemon/` (builds `Writerdeck-server`) and `third_party/keywriter/` (CI builds `Writerdeck`). Script names like `deploy-rmkbd.sh` are historical; they deploy the binaries above.

## Components

- **Writerdeck-server** (we build, `daemon/`): static Go binary on the tablet. Serves the capture page + WebSocket, file API, per-boot PIN, and the xochitl↔Writerdeck lifecycle.
- **Writerdeck** (upstream [remarkable-keywriter](https://github.com/dps/remarkable-keywriter), patched): full-screen editor; saves `.md` to `Writerdeck-user-documents/`. Input via Qt QPA + socket-injected `QKeyEvent`s (`main.cpp` patch).
- **Client** (embedded in Writerdeck-server): HTML/JS capture page in the phone/Mac browser.

## The Companion (the working appliance)

The tablet is the web server — the phone/Mac just open Safari, no app to install.

### Tablet-first controls (minimize phone duplicates)

Over time, move as much functionality as possible from the phone browser to the tablet Lobby. The phone stays the right surface for clipboard paste, file upload/download, GitHub token entry, and PIN/font prefs that are awkward on e-ink — but anything that works on the tablet should live there.

**Rule:** ship on tablet → device-verify → **remove the duplicate from the phone UI**. Do not maintain two pickers for the same setting. Current split: [browser-vs-tablet.md](browser-vs-tablet.md). Example migrated to tablet-only: **USB keyboard layout** (Lobby **Keyboard** tab).

- Always-on daemon, on-demand editor (the lifecycle split). Writerdeck-server keeps serving `:8000` even under the stock xochitl GUI, and summons Writerdeck on demand: stop xochitl → spawn Writerdeck → on Home, restart xochitl and keep serving. Boot auto-launches one editor session (power-on = typewriter). This is what lets the phone re-launch writing and manage files while the GUI is up.
- Notes = `/home/root/Writerdeck-user-documents/*.md`. Writerdeck-server (Go) does all file ops natively — list/read/create/rename/delete over an HTTP API (`/api/notes…`; `notesSafe()` rejects `/` and `..`). These are Writerdeck's notes, separate from xochitl's document store, so the API is safe to hit while the GUI is up.
- PIN-on-tablet auth (length owner-choosable). Writerdeck-server mints a random PIN per boot and shows it in the Lobby; the phone POSTs it (`/api/pin`) for an HttpOnly `writerdeck_token` cookie that gates the notes API and the WebSocket upgrade. The length is owner-choosable on the Preferences screen — 6 / 4 / none (none = LAN-open, UI-warned), and the PIN is runtime-mutable: changing it re-mints on the spot and re-shows it on the e-ink. A per-IP brute-force lockout (5 wrong / 60 s → 429) backs the 4- and 6-digit modes.
- Preferences screen (phone). A **Preferences** overlay (font, PIN, display rotate, exit) and a separate **Setup** overlay (GitHub sync) in the top bar; both dismiss via Done, ×, ESC, or backdrop click. Font, PIN length, and display rotation (`"rotation": 0|90|180|270`) persist to `/home/root/.Writerdeck/settings.json`. On editor socket connect, Writerdeck-server pushes saved font (`setfont`) and rotation (`setrotation`); USB Ctrl+←/→ in preview posts `{"t":"rotation","degrees":N}` back so the server can save without the phone. When sync is on and a repo is set, Setup links to `github.com/{repo}`. A missing tablet token shows a banner in the phone UI; the tablet **Sync** tab shows a black-bordered **TOKEN NEEDED** box (high contrast for e-ink).
- Lobby-on-demand. A second device that arrives *after* the owner is already editing finds the tablet showing the note, not the PIN; a pre-auth "Show PIN on tablet" button (`POST /api/lobby`, rate-limited) saves the open note and drops the tablet to the Lobby so the PIN is readable — it reveals nothing over the wire (PIN only on the e-ink).
- Lobby + two-level Home. Boot shows a six-tab Lobby pager (**Home · Files · Keyboard · Sync · Settings · Shortcuts** — touch or Tab/arrows/1–6), fed by `{"t":"info",…}` with IP, PIN, `syncOn`/`syncRepo`, note count, formatted last sync, and `keyboardLayout`; re-pushed when `wlan0` gets an address, after reconcile, or when notes/settings change. **Files** tab: list/create/open/rename/delete via trusted socket `{"t":"req","op":…}` (#23). **Keyboard** tab: USB layout picker (`setkeyboardlayout` → `settings.json`; applies on next editor launch). Home from the editor → save + return to the Lobby; Home from the Lobby → quit (Writerdeck-server restarts xochitl but stays serving `:8000`). **Launch from stock UI** (no active session): phone **Edit** without a note, **USB Escape**, **left+right page buttons together**, Mac `wd` / `bash scripts/lobby.sh`, or on-tablet `~/wd` → Lobby. **Tablet vs browser:** see [browser-vs-tablet.md](browser-vs-tablet.md).
- Two page modes. Browse (Lobby / note list / Read-preview — no key capture, no echo footer) vs Type (active editing — capture + echo footer). Tapping Edit on a note enters Type mode and opens that note on the e-ink via Writerdeck's `saveAndLoad(name)`.
- Tablet → phone sync. The WebSocket also pushes server→browser: pressing Home or deleting the open note broadcasts `exitedit`, so the phone drops out of the typing view back to Browse in step with the tablet.
- GitHub note-sync (optional, off by default). **Writerdeck-server** is the sync engine (GitHub Contents API). The phone saves the PAT in browser `localStorage` and POSTs it to `/api/sync/token`; the server holds a copy **in RAM** for the service lifetime — never on disk. On page load, if sync is on but tablet RAM is empty, the phone re-posts the saved token automatically. Settings hold non-secret `syncOn`/`syncRepo` and per-note `syncMeta` (SHA + fingerprint). Reconcile runs on boot, 3 min poll, Home/power exit, tablet CRUD, phone CRUD, tablet **Sync now**, and token verify. Same integrity contract: edit-lease gate, empty-push guard, copy-on-clash (#19, #24). Plan: [server-sync-implementation.md](server-sync-implementation.md).
- IP is detected dynamically (`wlan0` first, then any up interface) and re-pushed to the Lobby when it changes — survives DHCP delay on boot and lease renewals. Last sync time on the Lobby comes from `lastSyncAt` in settings, updated when the server finishes a reconcile (`markSyncComplete` after `reconcileAll`).

## Constraints (honor these)

- No jailbreak; preserve OTA firmware updates ⇒ avoid Toltec (it locks the OS to a fixed range; can soft-brick on unsupported versions).
- No on-device runtime deps ⇒ static Go binary (`CGO_ENABLED=0`, `GOOS=linux GOARCH=arm GOARM=7`). The tablet ships no Python; installing it implies Entware/Toltec + a firmware lock.
- Markdown is the save format — see § Document integrity; HTML/qrichtext on disk is a bug, not a format option.
- Executable / device files = ASCII-only + LF (`.sh`, `.service`, `Dockerfile`, `.go`, `.yml`): a stray non-ASCII byte or CRLF breaks the device shell / systemd. (`.md` prose may use Unicode.) `.gitattributes` normalizes line endings.
- Keep the tablet awake — it drops Wi-Fi on suspend, which breaks the dev SSH / WebSocket connection.
- Latency is the e-ink refresh, not the LAN — don't over-engineer the transport.

> Escape hatch: the rM1's micro-USB is OTG-capable, so a plain USB keyboard drives Writerdeck directly if the Wi-Fi path ever stalls.

---

## Environment & facts

| Item | Value |
|---|---|
| Device | reMarkable 1 (first gen), codename *zero-gravitas* |
| OS / kernel | `20260506100933` · kernel `5.4.70-v1.6.3-rm10x` |
| `/dev/uinput` | Absent & un-addable (open → `ENODEV`; kernel exports trimmed via `CONFIG_TRIM_UNUSED_KSYMS`, so no out-of-tree `uinput.ko` can bind). Gate permanently 🔴 RED. Don't retry — the editor is fed over a socket instead; see [decisions.md](decisions.md). |
| SSH path | `ssh root@<tablet-ip>` over Wi-Fi (key login works) — the working path. USB (`10.11.99.1`) is dead on the Mac (no DHCP lease). Wi-Fi IP is DHCP; set `RM_HOST_WIFI` in `secrets/remarkable.local.env` (currently `192.168.1.8`). Reserve the tablet's MAC on the router so the IP stays put for the iPhone. **Lobby:** Mac `wd` / `bash scripts/lobby.sh`; on tablet SSH `~/wd` (`/home/root/wd`). |
| SSH password | gitignored in [../secrets/remarkable.local.env](../secrets/remarkable.local.env). Source: device `Settings → Help → Copyrights and licenses → General information`. Regenerates after every firmware update — re-record then. |
| Notes dir | `/home/root/Writerdeck-user-documents/` (Writerdeck boots to the Lobby; open from phone **Edit**, Lobby **Files** tab, or **Ctrl-K**). Deploy the binary to `/home/root/Writerdeck` — not `/home/root/Writerdeck-user-documents` (that's the notes *directory*). |
| Buttons | On `/dev/input/event1` (value 1 = press): middle/Home = `KEY_HOME` 102, left 105, right 106, power 116. Readable with xochitl up (Qt doesn't `EVIOCGRAB`). Writerdeck-server watches: **left+right together** → launch Lobby when idle; Home → relay while editing; Power → sleep/wake; USB keyboard evdev (name contains `keyboard`) → **Escape** launch when idle. |
| Disk | `/` rootfs (~228 MB) is 96% full — but nothing we ship goes there. The binary + Qt5 sysroot (~14 MB) + notes live on `/home/root/`, a separate multi-GB partition. Don't resize rootfs (A/B OTA scheme; brick risk). |

## Reference projects / links

- Editor (primary): https://github.com/dps/remarkable-keywriter — QML/C++. Input arrives via Qt QPA, not a direct `open()` (the input-injection patch point).
- Build toolchain: `ghcr.io/toltec-dev/qt:v3.3` — Qt 5.15.1 reMarkable qtbase, glibc 2.31, matches the device; bundles the closed-source `libqsgepaper`. `latest`=v4.0 is a trap (Qt6 + glibc 2.35).
- Input subsystem docs (evdev `input_event`): https://remarkable.guide/devel/device/input.html
- Prior art — uinput-free injection on rM1 (a documented fallback): https://github.com/machinelevel/sp425-crazy-cow — writes `input_event` to the *existing* Wacom pen node to draw glyphs → HWR ink, not clean Markdown.
- Own e-ink editor (last-resort fallback): https://github.com/canselcik/libremarkable (Rust framebuffer).
- Broad index: https://github.com/reHackable/awesome-reMarkable

---

## Build & deploy

Cross-compile + deploy Writerdeck-server from a host that can reach the tablet over Wi-Fi:

```bash
bash scripts/deploy-rmkbd.sh    # builds Writerdeck-server → /home/root/Writerdeck-server
```

Writerdeck (QML baked into binary — **rebuild before deploy**):

```bash
# CI
git push && bash scripts/fetch-keywriter-dist.sh && bash scripts/deploy-keywriter.sh -b
# or local Docker (Apple Silicon: --platform linux/amd64 on both docker commands — see scripts/README.md)
```

**`deploy-keywriter.sh` pushes `dist/Writerdeck`; it does not rebuild.** After `build-keywriter.sh` or `lobby/` edits, fetch or Docker-build first. `systemctl restart writerdeck` reloads the server only; relaunch the editor so the new binary loads. Check `journalctl -u writerdeck` for QML parse errors after deploy.

Requires Go (`brew install go` on macOS). `deploy-rmkbd.sh` cross-builds and deploys. Writerdeck is cross-built in CI (`third_party/keywriter/`, toltec Qt sysroot), not on a host toolchain.

> **After `deploy-rmkbd.sh`:** calls `POST /api/flush-save` then SIGTERM-waits up to ~12 s for graceful shutdown (save + quit). Restart with `systemctl start writerdeck`. Server-only — no `test-edit-session.sh` ([decisions.md](decisions.md) #21).

## Dev-loop shortcuts (aliases via `bash scripts/install-alias.sh`)

- `rmkw` (= `deploy-keywriter.sh -b`) — binary-only Writerdeck redeploy (~1 s): pushes `dist/Writerdeck`, skips the Qt5 sysroot tarball. **Requires a fresh binary in `dist/`** (CI fetch or local Docker rebuild after any `build-keywriter.sh` / `lobby/` change). **Then:** relaunch editor + `journalctl -u writerdeck` or `bash scripts/test-edit-session.sh`.
- `bash scripts/test-e2e.sh -s` — full browser→e-ink pipeline test, skipping the Writerdeck-server build+scp (~2 s; server already on device). Drop `-s` to rebuild+redeploy first.
- `bash scripts/test-edit-session.sh` — **Writerdeck/QML only** ([decisions.md](decisions.md) #21): phone **Edit** must keep Writerdeck up; not for `deploy-rmkbd.sh`-only changes.
- `rmpush "msg"` (= `push.sh`) — commit+push.
- `/home/root/Writerdeck-server -v` — per-key inject logging for keymap debugging; default is terse (connects + a count every 25 keys).
- SSH preflight pings first, so the scripts tell *tablet asleep* from *missing key*.

> Deploy transport: scp deadlocks at a fixed offset on this Mac→Wi-Fi→tablet link, so `rm_send_file` ([scripts/_env.sh](../scripts/_env.sh)) streams files gzip-over-ssh with a post-copy size check. See [decisions.md](decisions.md).

## Testing — the inner loop

Iterate over Wi-Fi SSH — set `RM_HOST_WIFI` in `secrets/remarkable.local.env`. Keep the tablet awake (it drops Wi-Fi on suspend). With `writerdeck.service`: logs in `journalctl -u writerdeck.service`.

| What you changed | Device verification |
|---|---|
| Writerdeck binary / QML patches | `bash scripts/test-edit-session.sh` |
| Writerdeck-server / `sync.js` / `app.js` only | Restart server (`systemctl start writerdeck`); API or browser smoke — **not** test-edit-session |
| Both binaries | test-edit-session **and** server smoke |

Manual ad-hoc: `systemctl stop xochitl` → run Writerdeck + Writerdeck-server → test → restore xochitl. Verify characters land *in Writerdeck*, not just that Writerdeck-server forwarded them.
