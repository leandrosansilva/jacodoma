package jacodoma

import (
	. "../src/jacodoma"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

func TestConfigLoading(t *testing.T) {
	fileContent := `
[Session]
ExerciseReferences = http://problems.example.com/puzzle-1
NotifyBadBehaviour = true
TurnTime = 5min
Critical = 4min
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
Skin = Default
  `

	ioutil.WriteFile("/tmp/jacodoma_jcdmarc.txt_test", []byte(fileContent), 0644)
	defer os.Remove("/tmp/jacodoma_jcdmarc.txt_test")

	Convey("Load example file", t, func() {
		config, err := LoadProjectConfigFile("/tmp/jacodoma_jcdmarc.txt_test")

		So(err, should.Equal, nil)
		So(config, should.NotEqual, ProjectConfig{})

		So(config.Tests.Command, should.Equal, "go test")
	})
}
