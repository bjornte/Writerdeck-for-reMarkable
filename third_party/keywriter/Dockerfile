# keywriter cross-build environment.
# Image : ghcr.io/toltec-dev/qt:v3.3
#           - Qt 5.15.1 (reMarkable qtbase)
#           - glibc 2.31  (matches rM1 device)
#           - Linux headers 5.4.70  (matches device kernel)
#           - closed-source libqsgepaper plugin bundled (required by
#             Q_IMPORT_PLUGIN(QsgEpaperPlugin) in keywriter main.cpp)
#           - qmake in the build system root
#           - ARMv7 hard-float cross-compiler in $PATH
#
# Output (written to /out by build-keywriter.sh):
#   /out/keywriter          -- the ARM binary
#   /out/qt5/lib/           -- Qt5 shared libs needed at runtime
#   /out/qt5/lib/fonts/     -- DejaVu TTFs (Qt 5.15 ships none; needed for text)
#   /out/qt5/qml/           -- QML import modules (QtQuick, folderlistmodel, ...)
#   /out/qt5/plugins/       -- platform + imageformat plugins
#
# Usage: see third_party/keywriter/build.sh (Linux/macOS) or
#        .github/workflows/build-keywriter.yml (CI).
# The caller mounts a host directory to /out.
FROM ghcr.io/toltec-dev/qt:v3.3

ARG DEBIAN_FRONTEND=noninteractive

# The toltec image already has the ARM cross-compiler and qmake.
# We need git (to clone keywriter + sundown submodule), a host C++ compiler
# for any native build steps qmake may invoke, and fonts-dejavu-core -- Qt 5.15
# ships NO fonts, so build-keywriter.sh bundles DejaVu TTFs into the sysroot
# (without them the editor paints chrome but renders no text on device).
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    ca-certificates \
    build-essential \
    python3 \
    file \
    fonts-dejavu-core \
    fonts-ebgaramond \
    curl \
    unzip \
    fontconfig \
 && rm -rf /var/lib/apt/lists/*

# Reading-view fonts (OFL/free licences). Qt resolves font families by the
# TTF's *internal* family name, not the filename, so each block ends with an
# fc-list assertion: if the family is missing the BUILD fails here (with the
# full curl/unzip log visible) instead of 20 min later in build-keywriter.sh.
# Each font is its own RUN with `set -eux` -- no `&&...|| true` chain, which
# previously let a failed download/extract cascade to `|| true` and exit 0,
# masking the failure (the bug that shipped only DejaVu).
#
#   Inter -- https://rsms.me/inter (SIL OFL 1.1). The v4.0 zip keeps static
#   TTFs under "Inter Hinted for Windows/"; "Inter Desktop/" holds OTF. Match
#   by *name* (path-agnostic), with the variable font as a last-resort fallback.
RUN set -eux; \
    mkdir -p /usr/share/fonts/truetype/inter; \
    curl -fsSL https://github.com/rsms/inter/releases/download/v4.0/Inter-4.0.zip -o /tmp/inter.zip; \
    unzip -jo /tmp/inter.zip '*Inter-Regular.ttf' '*Inter-Bold.ttf' '*Inter-Italic.ttf' '*Inter-BoldItalic.ttf' -d /usr/share/fonts/truetype/inter/ || true; \
    if [ -z "$(ls -A /usr/share/fonts/truetype/inter/ 2>/dev/null)" ]; then \
        echo "static Inter TTFs not matched -- falling back to InterVariable.ttf"; \
        unzip -jo /tmp/inter.zip '*InterVariable.ttf' -d /usr/share/fonts/truetype/inter/; \
    fi; \
    rm -f /tmp/inter.zip; \
    fc-cache -f; \
    fc-list "Inter" | grep -qi "inter" || { echo "ERROR: Inter family missing after install" >&2; exit 1; }

#   Literata -- https://github.com/googlefonts/literata (SIL OFL 1.1). Latest
#   release is 3.103 (asset 3.103.zip); there is no v7.000 tag. Extract the
#   whole zip and pull the static TTFs by name, falling back to any Literata*.ttf.
RUN set -eux; \
    mkdir -p /usr/share/fonts/truetype/literata /tmp/lit; \
    curl -fsSL https://github.com/googlefonts/literata/releases/download/3.103/3.103.zip -o /tmp/lit.zip; \
    unzip -qo /tmp/lit.zip -d /tmp/lit; \
    find /tmp/lit \( -iname 'Literata-Regular.ttf' -o -iname 'Literata-Bold.ttf' -o -iname 'Literata-Italic.ttf' -o -iname 'Literata-BoldItalic.ttf' \) -exec cp {} /usr/share/fonts/truetype/literata/ \; ; \
    if [ -z "$(ls -A /usr/share/fonts/truetype/literata/ 2>/dev/null)" ]; then \
        echo "named static Literata not found -- copying any Literata*.ttf"; \
        find /tmp/lit -iname 'Literata*.ttf' -exec cp {} /usr/share/fonts/truetype/literata/ \; ; \
    fi; \
    rm -rf /tmp/lit /tmp/lit.zip; \
    fc-cache -f; \
    fc-list "Literata" | grep -qi "literata" || { echo "ERROR: Literata family missing after install" >&2; exit 1; }

#   EB Garamond -- apt fonts-ebgaramond (installed above; lands in opentype/).
#   Assert here so a missing/renamed package fails fast too.
RUN set -eux; \
    fc-cache -f; \
    fc-list "EB Garamond" | grep -qi "garamond" || { echo "ERROR: EB Garamond missing (fonts-ebgaramond apt?)" >&2; exit 1; }

COPY socket-inject.patch /socket-inject.patch
COPY rotation_watcher.h /rotation_watcher.h
COPY rotation_watcher.cpp /rotation_watcher.cpp
COPY lobby_bridge.h /lobby_bridge.h
COPY lobby_bridge.cpp /lobby_bridge.cpp
COPY lobby_subpages.qml.inc /lobby_subpages.qml.inc
COPY build-keywriter.sh /build-keywriter.sh
RUN chmod +x /build-keywriter.sh

ENTRYPOINT ["/build-keywriter.sh"]
