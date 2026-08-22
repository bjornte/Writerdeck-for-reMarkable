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
OUT="$CACHE/out"
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
  if [ -d "$FORK_DIR/.git" ]; then
    echo "Updating fork checkout ..."
    git -C "$FORK_DIR" fetch --depth 1 origin "$FORK_REF"
    git -C "$FORK_DIR" checkout -f FETCH_HEAD
  else
    echo "Cloning $FORK_REPO ($FORK_REF) ..."
    git clone --depth 1 --branch "$FORK_REF" "$FORK_REPO" "$FORK_DIR"
  fi
  git -C "$FORK_DIR" submodule update --init --recursive
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
  git -C "$FORK_DIR" checkout -f -- edit_utils.h main.cpp 2>/dev/null || true
  rm -f "$FORK_DIR/edit_utils.cpp" "$FORK_DIR/main.cpp.bak" "$FORK_DIR/edit_utils.h.bak"
  # Qt6 probe: drop Qt5 epaper plugin import; use env from launcher.
  if grep -q 'Q_IMPORT_PLUGIN(QsgEpaperPlugin)' "$FORK_DIR/main.cpp"; then
    sed -i.bak '/Q_IMPORT_PLUGIN(QsgEpaperPlugin)/d' "$FORK_DIR/main.cpp"
    sed -i.bak '/#include <QtPlugin>/d' "$FORK_DIR/main.cpp"
  fi
  # Qt6 cmake: EditUtils is header-only with Q_OBJECT; needs a .cpp in the fork.
  # Probe stub until fork adds edit_utils.cpp.
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
  python3 - <<'PY' "$FORK_DIR/main.cpp"
import sys
path = sys.argv[1]
text = open(path).read()
needle = "    return app.exec();"
hook = '''    QTimer::singleShot(8000, &app, [&engine]() {
        const char *shotPath = "/home/root/spike-rm2-fork-probe/screen.png";
        const QObjectList roots = engine.rootObjects();
        if (roots.isEmpty())
            return;
        auto *win = qobject_cast<QQuickWindow *>(roots.first());
        if (!win)
            return;
        win->grabWindow().save(QString::fromUtf8(shotPath), "PNG");
    });

    return app.exec();'''
if needle not in text or "spike-rm2-fork-probe" in text:
    sys.exit(0)
if "#include <QTimer>" not in text:
    text = text.replace("#include <QQuickItem>", "#include <QQuickItem>\n#include <QTimer>")
open(path, 'w').write(text.replace(needle, hook, 1))
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
  "$native_cmake" -S "$PROBE" -B "$BUILD" \
    -DCMAKE_TOOLCHAIN_FILE="$toolchain" \
    -DFORK_DIR="$FORK_DIR" \
    -G Ninja
  jobs="$(nproc 2>/dev/null || echo 4)"
  if ! "$native_cmake" --build "$BUILD" -j"$jobs" 2>&1 | tee "$CACHE/fork-probe-build.log"; then
    echo "ERROR: fork probe compile failed (see $CACHE/fork-probe-build.log)" >&2
    tail -30 "$CACHE/fork-probe-build.log" >&2
    exit 1
  fi
  cp -f "$BUILD/Writerdeck_fork_probe" "$OUT/Writerdeck_fork_probe"
  file "$OUT/Writerdeck_fork_probe"
  echo "OK: $OUT/Writerdeck_fork_probe"
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
