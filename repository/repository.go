package repository

import (
	"sync"

	"github.com/cleung2010/go-git2consul/config"
	"gopkg.in/libgit2/git2go.v24"
)

type Repository struct {
	sync.Mutex

	*git.Repository
	repoConfig *config.Repo
	store      string

	// Channel to notify repo change
	changeCh chan struct{}
}

type Repositories []*Repository

func (r *Repository) Name() string {
	return r.repoConfig.Name
}

func (r *Repository) Store() string {
	return r.store
}

func (r *Repository) ChangeCh() <-chan struct{} {
	return r.changeCh
}
