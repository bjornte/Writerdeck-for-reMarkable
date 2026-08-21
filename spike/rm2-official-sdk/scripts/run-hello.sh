#!/usr/bin/env bash
# Run hello on the rM2 using the official epaper QPA. Restores xochitl on exit.
#
# Note (3.11.2.5 / Codex 4.0.447): stock has libepaper.so but no
# plugins/scenegraph/libqsgepaper.so, and SDK 4.0.367 likewise. Setting
# QT_QUICK_BACKEND=epaper aborts. Default/software Quick backends stay up.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="${RM2_HOST:-${RM2_HOST_WIFI:-}}"
if [ -z "$HOST" ]; then
  err "RM2_HOST_WIFI not set"
  exit 1
fi
SECS="${SPIKE_HELLO_SECS:-45}"
# Override: epaper | software | unset
BACKEND="${SPIKE_QUICK_BACKEND:-}"
REMOTE_DIR="/home/root/spike-rm2-hello"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "SPIKE_HELLO_SECS='$SECS' SPIKE_QUICK_BACKEND='$BACKEND' REMOTE_DIR='$REMOTE_DIR' bash -s" <<'REMOTE'
set -euo pipefail
cd "$REMOTE_DIR"
if [ ! -x ./hello_remarkable ]; then
  echo "missing $REMOTE_DIR/hello_remarkable - deploy first" >&2
  exit 1
fi
mkdir -p ./plugins/scenegraph
if [ -f ./libqsgepaper.so ]; then
  cp -f ./libqsgepaper.so ./plugins/scenegraph/
fi
export QT_PLUGIN_PATH="$PWD/plugins${QT_PLUGIN_PATH:+:$QT_PLUGIN_PATH}"
export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS="rotate=180:invertx"
export HOME=/home/root
export LC_ALL=C.UTF-8
export LANG=C.UTF-8
if [ -n "${SPIKE_QUICK_BACKEND:-}" ]; then
  export QT_QUICK_BACKEND="$SPIKE_QUICK_BACKEND"
  echo "QT_QUICK_BACKEND=$QT_QUICK_BACKEND"
else
  echo "QT_QUICK_BACKEND unset (default scene graph)"
fi

cleanup() {
  echo "restarting xochitl..."
  killall hello_remarkable 2>/dev/null || true
  sleep 1
  systemctl start xochitl || true
}
trap cleanup EXIT

systemctl stop xochitl
sleep 1
echo "starting hello_remarkable for ${SPIKE_HELLO_SECS}s - look at the tablet"
./hello_remarkable -platform epaper &
HPID=$!
sleep 2
if ! kill -0 "$HPID" 2>/dev/null; then
  echo "hello exited early" >&2
  wait "$HPID" || true
  exit 1
fi
echo "hello alive pid=$HPID"
sleep "$SPIKE_HELLO_SECS"
if kill -0 "$HPID" 2>/dev/null; then
  echo "stopping hello after ${SPIKE_HELLO_SECS}s"
  kill "$HPID" 2>/dev/null || true
  wait "$HPID" 2>/dev/null || true
fi
echo "hello run finished"
REMOTE
