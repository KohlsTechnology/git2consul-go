package watch

import (
	"time"

	"github.com/cleung2010/go-git2consul/repository"
	"gopkg.in/libgit2/git2go.v24"
)

func pollByInterval(repo *repository.Repository) {
	hooks := repo.repoConfig.Hooks
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
	defer ticker.Stop()

	for {
		err := repo.PollBranches()
		if err != nil {
			errCh <- err
		}

		select {
		case <-ticker.C:
		}
	}
}

func pollBranches(repo *repository.Repository) {
	itr, err := repo.NewBranchIterator(git.LocalBranch)

	return nil
}
