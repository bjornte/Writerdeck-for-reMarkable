#!/usr/bin/env bash
# scripts/fix-hostkey.sh -- clear a stale device SSH host key, then reconnect.
# macOS / Linux. Idempotent. LOCAL-ONLY: it edits ~/.ssh/known_hosts on the
# Mac and never changes anything on the tablet.
#
# Symptom this fixes (from `ssh root@<device>` or any deploy script):
#   "WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!"
#   "Host key verification failed."
# Cause: the reMarkable regenerates its SSH host key on a firmware update or
# reflash, so the key your Mac cached no longer matches. On a home LAN this is
# almost never an attack -- but if you did NOT just update firmware, eyeball the
# new fingerprint before trusting it.
#
# Hosts are read from secrets (RM_HOST_WIFI / RM_HOST_USB) via _env.sh, so the
# IPs are never hardcoded here and a DHCP change is handled automatically.
#
# Usage:  bash scripts/fix-hostkey.sh        (run from repo root, on the Mac)
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"

echo
echo "[1/2] Forgetting stale host keys in ~/.ssh/known_hosts"
for host in "${RM_HOST_WIFI:-}" "${RM_HOST_USB:-}"; do
  [ -n "$host" ] || continue
  if ssh-keygen -F "$host" >/dev/null 2>&1; then
    ssh-keygen -R "$host" >/dev/null 2>&1 && echo "  removed: $host"
  else
    echo "  (none cached: $host)"
  fi
done

echo
echo "[2/2] Reconnecting to $RM_HOST (accept-new re-adds the fresh key)"
if rm_test_key; then
  echo "  Key login OK -- you're back in, new host key now trusted."
else
  echo "  Stale key cleared, but key login did not succeed."
  echo "  A firmware update also resets the root password and may wipe your"
  echo "  installed SSH key and the systemd unit. Next steps:"
  echo "    1. Re-record RM_ROOT_PASSWORD in secrets/remarkable.local.env"
  echo "       (tablet: Settings > Help > Copyrights and licenses > General information)."
  echo "    2. bash scripts/bootstrap.sh      # re-install your SSH key"
  echo "    3. If boot-to-typewriter stopped: systemctl enable rm1-writerdeck (on device)."
fi
