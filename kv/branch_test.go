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
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/KohlsTechnology/git2consul-go/config"
	"github.com/KohlsTechnology/git2consul-go/kv/mocks"
	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/apex/log"
	"github.com/stretchr/testify/assert"
)

//TestPutBranch verifies putBranch function.
func TestPutBranch(t *testing.T) {
	var repo repository.Repo
	_, path, _, _ := runtime.Caller(0)
	repo = &mocks.Repo{Path: filepath.Dir(path), Config: &config.Repo{}}
	handler := &KVHandler{
		API: &mocks.KV{T: t},
		logger: log.WithFields(log.Fields{
			"caller": "consul",
		}),
	}

	handler.putBranch(repo, repo.Branch()) //nolint:errcheck
	handler.Commit()                       //nolint:errcheck

	err := filepath.Walk(repository.WorkDir(repo), func(path string, f os.FileInfo, err error) error { //nolint:staticcheck
		// Skip the .git directory
		if f.IsDir() && f.Name() == ".git" {
			return filepath.SkipDir
		}

		// Do not push directories
		if f.IsDir() {
			return nil
		}

		key := strings.TrimPrefix(path, repository.WorkDir(repo))
		kvPath := filepath.Join(repo.Name(), repo.Branch().Short(), key)
		kvContent, _, err := handler.Get(kvPath, nil) //nolint:ineffassign,staticcheck
		fileContent, err := ioutil.ReadFile(path)     //nolint:ineffassign,staticcheck

		assert.Equal(t, fileContent, kvContent.Value)
		return nil
	})
	assert.NoError(t, err)
}
