#!/usr/bin/env bash
# Copy built hello binary (and optional libqsgepaper.so) to the rM2.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="${RM2_HOST:-${RM2_HOST_WIFI:-}}"
if [ -z "$HOST" ]; then
  err "RM2_HOST_WIFI not set in secrets/remarkable.local.env"
  exit 1
fi
BIN="$SPIKE/.cache/out/hello_remarkable"
if [ ! -f "$BIN" ]; then
  err "missing $BIN - run scripts/build-hello.sh first"
  exit 1
fi
REMOTE_DIR="/home/root/spike-rm2-hello"
rm_ssh "mkdir -p '$REMOTE_DIR'" "$HOST"
rm_send_file "$BIN" "$REMOTE_DIR/hello_remarkable" "$HOST"
rm_ssh "chmod +x '$REMOTE_DIR/hello_remarkable'" "$HOST"
if [ -f "$SPIKE/.cache/out/libqsgepaper.so" ]; then
  rm_send_file "$SPIKE/.cache/out/libqsgepaper.so" "$REMOTE_DIR/libqsgepaper.so" "$HOST"
fi
echo "Deployed to root@$HOST:$REMOTE_DIR/"
