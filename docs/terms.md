# Terms

Short definitions. Fork glossary: [Writerdeck-keywriter docs/terms.md](https://github.com/bjornte/Writerdeck-keywriter/blob/master/docs/terms.md).

## Product

Writerdeck. A Markdown typewriter on a first-gen reMarkable for USB and Bluetooth keyboards. Bluetooth pairs to the phone; typing reaches the tablet over Wi-Fi. USB plugs in with an OTG cable.

Writerdeck-server. Always-on program on the tablet: files, sync, PIN, launching the editor. No screen of its own.

Writerdeck (the app). The full-screen editor you see. Built from our fork of Dave’s keywriter.

Lobby. In-app home on the tablet — files, settings, sleep — not the stock reMarkable UI.

Document integrity. Your prose must survive as plain Markdown on disk.

## Editor

QML. What you see and how edits are applied on screen.

C++ / EditHelper. Startup, display, socket keys, and the math behind chords, wrap, and undo.

Chord. A shortcut that holds Ctrl, Alt, or Shift with another key.

Visual line. One wrapped row on screen — not the same as a line ending in a newline.

Goal column. The horizontal spot Up/Down tries to keep across uneven lines.

## Testing and ops

Keyboard harness. Scripted typing on the real tablet over the same path the phone uses.

Critical / full suite. Thirty-eight scenarios for “basic editing works”; one hundred ten for product sign-off.

Edit-session test. Opens a note and checks the editor stays up — catches broken QML.

Deploy. Copy a new binary to the tablet and relaunch the editor. Restarting the server alone does not reload the editor.
