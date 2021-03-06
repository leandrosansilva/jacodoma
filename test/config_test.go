package jacodoma

import (
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	. "jacodoma/src"
	"os"
	"testing"
	"time"
)

func TestCompleteConfigLoading(t *testing.T) {
	completeFileContent := `
[Session]
ExerciseReference = http://problems.example.com/puzzle-1
ExerciseReference = http://problems.example2.com/puzzle-2
NotifyBadBehaviour = true
TurnTime = 5m
Critical = 1m
UseSoundNotification = true
UseSystemNotification = true
SoundNotificationFilename = beep.ogg
LockScreenOnTimeout = true
ShuffleUsersOrder = true
  
[Tests]
Command = go test
OnEveryChange = false
OnTimeout = 10s
Files = tests/*.go

[Project]
VC = Git
CommitOnEveryChange = On
SourceFiles = src/**.go

[Report]
DbFile = db.jcdmdb   

[UI]
Type = QML
Skin = Default`

	incompleteFileContent := `
[Session]
ExerciseReference = http://problems.example3.com/puzzle-3
TurnTime = 4m
  
[Project]
VC = Git
CommitOnEveryChange = On
SourceFiles = src/**.go`

	Convey("Load complete file", t, func() {
		ioutil.WriteFile("/tmp/jacodoma_jcdmarc.txt_test", []byte(completeFileContent), 0644)
		defer os.Remove("/tmp/jacodoma_jcdmarc.txt_test")

		config, err := LoadProjectConfigFile("/tmp/jacodoma_jcdmarc.txt_test")

		So(err, should.Equal, nil)
		So(config, should.NotEqual, ProjectConfig{})

		So(config.Tests.Command, should.Equal, "go test")

		So(config.Session.TurnTime, should.Equal, 5*time.Minute)

		So(len(config.Session.ExerciseReference), should.Equal, 2)
		So(config.Session.ExerciseReference[0], should.Equal, "http://problems.example.com/puzzle-1")
		So(config.Session.ExerciseReference[1], should.Equal, "http://problems.example2.com/puzzle-2")
	})

	Convey("Load incomplete file", t, func() {
		ioutil.WriteFile("/tmp/jacodoma_jcdmarc.txt_test", []byte(incompleteFileContent), 0644)
		defer os.Remove("/tmp/jacodoma_jcdmarc.txt_test")

		config, err := LoadProjectConfigFile("/tmp/jacodoma_jcdmarc.txt_test")

		So(err, should.Equal, nil)
		So(config, should.NotEqual, ProjectConfig{})

		So(config.Session.TurnTime, should.Equal, 4*time.Minute)
		So(config.Session.Critical, should.Equal, 1*time.Minute)
		So(config.Tests.OnTimeout, should.Equal, 10*time.Second)
		So(config.Report.DbFile, should.Equal, "db.jcdmdb")

		So(config.Session.ShuffleUsersOrder, should.BeTrue)

		So(len(config.Session.ExerciseReference), should.Equal, 1)
		So(config.Session.ExerciseReference[0], should.Equal, "http://problems.example3.com/puzzle-3")
	})
}
