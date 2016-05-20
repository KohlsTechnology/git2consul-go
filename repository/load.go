package repository

import (
	"fmt"
	"path/filepath"

	"github.com/cleung2010/go-git2consul/config"
)

// Populates Repository slice from configuration. It also
// handles cloning of the repository if not present
func LoadRepos(cfg *config.Config) ([]*Repository, error) {
	repos := []*Repository{}

	// Create Repository object for each repo
	for _, repoConfig := range cfg.Repos {
		repoPath := filepath.Join(cfg.LocalStore, repoConfig.Name)

		r, state, err := New(repoPath, repoConfig)
		if err != nil {
			return repos, err
		}

		switch state {
		case RepositoryCloned:
			//log.Infof("Cloned repository %s", r.Name())
		case RepositoryOpened:
			//log.Infof("Loaded repository %s", r.Name())
		}

		repos = append(repos, r)
	}

	if len(repos) == 0 {
		return repos, fmt.Errorf("No repositories provided in the configuration")
	}

	return repos, nil
}
