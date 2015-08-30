import QtQuick 2.0
import QtQuick.Controls 1.1

Rectangle {
  id: root
  width: 320; height: 480
  color: "lightgray"

  function padTime(time) {
    return time.length == 1 ? "0" + time : time
  }

  function formatTime(duration) {
    // duration is in nanosecond
    var secTotal = duration / 1000000000
    var sec = (secTotal % 60).toFixed()
    var min = Math.floor(secTotal / 60)
    var hour = Math.floor(secTotal / 3600)

    return (hour == 0 
              ? "" 
              : padTime(hour) + ":") 
           + padTime(min.toString()) + ":" + padTime(sec.toString())
  }

  Text {
    id: "timer"
    text: formatTime(ctrl.turnDuration)
    y: 30
    font.pointSize: 24; font.bold: true
  }

  Text {
    text: formatTime(ctrl.sessionDuration)
    y: 150
  }

  Text {
    text: ctrl.state
    y: 200
  }

  Text {
    text: ctrl.participants.get(ctrl.currentParticipantIndex).email
    y: 250
  }

  Button {
    function buttonLabel() {
      return "Start: " + ctrl.participants.get(ctrl.currentParticipantIndex).name
    }

    y: 300
    x: 100
    text: buttonLabel()
    onClicked: ctrl.setParticipantReady()    
  }
}
