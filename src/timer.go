package jacodoma

import (
	"fmt"
	"time"
)

type TurnTimeInfo struct {
	// Duration of the "Non Critical" time,
	// when the participant can relax and code
	RelaxAndCodeDuration time.Duration

	// Duration of Time the participant has to hurry up
	// because the time is ending
	HurryUpDuration time.Duration
}

type ITurnLogic interface {
	OnTimeGetsCritical(time.Time)
	OnNextParticipantStarts(time.Time, Participant)
	OnTimeIsOver(time.Time)
	OnStartsWaitingNextParticipant(time.Time)
	NextParticipant() Participant
	CurrentParticipant() Participant
	TurnTimeInfo() *TurnTimeInfo
	NextParticipantIsReady() bool
}

type TimerInternalStateLabel int

const (
	STATE_WAITING_NEXT_PARTICIPANT TimerInternalStateLabel = iota
	STATE_TIME_IS_OK               TimerInternalStateLabel = iota
	STATE_TIME_IS_CRITICAL         TimerInternalStateLabel = iota
	STATE_TIME_IS_OVER             TimerInternalStateLabel = iota
)

type ITimerIntenalState interface {
	ChangeToState(ITurnLogic, time.Time) TimerInternalStateLabel
}

type StatesMap map[TimerInternalStateLabel]ITimerIntenalState

type Timer struct {
	TurnLogic         ITurnLogic
	CurrentStateLabel TimerInternalStateLabel
	States            StatesMap
}

func (timer *Timer) CurrentState() ITimerIntenalState {
	return timer.States[timer.CurrentStateLabel]
}

func (timer *Timer) Step(time time.Time) {
	currentState := timer.CurrentState()
	timer.CurrentStateLabel = currentState.ChangeToState(timer.TurnLogic, time)
}

// Implementing states
type TimerStartedStated struct {
}

type TimerWaitingNextParticipant struct {
}

type TimerTimeIsOk struct {
	Begin time.Time
}

type TimerTimeIsCritical struct {
	Begin time.Time
}

type TimerTimeIsOver struct {
}

func NewTimer(logic ITurnLogic) *Timer {
	timer := &Timer{logic, STATE_WAITING_NEXT_PARTICIPANT, StatesMap{}}

	timer.States[STATE_WAITING_NEXT_PARTICIPANT] = &TimerWaitingNextParticipant{}
	timer.States[STATE_TIME_IS_OK] = &TimerTimeIsOk{}
	timer.States[STATE_TIME_IS_CRITICAL] = &TimerTimeIsCritical{}
	timer.States[STATE_TIME_IS_OVER] = &TimerTimeIsOver{}

	return timer
}

func (state *TimerStartedStated) ChangeToState(logic ITurnLogic, time time.Time) TimerInternalStateLabel {
	logic.OnStartsWaitingNextParticipant(time)
	return STATE_WAITING_NEXT_PARTICIPANT
}

func (state *TimerWaitingNextParticipant) ChangeToState(logic ITurnLogic, time time.Time) TimerInternalStateLabel {
	participant := logic.CurrentParticipant()

	if participant.Valid() {
		logic.OnNextParticipantStarts(time, logic.CurrentParticipant())
		return STATE_TIME_IS_OK
	}

	return STATE_WAITING_NEXT_PARTICIPANT
}

func (state *TimerTimeIsOk) ChangeToState(logic ITurnLogic, time time.Time) TimerInternalStateLabel {
	if state.Begin.IsZero() {
		fmt.Printf("ok started in %s\n", time)
		state.Begin = time
	}

	if state.Begin.Add(logic.TurnTimeInfo().RelaxAndCodeDuration).After(time) {
		logic.OnTimeGetsCritical(time)
		return STATE_TIME_IS_CRITICAL
	}

	return STATE_TIME_IS_OK
}

func (state *TimerTimeIsCritical) ChangeToState(logic ITurnLogic, time time.Time) TimerInternalStateLabel {
	if state.Begin.IsZero() {
		fmt.Printf("critical started in %s\n", time)
		state.Begin = time
	}

	if state.Begin.Add(logic.TurnTimeInfo().HurryUpDuration).After(time) {
		logic.OnTimeIsOver(time)
		return STATE_TIME_IS_OVER
	}

	return STATE_TIME_IS_CRITICAL
}

func (state *TimerTimeIsOver) ChangeToState(logic ITurnLogic, time time.Time) TimerInternalStateLabel {
	if logic.NextParticipantIsReady() {
		return STATE_WAITING_NEXT_PARTICIPANT
	}

	return STATE_TIME_IS_OVER
}
