package jacodoma

import (
	"bytes"
	"gopkg.in/gcfg.v1"
	"gopkg.in/mcuadros/go-defaults.v1"
	"time"
)

type ConfigDuration time.Duration

func (this *ConfigDuration) UnmarshalText(text []byte) error {
	var b bytes.Buffer

	b.Write(text)

	if len(b.String()) == 0 {
		return nil
	}

	duration, err := time.ParseDuration(b.String())

	if err == nil {
		*this = ConfigDuration(duration)
	}

	return err
}

type ProjectConfig struct {
	Session struct {
		ExerciseReference         []string
		NotifyBadBehaviour        bool           `default:false`
		TurnTime                  ConfigDuration `default:"4m"`
		Critical                  ConfigDuration `default:"1m"`
		UseSoundNotification      bool           `default:true`
		UseSystemNotification     bool           `default:true`
		SoundNotificationFilename string         `default:"sound.ogg"`
		LockScreenOnTimeout       bool           `default:true`
		ShuffleUsersOrder         bool           `default:true`
	}

	Tests struct {
		Command       string
		OnEveryChange bool           `default:"true"`
		OnTimeout     ConfigDuration `default:"10s"`
		Files         string
	}

	Project struct {
		VC                  string `default:"Git"`
		CommitOnEveryChange bool   `default:true`
		SourceFiles         string
	}

	Report struct {
		DbFile string `default:"db.jcdmdb"`
	}

	UI struct {
		Type string `default:"QML"`
		Skin string `default:"default"`
	}
}

func LoadProjectConfigFile(filename string) (ProjectConfig, error) {
	var config ProjectConfig
	defaults.SetDefaults(&config)
	err := gcfg.ReadFileInto(&config, filename)
	return config, err
}
