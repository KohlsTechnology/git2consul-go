package runner

import (
	"fmt"

	"github.com/Cimpress-MCP/go-git2consul/config"
	"github.com/Cimpress-MCP/go-git2consul/kv"
	"github.com/Cimpress-MCP/go-git2consul/repository"
	"github.com/Cimpress-MCP/go-git2consul/watcher"
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

	kvHandler *kv.KVHandler

	watcher *watch.Watcher
}

// NewRunner creates a new runner instance
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
			err := r.kvHandler.HandleUpdate(repo)
			if err != nil {
				r.ErrCh <- err
				return
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
