                    // 1 Files
                    Item {
                        visible: lobbyPage === 0
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
                                text: lobbyFilesInputDisplay()
                                font.family: "Noto Mono"
                                font.pointSize: 12
                                color: "#333"
                                width: parent.width
                            }

                            Text {
                                visible: lobbyVaultError !== "" && lobbyFilesMode === ""
                                text: lobbyVaultError
                                font.family: "Noto Sans"
                                font.pointSize: 11
                                color: "#aa0000"
                                width: parent.width
                                wrapMode: Text.WordWrap
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
                                        text: lobbyFilesStripSuffix(model.name) + (model.encrypted ? " [private]" : "")
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
