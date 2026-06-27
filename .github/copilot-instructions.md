# Copilot instructions — Writerdeck for reMarkable 1

rM1-Writerdeck turns a reMarkable 1 e-paper tablet into a distraction-free Markdown typewriter by forwarding keystrokes from an iPhone keyboard over Wi-Fi. A static Go daemon (`rmkbd`) on the tablet receives key events and feeds them to the third-party keywriter editor (patched to read a local socket), which saves `.md`.

> The architecture in one line: keystrokes reach keywriter over a local socket, *not* `/dev/uinput` — this kernel can't load uinput (exports trimmed). Don't retry the uinput build. Why: DONE.md.

Read first, every session: [TODO.md](../TODO.md) (the open plan — what's left to build), [docs/architecture.md](../docs/architecture.md) (how it works), [docs/decisions.md](../docs/decisions.md) (the why — the ADR), and [DONE.md](../DONE.md) (what's done). Work the first unchecked item; don't redo finished ones.

## Model usage protocol (Opus <-> Sonnet) — the operating rule

Match the model to the task, and proactively tell the user when to switch — state the trigger and what the next model should do first. Never switch silently; the user drives the model picker. When switching models within the same Claude session/chat, the context is preserved. Only the model changes; the conversation history/context remains.

- Opus — plan & inspect. Research, architecture, risk analysis, writing/expanding the plans (TODO/DONE), reviewing & inspecting code, diagnosing non-obvious failures.
- Sonnet — write code. Once a phase's spec is settled, author the files: Dockerfile, build/install scripts, Go source, systemd units, keymap tables, boilerplate, repetitive edits.

| Switch | When |
|---|---|
| -> Sonnet | The active phase's Done-when spec is settled and the remaining work is typing code/files. |
| -> Opus | Code needs review/inspection; or a decision / ambiguity / architecture / risk question arises; or a test fails in a non-obvious way. |

Example: *"The socket protocol + keywriter input-source patch are specced — prudent to switch to Sonnet to write the Go feeder + the keywriter patch; have it read the active TODO.md slice + [docs/decisions.md](../docs/decisions.md) first. Switch back to Opus to inspect the result."*

## Constraints (honor these)
- No jailbreak; preserve OTA firmware updates => avoid Toltec.
- No on-device runtime deps => static Go binary (`CGO_ENABLED=0`, `GOOS=linux GOARCH=arm GOARM=7`).
- Markdown is the save format.
- rmkbd (static Go binary) cross-compiles on the Mac (the only device-reachable host) — `GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0`. Only keywriter needs CI/Docker (the toltec Qt sysroot); don't build it on a host toolchain.

## Two-machine dev setup
Development is split across two machines bridged by git: the assistant + commits run on one machine with no tablet access; the other does all device SSH/deploy over Wi-Fi `192.168.1.8`. Honor the device conventions every session: script device actions (don't hand-type long `ssh` one-liners), `tee` device output to `docs/recon/`, and never log a secret there.

## Doc hygiene
- Update [TODO.md](../TODO.md) (the open plan) and [DONE.md](../DONE.md) (newest log entry on top, grouped under one date heading per day) as work progresses; durable *how-it-works* goes in [docs/architecture.md](../docs/architecture.md) and *why* in [docs/decisions.md](../docs/decisions.md). Do not create new summary markdown files for changes — fold them into these standing docs.
- Keep TODO/DONE terse but human-friendly: plain language, durable lessons over blow-by-blow. Prune step-by-step narration, commit hashes, and durations once a lesson is settled; don't repeat the date within a day.
- Prune dead files as paths close — not just prose. When a phase finishes or a path is declared closed-failed, delete the machinery that only served it (one-off scripts, generated logs, build harnesses, committed binaries) — but capture the durable lesson in TODO/DONE first. Git history is the archive: keep the gravestone (why it failed), drop the tooling. Before deleting, grep for references and fix link rot; never remove a file a live script still calls. Regenerable outputs (e.g. `docs/recon/` logs) are always safe to drop — keep the folder via `.gitkeep`.

## Encoding & line endings (this bit us before — honor exactly)
- Executable / device files = ASCII-only + LF. `.sh`, `.service`, `Dockerfile`, `.go`, `.yml` run on the reMarkable (Linux) or in CI; a stray non-ASCII byte or a CRLF breaks the device shell / systemd. `.gitattributes` already normalizes these to LF in the repo (every file reads `i/lf` under `git ls-files --eol`, even if a Windows checkout shows `w/crlf`) — keep it that way; never add a BOM.
- PowerShell `.ps1` = ASCII-only too — a single non-ASCII char once broke PowerShell parsing (wrong-encoding read). `.ps1` deliberately keep CRLF (Windows-native, per `.gitattributes`); that's correct, leave it. _(Known straggler: `push.ps1` carries the author's accented name — clean to ASCII when next touched, though the name is a deliberate exception.)_
- Markdown `.md` = Unicode is fine and intentional (em-dashes, arrows, status emoji such as ✅ 🔴 ⬜). The ASCII rule is for code, not prose.
- Before committing a script, grep it for `[^\x00-\x7F]` — should be empty for every `.sh`/`.ps1`.
