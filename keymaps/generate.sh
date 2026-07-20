#!/usr/bin/env bash
# keymaps/generate.sh -- Build .qmap files from Linux kmap text via Docker.
#
# Requires: docker (Colima on Mac: `colima start`).
# Output: keymaps/us.qmap, keymaps/no.qmap (committed to the repo).
#
# kmap sources (from device /usr/share/keymaps/i386/qwerty/):
#   us.map        -- US QWERTY baseline
#   no-latin1.map -- Norwegian (æøå, AltGr symbols)
#
# Usage (from repo root):
#   bash keymaps/generate.sh
#   bash keymaps/generate.sh de   # extra layout (needs .map in keymaps/src/)

set -euo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
OUT="$DIR"
SRC="$DIR/src"
KMAP2QMAP_SRC="$DIR/kmap2qmap-main.cpp"

# Default layouts shipped in-repo.
LAYOUTS=("${@:-us no}")

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required (install Colima + docker on Mac: brew install colima docker && colima start)." >&2
  exit 1
fi

mkdir -p "$SRC/i386/qwerty"

# Fetch kmap files from device if missing (reMarkable ships them in /usr/share/keymaps).
if [ ! -f "$SRC/i386/qwerty/us.map" ] || [ ! -f "$SRC/i386/qwerty/no-latin1.map" ]; then
  if [ -f "$REPO/secrets/remarkable.local.env" ]; then
    # shellcheck source=/dev/null
    . "$REPO/scripts/_env.sh" 2>/dev/null || true
  fi
  if [ -n "${RM_HOST:-}" ] && ping -c1 -W2 "$RM_HOST" >/dev/null 2>&1; then
    echo "Fetching kmap files from root@$RM_HOST ..."
    scp -o StrictHostKeyChecking=accept-new \
      "root@$RM_HOST:/usr/share/keymaps/i386/qwerty/us.map.gz" \
      "root@$RM_HOST:/usr/share/keymaps/i386/qwerty/no-latin1.map.gz" \
      "$SRC/i386/qwerty/" 2>/dev/null || true
    gunzip -f "$SRC/i386/qwerty"/*.gz 2>/dev/null || true
  fi
fi

if [ ! -f "$SRC/i386/qwerty/us.map" ] || [ ! -f "$SRC/i386/qwerty/no-latin1.map" ]; then
  echo "ERROR: missing keymaps/src/i386/qwerty/{us,no-latin1}.map" >&2
  echo "  Copy from the device or run with SSH access." >&2
  exit 1
fi

if [ ! -d "$SRC/i386/include" ]; then
  echo "ERROR: missing keymaps/src/i386/include/ (needed for kmap includes)" >&2
  exit 1
fi

echo "=== Generating qmaps: ${LAYOUTS[*]} ==="

docker run --rm \
  -v "$OUT:/out" \
  -v "$SRC:/src:ro" \
  -v "$KMAP2QMAP_SRC:/kmap2qmap.cpp:ro" \
  -v "$DIR/qevdevkeyboardhandler_p.h:/qevdevkeyboardhandler_p.h:ro" \
  debian:bookworm-slim \
  bash -c '
    set -euo pipefail
    apt-get update -qq
    apt-get install -qq -y --no-install-recommends \
      g++ make pkg-config qtbase5-dev >/dev/null

    work=$(mktemp -d)
    cd "$work"
    cp /kmap2qmap.cpp kmap2qmap.cpp
    cp /qevdevkeyboardhandler_p.h .
    sed -i "s|<QtInputSupport/private/qevdevkeyboardhandler_p.h>|\"qevdevkeyboardhandler_p.h\"|" kmap2qmap.cpp
    g++ -O2 -o kmap2qmap kmap2qmap.cpp -I. \
      $(pkg-config --cflags --libs Qt5Core)
    chmod +x kmap2qmap

    for layout in '"$(printf '%s ' "${LAYOUTS[@]}")"'; do
      case "$layout" in
        us) kmap="us.map" ;;
        no) kmap="no-latin1.map" ;;
        *)  kmap="$layout.map" ;;
      esac
      echo "  kmap2qmap /src/i386/qwerty/$kmap -> /out/$layout.qmap"
      cd /src/i386/qwerty
      "$work/kmap2qmap" "$kmap" "/out/$layout.qmap"
      ls -lh "/out/$layout.qmap"
    done
  '

echo
echo "Done. Commit keymaps/*.qmap and deploy with deploy-keywriter.sh."
