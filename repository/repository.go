package repository

import (
	"path/filepath"
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

// Populates Repository slice from configuration
// Handles cloning of the repository if not present
func LoadRepos(cfg *config.Config) (Repositories, error) {
	repos := []*Repository{}

	// Create Repository object for each repo
	for _, repo := range cfg.Repos {
		store := filepath.Join(cfg.LocalStore, repo.Name)

		r := &Repository{
			repoConfig: repo,
			store:      store,
			cloneCh:    make(chan struct{}, 1),
			changeCh:   make(chan struct{}, 1),
		}

		repo, err := git.OpenRepository(store)
		if err != nil {
			log.Infof("Repository %s not cached, cloning to %s", r.repoConfig.Name, r.store)
			err = r.Clone()
			if err != nil {
				return nil, err
			}
		} else {
			r.Repository = repo
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
