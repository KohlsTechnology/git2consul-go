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
	"os"
	"path/filepath"
	"testing"

	"github.com/KohlsTechnology/git2consul-go/config/mock"
	"github.com/KohlsTechnology/git2consul-go/repository/mocks"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func TestGetAuthWithPlainAuth(t *testing.T) {
	_, remotePath := mocks.InitRemote(t)
	defer os.RemoveAll(remotePath)

	cfg := mock.Config(remotePath)
	defer os.RemoveAll(cfg.LocalStore)

	repoConfig := cfg.Repos[0]
	repoConfig.Credentials.Username = "foo"
	repoConfig.Credentials.Password = "bar"

	expectedAuth := &http.BasicAuth{
		Username: "foo",
		Password: "bar",
	}

	auth, err := GetAuth(repoConfig)
	assert.NoError(t, err)
	assert.Equal(t, expectedAuth, auth)
}

func TestGetAuthWithKeyAuth(t *testing.T) {
	wd, _ := os.Getwd()
	key := filepath.Join(wd, "../config/test-fixtures", "test_ssh_key")
	_, remotePath := mocks.InitRemote(t)
	defer os.RemoveAll(remotePath)
	cfg := mock.Config(remotePath)
	defer os.RemoveAll(cfg.LocalStore)
	repoConfig := cfg.Repos[0]
	repoConfig.Credentials.PrivateKey.Key = key
	repoConfig.Credentials.PrivateKey.Username = "foo"

	expectedAuth := &ssh.PublicKeys{
		User:   "foo",
		Signer: nil,
	}

	auth, err := GetAuth(repoConfig)
	assert.NoError(t, err)

	assert.Equal(t, expectedAuth.String(), auth.String())
}

func TestGetAuthWithoutCred(t *testing.T) {
	_, remotePath := mocks.InitRemote(t)
	defer os.RemoveAll(remotePath)
	cfg := mock.Config(remotePath)
	defer os.RemoveAll(cfg.LocalStore)
	repoConfig := cfg.Repos[0]

	auth, err := GetAuth(repoConfig)
	assert.NoError(t, err)
	assert.Nil(t, auth)
}
