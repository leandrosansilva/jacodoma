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
	OnNextParticipantStarts(time.Time)
	OnTimeIsOver(time.Time)
	OnStartsWaitingNextParticipant(time.Time, int)
	NextParticipantIndex() int
	TurnTimeInfo() *TurnTimeInfo
	NextParticipantIsReady() bool
	BlockSession(time.Time)
}

type TimerInternalStateLabel int

const (
	STATE_INITIAL                  TimerInternalStateLabel = iota
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
	SessionBegin        time.Time
	TurnBegin           time.Time
	CriticalTime        time.Time
	LastTurnDuration    time.Duration
	LastSessionDuration time.Duration
}

func (this *TurnContext) SetTurnBeginIfNotDefined(t time.Time) {
	if this.TurnBegin.IsZero() {
		this.TurnBegin = t
	}
}

func (this *TurnContext) SetCriticalBeginIfNotDefined(t time.Time) {
	if this.CriticalTime.IsZero() {
		this.CriticalTime = t
		this.LastTurnDuration = 0
	}
}

func (this *TurnContext) Reset() {
	this.CriticalTime = time.Time{}
	this.TurnBegin = time.Time{}
}

func tryToUpdateTimeAndSendToChannel(channel DurationChannel, begin time.Time, t time.Time, duration time.Duration) time.Duration {
	d := t.Sub(begin)

	seconds := d / time.Second

	if seconds != duration/time.Second || duration == 0 {
		channel <- d
	}

	return d
}

func (this *TurnContext) UpdateTurnTime(turnChannel DurationChannel, t time.Time) {
	if this.TurnBegin.IsZero() {
		return
	}

	this.LastTurnDuration = tryToUpdateTimeAndSendToChannel(
		turnChannel, this.TurnBegin, t, this.LastTurnDuration)
}

func (this *TurnContext) UpdateSessionTime(sessionChannel DurationChannel, t time.Time) {
	if this.SessionBegin.IsZero() {
		this.SessionBegin = t
	}

	this.LastSessionDuration = tryToUpdateTimeAndSendToChannel(
		sessionChannel, this.SessionBegin, t, this.LastSessionDuration)
}

func (this *TurnContext) Update(turnChannel DurationChannel, sessionChannel DurationChannel, t time.Time) {
	this.UpdateSessionTime(sessionChannel, t)
	this.UpdateTurnTime(turnChannel, t)
}

type Timer struct {
	Context                *TurnContext
	TurnLogic              ITurnLogic
	CurrentStateLabel      TimerInternalStateLabel
	States                 StatesMap
	TurnDurationChannel    DurationChannel
	SessionDurationChannel DurationChannel
}

func (timer *Timer) CurrentState() ITimerIntenalState {
	return timer.States[timer.CurrentStateLabel]
}

func (timer *Timer) Step(t time.Time) {
	currentState := timer.CurrentState()
	timer.CurrentStateLabel = currentState.ChangeToState(timer.TurnLogic, t)
	timer.Context.Update(timer.TurnDurationChannel, timer.SessionDurationChannel, t)
}

// Implementing states
type TimerInitialState struct {
}

type TimerStartedState struct {
}

type TimerWaitingNextParticipantState struct {
}

type TimerTimeIsOkState struct {
	Context *TurnContext
}

type TimerTimeIsCriticalState struct {
	Context *TurnContext
}

type TimerTimeIsOverState struct {
}

func NewTimer(logic ITurnLogic, turnChannel DurationChannel, sessionChannel DurationChannel) *Timer {
	context := &TurnContext{}
	timer := &Timer{context, logic, STATE_INITIAL, StatesMap{
		STATE_INITIAL:                  &TimerInitialState{},
		STATE_WAITING_NEXT_PARTICIPANT: &TimerWaitingNextParticipantState{},
		STATE_TIME_IS_OK:               &TimerTimeIsOkState{context},
		STATE_TIME_IS_CRITICAL:         &TimerTimeIsCriticalState{context},
		STATE_TIME_IS_OVER:             &TimerTimeIsOverState{},
	}, turnChannel, sessionChannel}

	return timer
}

func (this *TimerInitialState) ChangeToState(logic ITurnLogic, time time.Time) TimerInternalStateLabel {
	logic.OnStartsWaitingNextParticipant(time, logic.NextParticipantIndex())
	return STATE_WAITING_NEXT_PARTICIPANT
}

func (this *TimerWaitingNextParticipantState) ChangeToState(logic ITurnLogic, time time.Time) TimerInternalStateLabel {
	if logic.NextParticipantIsReady() {
		logic.OnNextParticipantStarts(time)
		return STATE_TIME_IS_OK
	}

	return STATE_WAITING_NEXT_PARTICIPANT
}

func (this *TimerTimeIsOkState) ChangeToState(logic ITurnLogic, t time.Time) TimerInternalStateLabel {
	this.Context.SetTurnBeginIfNotDefined(t)

	if this.Context.TurnBegin.Add(logic.TurnTimeInfo().RelaxAndCodeDuration).Before(t) {
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
	logic.OnStartsWaitingNextParticipant(time, logic.NextParticipantIndex())
	logic.BlockSession(time)

	return STATE_WAITING_NEXT_PARTICIPANT
}
