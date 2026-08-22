# rM2 port — session handoff

Copy this file (or point a fresh agent at it) to continue Writerdeck on reMarkable 2.

Branch: `spike/rm2-official-sdk`. Substantial uncommitted work may exist — do not commit unless the owner asks.

Related: [README.md](README.md), [docs/decisions.md](../../docs/decisions.md) §33, `.cursor/rules/writerdeck.mdc`.

## Goal

Port Writerdeck to reMarkable 2 using reMarkable's official Codex SDK (Qt6 + epaper QPA). No Toltec, no rm2fb — preserve over-the-air updates.

Work autonomously: build, deploy, verify with scripts. Deploy success is not tested. PNG verify beats eyeballing.

## Device and secrets

| Fact | Value |
|------|--------|
| Firmware | 3.27.3.0 (Codex 5.7.126) — required for official epaper path |
| USB host | `10.11.99.1` |
| Wi-Fi host | `192.168.1.115` (was `.114`; check env) |
| Secrets file | `secrets/remarkable.local.env` — `RM2_HOST_USB`, `RM2_HOST_WIFI`, `RM2_ROOT_PASSWORD` |

Host pick tries USB then Wi-Fi (`scripts/_env.sh` → `rm2_pick_host()`, `spike/rm2-official-sdk/scripts/rm2-pick-host.sh`).

Run spike scripts with `bash`, not bare zsh — word-split breaks `$RM_SSH_OPTS`.

Mac PNG verify needs `pip3 install Pillow`.

After firmware OTA: `ssh-keygen -R` for changed host keys; password may change in env.

## What passed (Aug 2026)

All four spike phases verified on device with autonomous PNG checks (`verify-png.sh`: 1404×1872, minimum dark pixel count).

| Phase | Proves | Command |
|-------|--------|---------|
| 1 hello | Official epaper QPA draws | `bash spike/rm2-official-sdk/scripts/verify-spike.sh hello` |
| 2 textedit | Wrapped prose on epaper | `bash spike/rm2-official-sdk/scripts/verify-spike.sh textedit` |
| 3 socket | NDJSON keys → on-screen text | `bash spike/rm2-official-sdk/scripts/verify-spike.sh socket` |
| 4 fork-probe | Full Writerdeck-keywriter binary | `bash spike/rm2-official-sdk/scripts/verify-spike.sh fork-probe` |

Phase 4 screenshot (`screenshots/rm2-spike-2026-08-22-fork-probe-verify.png`) shows the real Lobby UI — Documents tab, New/Edit/Read/Rename/Delete/Download.

Quick re-check of phases 1–3:

```bash
bash spike/rm2-official-sdk/scripts/verify-all-spike.sh
```

## Proven launch stack

```bash
export QT_QUICK_BACKEND=epaper
export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS=rotate=180:invertx
# binary args: -platform epaper
```

Production launcher stub (deploy when fork ships): `scripts/Writerdeck-launcher-rm2.sh`.

Screenshot capture: `/dev/fb0` is not a readable mirror on rM2 software epaper. Apps write PNG via `QQuickWindow::grabWindow()`; pull with `capture-screenshot.sh` or `run-and-capture*.sh`.

## Infrastructure in this repo

| Path | Role |
|------|------|
| `spike/rm2-official-sdk/scripts/_spike-sdk.sh` | SDK symlink fix, Docker/cmake build helpers |
| `spike/rm2-official-sdk/scripts/build-*.sh` | hello, textedit, socket, fork-probe builds |
| `spike/rm2-official-sdk/scripts/verify-spike.sh` | Deploy + run + PNG verify per phase |
| `spike/rm2-official-sdk/scripts/verify-png.sh` | Size and dark-pixel thresholds |
| `spike/rm2-official-sdk/fork_probe/` | CMakeLists reference for fork Qt6 build |
| `spike/rm2-official-sdk/.cache/writerdeck-keywriter` | Fork checkout (gitignored) |
| `spike/rm2-official-sdk/.cache/out/` | Build artifacts (gitignored) |

Mac builds use Docker linux/amd64. SDK installs in container `/tmp/codex-sdk-rm2` — host-mounted `.cache/sdk-*` had permission and tar corruption issues.

Cross-compile: SDK `OEToolchainConfig.cmake` + SDK `cmake` + Ninja, not host gmake.

## Fork probe patches (build time only — not in fork yet)

`build-fork-probe.sh` clones `https://github.com/bjornte/Writerdeck-keywriter.git` and applies:

1. Remove `Q_IMPORT_PLUGIN(QsgEpaperPlugin)` and `#include <QtPlugin>` from `main.cpp` (Qt5 epaper plugin import).
2. Split header-only `EditUtils` into `edit_utils.cpp` + trimmed header (MOC/link fix for Qt6 CMake).
3. Insert `QTimer` + `grabWindow()` screenshot hook → `/home/root/spike-rm2-fork-probe/screen.png`.

The fork already sets `QMLSCENE_DEVICE=epaper` and `QT_QPA_PLATFORM=epaper:enable_fonts` on `__arm__` when unset; production rM2 also uses `QT_QUICK_BACKEND=epaper` from the launcher.

## What remains

### 1. Formalize Qt6 build in Writerdeck-keywriter (highest priority)

Move probe patches into the fork permanently:

- CMake/Qt6 Quick build replacing qmake/Toltec path.
- Permanent `edit_utils.cpp`.
- Drop Qt5 epaper plugin import.
- Official SDK CI for rm2 (and eventually rm1).

Use `fork_probe/CMakeLists.txt` and `build-fork-probe.sh` as reference.

### 2. Production deploy path (this repo)

- `deploy-keywriter-rm2.sh` or extend existing deploy.
- Wire `Writerdeck-launcher-rm2.sh`.
- Go Writerdeck-server via `deploy-rmkbd.sh` with `RM2_HOST` — likely already works on rM2.

### 3. End-to-end on rM2

- Phone WebSocket → daemon → `/run/Writerdeck.sock` → real editor (not spike `socket_spike`).
- SSH `journalctl -u writerdeck` — fail on QML parse errors or instant editor exit.
- `bash scripts/test-edit-session.sh` on rM2.

### 4. Editor validation on Qt6

- Typing harness / EditHelper on rM2.
- Update `docs/architecture.md` device facts (still rM1-centric in places).

### 5. CI and verify-all

- Add fork-probe to `.github/workflows/spike-rm2-hello.yml` if feasible (slow build).
- Optionally add fork-probe to `verify-all-spike.sh`.

## Constraints

- No jailbreak / Toltec.
- Document integrity is paramount — see `docs/decisions.md`, `docs/integrity-audit.md`.
- Editor behavior lives in fork C++ (`EditHelper`), not `build-keywriter.sh`.
- Static Go server: `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
- Cross-compile Writerdeck on Mac via Docker/GitHub Actions — not local docker for production binary.
- Do not `pkill -f /home/root/Writerdeck` (matches Writerdeck-server).

## Suggested start for a fresh session

1. Read this file and [README.md](README.md).
2. Confirm device: `bash spike/rm2-official-sdk/scripts/recon-rm2.sh`.
3. Baseline: `bash spike/rm2-official-sdk/scripts/verify-spike.sh fork-probe`.
4. Begin fork Qt6 CMake port, then deploy pipeline, then edit-session + server integration.

Verify before calling work done. Report tersely to the owner unless they ask for depth.
