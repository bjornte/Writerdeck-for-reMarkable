#!/usr/bin/env bash
# scripts/install-service.sh -- Install writerdeck.service on the device.
#
# Copies the unit to /etc/systemd/system/ and runs daemon-reload.
# Does NOT enable or start the unit (boot-loop guard).
#
# Safe sequence:
#   1. bash scripts/install-service.sh
#   2. systemctl start writerdeck        (manual test -- from SSH)
#   3. systemctl enable writerdeck       (only after step 2 passes)
# Recovery: systemctl disable --now writerdeck && systemctl start xochitl
#
# Usage (run from repo root on the Mac):
#   bash scripts/install-service.sh
#   bash scripts/install-service.sh 192.168.1.8

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
# shellcheck source=/dev/null
. "$DIR/migrate-device-layout.sh"

TARGET="${1:-$RM_HOST}"
UNIT_SRC="$DIR/writerdeck.service"
UNIT_DST="$SYSTEMD_UNIT_PATH"

echo "=== install-service  target=$TARGET ==="
echo

if [ ! -f "$UNIT_SRC" ]; then
    echo "ERROR: $UNIT_SRC not found." >&2
    exit 1
fi
echo "--- Testing SSH key login to $TARGET ---"
if ! ping -c1 -W2 "$TARGET" >/dev/null 2>&1; then
    echo "ERROR: $TARGET is unreachable (ping failed)." >&2
    exit 1
fi
if ! rm_test_key "$TARGET"; then
    echo "ERROR: key-based SSH to root@$TARGET failed." >&2
    exit 1
fi
echo "  OK"
echo

echo "--- Migrating legacy layout (notes, settings, units) ---"
migrate_device_layout "$TARGET"
echo

echo "--- Copying unit to $TARGET:$UNIT_DST ---"
rm_scp_to "$UNIT_SRC" "$UNIT_DST" "$TARGET"
echo "  copied."

echo "--- systemctl daemon-reload ---"
rm_ssh "systemctl daemon-reload" "$TARGET"
echo "  done."
echo

echo "--- Unit status (should show 'loaded' / inactive) ---"
rm_ssh "systemctl status $SYSTEMD_UNIT --no-pager 2>&1 || true" "$TARGET"
echo

echo "======================================"
echo "  INSTALL DONE  (unit installed, NOT enabled)"
echo "======================================"
echo ""
echo "  Step 1 -- manual test:"
echo "    ssh root@$TARGET"
echo "    systemctl start writerdeck"
echo ""
echo "  Step 2 -- enable autostart (only after step 1 passes):"
echo "    systemctl enable writerdeck"
echo ""
echo "  Recovery:"
echo "    systemctl disable --now writerdeck"
echo "    systemctl start xochitl"
echo "======================================"
