# third_party/keywriter/dist/

CI-built artifacts, produced by `.github/workflows/build-keywriter.yml`.

| File | Description |
|---|---|
| `Writerdeck` | ARM editor binary (patched remarkable-keywriter) |
| `qt5.tar.gz` | Qt5 runtime sysroot subset (libs + QML modules + plugins) |

Deployed to `/home/root/Writerdeck` on the tablet.

```bash
bash scripts/fetch-keywriter-dist.sh   # if not built locally
bash scripts/deploy-keywriter.sh -b
```
