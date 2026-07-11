#!/usr/bin/env bash
# scripts/deploy-rmkbd.sh -- Cross-build Writerdeck-server (ARMv7 static) and deploy.
# Run from the repo root on the Mac (the only machine that can reach the tablet).
#
# Usage:
#   bash scripts/deploy-rmkbd.sh               # build + deploy
#   bash scripts/deploy-rmkbd.sh --build-only  # just build; no device connection needed
#
# Requires: go (1.21+) on the Mac.

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"

BUILD_ONLY=0
[ "${1:-}" = "--build-only" ] && BUILD_ONLY=1

BINARY="${REPO}/Writerdeck-server"

echo "=== Writerdeck-server: cross-build (ARMv7 static) ==="
cd "${REPO}/daemon"
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
# shellcheck source=/dev/null
. "${DIR}/migrate-device-layout.sh"

echo "=== Writerdeck-server: deploying to ${RM_HOST} ==="
migrate_device_layout
rm_ssh "pkill -f ${LEGACY_SERVER} 2>/dev/null; pkill -f ${DEVICE_SERVER} 2>/dev/null; for p in \$(pidof rmkbd Writerdeck-server 2>/dev/null); do kill \"\$p\" 2>/dev/null; done; sleep 0.5; true" >/dev/null
echo "  any old server stopped."
printf '    '; with_ticker 5 rm_send_file "${BINARY}" "${DEVICE_SERVER}.new"
rm_ssh "mv -f ${DEVICE_SERVER}.new ${DEVICE_SERVER} && chmod +x ${DEVICE_SERVER}" >/dev/null
echo "  ${DEVICE_SERVER} updated."
rm_deploy_wd
echo "  ${DEVICE_WD} updated."
echo

echo "======================================"
echo "  DEPLOY DONE"
echo "======================================"
echo "  Binary : ${DEVICE_SERVER}"
echo "  Next   : bash scripts/test-e2e.sh"
echo "======================================"
