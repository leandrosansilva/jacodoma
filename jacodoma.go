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

func (info *TurnInformation) StartsWaitingNextParticipant(p Participant) {
	info.State <- "waiting_participant"
	info.ParticipantChannel <- p
}

type TurnLogic struct {
	info *TurnInformation
}

func (logic *TurnLogic) OnTimeGetsCritical(t time.Time) {
	logic.info.HurryUp()
}

func (logic *TurnLogic) OnNextParticipantStarts(t time.Time) {
	logic.info.ParticipantStarts()
}

func (logic *TurnLogic) OnTimeIsOver(t time.Time) {
	logic.info.TimeIsOver()
}

func (logic *TurnLogic) OnStartsWaitingNextParticipant(t time.Time, p Participant) {
	logic.info.StartsWaitingNextParticipant(p)
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
	Info            *TurnInformation
	TurnDuration    int64
	SessionDuration int64
	State           string
	Participant     Participant
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
	info               *TurnInformation
	turnTimeChannel    DurationChannel
	sessionTimeChannel DurationChannel
	ctrl               *Control
}

func (this *QmlGui) Run() error {
	return qml.Run(func() error {
		engine := qml.NewEngine()

		component, err := engine.LoadFile("main.qml")

		if err != nil {
			return err
		}

		engine.Context().SetVar("ctrl", this.ctrl)

		// ui loop
		go func() {
			for {
				select {
				case d := <-this.turnTimeChannel:
					this.ctrl.TurnDuration = int64(d)
					qml.Changed(this.ctrl, &this.ctrl.TurnDuration)
				case d := <-this.sessionTimeChannel:
					this.ctrl.SessionDuration = int64(d)
					qml.Changed(this.ctrl, &this.ctrl.SessionDuration)
				case this.ctrl.State = <-this.info.State:
					qml.Changed(this.ctrl, &this.ctrl.State)
				case this.ctrl.Participant = <-this.info.ParticipantChannel:
					qml.Changed(this.ctrl, &this.ctrl.Participant)
				}
			}
		}()

		window := component.CreateWindow(nil)

		window.Show()
		window.Wait()

		return nil
	})
}

func NewQmlGui(info *TurnInformation, turnTimeChannel, sessionTimeChannel DurationChannel) *QmlGui {
	control := &Control{info, 0, 0, "", Participant{}}
	return &QmlGui{info, turnTimeChannel, sessionTimeChannel, control}
}

func main() {
	participants, _ := LoadParticipantsFromFile("users.jcdm")

	turnInfo := TurnTimeInfo{20 * time.Second, 10 * time.Second}

	info := &TurnInformation{
		turnInfo,
		participants, 0, false,
		make(chan string),
		make(chan Participant)}

	logic := &TurnLogic{info}

	turnTimeChannel := make(DurationChannel, 0)
	sessionTimeChannel := make(DurationChannel, 0)

	timer := NewTimer(logic, turnTimeChannel, sessionTimeChannel)

	gui := NewQmlGui(info, turnTimeChannel, sessionTimeChannel)

	// timer ticker loop
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
