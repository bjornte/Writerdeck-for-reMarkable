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
| Wi-Fi host | `192.168.1.117` — refresh: `bash scripts/discover-rm2-wifi.sh --write-secrets` |
| Secrets file | `secrets/remarkable.local.env` — `RM2_HOST_USB`, `RM2_HOST_WIFI`, `RM2_ROOT_PASSWORD` |

**Find Wi-Fi IP:** USB query first (best), LAN scan fallback:

```bash
bash scripts/discover-rm2-wifi.sh              # print IP
bash scripts/discover-rm2-wifi.sh --write-secrets # also update secrets
```

Deploy SSH still uses USB (`10.11.99.1`) — rM2 often refuses SSH on Wi-Fi. Phone UI uses Wi-Fi IP (`http://192.168.1.117:8000/`).

Host pick for deploy: `rm2-pick-host.sh` (USB then Wi-Fi SSH).

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

**Production stack on rM2 (Aug 22):** `deploy-keywriter-rm2.sh` + `deploy-rmkbd.sh` + `install-service.sh --start` on USB `10.11.99.1`. `test-edit-session.sh` PASS (Writerdeck stays up 8s, xochitl down, editorActive=true). Phone UI at `:8000/` loads JS (PIN screen). Wi-Fi IP was `.115` (timeout); device reported `.117` on LAN — update `RM2_HOST_WIFI` in secrets when Wi-Fi is the primary path.

Quick re-check of phases 1–3:

```bash
bash spike/rm2-official-sdk/scripts/verify-all-spike.sh
```

Production deploy (editor + server):

```bash
RM_HOST=10.11.99.1 bash scripts/deploy-keywriter-rm2.sh
RM_HOST=10.11.99.1 bash scripts/deploy-rmkbd.sh --deploy-only
RM_HOST=10.11.99.1 bash scripts/install-service.sh --start
RM_HOST=10.11.99.1 bash scripts/test-edit-session.sh
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
| `spike/rm2-official-sdk/fork_probe/` | Spike wrapper CMake + `CMakeLists.fork.txt` template |
| `scripts/discover-rm2-wifi.sh` | USB Wi-Fi IP query + LAN scan fallback |
| `scripts/deploy-keywriter-rm2.sh` | Production Qt6 editor deploy to rM2 |
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

### 1. Formalize Qt6 build in Writerdeck-keywriter (in progress)

Fork-side CMakeLists.txt + `edit_utils.cpp` drafted in spike cache; `build-fork-probe.sh` applies patches idempotently until pushed to `Writerdeck-keywriter` master. Push fork commit with:

- `CMakeLists.txt` (Qt6 Quick)
- `edit_utils.cpp` + trimmed `edit_utils.h`
- `#ifdef WRITERDECK_RM2_QT6_PROBE` screenshot hook in `main.cpp`
- Drop Qt5 epaper plugin import if still present on old branches

Reference: `spike/rm2-official-sdk/fork_probe/CMakeLists.fork.txt`, `build-fork-probe.sh`.

### 2. Production deploy path (this repo) — done for spike binary

- `scripts/deploy-keywriter-rm2.sh` — deploy Qt6 binary + `Writerdeck-launcher-rm2.sh` as `/home/root/Writerdeck-launcher.sh`
- `scripts/Writerdeck-launcher-rm2.sh` — epaper QPA launcher
- Go Writerdeck-server via `deploy-rmkbd.sh` with `RM_HOST=10.11.99.1` (USB) or reachable Wi-Fi IP

### 3. End-to-end on rM2 — PASS (edit-session)

- Phone WebSocket → daemon → `/run/Writerdeck.sock` → real editor: verified via `test-edit-session.sh`
- Remaining: `test-keyboard-harness.sh` on rM2, phone keyboard typing loop

### 4. Editor validation on Qt6

- Typing harness on rM2 — critical **55/57** (`editFontScale=3.0`). Lobby fonts: `lobbyFontScale=2.35` (point sizes only; layout from lobby-ui.json). Compare: `bash spike/rm2-official-sdk/scripts/compare-lobby-text.sh`; capture: `capture-lobby-screenshot.sh`. PNG: `docs/screenshots/writerdeck-rm2-lobby-keyboard-2026-08-22.png`.
- Update `docs/architecture.md` device facts (still rM1-centric in places)

### 5. CI and verify-all

- fork-probe added to `.github/workflows/spike-rm2-hello.yml`
- Optionally add fork-probe to `verify-all-spike.sh`

## Suggested start for a fresh session

1. Read this file and [README.md](README.md).
2. Wi-Fi IP: `bash scripts/discover-rm2-wifi.sh --write-secrets`
3. If harness 401: `bash scripts/configure-sync.sh 10.11.99.1` (needs `PIN_DIGITS=none` in secrets).
4. Wrap calibration: run `-s wrap-down-one-visual-line -v` on rM2; compare cursor after one Down from Ctrl+Home vs `wrap_fixtures.go` (expect ~10 at W=320, device reports ~30).
5. Fix fork `harnessSetWidth` / recalibrate `wrap_fixtures.go` for rM2; redeploy via `deploy-keywriter-rm2.sh`; re-run critical then full harness.

## Constraints

- No jailbreak / Toltec.
- Document integrity is paramount — see `docs/decisions.md`, `docs/integrity-audit.md`.
- Editor behavior lives in fork C++ (`EditHelper`), not `build-keywriter.sh`.
- Static Go server: `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7`.
- Cross-compile Writerdeck on Mac via Docker/GitHub Actions — not local docker for production binary.
- Do not `pkill -f /home/root/Writerdeck` (matches Writerdeck-server).

### Power button / sleep (rM2)

rM2 power is `snvs-powerkey` (often `/dev/input/event0`), not rM1 `gpio-keys` on `event1`. Writerdeck resolves button devices by sysfs name. Short-press power opens an on-screen menu: **Sleep / Lobby / Exit** while editing, **Sleep / Exit** in the Lobby (no Home button on rM2). A second short-press while the menu is open chooses Sleep. Sleep still uses `systemctl suspend`; this firmware can return at Sleep target before the device has actually slept, so the server waits for `systemd-suspend.service` to finish, then relaunches the editor and reopens the note. rM2 often resumes without delivering `KEY_POWER` to userspace (same class of bug KOReader fixed).

Verify: short-press power while editing → menu → Sleep → sleep screen → short-press → note back; Lobby and Exit from the menu; second power on an open menu = Sleep.
