#!/usr/bin/env bash
# scripts/fetch-keywriter-dist.sh -- pull CI-built keywriter artifacts into dist/.
# Mac deploy path after third_party/keywriter/ changes (Docker runs in GHA, not locally).
# Requires gh (brew install gh && gh auth login).
#
# Usage (run from repo root on the Mac):
#   bash scripts/fetch-keywriter-dist.sh            # latest successful run
#   bash scripts/fetch-keywriter-dist.sh <run-id>   # a specific run
set -euo pipefail
if [ -t 2 ]; then R=$'\033[1;31m'; Z=$'\033[0m'; else R=''; Z=''; fi
err() { printf '%sERROR:%s %s\n' "$R" "$Z" "$*" >&2; }
command -v gh >/dev/null || { err "gh not found -- run: brew install gh && gh auth login"; exit 1; }
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
cp -f "$TMP/Writerdeck" "$TMP/qt5.tar.gz" "$DEST/"

echo "Fetched into $DEST:"
ls -lh "$DEST/Writerdeck" "$DEST/qt5.tar.gz"
echo "Next: bash scripts/deploy-keywriter.sh -b"
