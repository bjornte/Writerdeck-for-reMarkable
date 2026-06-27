# third_party/keywriter/dist/

CI-built artifacts, produced by `.github/workflows/build-keywriter.yml`.

In this public mirror the binaries are NOT committed -- clone and run the
workflow (or build locally) to produce them.

| File | Description |
|---|---|
| `keywriter` | ARM binary cross-built from source (`dps/remarkable-keywriter`) |
| `qt5.tar.gz` | Qt5 runtime sysroot subset (libs + QML modules + plugins) |

## Build

Push a change to the keywriter sources (CI runs `build-keywriter.sh` in the
toltec Qt container `ghcr.io/toltec-dev/qt:v3.3`), or run that script locally
inside the same container.

## Deployment

```bash
bash scripts/deploy-keywriter.sh
```

Unpacks `qt5.tar.gz` transiently, copies everything to the device, and runs
the alive-8s launch check. The unpacked `qt5/` directory is removed after the
script exits.