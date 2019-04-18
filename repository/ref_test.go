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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/KohlsTechnology/git2consul-go/config/mock"
	"github.com/KohlsTechnology/git2consul-go/repository/mocks"
	git "gopkg.in/src-d/go-git.v4"
)

func TestCheckRef(t *testing.T) {
	remoteRepo, remotePath := mocks.InitRemote(t)
	defer os.RemoveAll(remotePath)

	repoConfig := mock.RepoConfig(remotePath)
	dstPath, err := ioutil.TempDir("", repoConfig.Name)
	assert.Nil(t, err)
	defer os.RemoveAll(dstPath)

	localRepo, err := git.PlainClone(dstPath, false, &git.CloneOptions{URL: repoConfig.URL})
	assert.Nil(t, err)

	repo := &Repository{
		Repository: localRepo,
		Config:     repoConfig,
	}

	h, err := repo.Head()
	assert.Nil(t, err)

	ref := h.Name().Short()
	// Push a commit to the repository
	mocks.Add(t, remoteRepo, "tree/test.yml", []byte("foo"))
	mocks.Commit(t, remoteRepo, "Add test.yml file.")

	err = repo.Pull("master")
	assert.Nil(t, err)

	err = repo.CheckRef(ref)
	assert.Nil(t, err)

}
