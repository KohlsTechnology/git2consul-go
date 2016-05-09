package repository

import "time"

// Watch for changes on the remotes
func (rs Repositories) WatchRepos() error {
	errCh := make(chan error, 1)

	// Poll repository by interval, or webhook
	for _, repo := range rs {
		// Initial poll
		err := repo.pollBranches()
		if err != nil {
			return err
		}

		go repo.pollRepoByInterval(errCh)
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

func (r *Repository) pollRepoByInterval(errCh chan error) {
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
				errCh <- err
				//return
			}
		}
	}
}

func (r *Repository) pollRepoByWebhook() {

}
