package repository

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cleung2010/go-git2consul/config"
	"gopkg.in/libgit2/git2go.v24"
)

type Repository struct {
	sync.Mutex

	*git.Repository

	repoConfig *config.Repo
	basePath   string
}

const (
	RepositoryUnchanged = 0 + iota
	RepositoryCloned
	RepositoryOpened
)

func (r *Repository) Name() string {
	name := filepath.Base(r.Workdir())
	return name
}

// Create repository object, whether by cloning or opening
func (r *Repository) init(store string) (int, error) {
	fi, err := os.Stat(store)

	// Case: Directory doesn't exist
	if os.IsNotExist(err) || fi.IsDir() == false {
		// log.Infof("(git): Repository %s not cached, cloning to %s", cRepo.Name, store)
		err := os.Mkdir(store, 0755)
		if err != nil {
			return RepositoryUnchanged, err
		}

		err = r.Clone()
		if err != nil {
			return RepositoryUnchanged, err
		}

		return RepositoryCloned, nil
	} else if err != nil {
		// Some other stat error
		return RepositoryUnchanged, err
	}

	// Case: Not a git repository, remove directory and clone
	_, err = os.Stat(filepath.Join(store, ".git"))
	if os.IsNotExist(err) {
		// log.Warnf("(git): %s exists locally, overwritting", cRepo.Name)
		err := os.RemoveAll(store)
		if err != nil {
			return RepositoryUnchanged, err
		}

		err = os.Mkdir(store, 0755)
		if err != nil {
			return RepositoryUnchanged, err
		}

		err = r.Clone()
		if err != nil {
			return RepositoryUnchanged, err
		}

		return RepositoryCloned, nil
	} else if err != nil {
		// Some other stat error
		return RepositoryUnchanged, err
	}

	// Open repository otherwise
	repo, err := git.OpenRepository(store)
	if err != nil {
		return RepositoryUnchanged, err
	}

	// Check if config repo and cached repo matches
	rm, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return RepositoryUnchanged, err
	}

	absPath, err := filepath.Abs(r.repoConfig.Url)
	if err != nil {
		return RepositoryUnchanged, err
	}
	// If not equal attempt to recreate the repo
	if strings.Compare(rm.Url(), absPath) != 0 {
		// log.Warnf("(git): Diffrent %s repository exists locally, overwritting", cRepo.Name)
		os.RemoveAll(store) // Potentially dangerous?
		err := os.Mkdir(store, 0755)
		if err != nil {
			return RepositoryUnchanged, err
		}

		err = r.Clone()
		if err != nil {
			return RepositoryUnchanged, err
		}
	} else {
		r.Repository = repo

		return r, nil
	}
}
