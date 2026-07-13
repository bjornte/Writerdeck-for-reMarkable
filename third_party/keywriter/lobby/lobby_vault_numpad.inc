        Rectangle {
            id: vaultOverlay
            anchors.fill: parent
            color: "#f8f8f8"
            visible: vaultOverlayMode !== ""
            z: 25

            Column {
                anchors.centerIn: parent
                width: parent.width * 0.7
                spacing: lobby.contentSpacing

                Text {
                    width: parent.width
                    horizontalAlignment: Text.AlignHCenter
                    wrapMode: Text.WordWrap
                    font.family: "Noto Sans"
                    font.pointSize: 14
                    color: "black"
                    text: vaultOverlayMode === "setup" ? "Choose a 6-digit private PIN"
                        : vaultOverlayMode === "confirm" ? "Confirm private PIN"
                        : vaultOverlayMode === "change-old" ? "Enter current private PIN"
                        : vaultOverlayMode === "change-new" ? "Enter new private PIN"
                        : vaultOverlayMode === "change-confirm" ? "Confirm new private PIN"
                        : vaultOverlayReason !== "" ? vaultOverlayReason
                        : "Enter private PIN"
                }

                Text {
                    width: parent.width
                    horizontalAlignment: Text.AlignHCenter
                    font.family: "Noto Mono"
                    font.pointSize: 18
                    color: "#333"
                    text: vaultPinDisplay()
                }

                Grid {
                    id: vaultPad
                    width: parent.width
                    columns: 3
                    rowSpacing: 8
                    columnSpacing: 8
                    Repeater {
                        model: ["1","2","3","4","5","6","7","8","9","Bksp","0","Done"]
                        delegate: Rectangle {
                            width: (vaultPad.width - 16) / 3
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
                                onClicked: vaultNumpadTap(modelData)
                            }
                        }
                    }
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
                        text: "Cancel"
                        font.family: "Noto Sans"
                        font.pointSize: 12
                    }
                    MouseArea {
                        anchors.fill: parent
                        onClicked: vaultNumpadCancel()
                    }
                }
            }
        }
