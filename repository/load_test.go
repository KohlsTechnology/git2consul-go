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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/KohlsTechnology/git2consul-go/config/mock"
	"github.com/KohlsTechnology/git2consul-go/repository/mocks"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

func init() {
	log.SetHandler(discard.New())
}

//TestLoadRepos test load repos from the configuration file.
func TestLoadRepos(t *testing.T) {
	_, remotePath := mocks.InitRemote(t)
	defer os.RemoveAll(remotePath)

	cfg := mock.Config(remotePath)
	defer os.RemoveAll(cfg.LocalStore)
	_, err := LoadRepos(cfg)
	assert.Nil(t, err)
}

//TestLoadReposExistingDir tests load to existing repo.
func TestLoadReposExistingDir(t *testing.T) {
	bareDir, err := ioutil.TempDir("", "bare-dir")
	defer os.RemoveAll(bareDir)
	assert.Nil(t, err)

	cfg := mock.Config(bareDir)
	defer os.RemoveAll(cfg.LocalStore)

	_, err = LoadRepos(cfg)

	assert.NotNil(t, err)
}

//TestLoadReposInvalidRepo verifies failure in case wrong url provided.
func TestLoadReposInvalidRepo(t *testing.T) {
	cfg := mock.Config("bogus-url")
	defer os.RemoveAll(cfg.LocalStore)
	_, err := LoadRepos(cfg)
	assert.NotNil(t, err)
}

func TestLoadReposExistingRepo(t *testing.T) {
	_, remotePath := mocks.InitRemote(t)
	defer os.RemoveAll(remotePath)

	cfg := mock.Config(remotePath)
	defer os.RemoveAll(cfg.LocalStore)
	localRepoPath := filepath.Join(cfg.LocalStore, cfg.Repos[0].Name)

	// Init a repo in the local store, with same name are the "remote"
	err := os.Mkdir(localRepoPath, 0755)
	assert.Nil(t, err)

	repo, err := git.PlainInit(localRepoPath, false)
	assert.Nil(t, err)

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"/foo/bar"},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadRepos(cfg)
	assert.Nil(t, err)
}
