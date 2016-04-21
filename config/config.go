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

type Hooks []Hook

type Repo struct {
	Name     string   `json:"name"`
	Url      string   `json:"url"`
	Branches []string `json:"branches"`
	Hooks    *Hooks   `json:"hooks"`
}

type Repos []Repo

type Config struct {
	Repos      *Repos `json:"repos"`
	LocalStore string `json:"local_store,omitempty"`
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

	log.Debugf("Config: %+v", config)

	log.Info("Setting configuration with sane defaults")
	err = config.setDefaultConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Return a configuration with sane defaults
func (c *Config) setDefaultConfig() error {
	return nil
}
