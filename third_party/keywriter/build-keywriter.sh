#!/usr/bin/env bash
# third_party/keywriter/build-keywriter.sh -- Runs INSIDE the Docker container.
# Clones keywriter, builds it with qmake, then collects the Qt5 runtime
# sysroot subset (libs + qml modules + plugins) that the binary needs.
# Everything is written to /out, which the caller mounts from the host.
#
# Environment variables (all have sensible defaults):
#   KEYWRITER_REPO   git repo to clone (default: dps/remarkable-keywriter)
#   KEYWRITER_REF    branch/tag/sha to build (default: master)

set -euo pipefail

KEYWRITER_REPO="${KEYWRITER_REPO:-https://github.com/dps/remarkable-keywriter.git}"
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

# Rewrite the hardcoded OE spec inside edit.pro to the toltec device spec.
echo "  edit.pro spec line(s) before sed:"
grep -n 'linux-oe-g++' edit.pro || echo "    (no linux-oe-g++ found -- pro file may already be patched)"
sed -i 's/linux-oe-g++/linux-arm-remarkable-g++/' edit.pro

# Patch main.cpp: guard the two display qputenv() calls so an exported env var
# wins, while the stock 'epaper' default is preserved when nothing is set.
# This is required because:
#  (a) the toltec toolchain has no 'epaper' QPA platform plugin -- only linuxfb
#      (and later rm2fb) are available, so we need QT_QPA_PLATFORM to be
#      overridable from the outside.
#  (b) the same file will be patched again in Phase 2 (socket injection),
#      so doing it here is consistent and cheap.
# qEnvironmentVariableIsEmpty() is in <QtGlobal>, already included via
# <QGuiApplication>.
echo "  Patching main.cpp: guarding qputenv with qEnvironmentVariableIsEmpty ..."
grep -n 'qputenv.*QT_QPA_PLATFORM\|qputenv.*QMLSCENE_DEVICE' main.cpp \
    || echo "    (no matching qputenv lines found -- may already be patched)"
sed -i \
    's|qputenv("QT_QPA_PLATFORM",|if (qEnvironmentVariableIsEmpty("QT_QPA_PLATFORM")) qputenv("QT_QPA_PLATFORM",|' \
    main.cpp
sed -i \
    's|qputenv("QMLSCENE_DEVICE",|if (qEnvironmentVariableIsEmpty("QMLSCENE_DEVICE")) qputenv("QMLSCENE_DEVICE",|' \
    main.cpp
echo "  main.cpp after patch:"
grep -n 'QT_QPA_PLATFORM\|QMLSCENE_DEVICE' main.cpp || true
echo

# Apply the Phase 2 socket-injection patch (socket reader thread).
echo "=== Applying socket-injection patch ==="
# --recount: the patch is hand-authored, so trust the diff body, not the @@
#   line counts (git apply is strict and does no fuzz; --recount infers counts).
# --ignore-whitespace: tolerate any whitespace/CRLF drift in upstream main.cpp.
git apply --recount --ignore-whitespace /socket-inject.patch
echo "  Patch applied."

# Rotation persistence: moc'd helper for rotationChanged -> server notify.
cp /rotation_watcher.h /rotation_watcher.cpp /keywriter/
cp /lobby_bridge.h /lobby_bridge.cpp /keywriter/
printf '\nHEADERS += rotation_watcher.h\nSOURCES += rotation_watcher.cpp\n' >> edit.pro
printf '\nHEADERS += lobby_bridge.h\nSOURCES += lobby_bridge.cpp\n' >> edit.pro
echo "  rotation_watcher + lobby_bridge added to edit.pro."

# Add -pthread to edit.pro for std::thread (socket reader thread).
printf '\nQMAKE_CXXFLAGS += -pthread\nQMAKE_LFLAGS += -pthread\n' >> edit.pro
echo "  -pthread added to edit.pro."
echo

# Apply Phase 4+7 QML edits via Python3 (not git apply / sed: Python does exact
# string replacement -- no line-number dependency, robust to upstream whitespace).
# Changes: boot in edit mode, saveAndQuit() for home-button exit, Ctrl-K/Q fix.
echo "=== Patching main.qml (Python3) ==="
python3 - << 'PYEOF'
import re, sys

with open('main.qml') as f:
    s = f.read()

# 0. Notes directory: upstream uses $HOME/edit/; Writerdeck uses Writerdeck-user-documents/.
assert 'file://%1/edit/' in s, "notes folder path not found in main.qml"
s = s.replace('file://%1/edit/', 'file://%1/Writerdeck-user-documents/', 1)

# 1. Boot in edit mode (unique string; safe global replace).
assert 'property int mode: 0' in s, "mode property not found in main.qml"
s = s.replace('property int mode: 0', 'property int mode: 1', 1)

# 1b. Plain-markdown helpers: detect Qt qrichtext/HTML and sanitize corrupted loads.
old1b = '    function toggleMode() {'
new1b = (
    '    function isHtmlPayload(t) {\n'
    '        if (!t || t.length < 9) return false\n'
    '        var head = t.substring(0, Math.min(256, t.length)).toLowerCase()\n'
    '        return head.indexOf("<!doctype html") === 0\n'
    '            || head.indexOf("<html") === 0\n'
    '            || t.indexOf(\'name="qrichtext"\') >= 0\n'
    '    }\n'
    '\n'
    '    function sanitizeLoadedNote(t) {\n'
    '        if (!isHtmlPayload(t)) return t\n'
    '        console.log("sanitize: stripping HTML wrapper from loaded note")\n'
    '        var plain = t\n'
    '        plain = plain.replace(/<br\\s*\\/?>/gi, "\\n")\n'
    '        plain = plain.replace(/<\\/p>/gi, "\\n")\n'
    '        plain = plain.replace(/<\\/div>/gi, "\\n")\n'
    '        plain = plain.replace(/<\\/li>/gi, "\\n")\n'
    '        plain = plain.replace(/<[^>]+>/g, "")\n'
    '        plain = plain.replace(/&nbsp;/g, " ")\n'
    '        plain = plain.replace(/&amp;/g, "&")\n'
    '        plain = plain.replace(/&lt;/g, "<")\n'
    '        plain = plain.replace(/&gt;/g, ">")\n'
    '        plain = plain.replace(/&quot;/g, \'"\')\n'
    '        plain = plain.replace(/&#(\\d+);/g, function(_, n) { return String.fromCharCode(parseInt(n, 10)) })\n'
    '        plain = plain.replace(/\\n{3,}/g, "\\n\\n")\n'
    '        return plain.replace(/^\\s+|\\s+$/g, "")\n'
    '    }\n'
    '\n'
    '    function toggleMode() {'
)
assert old1b in s, "toggleMode not found (edit 1b)"
s = s.replace(old1b, new1b, 1)

# 1c. toggleMode: save only on edit->preview; never doc=query.text across RichText.
old1c = (
    '    function toggleMode() {\n'
    '        if (mode == 0) {\n'
    '            mode = 1\n'
    '            query.cursorPosition = lastCursorPostion == -1 ? query.length : lastCursorPostion\n'
    '        } else {\n'
    '            doc = query.text\n'
    '            lastCursorPostion = query.cursorPosition\n'
    '            mode = 0\n'
    '        }\n'
    '        saveFile()\n'
    '    }'
)
new1c = (
    '    function toggleMode() {\n'
    '        if (mode == 0) {\n'
    '            mode = 1\n'
    '            syncQueryDisplay()\n'
    '            query.cursorPosition = lastCursorPostion == -1 ? query.length : lastCursorPostion\n'
    '            if (currentFile !== "") writerdeck.notifyOpen(currentFile)\n'
    '        } else {\n'
    '            doc = query.text\n'
    '            lastCursorPostion = query.cursorPosition\n'
    '            mode = 0\n'
    '            syncQueryDisplay()\n'
    '            saveFile()\n'
    '            if (currentFile !== "") writerdeck.notifyReadOpen(currentFile)\n'
    '        }\n'
    '    }'
)
assert old1c in s, "toggleMode body not found (edit 1c)"
s = s.replace(old1c, new1c, 1)

# 2. doLoad resets mode on every load -- keep it edit (mode=1).
#    Match the first 'mode = 0' on the line immediately after 'isOmni = false';
#    leaves toggleMode's own 'mode = 0' (preview+save) untouched.
s, n = re.subn(
    r'(isOmni\s*=\s*false[^\n]*\n[^\n]*)mode\s*=\s*0',
    lambda m: m.group(0).replace('mode = 0', 'mode = 1', 1),
    s, count=1
)
assert n == 1, "doLoad mode=0 pattern not found in main.qml"

# 2a. doLoad: sanitize Qt qrichtext/HTML accidentally saved as .md.
old2a = '                var response = xhr.responseText\n'
new2a = '                var response = sanitizeLoadedNote(xhr.responseText)\n'
assert old2a in s, "doLoad response line not found (edit 2a)"
s = s.replace(old2a, new2a, 1)

# 2b. doLoad must push loaded bytes into query.text. handleHome/showLobby clear
#     query.text for the Lobby overlay, which breaks the `text: doc` binding; the
#     next Home save then copies empty query.text -> doc and wipes the file.
old2b = '                currentFile = name\n                doc = response\n'
new2b = (
    '                currentFile = name\n'
    '                doc = response\n'
    '                autosaveSnapshot = response\n'
    '                if (lobbyOpenInReadMode) {\n'
    '                    mode = 0\n'
    '                    lobbyOpenInReadMode = false\n'
    '                    syncQueryDisplay()\n'
    '                    writerdeck.notifyReadOpen(name)\n'
    '                } else {\n'
    '                    syncQueryDisplay()\n'
    '                    writerdeck.notifyOpen(name)\n'
    '                }\n'
)
assert old2b in s, "doLoad doc=response block not found (edit 2b)"
s = s.replace(old2b, new2b, 1)

# 2c. doLoad must use loopback HTTP (decrypt .md.enc server-side), not file://.
old2c_open = '        xhr.open("GET", folder + name)'
new2c_open = '        xhr.open("GET", "http://127.0.0.1:8000/api/notes/" + encodeURIComponent(name))'
assert old2c_open in s, "doLoad xhr.open(folder) not found (edit 2c)"
s = s.replace(old2c_open, new2c_open, 1)

old2c_stat = (
    '            if (xhr.readyState === XMLHttpRequest.DONE) {\n'
    '                var response = sanitizeLoadedNote(xhr.responseText)'
)
new2c_stat = (
    '            if (xhr.readyState === XMLHttpRequest.DONE) {\n'
    '                if (xhr.status === 423) {\n'
    '                    vaultPendingLoad = name\n'
    '                    vaultBeginPIN("Enter PIN to open this note", true)\n'
    '                    return\n'
    '                }\n'
    '                if (xhr.status !== 200) {\n'
    '                    var errMsg = "Could not open note"\n'
    '                    if (xhr.status === 500 && name.indexOf(".md.enc") >= 0)\n'
    '                        errMsg = "Cannot decrypt: wrong vault key or corrupted file"\n'
    '                    vaultOpFailed(errMsg)\n'
    '                    return\n'
    '                }\n'
    '                var response = sanitizeLoadedNote(xhr.responseText)'
)
assert old2c_stat in s, "doLoad DONE handler not found (edit 2c)"
s = s.replace(old2c_stat, new2c_stat, 1)

# 3. Add saveAndQuit() and saveAndLoad(name) before initFile.
#    saveAndQuit(): syncs query.text->doc so a home-button exit does not
#    lose the live edit buffer (saveFile() writes doc, not query.text).
#    saveAndLoad(name): same sync+save, then loads a new note -- used by the
#    open-note command (slice 8d) so the current note is always saved first.
old3 = '    function initFile(name) {'
new3 = (
    '    function saveAndLoad(name) {\n'
    '        isLobby = false\n'
    '        if (mode == 1) doc = query.text\n'
    '        saveFile()\n'
    '        doLoad(name)\n'
    '    }\n'
    '\n'
    '    function saveAndQuit() {\n'
    '        if (mode == 1) doc = query.text\n'
    '        saveFile()\n'
    '        Qt.quit()\n'
    '    }\n'
    '\n'
    '    function initFile(name) {'
)
assert old3 in s, "function initFile not found in main.qml"
s = s.replace(old3, new3, 1)

# 3b. openNotePicker(): show omni over the Lobby (keep isLobby true; z-order in 3c).
old3b = '    function saveAndLoad(name) {'
new3b = (
    '    function openNotePicker() {\n'
    '        isOmni = true\n'
    '        omniQuery = ""\n'
    '    }\n'
    '\n'
    '    function saveAndLoad(name) {'
)
assert old3b in s, "function saveAndLoad not found (edit 3b)"
s = s.replace(old3b, new3b, 1)

# 3d. saveAndLoad must match Lobby Files Edit: encrypted note switch -> PIN overlay.
old3d = (
    '    function saveAndLoad(name) {\n'
    '        isLobby = false\n'
    '        if (mode == 1) doc = query.text\n'
    '        saveFile()\n'
    '        doLoad(name)\n'
    '    }\n'
)
new3d = (
    '    function saveAndLoad(name) {\n'
    '        if (name && name.indexOf(".md.enc") >= 0 && currentFile !== name) {\n'
    '            vaultPendingLoad = name\n'
    '            vaultBeginPIN("Enter PIN to edit encrypted note", true)\n'
    '            return\n'
    '        }\n'
    '        var wasLobby = isLobby\n'
    '        isLobby = false\n'
    '        if (!wasLobby) {\n'
    '            if (mode == 1) doc = query.text\n'
    '            if (currentFile !== "") saveFile()\n'
    '        } else if (currentFile !== name) {\n'
    '            currentFile = ""\n'
    '        }\n'
    '        doLoad(name)\n'
    '    }\n'
)
assert old3d in s, "saveAndLoad body not found (edit 3d)"
s = s.replace(old3d, new3d, 1)

# 3c. Omni z-order above Lobby so Open/Ctrl-K from Lobby never exposes stale doc.
old3c = '        Rectangle {\n            id: quick\n'
new3c = '        Rectangle {\n            id: quick\n            z: isOmni ? 10 : 0\n'
assert old3c in s, "quick Rectangle not found (edit 3c)"
s = s.replace(old3c, new3c, 1)

# 4. Ctrl-K note-switcher: also accept event.modifiers.
#    Our injector sets the modifier FLAG on the K event but never sends a
#    standalone Key_Control press, so ctrlPressed is always false for injected
#    keys. Accepting event.modifiers fixes Ctrl-K without breaking real keyboards.
old4k = 'event.key === Qt.Key_K && ctrlPressed'
new4k = 'event.key === Qt.Key_K && (ctrlPressed || (event.modifiers & Qt.ControlModifier))'
assert old4k in s, "Ctrl-K handler not found in main.qml"
s = s.replace(old4k, new4k)

# 5. Ctrl-Q quit: modifier fix + save before exit (saveAndQuit, not bare Qt.quit).
old4q = (
    '        } else if (event.key === Qt.Key_Q && ctrlPressed) {\n'
    '            Qt.quit()\n'
    '        }'
)
new4q = (
    '        } else if (event.key === Qt.Key_Q && (ctrlPressed || (event.modifiers & Qt.ControlModifier))) {\n'
    '            saveAndQuit()\n'
    '        }'
)
assert old4q in s, "Ctrl-Q handler not found in main.qml"
s = s.replace(old4q, new4q, 1)

# 4b. Ctrl-K from Lobby opens the note picker; in editor it toggles omni.
old4b = (
    '        } else if (event.key === Qt.Key_K && (ctrlPressed || (event.modifiers & Qt.ControlModifier))) {\n'
    '            isOmni = !isOmni\n'
    '            event.accepted = true\n'
    '        }'
)
new4b = (
    '        } else if (event.key === Qt.Key_K && (ctrlPressed || (event.modifiers & Qt.ControlModifier))) {\n'
    '            if (isLobby) {\n'
    '                if (isOmni) isOmni = false\n'
    '                else { lobbyGoPage(0); openNotePicker() }\n'
    '            } else isOmni = !isOmni\n'
    '            event.accepted = true\n'
    '        }'
)
assert old4b in s, "Ctrl-K handler block not found (edit 4b)"
s = s.replace(old4b, new4b, 1)

# 4c. Ctrl-R in Lobby rotates the display 90° clockwise (preview uses Ctrl+Right).
old4c = (
    '            event.accepted = true\n'
    '        } else if (event.key === Qt.Key_Q && (ctrlPressed || (event.modifiers & Qt.ControlModifier))) {\n'
)
new4c = (
    '            event.accepted = true\n'
    '        } else if (event.key === Qt.Key_R && (ctrlPressed || (event.modifiers & Qt.ControlModifier))) {\n'
    '            if (isLobby) {\n'
    '                rotateScreen()\n'
    '                event.accepted = true\n'
    '            }\n'
    '        } else if (event.key === Qt.Key_Q && (ctrlPressed || (event.modifiers & Qt.ControlModifier))) {\n'
)
assert old4c in s, "handleKeyDown Ctrl-Q anchor not found (edit 4c)"
s = s.replace(old4c, new4c, 1)

# 5b. Ctrl-K note-switcher data-loss fix. The omni (note-switcher) Enter handler
#     saves the current note before switching, but it calls a bare saveFile(),
#     which writes `doc` -- and in edit mode (mode==1) the live text is in
#     query.text, with `doc` synced only by toggleMode(). So switching notes via
#     Ctrl-K wrote the STALE doc and lost everything typed since the last mode
#     toggle (the same class of bug as the home-button saveAndQuit fix). Sync
#     query.text -> doc before the save. Anchored on the unique
#     "Key_Return) {  saveFile()" pair so only the omni path is hit; [ \t] (not
#     \s) keeps the match on its own line. mode==0 (preview) already has a current
#     doc, so the mode==1 guard makes the sync a no-op there.
s, nk = re.subn(
    r'(\|\| event\.key === Qt\.Key_Return\) \{[ \t]*\n)([ \t]*)saveFile\(\)',
    r'\1\2if (mode == 1) doc = query.text\n\2saveFile()',
    s, count=1
)
assert nk == 1, "omni Enter saveFile pattern not found in main.qml (Ctrl-K data-loss fix)"

# 5c. Omni pick from Lobby: hide Lobby once a note is chosen.
s, n5c = re.subn(
    r'(\s+doLoad\(omniList\.currentItem\.text\)\n\s+\}\n)(\s+)isOmni = false',
    r'\1\2isLobby = false\n\2isOmni = false',
    s, count=1
)
assert n5c == 1, "omni Enter close block not found (edit 5c)"

# 6. Add isLobby / lobbyIP / lobbyPIN properties after isOmni.
#    isLobby starts true so the device shows the Lobby on boot; saveAndLoad()
#    (edit 3) sets it false when a note is opened from the phone.
old6 = '    property bool isOmni: false'
new6 = (
    '    property bool isOmni: false\n'
    '    property bool isLobby: true\n'
    '    property bool isSleeping: false\n'
    '    property string lobbyIP: ""\n'
    '    property string lobbyPIN: ""\n'
    '    property bool lobbySyncOn: false\n'
    '    property string lobbySyncRepo: ""\n'
    '    property int lobbyNoteCount: 0\n'
    '    property string lobbyLastSync: ""\n'
    '    property bool lobbySyncReady: false\n'
    '    property bool lobbySyncing: false\n'
    '    property string lobbySyncError: ""\n'
    '    property bool lobbyWifi: true\n'
    '    property string lobbyKeyboardLayout: "us"\n'
    '    property string lobbyPinDigits: "6"\n'
    '    property string lobbySettingsMode: ""\n'
    '    property int lobbyPage: 0\n'
    '    property var lobbyTabLabels: ["Files", "Keyboard", "Sync", "Settings", "Shortcuts", "Home"]\n'
    '    property int lobbyFilesIndex: 0\n'
    '    property string lobbyLastEditedFile: ""\n'
    '    property string lobbyFilesMode: ""\n'
    '    onLobbyFilesModeChanged: writerdeck.notifyLobbyInput(lobbyFilesMode)\n'
    '    property string lobbyFilesInput: ""\n'
    '    property int lobbyFilesInputPos: 0\n'
    '    property bool lobbyOpenInReadMode: false\n'
    '    property bool suppressNextHomeKey: false\n'
    '    property bool lobbyEncryptionEnabled: false\n'
    '    property string lobbyVaultError: ""\n'
    '    property string vaultOverlayMode: ""\n'
    '    property string vaultOverlayReason: ""\n'
    '    property string vaultPinInput: ""\n'
    '    property string vaultPinPending: ""\n'
    '    property bool vaultPinKeepSession: false\n'
    '    property string vaultPendingLoad: ""\n'
    '    property string vaultPendingAction: ""\n'
    '    property string vaultPendingNote: ""'
)
assert old6 in s, "isOmni property not found in main.qml"
s = s.replace(old6, new6, 1)

# 6b. In the Lobby, lobbyFocus owns keyboard input; query must not compete.
old6b = 'focus: !isOmni'
new6b = 'focus: !isOmni && !isLobby'
assert old6b in s, "query focus binding not found (edit 6b)"
s = s.replace(old6b, new6b, 1)

# 7. Add setLobbyInfo() and handleHome() before initFile.
#    setLobbyInfo: called from C++ via invokeMethod when rmkbd connects, so the
#    Lobby shows the current IP and PIN without hardcoding anything.
#    handleHome: two-level Home (slice 8e):
#      - editing -> save current note + return to Lobby (isLobby = true)
#      - Lobby   -> Qt.quit() -> cmd.Wait fires -> s.end() -> xochitl restarts
#    Physical Home delivers twice: gpio-keys -> daemon cmd (fromPhysicalCmd=true)
#    and Qt Key_Home. suppressNextHomeKey pairs off the duplicate; USB Home uses
#    the Key_Home path only (fromPhysicalCmd=false).
old7 = '    function initFile(name) {'
new7 = (
    '    function setLobbyInfo(ip, pin, syncOn, syncRepo, noteCount, lastSync, syncReady, syncing, keyboardLayout, pinDigits) {\n'
    '        lobbyIP = ip\n'
    '        lobbyPIN = pin\n'
    '        lobbySyncOn = !!syncOn\n'
    '        lobbySyncRepo = syncRepo || ""\n'
    '        lobbyNoteCount = noteCount || 0\n'
    '        lobbyLastSync = lastSync || ""\n'
    '        lobbySyncReady = !!syncReady\n'
    '        lobbySyncing = !!syncing\n'
    '        lobbyKeyboardLayout = keyboardLayout || "us"\n'
    '        lobbyPinDigits = pinDigits || "6"\n'
    '    }\n'
    '\n'
    '    function setEncryptionEnabled(enabled) {\n'
    '        lobbyEncryptionEnabled = !!enabled\n'
    '    }\n'
    '\n'
    '    function vaultOpFailed(msg) {\n'
    '        lobbyGoPage(0)\n'
    '        lobbyVaultError = msg || "Operation failed"\n'
    '    }\n'
    '\n'
    '    function vaultOnPINAccepted() {\n'
    '        lobbyVaultError = ""\n'
    '        if (vaultPendingAction === "encrypt") {\n'
    '            var encName = vaultPendingNote\n'
    '            vaultPendingAction = ""\n'
    '            vaultPendingNote = ""\n'
    '            if (encName) {\n'
    '                selectNoteByName(encName)\n'
    '                writerdeck.encryptNote(encName)\n'
    '            }\n'
    '            return\n'
    '        }\n'
    '        if (vaultPendingAction === "decrypt") {\n'
    '            var decName = vaultPendingNote\n'
    '            vaultPendingAction = ""\n'
    '            vaultPendingNote = ""\n'
    '            if (decName) {\n'
    '                selectNoteByName(decName)\n'
    '                writerdeck.decryptNote(decName)\n'
    '            }\n'
    '            return\n'
    '        }\n'
    '        if (vaultPendingAction === "new-encrypted") {\n'
    '            vaultPendingAction = ""\n'
    '            lobbyFilesMode = "new-encrypted"\n'
    '            lobbyFilesInput = ""\n'
    '            lobbyFilesInputPos = 0\n'
    '            return\n'
    '        }\n'
    '        if (vaultPendingLoad !== "") {\n'
    '            var pending = vaultPendingLoad\n'
    '            var readMode = lobbyOpenInReadMode\n'
    '            vaultPendingLoad = ""\n'
    '            if (readMode) lobbyOpenInReadMode = false\n'
    '            Qt.callLater(function() {\n'
    '                if (isLobby) {\n'
    '                    if (readMode) {\n'
    '                        isLobby = false\n'
    '                        if (mode == 1) doc = query.text\n'
    '                        saveFile()\n'
    '                        lobbyOpenInReadMode = true\n'
    '                        doLoad(pending)\n'
    '                    } else {\n'
    '                        isLobby = false\n'
    '                        currentFile = ""\n'
    '                        doLoad(pending)\n'
    '                    }\n'
    '                } else {\n'
    '                    doLoad(pending)\n'
    '                }\n'
    '            })\n'
    '        }\n'
    '    }\n'
    '\n'
    '    function setLobbySyncStatus(syncError, wifi) {\n'
    '        lobbySyncError = syncError || ""\n'
    '        lobbyWifi = !!wifi\n'
    '    }\n'
    '\n'
    '    function lobbyGoPage(idx) {\n'
    '        if (idx < 0 || idx >= lobbyTabLabels.length) return\n'
    '        lobbyPage = idx\n'
    '        lobbyFilesMode = ""\n'
    '        lobbyFilesInput = ""\n'
    '        lobbyFilesInputPos = 0\n'
    '        lobbySettingsMode = ""\n'
    '        if (idx === 0) lobbyRefreshNotes()\n'
    '    }\n'
    '\n'
    '    function lobbyRefreshNotes() {\n'
    '        writerdeck.requestNotesList()\n'
    '    }\n'
    '\n'
    '    function setNotesList(items) {\n'
    '        lobbyNotesModel.clear()\n'
    '        if (!items) return\n'
    '        for (var i = 0; i < items.length; i++) {\n'
    '            var it = items[i]\n'
    '            lobbyNotesModel.append({\n'
    '                name: it.name !== undefined ? it.name : "",\n'
    '                size: it.size !== undefined ? it.size : 0,\n'
    '                modified: it.modified !== undefined ? it.modified : "",\n'
    '                encrypted: !!it.encrypted\n'
    '            })\n'
    '        }\n'
    '        if (lobbyLastEditedFile !== "") {\n'
    '            if (!selectNoteByName(lobbyLastEditedFile))\n'
    '                lobbyFilesIndex = Math.max(0, lobbyNotesModel.count - 1)\n'
    '            lobbyLastEditedFile = ""\n'
    '        } else if (lobbyFilesIndex >= lobbyNotesModel.count) {\n'
    '            lobbyFilesIndex = Math.max(0, lobbyNotesModel.count - 1)\n'
    '        }\n'
    '    }\n'
    '\n'
    '    function selectNoteByName(name) {\n'
    '        for (var i = 0; i < lobbyNotesModel.count; i++) {\n'
    '            if (lobbyNotesModel.get(i).name === name) {\n'
    '                lobbyFilesIndex = i\n'
    '                return true\n'
    '            }\n'
    '        }\n'
    '        return false\n'
    '    }\n'
    '\n'
    '    function encryptNoteByName(name) {\n'
    '        if (!selectNoteByName(name)) return\n'
    '        lobbyEncryptSelected()\n'
    '    }\n'
    '\n'
    '    function decryptNoteByName(name) {\n'
    '        if (!selectNoteByName(name)) return\n'
    '        lobbyDecryptSelected()\n'
    '    }\n'
    '\n'
    '    function lobbyOpenSelected() {\n'
    '        if (lobbyNotesModel.count === 0) return\n'
    '        var row = lobbyNotesModel.get(lobbyFilesIndex)\n'
    '        if (!row || row.name === "") return\n'
    '        if (row.encrypted) {\n'
    '            vaultPendingLoad = row.name\n'
    '            vaultBeginPIN("Enter PIN to edit encrypted note", true); return }\n'
    '        saveAndLoad(row.name)\n'
    '    }\n'
    '\n'
    '    function lobbyReadSelected() {\n'
    '        if (lobbyNotesModel.count === 0) return\n'
    '        var row = lobbyNotesModel.get(lobbyFilesIndex)\n'
    '        if (!row || row.name === "") return\n'
    '        if (row.encrypted) {\n'
    '            vaultPendingLoad = row.name\n'
    '            lobbyOpenInReadMode = true\n'
    '            vaultBeginPIN("Enter PIN to read encrypted note", true); return }\n'
    '        isLobby = false\n'
    '        if (mode == 1) doc = query.text\n'
    '        saveFile()\n'
    '        lobbyOpenInReadMode = true\n'
    '        doLoad(row.name)\n'
    '    }\n'
    '\n'
    '    function lobbyFilesInputDisplay() {\n'
    '        var p = lobbyFilesInputPos\n'
    '        if (p < 0) p = 0\n'
    '        if (p > lobbyFilesInput.length) p = lobbyFilesInput.length\n'
    '        return lobbyFilesInput.slice(0, p) + "_" + lobbyFilesInput.slice(p)\n'
    '    }\n'
    '\n'
    '    function lobbyFilesStripSuffix(name) {\n'
    '        if (name.endsWith(".md.enc")) return name.slice(0, -7)\n'
    '        if (name.endsWith(".md")) return name.slice(0, -3)\n'
    '        return name\n'
    '    }\n'
    '\n'
    '    function lobbyFilesBeginNew() {\n'
    '        lobbyFilesMode = "new"\n'
    '        lobbyFilesInput = ""\n'
    '        lobbyFilesInputPos = 0\n'
    '    }\n'
    '\n'
    '    function lobbyFilesBeginRename() {\n'
    '        if (lobbyNotesModel.count === 0) return\n'
    '        var n = lobbyNotesModel.get(lobbyFilesIndex).name\n'
    '        lobbyFilesInput = lobbyFilesStripSuffix(n)\n'
    '        lobbyFilesInputPos = lobbyFilesInput.length\n'
    '        lobbyFilesMode = "rename"\n'
    '    }\n'
    '\n'
    '    function lobbyFilesBeginDelete() {\n'
    '        if (lobbyNotesModel.count === 0) return\n'
    '        lobbyFilesMode = "confirm-delete"\n'
    '    }\n'
    '\n'
    '    function lobbyFilesDoDelete() {\n'
    '        if (lobbyNotesModel.count === 0) { lobbyFilesMode = ""; return }\n'
    '        writerdeck.deleteNote(lobbyNotesModel.get(lobbyFilesIndex).name)\n'
    '        lobbyFilesMode = ""\n'
    '    }\n'
    '\n'
    '    function lobbySettingsBeginExit() {\n'
    '        lobbySettingsMode = "confirm-exit"\n'
    '    }\n'
    '\n'
    '    function lobbySettingsDoExit() {\n'
    '        lobbySettingsMode = ""\n'
    '        writerdeck.exitWriterdeck()\n'
    '    }\n'
    '\n'
    '    function lobbyFilesSubmitInput() {\n'
    '        var name = lobbyFilesInput.trim()\n'
    '        if (name === "") { lobbyFilesMode = ""; return }\n'
    '        if (lobbyFilesMode === "new") {\n'
    '            writerdeck.createNote(name)\n'
    '            lobbyFilesMode = ""\n'
    '            lobbyFilesInput = ""\n'
    '            lobbyFilesInputPos = 0\n'
    '        } else if (lobbyFilesMode === "rename") {\n'
    '            var oldName = lobbyNotesModel.get(lobbyFilesIndex).name\n'
    '            var newName = name\n'
    '            if (oldName.endsWith(".md.enc")) newName = name + ".md.enc"\n'
    '            writerdeck.renameNote(oldName, newName)\n'
    '            lobbyFilesMode = ""\n'
    '            lobbyFilesInput = ""\n'
    '            lobbyFilesInputPos = 0\n'
    '        } else if (lobbyFilesMode === "new-encrypted") {\n'
    '            writerdeck.createEncryptedNote(name)\n'
    '            lobbyFilesMode = ""\n'
    '            lobbyFilesInput = ""\n'
    '            lobbyFilesInputPos = 0\n'
    '        }\n'
    '    }\n'
    '\n'
    '    function lobbyFilesBeginNewEncrypted() {\n'
    '        vaultPendingAction = "new-encrypted"\n'
    '        vaultBeginPIN("Enter PIN to create encrypted note", false)\n'
    '    }\n'
    '\n'
    '    function lobbyEncryptSelected() {\n'
    '        if (lobbyNotesModel.count === 0) return\n'
    '        var row = lobbyNotesModel.get(lobbyFilesIndex)\n'
    '        if (!row || row.encrypted) return\n'
    '        vaultPendingNote = row.name\n'
    '        vaultPendingAction = "encrypt"\n'
    '        vaultBeginPIN("Enter PIN to encrypt note", false)\n'
    '    }\n'
    '\n'
    '    function lobbyDecryptSelected() {\n'
    '        if (lobbyNotesModel.count === 0) return\n'
    '        var row = lobbyNotesModel.get(lobbyFilesIndex)\n'
    '        if (!row || !row.encrypted) return\n'
    '        vaultPendingNote = row.name\n'
    '        vaultPendingAction = "decrypt"\n'
    '        vaultBeginPIN("Enter PIN to decrypt note", false)\n'
    '    }\n'
    '\n'
    '    function vaultBeginSetup() {\n'
    '        vaultPinInput = ""\n'
    '        vaultPinPending = ""\n'
    '        vaultOverlayReason = ""\n'
    '        vaultOverlayMode = "setup"\n'
    '    }\n'
    '\n'
    '    function vaultBeginPIN(reason, keepSession) {\n'
    '        vaultPinInput = ""\n'
    '        vaultPinPending = ""\n'
    '        vaultOverlayReason = reason || ""\n'
    '        vaultPinKeepSession = !!keepSession\n'
    '        vaultOverlayMode = "pin"\n'
    '    }\n'
    '\n'
    '    function vaultBeginChangePIN() {\n'
    '        vaultPinInput = ""\n'
    '        vaultPinPending = ""\n'
    '        vaultOverlayReason = ""\n'
    '        vaultOverlayMode = "change-old"\n'
    '    }\n'
    '\n'
    '    function requestVaultPIN(reason, name) {\n'
    '        var msg = "Phone download: " + (name || "encrypted note")\n'
    '        if (reason === "download") msg = "Enter PIN on tablet to allow phone download"\n'
    '        if (name) vaultPendingLoad = name\n'
    '        vaultBeginPIN(msg, false)\n'
    '    }\n'
    '\n'
    '    function vaultNumpadCancel() {\n'
    '        vaultOverlayMode = ""\n'
    '        vaultPinInput = ""\n'
    '        vaultPinPending = ""\n'
    '        vaultOverlayReason = ""\n'
    '        vaultPendingAction = ""\n'
    '        vaultPendingNote = ""\n'
    '        vaultPendingLoad = ""\n'
    '        lobbyOpenInReadMode = false\n'
    '        lobbyVaultError = ""\n'
    '    }\n'
    '\n'
    '    function vaultPinDisplay() {\n'
    '        var n = vaultPinInput.length\n'
    '        var out = ""\n'
    '        for (var i = 0; i < 6; i++) out += i < n ? "*" : "-"\n'
    '        return out\n'
    '    }\n'
    '\n'
    '    function vaultNumpadTap(label) {\n'
    '        if (label === "Bksp") {\n'
    '            vaultPinInput = vaultPinInput.slice(0, -1)\n'
    '            return\n'
    '        }\n'
    '        if (label === "Done") {\n'
    '            vaultNumpadSubmit()\n'
    '            return\n'
    '        }\n'
    '        if (vaultPinInput.length < 6) vaultPinInput += label\n'
    '        if (vaultPinInput.length === 6) vaultNumpadSubmit()\n'
    '    }\n'
    '\n'
    '    function vaultNumpadSubmit() {\n'
    '        if (vaultPinInput.length !== 6) return\n'
    '        if (vaultOverlayMode === "setup") {\n'
    '            vaultPinPending = vaultPinInput\n'
    '            vaultPinInput = ""\n'
    '            vaultOverlayMode = "confirm"\n'
    '            return\n'
    '        }\n'
    '        if (vaultOverlayMode === "confirm") {\n'
    '            if (vaultPinInput !== vaultPinPending) { vaultNumpadCancel(); return }\n'
    '            writerdeck.setVaultPin(vaultPinInput)\n'
    '            vaultNumpadCancel()\n'
    '            return\n'
    '        }\n'
    '        if (vaultOverlayMode === "pin") {\n'
    '            writerdeck.verifyVaultPin(vaultPinInput, vaultPinKeepSession)\n'
    '            vaultOverlayMode = ""\n'
    '            vaultPinInput = ""\n'
    '            vaultPinPending = ""\n'
    '            vaultOverlayReason = ""\n'
    '            return\n'
    '        }\n'
    '        if (vaultOverlayMode === "change-old") {\n'
    '            vaultPinPending = vaultPinInput\n'
    '            vaultPinInput = ""\n'
    '            vaultOverlayMode = "change-new"\n'
    '            return\n'
    '        }\n'
    '        if (vaultOverlayMode === "change-new") {\n'
    '            vaultOverlayReason = vaultPinInput\n'
    '            vaultPinInput = ""\n'
    '            vaultOverlayMode = "change-confirm"\n'
    '            return\n'
    '        }\n'
    '        if (vaultOverlayMode === "change-confirm") {\n'
    '            if (vaultPinInput !== vaultOverlayReason) { vaultNumpadCancel(); return }\n'
    '            writerdeck.changeVaultPin(vaultPinPending, vaultPinInput)\n'
    '            vaultNumpadCancel()\n'
    '        }\n'
    '    }\n'
    '\n'
    '    function vaultHandleDigitKey(digit) {\n'
    '        if (vaultOverlayMode === "") return false\n'
    '        vaultNumpadTap(String(digit))\n'
    '        return true\n'
    '    }\n'
    '\n'
    '    function vaultConsumeKey(event) {\n'
    '        if (vaultOverlayMode === "") return false\n'
    '        if (event.key === Qt.Key_Escape) { vaultNumpadCancel(); return true }\n'
    '        if (event.key === Qt.Key_Return) { vaultNumpadSubmit(); return true }\n'
    '        if (event.key === Qt.Key_Backspace) { vaultNumpadTap("Bksp"); return true }\n'
    '        if (event.key >= Qt.Key_0 && event.key <= Qt.Key_9)\n'
    '            return vaultHandleDigitKey(event.key - Qt.Key_0)\n'
    '        if (event.text && event.text.length === 1) {\n'
    '            var d = event.text.charCodeAt(0) - 48\n'
    '            if (d >= 0 && d <= 9) return vaultHandleDigitKey(d)\n'
    '        }\n'
    '        return true\n'
    '    }\n'
    '\n'
    '    function lobbyKeyChar(event) {\n'
    '        if (event.text && event.text.length === 1 && event.modifiers === Qt.NoModifier)\n'
    '            return event.text\n'
    '        if (event.modifiers !== Qt.NoModifier) return ""\n'
    '        if (event.key >= Qt.Key_Space && event.key <= Qt.Key_AsciiTilde)\n'
    '            return String.fromCharCode(event.key)\n'
    '        return ""\n'
    '    }\n'
    '\n'
    '    function lobbyHandleKey(event) {\n'
    '        if (isOmni) return false\n'
    '        if (vaultOverlayMode !== "") {\n'
    '            return vaultConsumeKey(event)\n'
    '        }\n'
    '        if (event.modifiers & Qt.ControlModifier) {\n'
    '            if (event.key === Qt.Key_Left || event.key === Qt.Key_Right)\n'
    '                return false\n'
    '        }\n'
    '        if (lobbyFilesMode === "confirm-delete") {\n'
    '            if (event.key === Qt.Key_Escape) { lobbyFilesMode = ""; return true }\n'
    '            if (event.key === Qt.Key_Return) { lobbyFilesDoDelete(); return true }\n'
    '            return true\n'
    '        }\n'
    '        if (lobbySettingsMode === "confirm-exit") {\n'
    '            if (event.key === Qt.Key_Escape) { lobbySettingsMode = ""; return true }\n'
    '            if (event.key === Qt.Key_Return) { lobbySettingsDoExit(); return true }\n'
    '            return true\n'
    '        }\n'
    '        if (lobbyFilesMode === "new" || lobbyFilesMode === "rename" || lobbyFilesMode === "new-encrypted") {\n'
    '            if (event.key === Qt.Key_Escape) {\n'
    '                lobbyFilesMode = ""\n'
    '                lobbyFilesInput = ""\n'
    '                lobbyFilesInputPos = 0\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_Return) {\n'
    '                lobbyFilesSubmitInput()\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_Backspace) {\n'
    '                if (lobbyFilesInputPos > 0) {\n'
    '                    var bp = lobbyFilesInputPos\n'
    '                    lobbyFilesInput = lobbyFilesInput.slice(0, bp - 1) + lobbyFilesInput.slice(bp)\n'
    '                    lobbyFilesInputPos = bp - 1\n'
    '                }\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_Left && event.modifiers === Qt.NoModifier) {\n'
    '                lobbyFilesInputPos = Math.max(0, lobbyFilesInputPos - 1)\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_Right && event.modifiers === Qt.NoModifier) {\n'
    '                lobbyFilesInputPos = Math.min(lobbyFilesInput.length, lobbyFilesInputPos + 1)\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_Home && event.modifiers === Qt.NoModifier) {\n'
    '                lobbyFilesInputPos = 0\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_End && event.modifiers === Qt.NoModifier) {\n'
    '                lobbyFilesInputPos = lobbyFilesInput.length\n'
    '                return true\n'
    '            }\n'
    '            var ch = lobbyKeyChar(event)\n'
    '            if (ch !== "") {\n'
    '                var ip = lobbyFilesInputPos\n'
    '                lobbyFilesInput = lobbyFilesInput.slice(0, ip) + ch + lobbyFilesInput.slice(ip)\n'
    '                lobbyFilesInputPos = ip + 1\n'
    '                return true\n'
    '            }\n'
    '            return true\n'
    '        }\n'
    '        if (event.key === Qt.Key_Tab) {\n'
    '            if (event.modifiers & Qt.ShiftModifier)\n'
    '                lobbyGoPage((lobbyPage + lobbyTabLabels.length - 1) % lobbyTabLabels.length)\n'
    '            else\n'
    '                lobbyGoPage((lobbyPage + 1) % lobbyTabLabels.length)\n'
    '            return true\n'
    '        }\n'
    '        if (event.key >= Qt.Key_1 && event.key <= Qt.Key_6) {\n'
    '            lobbyGoPage(event.key - Qt.Key_1)\n'
    '            return true\n'
    '        }\n'
    '        if (event.key === Qt.Key_Left && event.modifiers === Qt.NoModifier) {\n'
    '            lobbyGoPage((lobbyPage + lobbyTabLabels.length - 1) % lobbyTabLabels.length)\n'
    '            return true\n'
    '        }\n'
    '        if (event.key === Qt.Key_Right && event.modifiers === Qt.NoModifier) {\n'
    '            lobbyGoPage((lobbyPage + 1) % lobbyTabLabels.length)\n'
    '            return true\n'
    '        }\n'
    '        if (lobbyPage === 3 && lobbySettingsMode === "") {\n'
    '            if (!lobbyEncryptionEnabled && event.key === Qt.Key_E && event.modifiers === Qt.NoModifier) {\n'
    '                vaultBeginSetup(); return true }\n'
    '            if (lobbyEncryptionEnabled && event.key === Qt.Key_C && event.modifiers === Qt.NoModifier) {\n'
    '                vaultBeginChangePIN(); return true }\n'
    '        }\n'
    '        if (lobbyPage === 0 && lobbyFilesMode === "" && lobbyEncryptionEnabled) {\n'
    '            if (event.key === Qt.Key_X && event.modifiers === Qt.NoModifier) {\n'
    '                lobbyEncryptSelected(); return true }\n'
    '            if (event.key === Qt.Key_Y && event.modifiers === Qt.NoModifier) {\n'
    '                lobbyDecryptSelected(); return true }\n'
    '        }\n'
    '        if (lobbyPage === 0) {\n'
    '            if (event.key === Qt.Key_Up) {\n'
    '                lobbyFilesIndex = Math.max(0, lobbyFilesIndex - 1)\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_Down) {\n'
    '                lobbyFilesIndex = Math.min(Math.max(0, lobbyNotesModel.count - 1), lobbyFilesIndex + 1)\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_Return) {\n'
    '                lobbyOpenSelected()\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_N) {\n'
    '                lobbyFilesBeginNew()\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_D) {\n'
    '                lobbyFilesBeginDelete()\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_R && !(event.modifiers & Qt.ControlModifier)) {\n'
    '                lobbyFilesBeginRename()\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_V && !(event.modifiers & Qt.ControlModifier)) {\n'
    '                lobbyReadSelected()\n'
    '                return true\n'
    '            }\n'
    '        }\n'
    '        return false\n'
    '    }\n'
    '\n'
    '    function prepareSleep() {\n'
    '        if (mode == 1) doc = query.text\n'
    '        saveFile()\n'
    '        isLobby = false\n'
    '        isSleeping = true\n'
    '    }\n'
    '\n'
    '    function handleHome(fromPhysicalCmd) {\n'
    '        if (fromPhysicalCmd === undefined)\n'
    '            fromPhysicalCmd = false\n'
    '        if (isLobby) {\n'
    '            Qt.quit()\n'
    '        } else {\n'
    '            if (mode == 1) doc = query.text\n'
    '            saveFile()\n'
    '            var lastFile = currentFile\n'
    '            isLobby = true\n'
    '            currentFile = ""\n'
    '            doc = ""\n'
    '            query.text = ""\n'
    '            autosaveSnapshot = ""\n'
    '            lobbyFilesMode = ""\n'
    '            lobbyPage = 0\n'
    '            lobbyLastEditedFile = lastFile\n'
    '            if (fromPhysicalCmd)\n'
    '                suppressNextHomeKey = true\n'
    '            lobbyRefreshNotes()\n'
    '        }\n'
    '    }\n'
    '\n'
    '    function initFile(name) {'
)
assert old7 in s, "function initFile not found (edit 7)"
s = s.replace(old7, new7, 1)

# 7b. Editor margin: a little more breathing room. Upstream textMargin is 12 px;
#     44 px ~= 5 mm at the panel's 226 dpi (kept after device-look; revert = 12).
assert 'textMargin: 12' in s, "textMargin not found in main.qml"
s = s.replace('textMargin: 12', 'textMargin: 44', 1)
assert 'textMargin: 44' in s, "textMargin 44 not found (7b2)"
s = s.replace('textMargin: 44', 'textMargin: 44\n        objectName: "writerdeckQuery"', 1)

# 7c. Block cursor (half-width) instead of the thin underline. Upstream's
#     cursorDelegate is an 18 px-wide Item holding a 2 px bottom-anchored bar.
#     (i) narrow the Item 18->9 px (half-width block), anchoring on the trailing
#     'visible: query.cursorVisible' so only the cursor delegate is hit; (ii)
#     collapse the inner bar to a filled block (anchors.fill). Both regexes are
#     whitespace-tolerant and assert so a miss fails CI loudly. A half-width
#     block reads lighter than full-width on e-ink; blink is intentionally not
#     added (it ghosts/churns on e-ink). Kept after device-look; revert = drop 7c.
s, nwid = re.subn(
    r'width:\s*18[ \t]*\n(\s*visible:\s*query\.cursorVisible)',
    r'width: 9\n\1',
    s, count=1
)
assert nwid == 1, "cursor delegate Item width not found in main.qml"
s, ncur = re.subn(
    r'anchors\.bottom:\s*parent\.bottom\s+'
    r'anchors\.bottomMargin:\s*4\s+'
    r'color:\s*"black"\s+'
    r'height:\s*2\s+'
    r'width:\s*parent\.width',
    'anchors.fill: parent\n                            color: "black"',
    s, count=1
)
assert ncur == 1, "cursor underline delegate not found in main.qml"

# 7d. Scroll direction: the rM1 physical page buttons arrive as Key_Left (prev
#     page) and Key_Right (next page). Upstream groups Left with Down->scrollDown
#     and Right with Up->scrollUp, so pressing the right/next-page button scrolls
#     UP -- backwards. Split the cases and reverse Left/Right.
s, ndir = re.subn(
    r'case Qt\.Key_Down:\s+case Qt\.Key_Left:\s+'
    r'if \(mode == 0\)\s+flick\.scrollDown\(\)\s+break\s+'
    r'case Qt\.Key_Up:\s+case Qt\.Key_Right:\s+'
    r'if \(mode == 0\)\s+flick\.scrollUp\(\)\s+break',
    ('case Qt.Key_Down:\n'
     '                            if (mode == 0)\n'
     '                                flick.scrollDown()\n'
     '                            break\n'
     '                        case Qt.Key_Left:\n'
     '                            flick.scrollUp()\n'
     '                            break\n'
     '                        case Qt.Key_Up:\n'
     '                            if (mode == 0)\n'
     '                                flick.scrollUp()\n'
     '                            break\n'
     '                        case Qt.Key_Right:\n'
     '                            flick.scrollDown()\n'
     '                            break'),
    s, count=1
)
assert ndir == 1, "scroll direction case block not found in main.qml"

# 7e. Scroll amount: 4/5-screen step 400->1500 px (screen is 1872px tall).
#     Key_Left/Key_Right (rM1 physical page buttons) now also scroll in edit
#     mode -- the mode==0 guard was removed for those two cases in step 7d.
assert 'contentY -= 400' in s, "scrollUp contentY not found in main.qml"
assert 'contentY += 400' in s, "scrollDown contentY not found in main.qml"
s = s.replace('contentY -= 400', 'contentY -= 1500', 1)
s = s.replace('contentY += 400', 'contentY += 1500', 1)

# 7m. Read view: don't auto-scroll to bottom on Esc (toggle to preview).
#     Upstream onCursorRectangleChanged always calls flick.ensureVisible; when
#     mode flips to RichText the reflow chases the edit cursor (often the doc
#     end). Only auto-scroll while editing -- read mode uses arrow/page scroll.
old7m = (
    '                onCursorRectangleChanged: {\n'
    '                    flick.ensureVisible(cursorRectangle)\n'
    '                }'
)
new7m = (
    '                onCursorRectangleChanged: {\n'
    '                    if (mode == 1) {\n'
    '                        var margin = 120\n'
    '                        var viewTop = flick.contentY + margin\n'
    '                        var viewBot = flick.contentY + flick.height - margin\n'
    '                        var cy = cursorRectangle.y\n'
    '                        var cb = cy + cursorRectangle.height\n'
    '                        if (cy < viewTop || cb > viewBot)\n'
    '                            flick.ensureVisible(cursorRectangle, margin)\n'
    '                    }\n'
    '                }'
)
assert old7m in s, "onCursorRectangleChanged ensureVisible block not found (edit 7m)"
s = s.replace(old7m, new7m, 1)

# 7f. Paragraph spacing in Read (RichText) view: inject margin-bottom via a
#     root readHtml() helper so Edit (PlainText) is untouched. Risk: Qt RichText
#     CSS subset is partial -- margin-bottom on <p> is expected to work but must
#     be eyeballed on e-ink. Fallback if ignored: line-height or spacer paragraph.
assert '    property string lobbyPIN: ""' in s, "lobbyPIN property not found (edit 7f)"
s = s.replace('    property string lobbyPIN: ""',
              '    property string lobbyPIN: ""\n    property int paraSpacing: 28', 1)
old7f_fn = '    function initFile(name) {'
new7f_fn = (
    '    function readHtml(d) {\n' +
    '        return utils.markdown(d)\n' +
    '            .replace(/<p>/g, ' + chr(39) + '<p style="margin-bottom:' + chr(39) + ' + paraSpacing + ' + chr(39) + 'px">' + chr(39) + ')\n' +
    '            .replace(/<li>/g, ' + chr(39) + '<li style="margin-bottom:8px">' + chr(39) + ')\n' +
    '    }\n' +
    '\n' +
    '    function syncQueryDisplay() {\n' +
    '        if (mode == 0) {\n' +
    '            query.textFormat = TextEdit.RichText\n' +
    '            query.text = readHtml(doc)\n' +
    '        } else {\n' +
    '            query.textFormat = TextEdit.PlainText\n' +
    '            query.text = doc\n' +
    '        }\n' +
    '    }\n' +
    '\n' +
    '    function initFile(name) {'
)
assert old7f_fn in s, "function initFile not found (edit 7f)"
s = s.replace(old7f_fn, new7f_fn, 1)
# 7f only swaps utils.markdown -> readHtml in toggle paths; edit 7u drops the
# TextEdit `text:` binding (broken by typing and doLoad) in favour of syncQueryDisplay().
assert 'text: mode == 0 ? utils.markdown(doc) : doc' in s, "text binding not found (edit 7f)"
s = s.replace('text: mode == 0 ? utils.markdown(doc) : doc',
              'text: mode == 0 ? root.readHtml(doc) : doc', 1)

# 7g. Reading-view font: add readFont property (default Inter), a setReadFont()
#     function so the C++ setfont cmd can call it via invokeMethod, and change
#     the read-mode font.family binding to use it. Edit mode keeps Noto Mono.
assert 'property int paraSpacing: 28' in s, "paraSpacing not found (edit 7g)"
s = s.replace('property int paraSpacing: 28',
              'property int paraSpacing: 28\n    property string readFont: "Inter"', 1)
old7g_fn = '    function initFile(name) {'
new7g_fn = (
    '    function setReadFont(name) {\n'
    '        readFont = name\n'
    '        if (mode == 0) syncQueryDisplay()\n'
    '    }\n'
    '\n'
    '    function initFile(name) {'
)
assert old7g_fn in s, "function initFile not found (edit 7g)"
s = s.replace(old7g_fn, new7g_fn, 1)
assert 'font.family: mode == 0 ? ' + chr(34) + 'Noto Sans' + chr(34) + ' : ' + chr(34) + 'Noto Mono' + chr(34) in s, \
    "font.family Noto Sans binding not found (edit 7g)"
s = s.replace('font.family: mode == 0 ? ' + chr(34) + 'Noto Sans' + chr(34) + ' : ' + chr(34) + 'Noto Mono' + chr(34),
              'font.family: mode == 0 ? readFont : ' + chr(34) + 'Noto Mono' + chr(34), 1)

# 7l. Add showLobby() QML function: idempotent, saves if editing, sets
#     isLobby=true, never quits. Called by the C++ showlobby cmd (Feature L).
#     Distinct from handleHome which quits from the Lobby -- showLobby is safe
#     to call when already on the Lobby (it is a no-op in that case).
#     Anchors on initFile (the upstream function) like 7f/7g/7h -- noteDeleted
#     does not exist yet at this point (7h inserts it just below).
old7l = '    function initFile(name) {'
new7l = (
    '    function showLobby() {\n'
    '        var lastFile = ""\n'
    '        if (!isLobby) {\n'
    '            if (mode == 1) doc = query.text\n'
    '            saveFile()\n'
    '            lastFile = currentFile\n'
    '            currentFile = ""\n'
    '            doc = ""\n'
    '            query.text = ""\n'
    '            autosaveSnapshot = ""\n'
    '        }\n'
    '        isLobby = true\n'
    '        lobbyFilesMode = ""\n'
    '        lobbyPage = 0\n'
    '        if (lastFile !== "") lobbyLastEditedFile = lastFile\n'
    '        lobbyRefreshNotes()\n'
    '    }\n'
    '\n'
    '    function initFile(name) {'
)
assert old7l in s, "function initFile not found (edit 7l)"
s = s.replace(old7l, new7l, 1)

# 7h. noteDeleted(): go to the Lobby and CLEAR currentFile so the next save
#     can't resurrect the just-deleted file. Called by the C++ notedeleted cmd
#     when the server deletes the open note (phone or socket). No saveFile().
old7h = '    function initFile(name) {'
new7h = (
    '    function noteDeleted() {\n'
    '        currentFile = ""\n'
    '        doc = ""\n'
    '        query.text = ""\n'
    '        autosaveSnapshot = ""\n'
    '        isLobby = true\n'
    '        lobbyFilesMode = ""\n'
    '        lobbyRefreshNotes()\n'
    '    }\n'
    '\n'
    '    function noteRenamed(name) {\n'
    '        if (!isLobby) currentFile = name\n'
    '    }\n'
    '\n'
    '    function reloadNote() {\n'
    '        if (currentFile !== "") doLoad(currentFile)\n'
    '    }\n'
    '\n'
    '    function initFile(name) {'
)
assert old7h in s, "function initFile not found (edit 7h)"
s = s.replace(old7h, new7h, 1)

# 7i. saveFile contract: plain UTF-8 only; reject qrichtext/HTML; guard empty path.
#     In edit mode (mode==1) read query.text; never read query.text in preview
#     (RichText). Reject HTML payloads before write (content-fidelity slice 2).
old7i = (
    '    function saveFile() {\n'
    '        console.log("Save " + currentFile)\n'
    '        var fileUrl = folder + currentFile\n'
    '        console.log(fileUrl)\n'
    '        var request = new XMLHttpRequest()\n'
    '        request.open("PUT", fileUrl, false)\n'
    '        request.send(doc)\n'
    '        console.log("save -> " + request.status + " " + request.statusText)\n'
    '        return request.status\n'
    '    }'
)
new7i = (
    '    function saveFile() {\n'
    '        if (currentFile === "") return 0\n'
    '        var content = (mode == 1) ? query.text : doc\n'
    '        if (isHtmlPayload(content)) {\n'
    '            console.log("save rejected: HTML/qrichtext payload for " + currentFile)\n'
    '            return 0\n'
    '        }\n'
    '        if (mode == 1) doc = content\n'
    '        console.log("Save " + currentFile)\n'
    '        var url = "http://127.0.0.1:8000/api/notes/" + encodeURIComponent(currentFile)\n'
    '        console.log(url)\n'
    '        var request = new XMLHttpRequest()\n'
    '        request.open("PUT", url, false)\n'
    '        request.setRequestHeader("Content-Type", "application/json")\n'
    '        request.send(JSON.stringify({ content: content }))\n'
    '        console.log("save -> " + request.status + " " + request.statusText)\n'
    '        return request.status\n'
    '    }'
)
assert old7i in s, "saveFile body not found (edit 7i)"
s = s.replace(old7i, new7i, 1)

# 7t. Autosave: periodic flush of query.text -> disk while editing (slice 9).
#     Uses saveFile() (loopback HTTP -> atomic writeNoteFile on server).
#     autosaveSnapshot property is added in 7c2a (after readFont).
old7t = '    function reloadNote() {'
new7t = (
    '    function autosaveTick() {\n'
    '        if (harnessPrepareLock || isLobby || currentFile === "" || mode != 1) return\n'
    '        if (query.text === autosaveSnapshot) return\n'
    '        saveFile()\n'
    '        if (currentFile !== "" && !isHtmlPayload(query.text)) autosaveSnapshot = query.text\n'
    '    }\n'
    '\n'
    '    function reloadNote() {'
)
assert old7t in s, "reloadNote not found (edit 7t)"
s = s.replace(old7t, new7t, 1)

# 7j. Demote scratch to an ordinary note: clear the currentFile default so boot
#     starts in a clean no-file state. Combined with the 7i saveFile guard and
#     handleHome's clearance (edit 7), no path can resurrect a deleted note by
#     saving the last doc into a cleared-currentFile path.
assert 'property string currentFile: "scratch.md"' in s, "currentFile default not found (edit 7j)"
s = s.replace('property string currentFile: "scratch.md"',
              'property string currentFile: ""', 1)

# 7k. Guard Component.onCompleted boot load: with currentFile now "" by default,
#     the unconditional doLoad(currentFile) would attempt a GET to an empty URL.
#     Wrap it so a fresh boot (no note selected) skips the load cleanly.
s, nboot = re.subn(
    r'Component\.onCompleted:\s*\{\s*doLoad\(currentFile\)\s*\}',
    'Component.onCompleted: {\n        if (currentFile !== "") doLoad(currentFile)\n    }',
    s, count=1
)
assert nboot == 1, "Component.onCompleted doLoad pattern not found (edit 7k)"

# 7c2. Cursor: hidden while typing, strong block while idle / navigating.
#      Device-verified superior on e-ink: the block vanishes as you type (no
#      ghosting/smear trail) and reappears 500 ms after the last keystroke
#      (Timer, step 8b) or immediately on Up/Down. Blink is intentionally absent.

# 7c2a. Add cursor-state property after readFont.
assert '    property string readFont: "Inter"' in s, "readFont not found (7c2a)"
s = s.replace(
    '    property string readFont: "Inter"',
    '    property string readFont: "Inter"\n'
    '    property bool cursorStrong: true\n'
    '    property string autosaveSnapshot: ""\n'
    '    property int harnessTextWidth: 0\n'
    '    property int harnessDefaultQueryWidth: 0\n'
    '    property bool harnessPrepareLock: false\n'
    '    property int goalColumn: -1\n'
    '    property int macUndoMarkPos: -1\n'
    '    property string macUndoKind: ""\n'
    '    property int _prevTextLen: 0\n'
    '    property int _prevCursor: 0\n'
    '    property int pendingRedoCursor: -1',
    1
)

# 7c2b. Cursor delegate: hide the block while typing (visible only when strong).
s, ncwv = re.subn(
    r'width:\s*9\n(\s*visible:\s*)query\.cursorVisible',
    r'width: 9\n\1query.cursorVisible && cursorStrong',
    s, count=1
)
assert ncwv == 1, "cursor delegate visible not found (7c2b)"

# 7c2d. Arrow Down in edit mode -> strong cursor immediately.
old_kd = (
    'case Qt.Key_Down:\n'
    '                            if (mode == 0)\n'
    '                                flick.scrollDown()\n'
    '                            break'
)
new_kd = (
    'case Qt.Key_Down:\n'
    '                            if (mode == 1) { cursorStrong = true; cursorTimer.stop() }\n'
    '                            if (mode == 0)\n'
    '                                flick.scrollDown()\n'
    '                            break'
)
assert old_kd in s, "Key_Down case not found (7c2d)"
s = s.replace(old_kd, new_kd, 1)

# 7c2e. Arrow Up in edit mode -> strong cursor immediately.
old_ku = (
    'case Qt.Key_Up:\n'
    '                            if (mode == 0)\n'
    '                                flick.scrollUp()\n'
    '                            break'
)
new_ku = (
    'case Qt.Key_Up:\n'
    '                            if (mode == 1) { cursorStrong = true; cursorTimer.stop() }\n'
    '                            if (mode == 0)\n'
    '                                flick.scrollUp()\n'
    '                            break'
)
assert old_ku in s, "Key_Up case not found (7c2e)"
s = s.replace(old_ku, new_ku, 1)

# 7n. Cursor boundary niceties (edit mode): Down on the last line -> end of line;
#     Up on the first line -> start of line. Accept the key so QTextEdit does not
#     no-op or scroll the flickable instead.
old7n_kd = (
    'case Qt.Key_Down:\n'
    '                            if (mode == 1) { cursorStrong = true; cursorTimer.stop() }\n'
    '                            if (mode == 0)\n'
    '                                flick.scrollDown()\n'
    '                            break'
)
new7n_kd = (
    'case Qt.Key_Down:\n'
    '                            if (mode == 1) {\n'
    '                                cursorStrong = true; cursorTimer.stop()\n'
    '                                if (event.modifiers === Qt.NoModifier && root.cursorOnLastLine()) {\n'
    '                                    root.moveCursorEndOfLine()\n'
    '                                    event.accepted = true\n'
    '                                    break\n'
    '                                }\n'
    '                            }\n'
    '                            if (mode == 0)\n'
    '                                flick.scrollDown()\n'
    '                            break'
)
assert old7n_kd in s, "Key_Down case not found (7n)"
s = s.replace(old7n_kd, new7n_kd, 1)

old7n_ku = (
    'case Qt.Key_Up:\n'
    '                            if (mode == 1) { cursorStrong = true; cursorTimer.stop() }\n'
    '                            if (mode == 0)\n'
    '                                flick.scrollUp()\n'
    '                            break'
)
new7n_ku = (
    'case Qt.Key_Up:\n'
    '                            if (mode == 1) {\n'
    '                                cursorStrong = true; cursorTimer.stop()\n'
    '                                if (event.modifiers === Qt.NoModifier && root.cursorOnFirstLine()) {\n'
    '                                    root.moveCursorStartOfLine()\n'
    '                                    event.accepted = true\n'
    '                                    break\n'
    '                                }\n'
    '                            }\n'
    '                            if (mode == 0)\n'
    '                                flick.scrollUp()\n'
    '                            break'
)
assert old7n_ku in s, "Key_Up case not found (7n)"
s = s.replace(old7n_ku, new7n_ku, 1)

old7n_fn = '    function showLobby() {'
new7n_fn = (
    '    function harnessSetWidth(w) {\n'
    '        if (mode != 1) return\n'
    '        if (harnessDefaultQueryWidth <= 0 && query.width > 0)\n'
    '            harnessDefaultQueryWidth = query.width\n'
    '        if (w > 0) {\n'
    '            harnessTextWidth = w\n'
    '            query.width = w\n'
    '        } else if (harnessDefaultQueryWidth > 0) {\n'
    '            harnessTextWidth = 0\n'
    '            query.width = harnessDefaultQueryWidth\n'
    '        }\n'
    '    }\n'
    '\n'
    '    function harnessOpenNote(name) {\n'
    '        isLobby = false\n'
    '        mode = 1\n'
    '        currentFile = name\n'
    '    }\n'
    '\n'
    '    function harnessSandboxReset(widthPx) {\n'
    '        if (currentFile === "") return\n'
    '        harnessPrepareLock = true\n'
    '        isLobby = false\n'
    '        mode = 1\n'
    '        cursorStrong = true\n'
    '        cursorTimer.stop()\n'
    '        if (widthPx > 0)\n'
    '            harnessSetWidth(widthPx)\n'
    '        else if (harnessTextWidth > 0)\n'
    '            query.width = harnessTextWidth\n'
    '        var req = new XMLHttpRequest()\n'
    '        req.open("GET", "http://127.0.0.1:8000/api/notes/" + encodeURIComponent(currentFile), false)\n'
    '        req.send()\n'
    '        if (req.status !== 200) {\n'
    '            harnessPrepareLock = false\n'
    '            console.log("harnessSandboxReset: GET failed " + req.status)\n'
    '            return\n'
    '        }\n'
    '        var response = sanitizeLoadedNote(req.responseText)\n'
    '        doc = response\n'
    '        syncQueryDisplay()\n'
    '        autosaveSnapshot = response\n'
    '        if (widthPx > 0 && query.text.length > 0)\n'
    '            query.positionToRectangle(query.text.length)\n'
    '        query.deselect()\n'
    '        query.cursorPosition = 0\n'
    '        goalColumn = -1\n'
    '        macUndoKind = ""\n'
    '        macUndoMarkPos = -1\n'
    '        pendingRedoCursor = -1\n'
    '        _prevTextLen = query.text.length\n'
    '        _prevCursor = 0\n'
    '        ctrlPressed = false\n'
    '        query.forceActiveFocus()\n'
    '        if (query.text.length > 0) {\n'
    '            query.cursorPosition = query.text.length\n'
    '            query.cursorPosition = 0\n'
    '            query.deselect()\n'
    '        }\n'
    '        if (typeof flick !== "undefined")\n'
    '            flick.contentY = 0\n'
    '        try { if (query.undoStack) query.undoStack.clear() } catch (e) {}\n'
    '        harnessPrepareLock = false\n'
    '    }\n'
    '\n'
    '    function socketRouteKey(key, mods) {\n'
    '        if (mode != 1 || isLobby) return\n'
    '        query.forceActiveFocus()\n'
    '        key = parseInt(key)\n'
    '        mods = parseInt(mods)\n'
    '        var cmd = (mods & Qt.ControlModifier) !== 0\n'
    '        var alt = (mods & Qt.AltModifier) !== 0\n'
    '        var shift = (mods & Qt.ShiftModifier) !== 0\n'
    '        var text = query.text\n'
    '        var pos = query.cursorPosition\n'
    '        if (!shift && cmd && !alt) {\n'
    '            if (key === Qt.Key_Right) { moveCursorTo(lineEndPos(pos, text), false); return }\n'
    '            if (key === Qt.Key_Left) { moveCursorTo(lineStartPos(pos, text), false); return }\n'
    '            if (key === Qt.Key_Up) { moveCursorTo(0, false); return }\n'
    '            if (key === Qt.Key_Down) { moveCursorTo(text.length, false); return }\n'
    '            if (key === Qt.Key_End) { moveCursorTo(text.length, false); return }\n'
    '            if (key === Qt.Key_Home) { moveCursorTo(0, false); return }\n'
    '        }\n'
    '        if (!shift && !cmd && alt) {\n'
    '            if (key === Qt.Key_Right) { moveCursorTo(wordRightPos(pos, text), false); return }\n'
    '            if (key === Qt.Key_Left) { moveCursorTo(wordLeftPos(pos, text), false); return }\n'
    '            if (key === Qt.Key_Up) { moveCursorTo(paragraphUpPos(pos, text), false); return }\n'
    '            if (key === Qt.Key_Down) { moveCursorTo(paragraphDownPos(pos, text), false); return }\n'
    '        }\n'
    '        if (shift && cmd && !alt) {\n'
    '            if (key === Qt.Key_Left) { moveCursorTo(0, true); publishEditorState(); return }\n'
    '            if (key === Qt.Key_Right) { moveCursorTo(text.length, true); publishEditorState(); return }\n'
    '            if (key === Qt.Key_Up) { moveCursorTo(0, true); publishEditorState(); return }\n'
    '            if (key === Qt.Key_Down) { moveCursorTo(text.length, true); publishEditorState(); return }\n'
    '        }\n'
    '        if (shift && alt && !cmd) {\n'
    '            var ap = pos\n'
    '            if (key === Qt.Key_Left) ap = wordLeftPos(pos, text)\n'
    '            else if (key === Qt.Key_Right) ap = wordRightPos(pos, text)\n'
    '            else if (key === Qt.Key_Up) ap = paragraphUpPos(pos, text)\n'
    '            else if (key === Qt.Key_Down) ap = paragraphDownPos(pos, text)\n'
    '            else ap = -1\n'
    '            if (ap >= 0) { moveCursorTo(ap, true); publishEditorState(); return }\n'
    '        }\n'
    '        var event = { key: key, modifiers: mods, accepted: false }\n'
    '        if (handleMacEditKeys(event)) return\n'
    '        if (handleMacUndo(event)) return\n'
    '        if (handleMacBackspace(event)) return\n'
    '        if (handleMacArrow(event)) return\n'
    '        if (!shift && !cmd && !alt) {\n'
    '            if (key === Qt.Key_Right) { moveCursorTo(Math.min(pos + 1, query.text.length), false); publishEditorState(); return }\n'
    '            if (key === Qt.Key_Left) { moveCursorTo(Math.max(0, pos - 1), false); publishEditorState(); return }\n'
    '        }\n'
    '        publishEditorState()\n'
    '    }\n'
    '\n'
    '    function publishEditorState() {\n'
    '        writerdeck.publishState(query.cursorPosition, query.selectionStart,\n'
    '            query.selectionEnd, query.text.length, mode, isLobby ? 1 : 0,\n'
    '            vaultOverlayMode, currentFile, query.text)\n'
    '    }\n'
    '\n'
    '    function cursorOnLastLine() {\n'
    '        if (mode != 1) return false\n'
    '        var len = query.text.length\n'
    '        if (len === 0) return true\n'
    '        var curY = query.positionToRectangle(query.cursorPosition).y\n'
    '        var endY = query.positionToRectangle(len).y\n'
    '        return curY >= endY - 1\n'
    '    }\n'
    '\n'
    '    function cursorOnFirstLine() {\n'
    '        if (mode != 1) return false\n'
    '        var curY = query.positionToRectangle(query.cursorPosition).y\n'
    '        var topY = query.positionToRectangle(0).y\n'
    '        return curY <= topY + 1\n'
    '    }\n'
    '\n'
    '    function isSpaceChar(c) {\n'
    '        return c === " " || c === "\\t" || c === "\\n"\n'
    '    }\n'
    '\n'
    '    function lineStartPos(pos, text) {\n'
    '        var prev = text.lastIndexOf("\\n", pos - 1)\n'
    '        return prev === -1 ? 0 : prev + 1\n'
    '    }\n'
    '\n'
    '    function lineEndPos(pos, text) {\n'
    '        var nl = text.indexOf("\\n", pos)\n'
    '        return nl === -1 ? text.length : nl\n'
    '    }\n'
    '\n'
    '    function lineCharCount(lineStart, text) {\n'
    '        var end = lineEndPos(lineStart, text)\n'
    '        return end - lineStart\n'
    '    }\n'
    '\n'
    '    function rememberGoalColumn(pos, text) {\n'
    '        goalColumn = pos - lineStartPos(pos, text)\n'
    '    }\n'
    '\n'
    '    function deleteWordLeftPos(pos, text) {\n'
    '        if (pos <= 0) return 0\n'
    '        var start = wordLeftPos(pos, text)\n'
    '        if (start > 0) {\n'
    '            var s = start - 1\n'
    '            while (s >= 0 && isSpaceChar(text.charAt(s))) s--\n'
    '            start = s + 1\n'
    '        }\n'
    '        return start\n'
    '    }\n'
    '\n'
    '    function deleteLineLeftPos(pos, text) {\n'
    '        if (pos <= 0) return 0\n'
    '        var start = lineStartPos(pos, text)\n'
    '        if (pos > start) return start\n'
    '        if (start === 0) return 0\n'
    '        return lineStartPos(start - 1, text)\n'
    '    }\n'
    '\n'
    '    function wordLeftPos(pos, text) {\n'
    '        if (pos <= 0) return 0\n'
    '        if (pos > text.length) pos = text.length\n'
    '        pos--\n'
    '        while (pos > 0 && isSpaceChar(text.charAt(pos))) pos--\n'
    '        while (pos > 0 && !isSpaceChar(text.charAt(pos - 1))) pos--\n'
    '        if (pos < text.length && isSpaceChar(text.charAt(pos))) pos++\n'
    '        return pos\n'
    '    }\n'
    '\n'
    '    function wordRightPos(pos, text) {\n'
    '        var len = text.length\n'
    '        if (pos >= len) return len\n'
    '        while (pos < len && !isSpaceChar(text.charAt(pos))) pos++\n'
    '        while (pos < len && isSpaceChar(text.charAt(pos))) pos++\n'
    '        return pos\n'
    '    }\n'
    '\n'
    '    function paragraphUpPos(pos, text) {\n'
    '        var lineStart = lineStartPos(pos, text)\n'
    '        if (lineStart === 0) return 0\n'
    '        var i = lineStart - 1\n'
    '        while (i > 0) {\n'
    '            if (text.charAt(i) === "\\n" && text.charAt(i - 1) === "\\n")\n'
    '                return i + 1\n'
    '            i--\n'
    '        }\n'
    '        return 0\n'
    '    }\n'
    '\n'
    '    function paragraphDownPos(pos, text) {\n'
    '        var len = text.length\n'
    '        var lineEnd = lineEndPos(pos, text)\n'
    '        if (lineEnd >= len) return len\n'
    '        for (var i = lineEnd; i < len - 1; i++) {\n'
    '            if (text.charAt(i) === "\\n" && text.charAt(i + 1) === "\\n")\n'
    '                return i + 2\n'
    '        }\n'
    '        return lineEnd + 1\n'
    '    }\n'
    '\n'
    '    function moveCursorTo(newPos, extend, keepGoalColumn) {\n'
    '        var len = query.text.length\n'
    '        var text = query.text\n'
    '        newPos = Math.max(0, Math.min(newPos, len))\n'
    '        if (!extend) {\n'
    '            query.deselect()\n'
    '            query.cursorPosition = newPos\n'
    '            if (!keepGoalColumn) rememberGoalColumn(newPos, text)\n'
    '            return\n'
    '        }\n'
    '        var anchor = query.cursorPosition\n'
    '        if (query.selectionStart !== query.selectionEnd) {\n'
    '            anchor = (query.cursorPosition === query.selectionEnd)\n'
    '                ? query.selectionStart : query.selectionEnd\n'
    '        }\n'
    '        var a = Math.min(anchor, newPos)\n'
    '        var b = Math.max(anchor, newPos)\n'
    '        query.select(a, b)\n'
    '        query.cursorPosition = newPos\n'
    '    }\n'
    '\n'
    '    function extendSelectionHorizontal(newPos) {\n'
    '        var anchor\n'
    '        if (query.selectionStart === query.selectionEnd) {\n'
    '            anchor = query.cursorPosition\n'
    '        } else if (newPos >= query.cursorPosition) {\n'
    '            anchor = Math.min(query.selectionStart, query.selectionEnd)\n'
    '        } else {\n'
    '            anchor = Math.max(query.selectionStart, query.selectionEnd)\n'
    '        }\n'
    '        query.select(Math.min(anchor, newPos), Math.max(anchor, newPos))\n'
    '        query.cursorPosition = newPos\n'
    '    }\n'
    '\n'
    '    function visualLineDownPos(pos) {\n'
    '        var len = query.text.length\n'
    '        if (pos >= len) return len\n'
    '        var curRect = query.positionToRectangle(pos)\n'
    '        var nextLineY = curRect.y + curRect.height - 0.1\n'
    '        var goalX = curRect.x\n'
    '        var best = -1\n'
    '        var bestDist = 1e12\n'
    '        for (var p = pos + 1; p <= len; p++) {\n'
    '            var r = query.positionToRectangle(p)\n'
    '            if (r.y + 0.1 < nextLineY) continue\n'
    '            if (best >= 0 && r.y > query.positionToRectangle(best).y + 0.1) break\n'
    '            var dist = Math.abs(r.x - goalX)\n'
    '            if (best < 0 || dist < bestDist) { best = p; bestDist = dist }\n'
    '        }\n'
    '        return best >= 0 ? best : len\n'
    '    }\n'
    '\n'
    '    function visualLineUpPos(pos) {\n'
    '        if (pos <= 0) return 0\n'
    '        var curRect = query.positionToRectangle(pos)\n'
    '        var rowTop = curRect.y - curRect.height + 0.1\n'
    '        var goalX = curRect.x\n'
    '        var best = -1\n'
    '        var bestDist = 1e12\n'
    '        var targetY = -1\n'
    '        for (var p = pos - 1; p >= 0; p--) {\n'
    '            var r = query.positionToRectangle(p)\n'
    '            if (r.y > rowTop) continue\n'
    '            if (targetY < 0) targetY = r.y\n'
    '            if (Math.abs(r.y - targetY) > 0.5) break\n'
    '            var dist = Math.abs(r.x - goalX)\n'
    '            if (best < 0 || dist < bestDist) { best = p; bestDist = dist }\n'
    '        }\n'
    '        return best >= 0 ? best : 0\n'
    '    }\n'
    '\n'
    '    function visualLineStartPos(pos) {\n'
    '        if (pos <= 0) return 0\n'
    '        var curY = query.positionToRectangle(pos).y\n'
    '        var best = pos\n'
    '        for (var p = pos - 1; p >= 0; p--) {\n'
    '            if (Math.abs(query.positionToRectangle(p).y - curY) < 0.5) best = p\n'
    '            else break\n'
    '        }\n'
    '        return best\n'
    '    }\n'
    '\n'
    '    function visualLineEndPos(pos) {\n'
    '        var len = query.text.length\n'
    '        if (pos >= len) return len\n'
    '        var curY = query.positionToRectangle(pos).y\n'
    '        var best = pos\n'
    '        for (var p = pos + 1; p <= len; p++) {\n'
    '            if (Math.abs(query.positionToRectangle(p).y - curY) < 0.5) best = p\n'
    '            else break\n'
    '        }\n'
    '        return best\n'
    '    }\n'
    '\n'
    '    function onWrappedLine(pos, text) {\n'
    '        var s = lineStartPos(pos, text)\n'
    '        var nl = text.indexOf("\\n", s)\n'
    '        return nl === -1 || nl > lineEndPos(pos, text)\n'
    '    }\n'
    '\n'
    '    function macLineStartPos(pos, text) {\n'
    '        return onWrappedLine(pos, text) ? visualLineStartPos(pos)\n'
    '            : lineStartPos(pos, text)\n'
    '    }\n'
    '\n'
    '    function macLineEndPos(pos, text) {\n'
    '        return onWrappedLine(pos, text) ? visualLineEndPos(pos)\n'
    '            : lineEndPos(pos, text)\n'
    '    }\n'
    '\n'
    '    function lineDownPos(pos, text) {\n'
    '        var nl = text.indexOf("\\n", pos)\n'
    '        if (nl === -1) return visualLineDownPos(pos)\n'
    '        var col = goalColumn >= 0 ? goalColumn : (pos - lineStartPos(pos, text))\n'
    '        goalColumn = col\n'
    '        var nextStart = nl + 1\n'
    '        while (nextStart < text.length) {\n'
    '            var nextChars = lineCharCount(nextStart, text)\n'
    '            if (nextChars <= 0) return nextStart\n'
    '            if (col < nextChars) return nextStart + col\n'
    '            pos = nextStart + nextChars\n'
    '            if (pos >= text.length) return text.length\n'
    '            nl = text.indexOf("\\n", pos)\n'
    '            if (nl === -1) return text.length\n'
    '            nextStart = nl + 1\n'
    '        }\n'
    '        return text.length\n'
    '    }\n'
    '\n'
    '    function lineUpPos(pos, text) {\n'
    '        if (pos <= 0) return 0\n'
    '        var vis = visualLineUpPos(pos)\n'
    '        if (vis < pos) return vis\n'
    '        var curStart = lineStartPos(pos, text)\n'
    '        if (pos > curStart) {\n'
    '            var vis2 = visualLineUpPos(pos)\n'
    '            if (vis2 < pos) return vis2\n'
    '        }\n'
    '        if (curStart === 0) return 0\n'
    '        var col = goalColumn >= 0 ? goalColumn : (pos - curStart)\n'
    '        goalColumn = col\n'
    '        var prevStart = lineStartPos(curStart - 1, text)\n'
    '        while (prevStart >= 0) {\n'
    '            var prevChars = lineCharCount(prevStart, text)\n'
    '            if (prevChars <= 0) return prevStart\n'
    '            if (col < prevChars) return prevStart + col\n'
    '            if (prevStart === 0) return 0\n'
    '            prevStart = lineStartPos(prevStart - 1, text)\n'
    '        }\n'
    '        return 0\n'
    '    }\n'
    '\n'
    '    function lineUpForSelection(head, anchor, text) {\n'
    '        if (head === 0) return 0\n'
    '        if (head === lineStartPos(head, text) && head > anchor)\n'
    '            return head - 1\n'
    '        return lineUpPos(head, text)\n'
    '    }\n'
    '\n'
    '    function moveCursorVertical(down) {\n'
    '        var pos = query.cursorPosition\n'
    '        var text = query.text\n'
    '        var newPos = down ? lineDownPos(pos, text) : lineUpPos(pos, text)\n'
    '        if (newPos === pos) {\n'
    '            if (down) {\n'
    '                if (text.indexOf("\\n", pos) === -1) {\n'
    '                    var vis = visualLineDownPos(pos)\n'
    '                    newPos = (vis > pos) ? vis : lineEndPos(pos, text)\n'
    '                }\n'
    '            } else {\n'
    '                newPos = macLineStartPos(pos, text)\n'
    '            }\n'
    '        }\n'
    '        moveCursorTo(newPos, false, true)\n'
    '    }\n'
    '\n'
    '    function extendSelectionVertical(down) {\n'
    '        var text = query.text\n'
    '        var head = query.cursorPosition\n'
    '        var anchor\n'
    '        if (query.selectionStart === query.selectionEnd) {\n'
    '            anchor = head\n'
    '        } else if (query.selectionStart === head) {\n'
    '            anchor = query.selectionEnd\n'
    '        } else {\n'
    '            anchor = query.selectionStart\n'
    '        }\n'
    '        var newHead = down ? lineDownPos(head, text) : lineUpForSelection(head, anchor, text)\n'
    '        if (down && newHead > head && head > lineStartPos(head, text))\n'
    '            newHead = lineEndPos(newHead, text)\n'
    '        if (!down && newHead < head && head < lineEndPos(head, text))\n'
    '            newHead = lineStartPos(newHead, text)\n'
    '        if (!down && head === text.length && query.selectionStart === query.selectionEnd) {\n'
    '            var upOnce = lineUpPos(head, text)\n'
    '            anchor = lineStartPos(upOnce, text)\n'
    '            newHead = head\n'
    '        }\n'
    '        if (down && newHead === head && head < text.length\n'
    '                && text.indexOf("\\n", head) === -1 && onWrappedLine(head, text)) {\n'
    '            var vis = visualLineDownPos(head)\n'
    '            if (vis > head) newHead = vis\n'
    '        }\n'
    '        if (down && newHead === head && head === text.length && head > 0\n'
    '                && query.selectionStart === query.selectionEnd)\n'
    '            newHead = head - 1\n'
    '        query.select(Math.min(anchor, newHead), Math.max(anchor, newHead))\n'
    '    }\n'
    '\n'
    '    function moveCursorEndOfLine() {\n'
    '        moveCursorTo(lineEndPos(query.cursorPosition, query.text), false)\n'
    '    }\n'
    '\n'
    '    function moveCursorStartOfLine() {\n'
    '        moveCursorTo(lineStartPos(query.cursorPosition, query.text), false)\n'
    '    }\n'
    '\n'
    '    function handleMacUndo(event) {\n'
    '        if (mode != 1) return false\n'
    '        if (!(event.modifiers & Qt.ControlModifier)) return false\n'
    '        if (event.key === Qt.Key_Z && !(event.modifiers & Qt.ShiftModifier)) {\n'
    '            if (query.canUndo) {\n'
    '                var lenBefore = query.text.length\n'
    '                query.undo()\n'
    '                if (macUndoKind === "insert" && macUndoMarkPos >= 0) {\n'
    '                    pendingRedoCursor = macUndoMarkPos + (lenBefore - query.text.length)\n'
    '                    query.cursorPosition = macUndoMarkPos\n'
    '                } else if (macUndoKind === "delete" && macUndoMarkPos >= 0) {\n'
    '                    pendingRedoCursor = -1\n'
    '                    if (macUndoMarkPos === 0 && query.cursorPosition === 0\n'
    '                            && query.text.length > 0)\n'
    '                        query.cursorPosition = query.text.length\n'
    '                    else\n'
    '                        query.cursorPosition = macUndoMarkPos\n'
    '                } else if (query.cursorPosition > query.text.length)\n'
    '                    query.cursorPosition = query.text.length\n'
    '                macUndoKind = ""\n'
    '                macUndoMarkPos = -1\n'
    '            }\n'
    '            cursorStrong = true\n'
    '            cursorTimer.stop()\n'
    '            event.accepted = true\n'
    '            return true\n'
    '        }\n'
    '        if (event.key === Qt.Key_Y || (event.key === Qt.Key_Z && (event.modifiers & Qt.ShiftModifier))) {\n'
    '            if (query.canRedo) {\n'
    '                query.redo()\n'
    '                if (pendingRedoCursor >= 0 && pendingRedoCursor <= query.text.length)\n'
    '                    query.cursorPosition = pendingRedoCursor\n'
    '                else if (query.cursorPosition > query.text.length)\n'
    '                    query.cursorPosition = query.text.length\n'
    '                pendingRedoCursor = -1\n'
    '            }\n'
    '            cursorStrong = true\n'
    '            cursorTimer.stop()\n'
    '            event.accepted = true\n'
    '            return true\n'
    '        }\n'
    '        return false\n'
    '    }\n'
    '\n'
    '    function handleMacArrow(event) {\n'
    '        if (mode != 1) return false\n'
    '        var mods = event.modifiers\n'
    '        if (mods === Qt.NoModifier) {\n'
    '            if (event.key === Qt.Key_Left || event.key === Qt.Key_Right) {\n'
    '                if (query.selectionStart !== query.selectionEnd) {\n'
    '                    var c = (event.key === Qt.Key_Left)\n'
    '                        ? Math.min(query.selectionStart, query.selectionEnd)\n'
    '                        : Math.max(query.selectionStart, query.selectionEnd)\n'
    '                    query.deselect()\n'
    '                    query.cursorPosition = c\n'
    '                    cursorStrong = true\n'
    '                    cursorTimer.stop()\n'
    '                    event.accepted = true\n'
    '                    return true\n'
    '                }\n'
    '                return false\n'
    '            }\n'
    '        }\n'
    '        var shift = mods & Qt.ShiftModifier\n'
    '        var cmd = mods & Qt.ControlModifier\n'
    '        var alt = mods & Qt.AltModifier\n'
    '        var text = query.text\n'
    '        var pos = query.cursorPosition\n'
    '        var newPos = pos\n'
    '        if ((event.key === Qt.Key_Home || event.key === Qt.Key_End) && !(shift && !cmd && !alt)) {\n'
    '            if (alt) return false\n'
    '            newPos = (event.key === Qt.Key_Home)\n'
    '                ? (cmd ? 0 : macLineStartPos(pos, text))\n'
    '                : (cmd ? text.length : macLineEndPos(pos, text))\n'
    '        } else if (shift && !cmd && !alt) {\n'
    '            if (event.key === Qt.Key_Left) {\n'
    '                var anchor = (query.selectionStart === query.selectionEnd)\n'
    '                    ? pos : Math.max(query.selectionStart, query.selectionEnd)\n'
    '                var head = (query.selectionStart === query.selectionEnd)\n'
    '                    ? pos : Math.min(query.selectionStart, query.selectionEnd)\n'
    '                newPos = Math.max(0, head - 1)\n'
    '                cursorStrong = true\n'
    '                cursorTimer.stop()\n'
    '                query.select(newPos, anchor)\n'
    '                event.accepted = true\n'
    '                return true\n'
    '            } else if (event.key === Qt.Key_Right) {\n'
    '                var anchorR = (query.selectionStart === query.selectionEnd)\n'
    '                    ? pos : Math.min(query.selectionStart, query.selectionEnd)\n'
    '                var headR = (query.selectionStart === query.selectionEnd)\n'
    '                    ? pos : Math.max(query.selectionStart, query.selectionEnd)\n'
    '                newPos = Math.min(text.length, headR + 1)\n'
    '                cursorStrong = true\n'
    '                cursorTimer.stop()\n'
    '                query.select(anchorR, newPos)\n'
    '                event.accepted = true\n'
    '                return true\n'
    '            } else if (event.key === Qt.Key_Up || event.key === Qt.Key_Down) {\n'
    '                cursorStrong = true\n'
    '                cursorTimer.stop()\n'
    '                extendSelectionVertical(event.key === Qt.Key_Down)\n'
    '                event.accepted = true\n'
    '                return true\n'
    '            } else if (event.key === Qt.Key_Home) {\n'
    '                var anchorH = (query.selectionStart === query.selectionEnd)\n'
    '                    ? pos : Math.max(query.selectionStart, query.selectionEnd)\n'
    '                newPos = macLineStartPos(pos, text)\n'
    '                cursorStrong = true\n'
    '                cursorTimer.stop()\n'
    '                query.select(newPos, anchorH)\n'
    '                event.accepted = true\n'
    '                return true\n'
    '            } else if (event.key === Qt.Key_End) {\n'
    '                var anchorE = (query.selectionStart === query.selectionEnd)\n'
    '                    ? pos : Math.min(query.selectionStart, query.selectionEnd)\n'
    '                newPos = macLineEndPos(pos, text)\n'
    '                cursorStrong = true\n'
    '                cursorTimer.stop()\n'
    '                query.select(anchorE, newPos)\n'
    '                event.accepted = true\n'
    '                return true\n'
    '            } else {\n'
    '                return false\n'
    '            }\n'
    '        } else if (!cmd && !alt && !shift) {\n'
    '            if (event.key === Qt.Key_Up || event.key === Qt.Key_Down) {\n'
    '                cursorStrong = true\n'
    '                cursorTimer.stop()\n'
    '                moveCursorVertical(event.key === Qt.Key_Down)\n'
    '                event.accepted = true\n'
    '                return true\n'
    '            }\n'
    '            if (event.key === Qt.Key_Left || event.key === Qt.Key_Right) {\n'
    '                if (query.selectionStart !== query.selectionEnd) {\n'
    '                    var selLo = Math.min(query.selectionStart, query.selectionEnd)\n'
    '                    var selHi = Math.max(query.selectionStart, query.selectionEnd)\n'
    '                    newPos = (event.key === Qt.Key_Left) ? selLo : selHi\n'
    '                    moveCursorTo(newPos, false, false)\n'
    '                    cursorStrong = true\n'
    '                    cursorTimer.stop()\n'
    '                    event.accepted = true\n'
    '                    return true\n'
    '                }\n'
    '            }\n'
    '            return false\n'
    '        } else if (event.key === Qt.Key_Left) {\n'
    '            if (cmd && shift) newPos = 0\n'
    '            else if (cmd) newPos = lineStartPos(pos, text)\n'
    '            else if (alt) newPos = wordLeftPos(pos, text)\n'
    '            else newPos = lineStartPos(pos, text)\n'
    '        } else if (event.key === Qt.Key_Right) {\n'
    '            if (cmd && shift) newPos = text.length\n'
    '            else if (cmd) newPos = lineEndPos(pos, text)\n'
    '            else if (alt) newPos = wordRightPos(pos, text)\n'
    '            else newPos = lineEndPos(pos, text)\n'
    '        } else if (event.key === Qt.Key_Up) {\n'
    '            if (cmd) newPos = 0\n'
    '            else newPos = paragraphUpPos(pos, text)\n'
    '        } else if (event.key === Qt.Key_Down) {\n'
    '            if (cmd) newPos = text.length\n'
    '            else newPos = paragraphDownPos(pos, text)\n'
    '        } else {\n'
    '            return false\n'
    '        }\n'
    '        cursorStrong = true\n'
    '        cursorTimer.stop()\n'
    '        if (shift)\n'
    '            extendSelectionHorizontal(newPos)\n'
    '        else\n'
    '            moveCursorTo(newPos, false)\n'
    '        event.accepted = true\n'
    '        return true\n'
    '    }\n'
    '\n'
    '    function handleMacEditKeys(event) {\n'
    '        if (mode != 1) return false\n'
    '        var mods = event.modifiers\n'
    '        var cmd = mods & Qt.ControlModifier\n'
    '        if (cmd && !(mods & Qt.AltModifier) && event.key === Qt.Key_A) {\n'
    '            var len = query.text.length\n'
    '            query.select(0, len)\n'
    '            cursorStrong = true\n'
    '            cursorTimer.stop()\n'
    '            event.accepted = true\n'
    '            doc = query.text\n'
    '            return true\n'
    '        }\n'
    '        if (mods !== Qt.NoModifier && !(mods === Qt.ShiftModifier && event.key === Qt.Key_Backspace)) return false\n'
    '        var text = query.text\n'
    '        var pos = query.cursorPosition\n'
    '        if (event.key === Qt.Key_Backspace) {\n'
    '            if (query.selectionStart !== query.selectionEnd) {\n'
    '                var a = Math.min(query.selectionStart, query.selectionEnd)\n'
    '                var b = Math.max(query.selectionStart, query.selectionEnd)\n'
    '                query.cursorPosition = b\n'
    '                query.select(a, b)\n'
    '                query.insert("")\n'
    '                query.cursorPosition = a\n'
    '                query.deselect()\n'
    '                macUndoKind = "delete"\n'
    '                macUndoMarkPos = a\n'
    '            } else if (pos > 0) {\n'
    '                query.select(pos - 1, pos)\n'
    '                query.insert("")\n'
    '                query.cursorPosition = pos - 1\n'
    '                macUndoKind = "delete"\n'
    '                macUndoMarkPos = query.cursorPosition\n'
    '            }\n'
    '            cursorStrong = true\n'
    '            cursorTimer.stop()\n'
    '            event.accepted = true\n'
    '            doc = query.text\n'
    '            return true\n'
    '        }\n'
    '        if (event.key === Qt.Key_Delete) {\n'
    '            if (query.selectionStart !== query.selectionEnd) {\n'
    '                var sa = Math.min(query.selectionStart, query.selectionEnd)\n'
    '                var sb = Math.max(query.selectionStart, query.selectionEnd)\n'
    '                query.select(sa, sb)\n'
    '                query.insert("")\n'
    '                query.cursorPosition = sa\n'
    '                query.deselect()\n'
    '                macUndoKind = "delete"\n'
    '                macUndoMarkPos = sa\n'
    '            } else if (pos < text.length) {\n'
    '                query.select(pos, pos + 1)\n'
    '                query.insert("")\n'
    '                macUndoKind = "delete"\n'
    '                macUndoMarkPos = pos\n'
    '            }\n'
    '            cursorStrong = true\n'
    '            cursorTimer.stop()\n'
    '            event.accepted = true\n'
    '            doc = query.text\n'
    '            return true\n'
    '        }\n'
    '        if (event.key === Qt.Key_Return) {\n'
    '            query.text = text.slice(0, pos) + "\\n" + text.slice(pos)\n'
    '            query.cursorPosition = pos + 1\n'
    '            query.deselect()\n'
    '            cursorStrong = true\n'
    '            cursorTimer.stop()\n'
    '            event.accepted = true\n'
    '            doc = query.text\n'
    '            return true\n'
    '        }\n'
    '        return false\n'
    '    }\n'
    '\n'
    '    function handleMacBackspace(event) {\n'
    '        if (mode != 1) return false\n'
    '        var mods = event.modifiers\n'
    '        var cmd = mods & Qt.ControlModifier\n'
    '        var alt = mods & Qt.AltModifier\n'
    '        if (!cmd && !alt) return false\n'
    '        var text = query.text\n'
    '        var pos = query.cursorPosition\n'
    '        if (pos <= 0 && text.length === 0) {\n'
    '            event.accepted = true\n'
    '            return true\n'
    '        }\n'
    '        if (query.selectionStart !== query.selectionEnd) {\n'
    '            var selA = Math.min(query.selectionStart, query.selectionEnd)\n'
    '            var selB = Math.max(query.selectionStart, query.selectionEnd)\n'
    '            query.text = text.slice(0, selA) + text.slice(selB)\n'
    '            query.cursorPosition = selA\n'
    '            query.deselect()\n'
    '            cursorStrong = true\n'
    '            cursorTimer.stop()\n'
    '            event.accepted = true\n'
    '            return true\n'
    '        }\n'
    '        var start = cmd ? deleteLineLeftPos(pos, text) : deleteWordLeftPos(pos, text)\n'
    '        var end = pos\n'
    '        if (alt && !cmd) {\n'
    '            while (end < text.length && !isSpaceChar(text.charAt(end))) end++\n'
    '        }\n'
    '        if (start < end) {\n'
    '            query.text = text.slice(0, start) + text.slice(end)\n'
    '            query.cursorPosition = start\n'
    '            macUndoKind = "delete"\n'
    '            macUndoMarkPos = start\n'
    '            query.deselect()\n'
    '        }\n'
    '        cursorStrong = true\n'
    '        cursorTimer.stop()\n'
    '        event.accepted = true\n'
    '        return true\n'
    '    }\n'
    '\n'
    '    function showLobby() {'
)
assert old7n_fn in s, "function showLobby not found (7n)"
s = s.replace(old7n_fn, new7n_fn, 1)

# 7o. Mac-style modifier+arrow in edit mode; plain Left/Right keep page-scroll.
s, n7o_kp = re.subn(
    r'([ \t]+Keys\.onPressed:\s*\{\s*\n)([ \t]+)switch\(event\.key\)\{',
    lambda m: (
        m.group(1)
        + m.group(2) + 'if (mode == 1) {\n'
        + m.group(2) + '    if (root.handleMacEditKeys(event))\n'
        + m.group(2) + '        return\n'
        + m.group(2) + '    if (root.handleMacUndo(event))\n'
        + m.group(2) + '        return\n'
        + m.group(2) + '    if (root.handleMacBackspace(event))\n'
        + m.group(2) + '        return\n'
        + m.group(2) + '    if (root.handleMacArrow(event))\n'
        + m.group(2) + '        return\n'
        + m.group(2) + '}\n'
        + m.group(2) + 'switch(event.key){'
    ),
    s, count=1
)
assert n7o_kp == 1, "Keys.onPressed header not found (7o)"

old7o_l = (
    '                        case Qt.Key_Left:\n'
    '                            flick.scrollUp()\n'
    '                            break'
)
new7o_l = (
    '                        case Qt.Key_Left:\n'
    '                            if (event.modifiers === Qt.NoModifier) {\n'
    '                                if (mode == 1) { cursorStrong = true; cursorTimer.stop() }\n'
    '                                flick.scrollUp()\n'
    '                            }\n'
    '                            break'
)
assert old7o_l in s, "Key_Left scroll case not found (7o)"
s = s.replace(old7o_l, new7o_l, 1)

old7o_r = (
    '                        case Qt.Key_Right:\n'
    '                            flick.scrollDown()\n'
    '                            break'
)
new7o_r = (
    '                        case Qt.Key_Right:\n'
    '                            if (event.modifiers === Qt.NoModifier) {\n'
    '                                if (mode == 1) { cursorStrong = true; cursorTimer.stop() }\n'
    '                                flick.scrollDown()\n'
    '                            }\n'
    '                            break'
)
assert old7o_r in s, "Key_Right scroll case not found (7o)"
s = s.replace(old7o_r, new7o_r, 1)

old7o_rot_r = (
    '            case Qt.Key_Right:\n'
    '                if (ctrlPressed)\n'
    '                    root.rotation = (root.rotation + 90) % 360'
)
new7o_rot_r = (
    '            case Qt.Key_Right:\n'
    '                if (ctrlPressed || (event.modifiers & Qt.ControlModifier))\n'
    '                    root.rotation = (root.rotation + 90) % 360'
)
assert old7o_rot_r in s, "preview Ctrl+Right rotate not found (7o)"
s = s.replace(old7o_rot_r, new7o_rot_r, 1)

old7o_rot_l = (
    '            case Qt.Key_Left:\n'
    '                if (ctrlPressed)\n'
    '                    root.rotation = (root.rotation - 90) % 360'
)
new7o_rot_l = (
    '            case Qt.Key_Left:\n'
    '                if (ctrlPressed || (event.modifiers & Qt.ControlModifier))\n'
    '                    root.rotation = (root.rotation - 90) % 360'
)
assert old7o_rot_l in s, "preview Ctrl+Left rotate not found (7o)"
s = s.replace(old7o_rot_l, new7o_rot_l, 1)

# 7p. Lobby rotation: same Ctrl+arrow as preview (Lobby boots in edit mode).
old7p = (
    '        if (mode == 0)\n'
    '            switch (event.key) {\n'
    '            case Qt.Key_Home:\n'
    '                Qt.quit()\n'
    '                break\n'
    '            case Qt.Key_Right:\n'
    '                if (ctrlPressed || (event.modifiers & Qt.ControlModifier))\n'
    '                    root.rotation = (root.rotation + 90) % 360\n'
    '                break\n'
    '            case Qt.Key_Left:\n'
    '                if (ctrlPressed || (event.modifiers & Qt.ControlModifier))\n'
    '                    root.rotation = (root.rotation - 90) % 360\n'
    '                break\n'
    '            }'
)
new7p = (
    '        if (mode == 0 || isLobby) {\n'
    '            switch (event.key) {\n'
    '            case Qt.Key_Right:\n'
    '                if (ctrlPressed || (event.modifiers & Qt.ControlModifier))\n'
    '                    root.rotation = (root.rotation + 90) % 360\n'
    '                break\n'
    '            case Qt.Key_Left:\n'
    '                if (ctrlPressed || (event.modifiers & Qt.ControlModifier))\n'
    '                    root.rotation = (root.rotation - 90) % 360\n'
    '                break\n'
    '            }\n'
    '        }'
)
assert old7p in s, "handleKey rotate block not found (7p)"
s = s.replace(old7p, new7p, 1)

# 7q. Lobby subpage navigation at the top of handleKey.
old7q = '    function handleKey(event) {'
new7q = (
    '    function handleKey(event) {\n'
    '        if (event.key === Qt.Key_Home && event.modifiers === Qt.NoModifier) {\n'
    '            if (suppressNextHomeKey) {\n'
    '                suppressNextHomeKey = false\n'
    '                event.accepted = true\n'
    '                return\n'
    '            }\n'
    '            // Edit mode: line start is handled on press (handleMacArrow); do not\n'
    '            // treat Key_Home release as physical Home -> lobby.\n'
    '            if (mode == 1 && !isLobby) {\n'
    '                event.accepted = true\n'
    '                return\n'
    '            }\n'
    '            handleHome(false)\n'
    '            event.accepted = true\n'
    '            return\n'
    '        }\n'
    '        if (suppressNextHomeKey)\n'
    '            suppressNextHomeKey = false\n'
    '        if (vaultOverlayMode !== "") {\n'
    '            if (vaultConsumeKey(event)) { event.accepted = true; return }\n'
    '        }\n'
    '        if (isLobby && !isOmni) {\n'
    '            if (lobbyHandleKey(event)) {\n'
    '                event.accepted = true\n'
    '                return\n'
    '            }\n'
    '            if (!(event.modifiers & Qt.ControlModifier)) {\n'
    '                event.accepted = true\n'
    '                return\n'
    '            }\n'
    '        }'
)
assert old7q in s, "handleKey function not found (7q)"
s = s.replace(old7q, new7q, 1)

# 7v. Esc release toggles edit/preview; ignore Alt/Ctrl+Esc (NO qmap maps Alt+arrows to fake Esc).
old7v = (
    '        if (event.key === Qt.Key_Escape) {\n'
    '            if (isOmni) {\n'
    '                isOmni = false\n'
    '            } else {\n'
    '\n'
    '                toggleMode()\n'
    '            }\n'
    '        }'
)
new7v = (
    '        if (event.key === Qt.Key_Escape) {\n'
    '            if (isOmni) {\n'
    '                isOmni = false\n'
    '            } else if (!(event.modifiers & (Qt.AltModifier | Qt.ControlModifier))) {\n'
    '                toggleMode()\n'
    '            }\n'
    '        }'
)
assert old7v in s, "handleKey Escape toggle not found (7v)"
s = s.replace(old7v, new7v, 1)

# 8. Add ListModel + Lobby subpages (concat lobby/*.inc) and sleep screen.
import subprocess
subprocess.run(['/concat-lobby.sh'], check=True)
with open('/lobby_subpages.qml.inc', 'r') as lf:
    lobby_ui = lf.read()
with open('/lobby/lobby_vault_numpad.inc', 'r') as vf:
    vault_ui = vf.read()
lobby_rect = (
    '        ListModel {\n'
    '            id: lobbyNotesModel\n'
    '        }\n'
    + lobby_ui +
    vault_ui +
    '        Rectangle {\n'
    '            id: sleepScreen\n'
    '            anchors.fill: parent\n'
    '            color: "white"\n'
    '            visible: isSleeping\n'
    '            z: 10\n'
    '            Column {\n'
    '                anchors.centerIn: parent\n'
    '                width: sleepScreen.width * 0.75\n'
    '                spacing: 24\n'
    '                Text {\n'
    '                    text: "Writerdeck is sleeping.\\nWi-Fi is off. Press power to wake."\n'
    '                    color: "black"\n'
    '                    font.pointSize: 18\n'
    '                    font.family: "Noto Sans"\n'
    '                    width: parent.width\n'
    '                    wrapMode: Text.WordWrap\n'
    '                    horizontalAlignment: Text.AlignHCenter\n'
    '                }\n'
    '            }\n'
    '        }'
)
quick_close = '        }\n'
end_anchor = quick_close + '    }\n}'
assert end_anchor in s, "QML end structure (quick+body+Window close) not found"
last_pos = s.rfind(end_anchor)
s = s[:last_pos + len(quick_close)] + lobby_rect + '\n' + s[last_pos + len(quick_close):]

# 8b. Cursor-state Timer and Connections: siblings of the Lobby Rectangle,
#     inside the body Rectangle. Timer (500 ms, one-shot) sets cursorStrong=true
#     when typing stops. Connections fires on every text change to hide the block
#     and restart the timer.
body_end = '    }\n}'
assert body_end in s, "body+Window end not found (8b)"
last_body_pos = s.rfind(body_end)
cursor_state_block = (
    '        Timer {\n'
    '            id: autosaveTimer\n'
    '            interval: 45000\n'
    '            repeat: true\n'
    '            running: !isLobby && currentFile !== "" && mode == 1\n'
    '            onTriggered: autosaveTick()\n'
    '        }\n'
    '        Timer {\n'
    '            id: cursorTimer\n'
    '            interval: 500\n'
    '            repeat: false\n'
    '            onTriggered: cursorStrong = true\n'
    '        }\n'
    '        Connections {\n'
    '            target: query\n'
    '            onTextChanged: {\n'
    '                var len = query.text.length\n'
    '                if (len > _prevTextLen) {\n'
    '                    if (macUndoKind !== "insert")\n'
    '                        macUndoMarkPos = _prevCursor\n'
    '                    macUndoKind = "insert"\n'
    '                }\n'
    '                _prevTextLen = len\n'
    '                _prevCursor = query.cursorPosition\n'
    '                cursorStrong = false\n'
    '                cursorTimer.restart()\n'
    '            }\n'
    '        }\n'
)
s = s[:last_body_pos] + cursor_state_block + s[last_body_pos:]

# rotateScreen(): rotate the display 90 degrees clockwise. Called by the C++
# setrotation cmd: server pushes saved angle on connect or after USB rotation ack.
# rotate cmd (legacy): bump root.rotation 90 CW if ever sent; phone POST /api/rotate removed.
# 90 mod 360; the body Rectangle already swaps its width/height at 90/270 degrees.
old_rotate_fn = '    function initFile(name) {'
new_rotate_fn = (
    '    function rotateScreen() {\n'
    '        root.rotation = (root.rotation + 90) % 360\n'
    '    }\n'
    '\n'
    '    function initFile(name) {'
)
assert old_rotate_fn in s, "function initFile not found (rotateScreen)"
s = s.replace(old_rotate_fn, new_rotate_fn, 1)

# 7u. Drop the TextEdit `text:` binding -- broken by typing and doLoad; syncQueryDisplay()
#     (edit 7f) sets query.text imperatively on load and in toggleMode (edit 1c).
assert 'text: mode == 0 ? root.readHtml(doc) : doc' in s, "text binding not found (edit 7u drop)"
s = s.replace('                text: mode == 0 ? root.readHtml(doc) : doc\n', '', 1)

# Sanity: handleKey must close before Component.onCompleted (patch 7p regressed once).
hk = s.find('    function handleKey(event) {')
co = s.find('    Component.onCompleted: {')
assert hk >= 0 and co > hk, "handleKey / Component.onCompleted anchors missing"
assert s[hk:co].count('{') == s[hk:co].count('}'), "handleKey brace mismatch -- QML will fail to load"

with open('main.qml', 'w') as f:
    f.write(s)
print('  All QML edits applied (props + content-fidelity + setLobbyInfo + lobby-subpages + handleHome + doLoad-query-sync + prepareSleep + sleep-screen + openNotePicker + omni-z + saveAndLoad + saveAndQuit + boot-edit-mode + Ctrl-K/Q/R + margin + block-cursor + scroll-dir + scroll-4/5 + page-btn-edit-scroll + read-no-autoscroll + cursor-boundary + mac-arrows-home-end + para-spacing-28 + list-spacing + readFont + setReadFont + noteDeleted + noteRenamed + reloadNote + autosave + saveFile-guard + scratch-demote + showLobby + no-PIN-lobby + cursor-hidden-when-typing + rotateScreen + lobby-rotate + lobbyHandleKey + syncQueryDisplay).')
PYEOF
echo "  main.qml after edit:"
grep -n 'property int mode:\|saveAndQuit\|ControlModifier' main.qml || true
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
