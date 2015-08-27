package main

import (
	. "./src/jacodoma"
	"bufio"
	"fmt"
	"os"
	"time"
)

func init() {
}

type TurnLogic struct {
	info         TurnTimeInfo
	Participants Participants
	Index        int
	Ready        bool
}

func (logic *TurnLogic) OnTimeGetsCritical(t time.Time) {
	fmt.Println("Hurry up dude!")
}

func (logic *TurnLogic) OnNextParticipantStarts(t time.Time) {
	fmt.Println("Participant starts!")
}

func (logic *TurnLogic) OnTimeIsOver(t time.Time) {
	fmt.Println("Timeout :-(!")
}

func (logic *TurnLogic) OnStartsWaitingNextParticipant(t time.Time, p Participant) {
	fmt.Printf("Waiting for the next participant %s\n", p.Name)
}

func (logic *TurnLogic) BlockSession(t time.Time) {
	fmt.Println("Session blocked until the next person comes!")
}

func (logic *TurnLogic) NextParticipantIsReady() bool {
	// FIXME: data race!
	return logic.Ready
}

func (logic *TurnLogic) NextParticipant() Participant {
	p := logic.Participants.Get(logic.Index)
	logic.Index = (logic.Index + 1) % logic.Participants.Length()
	return p
}

func (logic *TurnLogic) TurnTimeInfo() *TurnTimeInfo {
	return &logic.info
}

func main() {
	participants, _ := LoadParticipantsFromFile("users.jcdm")

	turnInfo := TurnTimeInfo{10 * time.Second, 5 * time.Second}

	logic := &TurnLogic{turnInfo, participants, 0, false}

	channel := make(DurationChannel, 0)

	timer := NewTimer(logic, channel)

	timer.Step(time.Time{})

	ticker := time.NewTicker(100 * time.Millisecond)

	// ui loop
	go func() {
		for {
			d := <-channel
			r := turnInfo.RelaxAndCodeDuration + turnInfo.HurryUpDuration - d
			fmt.Printf("time: %s\n", r)
		}
	}()

	// ticker loop
	go func() {
		for t := range ticker.C {
			timer.Step(t)
		}
	}()

	// user input loop
	reader := bufio.NewReader(os.Stdin)
	for {
		// FIXME: data race!
		reader.ReadString('\n')
		logic.Ready = true
		time.Sleep(1 * time.Second)
		logic.Ready = false
	}

}
