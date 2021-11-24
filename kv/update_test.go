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
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KohlsTechnology/git2consul-go/config"
	"github.com/KohlsTechnology/git2consul-go/kv/mocks"
	"github.com/apex/log"
	"github.com/stretchr/testify/assert"
)

func TestUpdateToHead(t *testing.T) {
	handler := &KVHandler{
		API: &mocks.KV{T: t},
		logger: log.WithFields(log.Fields{
			"caller": "consul",
		},
		),
	}
	log.SetLevel(0)
	//_, repoPath, _, _ := runtime.Caller(1)
	repoPath, err := ioutil.TempDir("", "local-repo")
	defer os.RemoveAll(repoPath)
	assert.NoError(t, err)
	repo := &mocks.Repo{Path: repoPath, Config: &config.Repo{}, T: t}
	repo.Pull("master") //nolint:errcheck
	branch, err := repo.Head()
	assert.NoError(t, err)
	initialCommit := branch.Hash().String()
	repo.Pull(branch.Name().Short()) //nolint:errcheck
	//Make an initial load to the Consul KV store.
	handler.putBranch(repo, branch.Name())        //nolint:errcheck
	handler.putKVRef(repo, branch.Name().Short()) //nolint:errcheck
	//Fake commit
	f, err := ioutil.TempFile(repoPath, "example.txt")
	assert.NoError(t, err)
	f.Write([]byte("A content!")) //nolint:errcheck
	f.Close()
	fileName := strings.TrimPrefix(f.Name(), repoPath)
	repo.Add(fileName)
	//Pull the change.
	repo.Pull(branch.Name().Short()) //nolint:errcheck

	err = handler.UpdateToHead(repo)
	assert.NoError(t, err)
	branch, err = repo.Head()
	assert.NoError(t, err)

	//Ensure ref has been updated
	refFile := fmt.Sprintf("%s.ref", branch.Name().Short())
	key := path.Join(repo.Name(), refFile)

	kvBranch, _, err := handler.Get(key, nil) //nolint:ineffassign,staticcheck

	assert.Equal(t, string(kvBranch.Value), branch.Hash().String())

	assert.NotEqual(t, string(kvBranch.Value), initialCommit)

	kvPath := filepath.Join(repo.Name(), branch.Name().Short(), fileName)
	kvContent, _, err := handler.Get(kvPath, nil)                          //nolint:ineffassign,staticcheck
	fileContent, err := ioutil.ReadFile(filepath.Join(repoPath, fileName)) //nolint:ineffassign,staticcheck

	assert.Equal(t, kvContent.Value, fileContent)

	// 	return nil
	// })
}
