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
	"sync"

	"github.com/KohlsTechnology/git2consul-go/config"
	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/apex/log"
)

// Watcher is used to keep track of changes of the repositories
type Watcher struct {
	sync.Mutex
	logger *log.Entry

	Repositories []repository.Repo

	RepoChangeCh chan repository.Repo
	ErrCh        chan error
	RcvDoneCh    chan struct{}
	SndDoneCh    chan struct{}

	hookSvr *config.HookSvrConfig
	once    bool
}

// New create a new watcher, passing in the repositories, webhook
// listener config, and optional once flag
func New(repos []repository.Repo, hookSvr *config.HookSvrConfig, once bool) *Watcher {
	repoChangeCh := make(chan repository.Repo, len(repos))
	logger := log.WithField("caller", "watcher")

	return &Watcher{
		Repositories: repos,
		RepoChangeCh: repoChangeCh,
		ErrCh:        make(chan error),
		RcvDoneCh:    make(chan struct{}, 1),
		SndDoneCh:    make(chan struct{}, 1),
		logger:       logger,
		hookSvr:      hookSvr,
		once:         once,
	}
}

// Watch repositories available to the watcher
func (w *Watcher) Watch() {
	defer close(w.SndDoneCh)

	// Pass repositories to RepoChangeCh for initial update to the KV
	for _, repo := range w.Repositories {
		w.RepoChangeCh <- repo
	}

	// WaitGroup size is equal to number of interval goroutine plus webhook goroutine
	var wg sync.WaitGroup
	wg.Add(len(w.Repositories) + 1)

	for _, repo := range w.Repositories {
		go w.pollByInterval(repo, &wg)
	}

	go w.pollByWebhook(&wg)

	go func() {
		wg.Wait()
		// Only exit if it's -once, otherwise there might be webhook polling
		if w.once {
			w.Stop()
			return
		}
	}()

	for {
		select {
		case err := <-w.ErrCh:
			log.WithError(err).Error("Watcher error")
		case <-w.RcvDoneCh:
			w.logger.Info("Received finish")
			wg.Wait()
			return
		}
	}
}

// Stop watching for changes. It will stop interval and webhook polling
func (w *Watcher) Stop() {
	w.logger.Info("Stopping watcher...")
	close(w.RcvDoneCh)
}
