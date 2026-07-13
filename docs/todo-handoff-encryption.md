# Encrypted notes (shipped)

Implemented — see [decisions.md](decisions.md) §31 and [DONE.md](../DONE.md). This file is a short verify checklist for regressions.

Separate vault PIN (tablet only), `.md.enc` beside plain `.md`, GitHub `secret/pin` and `secret/vault`, phone download waits for tablet unlock. Vault locks on Lobby entry (`home` / `showlobby` ack), not on every editor-state poll while already in the Lobby. Harness notes use the `z-test-` prefix ([decisions.md](decisions.md) §32).

## Verify on device

```bash
bash scripts/test-vault.sh
bash scripts/test-vault-e2e.sh
```

`test-vault-e2e.sh` needs sync on with token and repo. It drives Settings PIN setup, Files encrypt/decrypt, PIN change via keyboard, edit encrypted note, and GitHub checks for `secret/pin`, opaque `.md.enc`, and decrypted `.md`.

Spot-check: phone download of an encrypted note while vault locked shows the tablet PIN overlay; after unlock, download succeeds without opening the note on e-ink.
