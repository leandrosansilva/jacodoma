package jacodoma

import (
	. "../src/"
	"fmt"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type FakeTimerAction struct {
	Time        time.Time
	Action      string
	Participant Participant
}

type FakeTurnLogic struct {
	info                    TurnTimeInfo
	Participants            Participants
	Actions                 []FakeTimerAction
	CurrentParticipantIndex int
}

func printTime(t time.Time, ev string) {
	//fmt.Printf("EV: %s TIME: %s SEC: %d and NS: %d\n", ev, t, t.Second(), t.Nanosecond())
}

func pt(t time.Time) string {
	return fmt.Sprintf("%d:%d.%d", t.Minute(), t.Second(), t.Nanosecond())
}

func (logic *FakeTurnLogic) OnTimeGetsCritical(t time.Time) {

	logic.Actions = append(logic.Actions, FakeTimerAction{t, "time_critical", Participant{}})
}

func (logic *FakeTurnLogic) OnNextParticipantStarts(t time.Time, p Participant) {
	printTime(t, "next starts")
	logic.Actions = append(logic.Actions, FakeTimerAction{t, "next_participant", p})
}

func (logic *FakeTurnLogic) OnTimeIsOver(t time.Time) {
	printTime(t, "time over")
	logic.Actions = append(logic.Actions, FakeTimerAction{t, "time_over", Participant{}})
}

func (logic *FakeTurnLogic) OnStartsWaitingNextParticipant(t time.Time) {
	printTime(t, "starts waiting")
	logic.Actions = append(logic.Actions, FakeTimerAction{t, "waiting_next_participant", Participant{}})
}

func (logic *FakeTurnLogic) NextParticipant() Participant {
	// FIXME: code repetition
	return logic.Participants.Get((logic.CurrentParticipantIndex + 1) % logic.Participants.Length())
}

func (logic *FakeTurnLogic) CurrentParticipant() Participant {
	return logic.Participants.Get(logic.CurrentParticipantIndex)
}

func (logic *FakeTurnLogic) TurnTimeInfo() *TurnTimeInfo {
	return &logic.info
}

func (logic *FakeTurnLogic) NextParticipantIsReady() bool {
	return true
}

func NewFakeTurnLogic(info TurnTimeInfo, participants Participants) *FakeTurnLogic {
	return &FakeTurnLogic{
		info, participants,
		make([]FakeTimerAction, 0), 0}
}

func ExecuteTimer(timer *Timer, begin, end time.Time, duration time.Duration) {
	for t := begin; t.Unix() < end.Unix(); t = t.Add(duration) {
		timer.Step(t)
	}
}

func TestTimer(t *testing.T) {
	Convey("One Turn Timer", t, func() {
		participants := BuildParticipantsFromArray([]Participant{
			{"Coding Dojo", "coding@do.jo"},
			{"Manoel Ribas", "manoel@ribas.go"},
			{"Juka Juke", "juka@ju.ke"},
			{"Jon Doe", "joe@doe.com"},
		})

		begin := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

		// Turn lasts 5min and the last 30sec are "critical"
		turnInfo := TurnTimeInfo{270 * time.Second, 30 * time.Second}

		logic := NewFakeTurnLogic(turnInfo, participants)
		timer := NewTimer(logic)
		So(timer, should.NotEqual, nil)

		// runs for 5:01 min
		ExecuteTimer(
			timer, begin, begin.Add(5*time.Minute+1*time.Second),
			100*time.Millisecond)

		// TODO: test the time when each event happened
		//So(len(logic.Actions), should.Equal, 3)

		Convey("User starts on 0sec", func() {
			So(logic.Actions[0].Action, should.Equal, "next_participant")
			So(logic.Actions[0].Participant.Email, should.Equal, "coding@do.jo")
			So(pt(logic.Actions[0].Time), should.Equal, pt(begin))
		})

		Convey("Time gets critical on 4:30", func() {
			So(logic.Actions[1].Action, should.Equal, "time_critical")
			So(pt(logic.Actions[1].Time), should.Equal, pt(begin.Add(270*time.Second)))
		})

		Convey("Time is over on 5:00", func() {
			So(logic.Actions[2].Action, should.Equal, "time_over")
			So(pt(logic.Actions[2].Time), should.Equal, pt(begin.Add(300*time.Second)))
		})
	})
}
