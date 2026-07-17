# Lessons learned

Things that bit us and still matter. Why we chose paths: [decisions.md](decisions.md). Words: [terms.md](terms.md).

## Sign-off vs debugging

While fixing, use the cheapest check that could prove you wrong. Before calling work done, run the full verify path in the project rules.

For keyboard work: one triage run, batch the fixes, one editor deploy. Do not push and redeploy per scenario. Do not run edit-session and the keyboard harness at the same time — they fight over the editor.

After about twenty minutes without sign-off, stop and report what you tried.

## Deploy traps

Three ways a change looks like it did nothing: CI has not built the new editor yet; the phone cached the page; the old editor process is still running.

`rmkw` pushes the binary only. Font changes need the full Qt tree forced onto the tablet. After deploying the server, start the service again before checking the phone UI.

`scp` stalls on this link — use gzip over SSH. Never `pkill -f /home/root/Writerdeck`; that also kills the server. Use `pidof Writerdeck`.

If SSH times out, check whether Mac and tablet are on an iPhone hotspot instead of home Wi-Fi.

## Saves and the screen

Every save path must copy the on-screen text into the note before writing the file. Clearing the screen box without syncing it back on load can save a zero-length note. Preview must never feed fancy HTML back into the note body.

Clear the open filename whenever you return to the Lobby, or a deleted note can come back.

## Keys and Home

Edit-mode keys from the socket must go through the QML router on the inject thread. Raw Qt events dropped keys or deadlocked. Block Ctrl/Alt navigation key-releases — Qt’s defaults could wipe the screen text while the file on disk stayed fine.

Physical Home is grabbed by the server while Writerdeck is open. Do not grab the whole gpio device from the USB launcher — that starves Home and Power. After Home from edit, Lobby focus must actually handle keys.

Alt+Left/Right on USB looked like Escape until the keymap was fixed. Qmap changes apply at editor launch, not mid-session. The automated harness does not exercise USB layouts — check those by hand.

Wrapped Up/Down must walk visual rows, not step by a tall caret rectangle. Shift selection needs a remembered anchor and head; do not trust the caret index after select. Page buttons are not arrow keys.

## Phone page

Stand down key capture when PIN, paste, or sync overlays are up. The GitHub token is per browser origin — a new tablet IP means enter it once for that address. After server restart, watch the journal for token restore after the client connects.

## Sync and vault

Never mass-delete from a failed remote list. Never push empty over a previously synced note. Disabling the vault while encrypted notes exist orphans them — refuse that; recover from GitHub history if it already happened. A failed decrypt must show an error on Files, not a blank editor.

## Build

QML parse errors look like “editor exited at once.” Balance braces before deploy. Calling caret geometry from C++ needs the right argument types — wrong types silently return zeros and send Up/Down to the ends of the document.
