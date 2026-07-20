#!/usr/bin/env bash
# scripts/test-settings-tab.sh -- tablet Settings socket ops + lobby navigation.
#
# Exercises setreadfont/setpindigits via the trusted socket path (test hook),
# navigates to Settings (Tab from Files) over WebSocket, restores prior settings.
#
# Run after Writerdeck QML deploy (deploy-keywriter.sh -b).
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
REPO="$(cd "$DIR/.." && pwd)"
TARGET="${RM_HOST}"

RECON_DIR="$REPO/docs/recon"
mkdir -p "$RECON_DIR"
TS="$(date +%Y-%m-%dT%H-%M-%S)"
LOG="$RECON_DIR/test-settings-tab-$TS.txt"

{
  echo "=== test-settings-tab  $TS  target=$TARGET ==="
  go run -C "$REPO/daemon" ./cmd/settings-tab-test -host "$TARGET"
} 2>&1 | tee "$LOG"
