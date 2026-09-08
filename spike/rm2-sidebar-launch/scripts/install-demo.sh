#!/usr/bin/env bash
# Install tethered XOVI + AppLoad + Writerdeck AppLoad entry on rM2 (arm32).
# Does NOT enable the UI mod until enable.sh passes the OS gate.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"

CACHE="$SPIKE/.cache"
XOVI_URL="${XOVI_URL:-https://github.com/asivery/rm-xovi-extensions/releases/download/v19-23052026/xovi-arm32.tar.gz}"
APPLOAD_URL="${APPLOAD_URL:-https://github.com/asivery/rm-appload/releases/download/v0.5.3/appload-arm32.zip}"
HOST="${1:-}"
if [ -z "$HOST" ]; then
  HOST="$(bash "$ROOT/spike/rm2-official-sdk/scripts/rm2-pick-host.sh")"
fi
export RM_HOST="$HOST"
echo "rM2 host: $HOST"

mkdir -p "$CACHE"
XOVI_TGZ="$CACHE/xovi-arm32.tar.gz"
APPLOAD_ZIP="$CACHE/appload-arm32.zip"

if [ ! -s "$XOVI_TGZ" ]; then
  echo "Downloading XOVI arm32..."
  curl -fsSL -o "$XOVI_TGZ" "$XOVI_URL"
fi
if [ ! -s "$APPLOAD_ZIP" ]; then
  echo "Downloading AppLoad arm32..."
  curl -fsSL -o "$APPLOAD_ZIP" "$APPLOAD_URL"
fi

echo "Extracting AppLoad on Mac (tablet may lack unzip)..."
APPLOAD_DIR="$CACHE/appload-extract"
rm -rf "$APPLOAD_DIR"
mkdir -p "$APPLOAD_DIR"
unzip -qo "$APPLOAD_ZIP" -d "$APPLOAD_DIR"
test -f "$APPLOAD_DIR/appload.so"

echo "Pushing XOVI archive..."
rm_send_file "$XOVI_TGZ" /tmp/xovi-arm32.tar.gz "$HOST"

echo "Installing XOVI (skip if already present)..."
rm_ssh '
set -e
if [ ! -x /home/root/xovi/start ]; then
  tar -xzf /tmp/xovi-arm32.tar.gz -C /home/root
fi
test -x /home/root/xovi/start
test -f /home/root/xovi/extensions.d/qt-resource-rebuilder.so
mkdir -p /home/root/xovi/exthome/qt-resource-rebuilder \
  /home/root/xovi/exthome/appload /home/root/shims
' "$HOST"

echo "Installing AppLoad .so + shims..."
rm_send_file "$APPLOAD_DIR/appload.so" /home/root/xovi/extensions.d/appload.so "$HOST"
rm_send_file "$APPLOAD_DIR/shims/qtfb-shim.so" /home/root/shims/qtfb-shim.so "$HOST"
rm_send_file "$APPLOAD_DIR/shims/qtfb-shim-32bit.so" /home/root/shims/qtfb-shim-32bit.so "$HOST"
rm_ssh 'ls -la /home/root/xovi/extensions.d/; ls -la /home/root/shims/' "$HOST"

echo "Deploying Writerdeck AppLoad entry..."
rm_ssh 'mkdir -p /home/root/xovi/exthome/appload/writerdeck' "$HOST"
rm_send_file "$SPIKE/appload-app/external.manifest.json" \
  /home/root/xovi/exthome/appload/writerdeck/external.manifest.json "$HOST"
rm_send_file "$SPIKE/appload-app/icon.png" \
  /home/root/xovi/exthome/appload/writerdeck/icon.png "$HOST"
rm_send_file "$SPIKE/appload-app/launch-writerdeck.sh" \
  /home/root/xovi/exthome/appload/writerdeck/launch-writerdeck.sh "$HOST"
rm_ssh 'chmod a+x /home/root/xovi/exthome/appload/writerdeck/launch-writerdeck.sh' "$HOST"

echo "Rebuilding QML hashtab non-interactively (AppLoad #68 / remagic#2)..."
echo "If the tablet shows a lock screen, unlock it once while this runs."
rm_ssh '
set -e
qrr_so=/home/root/xovi/extensions.d/qt-resource-rebuilder.so
qrr_dir=/home/root/xovi/exthome/qt-resource-rebuilder
hashtab="$qrr_dir/hashtab"
test -f "$qrr_so"
mkdir -p "$qrr_dir"
rm -f "$hashtab"
systemctl stop xochitl 2>/dev/null || true
killall xochitl 2>/dev/null || true
sleep 1
export XOVI_ROOT=/tmp/xovi-hashtab-build
rm -rf "$XOVI_ROOT"
mkdir -p "$XOVI_ROOT/extensions.d"
ln -sf "$qrr_so" "$XOVI_ROOT/extensions.d/qt-resource-rebuilder.so"
LOG=/tmp/hashtab-rebuild.log
rm -f "$LOG"
QMLDIFF_HASHTAB_CREATE="$hashtab" \
QML_DISABLE_DISK_CACHE=1 \
LD_PRELOAD=/home/root/xovi/xovi.so \
/usr/bin/xochitl >"$LOG" 2>&1 &
XOCH=$!
ok=0
for i in $(seq 1 90); do
  if grep -q "Hashtab saved to" "$LOG" 2>/dev/null; then
    ok=1
    break
  fi
  if ! kill -0 "$XOCH" 2>/dev/null; then
    break
  fi
  sleep 1
done
kill -15 "$XOCH" 2>/dev/null || true
killall xochitl 2>/dev/null || true
sleep 1
rm -rf "$XOVI_ROOT"
tail -30 "$LOG" || true
test "$ok" = 1
test -s "$hashtab"
ls -la "$hashtab"
systemctl start xochitl
' "$HOST"

echo
echo "Install done on $HOST."
echo "Next: ACCEPT_UNTESTED=1 bash spike/rm2-sidebar-launch/scripts/enable.sh"
echo "Then on tablet: unlock -> open left sidebar (hamburger) -> AppLoad -> Writerdeck."
echo "Revert anytime: reboot, or bash spike/rm2-sidebar-launch/scripts/disable.sh"
