package config

import "time"

// Hook is the configuration for hooks
type Hook struct {
	Type string `json:"type"`

	// Specific to polling
	Interval time.Duration `json:"interval"`

	// Specific to webhooks
	Url string `json:"url,omitempty"`
}

// Repo is the configuration for the repository
type Repo struct {
	Name     string   `json:"name"`
	Url      string   `json:"url"`
	Branches []string `json:"branches"`
	Hooks    []*Hook  `json:"hooks"`
}

// Config is used to represent the passed in configuration
type Config struct {
	LocalStore string         `json:"local_store"`
	HookSvr    *HookSvrConfig `json:"webhook"`
	Repos      []*Repo        `json:"repos"`
	Consul     *ConsulConfig  `json:"consul"`
}

// HookSvrConfig is the configuration for the git hoooks server
type HookSvrConfig struct {
	Address string `json:"address,omitempty"`
	Port    int    `json:"port"`
}

// ConsulConfig is the configuration for the Consul client
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
