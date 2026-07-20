# TODO

Writerdeck turns a first-gen reMarkable into a Markdown typewriter with USB and Bluetooth keyboards. Most of the product is finished — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md). Words: [docs/terms.md](docs/terms.md).

Verify on the tablet before checking anything off. Keys use a socket, not uinput (a fake keyboard device) ([decisions.md](docs/decisions.md) §3).

## Open for you

1. Physical Home — please press the middle button once from edit, once from read, and once from the Lobby. Scripts cannot do this. [user-should-test.md](docs/user-should-test.md); [decisions.md](docs/decisions.md) §16.

2. Tell people about Writerdeck — best places to post (reMarkable 1, no Toltec / OTA-safe angle):
   - **Reddit** — [r/reMarkableTablet](https://www.reddit.com/r/reMarkableTablet) (main user hub; reMarkable’s site points here)
   - **reMarkable Community Discord** — invite via [awesome-reMarkable](https://github.com/reHackable/awesome-reMarkable) / [remarkable.guide](https://remarkable.guide/) (technical + power users; say clearly this is SSH + systemd, not Toltec)
   - **Facebook** — official reMarkable user group (from [Join the community](https://remarkable.com/join-the-community))
   - **MobileRead** — [More E-Book Readers](https://www.mobileread.com/forums/forumdisplay.php?f=140) (e-ink audience; quieter than Reddit)
   - Optional launch posts: **Hacker News** (Show HN), **Lobsters** — one-shot, not ongoing forums

   Listed on [awesome-reMarkable](https://github.com/reHackable/awesome-reMarkable) Applications ([PR #268](https://github.com/reHackable/awesome-reMarkable/pull/268)).

## Open

- [ ] Windows installer (native or clearly supported path). Mac/Linux stays bash; Windows is missing today — [install-onboarding/todo-install-onboarding.md](docs/install-onboarding/todo-install-onboarding.md).
- [x] Lobby shortcuts on disk (tabs, Home path, sync/edit Enter, rotate; remove old Ctrl-K picker) — [todo-lobby-ui-shortcuts.md](docs/todo-lobby-ui-shortcuts.md).
- [ ] Lobby chrome still hardwired (labels, copy, fills, radii, type sizes) — [todo-lobby-ui-chrome.md](docs/todo-lobby-ui-chrome.md).

## Settled (kept for pointers)

Editor fork, EditHelper, wrap/undo keep, QML assembly, and linking git history to Singleton’s original are done. Policy: [decisions.md](docs/decisions.md) §4–§6. Automated typing tests: all 112 passed ([editor-testing/todo.md](docs/editor-testing/todo.md)). In-editor copy/cut/paste over Bluetooth: fork `df1d38b`. Mac/Linux installer credential memory and sync push: done (Windows still open). No-keyboard Lobby tip with phone QR: fork `80f568b`; phone path needs WebSocket `hello`, excluding Cursor/Electron ([decisions.md](docs/decisions.md) §34). Phone keyboard-first (no notes list); Lobby Download offers to open phones: fork `3cfff08`. Lobby Keyboard tab boxes + live `(connected)` status: fork `55da42b`. Lobby on-disk UI config (`lobby-ui.json`): fork `21ed25a` ([decisions.md](docs/decisions.md) §36).

## Open question

Stay on current firmware forever? Each OTA (over-the-air update) resets the SSH password and may wipe the boot service — redeploy and re-enable.

Install now auto-enables `writerdeck` on boot after a health check — inspect bricking risk: [install-onboarding/todo-install-onboarding.md](docs/install-onboarding/todo-install-onboarding.md) (Follow-up).

## Resume prompt

> Project: reMarkable 1 Markdown typewriter (USB and Bluetooth keyboards). Server in `daemon/`; editor from Writerdeck-keywriter fork; notes in `Writerdeck-user-documents/`. Mac deploys; phone types.
> Next: owner Physical Home check ([docs/user-should-test.md](docs/user-should-test.md)). Phone is keyboard-first (paste, sync token, Download here?); Lobby Download at fork `3cfff08`. Lobby Keyboard boxes + live `(connected)` at fork `55da42b`. Lobby UI config on disk (`lobby-ui.json`) at fork `21ed25a` ([docs/decisions.md](docs/decisions.md) §36). Typing tests all 112 passed at fork commit `df1d38b` (includes Ctrl+C/X/V). No-keyboard Lobby tip + QR; Cursor agent tabs do not count as a phone keyboard ([decisions.md](docs/decisions.md) §34). Stop short of replacing Qt’s text box ([decisions.md](docs/decisions.md) §5–§6).
> Read: architecture, decisions, DONE, lessons, terms, integrity-audit, browser-vs-tablet. Device: `secrets/remarkable.local.env`.
> Constraints: no jailbreak / Toltec; keep OTA (over-the-air updates); one static Go ARM binary.
