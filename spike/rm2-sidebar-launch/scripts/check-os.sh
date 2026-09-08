#!/usr/bin/env bash
# Print device OS and exit 0 only if it is in tested-os.json "tested".
# Usage: check-os.sh [--json PATH] [host]
# Env: ACCEPT_UNTESTED=1 allows candidates (demo first-run only).
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"

JSON="${SPIKE}/tested-os.json"
HOST=""
while [ $# -gt 0 ]; do
  case "$1" in
    --json) JSON="$2"; shift 2 ;;
    *) HOST="$1"; shift ;;
  esac
done
HOST="${HOST:-$(bash "$ROOT/spike/rm2-official-sdk/scripts/rm2-pick-host.sh")}"

OS_LINE="$(rm_ssh 'cat /etc/os-release; echo ---; cat /etc/version 2>/dev/null' "$HOST" | tr -d "\r")"
# Prefer software IMG_VERSION (e.g. 3.27.3.0), not Codex VERSION_ID (e.g. 5.7.126).
OS_VER="$(printf '%s\n' "$OS_LINE" | sed -n 's/^IMG_VERSION=//p' | head -1 | tr -d '"')"
if [ -z "$OS_VER" ]; then
  OS_VER="$(printf '%s\n' "$OS_LINE" | sed -n 's/^VERSION_ID=//p' | head -1 | tr -d '"')"
fi
if [ -z "$OS_VER" ]; then
  OS_VER="$(printf '%s\n' "$OS_LINE" | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+(\.[0-9]+)?' | head -1 || true)"
fi
if [ -z "$OS_VER" ]; then
  err "could not read OS version from $HOST"
  exit 2
fi

python3 - "$JSON" "$OS_VER" <<'PY'
import json, os, sys
path, ver = sys.argv[1], sys.argv[2]
data = json.load(open(path))
tested = set(data.get("tested") or [])
cands = set(data.get("candidates") or [])
print(f"os={ver}")
print(f"tested={','.join(sorted(tested)) or '(none)'}")
if ver in tested:
    print("status=tested")
    sys.exit(0)
if os.environ.get("ACCEPT_UNTESTED") == "1" and ver in cands:
    print("status=candidate-accepted")
    sys.exit(0)
print("status=untested")
print(f"refuse: {ver} is not in tested-os.json (set ACCEPT_UNTESTED=1 for candidates only)")
sys.exit(1)
PY
