package repository

import (
	"os"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/libgit2/git2go"
)

type Repository struct {
	*git.Repository
	repoConfig *config.Repo
	store      string

	// Channel to notify repo change
	changeCh chan struct{}
	sync.Mutex
}

type Repositories []*Repository

func LoadRepos(cfg *config.Config) (Repositories, error) {
	repos := []*Repository{}

	// Create Repository object for each repo
	for _, repo := range cfg.Repos {
		store := filepath.Join(cfg.LocalStore, repo.Name)

		r := &Repository{
			repoConfig: repo,
			store:      store,
			changeCh:   make(chan struct{}, 1),
		}

		repo, err := git.OpenRepository(store)
		if err != nil {
			log.Infof("Repository %s not cached, cloning to %s", r.repoConfig.Name, r.store)
			err := os.Mkdir(r.store, 0755)
			if err != nil {
				return nil, err
			}
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

func (r *Repository) Branches() []string {
	return r.repoConfig.Branches
}

func (r *Repository) ChangeLock() <-chan struct{} {
	return r.changeCh
}
