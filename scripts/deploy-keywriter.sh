#!/usr/bin/env bash
# scripts/deploy-keywriter.sh -- Phase 1B Step 4: deploy the source-built
# keywriter binary + Qt5 runtime sysroot to the device, launch it, and
# confirm the editor comes up.  This is the DEVICE round-trip that clears
# Unknown U4 (does the built binary + sysroot run and render?).
#
# Pre-conditions (all fulfilled by the CI build + git pull):
#   third_party/keywriter/dist/keywriter   -- ARM binary built from source
#   third_party/keywriter/dist/qt5.tar.gz  -- Qt5 runtime sysroot (compressed)
#   scripts/_env.sh, secrets/remarkable.local.env -- SSH credentials
#
# Device layout after deploy:
#   /home/root/keywriter          -- the binary (multi-GB /home, not rootfs)
#   /home/root/qt5/{lib,qml,plugins}/  -- Qt5 runtime sysroot
#   /home/root/edit/scratch.md    -- keywriter's notes dir (pre-seeded)
#
# Launch environment:
#   LD_LIBRARY_PATH=/home/root/qt5/lib
#   QML2_IMPORT_PATH=/home/root/qt5/qml
#   QT_PLUGIN_PATH=/home/root/qt5/plugins
#   QMLSCENE_DEVICE=epaper            (scene graph that drives the EPDC e-ink refresh)
#   QT_QPA_PLATFORM=linuxfb:fb=/dev/fb0:size=1404x1872:mmsize=158x210
#   QT_FONT_DPI=226                   (rM1 panel DPI; keeps point-sized fonts readable)
#   QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS=rotate=180
#   QT_QPA_GENERIC_PLUGINS=evdevtablet
#
# Verdict (launch keywriter, then classify the captured Qt stderr):
#   LIKELY-OK -- survived 8 s, no fatal Qt marker.  Human must confirm e-ink render.
#   FAILED    -- early exit or fatal marker; last 30 lines of Qt log printed.
#
# Usage (run from repo root on the Mac):
#   bash scripts/deploy-keywriter.sh
#   bash scripts/deploy-keywriter.sh 192.168.1.8   # explicit target
#   bash scripts/deploy-keywriter.sh -b                 # binary-only fast deploy (~1s); same as RM_BINARY_ONLY=1
#   RM_BINARY_ONLY=1 bash scripts/deploy-keywriter.sh   # push only the 205K binary + launch script (~1s); skip Qt5 sysroot AND the launch/verdict ceremony (terse); then: test-e2e.sh -s
#   RM_FORCE_SYSROOT=1 bash scripts/deploy-keywriter.sh # force re-push of Qt5 tarball (after Qt rebuild)

set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/_env.sh"
REPO="$(cd "$DIR/.." && pwd)"

# Args (any order): a target IP/host, or -b|--binary (binary-only fast deploy,
# = RM_BINARY_ONLY=1). Setting it here (before the knob defaults below) means
# the later RM_BINARY_ONLY="${RM_BINARY_ONLY:-0}" preserves it.
TARGET="$RM_HOST"
for arg in "$@"; do
    case "$arg" in
        -b|--binary) RM_BINARY_ONLY=1 ;;
        *)           TARGET="$arg" ;;
    esac
done

# --- Experiment knobs (override on the command line) -------------------------
# QPA platform RESOLVED (Phase 1 Step 4/5): the toltec toolchain has no epaper
# QPA *platform* plugin, so we render through **linuxfb** on the real /dev/fb0,
# pinning panel geometry + physical size (linuxfb, unlike the native epaper QPA,
# auto-detects NEITHER: without size= the window falls back to ~1024x768
# landscape upper-left; without mmsize= the DPI is guessed low and fonts come
# out ~4 px). rM1 = 1404x1872 px, ~158x210 mm => 226 DPI. The libqsgepaper
# SCENE GRAPH (QMLSCENE_DEVICE=epaper) still issues the EPDC e-ink refresh.
# Override to experiment:
#   RM_QPA=rm2fb / RM_QMLSCENE=software / RM_FONT_DPI=... bash scripts/deploy-keywriter.sh
# RM_SKIP_DEPLOY=1 skips the file transfer (binary+sysroot already on device)
# for a fast relaunch-only loop while experimenting with the platform.
# RM_BINARY_ONLY=1 pushes just the keywriter binary (~1s) and skips the ~90s
# Qt5 sysroot tarball -- use for socket-inject.patch iteration when /home/root/qt5/
# already exists on the device.  RM_FORCE_SYSROOT=1 re-pushes the full tarball
# even with RM_BINARY_ONLY=1 (use after a Qt/Dockerfile rebuild).
RM_QPA="${RM_QPA:-linuxfb:fb=/dev/fb0:size=1404x1872:mmsize=158x210}"
RM_QMLSCENE="${RM_QMLSCENE:-epaper}"
RM_FONT_DPI="${RM_FONT_DPI:-226}"
RM_SKIP_DEPLOY="${RM_SKIP_DEPLOY:-0}"
RM_BINARY_ONLY="${RM_BINARY_ONLY:-0}"
RM_FORCE_SYSROOT="${RM_FORCE_SYSROOT:-0}"

TS="$(date +%Y-%m-%dT%H-%M-%S)"
if [ "$RM_BINARY_ONLY" = "1" ]; then
    # Fast loop: no recon log file (avoids a committed docs/recon/ entry per run).
    LOG="(none -- binary-only)"
else
    RECON_DIR="$REPO/docs/recon"
    mkdir -p "$RECON_DIR"
    LOG="$RECON_DIR/deploy-keywriter-$TS.txt"
    exec > >(tee "$LOG") 2>&1
fi

echo "=== deploy-keywriter  $TS  target=$TARGET ==="
echo

# ---------------------------------------------------------------------------
# 0. Pre-flight: confirm artifacts exist from the CI build (pulled via git).
# ---------------------------------------------------------------------------
DIST="$REPO/third_party/keywriter/dist"
BINARY="$DIST/keywriter"
QT5_TGZ="$DIST/qt5.tar.gz"

echo "--- Pre-flight checks ---"
if [ ! -f "$BINARY" ]; then
    err "$BINARY not found."
    echo "  Run the CI build first (workflow: build-keywriter.yml)" >&2
    echo "  then: git pull" >&2
    exit 1
fi
if [ ! -f "$QT5_TGZ" ]; then
    err "$QT5_TGZ not found."
    echo "  Run the CI build first (workflow: build-keywriter.yml)" >&2
    echo "  then: git pull" >&2
    exit 1
fi
echo "  keywriter binary : $(ls -lh "$BINARY" | awk '{print $5, $NF}')"
echo "  Qt5 sysroot tgz  : $(ls -lh "$QT5_TGZ" | awk '{print $5, $NF}')"
echo

# Confirm device is reachable + key login works before stopping xochitl.
# Ping first: distinguishes "asleep / Wi-Fi dropped" from "missing SSH key".
echo "--- Testing SSH key login to $TARGET ---"
if ! ping -c1 -W2 "$TARGET" >/dev/null 2>&1; then
    err "$TARGET is unreachable (ping failed)."
    echo "  The reMarkable is probably asleep -- wake it and try again." >&2
    echo "  (If it's awake but on a different IP, update RM_HOST_WIFI in secrets/remarkable.local.env.)" >&2
    exit 1
fi
if ! rm_test_key "$TARGET"; then
    err "key-based SSH to root@$TARGET failed (device is reachable)."
    echo "  Run: bash scripts/bootstrap.sh  to install the key." >&2
    exit 1
fi
echo "  OK"
echo

# ---------------------------------------------------------------------------
# 1. No local unpack needed -- we push the single qt5.tar.gz to the device and
#    extract it there. That replaces a 457-file `scp -r` (one slow, hang-prone
#    round-trip per file) with ONE file transfer + a fast on-device untar.
#    QT5_LOCAL stays defined only so the EXIT trap's cleanup is set -u safe; we
#    never create it now.
# ---------------------------------------------------------------------------
QT5_LOCAL="$DIST/qt5"
if [ "$RM_BINARY_ONLY" != "1" ]; then
    echo "--- Qt5 sysroot will be pushed as a single tarball and extracted on-device ---"
    echo "  ($(du -h "$QT5_TGZ" | cut -f1) tarball; 457 files unpack on the tablet, not over Wi-Fi)"
    echo
fi

# ---------------------------------------------------------------------------
# 2. Deploy to device.
#    - Binary -> /home/root/keywriter (NOT /home/root/edit/ which is notes dir)
#    - Qt5 sysroot -> /home/root/qt5/
#    Both live on the multi-GB /home partition, not the 96%-full rootfs.
# ---------------------------------------------------------------------------
echo "--- Deploying binary -> /home/root/keywriter ---"
if [ "$RM_SKIP_DEPLOY" = "1" ]; then
    echo "  (RM_SKIP_DEPLOY=1 -- skipping transfer; using binary+sysroot already on device)"
else
# Stop any running keywriter FIRST: a live instance holds the executable busy
# (ETXTBSY -> scp "dest open Failure") -- e.g. a prior test left it running.
# pkill -x is unreliable on this device, so kill by full path + pidof.
rm_ssh 'pkill -f /home/root/keywriter 2>/dev/null; for p in $(pidof keywriter); do kill "$p" 2>/dev/null; done; sleep 0.5; true' "$TARGET"
# Stream to a temp name, then atomically mv into place (rename never hits ETXTBSY).
# rm_send_file = gzip-over-ssh stream (scp deadlocks on this link). See _env.sh.
rm_send_file "$BINARY" "/home/root/keywriter.new" "$TARGET"
rm_ssh 'mv -f /home/root/keywriter.new /home/root/keywriter && chmod +x /home/root/keywriter' "$TARGET"
echo "  OK"

# Deploy launch-keywriter.sh (the authoritative keywriter launch env, used by
# the systemd unit via --editor). Deployed alongside the binary so the script
# and binary versions stay in sync.
echo "--- Deploying launch-keywriter.sh -> /home/root/launch-keywriter.sh ---"
rm_send_file "$DIR/launch-keywriter.sh" "/home/root/launch-keywriter.sh" "$TARGET"
rm_ssh 'chmod +x /home/root/launch-keywriter.sh' "$TARGET"
echo "  OK"
echo

# Push the Qt5 sysroot tarball UNLESS RM_BINARY_ONLY=1 and the sysroot already
# exists on the device.  RM_FORCE_SYSROOT=1 always re-pushes (use after a Qt
# or Dockerfile rebuild).  Binary-only cuts the iteration loop from ~90s -> ~1s.
if [ "$RM_BINARY_ONLY" = "1" ] && [ "$RM_FORCE_SYSROOT" != "1" ]; then
    if rm_ssh '[ -d /home/root/qt5 ] && echo EXISTS || echo MISSING' "$TARGET" | grep -q EXISTS; then
        echo "--- Skipping Qt5 sysroot (RM_BINARY_ONLY=1 and /home/root/qt5 exists) ---"
        echo "  Use RM_FORCE_SYSROOT=1 to re-push after a Qt/Dockerfile rebuild."
        echo
        SKIP_SYSROOT=1
    fi
fi

if [ "${SKIP_SYSROOT:-0}" != "1" ]; then
echo "--- Deploying Qt5 sysroot -> /home/root/qt5/ (single tarball + on-device extract) ---"
# rm_send_file = gzip-over-ssh stream (scp's SFTP windowing deadlocks on this
# link); ServerAlive in RM_SSH_OPTS aborts a truly dead link in ~15s instead of
# hanging on the OS TCP retransmit (10-15+ min).
printf '    '; with_ticker 90 rm_send_file "$QT5_TGZ" "/home/root/qt5.tar.gz" "$TARGET"
echo "  tarball sent; extracting on device..."
# Clear any old sysroot, extract fresh, then drop the tarball to save space.
rm_ssh 'rm -rf /home/root/qt5 && mkdir -p /home/root/qt5 && tar -xzf /home/root/qt5.tar.gz -C /home/root/ && rm -f /home/root/qt5.tar.gz && echo extracted' "$TARGET"
echo "  OK"
fi # SKIP_SYSROOT
fi # RM_SKIP_DEPLOY
echo

# Binary-only fast loop: the binary + launch script are deployed. Skip the
# launch + verdict ceremony below -- it is a transient render-gate smoke test
# (launch, wait 8 s, kill, restore xochitl) that is redundant with
# test-e2e.sh, which launches keywriter AND rmkbd and holds them. Exiting
# here keeps the iteration loop ~1 s and the output to a few lines.
if [ "$RM_BINARY_ONLY" = "1" ]; then
    echo "Binary-only deploy done: /home/root/keywriter + launch-keywriter.sh updated."
    echo "Next: bash scripts/test-e2e.sh -s"
    exit 0
fi

# ---------------------------------------------------------------------------
# 3. Prep the notes directory.
#    keywriter hardcodes notes dir = /home/root/edit/ (main.qml).
# ---------------------------------------------------------------------------
echo "--- Preparing notes dir (/home/root/edit/) ---"
rm_ssh \
    'mkdir -p /home/root/edit; echo ok' \
    "$TARGET"
echo

# ---------------------------------------------------------------------------
# 4. Stop xochitl + launch keywriter with the full Qt5 env.
#    The trap fires on EXIT, INT, and TERM -- xochitl ALWAYS comes back.
# ---------------------------------------------------------------------------
XOCHITL_STOPPED=0

restore_xochitl() {
    if [ "$XOCHITL_STOPPED" = "1" ]; then
        echo
        echo "--- Restoring xochitl ---"
        # Kill keywriter first so it doesn't linger holding /home/root/keywriter
        # busy (ETXTBSY on the next deploy) or fight xochitl over the framebuffer.
        rm_ssh 'pkill -f /home/root/keywriter 2>/dev/null; for p in $(pidof keywriter); do kill "$p" 2>/dev/null; done; true' "$TARGET" || true
        rm_ssh 'systemctl start xochitl 2>/dev/null || true' "$TARGET" || true
        echo "  xochitl restored."
    fi
    # Clean up the locally-unpacked Qt5 tree (saved ~50 MB on the Mac).
    rm -rf "$QT5_LOCAL" 2>/dev/null || true
}
trap restore_xochitl EXIT INT TERM

echo "--- Stopping xochitl ---"
rm_ssh 'systemctl stop xochitl 2>/dev/null; killall xochitl 2>/dev/null || true' "$TARGET"
XOCHITL_STOPPED=1
sleep 1
echo "  Done."
echo

echo "--- Launching keywriter (backgrounded on device) ---"
echo "  QT_QPA_PLATFORM=$RM_QPA"
echo "  QMLSCENE_DEVICE=$RM_QMLSCENE   QT_FONT_DPI=$RM_FONT_DPI"
LAUNCH_CMD="
export LD_LIBRARY_PATH=/home/root/qt5/lib
export QML2_IMPORT_PATH=/home/root/qt5/qml
export QT_PLUGIN_PATH=/home/root/qt5/plugins
export QMLSCENE_DEVICE=$RM_QMLSCENE
export QT_QPA_PLATFORM=\"$RM_QPA\"
export QT_FONT_DPI=$RM_FONT_DPI
export QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS=rotate=180
export QT_QPA_GENERIC_PLUGINS=evdevtablet
nohup /home/root/keywriter >/tmp/kw.log 2>&1 </dev/null &
echo \$!
"
KW_PID="$(rm_ssh "$LAUNCH_CMD" "$TARGET" | tr -d '\r\n')"

if ! echo "$KW_PID" | grep -qE '^[0-9]+$'; then
    echo "ERROR: invalid or missing PID: '$KW_PID'" >&2
    echo
    echo "======================================"
    echo "  VERDICT: FAILED (did not start)"
    echo "======================================"
    exit 1
fi

echo "  PID=$KW_PID"
echo "  Waiting 8 s for Qt to initialise..."
sleep 8

# ---------------------------------------------------------------------------
# 5. Check whether the process survived Qt init.
# ---------------------------------------------------------------------------
ALIVE=0
CHECK="$(rm_ssh \
    "kill -0 $KW_PID 2>/dev/null && echo PROCESS_ALIVE || echo PROCESS_DEAD" \
    "$TARGET" | tr -d '\r')"
case "$CHECK" in *PROCESS_ALIVE*) ALIVE=1 ;; esac
echo "  Status after 8 s: $([ "$ALIVE" = "1" ] && echo alive || echo DIED)"
echo

# Done interrogating; kill the process so xochitl can reclaim the display.
rm_ssh "kill $KW_PID 2>/dev/null; sleep 1; kill -9 $KW_PID 2>/dev/null || true" \
    "$TARGET" || true

# ---------------------------------------------------------------------------
# 6. Pull Qt log and embed it in the tee'd recon log.
# ---------------------------------------------------------------------------
echo "--- Qt log (/tmp/kw.log) ---"
PULLED_LOG="/tmp/kw_deploy_pulled_$TS.log"
rm_scp_from '/tmp/kw.log' "$PULLED_LOG" "$TARGET" 2>/dev/null || true
echo "  ----------------------------------------"
if [ -f "$PULLED_LOG" ]; then
    cat "$PULLED_LOG"
else
    echo "  (log not available)"
fi
echo "  ----------------------------------------"
echo

# ---------------------------------------------------------------------------
# 7. Verdict.
# ---------------------------------------------------------------------------
VERDICT="LIKELY-OK"
FAIL_REASON=""

if [ "$ALIVE" = "0" ]; then
    VERDICT="FAILED"
    FAIL_REASON="process exited within 8 s of launch"
fi

if [ -f "$PULLED_LOG" ]; then
    FATAL="$(grep -Ei \
        'error while loading shared libraries|symbol lookup error|could not (find|load) the qt platform plugin|aborted|segmentation fault|qml module not found' \
        "$PULLED_LOG" | head -3 || true)"
    if [ -n "$FATAL" ]; then
        VERDICT="FAILED"
        FAIL_REASON="fatal marker in Qt log: $(echo "$FATAL" | head -1 | cut -c1-120)"
    fi
fi

echo "======================================"
echo "  VERDICT: $VERDICT"
if [ "$VERDICT" = "LIKELY-OK" ]; then
    echo
    echo "  keywriter (source-built) survived Qt init with no fatal errors."
    echo
    echo "  >> ACTION REQUIRED: glance at the tablet now."
    echo "     Does the e-ink show the Lobby screen (IP + PIN)?"
    echo "     If yes: the keywriter build is good -- Phase 1 is DONE."
    echo "     Bring the verdict back to Opus to proceed to Phase 2. <<<"
    echo
    echo "  Phase 1 DONE-WHEN gate:"
    echo "    [ ] Script says LIKELY-OK  (this run)"
    echo "    [ ] Human confirms e-ink render"
else
    echo "  REASON: $FAIL_REASON"
    echo
    echo "  Last ~30 lines of Qt log (evidence for Opus):"
    echo "  ----------------------------------------"
    if [ -f "$PULLED_LOG" ]; then
        tail -30 "$PULLED_LOG"
    else
        echo "  (no log)"
    fi
    echo "  ----------------------------------------"
    echo
    echo "  Possible causes:"
    echo "    - Qt shared lib missing from sysroot (check lib/ vs NEEDED)"
    echo "    - QML module version mismatch (check qml/)"
    echo "    - epaper QPA plugin not linked (was Q_IMPORT_PLUGIN compiled in?)"
    echo "    - glibc version mismatch (toolchain vs device)"
    echo "  Bring this log to Opus."
fi
echo "======================================"
echo
echo "  Full log : $LOG"
echo "  Sync back: git add docs/recon/ && git commit -m 'deploy-keywriter verdict' && git push"
echo
