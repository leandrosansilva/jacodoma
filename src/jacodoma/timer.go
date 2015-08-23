package jacodoma

import (
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
	TurnTimeInfo() *TurnTimeInfo
	NextParticipantIsReady() bool
	BlockSession(time.Time)
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
type TimerStartedState struct {
}

type TimerWaitingNextParticipant struct {
}

type TimerTimeIsOkState struct {
	Begin time.Time
}

type TimerTimeIsCriticalState struct {
	Begin time.Time
}

type TimerTimeIsOverState struct {
}

func NewTimer(logic ITurnLogic) *Timer {
	timer := &Timer{logic, STATE_WAITING_NEXT_PARTICIPANT, StatesMap{
		STATE_WAITING_NEXT_PARTICIPANT: &TimerWaitingNextParticipant{},
		STATE_TIME_IS_OK:               &TimerTimeIsOkState{},
		STATE_TIME_IS_CRITICAL:         &TimerTimeIsCriticalState{},
		STATE_TIME_IS_OVER:             &TimerTimeIsOverState{},
	}}

	return timer
}

func (this *TimerStartedState) ChangeToState(logic ITurnLogic, time time.Time) TimerInternalStateLabel {
	logic.OnStartsWaitingNextParticipant(time)
	return STATE_WAITING_NEXT_PARTICIPANT
}

func (this *TimerWaitingNextParticipant) ChangeToState(logic ITurnLogic, time time.Time) TimerInternalStateLabel {
	if logic.NextParticipantIsReady() {
		logic.OnNextParticipantStarts(time, logic.NextParticipant())
		return STATE_TIME_IS_OK
	}

	return STATE_WAITING_NEXT_PARTICIPANT
}

func (this *TimerTimeIsOkState) ChangeToState(logic ITurnLogic, t time.Time) TimerInternalStateLabel {
	if this.Begin.IsZero() {
		this.Begin = t
	}

	if this.Begin.Add(logic.TurnTimeInfo().RelaxAndCodeDuration).Before(t) {
		logic.OnTimeGetsCritical(t)
		this.Begin = time.Time{}
		return STATE_TIME_IS_CRITICAL
	}

	return STATE_TIME_IS_OK
}

func (this *TimerTimeIsCriticalState) ChangeToState(logic ITurnLogic, t time.Time) TimerInternalStateLabel {
	if this.Begin.IsZero() {
		this.Begin = t
	}

	if this.Begin.Add(logic.TurnTimeInfo().HurryUpDuration).Before(t) {
		logic.OnTimeIsOver(t)
		this.Begin = time.Time{}
		return STATE_TIME_IS_OVER
	}

	return STATE_TIME_IS_CRITICAL
}

func (this *TimerTimeIsOverState) ChangeToState(logic ITurnLogic, time time.Time) TimerInternalStateLabel {
	logic.BlockSession(time)
	return STATE_WAITING_NEXT_PARTICIPANT
}
