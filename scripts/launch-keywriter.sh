#!/usr/bin/env bash
# scripts/launch-keywriter.sh -- Launch keywriter with the proven linuxfb env.
#
# This is the ONE authoritative place for the keywriter launch environment.
# It is deployed to /home/root/launch-keywriter.sh and called:
#   - by rmkbd in supervisor mode: /home/root/rmkbd --editor /home/root/launch-keywriter.sh
#   - (replaces the inline KW_ENV blocks in test-phase*.sh)
#
# Deploy: bash scripts/deploy-keywriter.sh (copies this alongside the binary).
#
# Environment (all confirmed working -- see Phase 1 Step 5):
#   linuxfb + size/mmsize: rM1 panel is 1404x1872 px / 158x210 mm = 226 DPI.
#     Without size= the window falls back to ~1024x768 landscape upper-left.
#     Without mmsize= the DPI is guessed low and text comes out ~4 px.
#   QMLSCENE_DEVICE=epaper: the libqsgepaper scene graph drives the EPDC
#     refresh. The 'epaper' QPA *platform* is absent in the toltec toolchain;
#     only linuxfb is available as the QPA platform.
#   QT_FONT_DPI=226: overrides Qt's DPI calculation so point-sized fonts
#     render at the correct physical size on the e-ink panel.

set -euo pipefail

# HOME must be /home/root: keywriter builds its save/load folder from
# QStandardPaths::HomeLocation (= $HOME) -> file://$HOME/edit/scratch.md.
# Under systemd, root's HOME defaults to "/", which collapses the path to
# //edit/scratch.md and the save PUT fails (QNetworkReplyImplPrivate::error).
# Interactive SSH / the test scripts set HOME=/home/root so they never hit
# this; the systemd service does, so we pin it here (the one launch-env place).
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
    --why=rM1-Writerdeck \
    --mode=block \
    /home/root/keywriter
