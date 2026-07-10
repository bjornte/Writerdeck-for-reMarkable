#!/usr/bin/env bash
# scripts/_env.sh  --  SOURCE ONLY (`. scripts/_env.sh`); do not run directly.
# Loads RM_HOST_USB / RM_HOST_WIFI / RM_ROOT_PASSWORD from the gitignored
# secrets file, and defines ssh/scp helpers. macOS + Linux compatible.

# Resolve repo paths relative to this script.
_THIS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
# shellcheck source=/dev/null
. "$_THIS_DIR/paths.sh"
_SECRETS="$_THIS_DIR/../secrets/remarkable.local.env"

if [ ! -f "$_SECRETS" ]; then
  echo "ERROR: missing secrets file: $_SECRETS" >&2
  echo "  cp secrets/remarkable.local.env.example secrets/remarkable.local.env" >&2
  echo "  then fill in RM_ROOT_PASSWORD (and RM_HOST_WIFI)." >&2
  return 1 2>/dev/null || exit 1
fi

# Red error helper: bold-red on a TTY, plain when piped/tee'd to a recon log.
# Usage: err "message"   (prints to stderr)
if [ -t 2 ]; then RM_RED=$'\033[1;31m'; RM_RST=$'\033[0m'; else RM_RED=''; RM_RST=''; fi
err() { printf '%sERROR:%s %s\n' "$RM_RED" "$RM_RST" "$*" >&2; }
export RM_RED RM_RST

# Extract a KEY's value: strip inline comments + trailing whitespace.
_get_env() {
  sed -n -E "s/^[[:space:]]*$1[[:space:]]*=[[:space:]]*([^#]*).*/\1/p" "$_SECRETS" \
    | head -n1 | sed -E 's/[[:space:]]+$//'
}

RM_HOST_USB="$(_get_env RM_HOST_USB)";   : "${RM_HOST_USB:=10.11.99.1}"
RM_HOST_WIFI="$(_get_env RM_HOST_WIFI)"
RM_ROOT_PASSWORD="$(_get_env RM_ROOT_PASSWORD)"

# Default device target. USB-ethernet is currently inactive (no DHCP lease) on
# the Mac, so prefer Wi-Fi when an IP is recorded. Override per call with the
# optional target arg, or globally:  export RM_HOST=10.11.99.1
RM_HOST="${RM_HOST:-${RM_HOST_WIFI:-$RM_HOST_USB}}"
export RM_HOST_USB RM_HOST_WIFI RM_ROOT_PASSWORD RM_HOST

# Hardened SSH/SCP transport options shared by every helper below.
# (Word-split on purpose -- keep $RM_SSH_OPTS unquoted at the call sites.)
#
#   ServerAliveInterval=5 / CountMax=3 -- send a keepalive every 5s; abort after
#     3 misses (~15s). Keeps a marginal link alive AND fails fast (instead of
#     hanging forever) if it truly drops.
#   IPQoS=throughput -- drop the macOS-default 'lowdelay' DSCP marking that some
#     APs mishandle. (NOTE: on its own this did NOT cure the scp stall below --
#     see rm_send_file for the actual fix; kept because it's harmless + correct.)
RM_SSH_OPTS="-o StrictHostKeyChecking=accept-new -o IPQoS=throughput -o ServerAliveInterval=5 -o ServerAliveCountMax=3"
export RM_SSH_OPTS

# Run a command on the device. Usage: rm_ssh "<cmd>" [target]
rm_ssh() {
  local target="${2:-$RM_HOST}"
  ssh $RM_SSH_OPTS -o ConnectTimeout=8 "root@$target" "$1"
}

# Copy local->device. Usage: rm_scp_to <local> <remote> [target]
rm_scp_to() {
  local target="${3:-$RM_HOST}"
  scp $RM_SSH_OPTS "$1" "root@$target:$2"
}

# Copy device->local. Usage: rm_scp_from <remote> <local> [target]
rm_scp_from() {
  local target="${3:-$RM_HOST}"
  scp $RM_SSH_OPTS "root@$target:$1" "$2"
}

# Robustly copy a LOCAL file to the device.
# Usage: rm_send_file <local> <remote> [target]
#
# Transport = a gzip stream through the SSH channel itself (NOT scp/SFTP):
#     gzip -c <src> | ssh ... "gzip -dc > <dst>"
# Why not scp: on the Mac -> Wi-Fi -> tablet link, scp reliably wedges at a
# fixed ~255 KB offset and never recovers -- the SFTP subsystem's app-level
# windowing deadlocks on this AP / older-kernel TCP path (IPQoS + keepalives did
# NOT clear it; only Ctrl-C + this stream did). A raw exec stream has no SFTP
# windowing -- it just blasts bytes SSH frames itself -- and sails through.
# Bonus: gzip ~halves the bytes on the wire. The device has BusyBox gzip (the
# keywriter deploy already ships a .tar.gz). A final byte-count check catches a
# silently truncated stream before it can leave a corrupt binary in place.
rm_send_file() {
  local src="$1" dst="$2" target="${3:-$RM_HOST}"
  gzip -c "$src" | ssh $RM_SSH_OPTS -o ConnectTimeout=8 "root@$target" "gzip -dc > '$dst'" || return 1
  # Verify the on-device size matches the source (stat: macOS -f%z / Linux -c%s).
  local lsize rsize
  lsize=$(stat -f%z "$src" 2>/dev/null || stat -c%s "$src" 2>/dev/null)
  rsize=$(rm_ssh "wc -c < '$dst'" "$target" | tr -d '[:space:]')
  if [ -n "$lsize" ] && [ "$lsize" != "$rsize" ]; then
    echo "  ERROR: size mismatch after copy ($dst: local=$lsize remote=$rsize)" >&2
    return 1
  fi
  return 0
}

# Silent test: does passwordless key login work? Returns 0/1. Usage: rm_test_key [target]
rm_test_key() {
  local target="${1:-$RM_HOST}"
  ssh $RM_SSH_OPTS -o BatchMode=yes -o ConnectTimeout=5 "root@$target" "exit 0" >/dev/null 2>&1
}

# with_ticker <expected_secs> <cmd...>
#
# Runs <cmd...> in the FOREGROUND while a background ticker prints one '*'
# per (expected_secs / 10) seconds. When the command finishes, kills the
# ticker and appends the elapsed time + newline. Returns the command's real
# exit code -- safe under set -e (cmd is foreground, not a subshell).
#
# ASCII + append-only: no \r, so output tees cleanly into docs/recon/ logs.
#
# Asterisk count = health gauge:
#   ~10 = normal | 8-9 = fast | 12-15 = sluggish link | 20+ = something off
#
# Ctrl-C safety: in a non-job-control script, a backgrounded '&' child inherits
# SIG_IGN for SIGINT -- so a naive ticker would IGNORE the user's Ctrl-C, get
# orphaned when the script's trap exits, and keep printing '*' at the shell
# prompt forever. Two guards prevent that: the ticker traps INT/TERM to exit,
# and it self-terminates once the parent script PID ($parent) is gone.
with_ticker() {
  local expected="$1"; shift
  local interval parent=$$
  interval=$(awk "BEGIN{v=$expected/10; print (v<0.5)?0.5:v}")
  ( trap 'exit 0' INT TERM
    while kill -0 "$parent" 2>/dev/null; do printf '*'; sleep "$interval"; done ) &
  local _ticker=$!
  local _t0; _t0=$(date +%s)
  "$@"
  local rc=$?
  kill "$_ticker" 2>/dev/null; wait "$_ticker" 2>/dev/null || true
  printf ' %ds\n' "$(( $(date +%s) - _t0 ))"
  return "$rc"
}
