package repository

import (
	"fmt"

	"github.com/Cimpress-MCP/go-git2consul/config"
	"github.com/apex/log"
)

// LoadRepos populates Repository slice from configuration. It also
// handles cloning of the repository if not present
func LoadRepos(cfg *config.Config) ([]*Repository, error) {
	logger := log.WithFields(log.Fields{
		"caller": "repository",
	})
	repos := []*Repository{}

	// Create Repository object for each repo
	for _, repoConfig := range cfg.Repos {
		r, state, err := New(cfg.LocalStore, repoConfig)
		if err != nil {
			return nil, fmt.Errorf("Error loading %s: %s", repoConfig.Name, err)
		}

		switch state {
		case RepositoryCloned:
			logger.Infof("Cloned repository %s", r.Name())
		case RepositoryOpened:
			logger.Infof("Loaded repository %s", r.Name())
		}

		repos = append(repos, r)
	}

	if len(repos) == 0 {
		return repos, fmt.Errorf("No repositories provided in the configuration")
	}

	return repos, nil
}
