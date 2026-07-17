# TODO: Install and onboarding

Make it easier for visitors to get Writerdeck onto a reMarkable 1 without already knowing SSH, systemd, and CI artifacts. Identified in a docs/scripts review (2026-07).

Folder overview: [README.md](README.md). Cross-links: [README.md](../../README.md) install section, [secrets/README.md](../../secrets/README.md), [scripts/README.md](../../scripts/README.md), [architecture.md](../architecture.md) device facts.

## Problem

The install path worked for the repo owner but assumed expertise: six ordered scripts, a gitignored secrets file, CI artifacts that used to need `gh`, and manual systemd commands with boot-loop risk. A fresh clone still has no editor binary until fetch/install runs.

## Checklist

### README and secrets (low effort) — done 2026-07-18

- [x] Add a **Before you start** block: reMarkable 1 only, Mac or Linux on same Wi‑Fi, Go 1.25+ (see `daemon/go.mod`), tablet awake; editor binaries via Release (not mandatory `gh`).
- [x] Step 1 must include **Wi‑Fi IP** (`RM_HOST_WIFI`) — not just the root password. Point to tablet Wi‑Fi settings and router DHCP reservation ([architecture.md](../architecture.md)).
- [x] Fix `remarkable.local.env.example` copy instruction: `cp` for bash, not PowerShell `Copy-Item` as the primary path.
- [x] Add **You're done when**: Lobby on e‑ink, phone loads `http://<ip>:8000/` with a populated note list and connection bar not stuck on `connecting...`.
- [x] Add **After a firmware update**: password changes, systemd unit may be gone — re-run deploy + `install-service.sh`, re-enable only after manual `systemctl start` passes ([TODO.md](../../TODO.md) open question; [architecture.md](../architecture.md)).
- [x] Add **Recovery** one-liner: `systemctl disable --now writerdeck && systemctl start xochitl` (already in `install-service.sh` output; belongs in README).

### Binaries without `gh` (medium effort) — done 2026-07-18

- [x] Publish CI-built `Writerdeck` + `qt5.tar.gz` on **GitHub Releases** (rolling tag `keywriter` from `build-keywriter.yml` on main) so visitors can `curl` without `gh auth login`.
- [x] Document browser download from Actions artifacts as a fallback when Releases/`gh` unavailable.
- [x] Update `fetch-keywriter-dist.sh` to try release URL first, then `gh run download`, with a clear error when dist is empty.

### Scripts (medium effort) — done 2026-07-18

- [x] **`scripts/preflight.sh`** — check secrets file exists and has password + Wi‑Fi IP, ping tablet, `go version`, dist artifacts present (or `--fetch`).
- [x] **`scripts/install.sh`** — idempotent chain: preflight → `bootstrap.sh` → fetch if needed → `deploy-keywriter.sh` → `deploy-rmkbd.sh` → `install-service.sh`; print next manual steps (`systemctl start` / `enable`) and success URL. Supports `--start`.
- [x] Optional: `install-service.sh --start` runs `systemctl start writerdeck` over SSH and prints `journalctl` tail so visitors skip opening a second SSH session for the smoke test.

### Verification for visitors (low effort) — done 2026-07-18

- [x] Point install docs at `bash scripts/test-edit-session.sh` as an optional automated smoke test after first deploy (not required for sign-off by repo rules, but useful for self-service).
- [x] Mention `bash scripts/fix-hostkey.sh` in install section when OTA changes host key (bootstrap already references it).

## Out of scope (for now)

- Windows-native install path (scripts are bash; README should say Mac/Linux explicitly rather than imply Windows).
- reMarkable 2 support.
- Toltec or jailbreak-based distribution.
- Phone-only install with no computer.

## Done when

A new contributor can follow README only: configure secrets with IP + password, run one install command (or a short numbered list), load the phone UI, see the Lobby on the tablet, and recover from a bad autostart without reading `docs/lessons.md`.
