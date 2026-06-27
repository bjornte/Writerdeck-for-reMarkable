#!/usr/bin/env bash
# scripts/test-phase4.sh -- Phase 4 end-to-end device test.
#
# Proves the browser capture client: a human types from the Mac browser
# into keywriter on the e-ink.  No websocat needed -- the daemon now
# serves the capture page at http://<device>:<port>/.
#
# Procedure:
#   1. Build + deploy rmkbd (now embeds index.html)
#   2. Stop xochitl, launch patched keywriter, launch rmkbd daemon
#   3. Open the printed URL in the Mac browser, click to focus, type
#   4. Glance at the tablet; text should appear on the e-ink
#   5. (Optional) save: the editor saves scratch.md on mode toggle (Esc)
#   6. Restore xochitl (trap, always runs)
#
# Done-when: a paragraph typed in the Mac browser appears on the e-ink.
#
# Prerequisites:
#   - brew install go          (cross-build rmkbd)
#   - CI-built keywriter on device (run deploy-keywriter.sh at least once)
#
# Usage (from repo root on the Mac):
#   bash scripts/test-phase4.sh
#   bash scripts/test-phase4.sh -s            # skip build+scp (~2s loop); rmkbd already on device
#   bash scripts/test-phase4.sh -v            # verbose: log every translated key (keymap debugging)
#   bash scripts/test-phase4.sh 192.168.1.9   # explicit target

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
REPO="$(cd "$DIR/.." && pwd)"

# Args (any order): a target IP/host, -v|--verbose, -s|--skip (skip build+deploy).
VERBOSE=""
TARGET="$RM_HOST"
for arg in "$@"; do
    case "$arg" in
        -v|--verbose) VERBOSE="-v" ;;
        -s|--skip)    RM_SKIP_DEPLOY=1 ;;
        *)            TARGET="$arg" ;;
    esac
done
PORT=8000

RECON_DIR="$REPO/docs/recon"
mkdir -p "$RECON_DIR"
TS="$(date +%Y-%m-%dT%H-%M-%S)"
LOG="$RECON_DIR/test-phase4-$TS.txt"
exec > >(tee "$LOG") 2>&1

echo "=== test-phase4  $TS  target=$TARGET ==="
echo

# ---------------------------------------------------------------------------
# Trap: kill our processes cleanly so the next deploy is never ETXTBSY-blocked.
# ---------------------------------------------------------------------------
_restore() {
    echo
    echo "--- Cleanup ---"
    rm_ssh 'pkill -f /home/root/rmkbd 2>/dev/null; for p in $(pidof rmkbd); do kill "$p" 2>/dev/null; done; echo "  rmkbd stopped."' "$TARGET" || true
    rm_ssh 'pkill -f /home/root/keywriter 2>/dev/null; for p in $(pidof keywriter); do kill "$p" 2>/dev/null; done; echo "  keywriter stopped."' "$TARGET" || true
    rm_ssh 'systemctl start xochitl 2>/dev/null; echo "  xochitl restored."' "$TARGET" || true
}
trap _restore EXIT

RM_SKIP_DEPLOY="${RM_SKIP_DEPLOY:-0}"

if [ "$RM_SKIP_DEPLOY" = "1" ]; then
    echo "--- Skipping build+deploy (RM_SKIP_DEPLOY=1; using rmkbd already on device) ---"
    echo
else

# ---------------------------------------------------------------------------
# 1. Build rmkbd on the Mac (index.html is embedded at compile time).
# ---------------------------------------------------------------------------
echo "--- Building rmkbd (ARMv7 static, with embedded index.html) ---"
cd "$REPO/daemon"
go mod tidy
GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 \
    go build -trimpath -o "$REPO/rmkbd" .
echo "  built: $(file "$REPO/rmkbd" | sed 's/.*: //')"
cd "$REPO"
echo

# ---------------------------------------------------------------------------
# 2. Deploy rmkbd to device (kill-first + atomic mv to avoid ETXTBSY).
# ---------------------------------------------------------------------------
echo "--- Deploying rmkbd -> /home/root/rmkbd ---"
rm_ssh 'pkill -f /home/root/rmkbd 2>/dev/null; for p in $(pidof rmkbd); do kill "$p" 2>/dev/null; done; sleep 0.5; true' "$TARGET"
printf '  '; with_ticker 50 rm_scp_to "$REPO/rmkbd" /home/root/rmkbd.new "$TARGET"
rm_ssh 'mv -f /home/root/rmkbd.new /home/root/rmkbd && chmod +x /home/root/rmkbd' "$TARGET"
echo "  OK"
echo

fi # RM_SKIP_DEPLOY

# ---------------------------------------------------------------------------
# 3. Stop xochitl.
# ---------------------------------------------------------------------------
echo "--- Stopping xochitl ---"
rm_ssh 'systemctl stop xochitl 2>/dev/null; sleep 1; echo "  Done."' "$TARGET"
echo

# ---------------------------------------------------------------------------
# 4. Launch patched keywriter (background on device).
# ---------------------------------------------------------------------------
echo "--- Launching keywriter (backgrounded on device) ---"
KW_ENV='LD_LIBRARY_PATH=/home/root/qt5/lib'
KW_ENV="$KW_ENV QML2_IMPORT_PATH=/home/root/qt5/qml"
KW_ENV="$KW_ENV QT_PLUGIN_PATH=/home/root/qt5/plugins"
KW_ENV="$KW_ENV QMLSCENE_DEVICE=epaper"
KW_ENV="$KW_ENV QT_QPA_PLATFORM=linuxfb:fb=/dev/fb0:size=1404x1872:mmsize=158x210"
KW_ENV="$KW_ENV QT_FONT_DPI=226"
KW_ENV="$KW_ENV QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS=rotate=180"
KW_ENV="$KW_ENV QT_QPA_GENERIC_PLUGINS=evdevtablet"

rm_ssh "env $KW_ENV /home/root/keywriter >/tmp/kw-phase4.log 2>&1 &
        sleep 0.5; echo 'keywriter launched'" "$TARGET"
printf '  Qt init: '; with_ticker 5 sleep 5
echo

# ---------------------------------------------------------------------------
# 5. Launch rmkbd daemon (background on device, port 8000).
# ---------------------------------------------------------------------------
echo "--- Launching rmkbd daemon (port $PORT) ---"
rm_ssh 'pkill -f /home/root/rmkbd 2>/dev/null; for p in $(pidof rmkbd); do kill "$p" 2>/dev/null; done; true' "$TARGET"
rm_ssh "/home/root/rmkbd --port $PORT $VERBOSE >/tmp/rmkbd-phase4.log 2>&1 &
        sleep 1; echo 'rmkbd started'" "$TARGET"
printf '  WS ready: '; with_ticker 2 sleep 2
echo

BROWSER_URL="http://${TARGET}:${PORT}/"
WS_URL="ws://${TARGET}:${PORT}/ws"

echo "======================================"
echo "  >> OPEN THIS IN YOUR BROWSER <<"
echo ""
echo "     $BROWSER_URL"
echo ""
echo "  1. Click the page to focus it."
echo "  2. Type a few sentences."
echo "  3. Press Esc to toggle edit mode (first time enters edit mode)."
echo "  4. Text should appear on the e-ink screen."
echo ""
echo "  WS endpoint : $WS_URL"
echo "  Holding 3 min..."
echo "======================================"
printf '  '; with_ticker 180 sleep 180

echo
echo "--- rmkbd daemon log (/tmp/rmkbd-phase4.log) ---"
echo "  (terse by default: connections + a periodic key count; re-run with -v for per-key detail)"
rm_ssh 'cat /tmp/rmkbd-phase4.log 2>/dev/null || echo "(log not found)"' "$TARGET" \
    | sed 's/^/  /'
echo

echo "--- keywriter Qt log (last 20 lines) ---"
rm_ssh 'tail -20 /tmp/kw-phase4.log 2>/dev/null || echo "(log not found)"' "$TARGET" \
    | sed 's/^/  /'
echo

echo "--- scratch.md on device ---"
rm_ssh 'cat /home/root/edit/scratch.md 2>/dev/null || echo "(not found)"' "$TARGET" \
    | sed 's/^/  /'
echo

echo "--- Phase 4 verdict ---"
echo "  If typed text appeared on the e-ink: Phase 4 DONE-WHEN gate is GREEN."
echo "  Check scratch.md above for what keywriter saved."
echo "  Full log: $LOG"
echo "  Sync back: rmpush 'test-phase4 result'"
