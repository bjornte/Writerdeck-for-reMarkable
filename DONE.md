# What's shipped

“Writerdeck for reMarkable 1” expands the tablet's functionality to include word processing.

Normally, reMarkable 1 does not support keyboards or word processing. This app offers both. There are two ways to connect a keyboard: Using Bluetooth via e.g. your phone (over WiFi) or via USB using an OTG cable.

The words land on e-ink and save as Markdown-supported text documents (`.md`).

Below is the feature list — what works today. Open work lives in [TODO.md](TODO.md). How it works: [docs/architecture.md](docs/architecture.md). Why we built it this way: [docs/decisions.md](docs/decisions.md). Hard-won gotchas from building it: [docs/lessons.md](docs/lessons.md).

---



## The core loop

Power on, the tablet boots into a distraction-free editor. Open `http://<tablet-ip>:8000/` in a browser, pick a note, type. Keystrokes travel over WebSocket to Writerdeck-server, which feeds a patched Writerdeck over `/run/Writerdeck.sock`. No app to install — the tablet is the server.

Home from the editor saves and returns to the Lobby. Home from the Lobby quits to the stock reMarkable UI. `Writerdeck-server` keeps serving `:8000` either way — relaunch from the phone (**Edit**), **Esc on a USB keyboard**, or **left+right page buttons together** (stock UI only).

## Lobby

Six-tab pager on e-ink: **Home · Files · Keyboard · Sync · Settings · Shortcuts** (touch tabs or keyboard Tab / arrows / digits 1–6). Fed by `pushLobbyInfo` → `{"t":"info",…}` on socket connect — IP, PIN, `syncOn`/`syncRepo`, note count, formatted last sync — and re-pushed when `wlan0` gets an address, a reconcile finishes, or notes change.

**Files tab (tablet-side CRUD, device-verified 2026-07-11):** list notes from Writerdeck-server over a trusted socket `{"t":"req","op":…}`; touch or USB keys to select; **New / Open / Rename / Delete** buttons; `n` / Enter / `r` / `d` shortcuts; double-tap or Enter to open. Opens via `saveAndLoad` (same path as phone **Edit**). Ctrl-K still opens the omni picker from any Lobby page. **Show PIN on tablet** on the phone drops back to the Lobby when a second device needs the PIN.

**Launch from stock UI:** Mac `wd` / `bash scripts/lobby.sh`; on tablet SSH `~/wd`. USB **Esc** or **left+right page buttons** together also open the Lobby when idle.

## Phone companion

**Notes.** List, create, read, rename, delete, upload `.md` files, download, and copy to clipboard — all from the browser.

**Browse vs Type.** The page has two modes. Browse is the file manager: no key capture, no echo footer. Type is active editing: your keystrokes go to the tablet and echo at the bottom.

**Edit on tablet.** Tap Edit on a note; it opens on the e-ink. The phone shows a typing view. Press Home on the tablet and the phone drops back to the note list in step.

**Paste.** Create a note from clipboard text, or insert at the cursor in an open note.

**Dark type mode.** Near-black UI while typing — easier on OLED phones.

## Security

PIN on the tablet, shown on e-ink each boot. Choose 6 digits, 4 digits, or none (none warns you that anyone on your Wi-Fi can connect). Wrong guesses from one IP lock that IP for 60 seconds. The auth cookie lasts until 04:00 local time.

## Preferences & sync (browser)

**Preferences** — reading font, PIN length, display rotate, and **Exit Writerdeck** (stops the service and returns the tablet to stock UI). **Sync** — separate panel for optional GitHub two-way sync. Both dismiss with Done, ×, ESC, or a backdrop click.

The top bar shows a single connection indicator — **Connected · 85%**, **Connecting…**, or **Tablet offline** — refreshed via `GET /api/status` every 5 s.

Sync warnings appear in a banner when sync is on but the browser has no token. The repo link opens on GitHub when configured.

The GitHub token never leaves the browser (`localStorage`). The tablet holds only `syncOn` and `syncRepo`.

## Editor (e-ink)

Built from source (upstream remarkable-keywriter), deployed as Writerdeck and patched for socket input. Renders full-panel via linuxfb. Norwegian and other Unicode characters work through the browser path.

**Reading view.** Pick Inter, Literata, EB Garamond, or DejaVu from the phone. Page up/down in read and edit (about four-fifths of a screen per step). Esc from edit keeps your scroll position — no jump to the bottom.

**Editing.** Block cursor hides while you type, reappears after a pause. Ctrl-K note switcher saves before switching. Mac-style navigation in edit mode (device-verified): Home/End line start/end (Cmd+Home/End doc start/end); Option/Alt+←/→ word; Cmd+←/→ line end; Cmd+↑/↓ doc start/end; Shift extends selection; plain ←/→ scroll the page. **Power button** while editing: save, sleep screen (“Writerdeck is sleeping… Press power to wake.”), suspend; press power again to wake and resume. **USB Escape** from stock reMarkable UI (no editor session): launches Writerdeck to the Lobby. **Left+right page buttons together** (stock UI, no USB): same launch. Wider margins, paper-like Lobby theme.

**Rotate.** Preferences → Display → **Rotate tablet 90°** (global — affects Lobby, read, and edit). USB in preview/read: Ctrl+←/→. USB in Lobby: **Ctrl+R** (90° clockwise) or Ctrl+←/→ (same as preview). Angle is stored as `"rotation"` in `.Writerdeck/settings.json`, pushed to Writerdeck on connect via socket `setrotation`, and restored after exit/relaunch. USB changes sync back to the server via `rotationChanged`.

## GitHub sync

Optional, off by default. The phone reconciles tablet notes with a private repo — pull what's missing either way, push local-only notes, handle clashes by keeping both copies with clear names. **Safety nets (2026-07-11):** refuses to push a zero-byte file over a previously-synced note; empty-tablet vs non-empty-GitHub clash restores from GitHub without creating `(tablet copy)` duplicates.

**Marker-aware delete** — a note deleted on GitHub (VS Code, web UI, git) propagates to the tablet when the local copy is pristine and carries a stored `sha`. Unpushed local edits resurrect instead of deleting. External renames reconcile as delete-old + pull-new. **Tablet delete/rename** queues `pendingSync` and pairs to GitHub when the phone browser is connected (slice 7); otherwise on next connect/reconcile.

Triggers: connect, toggle on, three-minute poll, manual Sync now, **tablet Home or Power** (full reconcile via phone browser). Each successful reconcile POSTs `/api/sync/ack`, which stores `lastSyncAt` in settings and refreshes the Lobby. Skips the note the tablet is editing — `openNote` from `/api/status` (edit lease, slices 1+4).

## Document integrity (slices 1–11, 2026-07-11)

Non-negotiable contract: plain Markdown on disk, no silent overwrite of live edits, durable saves. Device-verified after `b1ce2bc`…`f72282d`.

| Slice | What |
|---|---|
| 1 | Edit lease — `notifyOpen`, `openedit` WS; reconcile skips open note |
| 2 | Content fidelity — markdown save contract, HTML guard, `toggleMode` fix |
| 3 | `notedeleted` / `noterenamed` on phone ops against open file |
| 4 | Reconcile gated on `openNote` in `/api/status` |
| 5 | OCC — ETag / `If-Match` on PUT |
| 6 | Atomic server writes — `writeNoteFile` temp+rename |
| 7 | Tablet CRUD → GitHub — `tabletcrud` WS + `pendingSync` |
| 8 | Disk drift — `diskchanged` WS, phone banner, `/api/reload` |
| 9 | 45 s autosave while editing |
| 10 | Tablet saves via loopback `PUT /api/notes` |
| 11 | Save before deploy/stop — `POST /api/flush-save`, graceful shutdown |

Residual risks and unknowns: [docs/integrity-audit.md](docs/integrity-audit.md).

## Infrastructure

Static Go binary (`Writerdeck-server` at `/home/root/Writerdeck-server`), no on-device runtime deps. Cross-built on the Mac, deployed over Wi-Fi. Writerdeck built in CI with the toltec Qt sysroot. Cold-boot autostart via `writerdeck.service`. Keep-awake during editor sessions only. On-device layout documented in [docs/architecture.md](docs/architecture.md).

**Regression test:** `bash scripts/test-edit-session.sh` — POST `/api/open` (phone **Edit**) must leave Writerdeck running, xochitl down, and `editorActive: true` for several seconds. Logs to `docs/recon/`.

`/dev/uinput` is dead on this kernel — we feed the editor over a socket instead. That path is settled; don't retry uinput.