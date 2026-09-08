#!/usr/bin/env bash
# Run hello, wait for grabWindow PNG, pull to screenshots/, restore xochitl.
# rM2 default. rM1: RM_SPIKE_HOST + SPIKE_TOUCH=rotate=180 + SPIKE_REMOTE_DIR.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
# shellcheck source=/dev/null
. "$SPIKE/scripts/sdk-url.sh"

HOST="${RM_SPIKE_HOST:-}"
if [ -z "$HOST" ]; then
  HOST="$(rm2_pick_host)" || { err "no reachable rM2"; exit 1; }
fi
LABEL="${1:-hello}"
if [ "${SDK_DEVICE:-rm2}" = "rm1" ]; then
  REMOTE_DIR="${SPIKE_REMOTE_DIR:-/home/root/spike-rm1-hello}"
  TOUCH="${SPIKE_TOUCH:-rotate=180}"
else
  REMOTE_DIR="${SPIKE_REMOTE_DIR:-/home/root/spike-rm2-hello}"
  TOUCH="${SPIKE_TOUCH:-rotate=180:invertx}"
fi
REMOTE="${SPIKE_SCREEN_REMOTE:-$REMOTE_DIR/screen.png}"
export SPIKE_SCREEN_REMOTE="$REMOTE"
export RM_SPIKE_HOST="$HOST"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "rm -f '$REMOTE'"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "REMOTE_DIR='$REMOTE_DIR' TOUCH='$TOUCH' bash -s" <<'REMOTE'
set -euo pipefail
# Stop the editor child only (not Writerdeck-server).
for p in $(pidof Writerdeck 2>/dev/null); do kill -TERM "$p" 2>/dev/null || true; done
killall hello_remarkable 2>/dev/null || true
systemctl stop xochitl
sleep 2
cd "$REMOTE_DIR"
rm -f screen.png
export HOME=/home/root
export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS="$TOUCH"
export QT_QUICK_BACKEND=epaper
./hello_remarkable -platform epaper > /tmp/hello.log 2>&1 &
echo $! > /tmp/hello.pid
REMOTE

SPIKE_SCREEN_WAIT=15 RM_SPIKE_HOST="$HOST" SPIKE_SCREEN_REMOTE="$REMOTE" \
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
