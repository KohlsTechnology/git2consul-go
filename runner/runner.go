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

package runner

import (
	"fmt"
	"time"

	"github.com/KohlsTechnology/git2consul-go/config"
	"github.com/KohlsTechnology/git2consul-go/kv"
	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/KohlsTechnology/git2consul-go/watcher"
	"github.com/apex/log"
)

// Runner is used to initialize a watcher and kvHandler
type Runner struct {
	logger *log.Entry
	ErrCh  chan error

	// Channel that receives done signal
	RcvDoneCh chan struct{}

	// Channel that sends done signal
	SndDoneCh chan struct{}

	once bool

	kvHandler kv.Handler

	watcher *watch.Watcher
}

// NewRunner creates a new runner instance
func NewRunner(config *config.Config, once bool) (*Runner, error) {
	// var repos repository.Repo
	logger := log.WithField("caller", "runner")

	// Create repos from configuration
	repos, err := repository.LoadRepos(config)
	if err != nil {
		return nil, fmt.Errorf("Cannot load repositories from configuration: %s", err)
	}
	var reposI = make([]repository.Repo, len(repos))
	for index, repo := range repos {
		reposI[index] = repo
	}
	// Create watcher to watch for repo changes
	watcher := watch.New(reposI, config.HookSvr, once)

	// Create the handler
	handler, err := kv.New(config.Consul)
	if err != nil {
		return nil, err
	}

	runner := &Runner{
		logger:    logger,
		ErrCh:     make(chan error),
		RcvDoneCh: make(chan struct{}, 1),
		SndDoneCh: make(chan struct{}, 1),
		once:      once,
		kvHandler: handler,
		watcher:   watcher,
	}

	return runner, nil
}

// Start the runner
func (r *Runner) Start() {
	defer close(r.SndDoneCh)

	go r.watcher.Watch()

	for {
		select {
		case repo := <-r.watcher.RepoChangeCh:
			// Handle change, and return if error on handler
			retry := 0
			var err error
			for ok := true; ok && retry < 3; retry++ {
				err = r.kvHandler.HandleUpdate(repo)
				_, ok = err.(*kv.TransactionIntegrityError)
				time.Sleep(1000 * time.Millisecond)
			}
			if err != nil {
				r.ErrCh <- err
			}
		case <-r.watcher.SndDoneCh: // This triggers when watcher gets an error that causes termination
			r.logger.Info("Watcher received finish")
			return
		case <-r.RcvDoneCh:
			r.logger.Info("Received finish")
			return
		}
	}
}

// Stop the runner, cleaning up any routines that it's running. In this case, it will stop
// the watcher before closing DoneCh
func (r *Runner) Stop() {
	r.logger.Info("Stopping runner...")
	r.watcher.Stop()
	<-r.watcher.SndDoneCh // NOTE: Might need a timeout to prevent blocking forever
	close(r.RcvDoneCh)
}
