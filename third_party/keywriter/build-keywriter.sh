#!/usr/bin/env bash
# third_party/keywriter/build-keywriter.sh -- Runs INSIDE the Docker container.
# Clones keywriter, builds it with qmake, then collects the Qt5 runtime
# sysroot subset (libs + qml modules + plugins) that the binary needs.
# Everything is written to /out, which the caller mounts from the host.
#
# Environment variables (all have sensible defaults):
#   KEYWRITER_REPO   git repo to clone (default: bjornte/Writerdeck-keywriter)
#   KEYWRITER_REF    branch/tag/sha to build (default: master)

set -euo pipefail

KEYWRITER_REPO="${KEYWRITER_REPO:-https://github.com/bjornte/Writerdeck-keywriter.git}"
KEYWRITER_REF="${KEYWRITER_REF:-master}"
OUT_DIR="/out"

echo "=== rM1-Writerdeck: build keywriter for reMarkable 1 ==="
echo "  Repo : ${KEYWRITER_REPO}"
echo "  Ref  : ${KEYWRITER_REF}"
echo "  SYSROOT from image: ${SYSROOT:-<not set; using image default>}"
echo

# ---------------------------------------------------------------------------
# 1. Clone keywriter + sundown markdown submodule.
# ---------------------------------------------------------------------------
echo "=== Cloning keywriter ==="
git clone --depth 1 --branch "${KEYWRITER_REF}" "${KEYWRITER_REPO}" /keywriter
cd /keywriter
git submodule update --init --recursive
echo

# ---------------------------------------------------------------------------
# 2. Build with qmake.
#    The toltec image (ghcr.io/toltec-dev/qt:v3.3) uses crosstool-ng -- there
#    is NO OE SDK environment-setup script. The cross-compiler and sysroot are
#    set as Docker ENV vars (CROSS_COMPILE, SYSROOT, PATH) at image build time.
#    Qt was configured with -device linux-arm-remarkable-g++, so that device
#    spec is the DEFAULT -- we do NOT pass -spec (it lives under
#    mkspecs/devices/, and the upstream recipe relies on the default).
#    keywriter's edit.pro hardcodes the OLD OE spec name (linux-oe-g++); the
#    upstream toltec recipe rewrites it to linux-arm-remarkable-g++ before
#    building. We replicate that exactly.
#    (ref: toltec package/keywriter/package build() function.)
# ---------------------------------------------------------------------------
echo "=== Building with qmake ==="
echo "  CROSS_COMPILE = ${CROSS_COMPILE:-<not set>}"
echo "  SYSROOT       = ${SYSROOT:-<not set>}"
echo "  qmake         = $(command -v qmake || echo '<not found>')"
echo "  default XSPEC = $(qmake -query QMAKE_XSPEC 2>/dev/null || echo '?')"
echo

# Phase 3: C++ infra lives in the Writerdeck-keywriter fork (not patched here).
# Assert so a bad KEYWRITER_REF fails loudly before qmake.
echo "=== Asserting fork C++ infra ==="
grep -q 'linux-arm-remarkable-g++' edit.pro \
    || { echo "ERROR: edit.pro missing linux-arm-remarkable-g++ (fork Phase 3)" >&2; exit 1; }
grep -q 'rotation_watcher.cpp' edit.pro \
    || { echo "ERROR: edit.pro missing rotation_watcher.cpp" >&2; exit 1; }
grep -q 'lobby_bridge.cpp' edit.pro \
    || { echo "ERROR: edit.pro missing lobby_bridge.cpp" >&2; exit 1; }
grep -q 'edit_helper.cpp' edit.pro \
    || { echo "ERROR: edit.pro missing edit_helper.cpp" >&2; exit 1; }
grep -q '\-pthread' edit.pro \
    || { echo "ERROR: edit.pro missing -pthread" >&2; exit 1; }
test -f rotation_watcher.h && test -f rotation_watcher.cpp \
    || { echo "ERROR: rotation_watcher.{h,cpp} missing from fork checkout" >&2; exit 1; }
test -f lobby_bridge.h && test -f lobby_bridge.cpp \
    || { echo "ERROR: lobby_bridge.{h,cpp} missing from fork checkout" >&2; exit 1; }
test -f edit_helper.h && test -f edit_helper.cpp \
    || { echo "ERROR: edit_helper.{h,cpp} missing from fork checkout" >&2; exit 1; }
grep -q 'clearUndoStacks' edit_helper.cpp \
    || { echo "ERROR: edit_helper.cpp missing clearUndoStacks (Phase A2 undo)" >&2; exit 1; }
grep -q 'dispatchMacArrow' edit_helper.cpp \
    || { echo "ERROR: edit_helper.cpp missing dispatchMacArrow (Phase B key dispatch)" >&2; exit 1; }
grep -q 'rmkbdSocketReader' main.cpp \
    || { echo "ERROR: main.cpp missing socket reader (fork Phase 3)" >&2; exit 1; }
grep -q 'qEnvironmentVariableIsEmpty("QT_QPA_PLATFORM")' main.cpp \
    || { echo "ERROR: main.cpp missing QT_QPA_PLATFORM guard" >&2; exit 1; }
grep -q 'qEnvironmentVariableIsEmpty("QMLSCENE_DEVICE")' main.cpp \
    || { echo "ERROR: main.cpp missing QMLSCENE_DEVICE guard" >&2; exit 1; }
echo "  fork C++ infra OK (socket + lobby_bridge + edit_helper + rotation_watcher + edit.pro)."
echo

# Fork owns QML assembly (committed main.qml from assemble-qml.sh). Assert only.
echo "=== Asserting fork Lobby/shell QML (assembled main.qml) ==="
grep -q 'property bool isLobby: true' main.qml \
    || { echo "ERROR: main.qml missing isLobby (fork Lobby/shell)" >&2; exit 1; }
grep -q 'function setLobbyInfo(' main.qml \
    || { echo "ERROR: main.qml missing setLobbyInfo" >&2; exit 1; }
grep -q 'function showLobby() {' main.qml \
    || { echo "ERROR: main.qml missing showLobby" >&2; exit 1; }
grep -q 'function handleHome(' main.qml \
    || { echo "ERROR: main.qml missing handleHome" >&2; exit 1; }
grep -q 'Writerdeck-user-documents' main.qml \
    || { echo "ERROR: main.qml missing Writerdeck-user-documents path" >&2; exit 1; }
grep -q 'handleMacKeysOnPressed' main.qml \
    || { echo "ERROR: main.qml missing handleMacKeysOnPressed" >&2; exit 1; }
grep -q 'function handleMacArrow' main.qml \
    || { echo "ERROR: main.qml missing handleMacArrow (run fork ./assemble-qml.sh before push)" >&2; exit 1; }
grep -q 'function handleMacBackspace' main.qml \
    || { echo "ERROR: main.qml missing handleMacBackspace" >&2; exit 1; }
grep -q 'editHelper.beginTextEdit' main.qml \
    || { echo "ERROR: main.qml missing editHelper.beginTextEdit" >&2; exit 1; }
grep -q 'editHelper.notifyTextChanged' main.qml \
    || { echo "ERROR: main.qml missing editHelper.notifyTextChanged" >&2; exit 1; }
grep -q 'editHelper.dispatchMacArrow' main.qml \
    || { echo "ERROR: main.qml missing editHelper.dispatchMacArrow" >&2; exit 1; }
grep -q 'id: cursorTimer' main.qml \
    || { echo "ERROR: main.qml missing cursorTimer" >&2; exit 1; }
grep -q 'id: autosaveTimer' main.qml \
    || { echo "ERROR: main.qml missing autosaveTimer" >&2; exit 1; }
grep -q 'id: sleepScreen' main.qml \
    || { echo "ERROR: main.qml missing sleepScreen (run fork ./assemble-qml.sh before push)" >&2; exit 1; }
grep -q 'id: lobbyNotesModel' main.qml \
    || { echo "ERROR: main.qml missing lobbyNotesModel" >&2; exit 1; }
test -f main.qml.in && test -f edit_mac_helpers.qml.inc && test -f assemble-qml.sh \
    || { echo "ERROR: fork main.qml.in / edit_mac_helpers.qml.inc / assemble-qml.sh missing" >&2; exit 1; }
test -d lobby && test -f concat-lobby.sh && test -f lobby/lobby_shell_top.inc \
    || { echo "ERROR: fork lobby/ or concat-lobby.sh missing" >&2; exit 1; }
# Sanity: handleKey must close before Component.onCompleted.
python3 - << 'PYEOF'
with open('main.qml') as f:
    s = f.read()
hk = s.find('    function handleKey(event) {')
co = s.find('    Component.onCompleted: {')
assert hk >= 0 and co > hk, "handleKey / Component.onCompleted anchors missing"
assert s[hk:co].count('{') == s[hk:co].count('}'), "handleKey brace mismatch -- QML will fail to load"
print('  handleKey brace balance OK.')
PYEOF
echo "  fork assembled QML OK (helpers + lobby + sleep in committed main.qml)."
echo

# No -spec: Qt's configured default XSPEC is already devices/linux-arm-remarkable-g++.
qmake edit.pro
make -j"$(nproc)"
echo

# Sanity: confirm the output binary is ARM.
echo "=== Binary check ==="
file edit
file edit | grep -q "ARM" || { echo "ERROR: binary is not ARM -- wrong toolchain?"; exit 1; }
echo

# ---------------------------------------------------------------------------
# 3. Copy the keywriter binary to /out.
# ---------------------------------------------------------------------------
mkdir -p "${OUT_DIR}"
cp edit "${OUT_DIR}/Writerdeck"
echo "  Writerdeck -> ${OUT_DIR}/Writerdeck"
echo

# ---------------------------------------------------------------------------
# 4. Collect the Qt5 runtime sysroot subset.
#    We bundle these from the toolchain (not the device's /usr/lib) so the
#    versions are guaranteed to match the binary we just built.
#    The device already has:
#      /usr/lib/qml/QtQuick + Qt/labs/folderlistmodel  (confirmed on device)
#    We still bundle the toolchain versions to stay self-contained and
#    avoid any Qt version mismatch should the device firmware change.
# ---------------------------------------------------------------------------
echo "=== Collecting Qt5 runtime sysroot subset ==="

QT5="${OUT_DIR}/qt5"
mkdir -p "${QT5}/lib" "${QT5}/qml" "${QT5}/plugins"

# The toltec Qt build used: configure -prefix /usr -sysroot "$SYSROOT" -hostprefix /usr
# So qmake -query returns paths like "/usr/lib" (prefix only, no sysroot).
# The ARM .so files actually live at ${SYSROOT}/usr/lib.
# We prepend $SYSROOT when the directory exists there; otherwise fall back to
# the bare query result (handles edge cases where qmake already includes it).
_resolve_qt_path() {
    local q
    q="$(qmake -query "$1")"
    if [ -n "${SYSROOT:-}" ] && [ -d "${SYSROOT}${q}" ]; then
        echo "${SYSROOT}${q}"
    else
        echo "$q"
    fi
}

QT_INSTALL_LIBS="$(_resolve_qt_path QT_INSTALL_LIBS)"
QT_INSTALL_QML="$(_resolve_qt_path  QT_INSTALL_QML)"
QT_INSTALL_PLUGINS="$(_resolve_qt_path QT_INSTALL_PLUGINS)"

echo "  QT_INSTALL_LIBS    = ${QT_INSTALL_LIBS}"
echo "  QT_INSTALL_QML     = ${QT_INSTALL_QML}"
echo "  QT_INSTALL_PLUGINS = ${QT_INSTALL_PLUGINS}"
echo

# 4a. Shared libraries.
#     Collect only the .so files referenced by the keywriter binary (and their
#     transitive deps within Qt). Use objdump to get the direct NEEDED list,
#     then copy all versioned Qt5 .so* files that match any of those names.
#     We include both the unversioned symlink (.so.5) and the real file (.so.5.x.y).
NEEDED="$(objdump -p "${OUT_DIR}/keywriter" 2>/dev/null \
    | grep NEEDED | awk '{print $2}' || true)"
echo "  Direct NEEDED libs: $(echo "$NEEDED" | tr '\n' ' ')"
echo

# Copy all libQt5*.so* from the toolchain lib dir that the binary needs.
# We cast a slightly wider net (all transitive Qt libs) because QML plugins
# load additional modules at runtime that objdump can't see.
for so in "${QT_INSTALL_LIBS}"/libQt5*.so*; do
    [ -f "$so" ] || [ -L "$so" ] || continue
    cp -P "$so" "${QT5}/lib/" 2>/dev/null || true
done

# Also copy non-Qt libs referenced in NEEDED that live in the Qt lib dir
# (e.g. libicui18n, libicuuc, libicudata if Qt was built with ICU).
for needed in $NEEDED; do
    for candidate in "${QT_INSTALL_LIBS}/${needed}" "${QT_INSTALL_LIBS}/${needed}.*"; do
        for f in $candidate; do
            [ -e "$f" ] && cp -P "$f" "${QT5}/lib/" 2>/dev/null || true
        done
    done
done

echo "  Copied $(ls "${QT5}/lib" | wc -l) files to qt5/lib/"

# 4b. QML import modules.
#     main.qml imports: QtQuick 2.11, QtQuick.Window 2.2,
#     Qt.labs.folderlistmodel 1.0 (+ io.singleton, compiled into the binary).
#
#     IMPORTANT: in Qt 5.15 the CORE QtQuick module lives in a directory
#     literally named "QtQuick.2" (libqtquick2plugin.so) -- this is what
#     "import QtQuick 2.x" resolves to, providing Rectangle/Text/TextEdit/
#     Flickable/Component etc. It is SEPARATE from the "QtQuick/" subfolder
#     tree (Window.2, Layouts, Controls.2, ...). Hand-listing modules missed
#     it once already (main.qml:114 "Component is not a type"), so we now copy
#     the ENTIRE qml/ tree -- a few extra MB buys certainty that no core or
#     transitive QML module (QtQuick.2, QtQml, QtQml/Models.2, ...) is missed.
if [ -d "${QT_INSTALL_QML}" ]; then
    cp -r "${QT_INSTALL_QML}/." "${QT5}/qml/"
    echo "  QML: copied entire tree from ${QT_INSTALL_QML}"
    # Sanity-check the one module whose absence bit us before.
    if [ -e "${QT5}/qml/QtQuick.2/libqtquick2plugin.so" ]; then
        echo "    OK: core QtQuick.2 plugin present."
    else
        echo "    ERROR: core QtQuick.2 plugin STILL missing after full copy!" >&2
        echo "    (import QtQuick 2.x will fail on device) -- check QT_INSTALL_QML." >&2
        exit 1
    fi
else
    echo "  ERROR: QML import dir not found: ${QT_INSTALL_QML}" >&2
    exit 1
fi

# 4c. Qt plugins.
#     platforms     -- the QPA platform plugins (linuxfb is the one we use).
#     imageformats  -- PNG/JPEG/etc decoders.
#     generic       -- input plugins; main.cpp requests evdevtablet via
#                      QT_QPA_GENERIC_PLUGINS (the "No such plugin for spec
#                      evdevtablet" warning comes from this being absent).
#                      Non-fatal for render, but cheap to include.
for plugindir in platforms imageformats generic; do
    src="${QT_INSTALL_PLUGINS}/${plugindir}"
    if [ -d "$src" ]; then
        cp -r "$src" "${QT5}/plugins/"
        echo "  Plugin dir: ${plugindir}"
    else
        echo "  (plugin dir not in toolchain, skipped: ${plugindir})"
    fi
done

echo

# ---------------------------------------------------------------------------
# 4d. Fonts.
#     Qt 5.15 ships NO fonts. The QPA basic font database looks for a "fonts"
#     directory beside the Qt libraries (on device: /home/root/qt5/lib/fonts).
#     Bundle DejaVu (guaranteed fallback) + the reading-view font set:
#     Inter (default), Literata, EB Garamond.
#     Copy each reading-view family by its *fontconfig-resolved path* rather
#     than a hardcoded directory: the Dockerfile fonts land in different trees
#     (Inter/Literata under truetype/, EB Garamond under opentype/ via apt), so
#     a fixed path silently missed EB Garamond before. fc-list is path-agnostic.
#     Qt resolves by the TTF internal family name (not filename); the fc-query
#     assertions below catch a name mismatch that would otherwise silently fall
#     back to DejaVu on device with no error.
# ---------------------------------------------------------------------------
echo "=== Bundling fonts (Qt 5.15 ships none) ==="
mkdir -p "${QT5}/lib/fonts"
fc-cache -f >/dev/null 2>&1 || true

# DejaVu -- guaranteed fallback (always present; fonts-dejavu-core in Dockerfile).
# Direct-copy the whole family (Sans/Serif/Mono) so the fallback is complete.
_font_src="/usr/share/fonts/truetype/dejavu"
if [ -d "${_font_src}" ]; then
    cp "${_font_src}"/*.ttf "${QT5}/lib/fonts/" 2>/dev/null || true
fi

# Reading-view families, located by fontconfig family name (path-independent,
# so truetype/ vs opentype/ no longer matters). cut -d: -f1 takes the file path
# (Linux paths have no ':'); the trailing fc-list column is dropped.
for _fam in "Inter" "Literata" "EB Garamond"; do
    fc-list "${_fam}" file 2>/dev/null | cut -d: -f1 | sort -u | while IFS= read -r _f; do
        [ -n "${_f}" ] && [ -f "${_f}" ] && cp "${_f}" "${QT5}/lib/fonts/" 2>/dev/null || true
    done
done

_nfonts="$(find "${QT5}/lib/fonts" -name '*.ttf' -o -name '*.otf' | wc -l | tr -d ' ')"
if [ "${_nfonts}" -ge 1 ]; then
    echo "  Bundled ${_nfonts} font file(s) into qt5/lib/fonts/:"
    ls "${QT5}/lib/fonts/" | sed 's/^/    /'
else
    echo "  ERROR: no fonts bundled -- text will NOT render on device." >&2
    exit 1
fi

# CI assertion: verify each expected Qt family name is actually present in the
# bundled files. A name typo silently falls back to DejaVu on device; fail
# loudly here instead.
echo "  Asserting font family names via fc-query..."
_assert_family() {
    local want="$1"
    local found
    found="$(fc-query --format='%{family}\n' "${QT5}/lib/fonts/"*.* 2>/dev/null \
             | tr ',' '\n' | sed 's/^ *//' | sort -u)"
    if echo "${found}" | grep -qiF "${want}"; then
        echo "    OK: '${want}' found."
    else
        echo "    ERROR: font family '${want}' NOT found in bundled fonts!" >&2
        echo "    (Qt will silently fall back to DejaVu -- check Dockerfile font install)" >&2
        echo "    Families found: $(echo "${found}" | tr '\n' ' ')" >&2
        exit 1
    fi
}
_assert_family "DejaVu Sans"
_assert_family "Inter"
_assert_family "Literata"
_assert_family "EB Garamond"
echo

# ---------------------------------------------------------------------------
# 5. Summary.
# ---------------------------------------------------------------------------
echo "=== Output summary ==="
du -sh "${OUT_DIR}/keywriter" 2>/dev/null || true
du -sh "${QT5}" 2>/dev/null || true
echo
find "${OUT_DIR}" -maxdepth 3 -name "*.so*" | wc -l | xargs printf "  Total .so files: %s\n"
echo
echo "=== Done. Deploy with: bash scripts/deploy-keywriter.sh ==="
