#!/usr/bin/env bash
# scripts/test-vault.sh -- loopback vault encrypt/save/PIN session on device or local.
#
# Usage (repo root):
#   bash scripts/test-vault.sh
#   bash scripts/test-vault.sh --local /tmp/wd-vault-test

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
REPO="$(cd "$DIR/.." && pwd)"

LOCAL=""
for arg in "$@"; do
  case "$arg" in
    --local) LOCAL="${2:-/tmp/wd-vault-test}"; shift 2 || true ;;
  esac
done

DEVICE_NOTES="/home/root/Writerdeck-user-documents"

if [[ -n "$LOCAL" ]]; then
  NOTES="$LOCAL/notes"
  SETTINGS="$LOCAL/settings.json"
  mkdir -p "$NOTES"
  PORT=18080
  PIDFILE="$LOCAL/server.pid"
  kill "$(cat "$PIDFILE" 2>/dev/null)" 2>/dev/null || true
  sleep 1
  (cd "$REPO/daemon" && go build -o "$LOCAL/Writerdeck-server" .)
  "$LOCAL/Writerdeck-server" --port "$PORT" --notes-dir "$NOTES" --settings-file "$SETTINGS" >/tmp/wd-vault-local.log 2>&1 &
  echo $! >"$PIDFILE"
  sleep 1
  BASE="http://127.0.0.1:$PORT"
  cleanup() { kill "$(cat "$PIDFILE" 2>/dev/null)" 2>/dev/null || true; rm -f "$PIDFILE"; }
  trap cleanup EXIT
else
  BASE="http://$RM_HOST:8000"
  NOTES="$DEVICE_NOTES"
fi

fail() { echo "FAIL: $*" >&2; exit 1; }

echo "=== test-vault base=$BASE ==="

# Remove harness notes before vault reset (disable refuses user .md.enc on disk).
if [[ -n "$LOCAL" ]]; then
  rm -f "$NOTES"/z-test-vault-plain.md "$NOTES"/z-test-vault-plain.md.enc
else
  ssh -o StrictHostKeyChecking=accept-new -o ConnectTimeout=8 "root@$RM_HOST" \
    "rm -f $DEVICE_NOTES/z-test-vault-plain.md $DEVICE_NOTES/z-test-vault-plain.md.enc"
fi

# Deterministic PIN: reset vault so a prior E2E run cannot leave a different PIN.
curl -s -o /dev/null -X POST "$BASE/api/test/tablet-req" \
  -H 'Content-Type: application/json' \
  -d '{"op":"disablevault"}'

# Setup vault via test tablet-req endpoint
code=$(curl -s -o /tmp/vault-setup.json -w '%{http_code}' -X POST "$BASE/api/test/tablet-req" \
  -H 'Content-Type: application/json' \
  -d '{"op":"setvaultpin","name":"123456"}')
[[ "$code" == "200" ]] || fail "setvaultpin HTTP $code"

TEST_NOTE="z-test-vault-plain.md"
TEST_ENC="${TEST_NOTE%.md}.md.enc"

if [[ -n "$LOCAL" ]]; then
  echo "test note" >"$NOTES/$TEST_NOTE"
else
  ssh -o StrictHostKeyChecking=accept-new -o ConnectTimeout=8 "root@$RM_HOST" \
    "echo 'test note' > $DEVICE_NOTES/$TEST_NOTE"
fi

code=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$BASE/api/test/tablet-req" \
  -H 'Content-Type: application/json' \
  -d '{"op":"verifyvaultpin","name":"123456","old":"once"}')
[[ "$code" == "200" ]] || fail "verifyvaultpin HTTP $code"

code=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$BASE/api/test/tablet-req" \
  -H 'Content-Type: application/json' \
  -d "{\"op\":\"encryptnote\",\"name\":\"$TEST_NOTE\"}")
[[ "$code" == "200" ]] || fail "encryptnote HTTP $code"

if [[ -n "$LOCAL" ]]; then
  [[ -f "$NOTES/$TEST_ENC" ]] || fail "$TEST_ENC missing"
  [[ ! -f "$NOTES/$TEST_NOTE" ]] || fail "$TEST_NOTE should be gone"
else
  ssh -o StrictHostKeyChecking=accept-new -o ConnectTimeout=8 "root@$RM_HOST" \
    "test -f $DEVICE_NOTES/$TEST_ENC" || fail "$TEST_ENC missing on device"
  ssh -o StrictHostKeyChecking=accept-new -o ConnectTimeout=8 "root@$RM_HOST" \
    "test ! -f $DEVICE_NOTES/$TEST_NOTE" || fail "$TEST_NOTE should be gone on device"
fi

if [[ -n "$LOCAL" ]]; then
  code=$(curl -s -o /dev/null -w '%{http_code}' "$BASE/api/notes/$TEST_ENC/download")
  [[ "$code" == "423" ]] || fail "download without PIN want 423 got $code"
else
  code=$(curl -s -o /dev/null -w '%{http_code}' "$BASE/api/notes/$TEST_ENC/download")
  [[ "$code" == "401" || "$code" == "423" ]] || fail "download without PIN want 401 or 423 got $code"
fi

code=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$BASE/api/test/tablet-req" \
  -H 'Content-Type: application/json' \
  -d '{"op":"verifyvaultpin","name":"123456","old":"once"}')
[[ "$code" == "200" ]] || fail "verifyvaultpin for download HTTP $code"

status=$(curl -s "$BASE/api/vault/status")
echo "$status" | grep -q '"enabled":true' || fail "vault not enabled after setup"
echo "$status" | grep -q '"locked":false' || fail "vault still needs PIN after verify"

if [[ -n "$LOCAL" ]]; then
  plain=$(curl -s "$BASE/api/notes/$TEST_ENC")
  [[ "$plain" == *"test note"* ]] || fail "loopback GET decrypt failed: $plain"
else
  # Loopback read is tablet-only; verify ciphertext is not plaintext on disk.
  ssh -o StrictHostKeyChecking=accept-new -o ConnectTimeout=8 "root@$RM_HOST" \
    "dd if=$DEVICE_NOTES/$TEST_ENC bs=1 count=6 2>/dev/null | grep -q WDENC1" || fail "encrypted file missing WDENC1 magic"
  ssh -o StrictHostKeyChecking=accept-new -o ConnectTimeout=8 "root@$RM_HOST" \
    "grep -q 'test note' $DEVICE_NOTES/$TEST_ENC && exit 1 || exit 0" || fail "plaintext visible in .md.enc on disk"
fi

echo "PASS: test-vault"
