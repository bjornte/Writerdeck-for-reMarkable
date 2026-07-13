        Rectangle {
            id: lobby
            anchors.fill: parent
            color: "white"
            visible: isLobby
            z: 5

            readonly property int pageMargin: 24
            readonly property int tabBtnHeight: 64
            readonly property int rowHeight: 72
            readonly property int actionBtnHeight: 72
            readonly property int tabSpacing: 12
            readonly property int contentSpacing: 12

            FocusScope {
                id: lobbyFocus
                anchors.fill: parent
                focus: isLobby && !isOmni
                Keys.enabled: isLobby && !isOmni
                Keys.onPressed: handleKeyDown(event)
                Keys.onReleased: {
                    handleKeyUp(event)
                    handleKey(event)
                }

                Connections {
                    target: root
                    function onIsLobbyChanged() {
                        if (isLobby) Qt.callLater(function() { lobbyFocus.forceActiveFocus() })
                    }
                }

                // ---- tab bar (touch + keyboard 1-6 / Tab / arrows) ----
                Row {
                    id: lobbyTabRow
                    anchors.top: parent.top
                    anchors.left: parent.left
                    anchors.right: parent.right
                    anchors.leftMargin: lobby.pageMargin
                    anchors.rightMargin: lobby.pageMargin
                    anchors.topMargin: lobby.pageMargin
                    height: lobby.tabBtnHeight + 8
                    spacing: lobby.tabSpacing

                    Repeater {
                        model: lobbyTabLabels
                        delegate: Rectangle {
                            width: Math.max(88, (lobby.width - lobby.pageMargin * 2 - lobby.tabSpacing * (lobbyTabLabels.length - 1)) / lobbyTabLabels.length)
                            height: lobby.tabBtnHeight
                            radius: 6
                            color: lobbyPage === index ? "#e0e0e0" : "#f5f5f5"
                            border.color: lobbyPage === index ? "#999" : "#ccc"
                            border.width: 1

                            Text {
                                anchors.centerIn: parent
                                text: (index + 1) + " " + modelData
                                font.family: "Noto Sans"
                                font.pointSize: 11
                                color: "black"
                            }

                            MouseArea {
                                anchors.fill: parent
                                onClicked: root.lobbyGoPage(index)
                            }
                        }
                    }
                }

                // ---- page stack ----
                Item {
                    anchors.top: lobbyTabRow.bottom
                    anchors.topMargin: lobby.contentSpacing
                    anchors.left: parent.left
                    anchors.right: parent.right
                    anchors.bottom: lobbyHint.top
                    anchors.bottomMargin: lobby.contentSpacing
                    anchors.leftMargin: lobby.pageMargin
                    anchors.rightMargin: lobby.pageMargin

                    // 0 Home
                    Item {
                        visible: lobbyPage === 0
                        anchors.fill: parent
                        Flickable {
                            anchors.fill: parent
                            contentWidth: width
                            contentHeight: homeCol.height
                            clip: true
                            Column {
                                id: homeCol
                                width: parent.width
                                spacing: lobby.contentSpacing
                                Text {
                                    text: "Writerdeck"
                                    color: "black"
                                    font.pointSize: 26
                                    font.family: "Noto Mono"
                                    width: parent.width
                                }
                                Text {
                                    text: "A text editor for use with a physical keyboard.\nWith Markdown support."
                                    color: "#555555"
                                    font.pointSize: 12
                                    font.family: "Noto Sans"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
                                }
                                Text {
                                    text: (lobbyNoteCount === 1 ? "1 note" : lobbyNoteCount + " notes") + " on this device."
                                    color: "#555555"
                                    font.pointSize: 11
                                    font.family: "Noto Sans"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
                                }
                                Text {
                                    text: "Open the Files tab (2) or press Ctrl-K.\nUse Tab / arrows / 1-6 to switch pages."
                                    color: "#888888"
                                    font.pointSize: 10
                                    font.family: "Noto Sans"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
                                }
                                Text {
                                    text: "Open sourced at github.com/bjornte/Writerdeck-for-reMarkable"
                                    color: "#888888"
                                    font.pointSize: 9
                                    font.family: "Noto Mono"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
                                }
                            }
                        }
                    }
                    // 1 Files
                    Item {
                        visible: lobbyPage === 1
                        anchors.fill: parent

                        Column {
                            anchors.fill: parent
                            spacing: lobby.contentSpacing

                            Text {
                                text: lobbyFilesMode === "new" ? "New note name:"
                                     : lobbyFilesMode === "rename" ? "Rename to:"
                                     : lobbyFilesMode === "confirm-delete" ? "Delete this note? Enter=yes  Esc=no"
                                     : "Notes"
                                font.family: "Noto Sans"
                                font.pointSize: 13
                                color: "black"
                                width: parent.width
                            }

                            Text {
                                visible: lobbyFilesMode === "new" || lobbyFilesMode === "rename"
                                text: lobbyFilesInput + "_"
                                font.family: "Noto Mono"
                                font.pointSize: 12
                                color: "#333"
                                width: parent.width
                            }

                            ListView {
                                id: lobbyFilesList
                                width: parent.width
                                height: {
                                    var reserved = lobbyFilesBar.height + lobby.contentSpacing * 2
                                    if (lobbyEncryptionEnabled && lobbyFilesMode === "" && lobbyNotesModel.count > 0)
                                        reserved += lobbyFilesVaultBar.height + lobby.contentSpacing
                                    return Math.max(lobby.rowHeight * 2, parent.height - reserved)
                                }
                                clip: true
                                model: lobbyNotesModel
                                currentIndex: lobbyFilesIndex
                                spacing: 4
                                visible: lobbyFilesMode === "" || lobbyFilesMode === "confirm-delete"
                                delegate: Rectangle {
                                    width: lobbyFilesList.width
                                    height: lobby.rowHeight
                                    color: index === lobbyFilesIndex ? "#e8e8e8" : "white"
                                    border.color: "#ddd"
                                    border.width: 1
                                    radius: 4
                                    Text {
                                        anchors.verticalCenter: parent.verticalCenter
                                        anchors.left: parent.left
                                        anchors.leftMargin: 16
                                        text: (model.encrypted ? "[lock] " : "") + model.name
                                        font.family: "Noto Mono"
                                        font.pointSize: 11
                                        color: "black"
                                        elide: Text.ElideRight
                                        width: parent.width - 32
                                    }
                                    MouseArea {
                                        anchors.fill: parent
                                        onClicked: {
                                            if (index === lobbyFilesIndex)
                                                root.lobbyOpenSelected()
                                            else {
                                                lobbyFilesIndex = index
                                                lobbyFilesList.currentIndex = index
                                            }
                                        }
                                        onDoubleClicked: root.lobbyOpenSelected()
                                    }
                                }
                                onCurrentIndexChanged: lobbyFilesIndex = currentIndex
                                highlight: Rectangle { color: "#d0d0d0"; radius: 4 }
                            }

                            Row {
                                id: lobbyFilesBar
                                width: parent.width
                                height: lobby.actionBtnHeight + 8
                                spacing: lobby.tabSpacing
                                visible: lobbyFilesMode === ""

                                Repeater {
                                    model: ["New", "Edit", "Read", "Rename", "Delete"]
                                    delegate: Rectangle {
                                        width: (lobbyFilesBar.width - lobby.tabSpacing * 4) / 5
                                        height: lobby.actionBtnHeight
                                        radius: 6
                                        color: "#f0f0f0"
                                        border.color: "#bbb"
                                        border.width: 1
                                        Text {
                                            anchors.centerIn: parent
                                            text: modelData
                                            font.family: "Noto Sans"
                                            font.pointSize: 11
                                        }
                                        MouseArea {
                                            anchors.fill: parent
                                            onClicked: {
                                                if (modelData === "New") root.lobbyFilesBeginNew()
                                                else if (modelData === "Edit") root.lobbyOpenSelected()
                                                else if (modelData === "Read") root.lobbyReadSelected()
                                                else if (modelData === "Rename") root.lobbyFilesBeginRename()
                                                else if (modelData === "Delete") root.lobbyFilesBeginDelete()
                                            }
                                        }
                                    }
                                }
                            }

                            Row {
                                id: lobbyFilesVaultBar
                                width: parent.width
                                height: lobby.actionBtnHeight + 4
                                spacing: lobby.tabSpacing
                                visible: lobbyFilesMode === "" && lobbyEncryptionEnabled && lobbyNotesModel.count > 0

                                property bool selectedEncrypted: {
                                    if (lobbyNotesModel.count === 0) return false
                                    var row = lobbyNotesModel.get(lobbyFilesIndex)
                                    return row && row.encrypted
                                }

                                Rectangle {
                                    visible: !lobbyFilesVaultBar.selectedEncrypted
                                    width: (lobbyFilesVaultBar.width - lobby.tabSpacing) / 2
                                    height: lobby.actionBtnHeight
                                    radius: 6
                                    color: "#f0f0f0"
                                    border.color: "#bbb"
                                    border.width: 1
                                    Text {
                                        anchors.centerIn: parent
                                        text: "Encrypt"
                                        font.family: "Noto Sans"
                                        font.pointSize: 11
                                    }
                                    MouseArea {
                                        anchors.fill: parent
                                        onClicked: root.lobbyEncryptSelected()
                                    }
                                }
                                Rectangle {
                                    visible: !lobbyFilesVaultBar.selectedEncrypted
                                    width: (lobbyFilesVaultBar.width - lobby.tabSpacing) / 2
                                    height: lobby.actionBtnHeight
                                    radius: 6
                                    color: "#f0f0f0"
                                    border.color: "#bbb"
                                    border.width: 1
                                    Text {
                                        anchors.centerIn: parent
                                        text: "New encrypted"
                                        font.family: "Noto Sans"
                                        font.pointSize: 11
                                    }
                                    MouseArea {
                                        anchors.fill: parent
                                        onClicked: root.lobbyFilesBeginNewEncrypted()
                                    }
                                }
                                Rectangle {
                                    visible: lobbyFilesVaultBar.selectedEncrypted
                                    width: lobbyFilesVaultBar.width
                                    height: lobby.actionBtnHeight
                                    radius: 6
                                    color: "#f0f0f0"
                                    border.color: "#bbb"
                                    border.width: 1
                                    Text {
                                        anchors.centerIn: parent
                                        text: "Decrypt"
                                        font.family: "Noto Sans"
                                        font.pointSize: 11
                                    }
                                    MouseArea {
                                        anchors.fill: parent
                                        onClicked: root.lobbyDecryptSelected()
                                    }
                                }
                            }
                        }
                    }
                    // 2 Keyboard
                    Item {
                        visible: lobbyPage === 2
                        anchors.fill: parent
                        Flickable {
                            anchors.fill: parent
                            contentWidth: width
                            contentHeight: kbCol.height
                            clip: true
                            Column {
                                id: kbCol
                                width: parent.width
                                spacing: lobby.contentSpacing
                                Text {
                                    text: "USB keyboard"
                                    font.pointSize: 14
                                    font.family: "Noto Sans"
                                    color: "black"
                                    width: parent.width
                                }
                                Text {
                                    text: "Connect with a USB OTG cable.\nChanging layout restarts Writerdeck."
                                    font.pointSize: 11
                                    font.family: "Noto Sans"
                                    color: "#555555"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
                                }
                                Text {
                                    text: "Layout"
                                    font.pointSize: 12
                                    font.family: "Noto Sans"
                                    color: "black"
                                    width: parent.width
                                }
                                Row {
                                    id: kbLayoutRow
                                    width: parent.width
                                    spacing: lobby.tabSpacing
                                    property var layoutOptions: [
                                        { id: "us", label: "US QWERTY" },
                                        { id: "no", label: "Norwegian" }
                                    ]
                                    Repeater {
                                        model: kbLayoutRow.layoutOptions
                                        delegate: Rectangle {
                                            width: (kbLayoutRow.width - lobby.tabSpacing * (kbLayoutRow.layoutOptions.length - 1)) / kbLayoutRow.layoutOptions.length
                                            height: lobby.actionBtnHeight
                                            radius: 6
                                            property bool selected: lobbyKeyboardLayout === modelData.id
                                            color: selected ? "#e8e8e8" : "#f0f0f0"
                                            border.color: selected ? "black" : "#bbb"
                                            border.width: selected ? 2 : 1
                                            Text {
                                                anchors.centerIn: parent
                                                text: modelData.label
                                                font.family: "Noto Sans"
                                                font.pointSize: 11
                                            }
                                            MouseArea {
                                                anchors.fill: parent
                                                onClicked: writerdeck.setKeyboardLayout(modelData.id)
                                            }
                                        }
                                    }
                                }
                                Text {
                                    text: "Bluetooth / phone"
                                    font.pointSize: 14
                                    font.family: "Noto Sans"
                                    color: "black"
                                    width: parent.width
                                }
                                Text {
                                    text: "Pair the keyboard to your phone, then open:\nhttp://" + lobbyIP + ":8000\nTyping is forwarded over Wi-Fi.\n" + (lobbyPIN !== "" ? ("PIN: " + lobbyPIN) : "PIN is not set")
                                    font.pointSize: 11
                                    font.family: "Noto Sans"
                                    color: "#555555"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
                                }
                            }
                        }
                    }
                    // 3 Sync
                    Item {
                        visible: lobbyPage === 3
                        anchors.fill: parent
                        Flickable {
                            anchors.fill: parent
                            contentWidth: width
                            contentHeight: syncCol.height
                            clip: true
                            Column {
                                id: syncCol
                                width: parent.width
                                spacing: lobby.contentSpacing
                                Text {
                                    text: "GitHub sync"
                                    font.pointSize: 14
                                    font.family: "Noto Sans"
                                    color: "black"
                                    width: parent.width
                                }
                                Text {
                                    visible: !(lobbySyncOn && lobbySyncRepo !== "" && !lobbySyncReady)
                                        && !(lobbySyncOn && lobbySyncRepo !== "" && lobbySyncReady && lobbySyncError !== "")
                                    text: lobbySyncOn && lobbySyncRepo !== ""
                                        ? (lobbyLastSync !== ""
                                            ? "Last sync was " + lobbyLastSync + ".\nNotes sync to github.com/" + lobbySyncRepo
                                            : "Notes sync to github.com/" + lobbySyncRepo)
                                        : ("Sync not configured.\nSet up in phone Sync setup:\nhttp://" + lobbyIP + ":8000")
                                    font.pointSize: 11
                                    font.family: "Noto Sans"
                                    color: lobbySyncOn && lobbySyncRepo !== "" ? "#1b5e20" : "#555555"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
                                }
                                Rectangle {
                                    visible: lobbySyncOn && lobbySyncRepo !== "" && lobbySyncReady && lobbySyncError !== ""
                                    width: parent.width
                                    height: syncErrCol.height + 20
                                    color: "white"
                                    border.color: "black"
                                    border.width: 2
                                    radius: 4
                                    Column {
                                        id: syncErrCol
                                        anchors.centerIn: parent
                                        width: parent.width - 20
                                        spacing: 10
                                        Text {
                                            text: !lobbyWifi ? "SYNC OFFLINE" : "SYNC FAILED"
                                            font.pointSize: 16
                                            font.bold: true
                                            font.family: "Noto Sans"
                                            color: "black"
                                            width: parent.width
                                        }
                                        Text {
                                            text: lobbySyncError
                                            font.pointSize: 13
                                            font.family: "Noto Sans"
                                            color: "black"
                                            width: parent.width
                                            wrapMode: Text.WordWrap
                                            lineHeight: 1.25
                                        }
                                    }
                                }
                                Rectangle {
                                    visible: lobbySyncOn && lobbySyncRepo !== "" && !lobbySyncReady
                                    width: parent.width
                                    height: tokenWarnCol.height + 20
                                    color: "white"
                                    border.color: "black"
                                    border.width: 2
                                    radius: 4
                                    Column {
                                        id: tokenWarnCol
                                        anchors.centerIn: parent
                                        width: parent.width - 20
                                        spacing: 10
                                        Text {
                                            text: "TOKEN NEEDED"
                                            font.pointSize: 16
                                            font.bold: true
                                            font.family: "Noto Sans"
                                            color: "black"
                                            width: parent.width
                                        }
                                        Text {
                                            text: "GitHub token is not on the tablet.\nOpen phone Sync setup and tap Save:\nhttp://" + lobbyIP + ":8000\nRepo: github.com/" + lobbySyncRepo
                                            font.pointSize: 13
                                            font.family: "Noto Sans"
                                            color: "black"
                                            width: parent.width
                                            wrapMode: Text.WordWrap
                                            lineHeight: 1.25
                                        }
                                    }
                                }
                                Rectangle {
                                    visible: lobbySyncOn && lobbySyncRepo !== ""
                                    width: parent.width
                                    height: lobby.actionBtnHeight
                                    radius: 6
                                    color: (lobbySyncReady && !lobbySyncing) ? "#f0f0f0" : "white"
                                    border.color: lobbySyncReady ? "#bbb" : "black"
                                    border.width: lobbySyncReady ? 1 : 2
                                    Text {
                                        anchors.centerIn: parent
                                        text: !lobbySyncReady ? "Token needed — phone Sync setup"
                                            : (lobbySyncing ? "Syncing…" : "Sync now")
                                        font.family: "Noto Sans"
                                        font.pointSize: !lobbySyncReady ? 14 : 12
                                        font.bold: !lobbySyncReady
                                        color: (lobbySyncReady && !lobbySyncing) ? "black" : (lobbySyncReady ? "#888888" : "black")
                                    }
                                    MouseArea {
                                        anchors.fill: parent
                                        enabled: lobbySyncReady && !lobbySyncing
                                        onClicked: writerdeck.syncNow()
                                    }
                                }
                                Text {
                                    text: "Sync also runs automatically on save, Home, and every few minutes."
                                    font.pointSize: 10
                                    font.family: "Noto Sans"
                                    color: "#888888"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
                                }
                            }
                        }
                    }
                    // 4 Settings
                    Item {
                        visible: lobbyPage === 4
                        anchors.fill: parent
                        Flickable {
                            anchors.fill: parent
                            contentWidth: width
                            contentHeight: setCol.height
                            clip: true
                            Column {
                                id: setCol
                                width: parent.width
                                spacing: lobby.contentSpacing
                                Text {
                                    text: lobbySettingsMode === "confirm-exit"
                                          ? "Stop Writerdeck? Enter=yes  Esc=no"
                                          : "Settings"
                                    font.pointSize: 14
                                    font.family: "Noto Sans"
                                    color: "black"
                                    width: parent.width
                                }

                                Column {
                                    width: parent.width
                                    spacing: lobby.contentSpacing
                                    visible: lobbySettingsMode === ""

                                    Text {
                                        text: "\nReading font"
                                        font.pointSize: 12
                                        font.family: "Noto Sans"
                                        color: "black"
                                        width: parent.width
                                    }
                                    Grid {
                                        id: fontGrid
                                        width: parent.width
                                        columns: 2
                                        rowSpacing: lobby.tabSpacing
                                        columnSpacing: lobby.tabSpacing
                                        property var fontOptions: [
                                            { id: "Inter", label: "Inter" },
                                            { id: "Literata", label: "Literata" },
                                            { id: "EB Garamond", label: "EB Garamond" },
                                            { id: "DejaVu Sans", label: "DejaVu Sans" }
                                        ]
                                        Repeater {
                                            model: fontGrid.fontOptions
                                            delegate: Rectangle {
                                                width: (fontGrid.width - lobby.tabSpacing) / 2
                                                height: lobby.actionBtnHeight
                                                radius: 6
                                                property bool selected: readFont === modelData.id
                                                color: selected ? "#e8e8e8" : "#f0f0f0"
                                                border.color: selected ? "black" : "#bbb"
                                                border.width: selected ? 2 : 1
                                                Text {
                                                    anchors.centerIn: parent
                                                    text: modelData.label
                                                    font.family: "Noto Sans"
                                                    font.pointSize: 11
                                                }
                                                MouseArea {
                                                    anchors.fill: parent
                                                    onClicked: writerdeck.setReadFont(modelData.id)
                                                }
                                            }
                                        }
                                    }

                                    Text {
                                        text: "\nPrivate notes"
                                        font.pointSize: 12
                                        font.family: "Noto Sans"
                                        color: "black"
                                        width: parent.width
                                    }
                                    Text {
                                        text: lobbyEncryptionEnabled
                                            ? "On — encrypted notes require PIN to open, read, or edit"
                                            : "Off — optional encryption with a separate 6-digit PIN. Recovery via GitHub secret/pin when sync is on."
                                        font.pointSize: 10
                                        font.family: "Noto Sans"
                                        color: "#666666"
                                        width: parent.width
                                        wrapMode: Text.WordWrap
                                    }
                                    Row {
                                        width: parent.width
                                        spacing: lobby.tabSpacing
                                        visible: !lobbyEncryptionEnabled
                                        Rectangle {
                                            width: (parent.width - lobby.tabSpacing) / 2
                                            height: lobby.actionBtnHeight
                                            radius: 6
                                            color: "#f0f0f0"
                                            border.color: "#bbb"
                                            border.width: 1
                                            Text {
                                                anchors.centerIn: parent
                                                text: "Enable"
                                                font.family: "Noto Sans"
                                                font.pointSize: 11
                                            }
                                            MouseArea {
                                                anchors.fill: parent
                                                onClicked: root.vaultBeginSetup()
                                            }
                                        }
                                    }
                                    Row {
                                        width: parent.width
                                        spacing: lobby.tabSpacing
                                        visible: lobbyEncryptionEnabled
                                        Rectangle {
                                            width: parent.width
                                            height: lobby.actionBtnHeight
                                            radius: 6
                                            color: "#f0f0f0"
                                            border.color: "#bbb"
                                            border.width: 1
                                            Text {
                                                anchors.centerIn: parent
                                                text: "Change PIN"
                                                font.family: "Noto Sans"
                                                font.pointSize: 11
                                            }
                                            MouseArea {
                                                anchors.fill: parent
                                                onClicked: root.vaultBeginChangePIN()
                                            }
                                        }
                                    }

                                    Text {
                                        text: "\nPIN for phone pairing"
                                        font.pointSize: 12
                                        font.family: "Noto Sans"
                                        color: "black"
                                        width: parent.width
                                    }
                                    Text {
                                        text: "Adding a PIN ensures that only intended devices can access your notes"
                                        font.pointSize: 10
                                        font.family: "Noto Sans"
                                        color: "#666666"
                                        width: parent.width
                                        wrapMode: Text.WordWrap
                                    }
                                    Column {
                                        width: parent.width
                                        spacing: lobby.tabSpacing
                                        Repeater {
                                            model: [
                                                { id: "6", label: "6 digits" },
                                                { id: "4", label: "4 digits" },
                                                { id: "none", label: "No PIN", warn: "Anyone on Wi-Fi can read and edit notes" }
                                            ]
                                            delegate: Rectangle {
                                                width: parent.width
                                                height: modelData.warn ? lobby.actionBtnHeight + 28 : lobby.actionBtnHeight
                                                radius: 6
                                                property bool selected: lobbyPinDigits === modelData.id
                                                color: selected ? "#e8e8e8" : "#f0f0f0"
                                                border.color: selected ? "black" : "#bbb"
                                                border.width: selected ? 2 : 1
                                                Column {
                                                    anchors.centerIn: parent
                                                    width: parent.width - 16
                                                    spacing: 2
                                                    Text {
                                                        anchors.horizontalCenter: parent.horizontalCenter
                                                        text: modelData.label
                                                        font.family: "Noto Sans"
                                                        font.pointSize: 11
                                                        color: "black"
                                                    }
                                                    Text {
                                                        visible: !!modelData.warn
                                                        anchors.horizontalCenter: parent.horizontalCenter
                                                        text: modelData.warn || ""
                                                        font.family: "Noto Sans"
                                                        font.pointSize: 9
                                                        color: "#666666"
                                                        horizontalAlignment: Text.AlignHCenter
                                                        wrapMode: Text.WordWrap
                                                        width: parent.width
                                                    }
                                                }
                                                MouseArea {
                                                    anchors.fill: parent
                                                    onClicked: writerdeck.setPinDigits(modelData.id)
                                                }
                                            }
                                        }
                                    }

                                    Text {
                                        text: "\nDisplay rotation"
                                        font.pointSize: 12
                                        font.family: "Noto Sans"
                                        color: "black"
                                        width: parent.width
                                    }
                                    Text {
                                        text: root.rotation + " degrees. Ctrl-R or Ctrl+arrows to rotate."
                                        font.pointSize: 10
                                        font.family: "Noto Sans"
                                        color: "#666666"
                                        width: parent.width
                                        wrapMode: Text.WordWrap
                                    }
                                    Rectangle {
                                        width: parent.width
                                        height: lobby.actionBtnHeight
                                        radius: 6
                                        color: "#f0f0f0"
                                        border.color: "#bbb"
                                        border.width: 1
                                        Text {
                                            anchors.centerIn: parent
                                            text: "Rotate 90"
                                            font.family: "Noto Sans"
                                            font.pointSize: 12
                                        }
                                        MouseArea {
                                            anchors.fill: parent
                                            onClicked: root.rotateScreen()
                                        }
                                    }

                                    Text {
                                        text: "\nService"
                                        font.pointSize: 12
                                        font.family: "Noto Sans"
                                        color: "black"
                                        width: parent.width
                                    }
                                    Text {
                                        text: "Stop Writerdeck and return the tablet to the stock reMarkable UI. Reconnect later via SSH or reboot."
                                        font.pointSize: 10
                                        font.family: "Noto Sans"
                                        color: "#666666"
                                        width: parent.width
                                        wrapMode: Text.WordWrap
                                    }
                                    Rectangle {
                                        width: parent.width
                                        height: lobby.actionBtnHeight
                                        radius: 6
                                        color: "#f0f0f0"
                                        border.color: "#bbb"
                                        border.width: 1
                                        Text {
                                            anchors.centerIn: parent
                                            text: "Exit Writerdeck"
                                            font.family: "Noto Sans"
                                            font.pointSize: 12
                                        }
                                        MouseArea {
                                            anchors.fill: parent
                                            onClicked: root.lobbySettingsBeginExit()
                                        }
                                    }
                                }
                            }
                        }
                    }
                    // 5 Shortcuts
                    Item {
                        visible: lobbyPage === 5
                        anchors.fill: parent
                        Flickable {
                            anchors.fill: parent
                            contentWidth: width
                            contentHeight: scCol.height
                            clip: true
                            Column {
                                id: scCol
                                width: parent.width
                                spacing: lobby.contentSpacing
                                Text {
                                    text: "Shortcuts"
                                    font.pointSize: 14
                                    font.family: "Noto Sans"
                                    color: "black"
                                    width: parent.width
                                }
                                Text {
                                    text: "Lobby: Tab / arrows / 1-6 switch pages\nFiles: Up/Down select  Enter edit  v read  n d r\nStock UI: Esc (USB) or L+R page buttons → Lobby\nCtrl-K: quick file picker\nCtrl-R: rotate  Ctrl-Q: quit\nHome: exit to reMarkable UI"
                                    font.pointSize: 10
                                    font.family: "Noto Mono"
                                    color: "#555555"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
                                }
                            }
                        }
                    }
                }
                Text {
                    id: lobbyHint
                    anchors.bottom: parent.bottom
                    anchors.left: parent.left
                    anchors.right: parent.right
                    anchors.leftMargin: lobby.pageMargin
                    anchors.rightMargin: lobby.pageMargin
                    anchors.bottomMargin: lobby.pageMargin
                    height: 48
                    verticalAlignment: Text.AlignVCenter
                    font.family: "Noto Mono"
                    font.pointSize: 9
                    color: "#888888"
                    text: lobbyPage === 1 && lobbyFilesMode === ""
                        ? "n new  Enter edit  v read  r rename  d delete"
                        : "Tab next page  1-6 jump  Ctrl-K files"
                }
            }
        }
