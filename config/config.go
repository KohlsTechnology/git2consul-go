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

package config

import "time"

// Credentials is the representation of git authentication
type Credentials struct {
	Username   string     `json:"username,omitempty"`
	Password   string     `json:"password,omitempty"`
	PrivateKey PrivateKey `json:"private_key,omitempty"`
}

// PrivateKey is the representation of private key used for the authentication
type PrivateKey struct {
	Key      string `json:"pk_key"`
	Username string `json:"pk_username,omitempty"`
	Password string `json:"pk_password,omitempty"`
}

// Hook is the configuration for hooks
type Hook struct {
	Type string `json:"type"`

	// Specific to polling
	Interval time.Duration `json:"interval"`

	// Specific to webhooks
	URL string `json:"url,omitempty"`
}

// Repo is the configuration for the repository
type Repo struct {
	Name           string      `json:"name"`
	URL            string      `json:"url"`
	Branches       []string    `json:"branches"`
	Hooks          []*Hook     `json:"hooks"`
	SourceRoot     string      `json:"source_root"`
	MountPoint     string      `json:"mount_point"`
	ExpandKeys     bool        `json:"expand_keys,omitempty"`
	SkipBranchName bool        `json:"skip_branch_name,omitempty"`
	SkipRepoName   bool        `json:"skip_repo_name,omitempty"`
	Credentials    Credentials `json:"credentials,omitempty"`
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
