#!/usr/bin/env bash
# scripts/discover-rm2-wifi.sh -- Find rM1 and rM2 Wi-Fi IPs.
#
# Prefer USB: ask the plugged-in tablet over SSH (fast, exact model + IP).
# Always also scan the LAN for Writerdeck on :8000 so a second tablet is not
# missed. Saved RM_HOST_WIFI / RM2_HOST_WIFI are used only to label hits.
#
# rM2 often refuses SSH on Wi-Fi; deploy still uses USB (rm2-pick-host).
# Phone UI URLs use the Wi-Fi IPs.
#
# Usage (repo root, bash):
#   bash scripts/discover-rm2-wifi.sh
#   bash scripts/discover-rm2-wifi.sh --write-secrets
#   bash scripts/discover-rm2-wifi.sh --scan-only

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"

WRITE_SECRETS=0
SCAN_ONLY=0
for arg in "$@"; do
  case "$arg" in
    --write-secrets) WRITE_SECRETS=1 ;;
    --scan-only)     SCAN_ONLY=1 ;;
    -h|--help)
      sed -n '2,16p' "$0"
      exit 0
      ;;
    *)
      err "unknown arg: $arg"
      exit 2
      ;;
  esac
done

USB="${RM_HOST_USB:-${RM2_HOST_USB:-10.11.99.1}}"

# Returns 0 when http://HOST:8000/ looks like Writerdeck (any HTTP response).
probe_writerdeck_http() {
  local host="$1"
  curl -fsS -m 2 "http://${host}:8000/" >/dev/null 2>&1
}

# Print rm1, rm2, or unknown from hardware over SSH.
classify_ssh() {
  local host="$1"
  local out
  out="$(rm_ssh 'm=$(cat /sys/devices/soc0/machine 2>/dev/null)
[ -z "$m" ] && m=$(tr -d "\0" < /proc/device-tree/model 2>/dev/null)
case "$m" in
  *2.0*|*"reMarkable 2"*) echo rm2 ;;
  *1.0*|*"reMarkable 1"*) echo rm1 ;;
  *) echo unknown ;;
esac' "$host" 2>/dev/null | tr -d '\r\n')"
  case "$out" in
    rm1|rm2) echo "$out" ;;
    *) echo unknown ;;
  esac
}

# Wi-Fi IPv4 of the USB-plugged tablet (status API, else wlan0).
usb_wifi_ip() {
  local ip
  rm_test_key "$USB" || return 1
  ip="$(rm_ssh 'wget -qO- http://127.0.0.1:8000/api/status 2>/dev/null \
    | sed -n "s/.*\"ip\":\"\\([^\"]*\\)\".*/\\1/p" | head -n1' "$USB" \
    | tr -d '\r\n')"
  if [ -n "$ip" ] && [ "$ip" != "null" ]; then
    echo "$ip"
    return 0
  fi
  rm_ssh 'ip -4 -o addr show wlan0 2>/dev/null \
    | awk "{print \$4}" | head -n1 | cut -d/ -f1' "$USB" \
    | tr -d '\r\n'
}

local_subnet_prefix() {
  local ip iface
  for iface in en0 en1 bridge0; do
    ip="$(ipconfig getifaddr "$iface" 2>/dev/null || true)"
    if [ -n "$ip" ]; then
      echo "$ip" | cut -d. -f1-3
      return 0
    fi
  done
  return 1
}

# Print every LAN IPv4 that answers Writerdeck on :8000 (one per line).
lan_scan_writerdeck() {
  local base host ip outdir
  base="$(local_subnet_prefix)" || return 1
  echo "Scanning ${base}.0/24 for Writerdeck on :8000 ..." >&2
  outdir="$(mktemp -d)"
  for host in $(seq 1 254); do
    (
      ip="${base}.${host}"
      if probe_writerdeck_http "$ip"; then
        printf '%s\n' "$ip" > "${outdir}/${ip}"
      fi
      exit 0
    ) &
  done
  wait || true
  if ls "$outdir"/* >/dev/null 2>&1; then
    cat "$outdir"/* | sort -u
  fi
  rm -rf "$outdir"
}

write_wifi_secret() {
  local key="$1" ip="$2"
  local secrets="$REPO/secrets/remarkable.local.env"
  if [ ! -f "$secrets" ]; then
    err "missing $secrets"
    return 1
  fi
  if grep -q "^${key}=" "$secrets"; then
    sed -i.bak -E "s/^${key}=.*/${key}=${ip}/" "$secrets"
    rm -f "${secrets}.bak"
  else
    printf '\n%s=%s\n' "$key" "$ip" >> "$secrets"
  fi
  echo "Updated ${key}=${ip} in secrets/remarkable.local.env"
}

FOUND_RM1=""
SRC_RM1=""
FOUND_RM2=""
SRC_RM2=""
UNKNOWN_IPS=()

already_found() {
  local ip="$1" u
  [ -n "$ip" ] || return 1
  [ "$ip" = "$FOUND_RM1" ] && return 0
  [ "$ip" = "$FOUND_RM2" ] && return 0
  for u in "${UNKNOWN_IPS[@]+"${UNKNOWN_IPS[@]}"}"; do
    [ "$u" = "$ip" ] && return 0
  done
  return 1
}

record() {
  local which="$1" ip="$2" src="$3"
  [ -n "$ip" ] || return 0
  already_found "$ip" && return 0
  case "$which" in
    rm1)
      if [ -z "$FOUND_RM1" ]; then
        FOUND_RM1="$ip"
        SRC_RM1="$src"
      fi
      ;;
    rm2)
      if [ -z "$FOUND_RM2" ]; then
        FOUND_RM2="$ip"
        SRC_RM2="$src"
      fi
      ;;
    *)
      UNKNOWN_IPS+=("$ip")
      ;;
  esac
}

# USB: whichever tablet is plugged in (both use 10.11.99.1).
if [ "$SCAN_ONLY" != "1" ]; then
  if rm_test_key "$USB"; then
    model="$(classify_ssh "$USB")"
    ip="$(usb_wifi_ip 2>/dev/null || true)"
    src="usb"
    if [ -n "$ip" ]; then
      if probe_writerdeck_http "$ip"; then
        src="usb-status"
      else
        src="usb-wlan0"
      fi
    fi
    if [ "$model" = "unknown" ]; then
      if [ "$ip" = "${RM_HOST_WIFI:-}" ]; then
        model="rm1"
      elif [ "$ip" = "${RM2_HOST_WIFI:-}" ]; then
        model="rm2"
      fi
    fi
    if [ -n "$ip" ]; then
      record "$model" "$ip" "$src"
    fi
  fi
fi

# LAN scan, then label each hit.
HITS=""
HITS="$(lan_scan_writerdeck || true)"

# Also consider saved IPs that still answer (in case they sit off this /24).
for hint in "${RM_HOST_WIFI:-}" "${RM2_HOST_WIFI:-}"; do
  [ -n "$hint" ] || continue
  printf '%s\n' "$HITS" | grep -qx "$hint" && continue
  if probe_writerdeck_http "$hint"; then
    HITS="$(printf '%s\n%s\n' "$HITS" "$hint")"
  fi
done

while IFS= read -r ip; do
  [ -n "$ip" ] || continue
  already_found "$ip" && continue

  which="unknown"
  if rm_test_key "$ip"; then
    which="$(classify_ssh "$ip")"
  fi
  if [ "$which" = "unknown" ]; then
    if [ "$ip" = "${RM_HOST_WIFI:-}" ]; then
      which="rm1"
    elif [ "$ip" = "${RM2_HOST_WIFI:-}" ]; then
      which="rm2"
    fi
  fi
  record "$which" "$ip" "lan-scan"
done <<< "$HITS"

# One unlabeled Writerdeck + one known model => the other model.
if [ "${#UNKNOWN_IPS[@]}" -eq 1 ]; then
  leftover="${UNKNOWN_IPS[0]}"
  if [ -z "$FOUND_RM1" ] && [ -n "$FOUND_RM2" ]; then
    UNKNOWN_IPS=()
    record rm1 "$leftover" "lan-scan"
  elif [ -z "$FOUND_RM2" ] && [ -n "$FOUND_RM1" ]; then
    UNKNOWN_IPS=()
    record rm2 "$leftover" "lan-scan"
  fi
fi

if [ -z "$FOUND_RM1" ] && [ -z "$FOUND_RM2" ] && [ "${#UNKNOWN_IPS[@]}" -eq 0 ]; then
  err "could not find Writerdeck on Wi-Fi (plug in USB or join the same LAN)"
  exit 1
fi

print_one() {
  local label="$1" ip="$2" src="$3"
  printf '%s=%s\n' "$label" "$ip"
  echo "${label}: ${ip}  (source: ${src})" >&2
  echo "  phone UI: http://${ip}:8000/" >&2
}

if [ -n "$FOUND_RM1" ]; then
  print_one RM_HOST_WIFI "$FOUND_RM1" "$SRC_RM1"
fi
if [ -n "$FOUND_RM2" ]; then
  print_one RM2_HOST_WIFI "$FOUND_RM2" "$SRC_RM2"
  echo "  deploy SSH: use USB ${USB} (Wi-Fi SSH often refused on rM2)" >&2
fi
for ip in "${UNKNOWN_IPS[@]+"${UNKNOWN_IPS[@]}"}"; do
  [ -n "$ip" ] || continue
  echo "unknown Writerdeck: ${ip}" >&2
  echo "  phone UI: http://${ip}:8000/" >&2
done

if [ -z "$FOUND_RM1" ]; then
  echo "rM1 not found (asleep, or Writerdeck not running)" >&2
fi
if [ -z "$FOUND_RM2" ]; then
  echo "rM2 not found (asleep, or Writerdeck not running)" >&2
fi
if [ "${#UNKNOWN_IPS[@]}" -gt 0 ]; then
  echo "Could not tell rM1 from rM2. Plug one in over USB and re-run." >&2
fi

if [ "$WRITE_SECRETS" = "1" ]; then
  [ -n "$FOUND_RM1" ] && write_wifi_secret RM_HOST_WIFI "$FOUND_RM1"
  [ -n "$FOUND_RM2" ] && write_wifi_secret RM2_HOST_WIFI "$FOUND_RM2"
fi
