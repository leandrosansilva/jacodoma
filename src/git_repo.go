package jacodoma

import (
	"errors"
	git "github.com/libgit2/git2go"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Repository struct {
	Repo *git.Repository
}

type CommitMetadata struct {
	Name  string
	Email string
	Time  time.Time
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

func resolveGlobs(globs []string, basePath string) ([]string, error) {
	filenames := make([]string, 0)

	for _, glob := range globs {
		p := path.Join(basePath, glob)

		paths, err := filepath.Glob(p)

		if err != nil {
			return []string{}, err
		}

		for _, filename := range paths {
			filenames = append(filenames, filepath.Base(filename))
		}
	}

	return filenames, nil
}

func (this *Repository) CommitFiles(globs []string, meta CommitMetadata) error {
	// TODO: refactor this method, which is doing too many things!

	message := "empty message"

	index, err := this.Repo.Index()

	if err != nil {
		return err
	}

	filenames, err := resolveGlobs(globs, filepath.Dir(filepath.Clean(this.Repo.Path())))

	if err != nil {
		return err
	}

	if len(filenames) == 0 {
		return errors.New("No files do add to the repository!")
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

	tree, err := this.Repo.LookupTree(treeId)

	sig := &git.Signature{
		Name:  meta.Name,
		Email: meta.Email,
		When:  meta.Time,
	}

	head, _ := this.Repo.Head()

	currentTip := (*git.Commit)(nil)

	if head != nil {
		currentTip, err = this.Repo.LookupCommit(head.Target())
	}

	if err != nil {
		return err
	}

	if currentTip != nil {
		_, err := this.Repo.CreateCommit("HEAD", sig, sig, message, tree, currentTip)
		return err
	}

	_, err = this.Repo.CreateCommit("HEAD", sig, sig, message, tree)
	return err
}

func CreateCommitMetadata(name, email string, t time.Time) CommitMetadata {
	return CommitMetadata{name, email, t}
}
