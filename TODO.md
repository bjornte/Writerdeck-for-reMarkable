# TODO

Writerdeck turns a first-gen reMarkable into a Wi-Fi Markdown typewriter. Most of the product is shipped — see [DONE.md](DONE.md).

How: [docs/architecture.md](docs/architecture.md). Why: [docs/decisions.md](docs/decisions.md). Gotchas: [docs/lessons.md](docs/lessons.md). Words: [docs/terms.md](docs/terms.md).

Verify on the tablet before checking anything off. Keys use a socket, not a fake keyboard device ([decisions.md](docs/decisions.md) §3).

## Open for you

1. Physical Home — please press the middle button once from edit, once from read, and once from the Lobby. Scripts cannot do this. [user-should-test.md](docs/user-should-test.md); [decisions.md](docs/decisions.md) §16.

## Settled (kept for pointers)

Editor fork, EditHelper, wrap/undo keep, QML assembly, and upstream ancestry are done. Policy: [decisions.md](docs/decisions.md) §4–§6. Keyboard harness green: **110/110/0** ([editor-testing/todo.md](docs/editor-testing/todo.md)).

## Open question

Stay on current firmware forever? Each OTA resets the SSH password and may wipe the systemd unit — redeploy and re-enable.

## Resume prompt

> Project: reMarkable 1 Wi-Fi Markdown typewriter. Server in `daemon/`; editor from Writerdeck-keywriter fork; notes in `Writerdeck-user-documents/`. Mac deploys; phone types.
> Next: owner Physical Home check ([docs/user-should-test.md](docs/user-should-test.md)). Harness **110/110/0** @ fork `0bb3b70`. Stop short of forking Qt TextEdit ([decisions.md](docs/decisions.md) §5–§6).
> Read: architecture, decisions, DONE, lessons, terms, integrity-audit. Device: `secrets/remarkable.local.env`.
> Constraints: no jailbreak / Toltec; keep OTA; static Go ARM binary.
