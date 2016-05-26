package watch

import (
	"sync"

	"github.com/apex/log"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/cleung2010/go-git2consul/repository"
)

type Watcher struct {
	sync.Mutex
	logger *log.Entry

	Repositories []*repository.Repository

	RepoChangeCh chan *repository.Repository
	ErrCh        chan error
	DoneCh       chan struct{}

	hookSvr *config.HookSvrConfig
	once    bool
}

// Create a new watcher, passing in the the repositories, webhook
// listener config, and optional once flag
func New(repos []*repository.Repository, hookSvr *config.HookSvrConfig, once bool) *Watcher {
	repoChangeCh := make(chan *repository.Repository, len(repos))
	errCh := make(chan error)
	doneCh := make(chan struct{}, 1)
	logger := log.WithField("caller", "watcher")

	return &Watcher{
		Repositories: repos,
		RepoChangeCh: repoChangeCh,
		ErrCh:        errCh,
		DoneCh:       doneCh,
		logger:       logger,
		hookSvr:      hookSvr,
		once:         once,
	}
}

// Watch repositories available to the watcher
func (w *Watcher) Watch() {
	for _, repo := range w.Repositories {
		go w.pollByInterval(repo)
	}

	go w.pollByWebhook()

	for {
		select {
		case err := <-w.ErrCh:
			log.WithError(err).Error("Watcher error")
		case <-w.DoneCh:
			w.logger.Info("Received finish")
			return
		}
	}

}

// Stop watching for changes. It will stop interval and webhook polling
func (w *Watcher) Stop() {
	w.logger.Info("Stopping watcher")
	close(w.DoneCh)
}
