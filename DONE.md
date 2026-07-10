# What's shipped

“Writerdeck for reMarkable 1” expands the tablet's functionality to include word processing.

Normally, reMarkable 1 does not support keyboards or word processing. This app offers both. There are two ways to connect a keyboard: Using Bluetooth via e.g. your phone (over WiFi) or via USB using an OTG cable.

The words land on e-ink and save as Markdown-supported text documents (`.md`).

Below is the feature list — what works today. Open work lives in [TODO.md](TODO.md). How it works: [docs/architecture.md](docs/architecture.md). Why we built it this way: [docs/decisions.md](docs/decisions.md). Hard-won gotchas from building it: [docs/lessons.md](docs/lessons.md).

---



## The core loop

Power on, the tablet boots into a distraction-free editor. Open `http://<tablet-ip>:8000/` in a browser, pick a note, type. Keystrokes travel over WebSocket to `rmkbd`, which feeds a patched keywriter over a local socket. No app to install — the tablet is the server.

Home from the editor saves and returns to the Lobby. Home from the Lobby quits to the stock reMarkable UI. `rmkbd` keeps serving `:8000` either way, so you can launch writing again without a reboot.

## Lobby

Full-screen welcome on e-ink: project title, connect URL, PIN (if enabled), and how-to text. **Open note… (Ctrl-K)** button opens the note switcher from the Lobby; Ctrl-K on a USB keyboard does the same. The URL tracks Wi-Fi — when `wlan0` gets an address after boot, the Lobby updates. A second device can ask the tablet to show the PIN again (`Show PIN on tablet` on the phone). When GitHub sync is on, the Lobby also shows the notes repo (`Sync: github.com/owner/repo`).

## Phone companion

**Notes.** List, create, read, rename, delete, upload `.md` files, download, and copy to clipboard — all from the browser.

**Browse vs Type.** The page has two modes. Browse is the file manager: no key capture, no echo footer. Type is active editing: your keystrokes go to the tablet and echo at the bottom.

**Edit on tablet.** Tap Edit on a note; it opens on the e-ink. The phone shows a typing view. Press Home on the tablet and the phone drops back to the note list in step.

**Paste.** Create a note from clipboard text, or insert at the cursor in an open note.

**Dark type mode.** Near-black UI while typing — easier on OLED phones.

## Security

PIN on the tablet, shown on e-ink each boot. Choose 6 digits, 4 digits, or none (none warns you that anyone on your Wi-Fi can connect). Wrong guesses from one IP lock that IP for 60 seconds. The auth cookie lasts until 04:00 local time.

## Settings & sync (browser)

**Settings** — reading font, PIN length, display rotate, and **Exit Writerdeck** (stops the service and returns the tablet to stock UI). **Sync** — separate panel for optional GitHub two-way sync. Both dismiss with Done, ×, ESC, or a backdrop click.

The top bar shows tablet battery and Wi-Fi (`96% · Wi-Fi`), refreshed every 30 seconds.

Sync warnings appear in a banner when sync is on but the browser has no token. The repo link opens on GitHub when configured.

The GitHub token never leaves the browser (`localStorage`). The tablet holds only `syncOn` and `syncRepo`.

## Editor (e-ink)

Built from source keywriter, patched for socket input. Renders full-panel via linuxfb. Norwegian and other Unicode characters work through the browser path.

**Reading view.** Pick Inter, Literata, EB Garamond, or DejaVu from the phone. Page up/down in read and edit (about four-fifths of a screen per step). Esc from edit keeps your scroll position — no jump to the bottom.

**Editing.** Block cursor hides while you type, reappears after a pause. Ctrl-K note switcher saves before switching. Mac-style navigation in edit mode (device-verified): Home/End line start/end (Cmd+Home/End doc start/end); Option/Alt+←/→ word; Cmd+←/→ line end; Cmd+↑/↓ doc start/end; Shift extends selection; plain ←/→ scroll the page. **Power button** while editing: save, sleep screen (“Writerdeck is sleeping…”), suspend; press power or **Esc** (USB keyboard) again to wake and resume. **USB Escape** when idle (stock UI / no editor session): launches Writerdeck to the Lobby. Wider margins, paper-like Lobby theme.

**Rotate.** Settings → Display → **Rotate tablet 90°** (global — affects Lobby, read, and edit). USB Ctrl+←/→ still works in preview mode. Calls `POST /api/rotate` → socket `rotate` cmd → `root.rotation`.

## GitHub sync

Optional, off by default. The phone reconciles tablet notes with a private repo — pull what's missing either way, push local-only notes, handle clashes by keeping both copies with clear names.

**Marker-aware delete** — a note deleted on GitHub (VS Code, web UI, git) propagates to the tablet when the local copy is pristine and carries a stored `sha`. Unpushed local edits resurrect instead of deleting. External renames reconcile as delete-old + pull-new. Tablet-only deletes still don't propagate to GitHub by design.

Triggers: connect, toggle on, three-minute poll, manual Sync now, **tablet Home or Power** (full reconcile via phone browser). Skips the note currently open on the tablet during reconcile — `tabletOpenNote` is cleared before sync runs.

## Infrastructure

Static Go binary (`rmkbd`), no on-device runtime deps. Cross-built on the Mac, deployed over Wi-Fi. keywriter built in CI with the toltec Qt sysroot. Cold-boot autostart via systemd. Keep-awake during editor sessions only.

`/dev/uinput` is dead on this kernel — we feed the editor over a socket instead. That path is settled; don't retry uinput.