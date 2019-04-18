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

package repository

import (
	"github.com/KohlsTechnology/git2consul-go/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// GetAuth returns AuthMethod based on the passed flags
func GetAuth(repo *config.Repo) (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	var err error
	auth = nil

	if len(repo.Credentials.Password) > 0 && len(repo.Credentials.Username) > 0 {
		auth = &http.BasicAuth{
			Username: repo.Credentials.Username,
			Password: repo.Credentials.Password,
		}
	} else if len(repo.Credentials.PrivateKey.Key) > 0 {
		if len(repo.Credentials.Username) == 0 {
			repo.Credentials.Username = "git"
		}
		auth, err = ssh.NewPublicKeysFromFile(repo.Credentials.PrivateKey.Username, repo.Credentials.PrivateKey.Key, repo.Credentials.PrivateKey.Password)
		if err != nil {
			return nil, err
		}
	}

	return auth, err
}
