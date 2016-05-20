package watch

import "github.com/cleung2010/go-git2consul/repository"

func (w *Watcher) watchRepo(repo *repository.Repository, errCh chan<- error) {
	go w.pollByInterval(repo, errCh)
	// pollByWebhook(repo)

}

// TODO: Should return error, which is an array of all the errors
func (w *Watcher) Watch() {
	//errsCh := make(chan error, len(w.Repositories)) // Error channel for all watching repos

	for _, repo := range w.Repositories {
		go w.watchRepo(repo, w.ErrCh)
	}

	// for {
	// 	select {
	// 	case err := <-errsCh:
	// 		log.WithError(err).Error("watch error")
	// 	}
	// }

}
