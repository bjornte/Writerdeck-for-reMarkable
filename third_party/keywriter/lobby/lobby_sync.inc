                    // 3 Sync
                    Item {
                        visible: lobbyPage === 2
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
