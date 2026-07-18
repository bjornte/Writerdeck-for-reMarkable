# scripts/

Deploy and test helpers. Run from the repo root. Credentials come from [secrets](../secrets/remarkable.local.env) via `_env.sh` (created/filled by `ensure-secrets.sh` when you install).

## What lands on the tablet

Writerdeck-server from `deploy-rmkbd.sh` (local Go build, or Release tag `server`). Writerdeck from CI via `deploy-keywriter.sh` (Release tag `keywriter`). Launcher, `wd`, and the systemd unit from this folder.

## First-time install

```bash
bash scripts/install.sh --start
```

Visitors: download the slim [installer ZIP](https://github.com/bjornte/Writerdeck-for-reMarkable/releases/download/installer/Writerdeck-installer.zip) (or rebuild with `bash scripts/pack-installer.sh`). Asks only for missing password / Wi-Fi / optional sync / PIN; reuses `secrets/remarkable.local.env`. Fetches binaries from Releases, deploys, health-checks, enables autostart, then `configure-sync.sh` if PIN/sync are saved. Checks only: `bash scripts/preflight.sh`. Remove from tablet (keeps Mac secrets + GitHub): `bash scripts/uninstall.sh` (`--reboot` optional). Reinstall notes: [reinstall-cheatsheet](../docs/install-onboarding/reinstall-cheatsheet.md). Regression: `bash scripts/test-install-reuse.sh`.

## Everyday commands

Ship a new editor binary after CI built it:

```bash
git push && bash scripts/fetch-keywriter-dist.sh && bash scripts/deploy-keywriter.sh -b
bash scripts/test-edit-session.sh
```

Ship a new server: `bash scripts/deploy-rmkbd.sh` (Go build if available, else Release), then restart the service.

Keyboard caret work: `bash scripts/test-keyboard-harness.sh --fast` ([editor-testing](../docs/editor-testing/todo.md)). Lobby/Home: `bash scripts/test-lobby-keyboard.sh`. Vault: `test-vault.sh` / `test-vault-e2e.sh`.

Show Lobby: `wd` or `bash scripts/lobby.sh`.

## Other useful scripts

ensure-secrets.sh — create/fill secrets (reuses saved values; optional sync + opens GitHub token URL). bootstrap.sh — SSH key and Wi-Fi SSH. preflight.sh — secrets / optional go / ping / dist. install.sh — first-time chain. uninstall.sh — remove Writerdeck from the tablet (`--reboot` optional). pack-installer.sh — slim end-user ZIP (`dist/Writerdeck-installer.zip`). configure-sync.sh — push SYNC_REPO + GH_TOKEN to a running tablet. fetch-keywriter-dist.sh / fetch-server-dist.sh — Release curl (gh fallback). install-service.sh — systemd unit (`--start` = start + health check + enable). install-alias.sh — Mac shortcuts. fix-hostkey.sh — after OTA host-key change. recover-orphaned-vault-notes.sh — after a vault key mistake.

Never log secrets. Prefer idempotent scripts.
