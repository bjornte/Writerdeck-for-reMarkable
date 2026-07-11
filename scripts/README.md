# scripts/

Automation for the dev/deploy loop. Run from the repo root.

> Device work is bash-only, on a machine that can reach the tablet over Wi-Fi. The bash scripts read credentials/addresses from [../secrets/remarkable.local.env](../secrets/remarkable.local.env) (gitignored) via `_env.sh`.
>
> The Windows `*.ps1` scripts cover the off-device parts: cross-build (`build-rmkbd.ps1`), commits (`push.ps1`), and the git-bridge watcher (`watch-pc.ps1`).

## On-device names vs repo script names

Deploy scripts push binaries with Writerdeck-branded names on the tablet. Repo script names like `deploy-rmkbd.sh` and `deploy-keywriter.sh` are historical:

| On tablet | Built from |
|---|---|
| `/home/root/Writerdeck-server` | `daemon/` via `deploy-rmkbd.sh` |
| `/home/root/Writerdeck` | `third_party/keywriter/` (CI) via `deploy-keywriter.sh` |
| `/home/root/Writerdeck-launcher.sh` | `scripts/Writerdeck-launcher.sh` |
| `/home/root/Writerdeck-user-documents/` | notes directory |
| `/home/root/.Writerdeck/settings.json` | persisted prefs |
| `/run/Writerdeck.sock` | daemon ↔ editor socket |
| `writerdeck.service` | `scripts/writerdeck.service` |

`migrate-device-layout.sh` runs automatically from deploy/install scripts: renames legacy paths (`rmkbd`, `keywriter`, `edit/`, `.rmkbd/`, etc.) and removes old binaries.

## Scripts

| Script | Does |
|---|---|
| `_env.sh` | Shared helper: dot-sourced by the bash device scripts. Loads secrets; defines ssh/scp + key-test helpers. |
| `paths.sh` | Canonical on-device path constants (sourced by other scripts). |
| `migrate-device-layout.sh` | One-time rename of legacy on-device paths + removal of old binaries. Auto-run from deploy/install scripts. |
| `bootstrap.sh` | Generate an SSH keypair if absent; install the pubkey on the device (one password prompt); enable Wi-Fi SSH (`rm-ssh-over-wlan on`); verify key login. |
| `recon.sh` | Snapshot device facts: OS version, `ip addr`, input devices, disk. Self-logs to `../docs/recon/`. Re-run after a firmware update to refresh the facts. |
| `deploy-keywriter.sh` | (Mac) Ship CI-built `Writerdeck` binary + `qt5/` sysroot to `/home/root/`, launch via `Writerdeck-launcher.sh`, print a verdict, trap-restore xochitl. Self-logs to `../docs/recon/`. |
| `build-rmkbd.ps1` / `deploy-rmkbd.sh` | Cross-build Writerdeck-server (ARMv7 static, `CGO_ENABLED=0`). `deploy-rmkbd.sh` (Mac) also ships to `/home/root/Writerdeck-server` and kills any running instance. `build-rmkbd.ps1` (ThinkPad) builds only — device steps require the Mac. |
| `Writerdeck-launcher.sh` | The proven linuxfb launch env (panel geometry/DPI + epaper scene graph) in one place — sourced by `deploy-keywriter.sh` and by Writerdeck-server to spawn Writerdeck as its child. |
| `test-e2e.sh` | (Mac) Full browser→e-ink pipeline test: build+deploy Writerdeck-server → stop xochitl → launch Writerdeck + server → print the browser URL → hold for a human to type → show daemon log + `scratch.md` → restore xochitl. `-s` skips the server build+scp. Self-logs to `../docs/recon/`. |
| `test-edit-session.sh` | (Mac) Regression: `POST /api/open` (phone **Edit**) from stock UI must keep Writerdeck running, xochitl stopped, and `editorActive: true` for ~8 s. Prep starts server if needed; cleanup returns to stock UI without killing the daemon. Self-logs to `../docs/recon/`. |
| `push.ps1` / `push.sh` | One-line stage+commit+push. `push.ps1` bakes in the personal git identity to prevent the work-email-leak footgun. On the Mac, `rmpush` is the alias. |
| `install-alias.sh` | One-time Mac setup: adds the `rmpush` alias to `~/.zshrc`. |
| `watch-mac.sh` | Git-bridge auto-sync (Mac side). Pulls everything; auto commits+pushes only new outputs under `docs/recon/` (scoped for safety — edits elsewhere are reported, not committed). macOS GUI banners on arm / each sync / stop. |
| `watch-pc.ps1` | Git-bridge auto-sync (PC side). Loops `git pull`; pops a Windows toast when a pull brings in new commits. Banners on arm / each pull / stop. No admin, no modules. |
| `install-service.sh` | (Mac) Install `writerdeck.service` on the device: scp unit → `/etc/systemd/system/`, `daemon-reload`. Migrates off legacy unit names if present. Does not enable (boot-loop guard); prints the manual-start → enable → recovery steps. |
| `writerdeck.service` | systemd unit — runs `/home/root/Writerdeck-server` under `systemd-inhibit` (keep-awake), stops/restores xochitl around it. Installed by `install-service.sh`. |

## Convention: device actions become committed scripts

Don't hand-type long `ssh root@…` one-liners — script every device action so it runs as one short line, and `tee` device output to `docs/recon/` so verdicts are captured (`recon.sh`, `deploy-keywriter.sh`, `test-e2e.sh` already do). Never log a secret there — `bootstrap.sh` echoes the root password, so it isn't logged. Optional auto-sync watchers: `bash scripts/watch-mac.sh` + `./scripts/watch-pc.ps1`.

## Conventions
- Iterate over Wi‑Fi (`192.168.1.8`) — the working path on the Mac (USB‑ethernet is dead there: no macOS RNDIS). Scripts default to it via `$RM_HOST`; override with `export RM_HOST=10.11.99.1` if USB ever revives.
- Scripts never hardcode the password; they read it from the secrets file at runtime.
- Keep scripts idempotent and re-runnable.
