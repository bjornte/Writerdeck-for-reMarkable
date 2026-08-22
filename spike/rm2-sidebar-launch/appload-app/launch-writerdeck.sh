#!/bin/sh
# AppLoad external entry: ask Writerdeck-server to open Lobby (same path as phone).
# Runs inside xochitl; /api/lobby stops stock UI and starts Writerdeck when idle.
set -eu

try_post() {
  url="$1"
  if command -v wget >/dev/null 2>&1; then
    wget -qO- --post-data='' --header='Content-Type: application/json' "$url" && return 0
  fi
  if command -v curl >/dev/null 2>&1; then
    curl -sf -X POST -H 'Content-Type: application/json' "$url" && return 0
  fi
  return 1
}

if try_post 'http://127.0.0.1:8000/api/lobby'; then
  exit 0
fi
if try_post 'http://127.0.0.1:8000/api/launch'; then
  exit 0
fi

echo "launch-writerdeck: Writerdeck-server not reachable on :8000" >&2
exit 1
