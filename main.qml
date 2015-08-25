// See http://qt-project.org/doc/qt-5.1/qtquick/qml-tutorial3.html

import QtQuick 2.0

Rectangle {
  id: page
  width: 320; height: 480
  color: "lightgray"

  Text {
    text: "00:00"
    model: timer_time.duration
    y: 30
    anchors.horizontalCenter: page.horizontalCenter
    font.pointSize: 24; font.bold: true
  }
}
