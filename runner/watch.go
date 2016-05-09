package runner

import "github.com/cleung2010/go-git2consul/repository"

func (r *Runner) watchKVUpdate() {
	// If there changes, push to KV
	for _, repo := range r.repos {
		go r.watchLocalRepo(repo)
	}
}

func (r *Runner) watchLocalRepo(repo *repository.Repository) {
	// Initial update to the KV
	err := r.initHandler(repo)
	if err != nil {
		r.ErrCh <- err
		return
	}

	for {
		select {
		case <-repo.ChangeCh():
			err := r.updateHandler(repo)
			if err != nil {
				r.ErrCh <- err
				return
			}
		}
	}
}

func (r *Runner) watchReposUpdate() {
	r.repos.WatchRepos()
}
