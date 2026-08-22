#!/usr/bin/env bash
# Launch Writerdeck on reMarkable 2 (official Qt6 epaper QPA).
# Deploy to /home/root/Writerdeck-launcher-rm2.sh when the Qt6 binary ships.
set -euo pipefail

export HOME=/home/root
export QT_QUICK_BACKEND=epaper
export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS="${QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS:-rotate=180:invertx}"

SETTINGS="/home/root/.Writerdeck/settings.json"
KEYMAPS="/home/root/keymaps"
DEFAULT_LAYOUT="us"
BUTTON_DEV="/dev/input/event1"

read_keyboard_layout() {
  local layout=""
  if [ -f "$SETTINGS" ]; then
    layout=$(grep -o '"keyboardLayout"[[:space:]]*:[[:space:]]*"[^"]*"' "$SETTINGS" 2>/dev/null \
      | sed 's/.*"keyboardLayout"[[:space:]]*:[[:space:]]*"//;s/"$//' | sed -n '1p')
  fi
  case "${layout:-$DEFAULT_LAYOUT}" in
    us|no|es|de|fr) echo "${layout:-$DEFAULT_LAYOUT}" ;;
    *)              echo "$DEFAULT_LAYOUT" ;;
  esac
}

find_usb_keyboard_dev() {
  local ev dev name
  for ev in /sys/class/input/event*; do
    [ -e "$ev" ] || continue
    dev="/dev/input/$(basename "$ev")"
    [ "$dev" = "$BUTTON_DEV" ] && continue
    name=$(cat "$ev/device/name" 2>/dev/null | tr '[:upper:]' '[:lower:]')
    case "$name" in
      *keyboard*) echo "$dev"; return 0 ;;
    esac
  done
  return 1
}

layout=$(read_keyboard_layout)
qmap="${KEYMAPS}/${layout}.qmap"
if [ -f "$qmap" ]; then
  if kb_dev=$(find_usb_keyboard_dev); then
    export QT_QPA_EVDEV_KEYBOARD_PARAMETERS="${kb_dev}:grab=1:keymap=${qmap}"
    echo "Writerdeck-launcher-rm2: USB layout=${layout} dev=${kb_dev}" >&2
  fi
fi

cd /home/root
exec systemd-inhibit \
  --what=sleep \
  --why=Writerdeck \
  --mode=block \
  /home/root/Writerdeck -platform epaper
