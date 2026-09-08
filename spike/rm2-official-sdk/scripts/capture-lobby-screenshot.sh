#!/usr/bin/env bash
# Capture Writerdeck lobby PNG on rM2 (Keyboard tab via probe hook).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
TARGET="${1:-10.11.99.1}"
OUT="${2:-$ROOT/docs/screenshots/writerdeck-rm2-lobby-keyboard-$(date +%Y-%m-%d).png}"
REMOTE="/home/root/spike-rm2-fork-probe/screen.png"

rm_ssh "killall Writerdeck 2>/dev/null || true; sleep 2; mkdir -p /home/root/spike-rm2-fork-probe; rm -f '$REMOTE'" "$TARGET"
curl -sf -X POST "http://$TARGET:8000/api/launch" >/dev/null
echo "Waiting for lobby screenshot hook (10s)..."
sleep 14
if ! rm_ssh "test -s '$REMOTE'" "$TARGET"; then
  err "no PNG at $REMOTE"
  exit 1
fi
rm_scp_from "$REMOTE" "$OUT" "$TARGET"
echo "Wrote $OUT"
