#!/usr/bin/env bash
# keymaps/generate.sh -- Build .qmap files from Linux kmap text via Docker.
#
# Requires: docker (Colima on Mac: `colima start`).
# Output: keymaps/{us,no,es,de,fr}.qmap (committed to the repo).
#
# kmap sources (from device /usr/share/keymaps/i386/):
#   qwerty/us.map          -- US QWERTY baseline
#   qwerty/no-latin1.map   -- Norwegian (ae/oe/aa, AltGr symbols)
#   qwerty/es.map          -- Spanish
#   qwertz/de-latin1.map   -- German
#   azerty/fr-latin1.map   -- French
#
# Usage (from repo root):
#   bash keymaps/generate.sh
#   bash keymaps/generate.sh de fr   # subset

set -euo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
OUT="$DIR"
SRC="$DIR/src"
KMAP2QMAP_SRC="$DIR/kmap2qmap-main.cpp"

# Default layouts shipped in-repo.
if [ "$#" -eq 0 ]; then
  LAYOUTS=(us no es de fr)
else
  LAYOUTS=("$@")
fi

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required (install Colima + docker on Mac: brew install colima docker && colima start)." >&2
  exit 1
fi

mkdir -p "$SRC/i386/qwerty" "$SRC/i386/qwertz" "$SRC/i386/azerty"

# Resolve layout id -> relative path under i386/
layout_kmap() {
  case "$1" in
    us) echo "qwerty/us.map" ;;
    no) echo "qwerty/no-latin1.map" ;;
    es) echo "qwerty/es.map" ;;
    de) echo "qwertz/de-latin1.map" ;;
    fr) echo "azerty/fr-latin1.map" ;;
    *)  echo "qwerty/$1.map" ;;
  esac
}

# Fetch missing kmap files from the tablet when possible.
need_fetch=0
for layout in "${LAYOUTS[@]}"; do
  rel=$(layout_kmap "$layout")
  if [ ! -f "$SRC/i386/$rel" ]; then
    need_fetch=1
    break
  fi
done

if [ "$need_fetch" -eq 1 ]; then
  if [ -f "$REPO/secrets/remarkable.local.env" ]; then
    # shellcheck source=/dev/null
    . "$REPO/scripts/_env.sh" 2>/dev/null || true
  fi
  if [ -n "${RM_HOST:-}" ] && ping -c1 -W2 "$RM_HOST" >/dev/null 2>&1; then
    echo "Fetching kmap files from root@$RM_HOST ..."
    for layout in "${LAYOUTS[@]}"; do
      rel=$(layout_kmap "$layout")
      if [ -f "$SRC/i386/$rel" ]; then
        continue
      fi
      mkdir -p "$SRC/i386/$(dirname "$rel")"
      scp -o StrictHostKeyChecking=accept-new \
        "root@$RM_HOST:/usr/share/keymaps/i386/${rel}.gz" \
        "$SRC/i386/$(dirname "$rel")/" 2>/dev/null || true
      gunzip -f "$SRC/i386/${rel}.gz" 2>/dev/null || true
      if [ -f "$SRC/i386/$rel" ] && ! grep -q writerdeck-alt-arrows "$SRC/i386/$rel"; then
        printf '\ninclude "writerdeck-alt-arrows.inc"\n' >> "$SRC/i386/$rel"
      fi
    done
  fi
fi

missing=0
for layout in "${LAYOUTS[@]}"; do
  rel=$(layout_kmap "$layout")
  if [ ! -f "$SRC/i386/$rel" ]; then
    echo "ERROR: missing keymaps/src/i386/$rel" >&2
    missing=1
  fi
done
if [ "$missing" -eq 1 ]; then
  echo "  Copy from the device or run with SSH access." >&2
  exit 1
fi

if [ ! -d "$SRC/i386/include" ]; then
  echo "ERROR: missing keymaps/src/i386/include/ (needed for kmap includes)" >&2
  exit 1
fi

echo "=== Generating qmaps: ${LAYOUTS[*]} ==="

# Build a shell snippet that maps layout -> relative kmap path for the container.
layout_cases=""
for layout in "${LAYOUTS[@]}"; do
  rel=$(layout_kmap "$layout")
  layout_cases="${layout_cases}${layout}) kmap=\"${rel}\" ;; "
done

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
        '"$layout_cases"'
        *) echo "unknown layout: $layout" >&2; exit 1 ;;
      esac
      echo "  kmap2qmap /src/i386/$kmap -> /out/$layout.qmap"
      cd "/src/i386/$(dirname "$kmap")"
      "$work/kmap2qmap" "$(basename "$kmap")" "/out/$layout.qmap"
      ls -lh "/out/$layout.qmap"
    done
  '

echo
echo "Done. Commit keymaps/*.qmap and deploy with deploy-keywriter.sh."
