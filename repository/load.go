/*
Copyright 2019 Kohl's Department Stores, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package repository

import (
	"fmt"

	"github.com/KohlsTechnology/git2consul-go/config"
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

		auth, err := GetAuth(repoConfig)
		if err != nil {
			return nil, fmt.Errorf("Error getting AuthMethod: %s", err)
		}

		r, state, err := New(cfg.LocalStore, repoConfig, auth)
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
