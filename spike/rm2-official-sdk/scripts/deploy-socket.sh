#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="$(rm2_pick_host)" || { err "no reachable rM2"; exit 1; }
BIN="$SPIKE/.cache/out/socket_spike"
INJ="$SPIKE/.cache/out/socket_inject"
if [ ! -f "$BIN" ] || [ ! -f "$INJ" ]; then
  err "missing build output - run scripts/build-socket.sh first"
  exit 1
fi
REMOTE_DIR="/home/root/spike-rm2-socket"
rm_ssh "killall socket_spike 2>/dev/null || true; mkdir -p '$REMOTE_DIR/plugins/scenegraph'" "$HOST"
sleep 1
rm_send_file "$BIN" "$REMOTE_DIR/socket_spike" "$HOST"
rm_send_file "$INJ" "$REMOTE_DIR/socket_inject" "$HOST"
rm_ssh "chmod +x '$REMOTE_DIR/socket_spike' '$REMOTE_DIR/socket_inject'" "$HOST"
if [ -f "$SPIKE/.cache/out/libqsgepaper.so" ]; then
  rm_send_file "$SPIKE/.cache/out/libqsgepaper.so" "$REMOTE_DIR/plugins/scenegraph/libqsgepaper.so" "$HOST"
fi
echo "Deployed to root@$HOST:$REMOTE_DIR/"
