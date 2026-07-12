# third_party/keywriter/dist/

Built artifacts — QML patches from `build-keywriter.sh` are **baked into** `Writerdeck`.

| File | Description |
|---|---|
| `Writerdeck` | ARM editor binary (patched remarkable-keywriter) |
| `qt5.tar.gz` | Qt5 runtime sysroot subset (libs + QML modules + plugins) |

Deployed to `/home/root/Writerdeck` on the tablet.

**CI:**

```bash
git push                               # triggers build-keywriter.yml
bash scripts/fetch-keywriter-dist.sh   # pull artifact to dist/
bash scripts/deploy-keywriter.sh -b
```

**Local rebuild** (when CI is slow or uncommitted):

```bash
docker build --platform linux/amd64 -t rm1-writerdeck-keywriter-builder third_party/keywriter/
docker run --rm --platform linux/amd64 -v "$PWD/third_party/keywriter/dist:/out" rm1-writerdeck-keywriter-builder
bash scripts/deploy-keywriter.sh -b
```

After deploy, relaunch the editor and confirm `journalctl -u writerdeck` shows framebuffer init — not `Expected token '}'`.
