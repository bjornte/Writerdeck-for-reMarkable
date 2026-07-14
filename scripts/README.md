# scripts/

Dev/deploy automation. Run from repo root. Credentials: [../secrets/remarkable.local.env](../secrets/remarkable.local.env) via `_env.sh`.

## On-device names

| On tablet | Built from |
|---|---|
| `/home/root/Writerdeck-server` | `daemon/` via `deploy-rmkbd.sh` |
| `/home/root/Writerdeck` | CI via `deploy-keywriter.sh` |
| `/home/root/Writerdeck-launcher.sh` | `scripts/Writerdeck-launcher.sh` |
| `/home/root/wd` | `scripts/wd` |
| `writerdeck.service` | `scripts/writerdeck.service` |

`migrate-device-layout.sh` renames legacy paths on deploy.

## Writerdeck deploy

`deploy-keywriter.sh` pushes `dist/Writerdeck` — does not rebuild.

```bash
git push && bash scripts/fetch-keywriter-dist.sh && bash scripts/deploy-keywriter.sh -b
bash scripts/test-edit-session.sh
```

Local Docker (Apple Silicon: `--platform linux/amd64` on both docker commands).

## Scripts

| Script | Purpose |
|---|---|
| `_env.sh` | Shared secrets, ssh, `rm_send_file` |
| `paths.sh` | On-device path constants |
| `bootstrap.sh` | SSH key install, enable Wi-Fi SSH |
| `recon.sh` | Device fact snapshot → `docs/recon/` (re-run after OTA) |
| `fetch-keywriter-dist.sh` | Pull CI artifacts (`gh` required) |
| `deploy-keywriter.sh` | Ship Writerdeck + qt5 sysroot |
| `deploy-rmkbd.sh` | Cross-build and ship Writerdeck-server (flush-save + graceful wait) |
| `Writerdeck-launcher.sh` | Qt linuxfb launch env |
| `test-e2e.sh` | Full pipeline test (`-s` skips server rebuild) |
| `test-edit-session.sh` | Writerdeck/QML regression — POST `/api/open` |
| `test-keyboard-harness.sh` | Modifier+arrow and selection on device (WebSocket path). Scenarios: [editor-testing/scenario-cookbook.md](../docs/editor-testing/scenario-cookbook.md). `-s NAME`, `-m PREFIX`, `--list`, `--unit`, `--fast`, `--no-prepare`, `--hard-reset`, `-v`. |
| `test-lobby-keyboard.sh` | Lobby keys after return from edit; Home-from-read must not quit. `POST /api/lobby`, `POST /api/test/home`. After Lobby, `handleHome`, or `lobbyFocus` QML changes ([decisions.md](../docs/decisions.md) §29). |
| `test-vault.sh` | Loopback vault encrypt, lock, unlock, decrypt on device or `--local`. Resets vault for deterministic PIN. After `daemon/vault.go` or vault API edits. |
| `recover-orphaned-vault-notes.sh` | Re-wrap `.md.enc` notes after a vault key rotation. Needs `--old-vault-ref` from GitHub `secret/vault` history, `--notes`, `--pin`. |
| `test-vault-e2e.sh` | Tablet vault UI, keyboard PIN, Files encrypt/decrypt, GitHub `secret/` and note bytes. Needs sync on. Logs: `docs/recon/test-vault-e2e-*.txt`. After vault QML or E2E harness edits. |
| `lobby.sh` / `wd` | Show Lobby on e-ink |
| `install-service.sh` | Install systemd unit (manual enable) |
| `install-alias.sh` | Mac aliases: `rmpush`, `rmkw`, `wd` |
| `fix-hostkey.sh` | Clear stale SSH host key |
| `push.sh` | Stage, commit, push |

Convention: script device actions; tee output to `docs/recon/` where useful. Never log secrets.

Wi-Fi via `$RM_HOST` (default from secrets). Keep scripts idempotent.
