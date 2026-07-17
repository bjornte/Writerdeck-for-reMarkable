# Terms

Short definitions for Writerdeck and its editor fork.

## Product

* **Writerdeck for reMarkable**: A typewriter for reMarkable 1 with a Bluetooth or USB keyboard. Bluetooth pairs via the user's phone. USB connects with an OTG cable. Made primarily in Go (server), QML & C++ (editor), and HTML/CSS/JavaScript (phone page).
* **Writerdeck-server**: Always-on program on the tablet: files, sync, PIN, launching the editor. No screen of its own. Made primarily in Go. The phone page it serves is HTML, CSS, and JavaScript.
* **Writerdeck (app)**: The full-screen editor you see. Built from our fork of Singleton’s keywriter: [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter). Made primarily in QML & C++.
* **keywriter**: Singleton’s original Qt Markdown notepad — the project that fork started from. Made primarily in QML & C++. Uses Qt’s TextEdit.
* **Lobby**: In-app home on the tablet — files, settings, sleep — not the stock reMarkable UI. Made primarily in QML.
* **Document integrity**: Your prose must survive as plain Markdown on disk.

## Editor

* **QML (the screen file)**: Screen language — layout and applying edits on screen (the main file is `main.qml`).
* **C++ / EditHelper**: Startup, display, socket keys, and the math behind shortcuts, wrap, and undo. Made primarily in C++.
* **TextEdit (Qt’s on-screen text box)**: Fine for drawing; weak for “which wrapped row am I on?” Shipped with Qt (C++); not ours.
* **Shortcut (chord)**: A key held with Ctrl, Alt, or Shift.
* **Visual line**: One wrapped row on screen — not the same as a line ending in a newline.
* **Goal column**: The horizontal spot Up/Down tries to keep across uneven lines.
* **assemble-qml.sh**: In the fork: builds committed `main.qml` from modular pieces. Run it after changing helpers or Lobby; CI does not stitch QML. Made primarily in shell.
* **uinput (fake keyboard device)**: Not used here — keys arrive on a unix socket instead. A Linux feature, not something we wrote.

## Testing and ops

* **Automated typing tests (`test-keyboard-harness.sh`)**: Scripted typing on the real tablet. Reaching the editor using the WebSocket, the mechanism used for Bluetooth keyboard typing. Made primarily in Go (the scenarios) and shell (the wrapper).
* **Basic set / full set**: Thirty-eight checks for “basic editing works”; one hundred ten before calling typing work done.
* **Edit-session check**: Opens a note and checks the editor stays up — catches broken QML that crashes on launch. Made primarily in shell.
* **Known-good commit**: A Writerdeck-keywriter revision that last passed those typing tests. Everyday builds usually follow `master`.
* **Deploy**: Copy a new binary to the tablet and relaunch the editor. Restarting the server alone does not reload the editor. Done with shell scripts.
* **OTA (over-the-air update)**: Tablet software update from reMarkable — may reset the SSH password and wipe our boot service. From reMarkable, not us.
* **The original / shared history**: Singleton’s remarkable-keywriter repo, and the shared git starting point that makes ordinary merges possible. Developers often nickname Singleton’s remote `upstream`; it still means the original. Made primarily in QML & C++ (same as keywriter).
