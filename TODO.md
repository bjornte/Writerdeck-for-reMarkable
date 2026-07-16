# TODO

Writerdeck for reMarkable 1 turns a first-gen tablet into a Wi-Fi Markdown typewriter. Phases 0–8, integrity slices 1–11, server-side GitHub sync, and encryption round 1 are shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md).

Keystrokes reach the editor over `/run/Writerdeck.sock`, not uinput ([docs/decisions.md](docs/decisions.md) §1). Verify on the device before checking anything off.

## Next unchecked

1. Physical Home — **done** (session `EVIOCGRAB` + fork `3be2de4` without `suppressNextHomeKey`; Writerdeck binary deployed). Please press physical Home once from edit, read, and Lobby to confirm. [docs/todo-handoff-physical-home-input.md](docs/todo-handoff-physical-home-input.md); [docs/decisions.md](docs/decisions.md) §28.
2. Keyboard editing — harness sign-off **105/105** @ `21-21-15`. Open: Shift+Up/Down from mid-sentence across wrapping paragraphs (not covered by current scenarios) — [docs/editor-testing/todo.md](docs/editor-testing/todo.md).
3. Keywriter fork migration — **done.** Owned fork [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) (`master`); helpers + C++ infra + Lobby/shell in-tree. Policy: [docs/decisions.md](docs/decisions.md) §3. Active rule: `.cursor/rules/writerdeck.mdc`.

## Open question

Stay firmware-update-current? Each OTA resets the SSH password and may wipe the systemd unit — re-deploy and re-enable ([docs/decisions.md](docs/decisions.md) open risks).

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Writerdeck-server (`daemon/` → `/home/root/Writerdeck-server`); Writerdeck-keywriter fork → Writerdeck (socket `/run/Writerdeck.sock`, notes in `Writerdeck-user-documents/`). Mac deploys; iPhone uses.
> Shipped: [DONE.md](DONE.md). Next: owner manual physical Home check; Shift+vertical mid-sentence across wrapping paragraphs. Keyboard harness **105/105**. Integrity: [integrity-audit.md](docs/integrity-audit.md).
> Read: architecture, decisions, DONE, lessons, browser-vs-tablet, integrity-audit. Device: `secrets/remarkable.local.env` (`RM_HOST_WIFI`).
> Constraints: no jailbreak/OTA/Toltec; `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
