#!/usr/bin/env bash
# scripts/lobby.sh -- from the Mac: SSH to the tablet and run wd (Lobby on e-ink).
#
# Usage:
#   bash scripts/lobby.sh
#   bash scripts/lobby.sh 192.168.1.8
#
# On the tablet (SSH session):  ~/wd
#
# Mac alias (after install-alias.sh):  wd

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"

TARGET="${1:-$RM_HOST}"

if ! ping -c1 -W2 "$TARGET" >/dev/null 2>&1; then
  err "Tablet unreachable at $TARGET (asleep or off Wi-Fi?)."
  exit 1
fi

rm_ssh "$DEVICE_WD" "$TARGET"
