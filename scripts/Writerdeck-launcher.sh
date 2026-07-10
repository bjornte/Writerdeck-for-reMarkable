#!/usr/bin/env bash
# scripts/Writerdeck-launcher.sh -- Launch Writerdeck with the proven linuxfb env.
#
# Authoritative launch environment for the e-ink editor binary (Writerdeck).
# Deployed to /home/root/Writerdeck-launcher.sh and called by Writerdeck-server:
#   /home/root/Writerdeck-server --editor /home/root/Writerdeck-launcher.sh
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

cd /home/root
exec systemd-inhibit \
    --what=sleep \
    --why=Writerdeck \
    --mode=block \
    /home/root/Writerdeck
