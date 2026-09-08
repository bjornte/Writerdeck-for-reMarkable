#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="$(rm2_pick_host)" || { err "no reachable rM2"; exit 1; }
BIN="$SPIKE/.cache/out/Writerdeck_fork_probe"
if [ ! -f "$BIN" ]; then
  err "missing $BIN - run scripts/build-fork-probe.sh first"
  exit 1
fi
REMOTE_DIR="/home/root/spike-rm2-fork-probe"
rm_ssh "killall Writerdeck_fork_probe 2>/dev/null || true; mkdir -p '$REMOTE_DIR/plugins/scenegraph'" "$HOST"
sleep 1
rm_send_file "$BIN" "$REMOTE_DIR/Writerdeck_fork_probe" "$HOST"
rm_ssh "chmod +x '$REMOTE_DIR/Writerdeck_fork_probe'" "$HOST"
if [ -f "$SPIKE/.cache/out/libqsgepaper.so" ]; then
  rm_send_file "$SPIKE/.cache/out/libqsgepaper.so" "$REMOTE_DIR/plugins/scenegraph/libqsgepaper.so" "$HOST"
fi
echo "Deployed to root@$HOST:$REMOTE_DIR/"
