#!/usr/bin/env bash
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="$(rm2_pick_host)" || { err "no reachable rM2"; exit 1; }
LABEL="${1:-fork-probe}"
REMOTE="/home/root/spike-rm2-fork-probe/screen.png"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "rm -f '$REMOTE'"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -euo pipefail
killall Writerdeck_fork_probe 2>/dev/null || true
systemctl stop xochitl
sleep 2
cd /home/root/spike-rm2-fork-probe
rm -f screen.png
export HOME=/home/root
export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS="rotate=180:invertx"
export QT_QUICK_BACKEND=epaper
./Writerdeck_fork_probe -platform epaper > /tmp/fork-probe.log 2>&1 &
echo $! > /tmp/fork-probe.pid
REMOTE

SPIKE_SCREEN_REMOTE="$REMOTE" SPIKE_SCREEN_WAIT=20 RM2_HOST="$HOST" \
  bash "$SPIKE/scripts/capture-screenshot.sh" "$LABEL"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -euo pipefail
if [ -f /tmp/fork-probe.pid ]; then
  kill "$(cat /tmp/fork-probe.pid)" 2>/dev/null || true
  rm -f /tmp/fork-probe.pid
fi
killall Writerdeck_fork_probe 2>/dev/null || true
systemctl start xochitl
REMOTE
