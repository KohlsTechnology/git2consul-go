package repository

import (
	"time"

	log "github.com/Sirupsen/logrus"
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
