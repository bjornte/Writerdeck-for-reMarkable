#!/usr/bin/env bash
# scripts/test-keyboard-harness.sh -- structured keyboard/selection tests.
#
# Requires Writerdeck with publishEditorState (build-keywriter.sh) and a
# Writerdeck-server build with /api/test/editor-state (deploy-rmkbd.sh).
#
# Usage:
#   bash scripts/test-keyboard-harness.sh              # all scenarios (soft reset)
#   bash scripts/test-keyboard-harness.sh -s NAME      # one scenario
#   bash scripts/test-keyboard-harness.sh --list       # scenario names
#   bash scripts/test-keyboard-harness.sh --unit       # translate tests only
#   bash scripts/test-keyboard-harness.sh --hard-reset # quit editor per scenario
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
REPO="$(cd "$DIR/.." && pwd)"
TARGET="${RM_HOST}"
SCENARIO=""
EXTRA=()

while [ $# -gt 0 ]; do
  case "$1" in
    -s) shift; SCENARIO="${1:-}"; shift || true ;;
    --list) EXTRA+=(-list); shift ;;
    --unit) EXTRA+=(-unit); shift ;;
    --hard-reset) EXTRA+=(-hard-reset); shift ;;
    -v) EXTRA+=(-v); shift ;;
    -h|--help)
      sed -n '2,12p' "$0"
      exit 0
      ;;
    *) TARGET="$1"; shift ;;
  esac
done

RECON_DIR="$REPO/docs/recon"
mkdir -p "$RECON_DIR"
TS="$(date +%Y-%m-%dT%H-%M-%S)"
LOG="$RECON_DIR/test-keyboard-harness-$TS.txt"

if printf '%s\n' "${EXTRA[@]}" | grep -q -- '-unit'; then
  exec go test -C "$REPO/daemon" -run TestTranslate -v . 2>&1 | tee "$LOG"
fi

ARGS=(-host "$TARGET" -port 8000)
if [ -n "$SCENARIO" ]; then
  ARGS+=(-scenario "$SCENARIO")
fi
if [ ${#EXTRA[@]} -gt 0 ]; then
  ARGS+=("${EXTRA[@]}")
fi

{
  echo "=== test-keyboard-harness  $TS  target=$TARGET ==="
  go run -C "$REPO/daemon" ./cmd/edit-harness "${ARGS[@]}"
} 2>&1 | tee "$LOG"
