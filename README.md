# Writerdeck for reMarkable 1

A distraction-free word processor for the first generation reMarkable. Supports Bluetooth and USB keyboards. Optionally syncs your notes to a private GitHub repository of your choice. Optionally encrypts files. Saves files as Markdown.

Natively, the reMarkable 1 supports the "draw", "write by hand" and "read" use cases. With this app, "use as typewriter" is also supported.

Bluetooth keyboards pair to your phone and bridge over Wi-Fi. USB keyboards use an [OTG cable](https://en.wikipedia.org/wiki/USB_On-The-Go#OTG_micro_cables).

![Writerdeck for reMarkable 1](docs/Writerdeck-for-reMarkable.jpg)

The reMarkable 1 has a large e-ink screen and a quiet OS, but no word processor and no keyboard support. This fills the gap.

The project is heavily LLM-assisted and partly human-reviewed. Primary sources: Singleton’s [keywriter](https://github.com/dps/remarkable-keywriter) (forked as [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter)) and ideas from [crazy-cow](https://github.com/machinelevel/sp425-crazy-cow).

## Status

A usable appliance. Plenty left to improve; the core loop works.

## Install

### Before you start

- reMarkable **1** only (not 2 / Paper Pro)
- Mac or Linux on the **same Wi-Fi** as the tablet (scripts are bash; no Windows-native path yet)
- **Go 1.25+** (`go version`) to build Writerdeck-server
- Tablet **awake** (sleep drops Wi-Fi)

Editor binaries come from the [`keywriter` Release](https://github.com/bjornte/Writerdeck-for-reMarkable/releases/tag/keywriter) (`bash scripts/fetch-keywriter-dist.sh` uses curl). If that fails: open [Actions → Build keywriter](https://github.com/bjornte/Writerdeck-for-reMarkable/actions/workflows/build-keywriter.yml), pick the latest green run, download the `keywriter-dist` artifact, and put `Writerdeck` + `qt5.tar.gz` in `third_party/keywriter/dist/`. Optional: `gh` for a specific Actions run id.

### Steps

1. Clone the repo. Copy the secrets template and fill in **password** and **Wi-Fi IP**:

   ```bash
   cp secrets/remarkable.local.env.example secrets/remarkable.local.env
   ```

   - `RM_ROOT_PASSWORD` — tablet: Settings → Help → Copyrights and licenses → General information (scroll down). Changes after every firmware update.
   - `RM_HOST_WIFI` — tablet Wi-Fi settings (or your router’s client list). Prefer a DHCP reservation so the address stays put. More: [docs/architecture.md](docs/architecture.md).

2. One-shot install (recommended):

   ```bash
   bash scripts/install.sh
   ```

   Or step by step: `preflight.sh` → `bootstrap.sh` → `fetch-keywriter-dist.sh` → `deploy-keywriter.sh` → `deploy-rmkbd.sh` → `install-service.sh`.

3. Start the service (boot-loop guard — do **not** enable until this works):

   ```bash
   bash scripts/install-service.sh --start
   ```

   Or SSH to the tablet and run `systemctl start writerdeck`. Only then: `systemctl enable writerdeck`.

Optional: `bash scripts/install.sh --start` runs the start smoke test at the end of the chain.

### You're done when

- Lobby Files is on the e-ink screen
- Phone loads `http://<RM_HOST_WIFI>:8000/` — notes list populated, connection bar shows **Connected** or **Tablet offline**, not stuck on `connecting...`

Optional smoke test from the Mac: `bash scripts/test-edit-session.sh`.

### After a firmware update (OTA)

The SSH password changes and the systemd unit may be gone. Update `RM_ROOT_PASSWORD`, then:

```bash
bash scripts/fix-hostkey.sh    # if SSH says "host key changed"
bash scripts/bootstrap.sh
bash scripts/deploy-keywriter.sh
bash scripts/deploy-rmkbd.sh
bash scripts/install-service.sh --start
```

Enable again only after a manual start works.

### Recovery (bad autostart)

On the tablet:

```bash
systemctl disable --now writerdeck && systemctl start xochitl
```

### Optional GitHub sync

Use a private personal repo. Conflicts keep both copies rather than overwrite. Set a fine-grained token with Contents read/write on that repo only. On the phone: Sync setup — turn sync on, enter `owner/repo`, paste the token. The token stays in the browser; a new Wi-Fi address is a new browser origin, so you may need to enter it again there.

![Create token](docs/create-token.png)

## Everyday use

Power on — Lobby Files shows the connect address and PIN (also on Home, digit 6). Open that address on the phone, enter the PIN, pair a keyboard to the phone if you like. Open a note on the tablet; the phone enters Type mode. Upload and Download live on the phone list. Paste from phone inserts at the cursor. Font, PIN length, and rotation live in Lobby Settings.

Show the Lobby from a Mac on the same Wi-Fi with `wd` (after `bash scripts/install-alias.sh`) or `bash scripts/lobby.sh`. On the tablet: `~/wd`.

Useful keys: Esc toggles edit and preview inside Writerdeck, or launches to Lobby from the stock UI with a USB keyboard. Left and right page buttons together do the same launch without USB. Ctrl-K switches notes. Ctrl-R rotates. Home from edit returns to Files; Home from Lobby quits to the stock UI.

## For developers

Start with [TODO.md](TODO.md) and [DONE.md](DONE.md). Credentials: [secrets/README.md](secrets/README.md). Keep the tablet awake and iterate over Wi-Fi. After editor source changes: push, fetch the CI binary, deploy, run `test-edit-session.sh`. First-time install planning: [docs/install-onboarding/](docs/install-onboarding/).

## Constraints

No jailbreak; keep OTA (over-the-air updates) — so no Toltec. One static Go binary on the tablet.

## Pieces

Writerdeck-server — always-on Go daemon: phone page, WebSocket, files, sync, PIN, key relay into `/run/Writerdeck.sock`.

Phone page — captures keys and talks to the server.

Writerdeck — full-screen editor from our keywriter fork. Saves Markdown under `Writerdeck-user-documents/`.

Keys use the socket because this kernel cannot load uinput. More: [docs/architecture.md](docs/architecture.md), [docs/decisions.md](docs/decisions.md).

## License

[MIT](LICENSE) © 2026 Bjørn Tennøe. Keywriter is third-party with its own license.
