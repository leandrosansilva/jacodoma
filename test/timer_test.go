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

func (logic *FakeTurnLogic) OnTimeGetsCritical(t time.Time) {
	fmt.Printf("Changed to critical on %s\n", t)
	logic.Actions = append(logic.Actions, FakeTimerAction{t, "time_critical", Participant{}})
}

func (logic *FakeTurnLogic) OnNextParticipantStarts(t time.Time, p Participant) {
	fmt.Printf("next starts on %s\n", t)
	logic.Actions = append(logic.Actions, FakeTimerAction{t, "next_participant", p})
	logic.CurrentParticipantIndex = (logic.CurrentParticipantIndex + 1) % logic.Participants.Length()
}

func (logic *FakeTurnLogic) OnTimeIsOver(t time.Time) {
	fmt.Printf("time is over on %s\n", t)
	logic.Actions = append(logic.Actions, FakeTimerAction{t, "time_over", Participant{}})
}

func (logic *FakeTurnLogic) OnStartsWaitingNextParticipant(t time.Time) {
	fmt.Printf("waits next on %s\n", t)
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

		begin := time.Unix(1000, 0)
		turnInfo := TurnTimeInfo{270 * time.Second, 30 * time.Second}

		logic := NewFakeTurnLogic(turnInfo, participants)
		timer := NewTimer(logic)
		So(timer, should.NotEqual, nil)

		// runs for 6min
		ExecuteTimer(
			timer, begin, begin.Add(6*time.Minute),
			100*time.Millisecond)

		// TODO: test the time when each event happened
		//So(len(logic.Actions), should.Equal, 3)

		// first user starts after 5sec
		So(logic.Actions[0].Time, should.Equal, begin.Add(5*time.Second))
		So(logic.Actions[0].Action, should.Equal, "next_participant")
		So(logic.Actions[0].Participant.Email, should.Equal, "coding@do.jo")

		// on 4:30 time gets critical
		So(logic.Actions[1].Time, should.Equal, begin.Add(270*time.Second))
		So(logic.Actions[1].Action, should.Equal, "time_critical")

		// on 5:00 time is over
		So(logic.Actions[2].Time, should.Equal, begin.Add(300*time.Second))
		So(logic.Actions[2].Action, should.Equal, "time_over")

		// on 5:20 time is over
		So(logic.Actions[2].Time, should.Equal, begin.Add(320*time.Second))
		So(logic.Actions[2].Action, should.Equal, "waiting_next_participant")
	})
}
