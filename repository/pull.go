package repository

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/libgit2/git2go"
)

func buildRefs(branches []string) ([]string, []string) {
	headBranchRef := []string{}
	remoteBranchRef := []string{}

	for _, branch := range branches {
		hbr := fmt.Sprintf("refs/heads/%s", branch)
		headBranchRef = append(headBranchRef, hbr)
		rbr := fmt.Sprintf("refs/remotes/origin/%s", branch)
		remoteBranchRef = append(remoteBranchRef, rbr)
	}

	return headBranchRef, remoteBranchRef
}

func mergeBranch() {

}

// Pull a repository branch
func (r *Repository) PullOne(branchName string) error {
	origin, err := r.Remotes.Lookup("origin")
	if err != nil {
		return err
	}
	defer origin.Free()

	rawLocalBranchRef := fmt.Sprintf("refs/heads/%s", branchName)
	rawRemoteBranchRef := fmt.Sprintf("refs/remotes/origin/%s", branchName)

	// Fetch
	err = origin.Fetch([]string{rawLocalBranchRef}, nil, "")
	if err != nil {
		return err
	}

	remoteBranchRef, err := r.References.Lookup(rawRemoteBranchRef)
	if err != nil {
		return err
	}

	// If the ref on the branch doesn't exist locally, create it
	// This also creates the branch
	localBranchRef, err := r.References.Lookup(rawLocalBranchRef)
	if err != nil {
		localBranchRef, err = r.References.Create(rawLocalBranchRef, remoteBranchRef.Target(), true, "")
		if err != nil {
			return err
		}
	}

	err = r.SetHead(rawLocalBranchRef)
	if err != nil {
		return err
	}
	err = r.CheckoutHead(&git.CheckoutOpts{})
	if err != nil {
		return err
	}

	// Create annotated commit
	annotatedCommit, err := r.AnnotatedCommitFromRef(remoteBranchRef)
	if err != nil {
		return err
	}

	// Merge analysis
	mergeHeads := []*git.AnnotatedCommit{annotatedCommit}
	analysis, _, err := r.MergeAnalysis(mergeHeads)
	if err != nil {
		return err
	}

	// Action on analysis
	if analysis&git.MergeAnalysisUpToDate != 0 { // On up-to-date merge
		log.Infof("Skipping pull on repository %s, branch %s. Already up to date", r.repoConfig.Name, branchName)
	} else if analysis&git.MergeAnalysisFastForward != 0 { // On fast-forward merge
		log.Infof("Changes detected on repository %s branch %s, Fast-forwarding", r.repoConfig.Name, branchName)

		err := r.Merge(mergeHeads, nil, nil)
		if err != nil {
			return err
		}

		localBranchRef.SetTarget(remoteBranchRef.Target(), "")
	} else if analysis&git.MergeAnalysisNormal != 0 { // On normal merge
		log.Infof("Changes detected on repository %s. Pulling commits from branch %s", r.repoConfig.Name, branchName)

		if err := r.Merge(mergeHeads, nil, nil); err != nil {
			return err
		}

		localBranchRef.SetTarget(remoteBranchRef.Target(), "")
	}

	r.StateCleanup()
	return nil
}

// Pull the repository, which is a fetch and merge
// It attempts to pull all branches specified in the repository configuration
func (r *Repository) Pull() error {
	origin, err := r.Remotes.Lookup("origin")
	if err != nil {
		return err
	}
	defer origin.Free()

	// Build ref paths
	rawLocalBranchRefs, rawRemoteBranchRefs := buildRefs(r.repoConfig.Branches)

	// Fetch
	err = origin.Fetch(rawLocalBranchRefs, nil, "")
	if err != nil {
		return err
	}

	// Iterate on remote refs for specified branches and do a merge
	for i, _ := range rawRemoteBranchRefs {
		remoteBranchRef, err := r.References.Lookup(rawRemoteBranchRefs[i])
		if err != nil {
			return err
		}

		// If the ref on the branch doesn't exist locally, create it
		// This also creates the branch
		localBranchRef, err := r.References.Lookup(rawLocalBranchRefs[i])
		if err != nil {
			localBranchRef, err = r.References.Create(rawLocalBranchRefs[i], remoteBranchRef.Target(), true, "")
			if err != nil {
				return err
			}
		}

		// Checkout correct branch
		// TODO: Thread safe, so that manual checkout does not mess up the analysis
		err = r.SetHead(rawLocalBranchRefs[i])
		if err != nil {
			return err
		}
		err = r.CheckoutHead(&git.CheckoutOpts{})
		if err != nil {
			return err
		}

		// Create annotated commit
		annotatedCommit, err := r.AnnotatedCommitFromRef(remoteBranchRef)
		if err != nil {
			return err
		}

		// Merge analysis
		mergeHeads := []*git.AnnotatedCommit{annotatedCommit}
		analysis, _, err := r.MergeAnalysis(mergeHeads)
		if err != nil {
			return err
		}

		// Action on analysis
		if analysis&git.MergeAnalysisUpToDate != 0 { // On up-to-date merge
			bn, _ := localBranchRef.Branch().Name()
			log.Infof("Skipping pull on repository %s, branch %s. Already up to date", r.repoConfig.Name, bn)
		} else if analysis&git.MergeAnalysisFastForward != 0 { // On fast-forward merge
			bn, _ := localBranchRef.Branch().Name()
			log.Infof("Changes detected on repository %s branch %s, Fast-forwarding", r.repoConfig.Name, bn)

			err := r.Merge(mergeHeads, nil, nil)
			if err != nil {
				return err
			}

			localBranchRef.SetTarget(remoteBranchRef.Target(), "")
		} else if analysis&git.MergeAnalysisNormal != 0 { // On normal merge
			// Just merge changes
			bn, _ := localBranchRef.Branch().Name()
			log.Infof("Changes detected on repository %s. Pulling commits from branch %s", r.repoConfig.Name, bn)

			if err := r.Merge(mergeHeads, nil, nil); err != nil {
				return err
			}

			localBranchRef.SetTarget(remoteBranchRef.Target(), "")
		}
	}

	r.StateCleanup()
	return nil
}
