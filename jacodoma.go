package main

import (
	. "./src"
	"bytes"
	"fmt"
	"github.com/ftrvxmtrx/gravatar"
	"gopkg.in/qml.v1"
	"image"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"time"
)

type TurnInformation struct {
	Info               TurnTimeInfo
	Participants       Participants
	Index              int
	Ready              bool
	State              chan string
	ParticipantChannel chan int
}

func (info *TurnInformation) ChangeToNextParticipantIndex() int {
	index := info.Index
	info.Index = (info.Index + 1) % info.Participants.Length()
	return index
}

func (info *TurnInformation) NextParticipantIndex() int {
	// FIXME: this "rotational" logic should be in Participants{}
	return (info.Index + 1) % info.Participants.Length()
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

func (info *TurnInformation) StartsWaitingNextParticipant(index int) {
	info.State <- "waiting_participant"
	info.ParticipantChannel <- index
}

func (info *TurnInformation) TotalTurnTime() int64 {
	return int64(info.Info.RelaxAndCodeDuration + info.Info.HurryUpDuration)
}

func (info *TurnInformation) HurryUpDuration() int64 {
	return int64(info.Info.HurryUpDuration)
}

func (info *TurnInformation) RelaxAndCodeDuration() int64 {
	return int64(info.Info.RelaxAndCodeDuration)
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

func (logic *TurnLogic) OnStartsWaitingNextParticipant(t time.Time, index int) {
	logic.info.StartsWaitingNextParticipant(index)
}

func (logic *TurnLogic) BlockSession(t time.Time) {
	fmt.Println("Session blocked until the next person comes!")
}

func (logic *TurnLogic) NextParticipantIsReady() bool {
	// FIXME: data race condition!
	return logic.info.Ready
}

func (logic *TurnLogic) NextParticipantIndex() int {
	return logic.info.ChangeToNextParticipantIndex()
}

func (logic *TurnLogic) TurnTimeInfo() *TurnTimeInfo {
	return &logic.info.Info
}

// Acts as model to the GUI
type Control struct {
	Info                    *TurnInformation
	TurnDuration            int64
	SessionDuration         int64
	State                   string
	CurrentParticipantIndex int
	Participants            *Participants
	ParticipantsLen         int
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

func gravatarImageProvider(email string, width, height int) image.Image {
	emailHash := gravatar.EmailHash(email)

	raw, err := gravatar.GetAvatar("https", emailHash, gravatar.DefaultMonster, width)

	if err != nil {
		panic(err)
	}

	if img, err := png.Decode(bytes.NewReader(raw)); err == nil {
		return img
	}

	if img, err := jpeg.Decode(bytes.NewReader(raw)); err == nil {
		return img
	}

	return nil
}

func (this *QmlGui) Run() error {
	return qml.Run(func() error {
		engine := qml.NewEngine()

		engine.AddImageProvider("gravatar", gravatarImageProvider)

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
				case this.ctrl.CurrentParticipantIndex = <-this.info.ParticipantChannel:
					qml.Changed(this.ctrl, &this.ctrl.CurrentParticipantIndex)
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
	control := &Control{info, 0, 0, "", 0, &info.Participants, info.Participants.Length()}
	return &QmlGui{info, turnTimeChannel, sessionTimeChannel, control}
}

func main() {
	participants, err := LoadParticipantsFromFile("users.jcdm")

	if err != nil {
		fmt.Printf("Error loading participants file: %s\n", err)
		os.Exit(1)
	}

	if participants.Length() == 0 {
		fmt.Printf("There is no participants :-(")
		os.Exit(1)
	}

	config, err := LoadProjectConfigFile("config.jcdm")

	if err != nil {
		fmt.Printf("Error loading config file: %s\n", err)
		os.Exit(1)
	}

	rand.Seed(time.Now().Unix())

	if config.Session.ShuffleUsersOrder {
		participants.Shuffle()
	}

	turnInfo := TurnTimeInfo{time.Duration(config.Session.TurnTime), time.Duration(config.Session.Critical)}

	info := &TurnInformation{
		turnInfo,
		participants, 0, false,
		make(chan string),
		make(chan int)}

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
