#!/usr/bin/env bash
# scripts/test-keyboard-harness.sh -- structured keyboard/selection tests.
#
# Requires Writerdeck with publishEditorState (build-keywriter.sh) and a
# Writerdeck-server build with /api/test/editor-state (deploy-rmkbd.sh).
#
# Usage:
#   bash scripts/test-keyboard-harness.sh              # all scenarios (soft reset)
#   bash scripts/test-keyboard-harness.sh -s NAME      # one scenario
#   bash scripts/test-keyboard-harness.sh -m PREFIX    # name substring match
#   bash scripts/test-keyboard-harness.sh --list       # scenario names
#   bash scripts/test-keyboard-harness.sh --unit       # translate + scenario lint (no device)
#   bash scripts/test-keyboard-harness.sh --fast       # shorter pauses (dev loop)
#   bash scripts/test-keyboard-harness.sh --no-prepare # skip PUT/reload (same scenario re-run)
#   bash scripts/test-keyboard-harness.sh --hard-reset # quit editor per scenario
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
REPO="$(cd "$DIR/.." && pwd)"
TARGET="${RM_HOST}"
SCENARIO=""
MATCH=""
EXTRA=()

while [ $# -gt 0 ]; do
  case "$1" in
    -s) shift; SCENARIO="${1:-}"; shift || true ;;
    -m|--match) shift; MATCH="${1:-}"; shift || true ;;
    --list) EXTRA+=(-list); shift ;;
    --unit) EXTRA+=(-unit); shift ;;
    --hard-reset) EXTRA+=(-hard-reset); shift ;;
    --fast) EXTRA+=(-fast); shift ;;
    --no-prepare) EXTRA+=(-no-prepare); shift ;;
    -v) EXTRA+=(-v); shift ;;
    -h|--help)
      sed -n '2,15p' "$0"
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
  {
    go test -C "$REPO/daemon" -run TestTranslate -v .
    go test -C "$REPO/daemon/cmd/edit-harness" -v .
  } 2>&1 | tee "$LOG"
  exit "${PIPESTATUS[0]}"
fi

ARGS=(-host "$TARGET" -port 8000)
if [ -n "$SCENARIO" ]; then
  ARGS+=(-scenario "$SCENARIO")
fi
if [ -n "$MATCH" ]; then
  ARGS+=(-match "$MATCH")
fi
if [ ${#EXTRA[@]} -gt 0 ]; then
  ARGS+=("${EXTRA[@]}")
fi

{
  echo "=== test-keyboard-harness  $TS  target=$TARGET ==="
  go run -C "$REPO/daemon" ./cmd/edit-harness "${ARGS[@]}"
} 2>&1 | tee "$LOG"
