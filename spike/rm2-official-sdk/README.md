# Spike: official Codex SDK on reMarkable 2

Branch: `spike/rm2-official-sdk`.

Goal: prove we can draw on an rM2 with reMarkable’s public Qt ePaper path (no Toltec, no `rm2fb`), then decide whether a Writerdeck port is worth it.

## Device under test (Aug 2026)

| Fact | Value |
|------|--------|
| Host | `RM2_HOST_WIFI` in `secrets/remarkable.local.env` |
| Software | `IMG_VERSION` 3.11.2.5 |
| Codex / os-release | 4.0.447 (kirkstone) |
| Kernel | 5.4.70-v1.3.4-rm11x |
| Qt on device | 6.5.2 |
| Platform plugin | `/usr/lib/plugins/platforms/libepaper.so` present |
| Framebuffer | `/dev/fb0` (`mxs-lcdif`) — not the rM1 EPDC path |

Closest public SDK installer found: OS image **4.0.367** rm2 x86_64 (no 4.0.447 upload). URL in `scripts/sdk-url.sh`.

## Layout

- `hello_remarkable/` — upstream [reMarkable developer example](https://github.com/reMarkable/remarkable-developer-examples) (Qt 6 Quick)
- `scripts/build-hello.sh` — download SDK (if needed), cmake, build (Linux x86_64 or Docker)
- `scripts/deploy-hello.sh` — copy binary to tablet
- `scripts/run-hello.sh` — stop `xochitl`, run with `-platform epaper`, restore `xochitl` on exit
- `scripts/recon-rm2.sh` — firmware / Qt / fb facts

Local SDK + build output live under `spike/rm2-official-sdk/.cache/` (gitignored).

## Mac note

Official SDK host is **Linux x86_64**. On this Mac, build via Docker/Colima (`scripts/build-hello.sh`) or GitHub Actions (`.github/workflows/spike-rm2-hello.yml`).

## Pass criteria

1. Hello text visible on the rM2 e-ink after `run-hello.sh`
2. Touch toggles text (official sample behaviour)
3. `systemctl start xochitl` restores stock UI
4. No Toltec installed

## Docs

Official: [Qt ePaper](https://developer.remarkable.com/documentation/qt_epaper), [SDK](https://developer.remarkable.com/documentation/sdk), [links / downloads](https://developer.remarkable.com/links).
