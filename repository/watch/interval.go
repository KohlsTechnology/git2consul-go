package watch

import (
	"fmt"
	"sync"
	"time"

	"github.com/cleung2010/go-git2consul/repository"
	"gopkg.in/libgit2/git2go.v24"
)

// Watch the repo by interval. This is called as a go routine since
// ticker blocks
func (w *Watcher) pollByInterval(repo *repository.Repository, wg *sync.WaitGroup) {
	defer wg.Done()

	hooks := repo.Hooks
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

	// Polling error should not stop polling by interval
	for {
		err := w.pollBranches(repo)
		if err != nil {
			w.ErrCh <- err
		}

		if w.once {
			return
		}

		select {
		case <-ticker.C:
		case <-w.RcvDoneCh:
			return
		}
	}
}

// Watch all branches of a repository
func (w *Watcher) pollBranches(repo *repository.Repository) error {
	itr, err := repo.NewBranchIterator(git.BranchLocal)
	if err != nil {
		return err
	}
	defer itr.Free()

	var checkoutBranchFn = func(b *git.Branch, _ git.BranchType) error {
		branchName, err := b.Name()
		if err != nil {
			return err
		}
		analysis, err := repo.Pull(branchName)
		if err != nil {
			return err
		}

		// If there is a change, send the repo RepoChangeCh
		switch {
		case analysis&git.MergeAnalysisUpToDate != 0:
			w.logger.Debugf("Up to date: %s/%s", repo.Name(), branchName)
		case analysis&git.MergeAnalysisNormal != 0, analysis&git.MergeAnalysisFastForward != 0:
			w.logger.Infof("Changed: %s/%s", repo.Name(), branchName)
			w.RepoChangeCh <- repo
		}

		return nil
	}

	// HACK: itr.ForEach() doesn't return nil, but instead an empty string
	//       See: https://github.com/libgit2/git2go/issues/315
	err = itr.ForEach(checkoutBranchFn)
	if err != nil && len(err.Error()) > 0 {
		fmt.Printf("====== itr err: %s", err.Error())
	}

	return nil
}
