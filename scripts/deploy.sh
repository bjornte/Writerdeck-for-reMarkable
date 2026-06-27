#!/usr/bin/env bash
# scripts/deploy.sh -- Cross-build rmkbd (ARMv7 static) and copy to the device.
# Run from the repo root on the Mac (the only machine that can reach the tablet).
#
# Usage:
#   bash scripts/deploy.sh               # build + deploy + kill any running rmkbd
#   bash scripts/deploy.sh --build-only  # just build; no device connection needed
#
# Requires: go (1.21+) on the Mac.

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"

BUILD_ONLY=0
[ "${1:-}" = "--build-only" ] && BUILD_ONLY=1

BINARY="${REPO}/rmkbd"

echo "=== rmkbd: cross-build (ARMv7 static) ==="
cd "${REPO}/daemon"
# Ensure go.sum is up to date (safe to run every time; needed after go.mod changes).
go mod tidy
GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 \
    go build -trimpath -o "${BINARY}" .
echo "  built: $(file "${BINARY}")"
echo

if [ "${BUILD_ONLY}" = "1" ]; then
    echo "  --build-only: skipping device deploy."
    exit 0
fi

# shellcheck source=/dev/null
. "${DIR}/_env.sh"

echo "=== rmkbd: deploying to ${RM_HOST} ==="
# Kill any running rmkbd FIRST: a live instance holds the executable busy
# (ETXTBSY -> scp "dest open Failure"). pkill -x proved unreliable on this
# device, so kill by full path + pidof.
rm_ssh 'pkill -f /home/root/rmkbd 2>/dev/null; for p in $(pidof rmkbd); do kill "$p" 2>/dev/null; done; sleep 0.5; true' >/dev/null
echo "  any old rmkbd stopped."
# Stream to a temp name, then atomically mv into place (rename never hits ETXTBSY).
# rm_send_file = gzip-over-ssh stream (scp deadlocks on this link). See _env.sh.
printf '    '; with_ticker 5 rm_send_file "${BINARY}" /home/root/rmkbd.new
rm_ssh 'mv -f /home/root/rmkbd.new /home/root/rmkbd && chmod +x /home/root/rmkbd' >/dev/null
echo "  /home/root/rmkbd updated."
echo

echo "======================================"
echo "  DEPLOY DONE"
echo "======================================"
echo "  Binary : /home/root/rmkbd"
echo "  Next   : bash scripts/test-phase4.sh"
echo "======================================"
