#!/usr/bin/env bash
# scripts/fetch-server-dist.sh -- pull CI-built Writerdeck-server (no go needed).
# Prefers the rolling GitHub Release (curl). Falls back to gh run download.
#
# Usage (run from repo root on the Mac):
#   bash scripts/fetch-server-dist.sh
#   bash scripts/fetch-server-dist.sh <run-id>
set -euo pipefail
if [ -t 2 ]; then R=$'\033[1;31m'; Z=$'\033[0m'; else R=''; Z=''; fi
err() { printf '%sERROR:%s %s\n' "$R" "$Z" "$*" >&2; }
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
DEST="$REPO/Writerdeck-server"
REPO_SLUG="${WRITERDECK_SERVER_REPO:-bjornte/Writerdeck-for-reMarkable}"
RELEASE_TAG="${WRITERDECK_SERVER_TAG:-server}"
RELEASE_URL="https://github.com/${REPO_SLUG}/releases/download/${RELEASE_TAG}/Writerdeck-server"

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
OUT="$TMP/Writerdeck-server"

fetch_release() {
    echo "Trying GitHub Release ${REPO_SLUG}@${RELEASE_TAG} ..."
    curl -fsSL -o "$OUT" "$RELEASE_URL" || return 1
    [ -s "$OUT" ]
}

fetch_gh_run() {
    if ! command -v gh >/dev/null; then
        return 1
    fi
    echo "Falling back to gh run download ..."
    if [ -n "${1:-}" ]; then
        gh run download "$1" -R "$REPO_SLUG" -n writerdeck-server -D "$TMP"
    else
        gh run download -R "$REPO_SLUG" -n writerdeck-server -D "$TMP"
    fi
    # Artifact may land as Writerdeck-server in TMP root
    [ -f "$OUT" ] || [ -f "$TMP/writerdeck-server/Writerdeck-server" ] || return 1
    if [ ! -f "$OUT" ]; then
        OUT="$TMP/writerdeck-server/Writerdeck-server"
    fi
    [ -s "$OUT" ]
}

print_help() {
    err "could not fetch Writerdeck-server"
    cat >&2 <<EOF
  Tried: ${RELEASE_URL}
  Also tried: gh run download (needs: brew install gh && gh auth login)

  Browser fallback:
    https://github.com/${REPO_SLUG}/releases/tag/${RELEASE_TAG}
    -- or Actions -> Build Writerdeck-server -> latest green run -> writerdeck-server artifact
  Save as Writerdeck-server in the repo root.
EOF
}

RUN_ID="${1:-}"
ok=0
if [ -n "$RUN_ID" ]; then
    if fetch_gh_run "$RUN_ID"; then ok=1; fi
else
    if fetch_release; then ok=1
    elif fetch_gh_run; then ok=1
    fi
fi

if [ "$ok" -ne 1 ] || [ ! -s "$OUT" ]; then
    print_help
    exit 1
fi

cp -f "$OUT" "$DEST"
chmod +x "$DEST"
echo "Fetched: $(ls -lh "$DEST")"
file "$DEST"
