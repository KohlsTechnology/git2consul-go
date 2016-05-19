package repository

import (
	"fmt"
	"path/filepath"

	"github.com/apex/log"
	"github.com/cleung2010/go-git2consul/config"
)

// Populates Repository slice from configuration. It also
// handles cloning of the repository if not present
func LoadRepos(cfg *config.Config) ([]*Repository, error) {
	repos := []*Repository{}

	// Create Repository object for each repo
	for _, cRepo := range cfg.Repos {
		store := filepath.Join(cfg.LocalStore, cRepo.Name)

		r := &Repository{
			repoConfig: cRepo,
			basePath:   cfg.LocalStore,
		}
		repoStatus, err := r.init(store)
		if err != nil {
			return nil, err
		}
		switch repoStatus {
		case RepositoryCloned:
			log.Info("Cloned")
		case RepositoryOpened:
			log.Info("Opened")
		}

		repos = append(repos, r)
	}

	if len(repos) == 0 {
		return repos, fmt.Errorf("No repositories provided in the configuration")
	}

	return repos, nil
}
