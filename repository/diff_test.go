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

	"github.com/KohlsTechnology/git2consul-go/config/mock"
	"github.com/KohlsTechnology/git2consul-go/repository/mocks"
	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
)

func TestDiffStatus(t *testing.T) {
	remoteRepo, remotePath := mocks.InitRemote(t)
	defer os.RemoveAll(remotePath)

	repoConfig := mock.RepoConfig(remotePath)
	dstPath, err := ioutil.TempDir("", repoConfig.Name)
	defer os.RemoveAll(dstPath)
	assert.Nil(t, err)

	localRepo, err := git.PlainClone(dstPath, false, &git.CloneOptions{URL: repoConfig.URL})
	assert.Nil(t, err)

	repo := &Repository{
		Repository: localRepo,
		Config:     repoConfig,
	}

	h, err := repo.Head()
	assert.Nil(t, err)

	oldRef := h.Hash().String()

	// Push a commit to the repository
	mocks.Add(t, remoteRepo, "tree/test.yml", []byte("foo"))
	mocks.Commit(t, remoteRepo, "Add test.yml file.")

	err = repo.Pull("master")
	assert.Nil(t, err)

	deltas, err := repo.DiffStatus(oldRef)
	assert.Nil(t, err)

	assert.Len(t, deltas, 1)

	action, err := deltas[0].Action()
	assert.Nil(t, err)

	assert.Equal(t, action, merkletrie.Insert)

}
