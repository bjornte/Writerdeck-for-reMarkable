#!/usr/bin/env bash
# Run textedit_spike, pull grabWindow PNG, restore xochitl.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="$(rm2_pick_host)" || { err "no reachable rM2"; exit 1; }
LABEL="${1:-textedit}"
REMOTE="/home/root/spike-rm2-textedit/screen.png"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "rm -f '$REMOTE'"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -euo pipefail
killall textedit_spike 2>/dev/null || true
systemctl stop xochitl
sleep 2
cd /home/root/spike-rm2-textedit
rm -f screen.png
export HOME=/home/root
export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS="rotate=180:invertx"
export QT_QUICK_BACKEND=epaper
./textedit_spike -platform epaper > /tmp/textedit.log 2>&1 &
echo $! > /tmp/textedit.pid
REMOTE

SPIKE_SCREEN_REMOTE="$REMOTE" SPIKE_SCREEN_WAIT=18 RM2_HOST="$HOST" \
  bash "$SPIKE/scripts/capture-screenshot.sh" "$LABEL"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -euo pipefail
if [ -f /tmp/textedit.pid ]; then
  kill "$(cat /tmp/textedit.pid)" 2>/dev/null || true
  rm -f /tmp/textedit.pid
fi
killall textedit_spike 2>/dev/null || true
systemctl start xochitl
REMOTE
