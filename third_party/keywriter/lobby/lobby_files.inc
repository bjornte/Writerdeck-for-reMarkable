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
                        }
                    }
