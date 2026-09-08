#!/usr/bin/env bash
# Enable tethered XOVI sidebar hook if OS is tested (or ACCEPT_UNTESTED=1 for candidates).
# Prefer: install-native.sh once, then writerdeck.service / Writerdeck-server re-arm automatically.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"

HOST="${1:-}"
if [ -z "$HOST" ]; then
  HOST="$(bash "$ROOT/spike/rm2-official-sdk/scripts/rm2-pick-host.sh")"
fi

echo "OS gate:"
bash "$SPIKE/scripts/check-os.sh" "$HOST"

if [ "${ACCEPT_UNTESTED:-}" = "1" ]; then
  rm_ssh 'touch /home/root/xovi/exthome/qt-resource-rebuilder/accept-untested' "$HOST"
fi

# Keep on-device gate file + ensure script + OS-matching .qmd current.
rm_send_file "$SPIKE/tested-os.json" \
  /home/root/xovi/exthome/qt-resource-rebuilder/tested-os.json "$HOST"
rm_send_file "$SPIKE/scripts/writerdeck-ensure-sidebar.sh" \
  /home/root/writerdeck-ensure-sidebar.sh "$HOST"
rm_send_file "$SPIKE/scripts/writerdeck-sidebar-launch.sh" \
  /home/root/writerdeck-sidebar-launch.sh "$HOST"
rm_ssh 'chmod a+x /home/root/writerdeck-ensure-sidebar.sh /home/root/writerdeck-sidebar-launch.sh' "$HOST"

OS_VER="$(rm_ssh 'sed -n "s/^IMG_VERSION=//p" /etc/os-release | tr -d "\""' "$HOST" | tr -d '\r')"
QMD_SRC="$SPIKE/qml/writerdeck-sidebar.qmd"
case "$OS_VER" in
  3.28.*) QMD_SRC="$SPIKE/qml/writerdeck-sidebar-3.28.qmd" ;;
esac
echo "qmd: $(basename "$QMD_SRC") (os=$OS_VER)"
test -f "$QMD_SRC"
rm_send_file "$QMD_SRC" \
  /home/root/xovi/exthome/qt-resource-rebuilder/writerdeck-sidebar.qmd "$HOST"

echo "Starting XOVI via writerdeck-ensure-sidebar.sh..."
rm_ssh '
set -e
test -x /home/root/writerdeck-ensure-sidebar.sh
test -x /home/root/xovi/start
test -f /home/root/xovi/extensions.d/qt-resource-rebuilder.so
test -f /home/root/xovi/extensions.d/qt-command-executor.so
test -f /home/root/xovi/exthome/qt-resource-rebuilder/writerdeck-sidebar.qmd
test -x /home/root/writerdeck-sidebar-launch.sh
test -s /home/root/xovi/exthome/qt-resource-rebuilder/hashtab
test ! -f /home/root/xovi/extensions.d/appload.so
/home/root/writerdeck-ensure-sidebar.sh
# ensure no-ops when the tmpfs drop-in is already mounted; reload qmd anyway
if pidof Writerdeck >/dev/null 2>&1; then
  echo "editor running -- skip xochitl restart"
else
  systemctl restart xochitl
  sleep 4
fi
echo "xochitl=$(systemctl is-active xochitl) NRestarts=$(systemctl show xochitl -p NRestarts --value)"
journalctl -u xochitl -n 50 --no-pager 2>/dev/null | grep -i qmldiff || true
' "$HOST"

echo
echo "Sidebar arm requested. Unlock -> hamburger -> Writerdeck."
echo "Writerdeck-server and writerdeck.service re-arm this on start / before xochitl."
echo "If xochitl crash-loops: bash spike/rm2-sidebar-launch/scripts/disable.sh"
