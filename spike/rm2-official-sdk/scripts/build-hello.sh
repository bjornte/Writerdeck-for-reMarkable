#!/usr/bin/env bash
# Build official hello_remarkable for rM2 using the Codex SDK.
# On macOS / non-x86_64: requires Docker (Colima) and runs the Linux x86_64 SDK
# inside ubuntu:22.04 --platform linux/amd64.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=/dev/null
. "$SPIKE/scripts/sdk-url.sh"

CACHE="$SPIKE/.cache"
SDK_DIR="$CACHE/sdk-$SDK_LABEL"
INSTALLER="$CACHE/$(basename "$SDK_URL")"
SRC="$SPIKE/hello_remarkable"
OUT="$CACHE/out"
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
  if ls "$SDK_DIR"/environment-setup-* >/dev/null 2>&1; then
    echo "SDK already installed under $SDK_DIR"
    return 0
  fi
  mkdir -p "$SDK_DIR"
  echo "Installing SDK into $SDK_DIR (non-interactive) ..."
  "$INSTALLER" -d "$SDK_DIR" -y
}

build_linux() {
  local setup jobs
  setup="$(ls "$SDK_DIR"/environment-setup-* | head -n1)"
  # shellcheck disable=SC1090
  set +u
  # shellcheck source=/dev/null
  . "$setup"
  set -u
  rm -rf "$OUT/build"
  cmake -S "$SRC" -B "$OUT/build"
  jobs="$(nproc 2>/dev/null || echo 4)"
  cmake --build "$OUT/build" -j"$jobs"
  cp -f "$OUT/build/hello_remarkable" "$OUT/hello_remarkable"
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
  if ! docker info >/dev/null 2>&1; then
    echo "ERROR: Docker is not running." >&2
    echo "  Start Colima (x86_64 for the official SDK):" >&2
    echo "    colima start --arch x86_64 --cpu 4 --memory 8 --disk 60" >&2
    exit 1
  fi
  download_sdk
  docker run --rm --platform linux/amd64 \
    -v "$SPIKE:/spike" \
    -w /spike \
    ubuntu:22.04 \
    bash -lc 'set -euo pipefail
      export DEBIAN_FRONTEND=noninteractive
      apt-get update -qq
      apt-get install -y -qq curl ca-certificates cmake ninja-build file python3 xz-utils build-essential
      bash scripts/build-hello.sh --inside-container'
}

MODE="${1:-}"
if need_docker "$MODE"; then
  build_via_docker
else
  download_sdk
  install_sdk_linux
  build_linux
fi
