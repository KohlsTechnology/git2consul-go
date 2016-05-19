package watch

import "github.com/cleung2010/go-git2consul/repository"

func (w *Watcher) watchRepo(repo *repository.Repository) {
	pollByInterval(repo)
	// pollByWebhook(repo)
}

func (w *Watcher) Watch(repos []*repository.Repository) {
	// errorsCh := make(chan error, len(repos))

	for _, repo := range repos {
		go w.watchRepo(repo)
	}
}
