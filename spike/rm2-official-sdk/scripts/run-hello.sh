#!/usr/bin/env bash
# Run hello on the rM2 using the official epaper QPA. Restores xochitl on exit.
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
REMOTE_DIR="/home/root/spike-rm2-hello"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "SPIKE_HELLO_SECS='$SECS' REMOTE_DIR='$REMOTE_DIR' bash -s" <<'REMOTE'
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
export QT_QUICK_BACKEND=epaper
export HOME=/home/root

cleanup() {
  echo "restarting xochitl..."
  systemctl start xochitl || true
}
trap cleanup EXIT

systemctl stop xochitl
sleep 1
echo "starting hello_remarkable for ${SPIKE_HELLO_SECS}s..."
set +e
timeout "$SPIKE_HELLO_SECS" ./hello_remarkable -platform epaper
ec=$?
set -e
if [ "$ec" -eq 124 ]; then
  echo "hello timed out after ${SPIKE_HELLO_SECS}s (expected for unattended run)"
  exit 0
fi
exit "$ec"
REMOTE
