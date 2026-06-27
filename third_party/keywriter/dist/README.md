# third_party/keywriter/dist/

CI-built artifacts, populated by `.github/workflows/build-keywriter.yml`.
Committed to the repo so the Mac can pull and deploy without a local toolchain.

| File | Description |
|---|---|
| `keywriter` | ARM binary cross-built from source (`dps/remarkable-keywriter`) |
| `qt5.tar.gz` | Qt5 runtime sysroot subset (libs + QML modules + plugins) |

## Deployment

```bash
bash scripts/deploy-keywriter.sh
```

This unpacks `qt5.tar.gz` transiently, copies everything to the device, and
runs the alive-8s launch check.  The unpacked `qt5/` directory is deleted after
the script exits.

## Not committed (gitignored)

- `qt5/` — the locally-unpacked sysroot tree (transient; created/removed by
  `deploy-keywriter.sh`).
