#!/usr/bin/env bash
# scripts/deploy-keywriter-rm2.sh -- Deploy Qt6 Writerdeck binary to reMarkable 2.
#
# Uses the Codex SDK cross-compile from spike/rm2-official-sdk (fork probe today;
# production CMake build in Writerdeck-keywriter when it lands).
#
# Device layout:
#   /home/root/Writerdeck              -- Qt6 ARM binary (stock Qt6 on device)
#   /home/root/Writerdeck-launcher.sh  -- epaper QPA env (Writerdeck-launcher-rm2.sh)
#
# Usage (repo root, bash):
#   bash scripts/deploy-keywriter-rm2.sh
#   bash scripts/deploy-keywriter-rm2.sh 10.11.99.1
#   bash scripts/deploy-keywriter-rm2.sh -b   # binary + launcher only (~fast)

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$DIR/.." && pwd)"
SPIKE="$REPO/spike/rm2-official-sdk"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
# shellcheck source=/dev/null
. "$DIR/migrate-device-layout.sh"

TARGET=""
RM_BINARY_ONLY="${RM_BINARY_ONLY:-0}"
for arg in "$@"; do
  case "$arg" in
    -b|--binary) RM_BINARY_ONLY=1 ;;
    *)           TARGET="$arg" ;;
  esac
done

if [ -z "$TARGET" ]; then
  TARGET="$(bash "$SPIKE/scripts/rm2-pick-host.sh")"
fi
export RM2_HOST="$TARGET"
export RM_HOST="$TARGET"

PROBE_BIN="$SPIKE/.cache/out/Writerdeck_fork_probe"
LAUNCHER_SRC="$DIR/Writerdeck-launcher-rm2.sh"

echo "=== deploy-keywriter-rm2  target=$TARGET ==="
echo

if [ ! -f "$PROBE_BIN" ]; then
  echo "Building Writerdeck Qt6 binary (fork probe) ..."
  bash "$SPIKE/scripts/build-fork-probe.sh"
fi
if [ ! -f "$PROBE_BIN" ]; then
  err "missing $PROBE_BIN after build"
  exit 1
fi
if [ ! -f "$LAUNCHER_SRC" ]; then
  err "missing $LAUNCHER_SRC"
  exit 1
fi

echo "--- Pre-flight ---"
echo "  binary : $(ls -lh "$PROBE_BIN" | awk '{print $5, $NF}')"
echo "  launcher: $LAUNCHER_SRC"
echo

echo "--- Testing SSH to $TARGET ---"
if ! ping -c1 -W2 "$TARGET" >/dev/null 2>&1; then
  err "$TARGET unreachable (ping failed). Wake tablet or use USB (10.11.99.1)."
  exit 1
fi
if ! rm_test_key "$TARGET"; then
  err "key SSH to root@$TARGET failed."
  exit 1
fi
echo "  OK"
echo

migrate_device_layout "$TARGET"

echo "--- Stopping running Writerdeck ---"
rm_ssh 'wget -q -O /dev/null --post-data="" http://127.0.0.1:8000/api/flush-save 2>/dev/null || true; for p in $(pidof Writerdeck 2>/dev/null); do kill -TERM "$p" 2>/dev/null; done; i=0; while pidof Writerdeck >/dev/null 2>&1 && [ "$i" -lt 30 ]; do sleep 0.2; i=$((i+1)); done; for p in $(pidof Writerdeck 2>/dev/null); do kill -KILL "$p" 2>/dev/null; done; sleep 0.3; true' "$TARGET"

echo "--- Deploying Writerdeck -> /home/root/Writerdeck ---"
rm_send_file "$PROBE_BIN" "/home/root/Writerdeck.new" "$TARGET"
rm_ssh 'mv -f /home/root/Writerdeck.new /home/root/Writerdeck && chmod +x /home/root/Writerdeck' "$TARGET"
echo "  OK"

echo "--- Deploying launcher -> /home/root/Writerdeck-launcher.sh ---"
rm_send_file "$LAUNCHER_SRC" "/home/root/Writerdeck-launcher.sh" "$TARGET"
rm_ssh 'chmod +x /home/root/Writerdeck-launcher.sh' "$TARGET"
rm_send_file "$LAUNCHER_SRC" "/home/root/Writerdeck-launcher-rm2.sh" "$TARGET"
rm_ssh 'chmod +x /home/root/Writerdeck-launcher-rm2.sh' "$TARGET"
echo "  OK"

echo "--- Deploying keymaps -> /home/root/keymaps/ ---"
rm_ssh 'mkdir -p /home/root/keymaps' "$TARGET"
for qmap in "$REPO/keymaps"/*.qmap; do
  [ -f "$qmap" ] || continue
  rm_send_file "$qmap" "/home/root/keymaps/$(basename "$qmap")" "$TARGET"
done
echo "  OK"

rm_deploy_wd "$TARGET"
echo "  ${DEVICE_WD} OK"

LOBBY_UI_SRC="$REPO/config/lobby-ui.json"
if [ -f "$LOBBY_UI_SRC" ]; then
  echo "--- Lobby UI config ---"
  rm_ssh "mkdir -p '${DEVICE_SETTINGS_DIR}'" "$TARGET"
  if rm_ssh "[ -f '${DEVICE_LOBBY_UI_FILE}' ] && echo EXISTS || echo MISSING" "$TARGET" | grep -q EXISTS; then
    echo "  already on tablet (left unchanged)"
  else
    rm_send_file "$LOBBY_UI_SRC" "${DEVICE_LOBBY_UI_FILE}" "$TARGET"
    echo "  seeded from config/lobby-ui.json"
  fi
fi

LOBBY_I18N_SRC="$REPO/config/lobby-ui-i18n"
if [ -d "$LOBBY_I18N_SRC" ] && [ "$RM_BINARY_ONLY" != "1" ]; then
  echo "--- Lobby UI i18n ---"
  rm_ssh "mkdir -p '${DEVICE_LOBBY_UI_I18N_DIR}'" "$TARGET"
  for pack in "$LOBBY_I18N_SRC"/*.json; do
    [ -f "$pack" ] || continue
    rm_send_file "$pack" "${DEVICE_LOBBY_UI_I18N_DIR}/$(basename "$pack")" "$TARGET"
  done
  echo "  OK"
fi

FONTS_SRC="$REPO/config/fonts"
if [ -d "$FONTS_SRC" ]; then
  echo "--- Lobby fonts (Serif + symbol fallback) ---"
  rm_ssh 'mkdir -p /usr/share/fonts/ttf/remarkable'
  for ff in reMarkableSerif.ttf NimbusSans-Regular.otf; do
    [ -f "$FONTS_SRC/$ff" ] || continue
    rm_send_file "$FONTS_SRC/$ff" "/usr/share/fonts/ttf/remarkable/$ff" "$TARGET"
  done
  rm_ssh 'fc-cache -f /usr/share/fonts/ttf/remarkable 2>/dev/null || true'
  echo "  OK"
fi

echo "--- Notes dir ---"
rm_ssh 'mkdir -p /home/root/Writerdeck-user-documents' "$TARGET"
echo "  OK"
echo

if [ "$RM_BINARY_ONLY" = "1" ]; then
  echo "Binary-only deploy done."
  echo "Next: bash scripts/install-service.sh --start $TARGET"
  echo "      bash scripts/test-edit-session.sh $TARGET"
  exit 0
fi

echo "======================================"
echo "  DEPLOY DONE (rM2 Qt6 editor)"
echo "======================================"
echo "  Binary  : /home/root/Writerdeck"
echo "  Launcher: /home/root/Writerdeck-launcher.sh"
echo "  Next    : bash scripts/deploy-rmkbd.sh (server)"
echo "            bash scripts/install-service.sh --start $TARGET"
echo "            bash scripts/test-edit-session.sh $TARGET"
echo "======================================"
