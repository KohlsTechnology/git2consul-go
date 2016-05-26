package watch

import (
	"sync"

	"github.com/apex/log"
	"github.com/cleung2010/go-git2consul/repository"
)

type Watcher struct {
	sync.Mutex
	logger *log.Entry

	Repositories []*repository.Repository

	RepoChangeCh chan *repository.Repository
	ErrCh        chan error
	DoneCh       chan struct{}

	webhookPort int
	once        bool
}

func New(repos []*repository.Repository, webhookPort int) *Watcher {
	repoChangeCh := make(chan *repository.Repository, len(repos))
	errCh := make(chan error)
	doneCh := make(chan struct{}, 1)

	logger := log.WithField("caller", "git")

	return &Watcher{
		Repositories: repos,
		RepoChangeCh: repoChangeCh,
		ErrCh:        errCh,
		DoneCh:       doneCh,
		logger:       logger,
		webhookPort:  webhookPort,
		once:         false,
	}
}

func (w *Watcher) Watch(once bool) {
	//errsCh := make(chan error, len(w.Repositories)) // Error channel for all watching repos

	for _, repo := range w.Repositories {
		go w.pollByInterval(repo)
	}

	go w.pollByWebhook()

	for {
		select {
		case err := <-w.ErrCh: // FIXME: This is already handled in runner!
			log.WithError(err).Error("Watch error")
		case <-w.DoneCh:
			log.Info("Received finish")
			return
		}
		if once {
			w.Stop()
		}
	}

}

// Stop watching for changes
func (w *Watcher) Stop() {
	w.logger.Info("Stopping watcher")
	close(w.DoneCh)
}
