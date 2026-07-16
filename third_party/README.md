# third_party/

Code/binaries from other projects. Not authored here; each carries its own license.

## keywriter (Component B — the editor engine, deployed as Writerdeck)

Upstream [remarkable-keywriter](https://github.com/dps/remarkable-keywriter): a **Qt 5** editor written in **C++** and **QML**. Writerdeck is our patched build of that engine (socket input, Lobby, Mac/Linux-style editing with Ctrl/Alt chords, and related behavior). Owned fork for migration: [bjornte/Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter).

- We patch `main.cpp` to inject synthetic `QKeyEvent`s from our local socket (input goes through Qt QPA, not a device `open()` — no fd to swap). Socket commands include `setfont`, `setrotation`, and editor→server `rotation` acks; `rotation_watcher.{h,cpp}` relays QML `rotationChanged` for USB persistence.
- Built from source — cross-compiled in `ghcr.io/toltec-dev/qt:v3.3` via **CI** ([build-keywriter.sh](keywriter/build-keywriter.sh) + [Dockerfile](keywriter/Dockerfile), workflow `build-keywriter.yml`). Mac: `git push` → `fetch-keywriter-dist.sh` → `deploy-keywriter.sh -b` — not local `docker build`.
- [../scripts/deploy-keywriter.sh](../scripts/deploy-keywriter.sh) ships the binary to `/home/root/Writerdeck`; notes live in the separate directory `/home/root/Writerdeck-user-documents/` — don't put the binary there.

```
third_party/
  keywriter/
    build-keywriter.sh  ← CI build (qmake + make in the toltec image)
    Dockerfile          ← toltec qt:v3.3 build image
    rotation_watcher.h  ← moc'd rotationChanged → server notify
    rotation_watcher.cpp
    socket-inject.patch ← main.cpp socket reader + setrotation
    dist/               ← CI-built Writerdeck + qt5.tar.gz (fetch via fetch-keywriter-dist.sh)
```
