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

                Connections {
                    target: root
                    onIsLobbyChanged: if (isLobby) Qt.callLater(function() { lobbyFocus.forceActiveFocus() })
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
                                height: parent.height - lobbyFilesBar.height - lobby.contentSpacing * 2
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
                                        text: model.name
                                        font.family: "Noto Mono"
                                        font.pointSize: 11
                                        color: "black"
                                        elide: Text.ElideRight
                                        width: parent.width - 32
                                    }
                                    MouseArea {
                                        anchors.fill: parent
                                        onClicked: {
                                            lobbyFilesIndex = index
                                            lobbyFilesList.currentIndex = index
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
                                    model: ["New", "Open", "Rename", "Delete"]
                                    delegate: Rectangle {
                                        width: (lobbyFilesBar.width - lobby.tabSpacing * 3) / 4
                                        height: lobby.actionBtnHeight
                                        radius: 6
                                        color: "#f0f0f0"
                                        border.color: "#bbb"
                                        border.width: 1
                                        Text {
                                            anchors.centerIn: parent
                                            text: modelData
                                            font.family: "Noto Sans"
                                            font.pointSize: 12
                                        }
                                        MouseArea {
                                            anchors.fill: parent
                                            onClicked: {
                                                if (modelData === "New") root.lobbyFilesBeginNew()
                                                else if (modelData === "Open") root.lobbyOpenSelected()
                                                else if (modelData === "Rename") root.lobbyFilesBeginRename()
                                                else if (modelData === "Delete") root.lobbyFilesBeginDelete()
                                            }
                                        }
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
                                    text: "Connect with a USB OTG cable.\nLayout is US QWERTY until locale qmaps ship."
                                    font.pointSize: 11
                                    font.family: "Noto Sans"
                                    color: "#555555"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
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
                                    text: lobbySyncOn && lobbySyncRepo !== ""
                                        ? (lobbyLastSync !== ""
                                            ? "Last sync was " + lobbyLastSync + ".\nNotes sync to github.com/" + lobbySyncRepo
                                            : "Notes sync to github.com/" + lobbySyncRepo)
                                        : ("Not configured. Set up in your phone's browser at http://" + lobbyIP + ":8000")
                                    font.pointSize: 11
                                    font.family: "Noto Sans"
                                    color: lobbySyncOn && lobbySyncRepo !== "" ? "#1b5e20" : "#555555"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
                                }
                                Text {
                                    text: "Sync runs from the phone browser when connected."
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
                                    text: "Settings"
                                    font.pointSize: 14
                                    font.family: "Noto Sans"
                                    color: "black"
                                    width: parent.width
                                }
                                Text {
                                    text: "Reading font: " + readFont + "\n(change on phone Preferences)"
                                    font.pointSize: 11
                                    font.family: "Noto Sans"
                                    color: "#555555"
                                    width: parent.width
                                    wrapMode: Text.WordWrap
                                }
                                Text {
                                    text: "Display rotation: " + root.rotation + " degrees\nCtrl-R or Ctrl+arrows to rotate."
                                    font.pointSize: 11
                                    font.family: "Noto Sans"
                                    color: "#555555"
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
                                    text: "Lobby: Tab / arrows / 1-6 switch pages\nFiles: Up/Down select  Enter open  n d r\nStock UI: Esc (USB) or L+R page buttons → Lobby\nCtrl-K: quick file picker\nCtrl-R: rotate  Ctrl-Q: quit\nHome: exit to reMarkable UI"
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
                        ? "n new  Enter open  r rename  d delete"
                        : "Tab next page  1-6 jump  Ctrl-K files"
                }
            }
        }
