# Lessons learned

Operational gotchas from building Writerdeck. Why things are the way they are: [decisions.md](decisions.md). What's shipped: [../DONE.md](../DONE.md).

## Device verify and iteration

Debugging and shipping are different jobs. While iterating, use the cheapest check that could disprove your guess. Before calling work done, run the full verify path in the project rules.

Rank checks by cost: unit/API, device harness, browser UI, full E2E with sync. Climb only when the cheaper step passes.

When something fails mid-flow, test from that step — not the whole story. Use harness scoping (`-s`, `-m`, `--no-prepare`) and the script that matches the layer you changed.

For keyboard harness work specifically: one full run to collect all failures, then batch fixes by layer (harness vs QML). Do not interleave push/CI/deploy with single-scenario retries. See [editor-testing/todo.md](editor-testing/todo.md) § Dev loop.

Verify rules are for sign-off. While fixing: server → `deploy-rmkbd.sh` + API/harness; QML → `deploy-keywriter.sh -b` + targeted harness; phone UI and sync → full pipeline last.

A failed run updates one hypothesis. Ask a sharper question; don't rerun the world to learn the same thing again.

After about 20 minutes on one task without sign-off, stop and report methodology (see [.cursor/rules/writerdeck.mdc](../.cursor/rules/writerdeck.mdc) § Time budget while debugging).

## Deploy and staleness

Three separate ways a change can look like it did nothing: the CI keywriter binary lags the git push; the browser caches the capture page (serve with `Cache-Control: no-store`); a live editor session keeps the old binary until you respawn it.

`rmkw` pushes the binary only. Font changes need the full Qt sysroot: `RM_FORCE_SYSROOT=1 bash scripts/deploy-keywriter.sh -b`.

`deploy-rmkbd.sh` flushes via `POST /api/flush-save`, then SIGTERM-waits about twelve seconds. The old pattern of `pkill` plus half a second could drop an unsaved buffer. Deploy stops the server — run `systemctl start writerdeck` (or restart) before verifying the phone UI.

scp deadlocks at a fixed offset on the Mac-to-Wi-Fi-to-tablet link. Use `rm_send_file` (gzip-over-ssh) in `_env.sh`. On ETXTBSY, kill by full path, stream to `.new`, then `mv`.

Browser rotate needed both binaries current when `POST /api/rotate` existed on the phone; that path is removed. The server still saves rotation in settings and pushes `setrotation` on connect; Writerdeck must handle it. USB Ctrl+arrow relays `rotationChanged` via `rotation_watcher`.

## systemd and device

With `writerdeck.service`, both server and Writerdeck Qt stderr land in `journalctl -u writerdeck.service`.

Any unit whose ExecStart lives on `/home/root` needs `RequiresMountsFor=/home/root` or cold boot races the mount.

Export `HOME=/home/root` in Writerdeck-launcher.sh — under systemd, root's home is `/`, which breaks the save path.

Never `pkill -f /home/root/Writerdeck`; that pattern matches Writerdeck-server too. Use `pidof Writerdeck` for the editor and `pkill -f Writerdeck-server` for the server.

The tablet drops Wi-Fi on suspend. Keep it awake during dev.

## Writerdeck and QML

Every save path in edit mode must sync `query.text` into `doc` before `saveFile()`. A bare `saveFile()` writes stale content.

Never clear `query.text` without re-syncing on load — it breaks the TextEdit binding. Call `syncQueryDisplay()` after load or mode switch. If you skip this, Home save can zero the file (`save -> 0` in the journal) or Esc-toggle can show corrupted preview.

Preview is imperative, not bound. `toggleMode()` and `doLoad` must call `syncQueryDisplay()` and must never read RichText back into `doc`.

The Lobby Home wipe bug (fixed): a binding bug plus sync pushing empties created `(tablet copy).md` junk. Prevention is slice 2 plus the empty-push guard. Recovery is GitHub history.

Python comments outside string literals in `build-keywriter.sh` heredocs cause SyntaxError in CI. Use `#` on the Python side only.

Socket-triggered saves ack back to the server — the server waits before exitedit, GitHub push, or HTTP 200 on open. Power sleep sends ready after the sleep screen paints (~800 ms).

Clear `currentFile` on every return to the Lobby or a stale name resurrects deleted notes.

Ctrl-K with injected keys must check `Qt.ControlModifier`, not only the standalone Control key bool.

USB Escape launch uses evdev watch with three-second hotplug rescan; idle only. Pin the USB keyboard device in the launcher — a bare `keymap=…:grab=1` grabs event1 and starves the Home and Power watcher. L+R page buttons use the same idle path with 800 ms debounce.

Standard Linux kmaps map Alt+Left/Right to `Decr_Console` / `Incr_Console`. Qt evdev turns those into fake `Key_Escape` events; Writerdeck toggles edit/preview on Esc release, so Alt+arrow looked like entering read mode. Override in `keymaps/src/i386/include/writerdeck-alt-arrows.inc`, regenerate with `bash keymaps/generate.sh`, redeploy keymaps, relaunch Writerdeck. `handleKey` also ignores Esc release when Alt or Ctrl is held. Qmap applies at editor launch, not mid-session.

## Keyboard and selection

Phone path and USB path are different inputs. The harness (`daemon/cmd/edit-harness`, `scripts/test-keyboard-harness.sh`) drives keys over WebSocket — same as Type mode, not Qt evdev. It will not catch qmap bugs; those need hardware or a future evdev probe.

Three tiers: `go test -C daemon -run TestTranslate` for modifier masks in `translate()`; the device harness for cursor/selection via `GET /api/test/editor-state` (Writerdeck publishes `editorstate` over the socket); manual USB checks after qmap or launcher changes.

After QML selection or arrow-handler edits: rebuild Writerdeck, relaunch, run `bash scripts/test-keyboard-harness.sh`. After server test API edits: `deploy-rmkbd.sh` too. Per-run logs: `docs/recon/test-keyboard-harness-*.{md,txt}` (each run); consolidated history: [recon/harness-runs.md](recon/harness-runs.md); milestone table: [editor-testing/milestone-runs.md](editor-testing/milestone-runs.md).

Sandbox-prepare (default): one editor launch per full run; between scenarios `PUT` note content plus `harnessprepare` (in-process reset, no quit). Do not use `POST /api/open` to reload the harness note — `saveAndLoad` writes the stale in-memory buffer over the `PUT` first. `--hard-reset` was removed from `test-keyboard-harness.sh`. Single scenario: `bash scripts/test-keyboard-harness.sh -s NAME --fast`.

Fast dev loop: [editor-testing/](editor-testing/) — add scenario, `--unit`, full triage run, batch fix, one deploy. Per-scenario: `-s NAME --fast --no-prepare` on the same binary. Harness changes need no Writerdeck deploy unless `/api/test/*` changed.

### Harness batch workflow

A keyboard harness session should produce one failure list, not a deploy per guess.

1. `--unit`, then `--fast` once (full suite, single session). Read per-run report; update [milestone-runs.md](editor-testing/milestone-runs.md).
2. Confirm each FAIL with `-s NAME --fast` on the current binary (no deploy between).
3. Fix all harness/prepare failures in `edit-harness` — no push for QML yet.
4. Batch QML fixes in one `build-keywriter.sh` diff. One push/CI/deploy for that batch.
5. Rerun full suite `--fast`; compare to [milestone-runs.md](editor-testing/milestone-runs.md).

Deploy budget while iterating: at most one Writerdeck binary deploy per agent session unless the tablet binary failed to launch (QML parse error, editor never connects). Harness-only and daemon-only changes never need `fetch-keywriter-dist.sh`.

Edit-mode socket keys must reach `handleMacArrow` via QML `socketRouteKey()` invoked from the **inject thread** (`BlockingQueuedConnection`, same as `harnessprepare`). Routing `QKeyEvent` to `activeFocusItem` or nesting `invokeMethod` on the GUI thread dropped keys or deadlocked. Block **Ctrl/Alt nav releases** — Qt TextEdit defaults wiped `query.text` while the file on disk stayed intact.

Harness `primeModifiedKeys` uses ArrowUp wake only — never End (EOF poison). Plain Ctrl+nav from cursor 0 is fixed in kernel (`22ad701`), not by scenario End steps.

Wrap harness sets `query.width` to 320 for calibration. `harnessSandboxReset` must call `harnessSetWidth(0)` when width is not requested — re-applying `harnessTextWidth` stuck the live editor in a narrow column after a harness run. Restore full width on Home, Lobby, and `doLoad` too.

On device Qt, plain Backspace via `query.select` + `query.insert("")` selects the previous character without deleting — use `query.text` slice (same as Alt/Ctrl backspace in `handleMacBackspace`).

Plain `Key_Home` **release** in edit mode used to call `handleHome()` → lobby and break `combo-shift-end-line`; skip lobby when `mode==1 && !isLobby` on Home release (line-start is press via `handleMacArrow`).

Never add a second `Keys.onPressed` on query TextEdit — patch 6c caused `Property value set multiple times` and a crash loop. Mac key routing belongs in patch 7o prepended to query's existing handler.

Anti-pattern that wasted a full session: fix one scenario → push → CI → deploy → `-s` one scenario → repeat. Correct pattern: triage all → batch fix → one deploy → rerun all failures.

Do not mark keyboard editing done when only newline-based harness scenarios pass. Wrapped paragraphs and Shift+Alt/Ctrl+arrow combos were explicit scope; `\n`-only tests do not cover them.

`TextEdit.moveCursorSelection` takes a character index, not `TextEdit.Down` / `TextEdit.Up`. Passing direction enums selects toward a low position and breaks shift+vertical. Use `lineDownPos` / `lineUpPos` and explicit anchor math (same model as horizontal `extendSelectionHorizontal`). Setting `query.cursorPosition` after `query.select()` collapses the selection.

Saved-file guessing for selection tests was unreliable. Assert `cursor`, `selStart`, `selEnd`, and `textLen` from the editor-state probe instead.

Qt RichText ignores margin-bottom on paragraphs and list items — use line-height. Font IDs must match Qt family names exactly. QML Text needs explicit width and wrapMode.

Apostrophes in Python patch heredocs: `' + chr(39) + '`. QML patch blocks must balance braces — symptom is `Expected token '}'` at EOF and Writerdeck exits on every launch. The build asserts balance in `handleKey()`; also sanity-check full patched QML before deploy.

Immediate `editor process exited` after start is almost always a QML parse error, not the USB watcher. Two-level Home looks like a crash in logs but is intentional: first Home to Lobby, second to stock UI. No cursor blink on e-ink — it ghosts.

Physical Home on gpio-keys delivers twice: Go sends `cmd home`, Qt sends `Key_Home`. Without `suppressNextHomeKey` pairing (or a future exclusive grab on event1), read → Home could quit Writerdeck instead of returning to the Lobby. See [decisions.md](decisions.md) §28 and [handoff-physical-home-input.md](handoff-physical-home-input.md).

After Home from edit, Lobby USB keys failed when `lobbyFocus.forceActiveFocus()` hit a `FocusScope` with no `Keys` handlers — delegate to `handleKey` / `handleKeyDown` / `handleKeyUp`, and set `query.focus: !isLobby` so the editor does not compete. Files is tab 1; vault e2e sends digit `1` for Files and `4` for Settings — update harness key numbers when tab order changes.

`cursorOnLastLine()` must use visual line position (`positionToRectangle`), not "no newline after cursor" — a wrapped last paragraph is one logical line but many visual lines; newline-only detection makes Down jump to end-of-line mid-paragraph. Auto-scroll via `ensureVisible` on every cursor move can feel like blanking or page-flips near the document end on e-ink; scroll only when the cursor nears the viewport edge.

## Browser and capture page

Capture must stand down when the PIN screen, paste modal, or settings/sync overlays are up — `followTabletOpen` checks the same gate before auto-entering Type mode on `openedit`. Setting `display: ''` does not restore visibility if CSS had `display:none` — set an explicit value. Inline onclick cannot reach IIFE closures; use addEventListener.

Clipboard API needs HTTPS; on plain LAN http, Copy falls back to execCommand.

Lobby last-sync needs the phone to POST sync ack after reconcile. Load sync flags at page init, not when Sync setup opens. Async functions must return their promises.

GitHub token is per browser origin (`localStorage ghToken`). A new tablet IP is a new origin — enter the token once in Notes sync setup, then it persists in that browser. After a service restart, tablet RAM is empty until the browser reposts (WebSocket `needtoken` or `refreshSyncStatus` on reconnect). Verify in journal: `sync reconcile (token)` after `client connected` — see [server-sync-implementation.md](server-sync-implementation.md).

Writerdeck deploy needs a fresh binary; QML lives inside it. After lobby edits: run `third_party/keywriter/concat-lobby.sh`, commit `lobby_subpages.qml.inc`, CI rebuild, deploy, relaunch, check journalctl. Restarting the server does not reload a running Writerdeck process.

Lobby Files vault row: when private notes is on, the note `ListView` must reserve height for the second button row — otherwise Encrypt/Decrypt renders below the visible area. Use explicit half-width `Rectangle` buttons (same pattern as Settings); a `Repeater` that sizes delegates with `parentRow.model.length` fails because `Row` has no `model` — labels draw at x=0 with zero-width chrome and overlap.

Lobby Files inline rename/new: handle printable keys on key release only in `lobbyHandleKey` — a parallel `Keys.onPressed` insert duplicates characters from the phone WebSocket path (press and release both carry text). Rename strips `.md.enc` before `.md` for the editable basename; re-append `.md.enc` on submit for encrypted notes. Use `lobbyFilesInputPos` for arrow/Home/End cursor movement in the prompt.

Vault disable+setup mints a new random data key. Same PIN afterward unlocks a different key — existing `.md.enc` files become unreadable without re-wrap. `disablevault` now refuses when non-`z-test-` encrypted notes exist; sync refuses to apply a different `secret/vault` wrap while user notes are on disk. Harnesses must delete `z-test-*.md.enc` before disable. Recovery: `bash scripts/recover-orphaned-vault-notes.sh` with an old `secret/vault` commit from GitHub. Failed decrypt on open must surface on the Files tab — a blank editor with no message is an integrity failure.

## Sync

Destructive sync needs per-note 404 confirmation — a failed remote list must not become mass-delete. Never push empty over a previously-synced note.

Open-file tracking shipped in slices 1, 3, and 4; residuals remain — see [integrity-audit.md](integrity-audit.md). Save and verify while editing skips only the open note, not the whole reconcile. Do not inject Escape on boot in edit mode; daemon, editor, and client have independent lifetimes.

## CI and patches

One patch file, one target file. Multi-file `git apply --recount` cannot tell where hunks end.

Font CI needs one hard-failing RUN per font with grep assertion — a trailing `|| true` swallows download failures.

`int(Uint32) % N` overflows 32-bit int on device; modulo in uint32 space first.

## Recon on BusyBox

This `od` is a stub — pull raw bytes to the Mac. No `timeout` — use `dd & sleep & kill`.
