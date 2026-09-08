#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="$(rm2_pick_host)" || { err "no reachable rM2"; exit 1; }
BIN="$SPIKE/.cache/out/textedit_spike"
if [ ! -f "$BIN" ]; then
  err "missing $BIN - run scripts/build-textedit.sh first"
  exit 1
fi
REMOTE_DIR="/home/root/spike-rm2-textedit"
rm_ssh "killall textedit_spike 2>/dev/null || true; mkdir -p '$REMOTE_DIR/plugins/scenegraph'" "$HOST"
sleep 1
rm_send_file "$BIN" "$REMOTE_DIR/textedit_spike" "$HOST"
rm_ssh "chmod +x '$REMOTE_DIR/textedit_spike'" "$HOST"
if [ -f "$SPIKE/.cache/out/libqsgepaper.so" ]; then
  rm_send_file "$SPIKE/.cache/out/libqsgepaper.so" "$REMOTE_DIR/plugins/scenegraph/libqsgepaper.so" "$HOST"
fi
echo "Deployed to root@$HOST:$REMOTE_DIR/"
