#!/usr/bin/env bash
# Launch Writerdeck on official Qt6 epaper QPA (no linuxfb, no bundled Qt5).
# Deployed to /home/root/Writerdeck-launcher.sh and called by Writerdeck-server:
#   /home/root/Writerdeck-server --editor /home/root/Writerdeck-launcher.sh
#
# USB keyboard layout: reads keyboardLayout from settings.json (default us),
# applies the matching .qmap from /home/root/keymaps/ on a USB keyboard device
# only (never gpio-keys /dev/input/event1 - grab=1 without a device path makes
# Qt evdev grab event1 and breaks Home/Power). If no USB keyboard is present at
# launch, keymap is skipped; Writerdeck-server restarts the editor when one is
# plugged in so the layout applies.
#
# Touch: rM1 rotate=180; rM2 rotate=180:invertx. Do not overwrite a full string
# already in the environment.
#
# Deploy: bash scripts/deploy-keywriter-rm1.sh (rM1) or deploy-keywriter-rm2.sh.

set -euo pipefail

export HOME=/home/root
export QT_QUICK_BACKEND=epaper

if [ -z "${QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS:-}" ]; then
  model="$(cat /sys/devices/soc0/machine 2>/dev/null || true)"
  case "$model" in
    *"reMarkable 2"*) export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS="rotate=180:invertx" ;;
    *) export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS="rotate=180" ;;
  esac
fi

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
    echo "Writerdeck-launcher: USB layout=${layout} dev=${kb_dev} keymap=${qmap}" >&2
  else
    echo "Writerdeck-launcher: USB layout=${layout} (no keyboard yet - keymap skipped)" >&2
  fi
fi

cd /home/root
exec systemd-inhibit \
    --what=sleep \
    --why=Writerdeck \
    --mode=block \
    /home/root/Writerdeck -platform epaper
