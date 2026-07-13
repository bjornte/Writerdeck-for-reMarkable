# Encrypted notes (round 1 shipped)

See [decisions.md](decisions.md) §31 and [DONE.md](../DONE.md). Regression checklist only — not open work.

Separate vault PIN (tablet only), `.md.enc` beside plain `.md`, GitHub `secret/pin` and `secret/vault`, phone download waits for tablet PIN entry. No persistent unlocked state — session key only while editing one encrypted note; Lobby return clears it. Harness notes use the `z-test-` prefix ([decisions.md](decisions.md) §32).

Lobby Files (private notes on): second touch row — Encrypt / New encrypted on a plain note, Decrypt on `.md.enc`. Settings tab is enable and change PIN only (no lock/unlock). USB x / y on Files still work.

## Verify on device

```bash
bash scripts/test-vault.sh
bash scripts/test-vault-e2e.sh
```

`test-vault-e2e.sh` needs sync on with token and repo. It drives Settings PIN setup, Files encrypt/decrypt, PIN change via keyboard, edit encrypted note, and GitHub checks for `secret/pin`, opaque `.md.enc`, and decrypted `.md`.

Spot-check: Settings shows no lock/unlock; every encrypted Edit/Read asks PIN; no “unlocked” anywhere. Phone download of an encrypted note shows the tablet PIN overlay; after PIN entry, download succeeds without opening the note on e-ink.
