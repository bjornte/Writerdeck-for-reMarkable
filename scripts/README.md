# scripts/

Automation for the dev/deploy loop. Run from the repo root.

> Device work is bash-only, on a machine that can reach the tablet over Wi-Fi. The bash scripts read credentials/addresses from [../secrets/remarkable.local.env](../secrets/remarkable.local.env) (gitignored) via `_env.sh`.
>
> The Windows `*.ps1` scripts cover the off-device parts: cross-build (`deploy.ps1`), commits (`push.ps1`), and the git-bridge watcher (`watch-pc.ps1`).

## Scripts

| Script | Does |
|---|---|
| `_env.sh` | Shared helper: dot-sourced by the bash device scripts. Loads secrets; defines ssh/scp + key-test helpers. |
| `bootstrap.sh` | Generate an SSH keypair if absent; install the pubkey on the device (one password prompt); enable Wi-Fi SSH (`rm-ssh-over-wlan on`); verify key login. |
| `recon.sh` | Snapshot device facts: OS version, `ip addr`, input devices, disk. Self-logs to `../docs/recon/`. Re-run after a firmware update to refresh the facts. |
| `deploy-keywriter.sh` | (Mac) `scp` the from-source `keywriter` binary + `qt5/` sysroot to `/home/root/`, launch via `launch-keywriter.sh`, print a verdict, trap-restore xochitl. Self-logs to `../docs/recon/`. |
| `deploy.ps1` / `deploy.sh` | Cross-build `rmkbd` (ARMv7 static, `CGO_ENABLED=0`). `deploy.sh` (Mac) also scps the binary to the device and kills any running instance. `deploy.ps1` (ThinkPad) builds only — device steps require the Mac. |
| `launch-keywriter.sh` | The proven linuxfb launch env (panel geometry/DPI + epaper scene graph) in one place — sourced by `deploy-keywriter.sh` and by `rmkbd` to spawn keywriter as its child. |
| `test-phase4.sh` | (Mac) Full browser→e-ink pipeline test: build+deploy `rmkbd` → stop xochitl → launch keywriter + `rmkbd` → print the browser URL → hold for a human to type → show daemon log + `scratch.md` → restore xochitl. `-s` skips the rmkbd build+scp. Self-logs to `../docs/recon/`. |
| `push.ps1` / `push.sh` | One-line stage+commit+push. `push.ps1` bakes in the personal git identity to prevent the work-email-leak footgun. On the Mac, `rmpush` is the alias. |
| `install-alias.sh` | One-time Mac setup: adds the `rmpush` alias to `~/.zshrc`. |
| `watch-mac.sh` | Git-bridge auto-sync (Mac side). Pulls everything; auto commits+pushes only new outputs under `docs/recon/` (scoped for safety — edits elsewhere are reported, not committed). macOS GUI banners on arm / each sync / stop. |
| `watch-pc.ps1` | Git-bridge auto-sync (PC side). Loops `git pull`; pops a Windows toast when a pull brings in new commits. Banners on arm / each pull / stop. No admin, no modules. |
| `install-service.sh` | (Mac) Install `rm1-writerdeck.service` on the device: scp unit → `/etc/systemd/system/`, `daemon-reload`. Migrates off the old `rmnetwriter.service` name if present. Does not enable (boot-loop guard); prints the manual-start → enable → recovery steps. |
| `rm1-writerdeck.service` | systemd unit — runs `rmkbd` under `systemd-inhibit` (keep-awake), stops/restores xochitl around it. Installed by `install-service.sh`. |

## Convention: device actions become committed scripts

Don't hand-type long `ssh root@…` one-liners — script every device action so it runs as one short line, and `tee` device output to `docs/recon/` so verdicts are captured (`recon.sh`, `deploy-keywriter.sh`, `test-phase4.sh` already do). Never log a secret there — `bootstrap.sh` echoes the root password, so it isn't logged. Optional auto-sync watchers: `bash scripts/watch-mac.sh` + `./scripts/watch-pc.ps1`.

## Conventions
- Iterate over Wi‑Fi (`192.168.1.8`) — the working path on the Mac (USB‑ethernet is dead there: no macOS RNDIS). Scripts default to it via `$RM_HOST`; override with `export RM_HOST=10.11.99.1` if USB ever revives.
- Scripts never hardcode the password; they read it from the secrets file at runtime.
- Keep scripts idempotent and re-runnable.
