# Writerdeck for reMarkable 1

A distraction-free word processor for the first generation reMarkable. Supports Bluetooth and USB keyboards. Optionally syncs your notes to a private GitHub repository of your choice. Saves files as Markdown.

Natively, the reMarkable 1 supports the "draw", "write by hand" and "read" use cases. With this app, "use as typewriter" is also supported.

Bluetooth keyboards pair to your phone and bridge over Wi-Fi. USB keyboards use an [OTG cable](https://en.wikipedia.org/wiki/USB_On-The-Go#OTG_micro_cables).

![Writerdeck for reMarkable 1](docs/Writerdeck-for-reMarkable.jpg)

The reMarkable 1 has a large e-ink screen and a quiet OS, but no word processor and no keyboard support. This fills the gap.

The project is heavily LLM-assisted and partly human-reviewed. Primary sources: Singleton’s [keywriter](https://github.com/dps/remarkable-keywriter) (forked as [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter)) and ideas from [crazy-cow](https://github.com/machinelevel/sp425-crazy-cow).

## Status

A usable appliance. Plenty left to improve; the core loop works.

## Install (expert)

1. Clone the repo. Copy [secrets/remarkable.local.env.example](secrets/remarkable.local.env.example) to `remarkable.local.env` and fill in the device password (tablet: Settings → Help → Copyrights and licenses → General information).
2. `bash scripts/bootstrap.sh` — installs your SSH key.
3. `bash scripts/fetch-keywriter-dist.sh` — pulls the CI-built editor (`gh` required).
4. `bash scripts/deploy-keywriter.sh` — ships the editor.
5. `bash scripts/deploy-rmkbd.sh` — ships Writerdeck-server.
6. `bash scripts/install-service.sh`, then on the tablet `systemctl start writerdeck` to test and `systemctl enable writerdeck` only after that works.

### Optional GitHub sync

Use a private personal repo. Conflicts keep both copies rather than overwrite. Set a fine-grained token with Contents read/write on that repo only. On the phone: Sync setup — turn sync on, enter `owner/repo`, paste the token. The token stays in the browser; a new Wi-Fi address is a new browser origin, so you may need to enter it again there.

![Create token](docs/create-token.png)

## Everyday use

Power on — Lobby Files shows the connect address and PIN (also on Home, digit 6). Open that address on the phone, enter the PIN, pair a keyboard to the phone if you like. Open a note on the tablet; the phone enters Type mode. Upload and Download live on the phone list. Paste from phone inserts at the cursor. Font, PIN length, and rotation live in Lobby Settings.

Show the Lobby from a Mac on the same Wi-Fi with `wd` (after `bash scripts/install-alias.sh`) or `bash scripts/lobby.sh`. On the tablet: `~/wd`.

Useful keys: Esc toggles edit and preview inside Writerdeck, or launches to Lobby from the stock UI with a USB keyboard. Left and right page buttons together do the same launch without USB. Ctrl-K switches notes. Ctrl-R rotates. Home from edit returns to Files; Home from Lobby quits to the stock UI.

## For developers

Start with [TODO.md](TODO.md) and [DONE.md](DONE.md). Credentials: [secrets/README.md](secrets/README.md). Keep the tablet awake and iterate over Wi-Fi. After editor source changes: push, fetch the CI binary, deploy, run `test-edit-session.sh`.

## Constraints

No jailbreak; keep OTA (over-the-air updates) — so no Toltec. One static Go binary on the tablet.

## Pieces

Writerdeck-server — always-on Go daemon: phone page, WebSocket, files, sync, PIN, key relay into `/run/Writerdeck.sock`.

Phone page — captures keys and talks to the server.

Writerdeck — full-screen editor from our keywriter fork. Saves Markdown under `Writerdeck-user-documents/`.

Keys use the socket because this kernel cannot load uinput. More: [docs/architecture.md](docs/architecture.md), [docs/decisions.md](docs/decisions.md).

## License

[MIT](LICENSE) © 2026 Bjørn Tennøe. Keywriter is third-party with its own license.
