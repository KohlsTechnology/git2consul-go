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
	"path"
	"sync"
	"time"

	"github.com/KohlsTechnology/git2consul-go/repository"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Watch the repo by interval. This is called as a go routine since
// ticker blocks
func (w *Watcher) pollByInterval(repo repository.Repo, wg *sync.WaitGroup) {
	defer wg.Done()
	config := repo.GetConfig()

	hooks := config.Hooks
	interval := time.Second

	// Find polling hook
	for _, h := range hooks {
		if h.Type == "polling" {
			interval = h.Interval
			break
		}
	}

	// If no polling found, don't poll
	if interval == 0 {
		return
	}

	ticker := time.NewTicker(interval * time.Second)
	defer ticker.Stop()

	// Polling error should not stop polling by interval
	for {
		err := w.pollBranches(repo)
		if err != nil {
			w.ErrCh <- err
		}

		if w.once {
			return
		}

		select {
		case <-ticker.C:
		case <-w.RcvDoneCh:
			return
		}
	}
}

func (w *Watcher) pollBranches(repo repository.Repo) error {
	storer := repo.GetStorer()
	config := repo.GetConfig()
	itr, err := repository.LocalBranches(storer)
	if err != nil {
		return err
	}
	changed := false

	var checkoutBranchFn = func(b *plumbing.Reference) error {
		branchOnRemote := repository.StringInSlice(path.Base(b.Name().String()), config.Branches)
		if branchOnRemote {
			branchName := b.Name().Short()
			err := repo.Pull(branchName)
			if err == git.NoErrAlreadyUpToDate {
				w.logger.Debugf("Up to date: %s/%s", repo.Name(), branchName)
			} else if err != nil {
				w.logger.Debugf("Unable to pull \"%s\" branch because of \"%s\"", branchName, err)
			} else {
				w.logger.Infof("Changed: %s/%s", repo.Name(), branchName)
				changed = true
			}
		}
		return nil
	}

	err = itr.ForEach(checkoutBranchFn)
	if err != nil {
		return err
	}

	if changed {
		w.RepoChangeCh <- repo
	}

	return nil
}
