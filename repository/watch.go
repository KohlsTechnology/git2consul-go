package repository

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/libgit2/git2go"
)

// Watch for changes on the remotes
func (rs Repositories) WatchRepos() error {
	// Poll repository by interval, or webhook
	for _, repo := range rs {
		go repo.pollRepoByInterval()
		// go r.PollRepoByWebhook()
	}

	return nil
}

// Poll repository once. Polling can either clone or update
func (r *Repository) poll() error {
	if _, err := git.OpenRepository(r.store); err != nil {
		log.Infof("Repository %s not cached, cloning to %s", r.repoConfig.Name, r.store)
		err := os.Mkdir(r.store, 0755)
		if err != nil {
			return err
		}

		err = r.Clone()
		if err != nil {
			return err
		}
	}

	for _, branch := range r.repoConfig.Branches {
		r.Pull(branch)
	}

	return nil
}

func (r *Repository) pollRepoByInterval() {
	hooks := r.repoConfig.Hooks
	interval := time.Second

	// Find polling hook
	for _, h := range hooks {
		if h.Type == "polling" {
			interval = h.Interval
			break
		}
	}

	// If no polling found, don't poll
	if interval == 0 {
		return
	}

	// Initial poll
	err := r.poll()
	if err != nil {
		log.Error(err)
	}

	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-ticker.C:
			err := r.poll()
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (r *Repository) pollRepoByWebhook() {

}
