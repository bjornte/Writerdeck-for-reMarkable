#!/usr/bin/env bash
# scripts/configure-sync.sh -- Push PIN_DIGITS / SYNC_REPO / GH_TOKEN from local
# secrets to the running Writerdeck-server (loopback on the tablet; PIN-safe).
# No-op sections when fields are empty.
#
# Usage (run from repo root after the service is up):
#   bash scripts/configure-sync.sh
#   bash scripts/configure-sync.sh 192.168.1.8
#
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"

TARGET="${1:-$RM_HOST}"
SECRETS="$DIR/../secrets/remarkable.local.env"

if [ -t 2 ]; then R=$'\033[1;31m'; Z=$'\033[0m'; else R=''; Z=''; fi
err() { printf '%sERROR:%s %s\n' "$R" "$Z" "$*" >&2; }

_get() {
  [ -f "$SECRETS" ] || return 0
  sed -n -E "s/^[[:space:]]*$1[[:space:]]*=[[:space:]]*([^#]*).*/\1/p" "$SECRETS" \
    | head -n1 | sed -E 's/[[:space:]]+$//'
}

PIN_DIGITS="$(_get PIN_DIGITS)"
SYNC_REPO="$(_get SYNC_REPO)"
GH_TOKEN="$(_get GH_TOKEN)"
SYNC_SKIP="$(_get SYNC_SKIP)"

echo "=== configure-sync  target=$TARGET ==="

do_pin=0
do_sync=0
case "$PIN_DIGITS" in
  6|4|none) do_pin=1 ;;
esac
if [ -n "$SYNC_REPO" ] && [ -n "$GH_TOKEN" ]; then
  do_sync=1
elif [ "$SYNC_SKIP" = "1" ] && [ -z "$SYNC_REPO" ]; then
  echo "  Sync skipped in secrets."
elif [ -z "$SYNC_REPO" ]; then
  echo "  No SYNC_REPO in secrets -- sync not pushed."
elif [ -z "$GH_TOKEN" ]; then
  err "SYNC_REPO is set but GH_TOKEN is empty"
  echo "  Run: bash scripts/ensure-secrets.sh" >&2
  exit 1
fi

if [ "$do_pin" -eq 0 ] && [ "$do_sync" -eq 0 ]; then
  echo "  Nothing to push (set PIN_DIGITS and/or SYNC_REPO + GH_TOKEN)."
  exit 0
fi

if ! ping -c1 -W2 "$TARGET" >/dev/null 2>&1; then
  err "$TARGET unreachable"
  exit 1
fi
if ! rm_test_key "$TARGET"; then
  err "SSH key login to root@$TARGET failed"
  exit 1
fi

post_json() {
  local label="$1" remote="$2" url="$3" localf
  localf="$(mktemp)"
  printf '%s' "$remote" >"$localf"
  rm_scp_to "$localf" /tmp/wd-cfg.json "$TARGET"
  rm -f "$localf"
  echo "  $label ..."
  if ! rm_ssh '
    wget -q -O /tmp/wd-cfg.out --header="Content-Type: application/json" \
      --post-file=/tmp/wd-cfg.json '"$url"'
    ec=$?
    rm -f /tmp/wd-cfg.json /tmp/wd-cfg.out
    exit $ec
  ' "$TARGET"; then
    err "$label failed (is writerdeck running?)"
    exit 1
  fi
}

if [ "$do_pin" -eq 1 ]; then
  PIN_JSON="$(PIN_DIGITS="$PIN_DIGITS" python3 - <<'PY'
import json, os
print(json.dumps({"pinDigits": os.environ["PIN_DIGITS"]}))
PY
)"
  post_json "Setting phone PIN length=$PIN_DIGITS" "$PIN_JSON" "http://127.0.0.1:8000/api/settings"
fi

if [ "$do_sync" -eq 1 ]; then
  SETTINGS_JSON="$(SYNC_REPO="$SYNC_REPO" python3 - <<'PY'
import json, os
print(json.dumps({"syncOn": True, "syncRepo": os.environ["SYNC_REPO"]}))
PY
)"
  TOKEN_JSON="$(GH_TOKEN="$GH_TOKEN" python3 - <<'PY'
import json, os
print(json.dumps({"token": os.environ["GH_TOKEN"]}))
PY
)"
  post_json "Setting syncOn + syncRepo=$SYNC_REPO" "$SETTINGS_JSON" "http://127.0.0.1:8000/api/settings"
  post_json "Verifying GitHub token via tablet" "$TOKEN_JSON" "http://127.0.0.1:8000/api/sync/token"
  echo "  Sync configured for $SYNC_REPO (token in tablet RAM only)."
  echo "  Phone browsers may still need the token once per Wi-Fi address (Sync setup)."
fi

echo "  Done."
