package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/hashicorp/consul/api"
)

// Load maps the configuration provided from a file to a Configuration object
func Load(file string) (*Config, error) {
	// log context
	logger := log.WithFields(log.Fields{
		"caller": "config",
	})

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Create Config object pointer and unmashal JSON into it
	config := &Config{
		Consul:  &ConsulConfig{},
		HookSvr: &HookSvrConfig{},
	}
	err = json.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}

	logger.Info("Setting configuration with sane defaults")
	config.setDefaultConfig()
	config.setDefaultConsulConfig()

	err = config.checkConfig()
	if err != nil {
		return nil, err
	}

	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Using configuration: %s", jsonConfig)

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
func (c *Config) setDefaultConfig() {

	// Set the default cache store to be the OS' temp dir
	if len(c.LocalStore) == 0 {
		c.LocalStore = os.TempDir()
	}

	// Set the default webhook port
	if c.HookSvr.Port == 0 {
		c.HookSvr.Port = 9000
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
}

// This is to return default values so that the returned JSON is correctly populated
func (c *Config) setDefaultConsulConfig() {
	defConfig := api.DefaultConfig()

	if c.Consul.Address == "" {
		c.Consul.Address = defConfig.Address
	}
}
