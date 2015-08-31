import QtQuick 2.0
import QtQuick.Controls 1.1

Rectangle {
  id: root
  width: 1000; 
  height: 700
  color: "lightgray"

  function padTime(time) {
    return time.length == 1 ? "0" + time : time
  }

  function formatTime(duration) {
    // duration is in nanosecond
    var secTotal = duration / 1000000000
    var sec = (secTotal % 60).toFixed()
    var min = Math.floor((secTotal / 60) % 60)
    var hour = Math.floor(secTotal / 3600)

    return (hour == 0 
              ? "" 
              : padTime(hour) + ":") 
           + padTime(min.toString()) + ":" + padTime(sec.toString())
  }

  function formatTimer(duration) {
    var remaining = ctrl.info.totalTurnTime() - duration
    // workaround for small negative values :-)
    return formatTime(remaining <= 0 ? 0 : remaining)
  }

  function buildParticipantAvatarSourceUrl(email) {
    return "image://gravatar/" + email
  }

  function colorFromState(state) {
    var map = {
      "start" : "lightgreen",
      "hurry_up" : "yellow",
      "waiting_participant" : "lightblue",
      "time_over" : "red"
    }

    return map[state]
  }

  Rectangle {
    y: 30
    width: childrenRect.width
    height: childrenRect.height
    color: colorFromState(ctrl.state) 
    Text {
      id: "timer"
      text: formatTimer(ctrl.turnDuration)
      font.pointSize: 24; 
      font.bold: true
    }
  }

  Text {
    text: formatTime(ctrl.sessionDuration)
    y: 150
  }

  Button {
    function buttonLabel() {
      return "Start"
    }

    y: 300
    x: 100
    text: buttonLabel()
    onClicked: ctrl.setParticipantReady()    
  }

  Image {
    source: buildParticipantAvatarSourceUrl(ctrl.participants.get(ctrl.currentParticipantIndex).email)
    width: 128
    height: 128
    y: 350
    x: 10
  }

  ListView {
    width: 200
    height: 500

    y: 0
    x: 400

    model: ctrl.participantsLen

    delegate: Rectangle {
      width: 100
      height: 50

      Image {
        width: parent.height
        height: parent.height
        y: 0
        x: 0
        source: buildParticipantAvatarSourceUrl(ctrl.participants.get(index).email)
      }

      Text {
        height: parent.height
        x: 60
        text: ctrl.participants.get(index).name
      }
    }
  }

}
