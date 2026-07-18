#!/usr/bin/env bash
# scripts/ensure-secrets.sh -- Create secrets/remarkable.local.env and fill
# missing fields. Reuses saved password, Wi-Fi IP, sync repo, and token.
# Prompts on a TTY; otherwise exits with a short hint when required fields lack.
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

_open_url() {
  local url="$1"
  if command -v open >/dev/null 2>&1; then
    open "$url" >/dev/null 2>&1 || true
  elif command -v xdg-open >/dev/null 2>&1; then
    xdg-open "$url" >/dev/null 2>&1 || true
  else
    echo "  Open this URL in a browser:"
    echo "  $url"
  fi
}

_ensure_keys() {
  # Older secrets files may lack newer keys; append blanks so _get/_set work.
  local k
  for k in PIN_DIGITS SYNC_REPO GH_TOKEN SYNC_SKIP; do
    if ! grep -qE "^[[:space:]]*${k}=" "$SECRETS" 2>/dev/null; then
      printf '%s=\n' "$k" >>"$SECRETS"
    fi
  done
}

_valid_repo() {
  printf '%s' "$1" | grep -Eq '^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$'
}

if [ ! -f "$EXAMPLE" ]; then
  err "missing template: $EXAMPLE"
  exit 1
fi

if [ ! -f "$SECRETS" ]; then
  cp "$EXAMPLE" "$SECRETS"
  echo "Created secrets/remarkable.local.env from template."
fi
_ensure_keys

PASS="$(_get RM_ROOT_PASSWORD)"
WIFI="$(_get RM_HOST_WIFI)"
PIN_DIGITS="$(_get PIN_DIGITS)"
SYNC_REPO="$(_get SYNC_REPO)"
GH_TOKEN="$(_get GH_TOKEN)"
SYNC_SKIP="$(_get SYNC_SKIP)"

echo
echo "Saved install details (secrets/remarkable.local.env -- not committed):"
if [ -n "$PASS" ]; then
  echo "  Tablet password: (saved)"
else
  echo "  Tablet password: (needed)"
fi
if [ -n "$WIFI" ]; then
  echo "  Wi-Fi IP:        $WIFI"
else
  echo "  Wi-Fi IP:        (needed)"
fi
if [ -n "$PIN_DIGITS" ]; then
  echo "  Phone PIN:       $PIN_DIGITS"
else
  echo "  Phone PIN:       (needed)"
fi
if [ -n "$SYNC_REPO" ]; then
  echo "  Notes repo:      $SYNC_REPO"
elif [ "$SYNC_SKIP" = "1" ]; then
  echo "  Notes sync:      skipped (saved choice)"
else
  echo "  Notes sync:      (optional)"
fi
if [ -n "$GH_TOKEN" ]; then
  echo "  GitHub token:    (saved)"
elif [ -n "$SYNC_REPO" ]; then
  echo "  GitHub token:    (needed for sync)"
fi
echo

need_pass=0
need_wifi=0
need_pin=0
[ -z "$PASS" ] && need_pass=1
[ -z "$WIFI" ] && need_wifi=1
[ -z "$PIN_DIGITS" ] && need_pin=1

need_sync_prompt=0
if [ -z "$SYNC_REPO" ] && [ "$SYNC_SKIP" != "1" ]; then
  need_sync_prompt=1
fi
need_token=0
if [ -n "$SYNC_REPO" ] && [ -z "$GH_TOKEN" ]; then
  need_token=1
fi

if [ "$need_pass" -eq 0 ] && [ "$need_wifi" -eq 0 ] && [ "$need_pin" -eq 0 ] \
  && [ "$need_sync_prompt" -eq 0 ] && [ "$need_token" -eq 0 ]; then
  echo "Nothing new to ask -- reusing saved values."
  echo
  exit 0
fi

if [ ! -t 0 ]; then
  if [ "$need_pass" -eq 1 ] || [ "$need_wifi" -eq 1 ]; then
    err "tablet password/Wi-Fi missing and no interactive terminal"
    echo "  Edit secrets/remarkable.local.env (RM_ROOT_PASSWORD, RM_HOST_WIFI)." >&2
    exit 1
  fi
  if [ "$need_token" -eq 1 ]; then
    err "SYNC_REPO is set but GH_TOKEN is empty (no TTY to paste a token)"
    echo "  Set GH_TOKEN in secrets, or run ensure-secrets.sh in a terminal." >&2
    exit 1
  fi
  # Optional PIN / sync questions -- skip quietly when non-interactive.
  echo "Skipping optional PIN/sync prompts (no TTY)."
  echo
  exit 0
fi

if [ "$need_pass" -eq 1 ]; then
  echo "Root password -- on the tablet: Settings > Help > Copyrights and licenses >"
  echo "General information (scroll down). Changes after every firmware update."
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

if [ "$need_wifi" -eq 1 ]; then
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

if [ "$need_pin" -eq 1 ]; then
  echo "Phone connection PIN (shown on the tablet Lobby; entered on the phone page)."
  echo "  6 = six digits (default)   4 = four digits   none = no PIN (anyone on your Wi-Fi)"
  while true; do
    printf 'PIN length [6/4/none]: '
    IFS= read -r PIN_DIGITS || true
    PIN_DIGITS="$(printf '%s' "$PIN_DIGITS" | tr '[:upper:]' '[:lower:]' | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    [ -z "$PIN_DIGITS" ] && PIN_DIGITS=6
    case "$PIN_DIGITS" in
      6|4|none) break ;;
      *) echo "  Choose 6, 4, or none" ;;
    esac
  done
  _set PIN_DIGITS "$PIN_DIGITS"
fi

if [ "$need_sync_prompt" -eq 1 ]; then
  echo "Optional: sync notes to a private GitHub repo?"
  printf 'Enable sync? [y/N]: '
  IFS= read -r ans || true
  ans="$(printf '%s' "$ans" | tr '[:upper:]' '[:lower:]')"
  case "$ans" in
    y|yes)
      while true; do
        printf 'GitHub repo (owner/repo): '
        IFS= read -r SYNC_REPO || true
        SYNC_REPO="$(printf '%s' "$SYNC_REPO" | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
        if _valid_repo "$SYNC_REPO"; then
          break
        fi
        echo "  Use the form owner/repo (example: alice/my-notes)"
      done
      _set SYNC_REPO "$SYNC_REPO"
      _set SYNC_SKIP ""
      need_token=1
      ;;
    *)
      _set SYNC_SKIP "1"
      _set SYNC_REPO ""
      echo "  Sync skipped (ask again by clearing SYNC_SKIP in secrets)."
      ;;
  esac
fi

if [ "$need_token" -eq 1 ]; then
  SYNC_REPO="$(_get SYNC_REPO)"
  owner="${SYNC_REPO%%/*}"
  # Prefill fine-grained PAT form (GitHub template URL). User still picks the repo
  # in the UI, generates once, then pastes here.
  pat_url="https://github.com/settings/personal-access-tokens/new"
  pat_url="${pat_url}?name=Writerdeck%20notes%20sync"
  pat_url="${pat_url}&description=Contents%20read%2Fwrite%20for%20Writerdeck%20notes%20sync"
  pat_url="${pat_url}&contents=write"
  if [ -n "$owner" ]; then
    pat_url="${pat_url}&target_name=${owner}"
  fi
  echo
  echo "Opening GitHub to create a fine-grained token (Contents: Read and write)"
  echo "for ${SYNC_REPO}. Generate the token, copy it, then paste below."
  _open_url "$pat_url"
  echo "  (If the browser did not open: $pat_url)"
  while true; do
    printf 'Paste GitHub token: '
    IFS= read -r -s GH_TOKEN || true
    echo
    GH_TOKEN="$(printf '%s' "$GH_TOKEN" | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    if [ -n "$GH_TOKEN" ]; then
      break
    fi
    echo "  (needed for sync -- try again, or Ctrl-C to abort)"
  done
  _set GH_TOKEN "$GH_TOKEN"
fi

echo "Saved."
echo
