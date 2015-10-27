package jacodoma

import (
	"fmt"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	. "jacodoma/src"
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
	ParticipantIsReady      bool
}

func pt(t time.Time) string {
	r := t.Round(200 * time.Millisecond)
	return fmt.Sprintf("%02d:%02d", r.Minute(), r.Minute())
}

func (logic *FakeTurnLogic) OnTimeGetsCritical(t time.Time) {
	logic.Actions = append(logic.Actions, FakeTimerAction{t, "time_critical", Participant{}})
}

func (logic *FakeTurnLogic) OnNextParticipantStarts(t time.Time) {
	logic.Actions = append(logic.Actions, FakeTimerAction{t, "next_participant", Participant{}})
}

func (logic *FakeTurnLogic) OnTimeIsOver(t time.Time) {
	logic.Actions = append(logic.Actions, FakeTimerAction{t, "time_over", Participant{}})
}

func (logic *FakeTurnLogic) OnStartsWaitingNextParticipant(t time.Time, index int) {
	logic.Actions = append(logic.Actions, FakeTimerAction{
		t, "waiting_next_participant",
		logic.Participants.Get(logic.CurrentParticipantIndex)})

	logic.CurrentParticipantIndex = index
}

func (logic *FakeTurnLogic) BlockSession(t time.Time) {
	logic.Actions = append(logic.Actions, FakeTimerAction{t, "block_session", Participant{}})
}

func (logic *FakeTurnLogic) NextParticipantIsReady() bool {
	return logic.ParticipantIsReady
}

func (logic *FakeTurnLogic) NextParticipantIndex() int {
	return (logic.CurrentParticipantIndex + 1) % logic.Participants.Length()
}

func (logic *FakeTurnLogic) TurnTimeInfo() *TurnTimeInfo {
	return &logic.info
}

func NewFakeTurnLogic(info TurnTimeInfo, participants Participants) *FakeTurnLogic {
	return &FakeTurnLogic{
		info, participants,
		make([]FakeTimerAction, 0),
		0, false}
}

type TimerExecutor struct {
	Time  time.Time
	Timer *Timer
	Logic *FakeTurnLogic
}

func NewTimerExecutor(begin time.Time, timer *Timer, logic *FakeTurnLogic) *TimerExecutor {
	return &TimerExecutor{begin, timer, logic}
}

func (this *TimerExecutor) Execute(end time.Time, readyTime time.Time) {
	t := this.Time
	readyIsSet := false

	for ; t.Before(end); t = t.Add(100 * time.Millisecond) {
		if readyTime.Second() == t.Second() && !readyIsSet {
			this.Logic.ParticipantIsReady = true
			this.Timer.Step(t)
			t = t.Add(100 * time.Millisecond)
			this.Timer.Step(t)
			this.Logic.ParticipantIsReady = false
			readyIsSet = true
		} else {
			this.Timer.Step(t)
		}
	}

	this.Time = t
}

type FakeTimerUi struct {
	turnTimes    []time.Duration
	sessionTimes []time.Duration
}

func NewFakeTimerUi() *FakeTimerUi {
	return &FakeTimerUi{}
}

func (ui *FakeTimerUi) UpdateTimer(d time.Duration) {
	ui.turnTimes = append(ui.turnTimes, d)
}

func (ui *FakeTimerUi) UpdateSessionTime(d time.Duration) {
	ui.sessionTimes = append(ui.sessionTimes, d)
}

func TestCodingDojoWithFourParticipants(t *testing.T) {
	// Turn lasts 5min and the last 30secs are "critical"
	turnInfo := TurnTimeInfo{270 * time.Second, 30 * time.Second}

	logic := NewFakeTurnLogic(turnInfo, BuildParticipantsFromArray([]Participant{
		{"Coding Dojo", "coding@do.jo"},
		{"Manoel Ribas", "manoel@ribas.go"},
		{"Juka Juke", "juka@ju.ke"},
		{"Jon Doe", "joe@doe.com"},
	}))

	turnTimerChannel := make(DurationChannel, 0)
	sessionTimerChannel := make(DurationChannel, 0)

	ui := NewFakeTimerUi()

	// user interface loop
	go func() {
		for {
			select {
			case s := <-turnTimerChannel:
				if s == -1 {
					return
				}
				ui.UpdateTimer(s)
			case s := <-sessionTimerChannel:
				ui.UpdateSessionTime(s)
			}
		}
	}()

	timer := NewTimer(logic, turnTimerChannel, sessionTimerChannel)

	genesis := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	ex := NewTimerExecutor(genesis, timer, logic)

	Convey("Test Timer with 4 users", t, func() {
		Convey("First Participant Turn", func() {
			begin := ex.Time

			// runs from 00:00 to 05:20
			ex.Execute(begin.Add(5*time.Minute+20*time.Second), begin.Add(10*time.Second))

			So(len(logic.Actions), should.Equal, 6)
			So(logic.ParticipantIsReady, should.Equal, false)

			So(logic.Actions[0].Action, should.Equal, "waiting_next_participant")
			So(logic.Actions[0].Participant.Email, should.Equal, "coding@do.jo")
			So(pt(logic.Actions[0].Time), should.Equal, pt(begin.Add(10*time.Second)))

			So(logic.Actions[1].Action, should.Equal, "next_participant")
			So(pt(logic.Actions[1].Time), should.Equal, pt(begin.Add(10*time.Second)))

			So(logic.Actions[2].Action, should.Equal, "time_critical")
			So(pt(logic.Actions[2].Time), should.Equal, pt(begin.Add(271*time.Second)))

			So(logic.Actions[3].Action, should.Equal, "time_over")
			So(pt(logic.Actions[3].Time), should.Equal, pt(begin.Add(301*time.Second)))

			So(logic.Actions[4].Action, should.Equal, "waiting_next_participant")
			So(logic.Actions[4].Participant.Email, should.Equal, "manoel@ribas.go")
			So(pt(logic.Actions[4].Time), should.Equal, pt(begin.Add(301*time.Second+100*time.Millisecond)))

			So(logic.Actions[5].Action, should.Equal, "block_session")
			So(pt(logic.Actions[5].Time), should.Equal, pt(begin.Add(301*time.Second+100*time.Millisecond)))
		})

		Convey("Second Participant Turn", func() {
			begin := ex.Time

			// the second one takes 40s to get ready
			ex.Execute(begin.Add(6*time.Minute), begin.Add(40*time.Second))

			So(len(logic.Actions), should.Equal, 11)
			So(logic.ParticipantIsReady, should.Equal, false)

			So(logic.Actions[6].Action, should.Equal, "next_participant")
			So(pt(logic.Actions[6].Time), should.Equal, pt(begin.Add(40*time.Second)))

			So(logic.Actions[7].Action, should.Equal, "time_critical")
			So(pt(logic.Actions[7].Time), should.Equal, pt(begin.Add(310*time.Second)))

			So(logic.Actions[8].Action, should.Equal, "time_over")
			So(pt(logic.Actions[8].Time), should.Equal, pt(begin.Add(340*time.Second)))

			So(logic.Actions[9].Action, should.Equal, "waiting_next_participant")
			So(pt(logic.Actions[9].Time), should.Equal, pt(begin.Add(340*time.Second+100*time.Millisecond)))

			So(logic.Actions[10].Action, should.Equal, "block_session")
			So(pt(logic.Actions[10].Time), should.Equal, pt(begin.Add(340*time.Second+100*time.Millisecond)))
		})

		Convey("Times received by the timer", func() {
			d := ex.Time.Sub(genesis)
			So(d, should.Equal, time.Second*680)
			So(len(ui.turnTimes), should.Equal, 604)
		})

		Convey("Session Duration", func() {
			So(len(ui.sessionTimes), should.Equal, 681)
			So(ui.sessionTimes[0], should.Equal, 0)
		})
	})
}
