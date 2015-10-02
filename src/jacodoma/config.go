package jacodoma

import (
	"gopkg.in/gcfg.v1"
)

type ProjectConfig struct {
	Session struct {
		ExerciseReference         string
		NotifyBadBehaviour        bool
		TurnTime                  string
		Critical                  string
		UseSoundNotification      bool
		SoundNotificationFilename string
		LockScreenOnTimeout       bool
		ShuffleUsersOrder         bool
	}

	Tests struct {
		Command       string
		OnEveryChange bool
		OnTimeout     string
		Files         string
	}

	Project struct {
		VC                  string
		CommitOnEveryChange bool
		SourceFiles         string
	}

	Report struct {
		DbFile string
	}

	UI struct {
		Type string
		Skin string
	}
}

func LoadProjectConfigFile(filename string) (ProjectConfig, error) {
	var config ProjectConfig
	err := gcfg.ReadFileInto(&config, filename)
	return config, err
}
