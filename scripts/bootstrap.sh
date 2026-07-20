#!/usr/bin/env bash
# scripts/bootstrap.sh -- Phase 0: install SSH key on device + enable Wi-Fi SSH.
# macOS / Linux. Idempotent (safe to re-run).
#
# Uses RM_ROOT_PASSWORD from secrets when key login is not yet active (via
# SSH_ASKPASS -- no interactive password typing). After this succeeds, later
# scripts use key login only.
#
# If you hit a "host key changed" error after a firmware update:
#   bash scripts/fix-hostkey.sh
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

# Run ssh/scp with the tablet password from secrets (non-interactive).
_ssh_with_password() {
  local pwfile askpass ec
  if [ -z "${RM_ROOT_PASSWORD:-}" ]; then
    echo "  ERROR: RM_ROOT_PASSWORD empty; cannot install SSH key." >&2
    return 1
  fi
  pwfile="$(mktemp)"
  askpass="$(mktemp)"
  printf '%s\n' "$RM_ROOT_PASSWORD" >"$pwfile"
  chmod 600 "$pwfile"
  printf '#!/bin/sh\ncat "%s"\n' "$pwfile" >"$askpass"
  chmod 700 "$askpass"
  set +e
  SSH_ASKPASS="$askpass" SSH_ASKPASS_REQUIRE=force DISPLAY="${DISPLAY:-none}" \
    "$@" </dev/null
  ec=$?
  set -e
  rm -f "$askpass" "$pwfile"
  return "$ec"
}

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
  echo "  Installing SSH public key (using saved tablet password)..."
  echo "  [1/2] scp public key to /tmp/laptop.pub ..."
  if ! _ssh_with_password scp -o StrictHostKeyChecking=accept-new \
      -o PreferredAuthentications=password -o PubkeyAuthentication=no \
      "$PUB" "root@$RM_HOST:/tmp/laptop.pub"; then
    echo "  ERROR: scp of public key failed (wrong password? tablet awake?)." >&2
    exit 1
  fi
  echo "  [2/2] installing into authorized_keys ..."
  if ! _ssh_with_password ssh -o StrictHostKeyChecking=accept-new \
      -o PreferredAuthentications=password -o PubkeyAuthentication=no \
      "root@$RM_HOST" \
      'mkdir -p ~/.ssh; cat /tmp/laptop.pub >> ~/.ssh/authorized_keys; chmod 600 ~/.ssh/authorized_keys; chmod 700 ~/.ssh; rm -f /tmp/laptop.pub; echo __installed__'; then
    echo "  ERROR: ssh key install failed." >&2
    exit 1
  fi
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

if [ "${WRITERDECK_INSTALL:-0}" != "1" ]; then
  echo
  echo "===  Bootstrap complete. Run:  bash scripts/recon.sh  ==="
  echo
fi
