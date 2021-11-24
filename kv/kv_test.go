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

package kv

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/KohlsTechnology/git2consul-go/config"
	"github.com/KohlsTechnology/git2consul-go/kv/mocks"
	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/apex/log"
	"github.com/stretchr/testify/assert"
)

//TestKV runs a test against KVPUT and DeleteKV handler functions.
func TestKV(t *testing.T) {
	handler := &KVHandler{
		API: &mocks.KV{T: t},
		logger: log.WithFields(log.Fields{
			"caller": "consul",
		})}
	repoPath, err := ioutil.TempDir("", "local-repo")
	defer os.RemoveAll(repoPath)
	assert.NoError(t, err)
	repo := &mocks.Repo{Path: repoPath, Config: &config.Repo{}, T: t}
	repoPath = repository.WorkDir(repo) //nolint:ineffassign,staticcheck
	t.Run("testPutKV", func(t *testing.T) { testPutKV(t, repo, handler) })
	t.Run("testDeleteKV", func(t *testing.T) { testDeleteKV(t, repo, handler) })
}

//testPutKV verifies the data pushed by putKV function.
func testPutKV(t *testing.T, repo repository.Repo, handler *KVHandler) {
	f, err := ioutil.TempFile(repository.WorkDir(repo), "example.txt")
	f.Write([]byte("Example content")) //nolint:errcheck
	f.Close()
	assert.NoError(t, err)
	value, err := ioutil.ReadFile(f.Name())

	assert.NoError(t, err)
	prefix := strings.TrimPrefix(f.Name(), repository.WorkDir(repo))
	err = handler.PutKV(repo, prefix, value)
	if err != nil {
		t.Fatal(err)
	}
	err = handler.Commit()
	assert.NoError(t, err)

	head, err := repo.Head()
	assert.NoError(t, err)

	pair, _, err := handler.Get(fmt.Sprintf("%s/%s%s", repo.Name(), head.Name().Short(), prefix), nil)
	assert.NoError(t, err)

	if assert.NotNil(t, pair) {
		assert.Equal(t, value, pair.Value)
	}
}

//testDeleteKV ensures data has been deleted.
func testDeleteKV(t *testing.T, repo repository.Repo, handler *KVHandler) {
	f, err := ioutil.TempFile(repository.WorkDir(repo), "example.txt")
	f.Write([]byte("Example content to delete")) //nolint:errcheck
	f.Close()
	assert.NoError(t, err)
	value, err := ioutil.ReadFile(f.Name()) //nolint:ineffassign,staticcheck

	prefix := strings.TrimPrefix(f.Name(), repository.WorkDir(repo))
	err = handler.PutKV(repo, prefix, value)
	assert.NoError(t, err)

	err = handler.Commit()
	assert.NoError(t, err)

	head, err := repo.Head()
	assert.NoError(t, err)

	pair, _, err := handler.Get(fmt.Sprintf("%s/%s%s", repo.Name(), head.Name().Short(), prefix), nil)
	assert.NoError(t, err)
	assert.NotNil(t, pair)

	handler.DeleteKV(repo, prefix) //nolint:errcheck
	err = handler.Commit()
	assert.NoError(t, err)

	pair, _, err = handler.Get(fmt.Sprintf("%s/%s%s", repo.Name(), head.Name().Short(), prefix), nil)
	assert.NoError(t, err)

	assert.Nil(t, pair)
}
