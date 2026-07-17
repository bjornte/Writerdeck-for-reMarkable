#!/usr/bin/env bash
# scripts/install.sh -- First-time (idempotent) install chain for Writerdeck.
#
# Runs: ensure-secrets -> preflight -> bootstrap -> fetch if needed -> deploy
# editor -> deploy server -> install-service. With --start: also start, health-
# check the phone UI, then enable autostart on boot.
#
# Usage (run from repo root on the Mac):
#   bash scripts/install.sh
#   bash scripts/install.sh --start    # start + enable after health check
#
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"

DO_START=0
for arg in "$@"; do
  case "$arg" in
    --start) DO_START=1 ;;
    -h|--help)
      sed -n '2,12p' "$0"
      exit 0
      ;;
    *)
      echo "ERROR: unknown arg: $arg (try --start)" >&2
      exit 1
      ;;
  esac
done

cd "$REPO"

echo "=== install.sh ==="
echo

bash "$DIR/ensure-secrets.sh"
bash "$DIR/preflight.sh" --skip-dist

echo
echo "--- bootstrap (SSH key) ---"
bash "$DIR/bootstrap.sh"

echo
echo "--- editor binary ---"
bash "$DIR/preflight.sh" --fetch

echo
echo "--- deploy Writerdeck (editor + Qt) ---"
bash "$DIR/deploy-keywriter.sh"

echo
echo "--- deploy Writerdeck-server ---"
bash "$DIR/deploy-rmkbd.sh"

echo
echo "--- systemd unit ---"
if [ "$DO_START" -eq 1 ]; then
  bash "$DIR/install-service.sh" --start
else
  bash "$DIR/install-service.sh"
fi

# shellcheck source=/dev/null
. "$DIR/_env.sh"

echo
echo "======================================"
echo "  INSTALL CHAIN DONE"
echo "======================================"
echo ""
echo "  Phone UI:  http://${RM_HOST}:8000/"
echo ""
if [ "$DO_START" -ne 1 ]; then
  echo "  Next -- start + enable autostart:"
  echo "    bash scripts/install.sh --start"
  echo "  Or: bash scripts/install-service.sh --start"
  echo ""
fi
echo "  You're done when: Lobby on e-ink; phone list populated;"
echo "  connection bar not stuck on connecting..."
echo ""
echo "  Recovery:"
echo "    systemctl disable --now writerdeck && systemctl start xochitl"
echo "======================================"
