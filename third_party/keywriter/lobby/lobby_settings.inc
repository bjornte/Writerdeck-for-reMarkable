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
