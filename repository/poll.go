package repository

import "time"

// Poll repository once. Polling can either clone or update
func (r *Repository) PollBranches() error {
	for _, branch := range r.repoConfig.Branches {
		err := r.pull(branch)
		if err != nil {
			return err
		}
	}

	return nil
}

// Poll a repository at a specified interval. This should be called as a
// goroutine since the time ticker blocks.
func (r *Repository) PollRepoByInterval(errCh chan error) {
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
			err := r.PollBranches()
			if err != nil {
				errCh <- err
				//return // Do not return to attempt retrying
			}
		}
	}
}

// Poll a repository by webhooks. This should be called as a go routine
func (r *Repository) PollRepoByWebhooks() {

}
