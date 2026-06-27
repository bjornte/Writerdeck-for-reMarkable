#!/usr/bin/env bash
# scripts/recon.sh -- snapshot device facts (OS version, Wi-Fi IP, input nodes,
# disk) to a timestamped log under docs/recon/. macOS / Linux.
# Re-run after a firmware update to refresh the facts. Requires key-based SSH
# (run bootstrap.sh first).
#
# Usage:  bash scripts/recon.sh            (run from repo root)
#         bash scripts/recon.sh --backup   (also tar /home/root to docs/recon/)

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
REPO="$(cd "$DIR/.." && pwd)"

BACKUP=0
[ "${1:-}" = "--backup" ] && BACKUP=1

RECON_DIR="$REPO/docs/recon"
mkdir -p "$RECON_DIR"
TS="$(date +%Y-%m-%dT%H-%M-%S)"
LOG="$RECON_DIR/recon-$TS.txt"

log() { echo "$*" | tee -a "$LOG"; }

log "# rM1-Writerdeck device recon  $TS  target=$RM_HOST"
log ""

echo "=== Collecting device facts ==="

log "## OS version"
OSVER="$(rm_ssh 'cat /etc/version' | tr -d '\r' | head -n1 | sed -E 's/[[:space:]]+$//')"
log "$OSVER"; log ""

log "## ip addr"
IPADDR="$(rm_ssh 'ip addr' | tr -d '\r')"
log "$IPADDR"; log ""
WLAN_IP="$(echo "$IPADDR" | sed -n -E 's@.*inet ([0-9.]+)/[0-9]+ .*wlan0.*@\1@p' | head -n1)"

log "## /dev/input/"
log "$(rm_ssh 'ls -la /dev/input/' | tr -d '\r')"; log ""

log "## /proc/bus/input/devices"
log "$(rm_ssh 'cat /proc/bus/input/devices 2>/dev/null || echo (not available)' | tr -d '\r')"; log ""

log "## Disk"
log "$(rm_ssh 'df -h /' | tr -d '\r')"; log ""

if [ "$BACKUP" = "1" ]; then
  echo "=== Backup /home/root/ (may take ~30s) ==="
  rm_ssh 'tar czf /tmp/rm_bkp.tar.gz /home/root/ 2>/dev/null; echo __done__' >/dev/null
  rm_scp_from "/tmp/rm_bkp.tar.gz" "$RECON_DIR/home_root_backup-$TS.tar.gz"
  rm_ssh 'rm -f /tmp/rm_bkp.tar.gz' >/dev/null
  echo "  Saved: $RECON_DIR/home_root_backup-$TS.tar.gz"
  log "## Backup: home_root_backup-$TS.tar.gz"
fi

log ""
log "## SUMMARY  osVer='$OSVER'  wlanIp='$WLAN_IP'"

echo
echo "======================================"
echo "  RECON SUMMARY"
echo "======================================"
echo "  OS version  : $OSVER"
echo "  Wi-Fi IP    : ${WLAN_IP:-(could not extract)}"
echo "  Log saved   : $LOG"
echo "======================================"
echo
