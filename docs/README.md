# docs/

Reference material and generated artifacts.

## Contents
- `architecture.md` — how the system works; **§ Document integrity** is the non-negotiable product contract every feature must satisfy.
- `decisions.md` — ADR: the *why* behind every choice; **§ Document integrity** states the same contract as foundational policy.
- `lessons.md` — operational gotchas (deploy staleness, QML patch regressions, journald logs, sync footguns).
- `improvements.md` — owner wish-list plus **document integrity risk matrix & fix backlog** (implementation under the architecture/decisions contract).
- `recon/` — regenerable device output, self-logged via `tee` by `scripts/*`. Timestamped `recon-*.txt` logs are committed, then pruned once their lesson lands in TODO/DONE — the folder persists via `.gitkeep`. Large `*.tar.gz` backups are gitignored. Includes `test-edit-session-*.txt` from the Edit-from-browser regression script.

## External references
- Editor (upstream keywriter → on-device `Writerdeck`): https://github.com/dps/remarkable-keywriter
- reMarkable input subsystem (evdev): https://remarkable.guide/devel/device/input.html
- Awesome-reMarkable index: https://github.com/reHackable/awesome-reMarkable
- ~~uinput Go lib~~ (not usable — this kernel can't load uinput; the editor is fed over a socket instead): https://github.com/bendahl/uinput
- libremarkable (fallback self-editor spike): https://github.com/canselcik/libremarkable
