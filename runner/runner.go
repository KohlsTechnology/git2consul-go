package runner

import (
	"fmt"

	"github.com/apex/log"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/cleung2010/go-git2consul/kv"
	"github.com/cleung2010/go-git2consul/repository"
	"github.com/cleung2010/go-git2consul/repository/watch"
)

type Runner struct {
	ErrCh  chan error
	DoneCh chan struct{}

	once bool

	kvHandler *kv.KVHandler

	repos []*repository.Repository
}

func NewRunner(config *config.Config, once bool) (*Runner, error) {
	// Create repos from configuration
	rs, err := repository.LoadRepos(config)
	if err != nil {
		return nil, fmt.Errorf("Cannot load repositories from configuration: %s", err)
	}

	// Create the handler
	handler, err := kv.New(config.Consul)
	if err != nil {
		return nil, err
	}

	runner := &Runner{
		ErrCh:     make(chan error),
		DoneCh:    make(chan struct{}, 1),
		once:      once,
		kvHandler: handler,
		repos:     rs,
	}

	return runner, nil
}

// Start the runner
func (r *Runner) Start() {
	rw := watch.New(r.repos)
	rw.Watch()

	// Grab changes from repositories. Do no stop the runner if there
	// are errors on the repo watcher
	for {
		select {
		case err := <-rw.ErrCh:
			log.WithError(err).Error("Watcher error")
		case repo := <-rw.RepoChangeCh:
			// Handle change, and return if error on handler
			err := r.kvHandler.HandleUpdate(repo)
			if err != nil {
				r.ErrCh <- err
				return
			}
		}
	}

	// FIXME: This doesn't work atm. Probably needs donCh on watches to block
	// until underlying goroutines are done before we can report back to r.DoneCh
	if r.once {
		r.DoneCh <- struct{}{}
	}
}
