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
   - **awesome-reMarkable** — open a PR to list the project under Applications (discovery for hackers)
   - **MobileRead** — [More E-Book Readers](https://www.mobileread.com/forums/forumdisplay.php?f=140) (e-ink audience; quieter than Reddit)
   - Optional launch posts: **Hacker News** (Show HN), **Lobsters** — one-shot, not ongoing forums

## Open

- [ ] Before actions that need a keyboard (edit note, new note, rename, and any similar Lobby prompts), check whether a Bluetooth (phone) or USB keyboard is present. If neither is available, show a short tip on how to connect either, including a QR code for the current phone-interface URL. Wishlist pointer: [docs/improvements.md](docs/improvements.md).
- [ ] Windows installer (native or clearly supported path). Mac/Linux stays bash; Windows is missing today — [install-onboarding/todo-install-onboarding.md](docs/install-onboarding/todo-install-onboarding.md).
- [x] Installer credential memory (Mac/Linux): `ensure-secrets.sh` reuses password / Wi-Fi / optional `SYNC_REPO` + `GH_TOKEN`; opens prefilled GitHub token page when needed; `configure-sync.sh` pushes sync to the tablet after start. Windows still open. Tablet still must not store the GitHub token on disk.

## Settled (kept for pointers)

Editor fork, EditHelper, wrap/undo keep, QML assembly, and linking git history to Singleton’s original are done. Policy: [decisions.md](docs/decisions.md) §4–§6. Automated typing tests: all 110 passed ([editor-testing/todo.md](docs/editor-testing/todo.md)).

## Open question

Stay on current firmware forever? Each OTA (over-the-air update) resets the SSH password and may wipe the boot service — redeploy and re-enable.

Install now auto-enables `writerdeck` on boot after a health check — inspect bricking risk: [install-onboarding/todo-install-onboarding.md](docs/install-onboarding/todo-install-onboarding.md) (Follow-up).

## Resume prompt

> Project: reMarkable 1 Markdown typewriter (USB and Bluetooth keyboards). Server in `daemon/`; editor from Writerdeck-keywriter fork; notes in `Writerdeck-user-documents/`. Mac deploys; phone types.
> Next: owner Physical Home check ([docs/user-should-test.md](docs/user-should-test.md)). Typing tests all 110 passed at fork commit `0bb3b70`. Stop short of replacing Qt’s text box ([decisions.md](docs/decisions.md) §5–§6).
> Read: architecture, decisions, DONE, lessons, terms, integrity-audit. Device: `secrets/remarkable.local.env`.
> Constraints: no jailbreak / Toltec; keep OTA (over-the-air updates); one static Go ARM binary.
