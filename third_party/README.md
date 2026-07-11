# third_party/

Code/binaries from other projects. Not authored here; each carries its own license.

## keywriter (Component B — the editor, deployed as Writerdeck)
- Source: https://github.com/dps/remarkable-keywriter
- We use it as the on-device Markdown editor: patch `main.cpp` to inject synthetic `QKeyEvent`s read from our local socket (it takes input via Qt QPA, not a device `open()` — no fd to swap). Socket commands include `setfont`, `setrotation`, and editor→server `rotation` acks; `rotation_watcher.{h,cpp}` relays QML `rotationChanged` for USB persistence.
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
    dist/               ← CI-built Writerdeck + qt5.tar.gz (git bridge)
```
