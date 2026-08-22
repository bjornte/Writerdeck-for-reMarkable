#!/usr/bin/env bash
# After a successful sidebar smoke test, record the device OS in tested-os.json.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"

HOST="${1:-}"
if [ -z "$HOST" ]; then
  HOST="$(bash "$ROOT/spike/rm2-official-sdk/scripts/rm2-pick-host.sh")"
fi
JSON="$SPIKE/tested-os.json"

OS_VER="$(rm_ssh 'grep ^IMG_VERSION= /etc/os-release | head -1 | cut -d= -f2 | tr -d "\""' "$HOST" | tr -d '\r\n')"
if [ -z "$OS_VER" ]; then
  OS_VER="$(rm_ssh 'grep ^VERSION_ID= /etc/os-release | head -1 | cut -d= -f2 | tr -d "\""' "$HOST" | tr -d '\r\n')"
fi
if [ -z "$OS_VER" ]; then
  err "empty IMG_VERSION/VERSION_ID from $HOST"
  exit 1
fi

python3 - "$JSON" "$OS_VER" <<'PY'
import json, sys
path, ver = sys.argv[1], sys.argv[2]
data = json.load(open(path))
tested = list(data.get("tested") or [])
cands = list(data.get("candidates") or [])
if ver not in tested:
    tested.append(ver)
if ver in cands:
    cands = [c for c in cands if c != ver]
data["tested"] = tested
data["candidates"] = cands
json.dump(data, open(path, "w"), indent=2)
print(f"marked tested: {ver}")
print(json.dumps(data, indent=2))
PY
