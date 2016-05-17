package config

import "time"

type Hook struct {
	Type string `json:"type"`

	// Specific to polling
	Interval time.Duration `json:"interval"`

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

func (r *Repo) String() string {
	return r.Name
}
