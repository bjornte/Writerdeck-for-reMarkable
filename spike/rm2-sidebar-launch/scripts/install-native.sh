#!/usr/bin/env bash
# Install native Writerdeck sidebar hook (XOVI + our .qmd). Removes AppLoad.
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$ROOT/scripts/_env.sh"

CACHE="$SPIKE/.cache"
XOVI_URL="${XOVI_URL:-https://github.com/asivery/rm-xovi-extensions/releases/download/v19-23052026/xovi-arm32.tar.gz}"
HOST="${1:-}"
if [ -z "$HOST" ]; then
  HOST="$(bash "$ROOT/spike/rm2-official-sdk/scripts/rm2-pick-host.sh")"
fi
export RM_HOST="$HOST"
echo "rM2 host: $HOST"

mkdir -p "$CACHE"
XOVI_TGZ="$CACHE/xovi-arm32.tar.gz"
if [ ! -s "$XOVI_TGZ" ]; then
  echo "Downloading XOVI arm32..."
  curl -fsSL -o "$XOVI_TGZ" "$XOVI_URL"
fi

ICON_SRC="$SPIKE/appload-app/icon.png"
QMD_SRC="$SPIKE/qml/writerdeck-sidebar.qmd"
OS_VER="$(rm_ssh 'sed -n "s/^IMG_VERSION=//p" /etc/os-release | tr -d "\""' "$HOST" | tr -d '\r' || true)"
case "$OS_VER" in
  3.28.*) QMD_SRC="$SPIKE/qml/writerdeck-sidebar-3.28.qmd" ;;
esac
echo "qmd: $(basename "$QMD_SRC") (os=${OS_VER:-unknown})"
RCC_SRC="$SPIKE/resources/writerdeck-icons.rcc"
ENSURE_SRC="$SPIKE/scripts/writerdeck-ensure-sidebar.sh"
LAUNCH_SRC="$SPIKE/scripts/writerdeck-sidebar-launch.sh"
TESTED_SRC="$SPIKE/tested-os.json"
test -f "$ICON_SRC"
test -f "$QMD_SRC"
test -f "$RCC_SRC"
test -f "$ENSURE_SRC"
test -f "$LAUNCH_SRC"
test -f "$TESTED_SRC"

echo "Pushing XOVI archive (if needed)..."
rm_send_file "$XOVI_TGZ" /tmp/xovi-arm32.tar.gz "$HOST"

echo "Installing XOVI baseline..."
rm_ssh '
set -e
if [ ! -x /home/root/xovi/start ]; then
  tar -xzf /tmp/xovi-arm32.tar.gz -C /home/root
fi
test -x /home/root/xovi/start
test -f /home/root/xovi/extensions.d/qt-resource-rebuilder.so
mkdir -p /home/root/xovi/exthome/qt-resource-rebuilder \
  /home/root/xovi/inactive-extensions
# Drop AppLoad -- we own the sidebar hook now
if [ -f /home/root/xovi/extensions.d/appload.so ]; then
  mv -f /home/root/xovi/extensions.d/appload.so \
    /home/root/xovi/inactive-extensions/appload.so
  echo "moved appload.so to inactive"
fi
# Enable CommandExecutor for QML shell-out
if [ -f /home/root/xovi/inactive-extensions/qt-command-executor.so ]; then
  mv -f /home/root/xovi/inactive-extensions/qt-command-executor.so \
    /home/root/xovi/extensions.d/qt-command-executor.so
fi
test -f /home/root/xovi/extensions.d/qt-command-executor.so
# Remove AppLoad app tree if present
rm -rf /home/root/xovi/exthome/appload
ls -la /home/root/xovi/extensions.d/
' "$HOST"

echo "Deploying Writerdeck sidebar .qmd + icon .rcc + ensure script..."
rm_send_file "$QMD_SRC" \
  /home/root/xovi/exthome/qt-resource-rebuilder/writerdeck-sidebar.qmd "$HOST"
rm_send_file "$RCC_SRC" \
  /home/root/xovi/exthome/qt-resource-rebuilder/writerdeck-icons.rcc "$HOST"
rm_send_file "$ICON_SRC" \
  /home/root/xovi/exthome/qt-resource-rebuilder/writerdeck-icon.png "$HOST"
rm_send_file "$TESTED_SRC" \
  /home/root/xovi/exthome/qt-resource-rebuilder/tested-os.json "$HOST"
rm_send_file "$ENSURE_SRC" /home/root/writerdeck-ensure-sidebar.sh "$HOST"
rm_send_file "$LAUNCH_SRC" /home/root/writerdeck-sidebar-launch.sh "$HOST"
rm_ssh 'chmod a+x /home/root/writerdeck-ensure-sidebar.sh /home/root/writerdeck-sidebar-launch.sh' "$HOST"
if [ "${ACCEPT_UNTESTED:-}" = "1" ]; then
  rm_ssh 'touch /home/root/xovi/exthome/qt-resource-rebuilder/accept-untested' "$HOST"
  echo "wrote accept-untested flag (candidates allowed)"
fi

# Hashtab: rebuild if missing
HAS="$(rm_ssh 'if test -s /home/root/xovi/exthome/qt-resource-rebuilder/hashtab; then echo yes; else echo no; fi' "$HOST" | tr -d '\r')"
if [ "$HAS" != "yes" ]; then
  echo "Rebuilding QML hashtab non-interactively..."
  echo "If the tablet shows a lock screen, unlock it once while this runs."
  rm_ssh '
set -e
qrr_so=/home/root/xovi/extensions.d/qt-resource-rebuilder.so
qrr_dir=/home/root/xovi/exthome/qt-resource-rebuilder
hashtab="$qrr_dir/hashtab"
systemctl stop xochitl 2>/dev/null || true
killall xochitl 2>/dev/null || true
sleep 1
export XOVI_ROOT=/tmp/xovi-hashtab-build
rm -rf "$XOVI_ROOT"
mkdir -p "$XOVI_ROOT/extensions.d" "$qrr_dir"
ln -sf "$qrr_so" "$XOVI_ROOT/extensions.d/qt-resource-rebuilder.so"
LOG=/tmp/hashtab-rebuild.log
rm -f "$LOG" "$hashtab"
QMLDIFF_HASHTAB_CREATE="$hashtab" \
QML_DISABLE_DISK_CACHE=1 \
LD_PRELOAD=/home/root/xovi/xovi.so \
/usr/bin/xochitl >"$LOG" 2>&1 &
XOCH=$!
ok=0
for i in $(seq 1 90); do
  if grep -q "Hashtab saved to" "$LOG" 2>/dev/null; then ok=1; break; fi
  if ! kill -0 "$XOCH" 2>/dev/null; then break; fi
  sleep 1
done
kill -15 "$XOCH" 2>/dev/null || true
killall xochitl 2>/dev/null || true
sleep 1
rm -rf "$XOVI_ROOT"
test "$ok" = 1
test -s "$hashtab"
systemctl start xochitl
' "$HOST"
else
  echo "hashtab already present"
fi

echo
echo "Native install done on $HOST."
echo "Arming sidebar (also runs on writerdeck service start / before xochitl)..."
rm_ssh '/home/root/writerdeck-ensure-sidebar.sh' "$HOST" || true
echo "Sidebar hooks Writerdeck start. Manual: bash spike/rm2-sidebar-launch/scripts/enable.sh"
echo "Tap Writerdeck in the stock sidebar to open Lobby."
