                    // 2 Keyboard
                    Item {
                        visible: lobbyPage === 1
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
