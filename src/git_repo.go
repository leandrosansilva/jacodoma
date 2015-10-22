package jacodoma

import (
	"errors"
	git "github.com/libgit2/git2go"
	"os"
)

type Repository struct {
	Repo *git.Repository
}

func CreateVcsRepository(dirName string) (Repository, error) {
	// the directory must already exist!
	repoStat, err := os.Stat(dirName)

	if err != nil {
		return Repository{nil}, err
	}

	if !repoStat.IsDir() {
		return Repository{nil}, errors.New(dirName + " must be a directory!")
	}

	gitRepo, err := git.InitRepository(dirName, false)

	if err != nil {
		return Repository{nil}, err
	}

	return Repository{gitRepo}, nil
}
