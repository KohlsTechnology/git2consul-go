package config

import "time"

type Hook struct {
	Type string `json:"type"`

	// Specific to polling
	Interval time.Duration `json:"interval"`

	// Specific to webhooks
	Url string `json:"url,omitempty"`
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
	LocalStore string         `json:"local_store"`
	HookSvr    *HookSvrConfig `json:"webhook"`
	Repos      []*Repo        `json:"repos"`
	Consul     *ConsulConfig  `json:"consul"`
}

type HookSvrConfig struct {
	Address string `json:"address,omitempty"`
	Port    int    `json:"port"`
}

type ConsulConfig struct {
	Address   string `json:"address"`
	Token     string `json:"token,omitempty"`
	SSLEnable bool   `json:"ssl"`
	SSLVerify bool   `json:"ssl_verify,omitempty"`
}

func (r *Repo) String() string {
	if r != nil {
		return r.Name
	}
	return ""
}
