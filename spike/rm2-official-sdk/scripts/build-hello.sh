#!/usr/bin/env bash
# Build official hello_remarkable for rM2 using the Codex SDK.
# On macOS / non-x86_64: requires Docker (Colima) and runs the Linux x86_64 SDK
# inside ubuntu:22.04 --platform linux/amd64.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=/dev/null
. "$SPIKE/scripts/sdk-url.sh"
# shellcheck source=/dev/null
. "$SPIKE/scripts/_spike-sdk.sh"

CACHE="$SPIKE/.cache"
SDK_DIR="${SPIKE_SDK_DIR:-$CACHE/sdk-$SDK_LABEL}"
INSTALLER="$CACHE/$(basename "$SDK_URL")"
SRC="$SPIKE/hello_remarkable"
if [ "${SDK_DEVICE:-rm2}" = "rm1" ]; then
  OUT="$CACHE/out-rm1"
else
  OUT="$CACHE/out"
fi
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
  _spike_cmake_build "$SDK_DIR" "$SRC" "$OUT/build" "$OUT/hello_remarkable" \
    || { echo "ERROR: SDK build failed" >&2; exit 1; }
  local sg
  sg="$(find "$SDK_DIR" -path '*/plugins/scenegraph/libqsgepaper.so' 2>/dev/null | head -n1 || true)"
  if [ -n "$sg" ]; then
    cp -f "$sg" "$OUT/libqsgepaper.so"
    echo "Copied libqsgepaper.so from SDK"
  fi
  file "$OUT/hello_remarkable"
  echo "OK: $OUT/hello_remarkable"
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
