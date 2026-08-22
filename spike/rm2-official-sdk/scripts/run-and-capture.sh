#!/usr/bin/env bash
# Run hello on rM2, capture framebuffer PNG, restore xochitl.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="${RM2_HOST:-${RM2_HOST_WIFI:-10.11.99.1}}"
LABEL="${1:-hello}"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -euo pipefail
cd /home/root/spike-rm2-hello
export HOME=/home/root
export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS="rotate=180:invertx"
export QT_QUICK_BACKEND=epaper
systemctl stop xochitl
sleep 1
./hello_remarkable -platform epaper &
echo $! > /tmp/hello.pid
REMOTE

sleep 4
bash "$SPIKE/scripts/capture-screenshot.sh" "$LABEL"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -euo pipefail
if [ -f /tmp/hello.pid ]; then
  kill "$(cat /tmp/hello.pid)" 2>/dev/null || true
  rm -f /tmp/hello.pid
fi
killall hello_remarkable 2>/dev/null || true
systemctl start xochitl
REMOTE
