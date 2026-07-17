#!/usr/bin/env bash
# scripts/ensure-secrets.sh -- Create secrets/remarkable.local.env and fill
# RM_ROOT_PASSWORD + RM_HOST_WIFI if empty. Prompts on a TTY; otherwise exits
# with a short hint.
#
# Usage (from repo root, or via install.sh / preflight.sh):
#   bash scripts/ensure-secrets.sh
#
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
EXAMPLE="$REPO/secrets/remarkable.local.env.example"
SECRETS="$REPO/secrets/remarkable.local.env"

if [ -t 2 ]; then R=$'\033[1;31m'; Z=$'\033[0m'; else R=''; Z=''; fi
err() { printf '%sERROR:%s %s\n' "$R" "$Z" "$*" >&2; }

_get() {
  [ -f "$SECRETS" ] || return 0
  sed -n -E "s/^[[:space:]]*$1[[:space:]]*=[[:space:]]*([^#]*).*/\1/p" "$SECRETS" \
    | head -n1 | sed -E 's/[[:space:]]+$//'
}

# Replace KEY=... in place (or append). Value must not contain newlines.
_set() {
  local key="$1" val="$2"
  if grep -qE "^[[:space:]]*${key}=" "$SECRETS"; then
    # Avoid sed -i portability issues: rewrite via temp file.
    local tmp
    tmp="$(mktemp)"
    awk -v k="$key" -v v="$val" '
      BEGIN { done=0 }
      $0 ~ "^[[:space:]]*"k"=" {
        print k"="v
        done=1
        next
      }
      { print }
      END { if (!done) print k"="v }
    ' "$SECRETS" >"$tmp"
    mv "$tmp" "$SECRETS"
  else
    printf '%s=%s\n' "$key" "$val" >>"$SECRETS"
  fi
}

if [ ! -f "$EXAMPLE" ]; then
  err "missing template: $EXAMPLE"
  exit 1
fi

if [ ! -f "$SECRETS" ]; then
  cp "$EXAMPLE" "$SECRETS"
  echo "Created secrets/remarkable.local.env from template."
fi

PASS="$(_get RM_ROOT_PASSWORD)"
WIFI="$(_get RM_HOST_WIFI)"

need_prompt=0
[ -z "$PASS" ] && need_prompt=1
[ -z "$WIFI" ] && need_prompt=1

if [ "$need_prompt" -eq 0 ]; then
  exit 0
fi

if [ ! -t 0 ]; then
  err "secrets incomplete and no interactive terminal"
  echo "  Edit secrets/remarkable.local.env and set RM_ROOT_PASSWORD + RM_HOST_WIFI." >&2
  exit 1
fi

echo
echo "Tablet login details (saved only in secrets/remarkable.local.env, not committed):"
echo

if [ -z "$PASS" ]; then
  echo "Root password -- on the tablet: Settings > Help > Copyrights and licenses >"
  echo "General information (scroll down). Changes after every firmware update."
  # -s hides typing; -r keeps backslashes literal
  while true; do
    printf 'Password: '
    IFS= read -r -s PASS || true
    echo
    if [ -n "$PASS" ]; then
      break
    fi
    echo "  (needed -- try again)"
  done
  _set RM_ROOT_PASSWORD "$PASS"
fi

if [ -z "$WIFI" ]; then
  echo "Wi-Fi address of the tablet (tablet Wi-Fi settings, or your router device list)."
  while true; do
    printf 'Wi-Fi IP: '
    IFS= read -r WIFI || true
    WIFI="$(printf '%s' "$WIFI" | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    if [ -n "$WIFI" ]; then
      break
    fi
    echo "  (needed -- try again)"
  done
  _set RM_HOST_WIFI "$WIFI"
fi

echo "Saved."
echo
