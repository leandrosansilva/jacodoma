package jacodoma

import (
	. "../src/"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type FakeTimerActions struct {
	Time   time.Time
	Action string
}

type FakeTimerLogic struct {
	TurnTimeInfo TurnTimeInfo
	Participants Participants
	timeChannel  chan time.Time
	Actions      []FakeTimerActions
}

func (logic *FakeTimerLogic) OnStarted(t time.Time) {
}

func (logic *FakeTimerLogic) OnTimeGetsCritical(t time.Time) {
}

func (logic *FakeTimerLogic) OnNextParticipantStarts(t time.Time, p Participant) {
}

func (logic *FakeTimerLogic) OnCriticalTimeHasReached(t time.Time) {
}

func (logic *FakeTimerLogic) OnTimeIsOver(t time.Time) {
}

func (logic *FakeTimerLogic) OnStartsWaitingNextParticipant(t time.Time) {
}

func (logic *FakeTimerLogic) TimeChannel() chan time.Time {
	return logic.timeChannel
}

func (logic *FakeTimerLogic) HasFinished() bool {
	return true
}

func (logic *FakeTimerLogic) NextParticipant() Participant {
	return Participant{}
}

func (logic *FakeTimerLogic) CurrentParticipant() Participant {
	return Participant{}
}

func NewFakeTimerLogic(info TurnTimeInfo, participants Participants) *FakeTimerLogic {
	return &FakeTimerLogic{
		info, participants,
		make(chan time.Time),
		make([]FakeTimerActions, 0)}
}

func runFakeTurns(logic *FakeTimerLogic, numberOfTurns int) {

}

func TestTimer(t *testing.T) {
	Convey("One Turn Timer", t, func() {
		participants := BuildParticipantsFromArray([]Participant{
			{"Coding Dojo", "coding@do.jo"},
			{"Manoel Ribas", "manoel@ribas.go"},
			{"Juka Juke", "juka@ju.ke"},
			{"Jon Doe", "joe@doe.com"},
		})

		turnInfo := TurnTimeInfo{5 * 60, 4.5 * 60}

		logic := NewFakeTimerLogic(turnInfo, participants)
		timer := NewTimer(logic)
		So(timer, should.NotEqual, nil)

		go runFakeTurns(logic, 1)

		timer.Run()

		// TODO: test the time when each event happened
		//So(logic.)
	})
}
