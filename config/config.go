package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

type Hook struct {
	Type string `json:"type"`

	// Specific to polling
	Interval int `json:"interval"`

	// Specific to webhooks
	Url  string `json:"url,omitempty"`
	Port int    `json:"port,omitempty"`
}

type Hooks []*Hook

type Repo struct {
	Name     string   `json:"name"`
	Url      string   `json:"url"`
	Branches []string `json:"branches"`
	Hooks    []*Hook  `json:"hooks"`
}

type Repos []*Repo

type Config struct {
	Repos      []*Repo `json:"repos"`
	LocalStore string  `json:"local_store,omitempty"`
}

// Create configuration from a provided file path
func Load(file string) (*Config, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		os.Exit(1)
	}

	// Create Config object pointer and unmashal JSON into it
	config := &Config{}
	err = json.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}

	log.Info("Setting configuration with sane defaults")
	err = config.setDefaultConfig()
	if err != nil {
		return nil, err
	}

	log.Debugf("Using configuration: %+v", config)

	return config, nil
}

// Return a configuration with sane defaults
func (c *Config) setDefaultConfig() error {

	// Set the default cache store to be the OS' temp dir
	if len(c.LocalStore) == 0 {
		c.LocalStore = os.TempDir()
	}

	//For each repo, set default branch and hook
	for _, repo := range c.Repos {
		branch := []string{"master"}
		// If there are no branches, set it to master
		if len(repo.Branches) == 0 {
			repo.Branches = branch
		}

		// If there are no hooks, set a 60s polling hook
		if len(repo.Hooks) == 0 {
			hook := &Hook{
				Type:     "polling",
				Interval: 60,
			}

			repo.Hooks = append(repo.Hooks, hook)
		}
	}
	return nil
}
