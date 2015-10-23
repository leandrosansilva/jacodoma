package jacodoma

import (
	"errors"
	git "github.com/libgit2/git2go"
	"os"
	"time"
)

type Repository struct {
	Repo *git.Repository
}

type CommitMetadata struct {
	Name  string
	Email string
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

func (this *Repository) CommitFiles(filenames []string, meta CommitMetadata) error {
	message := "empty message"

	index, err := this.Repo.Index()

	if err != nil {
		return err
	}

	for _, filename := range filenames {
		err := index.AddByPath(filename)
		if err != nil {
			return err
		}
	}

	treeId, err := index.WriteTree()

	if err != nil {
		return err
	}

  head, err := tis.Repo.Head()

	if err != nil {
		return err
	}

  treeToCommit, err := func(), git.Reference*, error {
    if head == nil {
      return this.Repo.LookupTree(treeId)
    }

    currentTip, err := 
  }()

	

	if err != nil {
		return err
	}

	sig := &git.Signature{
		Name:  meta.Name,
		Email: meta.Email,
		When:  time.Unix(0, 0),
	}

	_, err = this.Repo.CreateCommit("HEAD", sig, sig, message, tree)

	if err != nil {
		return err
	}

	return nil
}

func CreateCommitMetadata(name, email string) CommitMetadata {
	return CommitMetadata{name, email}
}
