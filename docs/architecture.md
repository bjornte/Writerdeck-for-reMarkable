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
   - rmkbd daemon (Go): serves capture page + WebSocket, feeds keys to a local socket
   - keywriter (patched): reads the socket, full-screen Markdown editor, saves .md
```

Keystrokes reach the editor through a local socket rather than `/dev/uinput`: this tablet's kernel can't load uinput, so the daemon feeds a patched keywriter instead (the reasoning is in [decisions.md](decisions.md)).

## Components

- Component A — `rmkbd` daemon (we build): a static Go binary on the tablet. Serves a tiny HTML capture page + a WebSocket, and forwards received key events into a local socket the editor reads — no `/dev/uinput`. It is always-on (its own systemd service) and also owns the file API, the per-boot PIN, and the xochitl↔keywriter lifecycle (below).
- Component B — editor ([`dps/remarkable-keywriter`](https://github.com/dps/remarkable-keywriter), patched): a full-screen, distraction-free Markdown editor; saves `.md` to `/home/root/edit/`. Input arrives via Qt QPA (there is no `/dev/input` fd to swap), so we inject synthetic `QKeyEvent`s from a small socket-reader thread added to `main.cpp`.
- Component C — client (we build, trivial): an HTML/JS page served by the daemon that captures `keydown`/`keyup` and pushes them over the WebSocket. The Mac plays this role during dev; the iPhone + a Bluetooth/Lightning keyboard at the end.

## The Companion (the working appliance)

The tablet is the web server — the phone/Mac just open Safari, no app to install.

- Always-on daemon, on-demand editor (the lifecycle split). `rmkbd` keeps serving `:8000` even under the stock xochitl GUI, and summons keywriter on demand: stop xochitl → spawn keywriter → on Home, restart xochitl and keep serving. Boot auto-launches one editor session (power-on = typewriter). This is what lets the phone re-launch writing and manage files while the GUI is up.
- Notes = `/home/root/edit/*.md`. `rmkbd` (Go) does all file ops natively — list/read/create/rename/delete over an HTTP API (`/api/notes…`; `notesSafe()` rejects `/` and `..`). These are keywriter's notes, separate from xochitl's document store, so the API is safe to hit while the GUI is up.
- PIN-on-tablet auth (length owner-choosable). `rmkbd` mints a random PIN per boot and shows it in the Lobby; the phone POSTs it (`/api/pin`) for an HttpOnly `rmkbd_token` cookie that gates the notes API and the WebSocket upgrade. The length is owner-choosable on the settings screen — 6 / 4 / none (none = LAN-open, UI-warned), and the PIN is runtime-mutable: changing it re-mints on the spot and re-shows it on the e-ink. A per-IP brute-force lockout (5 wrong / 60 s → 429) backs the 4- and 6-digit modes.
- Settings screen (phone). A **Settings** overlay (font, PIN, display rotate) and a separate **Sync** overlay (GitHub sync) in the top bar; both dismiss via Done, ×, ESC, or backdrop click. Font choice (Inter / Literata / EB Garamond / DejaVu, pushed live to the editor) and PIN length persist to `/home/root/.rmkbd/settings.json`. When sync is on and a repo is set, the Sync panel links to `github.com/{repo}`. A missing token shows a yellow banner below the status bar.
- Lobby-on-demand. A second device that arrives *after* the owner is already editing finds the tablet showing the note, not the PIN; a pre-auth "Show PIN on tablet" button (`POST /api/lobby`, rate-limited) saves the open note and drops the tablet to the Lobby so the PIN is readable — it reveals nothing over the wire (PIN only on the e-ink).
- Lobby + two-level Home. Boot shows a full-screen Lobby (IP + PIN + how-to, fed by `{"t":"info",…}` on socket connect and re-pushed when `wlan0` gets an address). Home from the editor → save + return to the Lobby; Home from the Lobby → quit (rmkbd restarts xochitl but stays serving `:8000`).
- Two page modes. Browse (Lobby / note list / Read-preview — no key capture, no echo footer) vs Type (active editing — capture + echo footer). Tapping Edit on a note enters Type mode and opens that note on the e-ink via keywriter's existing `doLoad(name)`.
- Tablet → phone sync. The WebSocket also pushes server→browser: pressing Home or deleting the open note broadcasts `exitedit`, so the phone drops out of the typing view back to Browse in step with the tablet.
- GitHub note-sync (optional, off by default). The phone browser is the sync engine — the GitHub token lives in its `localStorage`, never on the tablet (which holds only non-secret `syncOn`/`syncRepo`). It reconciles by *unioning* the tablet + repo note lists and copying whatever's missing either way; it never deletes on its own (safer: it can't lose a note, and it dodges real-git-on-mobile instability by using GitHub's plain Contents API). Caveat: delete/rename must be done *in the phone browser* (which pairs the op to GitHub); a note deleted or renamed elsewhere — VS Code, `git`, the GitHub web UI — resurrects or duplicates on the next sync. See [decisions.md](decisions.md) #19.
- IP is detected dynamically (`wlan0` first, then any up interface) and re-pushed to the Lobby when it changes — survives DHCP delay on boot and lease renewals.

## Constraints (honor these)

- No jailbreak; preserve OTA firmware updates ⇒ avoid Toltec (it locks the OS to a fixed range; can soft-brick on unsupported versions).
- No on-device runtime deps ⇒ static Go binary (`CGO_ENABLED=0`, `GOOS=linux GOARCH=arm GOARM=7`). The tablet ships no Python; installing it implies Entware/Toltec + a firmware lock.
- Markdown is the save format.
- Executable / device files = ASCII-only + LF (`.sh`, `.service`, `Dockerfile`, `.go`, `.yml`): a stray non-ASCII byte or CRLF breaks the device shell / systemd. (`.md` prose may use Unicode.) `.gitattributes` normalizes line endings.
- Keep the tablet awake — it drops Wi-Fi on suspend, which breaks the dev SSH / WebSocket connection.
- Latency is the e-ink refresh, not the LAN — don't over-engineer the transport.

> Escape hatch: the rM1's micro-USB is OTG-capable, so a plain USB keyboard drives keywriter directly if the Wi-Fi path ever stalls.

---

## Environment & facts

| Item | Value |
|---|---|
| Device | reMarkable 1 (first gen), codename *zero-gravitas* |
| OS / kernel | `20260506100933` · kernel `5.4.70-v1.6.3-rm10x` |
| `/dev/uinput` | Absent & un-addable (open → `ENODEV`; kernel exports trimmed via `CONFIG_TRIM_UNUSED_KSYMS`, so no out-of-tree `uinput.ko` can bind). Gate permanently 🔴 RED. Don't retry — the editor is fed over a socket instead; see [decisions.md](decisions.md). |
| SSH path | `ssh root@<tablet-ip>` over Wi-Fi (key login works) — the working path. USB (`10.11.99.1`) is dead on the Mac (no DHCP lease). Wi-Fi IP is DHCP; set `RM_HOST_WIFI` in `secrets/remarkable.local.env` (currently `10.0.0.20`). Reserve the tablet's MAC on the router so the IP stays put for the iPhone. |
| SSH password | gitignored in [../secrets/remarkable.local.env](../secrets/remarkable.local.env). Source: device `Settings → Help → Copyrights and licenses → General information`. Regenerates after every firmware update — re-record then. |
| Notes dir | `/home/root/edit/` (keywriter boots to the Lobby; files are opened from the phone via the companion page). Deploy the binary to `/home/root/keywriter` — not `/home/root/edit` (that's the notes *directory*). |
| Buttons | On `/dev/input/event1` (value 1 = press): middle/Home = `KEY_HOME` 102, left 105, right 106, power 116. Readable with xochitl up (Qt doesn't `EVIOCGRAB`). |
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

Cross-compile + deploy `rmkbd` from a host that can reach the tablet over Wi-Fi:

```bash
bash scripts/deploy-rmkbd.sh    # GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build + gzip-over-ssh to /home/root/
```

Requires Go (`brew install go` on macOS). `deploy-rmkbd.sh` cross-builds and deploys; `build-rmkbd.ps1` cross-builds only (no device step). keywriter is cross-built in CI (the toltec Qt sysroot), not on a host toolchain.

## Dev-loop shortcuts (aliases via `bash scripts/install-alias.sh`)

- `rmkw` (= `deploy-keywriter.sh -b`) — binary-only keywriter redeploy (~1 s): pushes just the ~205 K binary, skips re-throwing the 14 MB Qt5 sysroot (static; only changes on a Qt rebuild — then pass `RM_FORCE_SYSROOT=1`). Use after any `socket-inject.patch` / `build-keywriter.sh` change once CI has rebuilt.
- `bash scripts/test-e2e.sh -s` — full browser→e-ink pipeline test, skipping the rmkbd build+scp (~2 s; rmkbd already on device). Drop `-s` to rebuild+redeploy rmkbd first.
- `rmpush "msg"` (= `push.sh`/`push.ps1`) — commit+push under the personal identity.
- `rmkbd -v` — per-key inject logging for keymap debugging; default is terse (connects + a count every 25 keys).
- SSH preflight pings first, so the scripts tell *tablet asleep* from *missing key*.

> Deploy transport: scp deadlocks at a fixed offset on this Mac→Wi-Fi→tablet link, so `rm_send_file` ([scripts/_env.sh](../scripts/_env.sh)) streams files gzip-over-ssh with a post-copy size check. See [decisions.md](decisions.md).

## Testing — the inner loop

Iterate over Wi-Fi SSH — set `RM_HOST_WIFI` in `secrets/remarkable.local.env`. Keep the tablet awake (it drops Wi-Fi on suspend). `systemctl stop xochitl` → run `keywriter` + `rmkbd` → test → `systemctl start xochitl` to restore the stock UI. Run the daemon in the foreground with `-v` logging of every received/injected event. Verify one capability on the device — confirm characters land *in keywriter*, not just that the daemon "sent" them; when one is wrong, log incoming vs emitted and fix the keymap entry.
