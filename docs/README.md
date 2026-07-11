# docs/

Reference material and generated artifacts.

## Contents
- `architecture.md` — how the system works: architecture, environment facts, the dev/deploy loop.
- `decisions.md` — the Architecture Decision Record (ADR): the *why* behind every choice.
- `lessons.md` — operational gotchas (deploy staleness, QML patch regressions, journald logs, sync footguns).
- `improvements.md` — owner wish-list and design notes (browser vs tablet parity, Lobby subpages, USB locales, encrypted notes); actionable items land in [../TODO.md](../TODO.md).
- `recon/` — regenerable device output, self-logged via `tee` by `scripts/*`. Timestamped `recon-*.txt` logs are committed, then pruned once their lesson lands in TODO/DONE — the folder persists via `.gitkeep`. Large `*.tar.gz` backups are gitignored. Includes `test-edit-session-*.txt` from the Edit-from-browser regression script.

## External references
- Editor (upstream keywriter → on-device `Writerdeck`): https://github.com/dps/remarkable-keywriter
- reMarkable input subsystem (evdev): https://remarkable.guide/devel/device/input.html
- Awesome-reMarkable index: https://github.com/reHackable/awesome-reMarkable
- ~~uinput Go lib~~ (not usable — this kernel can't load uinput; the editor is fed over a socket instead): https://github.com/bendahl/uinput
- libremarkable (fallback self-editor spike): https://github.com/canselcik/libremarkable
