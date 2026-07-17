# Jargon to avoid

Plain speech for the owner (a designer). Real technology names stay. Process slang goes.

Use this list when searching docs later. Prefer the pattern plain gloss (RealName) on first mention in a file — for example “the screen file (QML)” or “fake keyboard device (uinput)” — then RealName alone is fine.

## Keep these names

QML, C++, Qt, TextEdit, EditHelper, Go, CI, GitHub Actions, Docker, systemd, SSH, WebSocket, AES-GCM, scrypt, gpio, evdev, qmap, uinput, Toltec, OTA, linuxfb, sysroot.

Do not replace them with euphemisms only. Wrong: “the screen file” forever. Right: “the screen file (QML)” or just “QML” after that.

## Avoid these process words in prose

upstream — say Dave’s original (or the original keywriter). Keep `upstream` only as a git remote name, with a gloss.

downstream, ship tip, tip (for a commit) — say known-good commit or commit hash.

harness, green, gate, sign-off, smoke test, triage, flake, LOC, Patch LOC, ADR, hygiene (for docs), artifact (for “look at this file”), acceptance surface, false completion, behavior-identical, bit-identical, cheapest proof, brain/hands/spine (metaphors), fudge, calibrated fudge, CRUD (spell out create/rename/delete), opaque (say “copied without decrypting”), reconcile (say two-way sync / copy missing notes), inject/synthetic as unexplained slang (say feed keystrokes / fake Qt key events, and name QKeyEvent when useful).

Score slang like “harness green” — say “all 110 typing checks passed.”

## Script and path names

Keep real script names (`test-keyboard-harness.sh`, `build-keywriter.sh`). Gloss once if helpful: “automated typing tests (`test-keyboard-harness.sh`)”.

## Audience reminder

Standing docs and READMEs: John Gruber — whole sentences, little markup, readable by a layperson. Agent rules and historical handoffs may keep denser shop talk; mark archives as historical.
