#!/usr/bin/env bash
# Official Codex SDK 3.27.0.97 (device software 3.27.x).
# SDK_DEVICE=rm1 or rm2 (default rm2 so existing spike scripts stay put).
SDK_DEVICE="${SDK_DEVICE:-rm2}"
SDK_VERSION="${SDK_VERSION:-3.27.0.97}"
SDK_CODEX="${SDK_CODEX:-5.7.119}"
SDK_URL="${SDK_URL:-https://storage.googleapis.com/remarkable-codex-toolchain/${SDK_VERSION}/${SDK_DEVICE}/remarkable-production-image-${SDK_CODEX}-${SDK_DEVICE}-public-x86_64-toolchain.sh}"
SDK_LABEL="${SDK_LABEL:-${SDK_VERSION}-${SDK_DEVICE}}"
export SDK_DEVICE SDK_VERSION SDK_CODEX SDK_URL SDK_LABEL
