package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Cimpress-MCP/go-git2consul/config"
	"gopkg.in/libgit2/git2go.v24"
)

type Repository struct {
	sync.Mutex

	*git.Repository
	Config *config.Repo
}

const (
	RepositoryError = iota
	RepositoryCloned
	RepositoryOpened
)

// Returns the repository name
func (r *Repository) Name() string {
	return r.Config.Name
}

// Returns the branch name
func (r *Repository) Branch() string {
	head, err := r.Head()
	if err != nil {
		return ""
	}
	bn, err := head.Branch().Name()
	if err != nil {
		return ""
	}

	return bn
}

func New(repoBasePath string, repoConfig *config.Repo) (*Repository, int, error) {
	repoPath := filepath.Join(repoBasePath, repoConfig.Name)

	r := &Repository{
		Repository: &git.Repository{},
		Config:     repoConfig,
	}

	state, err := r.init(repoPath)
	if err != nil {
		return nil, RepositoryError, err
	}

	if r.Repository == nil {
		return nil, RepositoryError, fmt.Errorf("Could not find git repostory")
	}

	return r, state, nil
}

// Attempt to create *git.Repository object, whether by cloning or opening.
// This method also returns the state of the repository creation for loggging
// TODO: Refactor this
func (r *Repository) init(repoPath string) (int, error) {
	fi, err := os.Stat(repoPath)

	// Case: Directory doesn't exist
	if os.IsNotExist(err) || fi.IsDir() == false {
		// log.Printf("(git): Repository %s not cached, cloning to %s", r.Name(), r.store)
		err := os.Mkdir(repoPath, 0755)
		if err != nil {
			return RepositoryError, err
		}

		err = r.Clone(repoPath)
		if err != nil {
			return RepositoryError, err
		}

		return RepositoryCloned, nil
	} else if err != nil {
		// Some other stat error
		return RepositoryError, err
	}

	// Case: Not a git repository, remove directory and clone
	_, err = os.Stat(filepath.Join(repoPath, ".git"))
	if os.IsNotExist(err) {
		// log.Printf("(git): %s exists locally, overwritting", r.Name())
		err := os.RemoveAll(repoPath)
		if err != nil {
			return RepositoryError, err
		}

		err = os.Mkdir(repoPath, 0755)
		if err != nil {
			return RepositoryError, err
		}

		err = r.Clone(repoPath)
		if err != nil {
			return RepositoryError, err
		}

		return RepositoryCloned, nil
	} else if err != nil {
		// Some other stat error
		return RepositoryError, err
	}

	// Open repository otherwise
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return RepositoryError, err
	}

	// Check if config repo and cached repo matches
	rm, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return RepositoryError, err
	}

	absPath, err := filepath.Abs(r.Config.Url)
	if err != nil {
		return RepositoryError, err
	}
	// If not equal attempt to recreate the repo
	if strings.Compare(rm.Url(), absPath) != 0 {
		os.RemoveAll(repoPath) // Potentially dangerous?
		err := os.Mkdir(repoPath, 0755)
		if err != nil {
			return RepositoryError, err
		}

		err = r.Clone(repoPath)
		if err != nil {
			return RepositoryError, err
		}
	} else {
		r.Repository = repo
	}

	return RepositoryOpened, nil
}
