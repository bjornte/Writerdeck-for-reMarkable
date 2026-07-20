#!/usr/bin/env bash
# scripts/test-vault-e2e.sh -- vault UI + keyboard PIN + GitHub sync on device.
#
# Usage (repo root):
#   bash scripts/test-vault-e2e.sh
#   bash scripts/test-vault-e2e.sh 192.168.1.8
#   bash scripts/test-vault-e2e.sh --skip-cleanup

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
REPO="$(cd "$DIR/.." && pwd)"
RECON_DIR="$REPO/docs/recon"
mkdir -p "$RECON_DIR"

TARGET="${RM_HOST}"
SKIP_CLEANUP=""
for arg in "$@"; do
  case "$arg" in
    --skip-cleanup) SKIP_CLEANUP=1 ;;
    -*) ;;
    *) TARGET="$arg" ;;
  esac
done

TS="$(date +%Y-%m-%dT%H-%M-%S)"
LOG="$RECON_DIR/test-vault-e2e-$TS.txt"

ARGS=(-host "$TARGET")
[[ -n "$SKIP_CLEANUP" ]] && ARGS+=(-skip-cleanup)

{
  echo "=== test-vault-e2e  $TS  target=$TARGET ==="
  go run -C "$REPO/daemon" ./cmd/vault-e2e-test "${ARGS[@]}"
} 2>&1 | tee "$LOG"

if grep -q '^PASS$' "$LOG"; then
  echo ""
  echo "=== verdict: PASS ==="
else
  echo ""
  echo "=== verdict: FAIL ==="
  exit 1
fi
echo "Full log: $LOG"
