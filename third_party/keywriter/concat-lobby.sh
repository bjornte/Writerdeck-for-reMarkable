#!/bin/sh
# Concatenate lobby/*.inc into lobby_subpages.qml.inc (source of truth: lobby/ dir).
set -e
cd "$(dirname "$0")"
OUT=lobby_subpages.qml.inc
cat lobby/lobby_shell_top.inc \
    lobby/lobby_home.inc \
    lobby/lobby_files.inc \
    lobby/lobby_keyboard.inc \
    lobby/lobby_sync.inc \
    lobby/lobby_settings.inc \
    lobby/lobby_shortcuts.inc \
    lobby/lobby_shell_bottom.inc > "$OUT"
echo "wrote $OUT ($(wc -l < "$OUT") lines)"
