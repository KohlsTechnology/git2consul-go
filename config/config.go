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
	Hooks    Hooks    `json:"hooks"`
}

type Repos []Repo

type Config struct {
	Repos Repos `json:"repos"`
}

func Load(file string) *Config {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		os.Exit(1)
	}

	config := &Config{}

	json.Unmarshal(content, config)

	log.Debugf("Config: %+v", config)

	return config
}
