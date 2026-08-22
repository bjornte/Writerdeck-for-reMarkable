# Spike: official Codex SDK on reMarkable 2

Branch: `spike/rm2-official-sdk`.

Goal: prove we can draw on an rM2 with reMarkable's public Qt ePaper path (no Toltec, no `rm2fb`), then decide whether a Writerdeck port is worth it.

## Result: PASS (Aug 2026)

On **software 3.27.3.0** (Codex 5.7.126), after OTA:

- Stock ships `libqsgepaper.so` under `/usr/lib/plugins/scenegraph/`.
- Build hello with public SDK **3.27.0.97** rm2 (`scripts/sdk-url.sh`).
- Run: `QT_QUICK_BACKEND=epaper`, `-platform epaper`, touch params `rotate=180:invertx`.
- **"Hello reMarkable!" visible on device** (owner confirmed).
- `systemctl start xochitl` restores stock UI.
- No Toltec.

**3.11.2.5 was too old** for this path (no `libqsgepaper` on device; epaper Quick backend aborts).

**Framebuffer PNG capture:** `/dev/fb0` is not a readable mirror on rM2 software epaper. Apps write `/home/root/spike-rm2-*/screen.png` via `QQuickWindow::grabWindow()`; pull with `capture-screenshot.sh` or `run-and-capture*.sh`.

**Autonomous verify (no eyeballing):**

```bash
bash spike/rm2-official-sdk/scripts/verify-spike.sh hello
bash spike/rm2-official-sdk/scripts/verify-spike.sh textedit
bash spike/rm2-official-sdk/scripts/verify-spike.sh socket
```

Picks USB then Wi-Fi (`rm2-pick-host.sh`), runs the app, pulls PNG, checks 1404x1872 and enough dark pixels (`verify-png.sh`).

## Port phases (spike tree only)

| Phase | Proves | Status |
|-------|--------|--------|
| 1 hello | Official epaper QPA draws | PASS |
| 2 textedit | Wrapped prose on epaper | PASS |
| 3 socket | NDJSON keys into on-screen text | PASS |
| 4 Writerdeck fork probe | Full fork binary on epaper | PASS (Lobby UI) |

Full port is Qt6 + `-platform epaper` in the fork (CMake CI, not Toltec). Probe: `bash spike/rm2-official-sdk/scripts/build-fork-probe.sh` then `verify-spike.sh fork-probe`.

```bash
bash spike/rm2-official-sdk/scripts/verify-all-spike.sh
bash spike/rm2-official-sdk/scripts/verify-spike.sh fork-probe
```

Production launcher for rM2 when the fork ships: `scripts/Writerdeck-launcher-rm2.sh`.

## Device under test

| Fact | Value |
|------|--------|
| Host | `RM2_HOST_WIFI` / USB `10.11.99.1` in `secrets/remarkable.local.env` |
| Software | 3.27.3.0 |
| Codex | 5.7.126 (scarthgap) |
| Qt | 6.x + `libepaper.so` + `libqsgepaper.so` |
| SDK used | 3.27.0.97 rm2 x86_64 (CI + `scripts/sdk-url.sh`) |

## Commands

```bash
# Build (Mac: Docker/Colima, or CI artifact)
bash spike/rm2-official-sdk/scripts/build-hello.sh

# Deploy + run (~45s on screen)
RM2_HOST=10.11.99.1 bash spike/rm2-official-sdk/scripts/deploy-hello.sh
RM2_HOST=10.11.99.1 bash spike/rm2-official-sdk/scripts/run-hello.sh

# Facts
bash spike/rm2-official-sdk/scripts/recon-rm2.sh
```

After firmware OTA: SSH host key changes — `ssh-keygen -R 10.11.99.1` (and Wi-Fi IP). Password may change; update `RM2_ROOT_PASSWORD`.

## Layout

- `hello_remarkable/` -- upstream [reMarkable developer example](https://github.com/reMarkable/remarkable-developer-examples)
- `scripts/build-hello.sh` -- Codex SDK cmake build (Linux x86_64 or Docker)
- `scripts/deploy-hello.sh` -- copy binary (+ optional bundled `libqsgepaper.so` from CI)
- `scripts/run-hello.sh` -- stop `xochitl`, epaper hello, restore `xochitl`
- `scripts/recon-rm2.sh` -- firmware / Qt facts
- `scripts/capture-screenshot.sh` -- experimental fb0 PNG (not reliable for epaper apps)

Local SDK + build output: `spike/rm2-official-sdk/.cache/` (gitignored).

## What this does not prove

Writerdeck itself. Next step would be a Qt6 port of the editor (EditHelper, typing harness, launch path) on the same display stack -- large, separate project.

## Docs

Continue in a fresh session: [HANDOFF.md](HANDOFF.md).

[Qt ePaper](https://developer.remarkable.com/documentation/qt_epaper), [SDK](https://developer.remarkable.com/documentation/sdk), [links](https://developer.remarkable.com/links).
