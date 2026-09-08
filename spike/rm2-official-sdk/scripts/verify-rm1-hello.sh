#!/usr/bin/env bash
# Prove official Qt6 epaper QPA on reMarkable 1 (no linuxfb, no Toltec).
# Usage (repo root, bash):
#   bash spike/rm2-official-sdk/scripts/verify-rm1-hello.sh
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"

export SDK_DEVICE=rm1
export RM_SPIKE_HOST="${RM_SPIKE_HOST:-${RM_HOST_WIFI:-}}"
if [ -z "$RM_SPIKE_HOST" ]; then
  err "RM_HOST_WIFI empty; run: bash scripts/discover-rm2-wifi.sh"
  exit 1
fi
export SPIKE_REMOTE_DIR="${SPIKE_REMOTE_DIR:-/home/root/spike-rm1-hello}"
export SPIKE_TOUCH="${SPIKE_TOUCH:-rotate=180}"
export SPIKE_SCREEN_REMOTE="$SPIKE_REMOTE_DIR/screen.png"

echo "rM1 hello host: $RM_SPIKE_HOST"
echo "SDK_DEVICE=$SDK_DEVICE"

if [ ! -f "$SPIKE/.cache/out-rm1/hello_remarkable" ]; then
  echo "Building hello_remarkable with rm1 SDK ..."
  bash "$SPIKE/scripts/build-hello.sh"
fi
bash "$SPIKE/scripts/deploy-hello.sh"
bash "$SPIKE/scripts/run-and-capture.sh" "rm1-hello"

PNG="$(ls -t "$SPIKE/screenshots/rm2-spike-"*"rm1-hello"*.png 2>/dev/null | head -n1)"
if [ -z "$PNG" ]; then
  err "no screenshot written for rm1-hello"
  exit 1
fi
bash "$SPIKE/scripts/verify-png.sh" "$PNG" 400
echo "PASS: hello on rM1 $RM_SPIKE_HOST -> $PNG"
