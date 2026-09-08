#!/usr/bin/env bash
# Copy built hello binary (and optional libqsgepaper.so) to the tablet.
# Default: rM2 via USB/Wi-Fi pick. Override: RM_SPIKE_HOST + SPIKE_REMOTE_DIR.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
# shellcheck source=/dev/null
. "$SPIKE/scripts/sdk-url.sh"

HOST="${RM_SPIKE_HOST:-}"
if [ -z "$HOST" ]; then
  HOST="$(rm2_pick_host)" || { err "RM2 unreachable (set RM2_HOST_WIFI in secrets?)"; exit 1; }
fi
if [ "${SDK_DEVICE:-rm2}" = "rm1" ]; then
  BIN="$SPIKE/.cache/out-rm1/hello_remarkable"
  REMOTE_DIR="${SPIKE_REMOTE_DIR:-/home/root/spike-rm1-hello}"
else
  BIN="$SPIKE/.cache/out/hello_remarkable"
  REMOTE_DIR="${SPIKE_REMOTE_DIR:-/home/root/spike-rm2-hello}"
fi
if [ ! -f "$BIN" ]; then
  err "missing $BIN - run scripts/build-hello.sh first"
  exit 1
fi
rm_ssh "killall hello_remarkable 2>/dev/null || true; mkdir -p '$REMOTE_DIR/plugins/scenegraph'" "$HOST"
sleep 1
rm_send_file "$BIN" "$REMOTE_DIR/hello_remarkable" "$HOST"
rm_ssh "chmod +x '$REMOTE_DIR/hello_remarkable'" "$HOST"
SG_SRC=""
if [ -f "$SPIKE/.cache/out-rm1/libqsgepaper.so" ] && [ "${SDK_DEVICE:-rm2}" = "rm1" ]; then
  SG_SRC="$SPIKE/.cache/out-rm1/libqsgepaper.so"
elif [ -f "$SPIKE/.cache/out-ci327/libqsgepaper.so" ]; then
  SG_SRC="$SPIKE/.cache/out-ci327/libqsgepaper.so"
elif [ -f "$SPIKE/.cache/out/libqsgepaper.so" ]; then
  SG_SRC="$SPIKE/.cache/out/libqsgepaper.so"
fi
if [ -n "$SG_SRC" ]; then
  rm_send_file "$SG_SRC" "$REMOTE_DIR/plugins/scenegraph/libqsgepaper.so" "$HOST"
fi
echo "Deployed to root@$HOST:$REMOTE_DIR/"
