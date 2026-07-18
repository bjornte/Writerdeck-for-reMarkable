# Reinstall cheat sheet

Wipe the tablet (or fresh device), then get notes and vault back. Installer: [README.md](../../README.md). Regression: `bash scripts/test-install-reuse.sh`.

## Before you wipe

Keep a copy of `secrets/remarkable.local.env` on the computer (gitignored). It should contain:

| Field | What it is |
|---|---|
| `RM_ROOT_PASSWORD` | Tablet SSH password (**changes after wipe / OTA** — read the new one from the tablet) |
| `RM_HOST_WIFI` | Tablet Wi-Fi address |
| `PIN_DIGITS` | `6`, `4`, or `none` |
| `SYNC_REPO` | Notes repo, e.g. `owner/my-notes` |
| `GH_TOKEN` | Fine-grained GitHub token (Contents read/write on that repo) |

Also confirm the GitHub repo still has your notes, any `.md.enc` files, and `secret/vault` + `secret/pin`.

## After wipe / on a new tablet

1. On the tablet: note the **new** root password and Wi-Fi IP. Put them in `secrets/remarkable.local.env` (or let the installer ask).
2. Download ZIP (or use your existing clone), then:

   ```bash
   bash scripts/install.sh --start
   ```

3. Answer only what is missing. If secrets already have PIN + sync repo + token, it reuses them and pushes PIN + sync to the tablet.
4. Wait for sync (or open Lobby → Sync). Notes and encrypted files come from GitHub. Vault key material comes from `secret/vault` — open an encrypted note and enter your **vault PIN** (same as before wipe).

## What comes back automatically

- Plain notes (`.md`) and encrypted notes (`.md.enc`) from the GitHub repo
- Vault unlock material (`secret/vault`) so the old vault PIN still works
- Phone PIN length and sync settings, if saved in computer secrets and applied by `configure-sync.sh`

## What does **not** live in the notes repo

Re-set in Lobby Settings if you care: reading font, USB keyboard layout, rotation. Those stay in tablet `settings.json` and are not mirrored as note files.

## If the GitHub token is gone

Installer can open the prefilled token page; paste a new fine-grained token into secrets, then:

```bash
bash scripts/configure-sync.sh
```
