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
	watcher := watch.New(repos, config.HookSvr, once)

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
	go r.watcher.Watch()

	for {
		select {
		case repo := <-r.watcher.RepoChangeCh:
			// Handle change, and return if error on handler
			err := r.kvHandler.HandleUpdate(repo)
			if err != nil {
				r.ErrCh <- err
				return
			}
		case <-r.DoneCh:
			r.logger.Info("Received finish")
			return
		case <-r.watcher.DoneCh: // Mainly for -once, in this case we don't need to call w.Stop()
			close(r.DoneCh)
		}
	}
}

// Stop the runner, cleaning up any routines that it's running. In this case, it will stop
// the watcher before closing DoneCh
func (r *Runner) Stop() {
	r.logger.Info("Stopping runner")
	r.watcher.Stop()
	close(r.DoneCh)
}
