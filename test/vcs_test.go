package jacodoma

import (
	. "../src"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

func TestVcs(t *testing.T) {
	Convey("Git Reposiotry Test", t, func() {
		dirName := "/tmp/jacodoma_dir"
		defer os.RemoveAll(dirName)
		os.Mkdir(dirName, os.ModeDir|0755)

		Convey("Create repository", func() {
			_, err := CreateVcsRepository(dirName)

			So(err, should.Equal, nil)

			repoStat, err := os.Stat(dirName + "/.git")

			So(err, should.Equal, nil)

			So(repoStat.IsDir(), should.BeTrue)
		})

		Convey("Commit files on new repository", func() {
			repo, err := CreateVcsRepository(dirName)

			So(err, should.Equal, nil)

			ioutil.WriteFile(dirName+"/README", []byte(`Hi, I am a file!`), 0644)
			ioutil.WriteFile(dirName+"/TODO", []byte(`Hi again, never forget me!`), 0644)

			meta := CreateCommitMetadata("Leandro", "leandrosansilva@gmail.com")

			err = repo.CommitFiles([]string{"README", "TODO"}, meta)

			So(err, should.Equal, nil)
		})

		Convey("Change files in the repository", func() {
			repo, err := CreateVcsRepository(dirName)

			So(err, should.Equal, nil)

			{
				ioutil.WriteFile(dirName+"/README", []byte(`Hi, I am a file!`), 0644)
				ioutil.WriteFile(dirName+"/TODO", []byte(`Hi again, never forget me!`), 0644)

				meta := CreateCommitMetadata("Leandro", "leandrosansilva@gmail.com")

				err = repo.CommitFiles([]string{"README", "TODO"}, meta)

				So(err, should.Equal, nil)
			}

			{
				ioutil.WriteFile(dirName+"/README", []byte(`Hi, I am a file!\nAnd I am another line!`), 0644)
				ioutil.WriteFile(dirName+"/LICENCE", []byte(`Copycenter, no extremism`), 0644)

				meta := CreateCommitMetadata("Leandro", "leandrosansilva@gmail.com")

				err = repo.CommitFiles([]string{"README", "LICENCE"}, meta)

				So(err, should.Equal, nil)
			}
		})
	})

}
