# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, and server-side GitHub sync are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. USB Norwegian keyboard — æ ø å Æ Ø Å, AltGr, `@`, `{` `}` on a physical NO keyboard. Qmaps and Lobby picker are shipped; Alt+Left/Right no longer flip to preview ([lessons.md](docs/lessons.md)). Remaining checks are hardware-only.
2. Lobby Ctrl-K on USB keyboard — device verify.
3. Power button sleep/wake — device verify. Implementation is in [DONE.md](DONE.md); test is outstanding.

## Phase 10 — locales and protection

Design notes: [docs/improvements.md](docs/improvements.md).

USB keyboard locales:

- [x] `no.qmap` and `us.qmap` via ckbcomp and kmap2qmap; ship in `keymaps/`.
- [x] `Writerdeck-launcher.sh` reads `keyboardLayout` from settings.
- [x] Hotplug-safe keyboard path; Lobby Keyboard tab picker.
- [ ] Device test (item 1 above).

`loadkeys` and `setxkbmap` do not work for Qt apps on rM — see [remarkable-keywriter#1](https://github.com/dps/remarkable-keywriter/issues/1).

Encrypted note subset:

- [ ] Design ADR: encrypted subfolder, passphrase-derived key, session unlock, sync exclusion.
- [ ] Implement after sign-off.

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); patched keywriter → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next unchecked: Norwegian USB device test (æøå, AltGr — Alt+arrow fixed in qmap); Ctrl-K USB verify; power button device test. Phase 10 encryption: [improvements.md](docs/improvements.md). Integrity: [integrity-audit.md](docs/integrity-audit.md). After QML edits: `bash scripts/test-edit-session.sh` ([decisions](docs/decisions.md) §21); after arrow/selection QML: `bash scripts/test-keyboard-harness.sh` (§22).
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
