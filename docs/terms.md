# Terms

Short definitions. Fork twin: [Writerdeck-keywriter docs/terms.md](https://github.com/bjornte/Writerdeck-keywriter/blob/master/docs/terms.md).

## Product

* **Writerdeck**:  A Markdown typewriter on a first-gen reMarkable for USB and Bluetooth keyboards. Bluetooth pairs to the phone; typing reaches the tablet over Wi-Fi. USB plugs in with an OTG cable.
* **Writerdeck-server**: Always-on program on the tablet: files, sync, PIN, launching the editor. No screen of its own.
* **Writerdeck (the app)**: The full-screen editor you see. Built from our fork of Dave’s keywriter.
* **Lobby**: In-app home on the tablet — files, settings, sleep — not the stock reMarkable UI.
* **Document integrity**: Your prose must survive as plain Markdown on disk.

## Editor

* **QML**: Screen language — layout and applying edits on screen (the main screen file is `main.qml`).
* **C++ / EditHelper**: Startup, display, socket keys, and the math behind shortcuts, wrap, and undo.
* **Shortcut (chord)**: A key held with Ctrl, Alt, or Shift.
* **Visual line**: One wrapped row on screen — not the same as a line ending in a newline.
* **Goal column**: The horizontal spot Up/Down tries to keep across uneven lines.

## Testing and ops

* **Automated typing tests**: Scripted typing on the real tablet over the same path the phone uses (`test-keyboard-harness.sh`).
* **Basic set / full set**: Thirty-eight checks for “basic editing works”; one hundred ten before calling typing work done.
* **Edit-session check**: Opens a note and checks the editor stays up — catches broken QML that crashes on launch.
* **Deploy**: Copy a new binary to the tablet and relaunch the editor. Restarting the server alone does not reload the editor.
* **OTA (over-the-air update)**: Tablet software update from reMarkable — may reset the SSH password and wipe our boot service.
