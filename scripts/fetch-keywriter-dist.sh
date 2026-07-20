#!/usr/bin/env bash
# scripts/fetch-keywriter-dist.sh -- pull CI-built keywriter artifacts into dist/.
# Prefers the rolling GitHub Release (curl, no login). Falls back to gh run download.
#
# Usage (run from repo root on the Mac):
#   bash scripts/fetch-keywriter-dist.sh            # release, else latest successful run
#   bash scripts/fetch-keywriter-dist.sh <run-id>   # a specific Actions run (needs gh)
#
# Override repo/tag if needed:
#   KEYWRITER_DIST_REPO=owner/repo KEYWRITER_DIST_TAG=keywriter bash scripts/fetch-keywriter-dist.sh
set -euo pipefail
if [ -t 2 ]; then R=$'\033[1;31m'; Z=$'\033[0m'; else R=''; Z=''; fi
err() { printf '%sERROR:%s %s\n' "$R" "$Z" "$*" >&2; }
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
DEST="$REPO/third_party/keywriter/dist"
REPO_SLUG="${KEYWRITER_DIST_REPO:-bjornte/Writerdeck-for-reMarkable}"
RELEASE_TAG="${KEYWRITER_DIST_TAG:-keywriter}"
RELEASE_BASE="https://github.com/${REPO_SLUG}/releases/download/${RELEASE_TAG}"

mkdir -p "$DEST"
# Download to a temp dir first: the artifact carries a README.md that would
# collide with dist/README.md and abort gh's extraction. Copy only the two
# artifacts we need (Writerdeck + qt5.tar.gz) over the top.
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

have_artifacts() {
    [ -f "$1/Writerdeck" ] && [ -f "$1/qt5.tar.gz" ]
}

fetch_release() {
    echo "Trying GitHub Release ${REPO_SLUG}@${RELEASE_TAG} ..."
    curl -fsSL -o "$TMP/Writerdeck" "${RELEASE_BASE}/Writerdeck" || return 1
    curl -fsSL -o "$TMP/qt5.tar.gz" "${RELEASE_BASE}/qt5.tar.gz" || return 1
    chmod +x "$TMP/Writerdeck"
    have_artifacts "$TMP"
}

fetch_gh_run() {
    if ! command -v gh >/dev/null; then
        return 1
    fi
    echo "Falling back to gh run download ..."
    if [ -n "${1:-}" ]; then
        gh run download "$1" -R "$REPO_SLUG" -n keywriter-dist -D "$TMP"
    else
        gh run download -R "$REPO_SLUG" -n keywriter-dist -D "$TMP"
    fi
    have_artifacts "$TMP"
}

print_empty_help() {
    err "could not populate dist/ (Writerdeck + qt5.tar.gz missing)"
    cat >&2 <<EOF
  Tried: ${RELEASE_BASE}/
  Also tried: gh run download (needs: brew install gh && gh auth login)

  Browser fallback (no gh):
    1. https://github.com/${REPO_SLUG}/releases/tag/${RELEASE_TAG}
       -- or Actions -> Build keywriter -> latest green run -> keywriter-dist artifact
    2. Save Writerdeck and qt5.tar.gz into third_party/keywriter/dist/
EOF
}

RUN_ID="${1:-}"
ok=0
if [ -n "$RUN_ID" ]; then
    # Explicit run id always uses Actions (release is rolling-latest only).
    if fetch_gh_run "$RUN_ID"; then
        ok=1
    fi
else
    if fetch_release; then
        ok=1
    elif fetch_gh_run; then
        ok=1
    fi
fi

if [ "$ok" -ne 1 ] || ! have_artifacts "$TMP"; then
    print_empty_help
    exit 1
fi

cp -f "$TMP/Writerdeck" "$TMP/qt5.tar.gz" "$DEST/"
chmod +x "$DEST/Writerdeck"

if ! have_artifacts "$DEST"; then
    print_empty_help
    exit 1
fi

echo "Fetched into $DEST:"
ls -lh "$DEST/Writerdeck" "$DEST/qt5.tar.gz"
echo "Next: bash scripts/deploy-keywriter.sh -b"
