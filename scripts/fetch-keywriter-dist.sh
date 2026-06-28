#!/usr/bin/env bash
# scripts/fetch-keywriter-dist.sh -- pull the CI-built keywriter artifacts
# (keywriter binary + qt5.tar.gz) into third_party/keywriter/dist/ on the Mac.
#
# Why this exists: the public mirror ships source-only -- dist/keywriter and
# dist/qt5.tar.gz are gitignored, so `git pull` never brings them. CI builds
# them and uploads them as the "keywriter-dist" artifact instead. This script
# downloads that artifact so deploy-keywriter.sh can find it. Requires gh
# (brew install gh) authenticated as bjornte.
#
# Usage (run from repo root on the Mac):
#   bash scripts/fetch-keywriter-dist.sh            # latest successful run
#   bash scripts/fetch-keywriter-dist.sh <run-id>   # a specific run
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
DEST="$REPO/third_party/keywriter/dist"

mkdir -p "$DEST"
# Download to a temp dir first: the artifact carries a README.md that would
# collide with dist/README.md and abort gh's extraction. Copy only the two
# artifacts we need (keywriter + qt5.tar.gz) over the top.
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
if [ -n "${1:-}" ]; then
    gh run download "$1" -n keywriter-dist -D "$TMP"
else
    gh run download -n keywriter-dist -D "$TMP"
fi
cp -f "$TMP/keywriter" "$TMP/qt5.tar.gz" "$DEST/"

echo "Fetched into $DEST:"
ls -lh "$DEST/keywriter" "$DEST/qt5.tar.gz"
echo "Next: bash scripts/deploy-keywriter.sh -b"
