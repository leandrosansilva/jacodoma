import QtQuick 2.0

Rectangle {
  id: root
  width: 320; height: 480
  color: "lightgray"

  Text {
    function padTime(time) {
      return time.length == 1 ? "0" + time : time
    }

    function formatTimer(duration) {
      // duration is in nanosecond
      var secTotal = duration / 1000000000
      var sec = (secTotal % 60).toFixed()
      var min = Math.floor(secTotal / 60)

      return padTime(min.toString()) + ":" + padTime(sec.toString())
    }
    id: "timer"
    itext: formatTimer(ctrl.duration)
    y: 30
    font.pointSize: 24; font.bold: true
  }
}
