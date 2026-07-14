# Document integrity

Last updated 2026-07-13. Contract: [architecture.md](architecture.md), [decisions.md](decisions.md). Shipped features: [../DONE.md](../DONE.md).

Writerdeck is a typewriter. Integrity means your words end up as real Markdown on the tablet, stay readable, and are not silently replaced, emptied, or forked without you noticing.

For normal solo use — one person, one note at a time, saving with Home or waiting for autosave — text is in good shape. Saves are plain UTF-8. Sync does not mass-delete. When sync disagrees, it tends to keep both copies rather than pick a loser.

That bar is not the same as bank-grade durability or real-time collaborative editing, and we are not trying to be either.

What is still open: sudden death — pull power or kill the process hard and you can lose up to about forty-five seconds since the last autosave. The note on screen is protected from sync overwriting it; other notes still sync. If GitHub or sync changes a file while you have it open, the phone shows a drift banner; the tablet does not auto-refresh, and if you keep typing and save, your buffer can win over newer disk — you must choose reload or keep going. Clashes often leave `note (tablet copy).md` beside `note.md`; you reconcile by hand. Rapid note switching can occasionally interleave saves and loads. With PIN turned off, anyone on your LAN can read or change notes. Firmware updates and paths we have not hammered in testing could still surprise us.

## What shipped

Edit lease: reconcile skips the note the tablet is editing. Plain-markdown save contract with HTML guard and toggleMode fix. Tablet rename/delete of the open file notifies the editor. Reconcile gated on `openNote` in status. ETag and If-Match on PUT. Atomic server writes via temp and rename. Tablet CRUD pairs to GitHub on the server. Disk drift WebSocket, phone banner, and reload endpoint. Forty-five-second autosave while editing. Tablet saves via loopback PUT. Save before deploy and stop via flush-save and graceful shutdown.

Also under contract: empty-push guard, Lobby Home wipe fix, server-side sync engine. Vault disable refuses while user `.md.enc` notes exist; sync will not apply a foreign `secret/vault` that would orphan them. Failed encrypted-note load shows a Files-tab error instead of a blank editor.

## Residual risks

SIGKILL, segfault, or power loss between autosave ticks. Flush-save fails if the editor socket is down, autosavenow ack times out, or Writerdeck binary is stale — loopback save and autosavenow need matching binaries on both sides.

`doLoad` is async; rapid Ctrl-K or Lobby Files switching can interleave saves and loads. `/api/open` may continue on save-ack timeout.

Clash handling keeps both copies but does not auto-reload the editor. Stale open-note tracking after phone-back can skip the wrong file in reconcile until status refreshes. External GitHub edits can still duplicate or resurrect notes.

PIN none on an untrusted LAN is a confidentiality and tampering risk, not just sync.

`test-edit-session.sh` guards "editor stays up on Edit" — not save under load, clash while typing, or power sleep with an open note. `scripts/test-vault.sh` covers loopback vault setup, per-note encrypt, and decrypt. `scripts/test-vault-e2e.sh` covers tablet PIN UI, Files encrypt/decrypt, PIN change, edit encrypted note, and GitHub sync (needs sync configured).

Encrypted `.md.enc` files are opaque on disk; plaintext exists in the editor buffer only while a vault session key is held. The UTF-8 Markdown contract applies to `.md` only. Returning to the Lobby clears open-note tracking for sync so reconciles are not stuck skipping the last edited file. Orphaned ciphertext after a vault key rotation (disable+setup, or sync applying the wrong `secret/vault`) is recoverable via `scripts/recover-orphaned-vault-notes.sh` and GitHub `secret/vault` history — not by the current PIN alone.

## Unknown

Firmware OTA side effects. Qt regressions beyond test-edit-session. Sleep-wake and rapid-switch edge cases. Novel failure modes we have not stress-tested.

Before shipping note, save, sync, or lifecycle changes: can this lose text, write wrong bytes, or overwrite without the user knowing? If yes, fix it or log an explicit acceptance here.

We borrow edit-lease thinking from local-first systems, optimistic concurrency from HTTP ETags, conflict copies in the Syncthing spirit, and atomic writes via temp-rename. CRDT or operational transform is not planned.
