# scripts/

Deploy and test helpers. Run from the repo root. Credentials come from [secrets](../secrets/remarkable.local.env) via `_env.sh` (created/filled by `ensure-secrets.sh` when you install).

## What lands on the tablet

Writerdeck-server from `deploy-rmkbd.sh` (local Go build, or Release tag `server`). Writerdeck from CI via `deploy-keywriter.sh` (Release tag `keywriter`). Launcher, `wd`, and the systemd unit from this folder.

## First-time install

```bash
bash scripts/install.sh --start
```

Asks for password + Wi-Fi IP if missing, fetches binaries from Releases, deploys, health-checks the phone page, enables autostart. Checks only: `bash scripts/preflight.sh`. Notes: [install-onboarding](../docs/install-onboarding/).

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

ensure-secrets.sh — create/fill secrets (prompts on a TTY). bootstrap.sh — SSH key and Wi-Fi SSH. preflight.sh — secrets / optional go / ping / dist. install.sh — first-time chain. fetch-keywriter-dist.sh / fetch-server-dist.sh — Release curl (gh fallback). install-service.sh — systemd unit (`--start` = start + health check + enable). install-alias.sh — Mac shortcuts. fix-hostkey.sh — after OTA host-key change. recover-orphaned-vault-notes.sh — after a vault key mistake.

Never log secrets. Prefer idempotent scripts.
