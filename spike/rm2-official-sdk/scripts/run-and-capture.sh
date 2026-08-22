#!/usr/bin/env bash
# Run hello, wait for grabToImage PNG, pull to screenshots/, restore xochitl.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="${RM2_HOST:-${RM2_HOST_WIFI:-10.11.99.1}}"
LABEL="${1:-hello}"
REMOTE="/home/root/spike-rm2-hello/screen.png"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "rm -f '$REMOTE'"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -euo pipefail
cd /home/root/spike-rm2-hello
rm -f screen.png
export HOME=/home/root
export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS="rotate=180:invertx"
export QT_QUICK_BACKEND=epaper
systemctl stop xochitl
sleep 1
./hello_remarkable -platform epaper > /tmp/hello.log 2>&1 &
echo $! > /tmp/hello.pid
REMOTE

SPIKE_SCREEN_WAIT=12 RM2_HOST="$HOST" bash "$SPIKE/scripts/capture-screenshot.sh" "$LABEL"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -euo pipefail
if [ -f /tmp/hello.pid ]; then
  kill "$(cat /tmp/hello.pid)" 2>/dev/null || true
  rm -f /tmp/hello.pid
fi
killall hello_remarkable 2>/dev/null || true
systemctl start xochitl
REMOTE
