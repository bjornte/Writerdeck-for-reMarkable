#!/usr/bin/env bash
# scripts/fetch-observation.sh -- Pull the latest Observe export from the tablet.
# Writes docs/recon/last-observation.json for Cursor / humans.
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=scripts/_env.sh
source "$DIR/_env.sh"

OUT="${1:-$DIR/../docs/recon/last-observation.json}"
mkdir -p "$(dirname "$OUT")"

if ! STATUS="$(curl -sf "http://${RM_HOST}:8000/api/observe/status")"; then
  echo "ERROR: cannot reach http://${RM_HOST}:8000/api/observe/status" >&2
  exit 1
fi

HAS="$(printf '%s' "$STATUS" | python3 -c 'import sys,json; print(json.load(sys.stdin).get("hasExport", False))')"
if [ "$HAS" != "True" ]; then
  echo "ERROR: no observation on the tablet yet (Stop observe first)." >&2
  exit 1
fi

curl -sf "http://${RM_HOST}:8000/api/observe/export" -o "$OUT"
echo "Wrote $OUT"
