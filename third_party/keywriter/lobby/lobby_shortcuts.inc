                    // 5 Shortcuts
                    Item {
                        visible: lobbyPage === 4
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
