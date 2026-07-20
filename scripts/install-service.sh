#!/usr/bin/env bash
# scripts/install-service.sh -- Install writerdeck.service on the device.
#
# Copies the unit to /etc/systemd/system/ and runs daemon-reload.
# With --start: starts the unit, checks it is healthy, then enables autostart.
# Without --start: installs only (not started, not enabled).
#
# Recovery: systemctl disable --now writerdeck && systemctl start xochitl
#
# Usage (run from repo root on the Mac):
#   bash scripts/install-service.sh
#   bash scripts/install-service.sh --start
#   bash scripts/install-service.sh 192.168.1.8
#   bash scripts/install-service.sh --start 192.168.1.8

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
# shellcheck source=/dev/null
. "$DIR/migrate-device-layout.sh"

DO_START=0
TARGET="$RM_HOST"
for arg in "$@"; do
  case "$arg" in
    --start) DO_START=1 ;;
    -h|--help)
      sed -n '2,16p' "$0"
      exit 0
      ;;
    *) TARGET="$arg" ;;
  esac
done

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

echo "--- Unit status (should show 'loaded' / inactive or active) ---"
rm_ssh "systemctl status $SYSTEMD_UNIT --no-pager 2>&1 || true" "$TARGET"
echo

if [ "$DO_START" -eq 1 ]; then
  echo "--- systemctl start $SYSTEMD_UNIT ---"
  rm_ssh "systemctl start $SYSTEMD_UNIT" "$TARGET"
  sleep 2
  echo "--- journalctl -u $SYSTEMD_UNIT -n 30 ---"
  rm_ssh "journalctl -u $SYSTEMD_UNIT -n 30 --no-pager 2>&1 || true" "$TARGET"
  echo

  echo "--- health check ---"
  ACTIVE="$(rm_ssh "systemctl is-active $SYSTEMD_UNIT 2>/dev/null || true" "$TARGET" | tr -d '\r')"
  HTTP_OK=0
  if curl -fsS -m 8 "http://${TARGET}:8000/" >/dev/null 2>&1; then
    HTTP_OK=1
  fi
  if [ "$ACTIVE" != "active" ] || [ "$HTTP_OK" -ne 1 ]; then
    echo "ERROR: service not healthy (active=${ACTIVE:-unknown}, phone_http=${HTTP_OK})." >&2
    echo "  Autostart NOT enabled. Fix, then re-run: bash scripts/install-service.sh --start" >&2
    echo "  Recovery: systemctl disable --now writerdeck && systemctl start xochitl" >&2
    exit 1
  fi
  echo "  OK  active + http://${TARGET}:8000/"

  echo "--- systemctl enable $SYSTEMD_UNIT (autostart on boot) ---"
  rm_ssh "systemctl enable $SYSTEMD_UNIT" "$TARGET"
  echo "  enabled."
  echo
  echo "======================================"
  echo "  UNIT INSTALLED + STARTED + ENABLED"
  echo "======================================"
  echo ""
  echo "  Phone UI:  http://${TARGET}:8000/"
  echo ""
  echo "  You're done when: stock UI on e-ink; phone list populated;"
  echo "  connection bar not stuck on connecting..."
  echo "  Open Lobby: both page buttons, USB Esc, or wd"
  echo ""
  echo "  Recovery (if something goes wrong after reboot):"
  echo "    systemctl disable --now writerdeck && systemctl start xochitl"
  echo "======================================"
else
  echo "======================================"
  echo "  INSTALL DONE  (unit installed, NOT started)"
  echo "======================================"
  echo ""
  echo "  Start + enable autostart:"
  echo "    bash scripts/install-service.sh --start"
  echo ""
  echo "  Recovery:"
  echo "    systemctl disable --now writerdeck && systemctl start xochitl"
  echo "======================================"
fi
