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

package mock

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/KohlsTechnology/git2consul-go/config"
)

// RepoConfig returns a mock Repo config object
func RepoConfig(repoURL string) *config.Repo {
	return &config.Repo{
		Name:     "git2consul-test-local",
		URL:      repoURL,
		Branches: []string{"master"},
		Hooks: []*config.Hook{
			{
				Type:     "polling",
				Interval: 5 * time.Second,
			},
		},
	}
}

// Config returns a mock Config object with one repository configuration
func Config(repoURL string) *config.Config {
	localStore, err := ioutil.TempDir("", "git2consul-test-local")
	if err != nil {
		localStore = os.TempDir()
	}
	return &config.Config{
		LocalStore: localStore,
		HookSvr: &config.HookSvrConfig{
			Port: 9000,
		},
		Repos: []*config.Repo{
			{
				Name:     "git2consul-test-local",
				URL:      repoURL,
				Branches: []string{"master"},
				Hooks: []*config.Hook{
					{
						Type:     "polling",
						Interval: 5 * time.Second,
					},
				},
			},
		},
		Consul: &config.ConsulConfig{
			Address: "127.0.0.1:8500",
		},
	}
}
