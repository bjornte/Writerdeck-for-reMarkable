# docs/

Reference material and generated artifacts.

## Contents
- `architecture.md` — how the system works: architecture, environment facts, the dev/deploy loop.
- `decisions.md` — the Architecture Decision Record (ADR): the *why* behind every choice.
- `recon/` — regenerable device output, self-logged via `tee` by `scripts/*`. Timestamped `recon-*.txt` logs are committed, then pruned once their lesson lands in TODO/DONE — the folder persists via `.gitkeep`. Large `*.tar.gz` backups are gitignored.

## External references
- Editor (keywriter): https://github.com/dps/remarkable-keywriter
- reMarkable input subsystem (evdev): https://remarkable.guide/devel/device/input.html
- Awesome-reMarkable index: https://github.com/reHackable/awesome-reMarkable
- ~~uinput Go lib~~ (not usable — this kernel can't load uinput; the editor is fed over a socket instead): https://github.com/bendahl/uinput
- libremarkable (fallback self-editor spike): https://github.com/canselcik/libremarkable
