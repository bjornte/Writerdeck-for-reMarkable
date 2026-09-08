#!/usr/bin/env bash
# Disable XOVI UI mods and return to stock xochitl (also: just reboot).
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"

HOST="${1:-}"
if [ -z "$HOST" ]; then
  HOST="$(bash "$ROOT/spike/rm2-official-sdk/scripts/rm2-pick-host.sh")"
fi

rm_ssh '
set -e
cd /home/root
if [ -x xovi/stock ]; then
  xovi/stock
else
  # Fallback: clear LD_PRELOAD path by restarting stock service clean
  systemctl restart xochitl
fi
echo "stock UI requested"
' "$HOST"
echo "Disabled on $HOST (reboot also clears tethered XOVI)."
