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
	"path"
	"testing"

	"github.com/KohlsTechnology/git2consul-go/config"
	"github.com/KohlsTechnology/git2consul-go/kv/mocks"
	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/apex/log"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

//TestPutKVRef test functionality of putKVRef function.
func TestKVRef(t *testing.T) {
	var repo repository.Repo = &mocks.Repo{Config: &config.Repo{}, T: t}
	repo.Pull("master") //nolint:errcheck
	handler := &KVHandler{
		API: &mocks.KV{T: t},
		logger: log.WithFields(log.Fields{
			"caller": "consul",
		})}
	branch, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}
	refFile := fmt.Sprintf("%s.ref", branch.Name().Short())
	key := path.Join(repo.Name(), refFile)
	commit := branch.Hash().String()

	t.Run("TestPutKVRef", func(t *testing.T) {
		testPutKVRef(t, branch.Name().Short(), key, commit, handler, repo)
	})
	t.Run("TestPutKVRefModifiedIndex", func(t *testing.T) {
		testPutKVRefModifiedIndex(t, branch.Name().Short(), key, commit, handler, repo)
	})
}

func testPutKVRef(t *testing.T, branch string, key string, commit string, handler *KVHandler, repo repository.Repo) {
	err := handler.putKVRef(repo, branch)
	if err != nil {
		t.Fatal(err)
	}
	kvBranch, _, _ := handler.Get(key, nil)
	assert.Equal(t, string(kvBranch.Value), commit)

}

func testPutKVRefModifiedIndex(t *testing.T, branch string, key string, commit string, handler *KVHandler, repo repository.Repo) {
	lastCommit, err := handler.getKVRef(repo, branch)
	if err != nil {
		t.Fatal(err)
	}
	handler.API.Put(&api.KVPair{Key: key, Value: []byte(lastCommit)}, nil) //nolint:errcheck

	err = handler.putKVRef(repo, branch)
	assert.IsType(t, &TransactionIntegrityError{}, err)
	t.Log(err)
}
