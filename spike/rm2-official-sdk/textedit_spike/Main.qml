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
        text: "Writerdeck text spike"
    }

    Text {
        id: body
        x: 40
        y: 120
        width: parent.width - 80
        font.pixelSize: 24
        wrapMode: Text.Wrap
        text: "The quick brown fox jumps over the lazy dog. Wrapped prose on epaper is the next gate before the full editor port."
    }
}
