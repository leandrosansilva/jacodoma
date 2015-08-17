package jacodoma

import (
	"time"
)

type ITurnLogic interface {
	OnStarted(t time.Time)
	OnTimeGetsCritical(t time.Time)
	OnNextParticipantStarts(t time.Time, p Participant)
	OnTimeIsOver(t time.Time)
	OnStartsWaitingNextParticipant(t time.Time)
	NextParticipant() Participant
	CurrentParticipant() Participant
	HasFinished() bool
}

type TurnTimeInfo struct {
	TurnTimeInSeconds     int
	CriticalTimeInSeconds int
}

type Timer struct {
	logic ITurnLogic
}

func (timer *Timer) Step(time time.Time) {
}

func NewTimer(logic ITurnLogic) *Timer {
	return &Timer{logic}
}
