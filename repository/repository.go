package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"gopkg.in/libgit2/git2go.v23"
)

type Repository struct {
	*git.Repository
	repoConfig *config.Repo
	store      string

	// Channel to notify repo clone
	cloneCh chan struct{}

	// Channel to notify repo change
	changeCh chan struct{}
	sync.Mutex
}

type Repositories []*Repository

// Populates Repository slice from configuration. It also
// handles cloning of the repository if not present
func LoadRepos(cfg *config.Config) (Repositories, error) {
	repos := []*Repository{}

	// Create Repository object for each repo
	for _, cRepo := range cfg.Repos {
		store := filepath.Join(cfg.LocalStore, cRepo.Name)

		r := &Repository{
			repoConfig: cRepo,
			store:      store,
			cloneCh:    make(chan struct{}, 1),
			changeCh:   make(chan struct{}, 1),
		}

		// Check if repository can be opened
		repo, err := git.OpenRepository(store)
		if err != nil {
			// If path does not exist or is not directory clone
			if f, err := os.Stat(store); os.IsNotExist(err) || f.IsDir() == false {
				err := os.Mkdir(r.store, 0755)
				if err != nil {
					return nil, err
				}
			}

			err = r.Clone()
			if err != nil {
				return nil, err
			}
			log.Infof("(git): Repository %s not cached, cloned to %s", r.repoConfig.Name, r.store)
		} else {
			r.Repository = repo
			// Check on Url
			if r.checkUrl(cRepo.Url) == false {
				return nil, fmt.Errorf("Repository %s exists locally on %s", r.repoConfig.Name, r.store)
			}
		}

		repos = append(repos, r)
	}

	return repos, nil
}

func (r *Repository) Name() string {
	return r.repoConfig.Name
}

func (r *Repository) Store() string {
	return r.store
}

func (r *Repository) ChangeCh() <-chan struct{} {
	return r.changeCh
}

func (r *Repository) CloneCh() <-chan struct{} {
	return r.cloneCh
}

func (r *Repository) checkUrl(url string) bool {
	rm, err := r.Remotes.Lookup("origin")
	if err != nil {
		return false
	}

	absPath, err := filepath.Abs(url)
	if err != nil {
		return false
	}
	if strings.Compare(rm.Url(), absPath) == 0 {
		return true
	}

	return false
}
