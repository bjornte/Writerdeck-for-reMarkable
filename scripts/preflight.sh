#!/usr/bin/env bash
# scripts/preflight.sh -- First-install checks before bootstrap/deploy.
#
# Checks: secrets file + password + Wi-Fi IP, tablet ping, editor dist artifacts.
# Go is optional (Release binary used when missing). Does not print secret values.
#
# Usage (run from repo root on the Mac):
#   bash scripts/preflight.sh
#   bash scripts/preflight.sh --skip-dist   # secrets/ping only
#   bash scripts/preflight.sh --fetch       # if dist missing, run fetch-keywriter-dist.sh
#
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
SECRETS="$REPO/secrets/remarkable.local.env"
DIST="$REPO/third_party/keywriter/dist"

SKIP_DIST=0
DO_FETCH=0
for arg in "$@"; do
  case "$arg" in
    --skip-dist) SKIP_DIST=1 ;;
    --fetch)     DO_FETCH=1 ;;
    -h|--help)
      sed -n '2,12p' "$0"
      exit 0
      ;;
    *)
      echo "ERROR: unknown arg: $arg (try --skip-dist or --fetch)" >&2
      exit 1
      ;;
  esac
done

if [ -t 2 ]; then R=$'\033[1;31m'; Z=$'\033[0m'; else R=''; Z=''; fi
err() { printf '%sERROR:%s %s\n' "$R" "$Z" "$*" >&2; }
ok()  { printf '  OK  %s\n' "$*"; }
fail=0

echo "=== preflight ==="
echo

bash "$DIR/ensure-secrets.sh"

# --- secrets file ------------------------------------------------------------
if [ ! -f "$SECRETS" ]; then
  err "missing $SECRETS after ensure-secrets"
  exit 1
fi
ok "secrets file present"

_get() {
  sed -n -E "s/^[[:space:]]*$1[[:space:]]*=[[:space:]]*([^#]*).*/\1/p" "$SECRETS" \
    | head -n1 | sed -E 's/[[:space:]]+$//'
}

PASS="$(_get RM_ROOT_PASSWORD)"
WIFI="$(_get RM_HOST_WIFI)"
if [ -z "$PASS" ]; then
  err "RM_ROOT_PASSWORD is empty in secrets/remarkable.local.env"
  fail=1
else
  ok "RM_ROOT_PASSWORD is set"
fi
if [ -z "$WIFI" ]; then
  err "RM_HOST_WIFI is empty (tablet Wi-Fi IP required on Mac)"
  fail=1
else
  ok "RM_HOST_WIFI is set ($WIFI)"
fi
if [ "$fail" -ne 0 ]; then
  exit 1
fi

# shellcheck source=/dev/null
. "$DIR/_env.sh"

# --- go (optional) -----------------------------------------------------------
if command -v go >/dev/null 2>&1 && GO_VER="$(go version 2>/dev/null)"; then
  ok "$GO_VER (local server builds)"
  GO_NUM="$(printf '%s\n' "$GO_VER" | sed -n -E 's/.*go([0-9]+\.[0-9]+).*/\1/p')"
  if [ -n "$GO_NUM" ]; then
    GO_MAJOR="${GO_NUM%%.*}"
    GO_MINOR="${GO_NUM#*.}"
    if [ "$GO_MAJOR" -lt 1 ] || { [ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 25 ]; }; then
      echo "  WARN  daemon/go.mod wants go 1.25+; found $GO_NUM" >&2
    fi
  fi
else
  ok "go not installed -- install will download Writerdeck-server from Releases"
fi

# --- ping --------------------------------------------------------------------
TARGET="$RM_HOST"
echo "--- ping $TARGET ---"
if ! ping -c1 -W2 "$TARGET" >/dev/null 2>&1; then
  err "$TARGET unreachable (tablet awake? correct RM_HOST_WIFI? same Wi-Fi?)"
  exit 1
fi
ok "tablet reachable"

# --- dist --------------------------------------------------------------------
if [ "$SKIP_DIST" -eq 1 ]; then
  echo "--- dist (skipped) ---"
  echo
  echo "preflight OK (dist not checked)"
  exit 0
fi

echo "--- editor dist ---"
need_fetch=0
if [ ! -f "$DIST/Writerdeck" ] || [ ! -f "$DIST/qt5.tar.gz" ]; then
  need_fetch=1
fi

if [ "$need_fetch" -eq 1 ]; then
  if [ "$DO_FETCH" -eq 1 ]; then
    echo "  dist incomplete -- running fetch-keywriter-dist.sh"
    bash "$DIR/fetch-keywriter-dist.sh"
  else
    err "missing $DIST/Writerdeck and/or qt5.tar.gz"
    echo "  Run: bash scripts/fetch-keywriter-dist.sh" >&2
    echo "  Or:  bash scripts/preflight.sh --fetch" >&2
    exit 1
  fi
fi

if [ ! -f "$DIST/Writerdeck" ] || [ ! -f "$DIST/qt5.tar.gz" ]; then
  err "dist still incomplete after fetch"
  exit 1
fi
ok "Writerdeck + qt5.tar.gz present"
echo
echo "preflight OK"
