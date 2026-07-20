# Lessons learned

Things that bit us and still matter. Why we chose paths: [decisions.md](decisions.md). Words: [terms.md](terms.md).

## Calling work done vs debugging

While fixing, use the smallest check that could prove you wrong. Before calling work done, run the full verify path in the project rules.

For keyboard work: one run to list all failures, fix several at once, then one editor deploy. Do not push and redeploy per failing case. Do not run edit-session and the typing tests at the same time — they fight over the editor.

After about twenty minutes without a clear finish, stop and report what you tried.

## Typing tests — strategy not proven

Standing decision: [decisions.md](decisions.md) **Typing-test strategy is failing**. Theory: [editor-testing/methodology-shortcomings.md](editor-testing/methodology-shortcomings.md). Claims inventory: [editor-testing/basic-claims.md](editor-testing/basic-claims.md). Do not treat critical-or-full green as “basic editing works,” and do not delete this lesson (or that decisions section, or the banner in [editor-testing/todo.md](editor-testing/todo.md)) until there is solid proof those misses have stopped recurring.

## Deploy traps

Three ways a change looks like it did nothing: GitHub has not built the new editor yet; the phone cached the page; the old editor process is still running.

`rmkw` pushes the binary only. Font changes need the full Qt tree forced onto the tablet. After deploying the server, start the service again before checking the phone UI. Lobby chrome edits go in `lobby-ui.json` on the tablet — do not expect a binary-only deploy to clobber that file (seed only when missing).

`scp` stalls on this link — use gzip over SSH. Never `pkill -f /home/root/Writerdeck`; that also kills the server. Use `pidof Writerdeck`.

If SSH times out, check whether Mac and tablet are on an iPhone hotspot instead of home Wi-Fi.

Product version is auto (`scripts/product-version.sh`). Do not hand-edit `VERSION` and expect About to stay honest — CI and `deploy-rmkbd.sh` own that file. After changing Lobby fragments in the fork, run `./assemble-qml.sh` so committed `main.qml` matches; editing only `main.qml` and then assembling from stale `.inc` files wiped About once. Do not stamp the editor with `-DPRODUCT_VERSION=YYYY-MM-DD` — the ARM toolchain treats the dashes as subtraction; write `product_version.h` instead ([decisions.md](decisions.md) §38).

## Saves and the screen

Every save path must copy the on-screen text into the note before writing the file. Clearing the screen box without syncing it back on load can save a zero-length note. Preview must never feed fancy HTML back into the note body.

Clear the open filename whenever you return to the Lobby, or a deleted note can come back.

## Keys and Home

Edit-mode keys from the socket must go through the QML router on the thread that feeds keys in. Raw Qt events dropped keys or deadlocked. Block Ctrl/Alt navigation key-releases — Qt’s defaults could wipe the screen text while the file on disk stayed fine. Escape toggles edit/preview on key-up; socket inject does not auto-release Escape (that double-fired the harness), so the phone must send an explicit Escape release.

Physical Home is taken over by the server while Writerdeck is open. Do not grab the whole button device from the USB launcher — that starves Home and Power. After Home from edit, Lobby focus must actually handle keys; Lobby now re-focuses after touch and pages the file list instead of flicking ([decisions.md](decisions.md) §35). Remember the open filename before the sync save on Home: that XHR can re-enter the event loop and deliver a noteslist early; keep the prefer-name until select succeeds. Do not bind a scrolling ListView.currentIndex to the selection property (model clear breaks the binding) — the Files list is a fixed page of rows keyed off lobbyFilesIndex, with page start `floor(index / pageSize) * pageSize`. Do not invent a sliding one-row window; that fights e-ink. Lobby chords run on key-press (not release) so phone inject’s press+auto-release cannot double-step the list.

Phone key capture must register listeners once (`initConnection` only). Duplicate `keydown` handlers forward every arrow twice and the Files selection jumps two rows per press.

Alt+Left/Right on USB looked like Escape until the keymap was fixed. qmap (USB keyboard map) changes apply at editor launch, not mid-session. The automated typing tests do not exercise USB layouts — check those by hand.

Wrapped Up/Down must walk visual rows, not step by a tall caret rectangle. Shift selection needs a remembered anchor and head; do not trust the caret index after select. Page buttons are not arrow keys.

After soft-wrap End, set affinity before remembering goal X — Qt’s rect at the wrap index is the next row’s left edge, so Down looked like a no-op and fell through to paragraph end. Down from that End must also snap to the next visual row’s exclusive end, not the last glyph.

## Lobby dialogs

A confirmation or other dialog must read as one piece — title, body, and actions together — not a prompt above the list and buttons far below it. Scattered chrome blends into the note list (especially when type size and weight match list rows) and people miss the question. At minimum put a clear divider between dialog and the rest of the UI; a floating black-on-white box is fine and expected. Do not copy the inherited Ctrl-K note picker (black panel, white type) for Lobby confirms. Shared chrome lives in `lobby/lobby_dialog.inc` (scrim + white box) so later changes (for example letting the list show through the scrim) apply to every Lobby dialog.

Lobby look, wording, and Ctrl-letter chords live in `/home/root/.Writerdeck/lobby-ui.json`. Writerdeck reloads that file when it changes (file/directory watch plus a short mtime poll — rename replaces often miss inotify on the tablet). Start that poll only after QApplication exists; a timer started from a static constructed before the app never fires. Bad or unreadable JSON keeps the last good load (or the baked-in defaults if nothing has loaded yet) — fix the file rather than restarting for a blank Lobby.

Name collisions on New / Rename / New encrypted stay in that dialog — keep it open with a short line under the typed name; clear the line when the user edits. Do not close the dialog and park the same message in the Files header. The header box is for other Files feedback (failed decrypt, Download without a phone page, and the like). The server still refuses the duplicate and can push `vaultopfailed` so a race reopens the dialog with the same sentence. Uniqueness is case-insensitive and shared across plain and encrypted forms of the same title; the Files list sorts that way too.

On the Files tab, use a fixed stack: header (feedback), list, footer (Prev / Page N/M / Next when notes spill a page, a hairline, then the action bars). The list only fills between header and footer. Do not lay feedback, list, and bars out in one Column while hand-computing the list height, and do not paint the list above the footer — either mistake pushes buttons off-screen or draws rows on top of them. Gray 9pt page text is invisible on e-ink; use black labels and a separator so the strip does not blend into New/Edit.

Typing actions from a touch tap (edit, new, rename, new encrypted) show a Connect-a-keyboard tip when neither a USB keyboard nor an open phone/laptop page is present. An open page counts as a keyboard path (Type field). Keyboard chords skip the tip. Continue, or any key while the tip is up, runs that one pending action once — never a sticky ready flag. Dead WebSocket clients kept phoneConnected true without an intentional page; hello plus ping/pong fixed that. Cursor’s embedded browser also sent hello and skipped the tip — the phone page skips hello there, and the server ignores hello from User-Agents containing Cursor/ or Electron/ ([decisions.md](decisions.md) §34).


## Phone page

Stand down key capture when PIN, paste, or sync overlays are up. The GitHub token is per browser address — a new tablet IP means enter it once for that address. After server restart, watch the journal for `sync: nothing to do (token)` (or a push/pull line) after the client connects — not a burst of duplicate token reconciles. Do not poll-check editor disk hash on the status timer — tablet autosave would false-alarm “Disk changed”; rely on the WebSocket diskchanged path for real external writes.

## Sync and vault

Never mass-delete from a failed remote list. Never push empty over a previously synced note. Do not PUT a note or vault secret whose local fingerprint already matches `syncMeta` — that was filling GitHub with empty commits every few minutes when a timer still ran. Sync is event-driven now; a clean tablet logs at most one `sync: nothing to do` per quiet streak, then a ×N summary when something finally changes — not a no-op commit. Re-POSTing the same GitHub token from extra browser tabs must not each start a reconcile. Disabling the vault while encrypted notes exist orphans them — refuse that; recover from GitHub history if it already happened. A failed decrypt must show an error on Files, not a blank editor. A wrong vault PIN must keep the pad open with a short message (for example Wrong PIN. Try again.), not dismiss silently.

`journalctl -u writerdeck` is Writerdeck’s slice only (tens of KB on a quiet device). Almost all of `/var/log` is the stock reMarkable journal — do not treat multi-megabyte disk use as Writerdeck log spam.

## Build

Screen-file parse errors look like “editor exited at once.” Balance braces before deploy. Calling caret geometry from C++ needs the right argument types — wrong types silently return zeros and send Up/Down to the ends of the document.
