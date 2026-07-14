#!/usr/bin/env bash
# Recover .md.enc notes orphaned by a vault key rotation (disable+setup or sync).
#
# Usage (repo root):
#   bash scripts/recover-orphaned-vault-notes.sh \
#     --old-vault-ref 4fbd7a93 \
#     --notes "til 1.md.enc,til 2.md.enc" \
#     --pin 123456
#
# --old-vault-ref is a git commit or ref in syncRepo where secret/vault still
# matches the key that encrypted the notes. Find via:
#   gh api repos/OWNER/REPO/commits?path=secret/vault

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"

OLD_REF=""
NOTES=""
PIN=""
REPO_OVERRIDE=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --old-vault-ref) OLD_REF="$2"; shift 2 ;;
    --notes) NOTES="$2"; shift 2 ;;
    --pin) PIN="$2"; shift 2 ;;
    --repo) REPO_OVERRIDE="$2"; shift 2 ;;
    *) echo "Unknown arg: $1" >&2; exit 1 ;;
  esac
done

[[ -n "$OLD_REF" && -n "$NOTES" && -n "$PIN" ]] || {
  echo "Usage: $0 --old-vault-ref COMMIT --notes 'a.md.enc,b.md.enc' --pin PIN" >&2
  exit 1
}

BASE="http://$RM_HOST:8000"
SYNC_REPO="$REPO_OVERRIDE"
if [[ -z "$SYNC_REPO" ]]; then
  SYNC_REPO=$(curl -sf "$BASE/api/settings" | python3 -c "import sys,json; print(json.load(sys.stdin).get('syncRepo',''))")
fi
[[ -n "$SYNC_REPO" ]] || { echo "syncRepo not configured; pass --repo owner/repo" >&2; exit 1; }

OLD_VAULT=$(mktemp)
trap 'rm -f "$OLD_VAULT"' EXIT
gh api "repos/${SYNC_REPO}/contents/secret/vault?ref=${OLD_REF}" --jq '.content' | base64 -d >"$OLD_VAULT"

IFS=',' read -ra NOTE_ARR <<< "$NOTES"
NOTES_JSON=$(python3 -c "import json,sys; print(json.dumps(sys.argv[1:]))" "${NOTE_ARR[@]}")

BODY=$(python3 -c '
import json,sys
old_vault=open(sys.argv[1]).read()
print(json.dumps({"op":"vaultrewrap","old":old_vault,"name":sys.argv[2],"notes":json.loads(sys.argv[3])}))
' "$OLD_VAULT" "$PIN" "$NOTES_JSON")

code=$(curl -s -o /tmp/vault-rewrap.out -w '%{http_code}' -X POST "$BASE/api/test/tablet-req" \
  -H 'Content-Type: application/json' -d "$BODY")
if [[ "$code" != "200" ]]; then
  echo "vaultrewrap failed HTTP $code:" >&2
  cat /tmp/vault-rewrap.out >&2
  exit 1
fi

echo "OK: re-wrapped ${#NOTE_ARR[@]} note(s). Open on tablet to verify."
echo "If a note still fails, it may need a different --old-vault-ref (check secret/vault commits near when that note was encrypted)."
