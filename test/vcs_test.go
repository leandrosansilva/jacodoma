package jacodoma

import (
	. "../src"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestVcs(t *testing.T) {
	dirName := "/tmp/jacodoma_dir"

	// create project in any case
	os.Mkdir(dirName, os.ModeDir|0755)

	Convey("Create repository", t, func() {
		_, err := CreateVcsRepository(dirName)

		So(err, should.Equal, nil)

		repoStat, err := os.Stat(dirName + "/.git")

		So(err, should.Equal, nil)

		So(repoStat.IsDir(), should.BeTrue)
	})

	Convey("Commit file", t, func() {
		repo, err := CreateVcsRepository(dirName)

		So(err, should.Equal, nil)

		fileContent := `Hi, I am a file!`
		filename := dirName + "/README"

		ioutil.WriteFile(filename, []byte(fileContent), 0644)

	})
}
