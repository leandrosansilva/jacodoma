package main

import (
	. "./src"
	"fmt"
	"gopkg.in/qml.v1"
	"time"
)

type TurnInformation struct {
	Info               TurnTimeInfo
	Participants       Participants
	Index              int
	Ready              bool
	State              chan string
	ParticipantChannel chan Participant
}

func (info *TurnInformation) ChangeToNextParticipant() Participant {
	p := info.Participants.Get(info.Index)
	info.Index = (info.Index + 1) % info.Participants.Length()
	return p
}

func (info *TurnInformation) NextParticipant() Participant {
	// FIXME: this "rotational" logic should be in Participants{}
	return info.Participants.Get((info.Index + 1) % info.Participants.Length())
}

func (info *TurnInformation) HurryUp() {
	info.State <- "hurry_up"
}

func (info *TurnInformation) TimeIsOver() {
	info.State <- "time_over"
}

func (info *TurnInformation) ParticipantStarts() {
	info.State <- "start"
}

func (info *TurnInformation) StartsWaitingNextParticipant() {
	fmt.Println("waiting for the next participant")
	info.ParticipantChannel <- info.NextParticipant()
}

type TurnLogic struct {
	info *TurnInformation
}

func (logic *TurnLogic) OnTimeGetsCritical(t time.Time) {
	logic.info.HurryUp()
}

func (logic *TurnLogic) OnNextParticipantStarts(t time.Time, p Participant) {
	logic.info.ParticipantStarts()
}

func (logic *TurnLogic) OnTimeIsOver(t time.Time) {
	logic.info.TimeIsOver()
}

func (logic *TurnLogic) OnStartsWaitingNextParticipant(t time.Time) {
	logic.info.StartsWaitingNextParticipant()
}

func (logic *TurnLogic) BlockSession(t time.Time) {
	fmt.Println("Session blocked until the next person comes!")
}

func (logic *TurnLogic) NextParticipantIsReady() bool {
	// FIXME: data race condition!
	return logic.info.Ready
}

func (logic *TurnLogic) NextParticipant() Participant {
	return logic.info.ChangeToNextParticipant()
}

func (logic *TurnLogic) TurnTimeInfo() *TurnTimeInfo {
	return &logic.info.Info
}

// Acts as model to the GUI
type Control struct {
	Info        *TurnInformation
	Duration    int64
	State       string
	Participant Participant
}

func (this *Control) SetParticipantReady() {
	this.Info.Ready = true

	// FIXME: workaround (with data race condition)
	go func() {
		time.Sleep(1 * time.Second)
		this.Info.Ready = false
	}()
}

type QmlGui struct {
	info    *TurnInformation
	channel DurationChannel
	ctrl    *Control
}

func (this *QmlGui) Run() error {
	setup := func() error {
		engine := qml.NewEngine()

		component, err := engine.LoadFile("main.qml")

		if err != nil {
			return err
		}

		engine.Context().SetVar("ctrl", this.ctrl)

		// timer ui loop
		go func() {
			for {
				d := <-this.channel
				this.ctrl.Duration = int64(d)
				qml.Changed(this.ctrl, &this.ctrl.Duration)
			}
		}()

		// turn state (ok, hurry up, time is over...) ui loop
		go func() {
			for {
				this.ctrl.State = <-this.info.State
				qml.Changed(this.ctrl, &this.ctrl.State)
			}
		}()

		// participant ui loop
		go func() {
			for {
				this.ctrl.Participant = <-this.info.ParticipantChannel
				qml.Changed(this.ctrl, &this.ctrl.Participant)
			}
		}()

		window := component.CreateWindow(nil)

		window.Show()
		window.Wait()

		return nil
	}

	return qml.Run(setup)
}

func NewQmlGui(info *TurnInformation, channel DurationChannel) *QmlGui {
	return &QmlGui{info, channel, &Control{info, 0, "", Participant{}}}
}

func main() {
	participants, _ := LoadParticipantsFromFile("users.jcdm")

	turnInfo := TurnTimeInfo{20 * time.Second, 10 * time.Second}

	info := &TurnInformation{turnInfo, participants, 0, false, make(chan string), make(chan Participant)}

	logic := &TurnLogic{info}

	channel := make(DurationChannel, 0)

	gui := NewQmlGui(info, channel)

	timer := NewTimer(logic, channel)

	timer.Step(time.Time{})

	// ticker loop
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for t := range ticker.C {
			timer.Step(t)
		}
	}()

	if err := gui.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Println("Exiting...")
}
