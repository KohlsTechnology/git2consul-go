package repository

import (
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/libgit2/git2go"
)

type Repository struct {
	*git.Repository
	repoConfig *config.Repo
	store      string

	UpdateCh chan bool
}

type Repositories []*Repository

func LoadRepos(cfg *config.Config) (Repositories, error) {
	repos := []*Repository{}
	for _, repo := range cfg.Repos {
		// Create Repository object for each repo
		store := filepath.Join(cfg.LocalStore, repo.Name)

		raw_repo, err := git.OpenRepository(store)
		if err != nil {
			log.Warnf("Cannot load repository: %s", err)
		}

		r := &Repository{
			raw_repo,
			repo,
			store,
			make(chan bool, 1),
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
