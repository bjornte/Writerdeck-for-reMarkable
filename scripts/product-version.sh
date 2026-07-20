#!/usr/bin/env bash
# scripts/product-version.sh -- Writerdeck product stamp (YYYY-MM-DD or YYYY-MM-DD.N).
#
# Default --write: if VERSION is already today's date (or today.N), keep it so
# server and editor builds the same day share one stamp. If VERSION is older or
# missing, set it to today.
#
# --bump: force the next stamp (today, or today.2 / .3 / ...). Use when you
# intentionally ship a second product build the same day.
#
# Usage:
#   bash scripts/product-version.sh           # print next/current without write
#   bash scripts/product-version.sh --write   # ensure VERSION is today; print it
#   bash scripts/product-version.sh --bump    # advance stamp and write VERSION
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "${DIR}/.." && pwd)"
VERSION_FILE="${REPO}/VERSION"
MODE="print"
if [ "${1:-}" = "--write" ]; then
  MODE="write"
elif [ "${1:-}" = "--bump" ]; then
  MODE="bump"
fi

TODAY="$(date +%Y-%m-%d)"
CUR=""
if [ -f "${VERSION_FILE}" ]; then
  CUR="$(tr -d '[:space:]' < "${VERSION_FILE}")"
fi

is_today() {
  local v="$1"
  [ "${v}" = "${TODAY}" ] && return 0
  [[ "${v}" =~ ^${TODAY}\.[0-9]+$ ]] && return 0
  return 1
}

next_after() {
  local v="$1"
  if [ -z "${v}" ]; then
    printf '%s\n' "${TODAY}"
    return
  fi
  if [ "${v}" = "${TODAY}" ]; then
    printf '%s\n' "${TODAY}.2"
    return
  fi
  if [[ "${v}" =~ ^${TODAY}\.([0-9]+)$ ]]; then
    printf '%s\n' "${TODAY}.$((BASH_REMATCH[1] + 1))"
    return
  fi
  printf '%s\n' "${TODAY}"
}

OUT=""
case "${MODE}" in
  print)
    if is_today "${CUR}"; then
      OUT="${CUR}"
    else
      OUT="${TODAY}"
    fi
    ;;
  write)
    if is_today "${CUR}"; then
      OUT="${CUR}"
    else
      OUT="${TODAY}"
    fi
    printf '%s\n' "${OUT}" > "${VERSION_FILE}"
    ;;
  bump)
    OUT="$(next_after "${CUR}")"
    printf '%s\n' "${OUT}" > "${VERSION_FILE}"
    ;;
esac

printf '%s\n' "${OUT}"
