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
                    text: lobbyPage === 0 && lobbyFilesMode === ""
                        ? "n new  Enter edit  v read  r rename  d delete"
                        : "Tab next page  1-6 jump  Ctrl-K files"
                }
            }
        }
