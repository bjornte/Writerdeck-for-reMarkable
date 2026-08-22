#!/usr/bin/env bash
# Shared Codex SDK helpers for spike build scripts.
set -euo pipefail

_spike_sdk_prepare() {
  local sdk_dir="$1"
  if ! ls "$sdk_dir"/environment-setup-* >/dev/null 2>&1 \
    || [ ! -d "$sdk_dir/sysroots" ]; then
    return 1
  fi
  # Installer hardcodes /usr/local/oe-sdk-hardcoded-buildpath in env scripts.
  mkdir -p /usr/local
  ln -sfn "$sdk_dir" /usr/local/oe-sdk-hardcoded-buildpath
  return 0
}

_spike_sdk_env() {
  local sdk_dir="$1"
  _spike_sdk_prepare "$sdk_dir" || return 1
  local setup
  setup="$(ls "$sdk_dir"/environment-setup-* | head -n1)"
  # shellcheck disable=SC1090
  set +u
  # shellcheck source=/dev/null
  . "$setup"
  set -u
}

# Configure + build one cmake Qt app with the Codex OE toolchain.
# Usage: _spike_cmake_build SDK_DIR SRC_DIR BUILD_DIR OUT_BINARY
_spike_cmake_build() {
  local sdk_dir="$1" src="$2" build_dir="$3" out_bin="$4"
  local toolchain native_cmake jobs base
  _spike_sdk_env "$sdk_dir" || return 1
  toolchain="$sdk_dir/sysroots/x86_64-codexsdk-linux/usr/share/cmake/OEToolchainConfig.cmake"
  native_cmake="$sdk_dir/sysroots/x86_64-codexsdk-linux/usr/bin/cmake"
  base="$(basename "$out_bin")"
  rm -rf "$build_dir"
  "$native_cmake" -S "$src" -B "$build_dir" \
    -DCMAKE_TOOLCHAIN_FILE="$toolchain" -G Ninja
  jobs="$(nproc 2>/dev/null || echo 4)"
  "$native_cmake" --build "$build_dir" -j"$jobs"
  cp -f "$build_dir/$base" "$out_bin"
}

# Mac/ARM: build inside linux/amd64 Docker with SDK in container /tmp (not host mount).
# Usage: _spike_docker_build SPIKE_DIR INSTALLER_PATH build-script-name
_spike_docker_build() {
  local spike="$1" installer="$2" script="$3"
  if ! docker info >/dev/null 2>&1; then
    echo "ERROR: Docker is not running." >&2
    echo "  Start Colima (x86_64 for the official SDK):" >&2
    echo "    colima start --arch x86_64 --cpu 4 --memory 8 --disk 60" >&2
    exit 1
  fi
  local inst_base
  inst_base="$(basename "$installer")"
  docker run --rm --platform linux/amd64 \
    -v "$spike:/spike" \
    -w /spike \
    ubuntu:22.04 \
    bash -lc "set -euo pipefail
      export DEBIAN_FRONTEND=noninteractive
      apt-get update -qq
      apt-get install -y -qq file python3 xz-utils
      export SPIKE_SDK_DIR=/tmp/codex-sdk-rm2
      rm -rf \"\$SPIKE_SDK_DIR\"
      bash .cache/${inst_base} -d \"\$SPIKE_SDK_DIR\" -y
      bash scripts/${script} --inside-container"
}
