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
	logger *log.Entry
	ErrCh  chan error
	DoneCh chan struct{}

	once bool

	kvHandler *kv.KVHandler

	watcher *watch.Watcher
}

func NewRunner(config *config.Config, once bool) (*Runner, error) {
	logger := log.WithField("caller", "runner")

	// Create repos from configuration
	repos, err := repository.LoadRepos(config)
	if err != nil {
		return nil, fmt.Errorf("Cannot load repositories from configuration: %s", err)
	}

	// Create watcher to watch for repo changes
	port := config.WebhookPort
	watcher := watch.New(repos, port)

	// Create the handler
	handler, err := kv.New(config.Consul)
	if err != nil {
		return nil, err
	}

	runner := &Runner{
		logger:    logger,
		ErrCh:     make(chan error),
		DoneCh:    make(chan struct{}, 1),
		once:      once,
		kvHandler: handler,
		watcher:   watcher,
	}

	return runner, nil
}

// Start the runner
func (r *Runner) Start() {
	r.watcher.Watch(r.once)

	for {
		select {
		case err := <-r.watcher.ErrCh:
			log.WithError(err).Error("Watcher error")
			// Do no stop the runner if there are errors on the repo watcher
		case repo := <-r.watcher.RepoChangeCh:
			// Handle change, and return if error on handler
			err := r.kvHandler.HandleUpdate(repo)
			if err != nil {
				r.ErrCh <- err
				return
			}
		case <-r.watcher.DoneCh:
			log.Info("Watcher reported finish")
			r.Stop()
		case <-r.DoneCh:
			log.Info("Received finish")
			return
		}
	}

	// FIXME: This doesn't work atm. Probably needs donCh on watches to block
	// until underlying goroutines are done before we can report back to r.DoneCh
	if r.once {
		r.DoneCh <- struct{}{}
	}
}

func (r *Runner) Stop() {
	r.logger.Info("Stopping runner")
	r.watcher.Stop()
	close(r.DoneCh)
}
