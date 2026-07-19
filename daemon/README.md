# daemon/ — Writerdeck-server

Always-on Go program on the tablet. Deploys as `/home/root/Writerdeck-server`. How it fits: [architecture.md](../docs/architecture.md).

Feeds the editor over `/run/Writerdeck.sock` — this kernel has no usable uinput (fake keyboard device).

Serves the phone page and WebSocket on port 8000. Handles notes, settings, PIN, Lobby launches, GitHub sync, and device tests under `/api/test/`. Phone UI is the keyboard shell (paste, sync, download offer) — [browser-vs-tablet.md](../docs/browser-vs-tablet.md). WebSocket `hello` marks a real browser as a keyboard path for the Lobby tip; Cursor/Electron User-Agents are excluded ([decisions.md](../docs/decisions.md) §34).

Build and ship: `bash ../scripts/deploy-rmkbd.sh`.

If the stock UI flashes right after open, the editor probably crashed — read the journal and run `test-edit-session.sh`. Do not `pkill -f /home/root/Writerdeck`; that hits the server too.
