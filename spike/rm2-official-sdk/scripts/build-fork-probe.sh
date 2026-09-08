#!/usr/bin/env bash
# Cross-compile Writerdeck-keywriter fork with Qt6 Codex SDK (compile probe).
set -euo pipefail
SPIKE="$(cd "$(dirname "$0")/.." && pwd)"
ROOT="$(cd "$SPIKE/../.." && pwd)"
# shellcheck source=/dev/null
. "$SPIKE/scripts/sdk-url.sh"
# shellcheck source=/dev/null
. "$SPIKE/scripts/_spike-sdk.sh"

FORK_REPO="${FORK_REPO:-https://github.com/bjornte/Writerdeck-keywriter.git}"
FORK_REF="${FORK_REF:-master}"
FORK_DIR="$SPIKE/.cache/writerdeck-keywriter"
PROBE="$SPIKE/fork_probe"
CACHE="$SPIKE/.cache"
SDK_DIR="${SPIKE_SDK_DIR:-$CACHE/sdk-$SDK_LABEL}"
INSTALLER="$CACHE/$(basename "$SDK_URL")"
if [ "${SDK_DEVICE:-rm2}" = "rm1" ]; then
  OUT="$CACHE/out-rm1"
else
  OUT="$CACHE/out"
fi
BUILD="$OUT/build-fork-probe"
mkdir -p "$CACHE" "$OUT"

need_docker() {
  if [ "${1:-}" = "--inside-container" ]; then
    return 1
  fi
  [ "$(uname -s)" = "Darwin" ] && return 0
  [ "$(uname -m)" != "x86_64" ] && return 0
  return 1
}

download_sdk() {
  if [ -f "$INSTALLER" ]; then
    echo "SDK installer already present: $INSTALLER"
    return 0
  fi
  echo "Downloading $SDK_URL ..."
  curl -fL --retry 3 -o "$INSTALLER.partial" "$SDK_URL"
  mv "$INSTALLER.partial" "$INSTALLER"
  chmod +x "$INSTALLER"
}

clone_fork() {
  # Always keep a full github checkout (submodules). Optional local tree overlays
  # edited sources for spike deploys without pushing the fork.
  if [ -d "$FORK_DIR/.git" ]; then
    echo "Updating fork checkout ..."
    git -C "$FORK_DIR" fetch --depth 1 origin "$FORK_REF"
    git -C "$FORK_DIR" checkout -f FETCH_HEAD
  else
    echo "Cloning $FORK_REPO ($FORK_REF) ..."
    git clone --depth 1 --branch "$FORK_REF" "$FORK_REPO" "$FORK_DIR"
  fi
  git -C "$FORK_DIR" submodule update --init --recursive --force
  if [ -n "${WRITERDECK_KEYWRITER_SRC:-}" ] && [ -d "$WRITERDECK_KEYWRITER_SRC" ]; then
    echo "Overlaying local fork edits from $WRITERDECK_KEYWRITER_SRC"
    for f in main.qml main.qml.in main.cpp lobby_bridge.h lobby_bridge.cpp lobby_ui_config.h lobby_ui_config.cpp assemble-qml.sh; do
      if [ -f "$WRITERDECK_KEYWRITER_SRC/$f" ]; then
        cp -f "$WRITERDECK_KEYWRITER_SRC/$f" "$FORK_DIR/$f"
      fi
    done
    if [ -d "$WRITERDECK_KEYWRITER_SRC/lobby" ]; then
      mkdir -p "$FORK_DIR/lobby"
      cp -f "$WRITERDECK_KEYWRITER_SRC"/lobby/*.inc "$FORK_DIR/lobby/" 2>/dev/null || true
    fi
  fi
}

write_product_version() {
  local ver
  ver="$(bash "$ROOT/scripts/product-version.sh" 2>/dev/null || echo "probe")"
  mkdir -p "$BUILD"
  printf '%s\n' \
    '#pragma once' \
    "/* Auto-generated for fork probe */" \
    "#define WRITERDECK_PRODUCT_VERSION \"${ver}\"" \
    > "$BUILD/product_version.h"
}

apply_probe_patch() {
  # Qt6 cmake: EditUtils needs a .cpp (header-only Q_OBJECT breaks MOC link).
  if grep -q 'Q_INVOKABLE QString markdown(QString input){' "$FORK_DIR/edit_utils.h" 2>/dev/null; then
    python3 - <<'PY' "$FORK_DIR/edit_utils.h"
import re, sys
path = sys.argv[1]
text = open(path).read()
text = re.sub(
    r'class EditUtils : public QObject\{\s*\n   Q_OBJECT\npublic:\n'
    r'    explicit EditUtils \(QObject\* parent = 0\) : QObject\(parent\) \{\}\n'
    r'    Q_INVOKABLE QString markdown\(QString input\)\{.*?\n    \}\n\};',
    'class EditUtils : public QObject {\n   Q_OBJECT\npublic:\n'
    '    explicit EditUtils(QObject *parent = nullptr);\n'
    '    Q_INVOKABLE QString markdown(QString input);\n};',
    text,
    count=1,
    flags=re.S,
)
open(path, 'w').write(text)
PY
  fi
  if [ ! -f "$FORK_DIR/edit_utils.cpp" ]; then
    cat > "$FORK_DIR/edit_utils.cpp" <<'EOF'
#include "edit_utils.h"

EditUtils::EditUtils(QObject *parent) : QObject(parent) {}

QString EditUtils::markdown(QString input)
{
    struct sd_callbacks callbacks;
    struct html_renderopt options;
    struct sd_markdown *markdown;

    struct buf *ob;
    ob = bufnew(64);
    sdhtml_renderer(&callbacks, &options, 0);
    markdown = sd_markdown_new(0, 16, &callbacks, &options);

    sd_markdown_render(ob, (const unsigned char *)input.toUtf8().constData(),
                       input.toUtf8().length(), markdown);
    sd_markdown_free(markdown);

    QString ret = QString(bufcstr(ob));
    bufrelease(ob);
    return ret;
}
EOF
  fi
  if ! grep -q 'WRITERDECK_RM2_QT6_PROBE' "$FORK_DIR/main.cpp" 2>/dev/null; then
    python3 - <<'PY' "$FORK_DIR/main.cpp"
import sys
path = sys.argv[1]
text = open(path).read()
needle = "    return app.exec();"
hook = '''#ifdef WRITERDECK_RM2_QT6_PROBE
    QTimer::singleShot(8000, &app, [&engine]() {
        const char *shotPath = "/home/root/spike-rm2-fork-probe/screen.png";
        const QObjectList roots = engine.rootObjects();
        if (roots.isEmpty())
            return;
        auto *win = qobject_cast<QQuickWindow *>(roots.first());
        if (!win)
            return;
        win->grabWindow().save(QString::fromUtf8(shotPath), "PNG");
    });
#endif

    return app.exec();'''
if needle not in text:
    sys.exit(0)
if "#include <QTimer>" not in text:
    text = text.replace("#include <QQuickItem>", "#include <QQuickItem>\n#include <QTimer>")
open(path, 'w').write(text.replace(needle, hook, 1))
PY
  fi
  if [ ! -f "$FORK_DIR/CMakeLists.txt" ]; then
    cp "$PROBE/CMakeLists.fork.txt" "$FORK_DIR/CMakeLists.txt"
  fi
  if grep -q 'Q_IMPORT_PLUGIN(QsgEpaperPlugin)' "$FORK_DIR/main.cpp" 2>/dev/null; then
    sed -i.bak '/Q_IMPORT_PLUGIN(QsgEpaperPlugin)/d' "$FORK_DIR/main.cpp"
    sed -i.bak '/#include <QtPlugin>/d' "$FORK_DIR/main.cpp"
  fi
  # Preserve launcher touch params (rM2 needs invertx). Unconditional qputenv
  # rotate=180 alone mirrors X -- Lobby top tabs fire the wrong button.
  if grep -q 'qputenv("QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS", "rotate=180");' "$FORK_DIR/main.cpp" 2>/dev/null; then
    python3 - <<'PY' "$FORK_DIR/main.cpp"
import sys
path = sys.argv[1]
text = open(path).read()
old = '''    qputenv("QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS", "rotate=180");
    qputenv("QT_QPA_GENERIC_PLUGINS", "evdevtablet");'''
new = '''    if (qEnvironmentVariableIsEmpty("QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS"))
        qputenv("QT_QPA_EVDEV_TOUCHSCREEN_PARAMETERS", "rotate=180");
    if (qEnvironmentVariableIsEmpty("QT_QPA_GENERIC_PLUGINS"))
        qputenv("QT_QPA_GENERIC_PLUGINS", "evdevtablet");'''
if old in text:
    open(path, 'w').write(text.replace(old, new, 1))
PY
  fi
  apply_rm2_ui_scale_patch
}

apply_rm2_ui_scale_patch() {
  # Qt6 epaper: same pointSize renders smaller than rM1 linuxfb. Scale fonts only;
  # lobby layout (tabBtnHeight, margins) stays at lobby-ui.json values.
  python3 - <<'PY' "$FORK_DIR/main.cpp" "$FORK_DIR/main.qml"
import re, sys
cpp, qml = sys.argv[1:3]

LOBBY_PSIZE_PROPS = (
    "labelPointSize", "badgePointSize", "titlePointSize", "sectionPointSize",
    "rowPointSize", "dialogTitlePointSize", "bannerPointSize", "helpPointSize",
    "fontPickerNamePointSize", "fontPickerSamplePointSize",
)

text = open(cpp).read()
if "writerdeckEditFontScale" not in text:
    text = text.replace(
        '    engine.rootContext()->setContextProperty("lobbyUi", &g_lobbyUi);\n',
        '    engine.rootContext()->setContextProperty("lobbyUi", &g_lobbyUi);\n'
        '    // Qt6 epaper point sizes read smaller than the old rM1 linuxfb+DPI path.\n'
        '    engine.rootContext()->setContextProperty(QStringLiteral("writerdeckEditFontScale"), 3.0);\n'
        '    engine.rootContext()->setContextProperty(QStringLiteral("writerdeckLobbyFontScale"), 2.35);\n',
    )
else:
    text = text.replace(
        '''#ifdef WRITERDECK_RM2_QT6_PROBE
    // rM1 ref tab cap ~34px at rowPointSize 14; Qt6 epaper ~3x smaller at same pt.
    engine.rootContext()->setContextProperty(QStringLiteral("writerdeckEditFontScale"), 3.0);
    engine.rootContext()->setContextProperty(QStringLiteral("writerdeckLobbyFontScale"), 1.0);
#else
    engine.rootContext()->setContextProperty(QStringLiteral("writerdeckEditFontScale"), 1.0);
    engine.rootContext()->setContextProperty(QStringLiteral("writerdeckLobbyFontScale"), 1.0);
#endif
''',
        '''    // Qt6 epaper point sizes read smaller than the old rM1 linuxfb+DPI path.
    engine.rootContext()->setContextProperty(QStringLiteral("writerdeckEditFontScale"), 3.0);
    engine.rootContext()->setContextProperty(QStringLiteral("writerdeckLobbyFontScale"), 2.35);
''',
    )
if "#include <QDir>" not in text:
    text = text.replace("#include <QTimer>\n", "#include <QTimer>\n#include <QDir>\n")
# Probe screenshot: Keyboard tab for visual compare to rM1 ref crop.
old_hook = """#ifdef WRITERDECK_RM2_QT6_PROBE
    QTimer::singleShot(8000, &app, [&engine]() {
        const char *shotPath = "/home/root/spike-rm2-fork-probe/screen.png";
        const QObjectList roots = engine.rootObjects();
        if (roots.isEmpty())
            return;
        auto *win = qobject_cast<QQuickWindow *>(roots.first());
        if (!win)
            return;
        win->grabWindow().save(QString::fromUtf8(shotPath), "PNG");
    });
#endif"""
new_hook = """#ifdef WRITERDECK_RM2_QT6_PROBE
    QTimer::singleShot(8500, &app, []() {
        if (g_rootObj && g_rootObj->property("isLobby").toBool())
            g_rootObj->setProperty("lobbyPage", 1);
    });
    QTimer::singleShot(10000, &app, [&engine]() {
        const char *shotPath = "/home/root/spike-rm2-fork-probe/screen.png";
        const QObjectList roots = engine.rootObjects();
        if (roots.isEmpty())
            return;
        auto *win = qobject_cast<QQuickWindow *>(roots.first());
        if (!win)
            return;
        QDir().mkpath(QStringLiteral("/home/root/spike-rm2-fork-probe"));
        win->grabWindow().save(QString::fromUtf8(shotPath), "PNG");
    });
#endif"""
if old_hook in text:
    text = text.replace(old_hook, new_hook)
elif "QDir().mkpath" not in text:
    text = text.replace(
        "        win->grabWindow().save(QString::fromUtf8(shotPath), \"PNG\");",
        "        QDir().mkpath(QStringLiteral(\"/home/root/spike-rm2-fork-probe\"));\n"
        "        win->grabWindow().save(QString::fromUtf8(shotPath), \"PNG\");",
    )
open(cpp, "w").write(text)

text = open(qml).read()
if "function lobbyPs(v)" not in text:
    text = text.replace(
        "    height: screen.height\n",
        "    height: screen.height\n"
        "    readonly property real editFontScale: (typeof writerdeckEditFontScale !== \"undefined\") ? writerdeckEditFontScale : 1.0\n"
        "    readonly property real lobbyFontScale: (typeof writerdeckLobbyFontScale !== \"undefined\") ? writerdeckLobbyFontScale : 1.0\n"
        "    function editPs(v) { return Math.max(1, Math.round(v * editFontScale)) }\n"
        "    function lobbyPs(v) { return Math.max(1, Math.round(v * lobbyFontScale)) }\n\n",
    )
elif "lobbyFontScale" not in text:
    text = text.replace(
        "    function editPs(v) { return Math.max(1, Math.round(v * editFontScale)) }\n",
        "    readonly property real lobbyFontScale: (typeof writerdeckLobbyFontScale !== \"undefined\") ? writerdeckLobbyFontScale : 1.0\n"
        "    function editPs(v) { return Math.max(1, Math.round(v * editFontScale)) }\n"
        "    function lobbyPs(v) { return Math.max(1, Math.round(v * lobbyFontScale)) }\n",
    )
if "harnessTextWidth > 0 ? harnessTextWidth : body.width" not in text:
    text = text.replace("width: body.width", "width: harnessTextWidth > 0 ? harnessTextWidth : body.width", 1)
text = re.sub(
    r"font\.pointSize: mode == 0 \? 12 : 10",
    "font.pointSize: editPs(mode == 0 ? 12 : 10)",
    text,
)
for prop in LOBBY_PSIZE_PROPS:
    text = text.replace(
        f"readonly property int {prop}: lobbyUi.{prop}",
        f"readonly property int {prop}: root.lobbyPs(lobbyUi.{prop})",
    )
# Hardcoded lobby literals (skip editPs / lobby. / pointSize variable refs).
text = re.sub(
    r"font\.pointSize: (\d+)\s*$",
    lambda m: f"font.pointSize: root.lobbyPs({m.group(1)})",
    text,
    flags=re.M,
)
# Loader caption sizes passed as raw ints (not lobby.* props) also need scale.
text = re.sub(
    r"property int pointSize: (\d+)\s*$",
    lambda m: f"property int pointSize: root.lobbyPs({m.group(1)})",
    text,
    flags=re.M,
)
text = text.replace("font.pointSize: root.lobbyPs(root.lobbyPs(", "font.pointSize: root.lobbyPs(")
text = text.replace("property int pointSize: root.lobbyPs(root.lobbyPs(", "property int pointSize: root.lobbyPs(")
text = re.sub(
    r"font\.pointSize: editPs\(([^)]+)\)",
    r"font.pointSize: editPs(\1)",
    text,
)
old_harness = """    function harnessSetWidth(w) {
        if (mode != 1) return
        if (harnessDefaultQueryWidth <= 0 && query.width > 0)
            harnessDefaultQueryWidth = query.width
        if (w > 0) {
            harnessTextWidth = w
            query.width = w
        } else if (harnessDefaultQueryWidth > 0) {
            harnessTextWidth = 0
            query.width = harnessDefaultQueryWidth
        }
    }"""
new_harness = """    function harnessSetWidth(w) {
        if (mode != 1) return
        if (harnessDefaultQueryWidth <= 0 && body.width > 0)
            harnessDefaultQueryWidth = body.width
        if (w > 0) {
            harnessTextWidth = w
        } else if (harnessDefaultQueryWidth > 0) {
            harnessTextWidth = 0
        }
    }"""
if old_harness in text:
    text = text.replace(old_harness, new_harness)
open(qml, "w").write(text)
PY
}

build_linux() {
  local skip_clone="${1:-}"
  if [ "$skip_clone" != "--skip-clone" ]; then
    clone_fork
  fi
  apply_probe_patch
  write_product_version
  _spike_sdk_env "$SDK_DIR" || { echo "ERROR: SDK not ready" >&2; exit 1; }
  local toolchain native_cmake jobs
  toolchain="$SDK_DIR/sysroots/x86_64-codexsdk-linux/usr/share/cmake/OEToolchainConfig.cmake"
  native_cmake="$SDK_DIR/sysroots/x86_64-codexsdk-linux/usr/bin/cmake"
  rm -rf "$BUILD"
  local cmake_src="$PROBE"
  local probe_flag=ON
  if [ "${WRITERDECK_PRODUCTION:-0}" = "1" ]; then
    cmake_src="$FORK_DIR"
    probe_flag=OFF
  fi
  "$native_cmake" -S "$cmake_src" -B "$BUILD" \
    -DCMAKE_TOOLCHAIN_FILE="$toolchain" \
    -DFORK_DIR="$FORK_DIR" \
    -DWRITERDECK_RM2_QT6_PROBE="$probe_flag" \
    -G Ninja
  jobs="$(nproc 2>/dev/null || echo 4)"
  # Qemu amd64 on Mac: parallel cc1plus can ICE. Keep the link step sequential.
  if [ "${SDK_DEVICE:-rm2}" = "rm1" ]; then
    jobs=2
  fi
  if ! "$native_cmake" --build "$BUILD" -j"$jobs" 2>&1 | tee "$CACHE/fork-probe-build.log"; then
    echo "ERROR: fork probe compile failed (see $CACHE/fork-probe-build.log)" >&2
    tail -30 "$CACHE/fork-probe-build.log" >&2
    exit 1
  fi
  if [ "${WRITERDECK_PRODUCTION:-0}" = "1" ]; then
    cp -f "$BUILD/Writerdeck" "$OUT/Writerdeck"
    file "$OUT/Writerdeck"
    echo "OK: $OUT/Writerdeck"
  else
    cp -f "$BUILD/fork_build/Writerdeck_fork_probe" "$OUT/Writerdeck_fork_probe" 2>/dev/null \
      || cp -f "$BUILD/Writerdeck_fork_probe" "$OUT/Writerdeck_fork_probe"
    file "$OUT/Writerdeck_fork_probe"
    echo "OK: $OUT/Writerdeck_fork_probe"
  fi
}

build_via_docker() {
  _spike_docker_build "$SPIKE" "$INSTALLER" "$(basename "$0")"
}

MODE="${1:-}"
if need_docker "$MODE"; then
  download_sdk
  clone_fork
  build_via_docker
else
  download_sdk
  if ! ls "$SDK_DIR"/environment-setup-* >/dev/null 2>&1; then
    mkdir -p "$SDK_DIR"
    "$INSTALLER" -d "$SDK_DIR" -y
  fi
  if [ "$MODE" = "--inside-container" ]; then
    build_linux --skip-clone
  else
    build_linux
  fi
fi
