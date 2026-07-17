# Terms

Layman definitions for Writerdeck. Editor-fork glossary (same spirit, engine-focused): [Writerdeck-keywriter docs/terms.md](https://github.com/bjornte/Writerdeck-keywriter/blob/master/docs/terms.md).

## Product

Writerdeck. Wi-Fi Markdown typewriter on a first-gen reMarkable: phone (or USB keyboard) to type; tablet to read and keep notes.

Writerdeck-server (daemon). Always-on Go program on the tablet. Wi-Fi, files, sync, PIN, launches the editor. No screen of its own.

Writerdeck (the app). Full-screen Qt editor on the tablet. Built from the [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter) fork of Dave Singleton’s keywriter.

keywriter / remarkable-keywriter. Upstream Qt Markdown notepad. We ship a maintained fork as Writerdeck.

Lobby. In-app tablet home (files, settings, sleep) — not the stock reMarkable UI.

Phone UI. Browser page served by the daemon for typing and sync token entry. Day-to-day file ops live on the tablet Lobby.

Document integrity. Owner prose must survive as plain UTF-8 Markdown on disk. Gates features; see [integrity-audit.md](integrity-audit.md).

## Editor

QML. Screen and apply layer: layout, Lobby, selection, timers.

C++. Startup, display, socket keys, and EditHelper (math, chords, wrap walk, undo).

EditHelper. C++ brain in the fork. QML TextEdit draws and applies.

Chord. Modifier + key shortcut (Ctrl/Alt/Shift), Mac/Linux style.

Visual line. One wrapped row on screen. Differs from a logical line (newline-separated).

Goal column. Sticky horizontal target for Up/Down across uneven lines.

Custom undo. EditHelper stacks; Qt’s TextEdit undo is sidelined so socket and chords share one history ([decisions.md](decisions.md) §30).

Socket / uinput. Keys use `/run/Writerdeck.sock`. This kernel cannot load uinput ([decisions.md](decisions.md) §1).

## Testing

Keyboard harness. `scripts/test-keyboard-harness.sh` — automated typing/selection on the real tablet over the socket path. Scenarios inspired by CodeMirror / Qt / Ace ([editor-testing/](editor-testing/)).

Critical. 38 scenarios — “basic editing works” gate. Must be **38/38/0** before trusting a keyboard deploy.

Full suite. 110 scenarios — product sign-off (**110/110/0**). Includes wrap and undo tags.

Edit-session. `scripts/test-edit-session.sh` — open a note; editor must stay up (~8s). Catches QML that crashes on launch.

Ship tip. Fork commit that last passed the harness. Documented in the fork README; CI usually tracks `master`.

## Ops

Deploy. Push code, build Writerdeck in CI when the fork changed, fetch the binary, copy to the tablet, relaunch the editor. Server restart alone is not enough after a binary change.

Fork assembly. Helpers and Lobby are assembled in the fork (`assemble-qml.sh` → committed `main.qml`). Writerdeck’s `build-keywriter.sh` clones, asserts, and builds — it does not stitch QML.
