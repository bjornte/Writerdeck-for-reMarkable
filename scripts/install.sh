#!/usr/bin/env bash
# scripts/install.sh -- First-time (idempotent) install chain for Writerdeck.
#
# Runs: ensure-secrets -> bootstrap -> fetch -> deploy editor -> deploy server
# -> install-service. With --start: also start, health-check the phone UI, then
# enable autostart on boot, then configure-sync when secrets have sync set.
#
# Usage (run from repo root on the Mac):
#   bash scripts/install.sh
#   bash scripts/install.sh --start    # start + enable after health check
#
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"

DO_START=0
for arg in "$@"; do
  case "$arg" in
    --start) DO_START=1 ;;
    -h|--help)
      sed -n '2,12p' "$0"
      exit 0
      ;;
    *)
      echo "ERROR: unknown arg: $arg (try --start)" >&2
      exit 1
      ;;
  esac
done

cd "$REPO"

# Friendlier prompts + quiet sub-scripts (their chatter goes to INSTALL_LOG).
export WRITERDECK_INSTALL=1

TS="$(date +%Y-%m-%dT%H-%M-%S)"
RECON_DIR="$REPO/docs/recon"
mkdir -p "$RECON_DIR"
INSTALL_LOG="$RECON_DIR/install-$TS.txt"

_progress_msg=""
_progress_start() {
  _progress_msg="$1"
  if [ -t 1 ]; then
    printf '  %s' "$_progress_msg"
  else
    printf '  ... %s\n' "$_progress_msg"
  fi
}

_progress_ok() {
  local done="$1"
  if [ -t 1 ]; then
    printf '\r\033[K  OK  %s\n' "$done"
  else
    printf '  OK  %s\n' "$done"
  fi
  _progress_msg=""
}

_progress_fail() {
  if [ -t 1 ]; then
    printf '\r\033[K  FAIL  %s\n' "${_progress_msg:-step}"
  else
    printf '  FAIL  %s\n' "${_progress_msg:-step}"
  fi
  echo
  echo "Details: $INSTALL_LOG" >&2
  echo "Last lines:" >&2
  tail -n 40 "$INSTALL_LOG" >&2 || true
}

# Run a command with output appended to INSTALL_LOG; update the in-place line.
_run_step() {
  local pending="$1" done_msg="$2"
  shift 2
  _progress_start "$pending"
  if "$@" >>"$INSTALL_LOG" 2>&1; then
    _progress_ok "$done_msg"
  else
    _progress_fail
    exit 1
  fi
}

_get() {
  sed -n -E "s/^[[:space:]]*$1[[:space:]]*=[[:space:]]*([^#]*).*/\1/p" \
    "$REPO/secrets/remarkable.local.env" | head -n1 | sed -E 's/[[:space:]]+$//'
}

{
  echo "=== install.sh $TS ==="
  echo "cwd=$REPO"
  echo
} >"$INSTALL_LOG"

echo "Writerdeck install"
echo "=================="
echo
echo "This will set up Writerdeck on your reMarkable over Wi-Fi."
echo "You'll answer a few questions, then the rest runs on its own"
echo "(a few minutes; mostly waiting on the tablet)."
echo
echo "Your answers are saved in secrets/remarkable.local.env"
echo "(local only -- never committed)."
echo

bash "$DIR/ensure-secrets.sh"

# shellcheck source=/dev/null
. "$DIR/_env.sh"

echo
echo "Checking your Mac and tablet"
echo "----------------------------"

_ok_line() { printf '  OK  %s\n' "$1"; }

if [ -f "$REPO/secrets/remarkable.local.env" ]; then
  _ok_line "Secrets file ready"
else
  echo "  FAIL  Secrets file missing" >&2
  exit 1
fi

if ping -c1 -W2 "$RM_HOST" >/dev/null 2>&1; then
  _ok_line "Tablet reachable"
else
  echo "  FAIL  Tablet not reachable at $RM_HOST" >&2
  exit 1
fi

if command -v go >/dev/null 2>&1; then
  _ok_line "Ready for local builds"
else
  _ok_line "Ready (will download binaries from GitHub)"
fi

_progress_start "Setting up SSH access..."
if bash "$DIR/bootstrap.sh" >>"$INSTALL_LOG" 2>&1; then
  _progress_ok "SSH access OK"
else
  _progress_fail
  exit 1
fi

echo
echo "Downloading and installing"
echo "--------------------------"

# Server first (draft order), then editor, then systemd.
if [ -f "$REPO/Writerdeck-server" ]; then
  _ok_line "Server downloaded"
else
  _run_step "Downloading server..." "Server downloaded" \
    bash "$DIR/deploy-rmkbd.sh" --build-only
fi

_run_step "Installing on tablet..." "Server installed on tablet" \
  bash "$DIR/deploy-rmkbd.sh" --deploy-only

DIST="$REPO/third_party/keywriter/dist"
if [ -f "$DIST/Writerdeck" ] && [ -f "$DIST/qt5.tar.gz" ]; then
  _ok_line "Editor downloaded"
else
  _run_step "Downloading editor..." "Editor downloaded" \
    bash "$DIR/fetch-keywriter-dist.sh"
fi

_run_step "Copying editor and Qt to tablet (~20s)..." "Editor and Qt copied to tablet" \
  bash "$DIR/deploy-keywriter.sh"

if [ "$DO_START" -eq 1 ]; then
  _run_step "Installing Always-on service (systemd unit)..." \
    "Always-on service (systemd unit) installed" \
    bash "$DIR/install-service.sh" --start
  _ok_line "Writerdeck started"
  _ok_line "Set to start on boot"
else
  _run_step "Installing Always-on service (systemd unit)..." \
    "Always-on service (systemd unit) installed" \
    bash "$DIR/install-service.sh"
fi

SYNC_REPO="$(_get SYNC_REPO)"
GH_TOKEN="$(_get GH_TOKEN)"
if [ "$DO_START" -eq 1 ]; then
  if [ -n "$SYNC_REPO" ] && [ -n "$GH_TOKEN" ]; then
    echo
    echo "Notes sync"
    echo "----------"
    _run_step "Enabling sync for $SYNC_REPO..." "Sync enabled for $SYNC_REPO" \
      bash "$DIR/configure-sync.sh"
    _ok_line "GitHub token checked from the tablet"
    echo
    echo "  On a new phone or Wi-Fi address you may still need to"
    echo "  paste the token once under Sync setup on the phone page."
  else
    # Still push PIN_DIGITS=none (and no-op sync) without a Notes sync section.
    bash "$DIR/configure-sync.sh" >>"$INSTALL_LOG" 2>&1
  fi
fi

echo
echo "Install complete"
echo "================"
echo
echo "  Phone UI:  http://${RM_HOST}:8000/"
echo
if [ "$DO_START" -ne 1 ]; then
  echo "  Next -- start + enable autostart:"
  echo "    bash scripts/install.sh --start"
  echo
fi
echo "  You're done when:"
echo "    - the tablet shows the normal reMarkable home screen"
echo "    - the phone page lists your notes (not an empty shell)"
echo "    - the connection bar says Connected (not stuck on connecting...)"
echo
echo "  To launch the Writerdeck:"
echo "    - press both page buttons simultaneously"
echo "    - on a USB keyboard, press Esc"
echo
echo "  If something goes wrong after a reboot:"
echo "    systemctl disable --now writerdeck && systemctl start xochitl"
echo
echo "  Full log: $INSTALL_LOG"
echo
