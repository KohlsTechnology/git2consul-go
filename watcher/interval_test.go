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

package watch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KohlsTechnology/git2consul-go/config/mock"
	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/KohlsTechnology/git2consul-go/repository/mocks"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetHandler(discard.New())
}

func TestPollBranches(t *testing.T) {
	remote, remotePath := mocks.InitRemote(t)
	defer os.RemoveAll(remotePath)

	cfg := mock.Config(remotePath)
	defer os.RemoveAll(cfg.LocalStore)
	repoConfig := cfg.Repos[0]

	repo, _, err := repository.New(cfg.LocalStore, repoConfig, nil)
	assert.NoError(t, err)

	mocks.Add(t, remote, "example/check_interval.txt", []byte("Example content for checke_interval"))
	mocks.Commit(t, remote, "Interval check")

	w := &Watcher{
		Repositories: []repository.Repo{repo},
		RepoChangeCh: make(chan repository.Repo, 1),
		ErrCh:        make(chan error),
		RcvDoneCh:    make(chan struct{}, 1),
		SndDoneCh:    make(chan struct{}, 1),
		logger:       log.WithField("caller", "watcher"),
		hookSvr:      nil,
		once:         true,
	}

	err = w.pollBranches(repo)
	assert.NoError(t, err)

	assert.FileExists(t, filepath.Join(repository.WorkDir(repo), "example", "check_interval.txt"))
}
