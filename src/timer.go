package jacodoma

import (
	"time"
)

type ITimerLogic interface {
	OnCriticalTimeHasReached(t time.Time)
	OnStarted(t time.Time)
	OnTimeGetsCritical(t time.Time)
	OnNextParticipantStarts(t time.Time, p Participant)
	OnTimeIsOver(t time.Time)
	OnStartsWaitingNextParticipant(t time.Time)
	NextParticipant() Participant
	CurrentParticipant() Participant
	TimeChannel() chan time.Time
	HasFinished() bool
}

type TurnTimeInfo struct {
	TurnTimeInSeconds     int
	CriticalTimeInSeconds int
}

type Timer struct {
	logic ITimerLogic
}

func (t *Timer) Run() {
}

func NewTimer(logic ITimerLogic) *Timer {
	return &Timer{logic}
}
