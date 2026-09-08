import QtQuick
import QtQuick.Window

Window {
    width: Screen.width
    height: Screen.height
    visible: true
    color: "white"

    Text {
        id: title
        x: 40
        y: 40
        width: parent.width - 80
        font.pixelSize: 28
        text: "Writerdeck socket spike"
    }

    Text {
        id: body
        x: 40
        y: 120
        width: parent.width - 80
        font.pixelSize: 24
        wrapMode: Text.Wrap
        text: socketFeed.displayText
    }
}
