# third_party/

Code/binaries from other projects. Not authored here; each carries its own license.

## keywriter (Component B — the editor engine, deployed as Writerdeck)

Upstream [remarkable-keywriter](https://github.com/dps/remarkable-keywriter): a **Qt 5** editor written in **C++** and **QML**. Writerdeck is our patched build of that engine (socket input, Lobby, Mac/Linux-style editing with Ctrl/Alt chords, and related behavior). Owned fork for migration: [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter).

**C++ vs QML:** QML = screen + typing/selection behavior (`main.qml`, fork `edit_mac_helpers.qml.inc`). C++ = startup, e-ink/display, socket keystroke inject (`main.cpp`). Editing improvements for the current migration almost always mean QML in the fork.

- Socket inject, `lobby_bridge`, and `rotation_watcher` live **in the fork** (`main.cpp` + helpers). CI asserts they are present; it no longer `git apply`s a patch or copies those sources from this tree.
- Built from source — cross-compiled in `ghcr.io/toltec-dev/qt:v3.3` via **CI** ([build-keywriter.sh](keywriter/build-keywriter.sh) + [Dockerfile](keywriter/Dockerfile), workflow `build-keywriter.yml`). Mac: `git push` → `fetch-keywriter-dist.sh` → `deploy-keywriter.sh -b` — not local `docker build`.
- [../scripts/deploy-keywriter.sh](../scripts/deploy-keywriter.sh) ships the binary to `/home/root/Writerdeck`; notes live in the separate directory `/home/root/Writerdeck-user-documents/` — don't put the binary there.

```
third_party/
  keywriter/
    build-keywriter.sh  ← CI build (assert fork C++ + QML Python patches + qmake)
    Dockerfile          ← toltec qt:v3.3 build image
    lobby/              ← Lobby QML fragments (still injected at build time)
    dist/               ← CI-built Writerdeck + qt5.tar.gz (fetch via fetch-keywriter-dist.sh)
```
