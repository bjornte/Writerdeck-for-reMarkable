#!/usr/bin/env bash
# scripts/bootstrap.sh -- Phase 0: install SSH key on device + enable Wi-Fi SSH.
# macOS / Linux. Idempotent (safe to re-run).
#
# You will be prompted for the device root password AT MOST TWICE during the
# key-install step (once for scp, once for ssh). It is shown on screen below,
# read from secrets/remarkable.local.env. After this succeeds, no more prompts.
#
# If you hit a "host key changed" error after a firmware update:
#   ssh-keygen -R 10.11.99.1 ; ssh-keygen -R 192.168.1.8
#
# Usage:  bash scripts/bootstrap.sh        (run from repo root)
#         bash scripts/bootstrap.sh --skip-wifi

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"

SKIP_WIFI=0
[ "${1:-}" = "--skip-wifi" ] && SKIP_WIFI=1

KEY="$HOME/.ssh/id_ed25519"
PUB="$KEY.pub"

echo
echo "[1/4] SSH keypair"
mkdir -p "$HOME/.ssh" && chmod 700 "$HOME/.ssh"
if [ -f "$KEY" ]; then
  echo "  Already exists: $KEY"
else
  echo "  Generating ed25519 keypair at: $KEY"
  ssh-keygen -t ed25519 -f "$KEY" -N "" -q
  echo "  Generated."
fi

echo
# Install over the working host. $RM_HOST prefers Wi-Fi (USB-ethernet is
# inactive on the Mac); override with `export RM_HOST=10.11.99.1` if USB revives.
echo "[2/4] Install pubkey on device (host: $RM_HOST)"
if rm_test_key; then
  echo "  Key-based login already active -- skipping."
else
  echo
  echo "  *** PASSWORD PROMPT BELOW -- type it when ssh/scp asks ***"
  if [ -n "${RM_ROOT_PASSWORD:-}" ]; then
    echo "  Password: $RM_ROOT_PASSWORD"
  else
    echo "  Password: see secrets/remarkable.local.env (RM_ROOT_PASSWORD)"
  fi
  echo
  echo "  [1/2] scp public key to /tmp/laptop.pub ..."
  scp -o StrictHostKeyChecking=accept-new "$PUB" "root@$RM_HOST:/tmp/laptop.pub"
  echo "  [2/2] installing into authorized_keys ..."
  ssh -o StrictHostKeyChecking=accept-new "root@$RM_HOST" \
    'mkdir -p ~/.ssh; cat /tmp/laptop.pub >> ~/.ssh/authorized_keys; chmod 600 ~/.ssh/authorized_keys; chmod 700 ~/.ssh; rm -f /tmp/laptop.pub; echo __installed__'
  if rm_test_key; then
    echo "  Key login verified."
  else
    echo "  ERROR: key login still failing after install." >&2
    exit 1
  fi
fi

echo
echo "[3/4] Enable SSH over Wi-Fi (rm-ssh-over-wlan on)"
rm_ssh "rm-ssh-over-wlan on 2>&1 || true; echo __wlan_done__" | sed 's/^/  /'
echo "  Waiting 3s for the service to come up ..."
sleep 3

echo
echo "[4/4] Verify key login over Wi-Fi (${RM_HOST_WIFI:-<unset>})"
if [ "$SKIP_WIFI" = "1" ] || [ -z "${RM_HOST_WIFI:-}" ]; then
  echo "  Skipped (no Wi-Fi IP or --skip-wifi)."
elif rm_test_key "$RM_HOST_WIFI"; then
  echo "  Wi-Fi SSH: OK"
else
  echo "  WARN: Wi-Fi SSH failed on $RM_HOST_WIFI (device asleep or IP changed?)."
  echo "        Continue with USB; validate Wi-Fi during recon."
fi

echo
echo "===  Bootstrap complete. Run:  bash scripts/recon.sh  ==="
echo
