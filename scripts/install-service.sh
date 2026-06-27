#!/usr/bin/env bash
# scripts/install-service.sh -- Install rm1-writerdeck.service on the device.
#
# Copies the unit to /etc/systemd/system/ and runs daemon-reload.
# Does NOT enable or start the unit (boot-loop guard -- see note below).
#
# Boot-loop guard: ExecStartPre stops xochitl, which can trip the firmware
# watchdog and reboot. An enabled unit on a bad install could loop:
#   boot -> stop xochitl -> watchdog reboot -> boot -> ...
# Safe sequence:
#   1. bash scripts/install-service.sh       (this script)
#   2. systemctl start rm1-writerdeck        (manual test -- from SSH)
#      Open Safari, type, press Home -> note saved, xochitl returns.
#   3. systemctl enable rm1-writerdeck       (only after step 2 passes)
# Recovery if stranded: ssh root@<ip>, then:
#   systemctl disable --now rm1-writerdeck
#   systemctl start xochitl
#
# Usage (run from repo root on the Mac):
#   bash scripts/install-service.sh
#   bash scripts/install-service.sh 192.168.1.8   # explicit target

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"

TARGET="${1:-$RM_HOST}"
UNIT_SRC="$DIR/rm1-writerdeck.service"
UNIT_DST="/etc/systemd/system/rm1-writerdeck.service"

echo "=== install-service  target=$TARGET ==="
echo

# ---------------------------------------------------------------------------
# Pre-flight
# ---------------------------------------------------------------------------
if [ ! -f "$UNIT_SRC" ]; then
    echo "ERROR: $UNIT_SRC not found." >&2
    exit 1
fi
echo "--- Testing SSH key login to $TARGET ---"
if ! ping -c1 -W2 "$TARGET" >/dev/null 2>&1; then
    echo "ERROR: $TARGET is unreachable (ping failed)." >&2
    echo "  The reMarkable is probably asleep -- wake it and try again." >&2
    exit 1
fi
if ! rm_test_key "$TARGET"; then
    echo "ERROR: key-based SSH to root@$TARGET failed." >&2
    echo "  Run: bash scripts/bootstrap.sh  to install the key." >&2
    exit 1
fi
echo "  OK"
echo

# ---------------------------------------------------------------------------
# Migrate: remove the old unit name if present (renamed rmnetwriter -> rm1-writerdeck).
# Stopping the old unit is safe -- rmkbd's ExecStopPost brings xochitl back.
# Idempotent: a no-op once the device is already on the new name.
# ---------------------------------------------------------------------------
echo "--- Removing any old rmnetwriter.service (renamed to rm1-writerdeck) ---"
rm_ssh 'if [ -f /etc/systemd/system/rmnetwriter.service ]; then systemctl disable --now rmnetwriter 2>/dev/null || true; rm -f /etc/systemd/system/rmnetwriter.service; systemctl daemon-reload; echo "  old unit removed"; else echo "  no old unit -- nothing to migrate"; fi' "$TARGET"
echo

# ---------------------------------------------------------------------------
# Install
# ---------------------------------------------------------------------------
echo "--- Copying unit to $TARGET:$UNIT_DST ---"
rm_scp_to "$UNIT_SRC" "$UNIT_DST" "$TARGET"
echo "  copied."

echo "--- systemctl daemon-reload ---"
rm_ssh "systemctl daemon-reload" "$TARGET"
echo "  done."
echo

# ---------------------------------------------------------------------------
# Verify it loaded cleanly (show any parse errors)
# ---------------------------------------------------------------------------
echo "--- Unit status (should show 'loaded' / inactive) ---"
rm_ssh "systemctl status rm1-writerdeck --no-pager 2>&1 || true" "$TARGET"
echo

# ---------------------------------------------------------------------------
# Instructions
# ---------------------------------------------------------------------------
echo "======================================"
echo "  INSTALL DONE  (unit installed, NOT enabled)"
echo "======================================"
echo ""
echo "  Step 1 -- manual test (do this first):"
echo "    ssh root@$TARGET"
echo "    systemctl start rm1-writerdeck"
echo "    # open Safari on $TARGET:8000, type, press Home"
echo "    # verify: note saved, xochitl UI returns"
echo ""
echo "  Step 2 -- enable autostart (only after step 1 passes):"
echo "    ssh root@$TARGET"
echo "    systemctl enable rm1-writerdeck"
echo ""
echo "  Recovery if stranded (boot loop or stuck):"
echo "    ssh root@$TARGET"
echo "    systemctl disable --now rm1-writerdeck"
echo "    systemctl start xochitl"
echo ""
echo "  To reinstall after a unit file change:"
echo "    bash scripts/install-service.sh"
echo "======================================"
