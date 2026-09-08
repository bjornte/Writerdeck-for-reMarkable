#!/usr/bin/env bash
# Run socket_spike, inject NDJSON keys, pull PNG, restore xochitl.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="$(rm2_pick_host)" || { err "no reachable rM2"; exit 1; }
LABEL="${1:-socket}"
REMOTE="/home/root/spike-rm2-socket/screen.png"
SOCK="/home/root/spike-rm2-socket/Writerdeck.sock"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "rm -f '$REMOTE' '$SOCK'"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -euo pipefail
killall socket_spike 2>/dev/null || true
systemctl stop xochitl
sleep 2
cd /home/root/spike-rm2-socket
rm -f screen.png Writerdeck.sock
export HOME=/home/root
export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS="rotate=180:invertx"
export QT_QUICK_BACKEND=epaper
./socket_spike -platform epaper > /tmp/socket-spike.log 2>&1 &
echo $! > /tmp/socket-spike.pid
REMOTE

deadline=$(( $(date +%s) + 12 ))
while [ "$(date +%s)" -lt "$deadline" ]; do
  if ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
    "test -S '$SOCK'" 2>/dev/null; then
    break
  fi
  sleep 1
done

if ! ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "test -S '$SOCK'" 2>/dev/null; then
  err "socket not listening at $SOCK"
  ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
    'tail -n 20 /tmp/socket-spike.log 2>/dev/null || true'
  exit 1
fi

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" \
  "/home/root/spike-rm2-socket/socket_inject '$SOCK'"

SPIKE_SCREEN_REMOTE="$REMOTE" SPIKE_SCREEN_WAIT=12 RM2_HOST="$HOST" \
  bash "$SPIKE/scripts/capture-screenshot.sh" "$LABEL"

ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -euo pipefail
if [ -f /tmp/socket-spike.pid ]; then
  kill "$(cat /tmp/socket-spike.pid)" 2>/dev/null || true
  rm -f /tmp/socket-spike.pid
fi
killall socket_spike 2>/dev/null || true
rm -f /home/root/spike-rm2-socket/Writerdeck.sock
systemctl start xochitl
REMOTE
