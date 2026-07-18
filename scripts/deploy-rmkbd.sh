#!/usr/bin/env bash
# scripts/deploy-rmkbd.sh -- Build or fetch Writerdeck-server (ARMv7 static) and deploy.
# Run from the repo root on the Mac (the only machine that can reach the tablet).
#
# Usage:
#   bash scripts/deploy-rmkbd.sh               # build-or-fetch + deploy
#   bash scripts/deploy-rmkbd.sh --build-only  # just obtain binary; no device connection
#   bash scripts/deploy-rmkbd.sh --deploy-only # deploy existing Writerdeck-server binary
#
# Prefer local go build when go is installed (dev loop). Otherwise curl the
# rolling GitHub Release tag "server" (visitors need no go).
#
# Requires: go 1.25+ OR network to download the Release binary.

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"

BUILD_ONLY=0
DEPLOY_ONLY=0
for arg in "$@"; do
  case "$arg" in
    --build-only)  BUILD_ONLY=1 ;;
    --deploy-only) DEPLOY_ONLY=1 ;;
    -h|--help)
      sed -n '2,14p' "$0"
      exit 0
      ;;
    *)
      echo "ERROR: unknown arg: $arg (try --build-only or --deploy-only)" >&2
      exit 1
      ;;
  esac
done
if [ "$BUILD_ONLY" -eq 1 ] && [ "$DEPLOY_ONLY" -eq 1 ]; then
  echo "ERROR: use only one of --build-only / --deploy-only" >&2
  exit 1
fi

BINARY="${REPO}/Writerdeck-server"

if [ "$DEPLOY_ONLY" -eq 0 ]; then
  echo "=== Writerdeck-server: obtain ARMv7 binary ==="
  # Prefer local go build for the full-repo dev loop. Slim installer ZIPs omit
  # daemon/, so fall through to the Release fetch even when go is installed.
  if command -v go >/dev/null 2>&1 && go version >/dev/null 2>&1 \
    && [ -d "${REPO}/daemon" ] && [ -f "${REPO}/daemon/go.mod" ]; then
      cd "${REPO}/daemon"
      go mod tidy
      VER="unknown"
      if [ -f "${REPO}/VERSION" ]; then
        VER="$(tr -d '[:space:]' < "${REPO}/VERSION")"
      fi
      GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 \
          go build -trimpath \
            -ldflags "-X main.productVersion=${VER}" \
            -o "${BINARY}" .
      echo "  built with go (version ${VER}): $(file "${BINARY}")"
  else
      if command -v go >/dev/null 2>&1 && [ ! -d "${REPO}/daemon" ]; then
        echo "  go present but no daemon/ -- fetching Release binary (installer package)"
      else
        echo "  go not found -- fetching Release binary"
      fi
      bash "${DIR}/fetch-server-dist.sh"
  fi
  echo

  if [ "$BUILD_ONLY" -eq 1 ]; then
      echo "  --build-only: skipping device deploy."
      exit 0
  fi
fi

if [ ! -f "$BINARY" ]; then
  echo "ERROR: missing $BINARY (run without --deploy-only, or --build-only first)" >&2
  exit 1
fi

# shellcheck source=/dev/null
. "${DIR}/_env.sh"
# shellcheck source=/dev/null
. "${DIR}/migrate-device-layout.sh"

# Flush open editor buffer before stopping the server (slice 11). deploy used to
# pkill and sleep 0.5s -- not long enough for saveAndQuit / loopback PUT to finish.
rm_graceful_stop_server() {
  rm_ssh '
    wget -q -O /dev/null --post-data="" http://127.0.0.1:8000/api/flush-save 2>/dev/null || true
    for p in $(pidof Writerdeck-server 2>/dev/null); do kill -TERM "$p" 2>/dev/null; done
    i=0
    while pidof Writerdeck-server >/dev/null 2>&1 && [ "$i" -lt 60 ]; do
      sleep 0.2
      i=$((i + 1))
    done
    for p in $(pidof Writerdeck-server 2>/dev/null); do kill -KILL "$p" 2>/dev/null; done
    sleep 0.3
    true
  ' >/dev/null
}

echo "=== Writerdeck-server: deploying to ${RM_HOST} ==="
migrate_device_layout
rm_graceful_stop_server
echo "  any old server stopped."
printf '    '; with_ticker 5 rm_send_file "${BINARY}" "${DEVICE_SERVER}.new"
rm_ssh "mv -f ${DEVICE_SERVER}.new ${DEVICE_SERVER} && chmod +x ${DEVICE_SERVER}" >/dev/null
echo "  ${DEVICE_SERVER} updated."
rm_deploy_wd
echo "  ${DEVICE_WD} updated."
echo

echo "======================================"
echo "  DEPLOY DONE"
echo "======================================"
echo "  Binary : ${DEVICE_SERVER}"
echo "  Next   : bash scripts/test-e2e.sh"
echo "======================================"
