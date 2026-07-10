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
printf '\nHEADERS += rotation_watcher.h\nSOURCES += rotation_watcher.cpp\n' >> edit.pro
echo "  rotation_watcher added to edit.pro."

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

# 2. doLoad resets mode on every load -- keep it edit (mode=1).
#    Match the first 'mode = 0' on the line immediately after 'isOmni = false';
#    leaves toggleMode's own 'mode = 0' (preview+save) untouched.
s, n = re.subn(
    r'(isOmni\s*=\s*false[^\n]*\n[^\n]*)mode\s*=\s*0',
    lambda m: m.group(0).replace('mode = 0', 'mode = 1', 1),
    s, count=1
)
assert n == 1, "doLoad mode=0 pattern not found in main.qml"

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
    '                else openNotePicker()\n'
    '            } else isOmni = !isOmni\n'
    '            event.accepted = true\n'
    '        }'
)
assert old4b in s, "Ctrl-K handler block not found (edit 4b)"
s = s.replace(old4b, new4b, 1)

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
    '    property string lobbySyncRepo: ""'
)
assert old6 in s, "isOmni property not found in main.qml"
s = s.replace(old6, new6, 1)

# 7. Add setLobbyInfo() and handleHome() before initFile.
#    setLobbyInfo: called from C++ via invokeMethod when rmkbd connects, so the
#    Lobby shows the current IP and PIN without hardcoding anything.
#    handleHome: two-level Home (slice 8e):
#      - editing -> save current note + return to Lobby (isLobby = true)
#      - Lobby   -> Qt.quit() -> cmd.Wait fires -> s.end() -> xochitl restarts
old7 = '    function initFile(name) {'
new7 = (
    '    function setLobbyInfo(ip, pin, syncOn, syncRepo) {\n'
    '        lobbyIP = ip\n'
    '        lobbyPIN = pin\n'
    '        lobbySyncOn = !!syncOn\n'
    '        lobbySyncRepo = syncRepo || ""\n'
    '    }\n'
    '\n'
    '    function prepareSleep() {\n'
    '        if (mode == 1) doc = query.text\n'
    '        saveFile()\n'
    '        isLobby = false\n'
    '        isSleeping = true\n'
    '    }\n'
    '\n'
    '    function handleHome() {\n'
    '        if (isLobby) {\n'
    '            Qt.quit()\n'
    '        } else {\n'
    '            if (mode == 1) doc = query.text\n'
    '            saveFile()\n'
    '            isLobby = true\n'
    '            currentFile = ""\n'
    '            doc = ""\n'
    '            query.text = ""\n'
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
    '                    if (mode == 1)\n'
    '                        flick.ensureVisible(cursorRectangle)\n'
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
    '    function initFile(name) {'
)
assert old7f_fn in s, "function initFile not found (edit 7f)"
s = s.replace(old7f_fn, new7f_fn, 1)
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
    '        if (!isLobby) {\n'
    '            if (mode == 1) doc = query.text\n'
    '            saveFile()\n'
    '        }\n'
    '        isLobby = true\n'
    '    }\n'
    '\n'
    '    function initFile(name) {'
)
assert old7l in s, "function initFile not found (edit 7l)"
s = s.replace(old7l, new7l, 1)

# 7h. noteDeleted(): go to the Lobby and CLEAR currentFile so the next save
#     can't resurrect the just-deleted file. Called by the C++ notedeleted cmd
#     when rmkbd tells the editor its open file was deleted from the phone.
#     Crucially this does NOT call saveFile() -- and because saveAndLoad() saves
#     BEFORE it switches files, clearing currentFile + the saveFile guard (7i)
#     makes that pre-switch save a no-op, so opening another note won't recreate X.
old7h = '    function initFile(name) {'
new7h = (
    '    function noteDeleted() {\n'
    '        currentFile = ""\n'
    '        isLobby = true\n'
    '    }\n'
    '\n'
    '    function initFile(name) {'
)
assert old7h in s, "function initFile not found (edit 7h)"
s = s.replace(old7h, new7h, 1)

# 7i. Guard saveFile() against an empty currentFile. After noteDeleted() clears
#     currentFile, the next saveAndLoad() still calls saveFile() BEFORE doLoad()
#     -- without this guard that write targets the deleted file's path and
#     recreates it on disk (the resurrection trap). The guard also fires at boot
#     (currentFile now defaults to "") and after any Lobby return (handleHome
#     clears it, edit 7j), making every Lobby entry a clean no-file state.
old7i = '    function saveFile() {\n        console.log("Save " + currentFile)'
new7i = '    function saveFile() {\n        if (currentFile === "") return 0\n        console.log("Save " + currentFile)'
assert old7i in s, "saveFile signature not found (edit 7i)"
s = s.replace(old7i, new7i, 1)

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
    '    property bool cursorStrong: true',
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
    '                                if (root.cursorOnLastLine()) {\n'
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
    '                                if (root.cursorOnFirstLine()) {\n'
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
    '    function cursorOnLastLine() {\n'
    '        if (mode != 1) return false\n'
    '        return query.text.indexOf("\\n", query.cursorPosition) === -1\n'
    '    }\n'
    '\n'
    '    function cursorOnFirstLine() {\n'
    '        if (mode != 1) return false\n'
    '        var pos = query.cursorPosition\n'
    '        return pos === 0 || query.text.lastIndexOf("\\n", pos - 1) === -1\n'
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
    '    function wordLeftPos(pos, text) {\n'
    '        if (pos <= 0) return 0\n'
    '        pos--\n'
    '        while (pos > 0 && isSpaceChar(text.charAt(pos))) pos--\n'
    '        while (pos > 0 && !isSpaceChar(text.charAt(pos - 1))) pos--\n'
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
    '        var i = lineEnd + 1\n'
    '        while (i < len - 1) {\n'
    '            if (text.charAt(i) === "\\n" && text.charAt(i + 1) === "\\n")\n'
    '                return i + 2\n'
    '            i++\n'
    '        }\n'
    '        return len\n'
    '    }\n'
    '\n'
    '    function moveCursorTo(newPos, extend) {\n'
    '        var len = query.text.length\n'
    '        newPos = Math.max(0, Math.min(newPos, len))\n'
    '        if (!extend) {\n'
    '            query.cursorPosition = newPos\n'
    '            return\n'
    '        }\n'
    '        var anchor = query.cursorPosition\n'
    '        if (query.selectionStart !== query.selectionEnd) {\n'
    '            anchor = (query.cursorPosition === query.selectionEnd)\n'
    '                ? query.selectionStart : query.selectionEnd\n'
    '        }\n'
    '        query.select(Math.min(anchor, newPos), Math.max(anchor, newPos))\n'
    '        query.cursorPosition = newPos\n'
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
    '    function handleMacArrow(event) {\n'
    '        if (mode != 1) return false\n'
    '        var mods = event.modifiers\n'
    '        var shift = mods & Qt.ShiftModifier\n'
    '        var cmd = mods & Qt.ControlModifier\n'
    '        var alt = mods & Qt.AltModifier\n'
    '        var text = query.text\n'
    '        var pos = query.cursorPosition\n'
    '        var newPos = pos\n'
    '        if (event.key === Qt.Key_Home || event.key === Qt.Key_End) {\n'
    '            if (alt) return false\n'
    '            newPos = (event.key === Qt.Key_Home)\n'
    '                ? (cmd ? 0 : lineStartPos(pos, text))\n'
    '                : (cmd ? text.length : lineEndPos(pos, text))\n'
    '        } else if (!cmd && !alt) {\n'
    '            return false\n'
    '        } else if (event.key === Qt.Key_Left) {\n'
    '            if (alt) newPos = wordLeftPos(pos, text)\n'
    '            else newPos = lineStartPos(pos, text)\n'
    '        } else if (event.key === Qt.Key_Right) {\n'
    '            if (alt) newPos = wordRightPos(pos, text)\n'
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
    '        moveCursorTo(newPos, !!shift)\n'
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
        + m.group(2) + '    if (root.handleMacArrow(event))\n'
        + m.group(2) + '        return\n'
        + m.group(2) + '    if (event.modifiers === Qt.ShiftModifier) {\n'
        + m.group(2) + '        cursorStrong = true; cursorTimer.stop()\n'
        + m.group(2) + '        event.accepted = false\n'
        + m.group(2) + '        return\n'
        + m.group(2) + '    }\n'
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

# 8. Add the Lobby Rectangle at the end of the body Rectangle, after the
#    quick (isOmni) overlay.  Anchor: the last "        }\n    }\n}" in the file
#    = quick close (8 sp) + body close (4 sp) + Window close (0 sp).
#    We insert the Lobby between quick's } and body's }.
lobby_rect = (
    '        Rectangle {\n'
    '            id: lobby\n'
    '            anchors.fill: parent\n'
    '            color: "white"\n'
    '            visible: isLobby\n'
    '            Column {\n'
    '                anchors.centerIn: parent\n'
    '                spacing: 28\n'
    '                Text {\n'
    '                    text: "Writerdeck for reMarkable 1"\n'
    '                    color: "black"\n'
    '                    font.pointSize: 28\n'
    '                    font.family: "Noto Mono"\n'
    '                }\n'
    '                Text {\n'
    '                    text: "A text editor for use with a physical keyboard. Markdown support.\\n\\nFor Bluetooth keyboards: Connect the keyboard to your phone, from there to reMarkable over Wi-Fi.\\n\\nFor USB keyboards: Connect using an OTG cable.\\n\\nOpen in your phone' + chr(39) + 's browser:"\n'
    '                    color: "#555555"\n'
    '                    font.pointSize: 13\n'
    '                    font.family: "Noto Sans"\n'
    '                    width: lobby.width * 0.8\n'
    '                    wrapMode: Text.WordWrap\n'
    '                    horizontalAlignment: Text.AlignLeft\n'
    '                }\n'
    '                Text {\n'
    '                    text: "http://" + lobbyIP + ":8000"\n'
    '                    color: "black"\n'
    '                    font.pointSize: 20\n'
    '                    font.family: "Noto Mono"\n'
    '                }\n'
    '                Text {\n'
    '                    text: lobbyPIN !== "" ? ("PIN:  " + lobbyPIN) : "No PIN needed \u2014 just open the address above"\n'
    '                    color: lobbyPIN !== "" ? "#1b5e20" : "#555555"\n'
    '                    font.pointSize: lobbyPIN !== "" ? 24 : 16\n'
    '                    font.family: "Noto Mono"\n'
    '                    width: lobby.width * 0.8\n'
    '                    wrapMode: Text.WordWrap\n'
    '                }\n'
    '                Rectangle {\n'
    '                    width: Math.min(lobby.width * 0.55, 420)\n'
    '                    height: 54\n'
    '                    color: "#f0f0f0"\n'
    '                    border.color: "#333333"\n'
    '                    border.width: 2\n'
    '                    radius: 4\n'
    '                    Text {\n'
    '                        anchors.centerIn: parent\n'
    '                        text: "Open note\\u2026  (Ctrl-K)"\n'
    '                        color: "black"\n'
    '                        font.pointSize: 14\n'
    '                        font.family: "Noto Sans"\n'
    '                    }\n'
    '                    MouseArea {\n'
    '                        anchors.fill: parent\n'
    '                        onClicked: openNotePicker()\n'
    '                    }\n'
    '                }\n'
    '                Text {\n'
    '                    text: "Home = exit to reMarkable UI"\n'
    '                    color: "#555555"\n'
    '                    font.pointSize: 11\n'
    '                    font.family: "Noto Sans"\n'
    '                }\n'
    '                Text {\n'
    '                    visible: lobbySyncOn && lobbySyncRepo !== ""\n'
    '                    text: "Your docs are synced to: github.com/" + lobbySyncRepo\n'
    '                    color: "#1b5e20"\n'
    '                    font.pointSize: 11\n'
    '                    font.family: "Noto Mono"\n'
    '                    width: lobby.width * 0.8\n'
    '                    wrapMode: Text.WordWrap\n'
    '                }\n'
    '                Text {\n'
    '                    text: "github.com/bjornte/Writerdeck-for-reMarkable"\n'
    '                    color: "#444444"\n'
    '                    font.pointSize: 9\n'
    '                    font.family: "Noto Mono"\n'
    '                }\n'
    '            }\n'
    '        }\n'
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
    '            id: cursorTimer\n'
    '            interval: 500\n'
    '            repeat: false\n'
    '            onTriggered: cursorStrong = true\n'
    '        }\n'
    '        Connections {\n'
    '            target: query\n'
    '            onTextChanged: {\n'
    '                cursorStrong = false\n'
    '                cursorTimer.restart()\n'
    '            }\n'
    '        }\n'
)
s = s[:last_body_pos] + cursor_state_block + s[last_body_pos:]

# rotateScreen(): rotate the display 90 degrees clockwise. Called by the C++
# rotate cmd sent by rmkbd on POST /api/rotate (C++ sets root.rotation directly;
# this function remains for symmetry / manual calls). Increments root.rotation by
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

with open('main.qml', 'w') as f:
    f.write(s)
print('  All QML edits applied (props + setLobbyInfo + handleHome + prepareSleep + sleep-screen + Lobby rect + openNotePicker + omni-z + saveAndLoad + saveAndQuit + boot-edit-mode + Ctrl-K/Q + margin + block-cursor + scroll-dir + scroll-4/5 + page-btn-edit-scroll + read-no-autoscroll + cursor-boundary + mac-arrows-home-end + para-spacing-28 + list-spacing + readFont + setReadFont + noteDeleted + saveFile-guard + scratch-demote + showLobby + no-PIN-lobby + cursor-hidden-when-typing + rotateScreen).')
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
