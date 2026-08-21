#!/usr/bin/env bash
# spike/rm2-official-sdk/scripts/recon-rm2.sh - read-only facts from the rM2.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"
HOST="${RM2_HOST:-${RM2_HOST_WIFI:-$RM_HOST}}"
ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=8 "root@$HOST" 'bash -s' <<'REMOTE'
set -e
echo "=== os-release ==="
cat /etc/os-release
echo "=== version ==="
cat /etc/version 2>/dev/null || true
echo "=== uname ==="
uname -a
echo "=== Qt ==="
ls /usr/lib/libQt6Core.so* 2>/dev/null || true
ls /usr/lib/plugins/platforms 2>/dev/null || true
ls /usr/lib/plugins/scenegraph 2>/dev/null || true
echo "=== fb ==="
ls -la /dev/fb* 2>/dev/null || true
cat /sys/class/graphics/fb0/name 2>/dev/null || true
cat /sys/class/graphics/fb0/virtual_size 2>/dev/null || true
echo "=== xochitl ==="
systemctl is-active xochitl 2>/dev/null || true
REMOTE
