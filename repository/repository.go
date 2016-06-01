package repository

import (
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
	Hooks []*config.Hook

	repoConfig *config.Repo
	store      string
}

const (
	RepositoryError = iota
	RepositoryCloned
	RepositoryOpened
)

// Returns the repository name
func (r *Repository) Name() string {
	return filepath.Base(r.Workdir())
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

func New(repoPath string, repoConfig *config.Repo) (*Repository, int, error) {
	r := &Repository{
		Hooks:      repoConfig.Hooks,
		repoConfig: repoConfig,
		store:      repoPath,
	}

	state, err := r.init()
	if err != nil {
		return nil, RepositoryError, err
	}

	return r, state, nil
}

// Attempt to create *git.Repository object, whether by cloning or opening.
// This method also returns the state of the repository creation for loggging
func (r *Repository) init() (int, error) {
	fi, err := os.Stat(r.store)

	// Case: Directory doesn't exist
	if os.IsNotExist(err) || fi.IsDir() == false {
		// log.Infof("(git): Repository %s not cached, cloning to %s", cRepo.Name, store)
		err := os.Mkdir(r.store, 0755)
		if err != nil {
			return RepositoryError, err
		}

		err = r.Clone()
		if err != nil {
			return RepositoryError, err
		}

		return RepositoryCloned, nil
	} else if err != nil {
		// Some other stat error
		return RepositoryError, err
	}

	// Case: Not a git repository, remove directory and clone
	_, err = os.Stat(filepath.Join(r.store, ".git"))
	if os.IsNotExist(err) {
		// log.Warnf("(git): %s exists locally, overwritting", cRepo.Name)
		err := os.RemoveAll(r.store)
		if err != nil {
			return RepositoryError, err
		}

		err = os.Mkdir(r.store, 0755)
		if err != nil {
			return RepositoryError, err
		}

		err = r.Clone()
		if err != nil {
			return RepositoryError, err
		}

		return RepositoryCloned, nil
	} else if err != nil {
		// Some other stat error
		return RepositoryError, err
	}

	// Open repository otherwise
	repo, err := git.OpenRepository(r.store)
	if err != nil {
		return RepositoryError, err
	}

	// Check if config repo and cached repo matches
	rm, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return RepositoryError, err
	}

	absPath, err := filepath.Abs(r.repoConfig.Url)
	if err != nil {
		return RepositoryError, err
	}
	// If not equal attempt to recreate the repo
	if strings.Compare(rm.Url(), absPath) != 0 {
		// log.Warnf("(git): Diffrent %s repository exists locally, overwritting", cRepo.Name)
		os.RemoveAll(r.store) // Potentially dangerous?
		err := os.Mkdir(r.store, 0755)
		if err != nil {
			return RepositoryError, err
		}

		err = r.Clone()
		if err != nil {
			return RepositoryError, err
		}
	} else {
		r.Repository = repo
	}

	return RepositoryOpened, nil
}
