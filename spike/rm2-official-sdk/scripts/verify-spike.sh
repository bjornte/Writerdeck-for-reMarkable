#!/usr/bin/env bash
# Autonomous spike check: pick rM2 host, run app, pull PNG, verify content.
# Usage: verify-spike.sh [hello|textedit|socket|fork-probe]
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"

APP="${1:-hello}"
HOST="$(bash "$SPIKE/scripts/rm2-pick-host.sh")"
export RM2_HOST="$HOST"

case "$APP" in
  hello)
    REMOTE="/home/root/spike-rm2-hello/screen.png"
    MIN_DARK=400
    LABEL="hello-verify"
    RUN="$SPIKE/scripts/run-and-capture.sh"
    ;;
  textedit)
    REMOTE="/home/root/spike-rm2-textedit/screen.png"
    MIN_DARK=800
    LABEL="textedit-verify"
    RUN="$SPIKE/scripts/run-and-capture-textedit.sh"
    ;;
  socket)
    REMOTE="/home/root/spike-rm2-socket/screen.png"
    MIN_DARK=2400
    LABEL="socket-verify"
    RUN="$SPIKE/scripts/run-and-capture-socket.sh"
    ;;
  fork-probe)
    REMOTE="/home/root/spike-rm2-fork-probe/screen.png"
    MIN_DARK=8000
    LABEL="fork-probe-verify"
    RUN="$SPIKE/scripts/run-and-capture-fork-probe.sh"
    ;;
  *)
    err "unknown app: $APP (hello|textedit|socket|fork-probe)"
    exit 1
    ;;
esac

echo "rM2 host: $HOST"
echo "app: $APP"

if [ "$APP" = "textedit" ]; then
  if [ ! -f "$SPIKE/.cache/out/textedit_spike" ]; then
    echo "Building textedit_spike ..."
    bash "$SPIKE/scripts/build-textedit.sh"
  fi
  bash "$SPIKE/scripts/deploy-textedit.sh"
elif [ "$APP" = "fork-probe" ]; then
  if [ ! -f "$SPIKE/.cache/out/Writerdeck_fork_probe" ]; then
    echo "Building Writerdeck fork probe ..."
    bash "$SPIKE/scripts/build-fork-probe.sh"
  fi
  bash "$SPIKE/scripts/deploy-fork-probe.sh"
elif [ "$APP" = "socket" ]; then
  if [ ! -f "$SPIKE/.cache/out/socket_spike" ] || [ ! -f "$SPIKE/.cache/out/socket_inject" ]; then
    echo "Building socket_spike ..."
    bash "$SPIKE/scripts/build-socket.sh"
  fi
  bash "$SPIKE/scripts/deploy-socket.sh"
else
  if [ ! -f "$SPIKE/.cache/out/hello_remarkable" ]; then
    echo "Building hello_remarkable ..."
    bash "$SPIKE/scripts/build-hello.sh"
  fi
  bash "$SPIKE/scripts/deploy-hello.sh"
fi

bash "$RUN" "$LABEL"
PNG="$(ls -t "$SPIKE/screenshots/rm2-spike-"*"-${LABEL}"*.png 2>/dev/null | head -n1)"
if [ -z "$PNG" ]; then
  err "no screenshot written for label $LABEL"
  exit 1
fi

bash "$SPIKE/scripts/verify-png.sh" "$PNG" "$MIN_DARK"
echo "PASS: $APP on $HOST -> $PNG"
