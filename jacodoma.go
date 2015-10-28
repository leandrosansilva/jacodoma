package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ftrvxmtrx/gravatar"
	"gopkg.in/qml.v1"
	"image"
	"image/jpeg"
	"image/png"
	. "jacodoma/src"
	"math/rand"
	"os"
	"path"
	"time"
)

type TurnControl struct {
	Info                    TurnTimeInfo
	Participants            Participants
	Index                   int
	Ready                   bool
	State                   chan string
	ParticipantIndexChannel chan int
	CommitChannel           chan Participant
}

func NewTurnControl(info TurnTimeInfo, participants Participants) *TurnControl {
	return &TurnControl{
		info,
		participants, 0, false,
		make(chan string),
		make(chan int),
		make(chan Participant),
	}
}

func modCalc(index, size, offset int) int {
	i := (index + size) % size

	if i < 0 {
		return i + size
	}

	return i
}

func (info *TurnControl) NextParticipantIndex() int {
	return modCalc(info.Index, info.Participants.Length(), 1)
}

func (info *TurnControl) HurryUp() {
	info.State <- "hurry_up"
}

func (info *TurnControl) TimeIsOver() {
	// FIXME: extremely dirty workaround due design errors
	index := modCalc(info.Index, info.Participants.Length(), -1)
	info.CommitChannel <- info.Participants.Get(index)

	info.State <- "time_over"
}

func (info *TurnControl) ParticipantStarts() {
	info.State <- "start"
}

func (info *TurnControl) StartsWaitingNextParticipant(index int) {
	info.State <- "waiting_participant"

	info.ParticipantIndexChannel <- info.Index

	info.Index = index
}

// func() -> int64 are wrappers to QML
func (info *TurnControl) TotalTurnTime() int64 {
	return int64(info.Info.RelaxAndCodeDuration + info.Info.HurryUpDuration)
}

func (info *TurnControl) HurryUpDuration() int64 {
	return int64(info.Info.HurryUpDuration)
}

func (info *TurnControl) RelaxAndCodeDuration() int64 {
	return int64(info.Info.RelaxAndCodeDuration)
}

type TurnLogic struct {
	info *TurnControl
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
	return logic.info.NextParticipantIndex()
}

func (logic *TurnLogic) TurnTimeInfo() *TurnTimeInfo {
	return &logic.info.Info
}

// Acts as model to the GUI
type Control struct {
	Info                    *TurnControl
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
	info               *TurnControl
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

func (this *QmlGui) Run(config *ProjectConfig) error {
	return qml.Run(func() error {
		engine := qml.NewEngine()

		engine.AddImageProvider("gravatar", gravatarImageProvider)

		component, err := engine.LoadFile(config.UI.Skin)

		if err != nil {
			return err
		}

		engine.Context().SetVar("ctrl", this.ctrl)

		engine.Context().SetVar("config", config)

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
				case this.ctrl.CurrentParticipantIndex = <-this.info.ParticipantIndexChannel:
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

func NewTurnTimeInfo(turnTime, criticalTime time.Duration) TurnTimeInfo {
	return TurnTimeInfo{turnTime - criticalTime, criticalTime}
}

func NewQmlGui(info *TurnControl, turnTimeChannel, sessionTimeChannel DurationChannel) *QmlGui {
	control := &Control{info, 0, 0, "", 0, &info.Participants, info.Participants.Length()}
	return &QmlGui{info, turnTimeChannel, sessionTimeChannel, control}
}

var (
	projectDirectory string
	showHelp         bool
)

func init() {
	flag.BoolVar(&showHelp, "help", false, "Show this usage message")
	flag.StringVar(&projectDirectory, "project", "", "Project directory (default to current directory)")
}

func parseCmdlineParams() {
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(2)
	}

	if len(projectDirectory) > 0 {
		return
	}

	var err error

	if projectDirectory, err = os.Getwd(); err != nil {
		fmt.Println("Invalid project directory")
		os.Exit(1)
	}
}

func main() {
	var config ProjectConfig
	var participants Participants
	var err error
	var repository Repository

	parseCmdlineParams()

	rand.Seed(time.Now().Unix())

	if config, err = LoadProjectConfigFile(path.Join(projectDirectory, "config.jcdm")); err != nil {
		fmt.Printf("Error loading config file: %s\n", err)
		os.Exit(1)
	}

	if participants, err = LoadParticipantsFromFile(path.Join(projectDirectory, "users.jcdm")); err != nil {
		fmt.Printf("Error loading participants file: %s\n", err)
		os.Exit(1)
	}

	if participants.Length() == 0 {
		fmt.Printf("There is no participants :-(")
		os.Exit(1)
	}

	if repository, err = CreateVcsRepository(projectDirectory); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if config.Session.ShuffleUsersOrder {
		participants.Shuffle()
	}

	turnInfo := NewTurnTimeInfo(time.Duration(config.Session.TurnTime), time.Duration(config.Session.Critical))

	control := NewTurnControl(turnInfo, participants)

	// Repository loop
	go func() {
		for {
			participant := <-control.CommitChannel

			meta := CreateCommitMetadata(participant.Name, participant.Email, time.Now())

			if err := repository.CommitFiles(config.Project.SourceFiles, meta); err != nil {
				fmt.Println(err)
			}
		}
	}()

	logic := &TurnLogic{control}

	turnTimeChannel := make(DurationChannel, 0)
	sessionTimeChannel := make(DurationChannel, 0)

	timer := NewTimer(logic, turnTimeChannel, sessionTimeChannel)

	gui := NewQmlGui(control, turnTimeChannel, sessionTimeChannel)

	// timer ticker loop
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for t := range ticker.C {
			timer.Step(t)
		}
	}()

	if err := gui.Run(&config); err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Println("Exiting...")
}
