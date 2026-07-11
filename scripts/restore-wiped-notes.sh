#!/usr/bin/env bash
# scripts/restore-wiped-notes.sh -- recover notes wiped by the Lobby Home bug (Jul 2026).
#
# 1. Restores zero-byte .md files on the tablet from the last good GitHub commit.
# 2. Restores the same files on github.com/bjornte/my-notes (Contents API).
# 3. Removes junk "(tablet copy).md" side files created by sync clash handling.
#
# Usage (from repo root):
#   bash scripts/restore-wiped-notes.sh
#   bash scripts/restore-wiped-notes.sh 192.168.1.8
#
# Requires: gh (authenticated), ssh to tablet, secrets/remarkable.local.env

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"

TARGET="${1:-$RM_HOST}"
REPO="bjornte/my-notes"
NOTES_DIR="/home/root/Writerdeck-user-documents"
# Parent of the Jul 11 13:27 empty push that wiped Additional improvements…
RESTORE_REF="${RESTORE_REF:-909ff1aea2a53907bc905e1cb9e3ab6b90c73fc9}"

if ! command -v gh >/dev/null 2>&1; then
  err "gh CLI required."
  exit 1
fi

if ! ping -c1 -W2 "$TARGET" >/dev/null 2>&1; then
  err "Tablet unreachable at $TARGET."
  exit 1
fi

echo "=== Restore wiped notes (ref=${RESTORE_REF:0:7}) ==="

# --- tablet: list zero-byte notes ---
EMPTY_ON_TABLET=()
while IFS= read -r line; do
  [ -n "$line" ] && EMPTY_ON_TABLET+=("$line")
done < <(
  ssh $RM_SSH_OPTS -o ConnectTimeout=8 "root@$TARGET" \
    "for f in ${NOTES_DIR}/*.md; do [ -f \"\$f\" ] || continue; s=\$(wc -c < \"\$f\" | tr -d ' '); [ \"\$s\" = 0 ] && basename \"\$f\"; done" \
    2>/dev/null || true
)

# --- always try to restore known victims even if size check missed ---
RESTORE_NAMES=(
  "Additional improvements to the Writerdeck.md"
)
for f in "${EMPTY_ON_TABLET[@]:-}"; do
  case "$f" in
    *'(tablet copy)'*) continue ;;
  esac
  case " ${RESTORE_NAMES[*]} " in *" $f "*) ;; *) RESTORE_NAMES+=("$f") ;; esac
done

echo "Files to restore on tablet + GitHub:"
printf '  %s\n' "${RESTORE_NAMES[@]}"

restore_one() {
  local name="$1"
  local enc
  enc=$(python3 -c 'import urllib.parse,sys; print(urllib.parse.quote(sys.argv[1]))' "$name")
  local blob
  blob=$(gh api "repos/${REPO}/contents/${enc}?ref=${RESTORE_REF}" --jq .content 2>/dev/null) || {
    echo "  SKIP $name (not in ref ${RESTORE_REF:0:7})"
    return 0
  }
  local content
  content=$(printf '%s' "$blob" | base64 -d)
  local bytes=${#content}
  if [ "$bytes" -eq 0 ]; then
    echo "  SKIP $name (ref also empty)"
    return 0
  fi

  echo "  RESTORE $name ($bytes bytes)"
  # Tablet
  printf '%s' "$content" | ssh $RM_SSH_OPTS -o ConnectTimeout=8 "root@$TARGET" \
    "cat > '${NOTES_DIR}/${name}'"

  # GitHub (current branch)
  local cur_sha cur_size
  cur_sha=$(gh api "repos/${REPO}/contents/${enc}" --jq .sha 2>/dev/null || true)
  cur_size=$(gh api "repos/${REPO}/contents/${enc}" --jq .size 2>/dev/null || echo "-1")
  if [ "${cur_size:-0}" -gt 0 ] 2>/dev/null; then
    echo "    GitHub already has content (${cur_size} bytes) — tablet only"
    return 0
  fi
  local b64
  b64=$(printf '%s' "$content" | base64 | tr -d '\n')
  if [ -z "${cur_sha:-}" ]; then
    gh api "repos/${REPO}/contents/${enc}" -X PUT \
      -f message="Restore ${name} after Lobby wipe bug" \
      -f content="$b64" \
      -f encoding="base64" >/dev/null
  else
    gh api "repos/${REPO}/contents/${enc}" -X PUT \
      -f message="Restore ${name} after Lobby wipe bug" \
      -f content="$b64" \
      -f encoding="base64" \
      -f sha="$cur_sha" >/dev/null
  fi
  echo "    GitHub restored"
}

for name in "${RESTORE_NAMES[@]}"; do
  restore_one "$name"
done

# --- remove junk (tablet copy) files ---
echo ""
echo "=== Remove (tablet copy) clash duplicates ==="
COPIES=()
while IFS= read -r line; do
  [ -n "$line" ] && COPIES+=("$line")
done < <(
  ssh $RM_SSH_OPTS -o ConnectTimeout=8 "root@$TARGET" \
    "for f in ${NOTES_DIR}/*'(tablet copy)'.md; do [ -f \"\$f\" ] && basename \"\$f\"; done" \
    2>/dev/null || true
)

delete_gh_file() {
  local name="$1"
  local enc sha
  enc=$(python3 -c 'import urllib.parse,sys; print(urllib.parse.quote(sys.argv[1]))' "$name")
  sha=$(gh api "repos/${REPO}/contents/${enc}" --jq .sha 2>/dev/null || echo "")
  [ -n "$sha" ] || return 0
  gh api "repos/${REPO}/contents/${enc}" -X DELETE \
    -f message="Remove clash duplicate ${name}" \
    -f sha="$sha" >/dev/null 2>&1 || true
}

for name in "${COPIES[@]:-}"; do
  echo "  DELETE $name"
  ssh $RM_SSH_OPTS -o ConnectTimeout=8 "root@$TARGET" "rm -f '${NOTES_DIR}/${name}'"
  delete_gh_file "$name"
done

echo ""
echo "=== Done. Redeploy Writerdeck (rmkw) so the doLoad fix is on-device. ==="
ssh $RM_SSH_OPTS -o ConnectTimeout=8 "root@$TARGET" \
  "for f in ${NOTES_DIR}/*.md; do [ -f \"\$f\" ] && echo \"\$(basename \"\$f\") \$(wc -c < \"\$f\")\"; done"
