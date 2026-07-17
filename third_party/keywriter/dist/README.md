# third_party/keywriter/dist/

Built artifacts — fork QML is baked into `Writerdeck` at CI build time (no runtime patching).

| File | Description |
|---|---|
| `Writerdeck` | ARM editor binary (patched remarkable-keywriter) |
| `qt5.tar.gz` | Qt5 runtime sysroot subset (libs + QML modules + plugins) |

Deployed to `/home/root/Writerdeck` on the tablet.

**Fetch (no `gh` required):** rolling Release tag `keywriter`, or Actions artifact as browser fallback.

```bash
bash scripts/fetch-keywriter-dist.sh   # curl Release, else gh run download
bash scripts/deploy-keywriter.sh -b
```

After a fork/CI change: wait for Build keywriter on main (updates the Release), then fetch + deploy.

**Local rebuild** (when CI is slow or uncommitted):

```bash
docker build --platform linux/amd64 -t rm1-writerdeck-keywriter-builder third_party/keywriter/
docker run --rm --platform linux/amd64 -v "$PWD/third_party/keywriter/dist:/out" rm1-writerdeck-keywriter-builder
bash scripts/deploy-keywriter.sh -b
```

After deploy, relaunch the editor and confirm `journalctl -u writerdeck` shows framebuffer init — not `Expected token '}'`.
