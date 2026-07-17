# Copilot instructions — Writerdeck for reMarkable 1

Writerdeck for reMarkable 1 turns a reMarkable 1 e-paper tablet into a distraction-free Markdown typewriter. Typing is enabled either over a USB keyboard with OTG cable or by forwarding keystrokes from an iPhone keyboard over Wi-Fi. Writerdeck-server (`/home/root/Writerdeck-server`, built from `daemon/`) receives key events and feeds the patched Writerdeck editor (upstream remarkable-keywriter, reads `/run/Writerdeck.sock`), which saves `.md` to `Writerdeck-user-documents/`.

For background, consult [todo](../TODO.md), [architecture](../docs/architecture.md), [decisions](../docs/decisions.md), [done](../DONE.md), and [lessons](../docs/lessons.md).

**Cursor agents:** project rules live in [`.cursor/rules/writerdeck.mdc`](../.cursor/rules/writerdeck.mdc) (`alwaysApply: true`). Keep that file in sync when you change operational instructions here.

## Constraints (honor these)

- No jailbreak; preserve OTA firmware updates => avoid Toltec.
- No on-device runtime deps => static Go binary (`CGO_ENABLED=0`, `GOOS=linux GOARCH=arm GOARM=7`).
- Markdown is the save format.
- Writerdeck-server cross-compiles on the Mac (the only device-reachable host) — `GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0`. Only Writerdeck needs CI/Docker (upstream keywriter + toltec Qt sysroot); don't build it on a host toolchain.

## Dev setup

Mac on the same Wi-Fi as the tablet: build, deploy, test locally. Secrets in `secrets/remarkable.local.env`; scripts source `scripts/_env.sh` for `RM_HOST`.

Daemon loop: `deploy-rmkbd.sh` (embeds `daemon/*` via `go:embed` → `/home/root/Writerdeck-server`) → `systemctl restart writerdeck` (deploy kills the running server). Verify: `curl http://$RM_HOST:8000/` + `/app.js`; UI at same URL (PIN `none` skips auth). GitHub token: browser `localStorage ghToken` + tablet RAM via `POST /api/sync/token`; after restart the server sends WebSocket `needtoken` and the browser reposts automatically — journal shows `sync reconcile (token)` after `client connected` (see [server-sync-implementation.md](../docs/server-sync-implementation.md)). **Writerdeck:** `git push` → `fetch-keywriter-dist.sh` → `deploy-keywriter.sh -b` (`rmkw`) — never local `docker build` on Mac; CI runs Docker in GHA; relaunch Writerdeck after binary deploy (server restart alone is not enough). After Writerdeck deploy, run `bash scripts/test-edit-session.sh` — `POST /api/open` must keep Writerdeck up (guards QML patch regressions that flash stock UI). Session logs: `journalctl -u writerdeck.service` (`editor started` / `editor process exited`, QML load errors, `sync reconcile (token)` for token restore).

## Verify on device (mandatory — do not skip)

**Deploy success ≠ tested.** Agents and humans must verify on the tablet before calling work done. After `daemon/`, `build-keywriter.sh`, or `lobby/` changes: rebuild Writerdeck if QML changed → deploy → SSH → `journalctl -u writerdeck -n 30` (fail on QML parse errors or instant `editor process exited`) → `bash scripts/test-edit-session.sh`; after Lobby or `handleHome` QML also `bash scripts/test-lobby-keyboard.sh` ([decisions](../docs/decisions.md) §15). You cannot see e-ink, but you **can** read logs and run these scripts — **do it every time**, without being asked.

On-device naming: see [architecture.md](../docs/architecture.md) and [decisions.md](../docs/decisions.md) #29. Do not `pkill -f /home/root/Writerdeck` (matches Writerdeck-server).

## Doc hygiene

- Update [todo](../TODO.md) and [done](../DONE.md) as work progresses; durable *how-it-works* → [architecture](../docs/architecture.md); *why* → [decisions](../docs/decisions.md); operational gotchas → [lessons](../docs/lessons.md). Do not create new summary markdown files — fold into these standing docs.
- Keep TODO/DONE terse but human-friendly: plain language, durable lessons over blow-by-blow. Prune step-by-step narration, commit hashes, and durations once a lesson is settled; don't repeat the date within a day.
- Prune dead files as paths close — not just prose. When a phase finishes or a path is declared closed-failed, delete the machinery that only served it (one-off scripts, generated logs, build harnesses, committed binaries) — but capture the durable lesson in TODO/DONE first. Git history is the archive: keep the gravestone (why it failed), drop the tooling. Before deleting, grep for references and fix link rot; never remove a file a live script still calls. Regenerable outputs (e.g. `docs/recon/` logs) are always safe to drop — keep the folder via `.gitkeep`.

## Encoding & line endings (this bit us before — honor exactly)

- Executable / device files = ASCII-only + LF. `.sh`, `.service`, `Dockerfile`, `.go`, `.yml` run on the reMarkable (Linux) or in CI; a stray non-ASCII byte or a CRLF breaks the device shell / systemd. `.gitattributes` already normalizes these to LF in the repo (every file reads `i/lf` under `git ls-files --eol`, even if a Windows checkout shows `w/crlf`) — keep it that way; never add a BOM.
- Markdown `.md` = Unicode is fine and intentional (em-dashes, arrows, status emoji such as ✅ 🔴 ⬜). The ASCII rule is for code, not prose.
- Before committing a script, grep it for `[^\x00-\x7F]` — should be empty for every `.sh`/`.ps1`.

