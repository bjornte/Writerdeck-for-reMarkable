#!/usr/bin/env bash
# Pull PNG written by hello_remarkable (grabToImage). fb0 is not readable on rM2 epaper.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="$(rm2_pick_host 2>/dev/null || true)"
HOST="${HOST:-${RM2_HOST:-${RM2_HOST_WIFI:-10.11.99.1}}}"
REMOTE="${SPIKE_SCREEN_REMOTE:-/home/root/spike-rm2-hello/screen.png}"
OUT_DIR="$ROOT/spike/rm2-official-sdk/screenshots"
mkdir -p "$OUT_DIR"
LABEL="${1:-hello}"
DATE="$(date +%Y-%m-%d)"
OUT="$OUT_DIR/rm2-spike-${DATE}-${LABEL}.png"
n=0
while [ -e "$OUT" ]; do
  n=$((n + 1))
  OUT="$OUT_DIR/rm2-spike-${DATE}-${LABEL}-${n}.png"
done

WAIT="${SPIKE_SCREEN_WAIT:-8}"
deadline=$(( $(date +%s) + WAIT ))
while [ "$(date +%s)" -lt "$deadline" ]; do
  if ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
    "test -s '$REMOTE'" 2>/dev/null; then
    break
  fi
  sleep 1
done

if ! ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "test -s '$REMOTE'" 2>/dev/null; then
  err "no screen PNG at $REMOTE (hello must run 2.5s+ first)"
  exit 1
fi

rm_scp_from "$REMOTE" "$OUT" "$HOST"
echo "Wrote $OUT"
