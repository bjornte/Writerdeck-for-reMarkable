# scripts/

Deploy and test helpers. Run from the repo root. Credentials come from [secrets](../secrets/remarkable.local.env) via `_env.sh`.

## What lands on the tablet

Writerdeck-server from `deploy-rmkbd.sh`. Writerdeck from CI via `deploy-keywriter.sh`. Launcher, `wd`, and the systemd unit from this folder.

## Everyday commands

Ship a new editor binary after CI built it:

```bash
git push && bash scripts/fetch-keywriter-dist.sh && bash scripts/deploy-keywriter.sh -b
bash scripts/test-edit-session.sh
```

Ship a new server: `bash scripts/deploy-rmkbd.sh`, then start the service.

Keyboard caret work: `bash scripts/test-keyboard-harness.sh --fast` ([editor-testing](../docs/editor-testing/todo.md)). Lobby/Home: `bash scripts/test-lobby-keyboard.sh`. Vault: `test-vault.sh` / `test-vault-e2e.sh` (full vault path including the tablet UI).

Show Lobby: `wd` or `bash scripts/lobby.sh`.

## Other useful scripts

bootstrap.sh — SSH key and Wi-Fi SSH. recon.sh — device snapshot after OTA. install-service.sh — systemd unit (enable only after a manual start works). install-alias.sh — Mac shortcuts. recover-orphaned-vault-notes.sh — after a vault key mistake.

Never log secrets. Prefer idempotent scripts.
