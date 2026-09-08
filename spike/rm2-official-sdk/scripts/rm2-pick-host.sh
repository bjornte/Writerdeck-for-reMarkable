#!/usr/bin/env bash
# Print first reachable rM2 SSH host (USB, then Wi-Fi). Exit 1 if none.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"

if [ -n "${RM2_HOST:-}" ]; then
  if rm_test_key "$RM2_HOST"; then
    echo "$RM2_HOST"
    exit 0
  fi
fi

for h in "$RM2_HOST_USB" "$RM2_HOST_WIFI"; do
  [ -n "$h" ] || continue
  if rm_test_key "$h"; then
    echo "$h"
    exit 0
  fi
done

err "no reachable rM2 (USB $RM2_HOST_USB, Wi-Fi ${RM2_HOST_WIFI:-unset})"
exit 1
