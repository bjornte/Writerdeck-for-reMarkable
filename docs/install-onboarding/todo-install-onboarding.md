# TODO: Install and onboarding

Make it easier for visitors to get Writerdeck onto a reMarkable 1 without already knowing SSH, systemd, and CI artifacts. Identified in a docs/scripts review (2026-07).

Folder overview: [README.md](README.md). Cross-links: [README.md](../../README.md) install section, [secrets/README.md](../../secrets/README.md), [scripts/README.md](../../scripts/README.md), [architecture.md](../architecture.md) device facts.

## Problem (was)

The install path worked for the repo owner but assumed expertise: six ordered scripts, a hand-edited secrets file, CI artifacts that needed `gh`, Go on the Mac to build the server, and a manual `systemctl enable` with boot-loop risk.

## Checklist

### README and secrets — done 2026-07-18

- [x] Short Install in README: reMarkable 1, Mac/Linux same Wi-Fi, tablet awake; Download ZIP; one command.
- [x] Prompt for root password + Wi-Fi IP (`ensure-secrets.sh`) instead of requiring a hand-edited secrets file first.
- [x] **You're done when**: Lobby on e-ink; phone list populated; connection bar not stuck on `connecting...`.
- [x] Recovery one-liner and OTA note in README.

### Binaries without `gh` / without Go — done 2026-07-18

- [x] Rolling Release `keywriter`: `Writerdeck` + `qt5.tar.gz` (`build-keywriter.yml` on main).
- [x] Rolling Release `server`: `Writerdeck-server` (`build-server.yml` on main).
- [x] `fetch-keywriter-dist.sh` / `fetch-server-dist.sh`: curl Release first, then `gh run download`, clear error + browser fallback.
- [x] `deploy-rmkbd.sh` builds with Go when present, else fetches the `server` Release.

### Scripts — done 2026-07-18

- [x] `ensure-secrets.sh`, `preflight.sh`, `install.sh` (chain + `--start`).
- [x] `install-service.sh --start`: start, health-check phone HTTP, then `systemctl enable`.

### Verification for visitors — done 2026-07-18

- [x] Optional `test-edit-session.sh` mentioned for self-service.
- [x] End-to-end install loop verified on a live tablet (prompted secrets, Release fetches, enable).

## Out of scope (for now)

- reMarkable 2 support — not started; open if the community wants it ([../decisions.md](../decisions.md) §33, [../improvements.md](../improvements.md)).
- Toltec or jailbreak-based distribution.
- Phone-only install with no computer.

## Follow-up

- [ ] Inspect dangers of bricking on boot now that `--start` auto-enables `writerdeck` after a health check. Confirm the check is strong enough, document failure modes, and whether enable should stay gated (human glance at e-ink, stricter checks, or opt-in flag).
- [ ] Windows installer — see [TODO.md](../../TODO.md) Open.
- [x] Mac/Linux installer remembers password / Wi-Fi / sync repo + token and can open GitHub’s token page — see [TODO.md](../../TODO.md) Open (Windows still open).

## Done when

A visitor can follow README only: Download ZIP, run one install command, answer two prompts, load the phone UI, see the Lobby, and recover from a bad autostart without reading `docs/lessons.md`. (Visitor path met; boot-risk follow-up still open.)
