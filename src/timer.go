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

type DurationChannel chan time.Duration

type StatesMap map[TimerInternalStateLabel]ITimerIntenalState

type TurnContext struct {
	Begin        time.Time
	CriticalTime time.Time
	LastDuration time.Duration
}

func (this *TurnContext) SetTurnBeginIfNotDefined(t time.Time) {
	if this.Begin.IsZero() {
		this.Begin = t
	}
}

func (this *TurnContext) SetCriticalBeginIfNotDefined(t time.Time) {
	if this.CriticalTime.IsZero() {
		this.CriticalTime = t
		this.LastDuration = 0
	}
}

func (this *TurnContext) Reset() {
	this.CriticalTime = time.Time{}
	this.Begin = time.Time{}
}

func (this *TurnContext) Update(channel DurationChannel, t time.Time) {
	if this.Begin.IsZero() {
		return
	}

	// time elapsed since the turn begin
	d := t.Sub(this.Begin)

	s := d / time.Second

	if s != this.LastDuration/time.Second || this.LastDuration == 0 {
		this.LastDuration = d
		channel <- d
	}
}

type Timer struct {
	Context           *TurnContext
	TurnLogic         ITurnLogic
	CurrentStateLabel TimerInternalStateLabel
	States            StatesMap
	DurationChannel   DurationChannel
}

func (timer *Timer) CurrentState() ITimerIntenalState {
	return timer.States[timer.CurrentStateLabel]
}

func (timer *Timer) Step(t time.Time) {
	currentState := timer.CurrentState()
	timer.CurrentStateLabel = currentState.ChangeToState(timer.TurnLogic, t)
	timer.Context.Update(timer.DurationChannel, t)
}

// Implementing states
type TimerStartedState struct {
}

type TimerWaitingNextParticipant struct {
}

type TimerTimeIsOkState struct {
	Context *TurnContext
}

type TimerTimeIsCriticalState struct {
	Context *TurnContext
}

type TimerTimeIsOverState struct {
}

func NewTimer(logic ITurnLogic, channel DurationChannel) *Timer {
	context := &TurnContext{}
	timer := &Timer{context, logic, STATE_WAITING_NEXT_PARTICIPANT, StatesMap{
		STATE_WAITING_NEXT_PARTICIPANT: &TimerWaitingNextParticipant{},
		STATE_TIME_IS_OK:               &TimerTimeIsOkState{context},
		STATE_TIME_IS_CRITICAL:         &TimerTimeIsCriticalState{context},
		STATE_TIME_IS_OVER:             &TimerTimeIsOverState{},
	}, channel}

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
	this.Context.SetTurnBeginIfNotDefined(t)

	if this.Context.Begin.Add(logic.TurnTimeInfo().RelaxAndCodeDuration).Before(t) {
		logic.OnTimeGetsCritical(t)
		return STATE_TIME_IS_CRITICAL
	}

	return STATE_TIME_IS_OK
}

func (this *TimerTimeIsCriticalState) ChangeToState(logic ITurnLogic, t time.Time) TimerInternalStateLabel {
	this.Context.SetCriticalBeginIfNotDefined(t)

	if this.Context.CriticalTime.Add(logic.TurnTimeInfo().HurryUpDuration).Before(t) {
		logic.OnTimeIsOver(t)
		this.Context.Reset()
		return STATE_TIME_IS_OVER
	}

	return STATE_TIME_IS_CRITICAL
}

func (this *TimerTimeIsOverState) ChangeToState(logic ITurnLogic, time time.Time) TimerInternalStateLabel {
	logic.BlockSession(time)
	return STATE_WAITING_NEXT_PARTICIPANT
}
