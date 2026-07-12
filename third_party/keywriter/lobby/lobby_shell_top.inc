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

