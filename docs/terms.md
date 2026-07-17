# Terms

Short definitions for Writerdeck and its editor fork.

## Product

Writerdeck. A Markdown typewriter on a first-gen reMarkable for USB and Bluetooth keyboards. Bluetooth pairs to the phone; typing reaches the tablet over Wi-Fi. USB plugs in with an OTG cable.

Writerdeck-server. Always-on program on the tablet: files, sync, PIN, launching the editor. No screen of its own.

Writerdeck (the app). The full-screen editor you see. Built from our fork of Dave’s keywriter: [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter).

keywriter. Dave Singleton’s original Qt Markdown notepad — the project that fork started from.

Lobby. In-app home on the tablet — files, settings, sleep — not the stock reMarkable UI.

Document integrity. Your prose must survive as plain Markdown on disk.

## Editor

QML (the screen file). Screen language — layout and applying edits on screen (the main file is `main.qml`).

C++ / EditHelper. Startup, display, socket keys, and the math behind shortcuts, wrap, and undo.

TextEdit (Qt’s on-screen text box). Fine for drawing; weak for “which wrapped row am I on?”

Shortcut (chord). A key held with Ctrl, Alt, or Shift.

Visual line. One wrapped row on screen — not the same as a line ending in a newline.

Goal column. The horizontal spot Up/Down tries to keep across uneven lines.

assemble-qml.sh. In the fork: builds committed `main.qml` from modular pieces. Run it after changing helpers or Lobby; CI does not stitch QML.

uinput (fake keyboard device). Not used here — keys arrive on a unix socket instead.

## Testing and ops

Automated typing tests (`test-keyboard-harness.sh`). Scripted typing on the real tablet over the same path the phone uses.

Basic set / full set. Thirty-eight checks for “basic editing works”; one hundred ten before calling typing work done.

Edit-session check. Opens a note and checks the editor stays up — catches broken QML that crashes on launch.

Known-good commit. A Writerdeck-keywriter revision that last passed those typing tests. Everyday builds usually follow `master`.

Deploy. Copy a new binary to the tablet and relaunch the editor. Restarting the server alone does not reload the editor.

OTA (over-the-air update). Tablet software update from reMarkable — may reset the SSH password and wipe our boot service.

The original / shared history. Dave’s remarkable-keywriter repo, and the shared git starting point that makes ordinary merges possible. Developers often nickname Dave’s remote `upstream`; it still means the original.
