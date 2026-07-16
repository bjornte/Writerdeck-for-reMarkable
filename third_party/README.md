# third_party/

Code/binaries from other projects. Not authored here; each carries its own license.

## keywriter (Component B — the editor engine, deployed as Writerdeck)

Upstream [remarkable-keywriter](https://github.com/dps/remarkable-keywriter): a **Qt 5** editor written in **C++** and **QML**. Writerdeck is our patched build of that engine (socket input, Lobby, Mac/Linux-style editing with Ctrl/Alt chords, and related behavior). Owned fork for migration: [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter).

**C++ vs QML:** QML = screen + most typing/selection wiring (`main.qml`, fork `edit_mac_helpers.qml.inc`). C++ = startup, e-ink/display, socket keystroke inject (`main.cpp`), plus (migration 2) pure text math / undo in `EditHelper`. See [../docs/editor-migration-2-to-cpp/](../docs/editor-migration-2-to-cpp/).

- Socket inject, `lobby_bridge`, `rotation_watcher`, Lobby/shell `main.qml`, and `lobby/*.inc` live **in the fork**. CI asserts they are present; it only inserts `edit_mac_helpers.qml.inc` and concatenates Lobby subpages + sleep screen.
- Built from source — cross-compiled in `ghcr.io/toltec-dev/qt:v3.3` via **CI** ([build-keywriter.sh](keywriter/build-keywriter.sh) + [Dockerfile](keywriter/Dockerfile), workflow `build-keywriter.yml`). Mac: `git push` → `fetch-keywriter-dist.sh` → `deploy-keywriter.sh -b` — not local `docker build`.
- [../scripts/deploy-keywriter.sh](../scripts/deploy-keywriter.sh) ships the binary to `/home/root/Writerdeck`; notes live in the separate directory `/home/root/Writerdeck-user-documents/` — don't put the binary there.

```
third_party/
  keywriter/
    build-keywriter.sh  ← CI build glue (assert fork + helpers/lobby concat + qmake)
    Dockerfile          ← toltec qt:v3.3 build image
    dist/               ← CI-built Writerdeck + qt5.tar.gz (fetch via fetch-keywriter-dist.sh)
```
