# third_party/

Upstream bits we do not author. Each keeps its own license.

## keywriter → Writerdeck

Dave Singleton’s [remarkable-keywriter](https://github.com/dps/remarkable-keywriter), maintained as [Writerdeck-keywriter](https://github.com/bjornte/Writerdeck-keywriter). CI builds the tablet binary here via `keywriter/build-keywriter.sh` (clone, assert, compile — no QML stitching). Fetch with `fetch-keywriter-dist.sh`, deploy with `deploy-keywriter.sh`. Put the binary in `/home/root/Writerdeck`, not in the notes folder.

Editor behavior and Lobby live in the fork. Policy: [docs/decisions.md](../docs/decisions.md) §4–§6.
