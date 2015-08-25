package main

import (
	. "./src"
	"bufio"
	"fmt"
	"gopkg.in/qml.v1"
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

func (logic *TurnLogic) OnNextParticipantStarts(t time.Time, p Participant) {
	fmt.Printf("%s starts!\n", p.Name)
}

func (logic *TurnLogic) OnTimeIsOver(t time.Time) {
	fmt.Println("Timeout :-(!")
}

func (logic *TurnLogic) OnStartsWaitingNextParticipant(t time.Time) {
	fmt.Println("Waiting for the next participant...")
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

type QmlGui struct {
	logic   *TurnLogic
	channel DurationChannel
}

func (this *QmlGui) Run() error {
	setup := func() error {
		engine := qml.NewEngine()

		component, err := engine.LoadFile("main.qml")

		if err != nil {
			return err
		}

		window := component.CreateWindow(nil)

		//root := window.Root()
		//root.On("lala", func(request *qml.Object) {
		//})

		/*var duration struct {
		      duration time.Duration
		    }

		    engine.Context().SetVar("timer_time", duration)

				go func() {
					for {
						duration.duration = <-this.channel
						qml.Changed(duration, &duration.duration)
					}
				}()*/

		window.Show()
		window.Wait()

		return nil
	}

	return qml.Run(setup)
}

func NewQmlGui(logic *TurnLogic, channel DurationChannel) *QmlGui {
	return &QmlGui{logic, channel}
}

func main() {
	participants, _ := LoadParticipantsFromFile("users.jcdm")

	turnInfo := TurnTimeInfo{10 * time.Second, 5 * time.Second}

	logic := &TurnLogic{turnInfo, participants, 0, false}

	channel := make(DurationChannel, 0)

	timer := NewTimer(logic, channel)

	timer.Step(time.Time{})

	ticker := time.NewTicker(100 * time.Millisecond)

	// user input loop
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			// FIXME: data race!
			reader.ReadString('\n')
			logic.Ready = true
			time.Sleep(1 * time.Second)
			logic.Ready = false
		}
	}()

	// ui loop

	// ticker loop
	go func() {
		for t := range ticker.C {
			timer.Step(t)
		}
	}()

	gui := NewQmlGui(logic, channel)

	if err := gui.Run(); err == nil {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Println("Exiting...")
}
