#!/usr/bin/env bash
# Build socket_spike + socket_inject for rM2 using the Codex SDK.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=/dev/null
. "$SPIKE/scripts/sdk-url.sh"
# shellcheck source=/dev/null
. "$SPIKE/scripts/_spike-sdk.sh"

CACHE="$SPIKE/.cache"
SDK_DIR="${SPIKE_SDK_DIR:-$CACHE/sdk-$SDK_LABEL}"
INSTALLER="$CACHE/$(basename "$SDK_URL")"
SRC="$SPIKE/socket_spike"
OUT="$CACHE/out"
BUILD="$OUT/build-socket"
mkdir -p "$CACHE" "$OUT"

need_docker() {
  if [ "${1:-}" = "--inside-container" ]; then
    return 1
  fi
  [ "$(uname -s)" = "Darwin" ] && return 0
  [ "$(uname -m)" != "x86_64" ] && return 0
  return 1
}

download_sdk() {
  if [ -f "$INSTALLER" ]; then
    echo "SDK installer already present: $INSTALLER"
    return 0
  fi
  echo "Downloading $SDK_URL ..."
  curl -fL --retry 3 -o "$INSTALLER.partial" "$SDK_URL"
  mv "$INSTALLER.partial" "$INSTALLER"
  chmod +x "$INSTALLER"
}

install_sdk_linux() {
  if ls "$SDK_DIR"/environment-setup-* >/dev/null 2>&1 \
    && [ -d "$SDK_DIR/sysroots" ]; then
    echo "SDK already installed under $SDK_DIR"
    return 0
  fi
  mkdir -p "$SDK_DIR"
  echo "Installing SDK into $SDK_DIR (non-interactive) ..."
  "$INSTALLER" -d "$SDK_DIR" -y
}

build_linux() {
  _spike_cmake_build "$SDK_DIR" "$SRC" "$BUILD" "$OUT/socket_spike" \
    || { echo "ERROR: SDK build failed" >&2; exit 1; }
  cp -f "$BUILD/socket_inject" "$OUT/socket_inject"
  file "$OUT/socket_spike"
  file "$OUT/socket_inject"
  echo "OK: $OUT/socket_spike $OUT/socket_inject"
}

build_via_docker() {
  _spike_docker_build "$SPIKE" "$INSTALLER" "$(basename "$0")"
}

MODE="${1:-}"
if need_docker "$MODE"; then
  download_sdk
  build_via_docker
else
  download_sdk
  install_sdk_linux
  build_linux
fi
