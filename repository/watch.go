package repository

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// Watch for changes on the remotes
func (rs Repositories) WatchRepos() error {
	// Poll repository by interval, or webhook
	for _, repo := range rs {
		// Initial poll
		// err := repo.pollBranches()
		// if err != nil {
		// 	log.Error(err)
		// }

		go repo.pollRepoByInterval()
		// go r.PollRepoByWebhook()
	}

	return nil
}

// Poll repository once. Polling can either clone or update
func (r *Repository) pollBranches() error {
	for _, branch := range r.repoConfig.Branches {
		r.pull(branch)
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

	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-ticker.C:
			err := r.pollBranches()
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (r *Repository) pollRepoByWebhook() {

}
