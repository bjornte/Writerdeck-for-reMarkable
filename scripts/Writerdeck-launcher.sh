#!/usr/bin/env bash
# scripts/Writerdeck-launcher.sh -- Launch Writerdeck with the proven linuxfb env.
#
# Authoritative launch environment for the e-ink editor binary (Writerdeck).
# Deployed to /home/root/Writerdeck-launcher.sh and called by Writerdeck-server:
#   /home/root/Writerdeck-server --editor /home/root/Writerdeck-launcher.sh
#
# USB keyboard layout: reads keyboardLayout from settings.json (default us),
# applies the matching .qmap from /home/root/keymaps/.  Omit the /dev/input/eventN
# prefix so Qt discovers the keyboard at runtime (hotplug-safe); pinning a device
# node at launch fails when the keyboard is plugged in after Writerdeck starts.
#
# Deploy: bash scripts/deploy-keywriter.sh (copies this alongside the binary).

set -euo pipefail

export HOME=/home/root

export LD_LIBRARY_PATH=/home/root/qt5/lib
export QML2_IMPORT_PATH=/home/root/qt5/qml
export QT_PLUGIN_PATH=/home/root/qt5/plugins
export QMLSCENE_DEVICE=epaper
export QT_QPA_PLATFORM=linuxfb:fb=/dev/fb0:size=1404x1872:mmsize=158x210
export QT_FONT_DPI=226
export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS=rotate=180
export QT_QPA_GENERIC_PLUGINS=evdevtablet

SETTINGS="/home/root/.Writerdeck/settings.json"
KEYMAPS="/home/root/keymaps"
DEFAULT_LAYOUT="us"

read_keyboard_layout() {
  local layout=""
  if [ -f "$SETTINGS" ]; then
    layout=$(grep -o '"keyboardLayout"[[:space:]]*:[[:space:]]*"[^"]*"' "$SETTINGS" 2>/dev/null \
      | sed 's/.*"keyboardLayout"[[:space:]]*:[[:space:]]*"//;s/"$//' | sed -n '1p')
  fi
  case "${layout:-$DEFAULT_LAYOUT}" in
    us|no) echo "${layout:-$DEFAULT_LAYOUT}" ;;
    *)     echo "$DEFAULT_LAYOUT" ;;
  esac
}

layout=$(read_keyboard_layout)
qmap="${KEYMAPS}/${layout}.qmap"
if [ -f "$qmap" ]; then
  export QT_QPA_EVDEV_KEYBOARD_PARAMETERS="keymap=${qmap}:grab=1"
  echo "Writerdeck-launcher: USB layout=${layout} keymap=${qmap}" >&2
fi

cd /home/root
exec systemd-inhibit \
    --what=sleep \
    --why=Writerdeck \
    --mode=block \
    /home/root/Writerdeck
