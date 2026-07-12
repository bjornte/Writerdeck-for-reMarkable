# Dev behind a corporate firewall (retired)

**Status:** Retired July 2026. The work laptop is no longer in use; all device work runs on the Mac on the same Wi-Fi as the tablet.

This file is kept for historical context only. Do not follow these steps.

---

## What this was

A **two-machine split**: corporate VPN kept the Windows work laptop off the home LAN, so the Mac (on Wi-Fi, reachable to the reMarkable) did deploy/recon/device scripts while the ThinkPad did cross-builds, commits, and code edits. **Git was the bridge** — push from one machine, pull on the other.

A reverse SSH tunnel was considered but deferred (infra + exposure). The intended clean fix was an IT split-tunnel exception.

## Machinery (removed)

| Piece | Role |
|---|---|
| `scripts/watch-mac.sh` | Mac: pull, auto-commit+push new `docs/recon/` logs |
| `scripts/watch-pc.ps1` | PC: pull loop + toast when Mac pushed |
| `scripts/build-rmkbd.ps1` | PC: cross-build Writerdeck-server only (no deploy) |
| `scripts/push.ps1` | PC: commit+push with personal git identity |
| CI commit of `dist/Writerdeck` | Optional git path for keywriter binary (Mac also uses `fetch-keywriter-dist.sh` from Actions) |

## Current setup

Single Mac on LAN: `deploy-rmkbd.sh`, `deploy-keywriter.sh`, `fetch-keywriter-dist.sh`, `test-edit-session.sh`. Secrets in `secrets/remarkable.local.env`. See [architecture.md](architecture.md) and [decisions.md](decisions.md) #10 (retired).

Recorded in [decisions.md](decisions.md) #10.
