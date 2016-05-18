package runner

import (
	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/repository"
)

// Method that watches for changes in local repositories and performs
// Consul operations to update the KV
func (r *Runner) watchKVUpdate() {
	// If there changes, push to KV
	for _, repo := range r.repos {
		go r.watchLocalRepo(repo)
	}
}

// Helper fuction for watchKVUpdate() to watch a specific
// local repository. This should be ran as a goroutine since it blocks and
// waits for the changeCh
func (r *Runner) watchLocalRepo(repo *repository.Repository) {
	// Initial update to the KV
	err := r.initHandler(repo)
	if err != nil {
		r.ErrCh <- err
		return
	}

	for {
		err := r.updateHandler(repo)
		if err != nil {
			r.ErrCh <- err
			return
		}

		// Block until changes
		select {
		case <-repo.ChangeCh():
		}
	}
}

// Method that watches for changes in remote repositories and performs
// git operations to update local copy
func (r *Runner) watchReposUpdate() {
	for _, repo := range r.repos {
		go r.watchRemoteRepo(repo)
	}
}

// Helper function for watchReposUpdate() to watch a specific
// remote repository changes. This should be ran as a goroutine
func (r *Runner) watchRemoteRepo(repo *repository.Repository) {
	errCh := make(chan error)
	go repo.PollRepoByInterval(errCh)
	// go r.PollRepoByWebhook()

	// FIXME: Handle errors better
	for {
		select {
		case err := <-errCh:
			log.Debugf("(git): %s", err)
		}
	}
}
