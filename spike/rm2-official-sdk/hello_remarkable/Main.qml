import QtQuick
import QtQuick.Window

Window {
    width: Screen.width
    height: Screen.height
    visible: true

    Text {
        id: hello_text
        x: 50
        y: 50
        font.pixelSize: 32
        text: "Hello reMarkable!"
    }

    MouseArea {
        anchors.fill: parent
        onPressed: {
            hello_text.visible = !hello_text.visible
        }
    }

    // Spike: write PNG for autonomous verify (scripts/capture-screenshot.sh).
    Timer {
        interval: 2500
        running: true
        repeat: false
        onTriggered: {
            contentItem.grabToImage(function(result) {
                if (!result) {
                    console.log("grabToImage failed")
                    return
                }
                var path = "/home/root/spike-rm2-hello/screen.png"
                if (result.saveToFile(path)) {
                    console.log("saved", path)
                } else {
                    console.log("saveToFile failed", path)
                }
            })
        }
    }
}
