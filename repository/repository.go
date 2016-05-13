package repository

import (
	"sync"

	"github.com/cleung2010/go-git2consul/config"
	"gopkg.in/libgit2/git2go.v24"
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
	rs, err := loadRepos(cfg)
	if err != nil {
		return nil, err
	}

	return rs, nil
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
