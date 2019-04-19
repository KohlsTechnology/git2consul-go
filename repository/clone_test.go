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
)

func TestClone(t *testing.T) {
	_, remotePath := mocks.InitRemote(t)
	defer os.RemoveAll(remotePath)

	cfg := mock.Config(remotePath)
	defer os.RemoveAll(cfg.LocalStore)

	repo := &Repository{
		Config: cfg.Repos[0],
	}

	localPath, err := ioutil.TempDir(cfg.LocalStore, repo.Config.Name)
	assert.Nil(t, err)
	defer os.RemoveAll(localPath)

	err = repo.Clone(localPath)
	assert.Nil(t, err)
}
