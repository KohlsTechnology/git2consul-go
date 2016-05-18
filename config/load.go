package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Create configuration from a provided file path
func Load(file string) (*Config, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Create Config object pointer and unmashal JSON into it
	config := &Config{}
	err = json.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}

	log.Info("(config): Setting configuration with sane defaults")
	err = config.setDefaultConfig()
	if err != nil {
		return nil, err
	}

	err = config.checkConfig()
	if err != nil {
		return nil, err
	}

	log.Debugf("(config): Using configuration: %+v", config)
	return config, nil
}

// Check for the validitiy of the configuration file
func (c *Config) checkConfig() error {
	for _, repo := range c.Repos {
		// Check on name
		if len(repo.Name) == 0 {
			return fmt.Errorf("Repository array object missing \"name\" value")
		}

		// Check on Url
		if len(repo.Url) == 0 {
			return fmt.Errorf("%s does no have a repository URL", repo.Name)
		}

		// Check on hooks
		for _, hook := range repo.Hooks {
			if hook.Type != "polling" && hook.Type != "webhook" {
				return fmt.Errorf("Invalid hook type: %s", hook.Type)
			}

			if hook.Type == "polling" && hook.Interval <= 0 {
				return fmt.Errorf("Invalid interval: %s. Hook interval must be greater than zero", hook.Interval)
			}
		}
	}

	return nil
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
				Interval: 60 * time.Second,
			}

			repo.Hooks = append(repo.Hooks, hook)
		}
	}
	return nil
}
