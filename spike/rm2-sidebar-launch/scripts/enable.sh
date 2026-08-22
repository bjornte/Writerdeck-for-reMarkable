#!/usr/bin/env bash
# Enable tethered XOVI/AppLoad only if OS is tested (or ACCEPT_UNTESTED=1 for candidates).
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

echo "Starting XOVI (restarts xochitl with mods)..."
rm_ssh '
set -e
test -x /home/root/xovi/start
test -f /home/root/xovi/extensions.d/qt-resource-rebuilder.so
test -f /home/root/xovi/extensions.d/qt-command-executor.so
test -f /home/root/xovi/exthome/qt-resource-rebuilder/writerdeck-sidebar.qmd
test -s /home/root/xovi/exthome/qt-resource-rebuilder/hashtab
# AppLoad must stay off for this demo
test ! -f /home/root/xovi/extensions.d/appload.so
echo "hashtab ok ($(wc -c < /home/root/xovi/exthome/qt-resource-rebuilder/hashtab) bytes)"
cd /home/root
./xovi/start
sleep 5
journalctl -u xochitl -n 60 --no-pager 2>/dev/null | tail -60 || true
' "$HOST"

echo
echo "XOVI start requested. On the tablet after unlock:"
echo "  1) Tap hamburger (top-left)"
echo "  2) Tap Writerdeck in the sidebar (launches Lobby)"
echo "If xochitl crash-loops: bash spike/rm2-sidebar-launch/scripts/disable.sh"
echo "When it works: bash spike/rm2-sidebar-launch/scripts/mark-tested.sh"
