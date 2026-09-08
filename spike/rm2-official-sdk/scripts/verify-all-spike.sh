#!/usr/bin/env bash
# Run all spike autonomous checks. Exit 1 on first failure.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
for app in hello textedit socket; do
  echo "=== verify $app ==="
  bash "$SPIKE/scripts/verify-spike.sh" "$app"
done
echo "ALL SPIKE CHECKS PASS"
