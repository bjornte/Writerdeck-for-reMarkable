# Writerdeck for reMarkable 1

A distraction-free word processor for the reMarkable 1, with support for hardware keyboards (USB & Bluetooth), optional syncing of notes, and Markdown.

* Bluetooth keyboards: Pair to your phone, then bridge to reMarkable over Wi-Fi
* USB keyboards: Use an [OTG cable](https://en.wikipedia.org/wiki/USB_On-The-Go#OTG_micro_cables)

![Writerdeck for reMarkable 1](docs/Writerdeck-for-reMarkable.jpg)

Background: The reMarkable 1 has a large, nice e-ink screen and a distraction-free OS, but no word processor, and no support for hardware keyboards. This Writerdeck fills the gap.

The project is 90% vibe coded and 50% human reviewed. I (the repo owner, Bjørn) have mostly just herded LLMs and perused the documentation, not the code itself. Primary sources: the [keywriter](https://github.com/dps/remarkable-keywriter) editor (slightly patched) and a keypress injection approach in [crazy-cow](https://github.com/machinelevel/sp425-crazy-cow).

## Status

A usable appliance. There's an abundance of possible improvements, but the core deliverable is working.

## Installation guide for (expert) users

1. Clone this repo. Copy [secrets/remarkable.local.env.example](secrets/remarkable.local.env.example) to `remarkable.local.env` and fill in the device password (tablet: Settings → Help → Copyrights and licenses → General information).
2. `bash scripts/bootstrap.sh` — installs your SSH key on the tablet.
3. `bash scripts/fetch-keywriter-dist.sh` — downloads the CI-built editor (keywriter binary + Qt runtime) into `third_party/keywriter/dist/`. Source-only mirror: these aren't committed, so fetch them from CI first (needs `gh`: `brew install gh && gh auth login`).
4. `bash scripts/deploy-keywriter.sh` — ships the editor to the tablet.
5. `bash scripts/deploy-rmkbd.sh` — cross-builds and ships Writerdeck-server to `/home/root/Writerdeck-server`.
6. `bash scripts/install-service.sh` (on the Mac) installs the systemd unit. Then SSH into the tablet (`ssh root@<ip>`) and run `systemctl start writerdeck` to test, then `systemctl enable writerdeck` to boot straight into it. Enable only after the test passes — see the script's boot-loop note.

### Optional: GitHub syncing of notes

Optionally, sync your notes towards GitHub. The assumed use case is to use a repo that's personal & private.

Edit conflicts never overwrite. The reMarkable's version is kept as `note (tablet copy).md`. A banner appears on the phone so you can reconcile.

If a note has edits on the reMarkable that haven't synced yet, deleting or renaming the note elsewhere keeps the note rather than removing it.

To set up the repo and enable syncing:

1. Here on GitHub, create a new private repo to hold your notes
2. Go to the [create token](https://github.com/settings/personal-access-tokens/new) page. Create a fine-grained personal access token with Repository access limited to just that repo and `Repository permissions` → `Contents: Read and write`. Copy the token.
3. On your phone, open **Sync** → GitHub sync: turn it on, enter the repo as `owner/repo`, and paste the token.
    * For security reasons, the token is kept in the phone's browser only, and never reaches the reMarkable. The browser ties the token to the URL. Therefore, the token must be reentered for every Wi-Fi network the reMarkable should sync over, e.g. home, work and mobile hotspot.
4. While the reMarkable never sees the token, it records whether sync is enabled and the name of the repo.

![Create token](docs/create-token.png)

## How-to for users incl. shortcuts

1. Power on the tablet — it boots into a Lobby showing the connect URL (`http://<ip>:8000`), PIN, and setup sections.
2. Open that address in the phone's browser, and enter the PIN.
3. Pair a physical keyboard to your phone.
4. Open a note on the tablet Files tab (Edit, Enter, or double-tap) — the phone enters Type mode and keystrokes land on e-ink. **Upload** and **Download** are on the phone note list. In Type mode, **Paste from phone** inserts clipboard text at the cursor. Reading font and PIN length live on the tablet Lobby Settings tab. Rotate the display there too (Ctrl-R or Ctrl-arrows).

**Show the Lobby**

| Where | Command |
|-------|---------|
| Mac (same Wi-Fi) | `wd` after `bash scripts/install-alias.sh`, or `bash scripts/lobby.sh` |
| SSH on the tablet | `~/wd` |

Both invoke `/home/root/wd`: start Writerdeck if needed, then open the Lobby. No PIN required. (`wd` is not on `$PATH`; use `~/wd` or the full path.)

Shortcuts:

- Esc — in Writerdeck: toggles edit / preview (or closes the note picker). From **stock reMarkable UI** with a USB keyboard: launches Writerdeck to the Lobby. **Left + right page buttons together** (no USB): same launch from stock UI.
- Ctrl-K — switch note from within edit view; from the Lobby, opens the note picker.
- Ctrl-R — rotate 90° clockwise in the Lobby (USB keyboard). Ctrl + side arrows rotate in preview/read mode and in the Lobby (angle is remembered across sessions).
- Ctrl-Q — quit Writerdeck from the Lobby.
- Home — from the editor: save and return to the Lobby; from the Lobby: exit to stock reMarkable UI.

## Getting started for devs

Development on the tablet is done over SSH from a machine on the same Wi-Fi. To get started:

1. [TODO.md](TODO.md), [DONE.md](DONE.md) etc. briefs on current status.
2. Create your local credentials: copy [secrets/remarkable.local.env.example](secrets/remarkable.local.env.example) to `remarkable.local.env` and fill in the device password — see [secrets/README.md](secrets/README.md).
3. Run `bash scripts/bootstrap.sh`, then `bash scripts/recon.sh`. Keep the tablet awake, and iterate over Wi-Fi.
4. After keywriter changes: `git push` → `fetch-keywriter-dist.sh` → `rmkw` → `test-edit-session.sh` (CI runs Docker; Mac does not).

## Design constraints

- No jailbreak, and preserve over-the-air firmware updates — so no Toltec.
- No runtime dependencies on the tablet — just one static Go binary (`CGO_ENABLED=0`, ARMv7).


## Main components

Three pieces — the server and client are built here, the editor is third-party (patched):

- **Writerdeck-server** — a small, static Go daemon at `/home/root/Writerdeck-server`. It serves an HTML capture page and a WebSocket, then forwards keystrokes into `/run/Writerdeck.sock`.
- **the client** — a browser page (served by Writerdeck-server) that captures keystrokes and sends them over the LAN.
- **Writerdeck** — the third-party [remarkable-keywriter](https://github.com/dps/remarkable-keywriter) editor, patched to read that socket. A full-screen, distraction-free Markdown editor that saves `.md` to `Writerdeck-user-documents/`.

Keystrokes reach the editor through a local socket rather than `/dev/uinput`: this tablet's kernel can't load uinput, so Writerdeck-server feeds the patched editor instead. The reasoning is in [docs/decisions.md](docs/decisions.md).


## Repo layout

- [Architecture](docs/architecture.md)
- [Architecture decision record (ADR)](docs/decisions.md)
- [Todo](TODO.md)
- [Done](DONE.md)

| Path | What's there |
|---|---|
| [daemon/](daemon/) | Go source for Writerdeck-server: WebSocket, editor-feed socket, embedded capture page |
| [third_party/](third_party/) | Upstream keywriter tree; CI builds the `Writerdeck` binary |
| [scripts/](scripts/) | Bash automation — bootstrap, recon, deploy, test |
| [docs/](docs/) | Architecture, decisions, lessons |
| [secrets/](secrets/) | Local credentials — gitignored; see [secrets/README.md](secrets/README.md) |


## License

[MIT](LICENSE) © 2026 Bjørn Tennøe — permissive, no warranty. [keywriter](https://github.com/dps/remarkable-keywriter) is third-party with its own license, not covered by this claim.
