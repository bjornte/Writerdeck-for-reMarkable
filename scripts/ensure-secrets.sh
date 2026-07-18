#!/usr/bin/env bash
# scripts/ensure-secrets.sh -- Create secrets/remarkable.local.env and fill
# missing fields. Reuses saved password, Wi-Fi IP, sync repo, and token.
# Prompts on a TTY; otherwise exits with a short hint when required fields lack.
# Phone PIN is not prompted: empty PIN_DIGITS becomes "none" (device-ack for
# sensitive phone actions is planned later; see docs/improvements.md).
#
# Usage (from repo root, or via install.sh / preflight.sh):
#   bash scripts/ensure-secrets.sh
#
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
EXAMPLE="$REPO/secrets/remarkable.local.env.example"
SECRETS="$REPO/secrets/remarkable.local.env"

# Set by install.sh so reuse exits stay quiet under the friendlier install UI.
INSTALL_UI="${WRITERDECK_INSTALL:-0}"

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

# Password SSH check (before key install). Uses SSH_ASKPASS so the password is
# not typed twice and does not appear in process lists as a flag.
_password_ssh_ok() {
  local host="$1" pass="$2" askpass pwfile out
  pwfile="$(mktemp)"
  askpass="$(mktemp)"
  printf '%s\n' "$pass" >"$pwfile"
  chmod 600 "$pwfile"
  printf '#!/bin/sh\ncat "%s"\n' "$pwfile" >"$askpass"
  chmod 700 "$askpass"
  out="$(
    SSH_ASKPASS="$askpass" SSH_ASKPASS_REQUIRE=force DISPLAY="${DISPLAY:-none}" \
      ssh -o StrictHostKeyChecking=accept-new -o ConnectTimeout=10 \
          -o PreferredAuthentications=password -o PubkeyAuthentication=no \
          -o NumberOfPasswordPrompts=1 \
          "root@$host" \
          'if grep -qi remarkable /etc/os-release 2>/dev/null; then echo __rm_ok__; else echo __not_rm__; fi' \
          </dev/null 2>/dev/null || true
  )"
  rm -f "$askpass" "$pwfile"
  printf '%s' "$out" | grep -q '__rm_ok__'
}

# Key login already works -- confirm the host looks like a reMarkable.
_key_ssh_is_remarkable() {
  local host="$1" out
  out="$(
    ssh -o StrictHostKeyChecking=accept-new -o ConnectTimeout=8 \
        -o BatchMode=yes \
        "root@$host" \
        'if grep -qi remarkable /etc/os-release 2>/dev/null; then echo __rm_ok__; fi' \
        2>/dev/null || true
  )"
  printf '%s' "$out" | grep -q '__rm_ok__'
}

if [ ! -f "$EXAMPLE" ]; then
  err "missing template: $EXAMPLE"
  exit 1
fi

if [ ! -f "$SECRETS" ]; then
  cp "$EXAMPLE" "$SECRETS"
  if [ "$INSTALL_UI" != "1" ]; then
    echo "Created secrets/remarkable.local.env from template."
  fi
fi
_ensure_keys

# Phone PIN: never prompt. Default to none (see docs/improvements.md).
PIN_DIGITS="$(_get PIN_DIGITS)"
if [ -z "$PIN_DIGITS" ]; then
  _set PIN_DIGITS "none"
  PIN_DIGITS=none
fi

PASS="$(_get RM_ROOT_PASSWORD)"
WIFI="$(_get RM_HOST_WIFI)"
SYNC_REPO="$(_get SYNC_REPO)"
GH_TOKEN="$(_get GH_TOKEN)"
SYNC_SKIP="$(_get SYNC_SKIP)"

need_pass=0
need_wifi=0
[ -z "$PASS" ] && need_pass=1
[ -z "$WIFI" ] && need_wifi=1

need_sync_prompt=0
if [ -z "$SYNC_REPO" ] && [ "$SYNC_SKIP" != "1" ]; then
  need_sync_prompt=1
fi
need_token=0
if [ -n "$SYNC_REPO" ] && [ -z "$GH_TOKEN" ]; then
  need_token=1
fi

if [ "$need_pass" -eq 0 ] && [ "$need_wifi" -eq 0 ] \
  && [ "$need_sync_prompt" -eq 0 ] && [ "$need_token" -eq 0 ]; then
  if [ "$INSTALL_UI" != "1" ]; then
    echo "Nothing new to ask -- reusing saved values."
    echo
  fi
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
  # Optional sync questions -- skip quietly when non-interactive.
  if [ "$INSTALL_UI" != "1" ]; then
    echo "Skipping optional sync prompts (no TTY)."
    echo
  fi
  exit 0
fi

# --- Wi-Fi first (so we can ping before asking for the password) --------------
if [ "$need_wifi" -eq 1 ]; then
  echo
  echo "Tablet Wi-Fi address"
  echo "--------------------"
  echo "Find it in the tablet's Wi-Fi settings, or in your router's"
  echo "list of connected devices. Example: 192.168.1.8"
  echo
  while true; do
    printf 'Wi-Fi IP: '
    IFS= read -r WIFI || true
    WIFI="$(printf '%s' "$WIFI" | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    if [ -z "$WIFI" ]; then
      echo "  (needed -- try again)"
      continue
    fi
    if ping -c1 -W2 "$WIFI" >/dev/null 2>&1; then
      echo
      echo "Check: A device was found"
      break
    fi
    echo "  No device answered at $WIFI. Check the address and that the tablet is awake."
  done
  _set RM_HOST_WIFI "$WIFI"
fi

# --- Root password -----------------------------------------------------------
if [ "$need_pass" -eq 1 ]; then
  echo
  echo "Tablet root password"
  echo "--------------------"
  echo "On the tablet: Settings -> Help -> Copyrights and licenses ->"
  echo "General information (scroll to the bottom)."
  echo
  echo "This password changes after every firmware update."
  echo "You can paste it into this window. It does not need to be typed by hand."
  echo
  WIFI="$(_get RM_HOST_WIFI)"
  while true; do
    printf 'Password: '
    IFS= read -r -s PASS || true
    echo
    if [ -z "$PASS" ]; then
      echo "  (needed -- try again)"
      continue
    fi
    if [ -n "$WIFI" ] && _key_ssh_is_remarkable "$WIFI"; then
      echo
      echo "Check: Password is correct, device identifies as a reMarkable"
      echo "  (key login already works; password saved for reference)"
      break
    fi
    if [ -n "$WIFI" ] && _password_ssh_ok "$WIFI" "$PASS"; then
      echo
      echo "Check: Password is correct, device identifies as a reMarkable"
      break
    fi
    echo "  Could not sign in with that password (or the device is not a reMarkable)."
    echo "  Try again, or check Wi-Fi / that the tablet is awake."
  done
  _set RM_ROOT_PASSWORD "$PASS"
fi

# --- Optional notes sync -----------------------------------------------------
if [ "$need_sync_prompt" -eq 1 ]; then
  echo
  echo "Notes sync (optional)"
  echo "---------------------"
  echo "You can keep a backup of your notes in a private GitHub repo."
  echo "Skip this if you only want notes on the tablet for now."
  echo "The repo can be set up later."
  echo
  printf 'Enable sync? [y/N]: '
  IFS= read -r ans || true
  ans="$(printf '%s' "$ans" | tr '[:upper:]' '[:lower:]')"
  case "$ans" in
    y|yes)
      echo
      echo "GitHub repo (owner/name), for example bjornte/my-notes"
      while true; do
        printf 'Repo: '
        IFS= read -r SYNC_REPO || true
        SYNC_REPO="$(printf '%s' "$SYNC_REPO" | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
        if _valid_repo "$SYNC_REPO"; then
          break
        fi
        echo "  Use the form owner/name (example: alice/my-notes)"
      done
      _set SYNC_REPO "$SYNC_REPO"
      _set SYNC_SKIP ""
      need_token=1
      ;;
    *)
      _set SYNC_SKIP "1"
      _set SYNC_REPO ""
      ;;
  esac
fi

if [ "$need_token" -eq 1 ]; then
  SYNC_REPO="$(_get SYNC_REPO)"
  owner="${SYNC_REPO%%/*}"
  # Prefill fine-grained PAT form (GitHub template URL). User still confirms
  # the repo in the UI, generates once, then pastes here.
  pat_url="https://github.com/settings/personal-access-tokens/new"
  pat_url="${pat_url}?name=Writerdeck%20notes%20sync"
  pat_url="${pat_url}&description=Contents%20read%2Fwrite%20for%20Writerdeck%20notes%20sync"
  pat_url="${pat_url}&contents=write"
  if [ -n "$owner" ]; then
    pat_url="${pat_url}&target_name=${owner}"
  fi
  echo
  echo "GitHub access token"
  echo "-------------------"
  echo "Writerdeck needs a fine-grained personal access token so it"
  echo "can read and write notes in ${SYNC_REPO}."
  echo
  echo "Your browser will open GitHub's token page with the right"
  echo "settings pre-filled. On that page:"
  echo
  echo "  1. Sign in if needed"
  echo "  2. Confirm the repo is ${SYNC_REPO}"
  echo "  3. Leave Contents on \"Read and write\""
  echo "  4. Click Generate token"
  echo "  5. Copy the token (you only see it once)"
  echo
  printf 'Press Enter to open the page...'
  IFS= read -r _ || true
  echo
  _open_url "$pat_url"
  echo
  echo "(If nothing opened, use this link:)"
  echo "$pat_url"
  echo
  while true; do
    printf 'Paste the token here (it will not be echoed): '
    IFS= read -r -s GH_TOKEN || true
    echo
    GH_TOKEN="$(printf '%s' "$GH_TOKEN" | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    if [ -n "$GH_TOKEN" ]; then
      break
    fi
    echo "  (needed for sync -- try again, or Ctrl-C to abort)"
  done
  _set GH_TOKEN "$GH_TOKEN"
  echo
  echo "Token saved."
fi

if [ "$INSTALL_UI" != "1" ]; then
  echo
  echo "Saved."
  echo
fi
