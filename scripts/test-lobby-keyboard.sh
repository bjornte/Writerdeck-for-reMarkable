#!/usr/bin/env bash
# scripts/test-lobby-keyboard.sh -- lobby USB/WS keys after return from edit.
#
# Opens a note, POST /api/lobby (same isLobby transition as Home), sends
# Files-tab + Enter over WebSocket, asserts the note reopens.
#
# Requires Writerdeck with publishEditorState and /api/test/editor-state.
# Run after Writerdeck QML deploy (deploy-keywriter.sh -b).
#
# Usage:
#   bash scripts/test-lobby-keyboard.sh
#   bash scripts/test-lobby-keyboard.sh 192.168.1.8
#   bash scripts/test-lobby-keyboard.sh -note z-test-keyboard-harness.md
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
REPO="$(cd "$DIR/.." && pwd)"

TARGET="${RM_HOST}"
NOTE="z-test-keyboard-harness.md"
while [ $# -gt 0 ]; do
  case "$1" in
    -note) shift; NOTE="${1:-z-test-keyboard-harness.md}"; shift || true ;;
    -h|--help)
      sed -n '2,14p' "$0"
      exit 0
      ;;
    *) TARGET="$1"; shift ;;
  esac
done

RECON_DIR="$REPO/docs/recon"
mkdir -p "$RECON_DIR"
TS="$(date +%Y-%m-%dT%H-%M-%S)"
LOG="$RECON_DIR/test-lobby-keyboard-$TS.txt"

{
  echo "=== test-lobby-keyboard  $TS  target=$TARGET  note=$NOTE ==="
  go run -C "$REPO/daemon" ./cmd/lobby-keyboard-test -host "$TARGET" -note "$NOTE"
} 2>&1 | tee "$LOG"
