# third_party/keywriter/dist/

CI-built artifacts, produced by `.github/workflows/build-keywriter.yml`.

| File | Description |
|---|---|
| `Writerdeck` | ARM editor binary (patched remarkable-keywriter) |
| `qt5.tar.gz` | Qt5 runtime sysroot subset (libs + QML modules + plugins) |

Deployed to `/home/root/Writerdeck` on the tablet.

```bash
git push                               # CI builds via Docker in GHA
bash scripts/fetch-keywriter-dist.sh   # pull artifact to dist/
bash scripts/deploy-keywriter.sh -b
```
