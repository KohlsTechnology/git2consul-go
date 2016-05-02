package repository

import (
	"errors"
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

// Pull the repository, which is a fetch and merge
// It attempts to pull all branches specified in the repository configuration
func (r *Repository) Pull() error {
	origin, err := r.Remotes.Lookup("origin")
	if err != nil {
		return err
	}
	defer origin.Free()

	headBranchRefs, remoteBranchRefs := buildRefs(r.repoConfig.Branches)
	err = origin.Fetch(remoteBranchRefs, nil, "") // Should I fetch from heads or from remotes?
	if err != nil {
		return err
	}

	// Iterate on remote refs, and do a merge
	// TODO: Figure out how to do multi-branch merge
	log.Debugf("All remote refs: %q", remoteBranchRefs)
	for i, _ := range remoteBranchRefs {
		log.Debugf("Current remote ref: %s | Index: %d", remoteBranchRefs[i], i)
		remoteBranchRef, err := r.References.Lookup(remoteBranchRefs[i])
		if err != nil {
			log.Debug("=== 1")
			return err
		}

		annotatedCommit, err := r.AnnotatedCommitFromRef(remoteBranchRef)
		if err != nil {
			return err
		}

		// If branch does not exist, create it
		// if !localBranchRef.IsBranch() {
		// bn, err := localBranchRef.Branch().Name()
		// if err != nil {
		// 	return err
		// }
		remoteCommit, err := r.LookupCommit(remoteBranchRef.Target())
		if err != nil {
			return err
		}

		_, err = r.CreateBranch(r.repoConfig.Branches[i], remoteCommit, true)
		if err != nil {
			return err
		}
		// }

		// Merge analysis
		mergeHeads := []*git.AnnotatedCommit{annotatedCommit}
		analysis, _, err := r.MergeAnalysis(mergeHeads)
		if err != nil {
			return err
		}

		localBranchRef, err := r.References.Lookup(headBranchRefs[i])
		if err != nil {
			return err
		}

		// Actions to take depending on analysis outcome
		if analysis&git.MergeAnalysisUpToDate != 0 {
			log.Infof("Skipping pull on repository %s. Already up to date", r.repoConfig.Name)
			// return nil
		} else if analysis&git.MergeAnalysisNormal != 0 {
			// Just merge changes
			log.Infof("Changes detected on repository %s. Pulling commits from branch %s", r.repoConfig.Name, remoteBranchRef)
			if err := r.Merge(mergeHeads, nil, nil); err != nil {
				return err
			}

			// Check for conflicts
			index, err := r.Index()
			if err != nil {
				return err
			}
			defer index.Free()

			if index.HasConflicts() {
				iter, err := index.ConflictIterator()
				if err != nil {
					return fmt.Errorf("Could not create iterator for conflicts: %s", err.Error())
				}
				defer iter.Free()

				for entry, err := iter.Next(); err != nil; entry, err = iter.Next() {
					fmt.Printf("CONFLICT: %s\n", entry.Our.Path)
				}
				return errors.New("Conflicts encountered. Please resolve them.")
			}

			// Make the merge commit
			sig, err := r.DefaultSignature()
			if err != nil {
				return err
			}

			// Get Write Tree
			treeId, err := index.WriteTree()
			if err != nil {
				return err
			}

			tree, err := r.LookupTree(treeId)
			if err != nil {
				return err
			}
			defer tree.Free()

			localCommit, err := r.LookupCommit(localBranchRef.Target())
			if err != nil {
				return err
			}
			defer localCommit.Free()

			remoteCommit, err := r.LookupCommit(remoteBranchRef.Target())
			if err != nil {
				return err
			}
			defer remoteCommit.Free()

			_, err = r.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)
			if err != nil {
				return fmt.Errorf("could not create commit after merge: %s", err.Error())
			}

			// Clean up
			r.StateCleanup()
		} else if analysis&git.MergeAnalysisFastForward != 0 {
			// Fast-forward changes
			// Get remote tree
			remoteTree, err := r.LookupTree(remoteBranchRef.Target())
			if err != nil {
				return err
			}

			// Checkout
			if err := r.CheckoutTree(remoteTree, nil); err != nil {
				return err
			}

			branchRef, err := r.References.Lookup(headBranchRefs[i])
			if err != nil {
				return err
			}

			// Point branch to the object
			branchRef.SetTarget(remoteBranchRef.Target(), "")
			if _, err := localBranchRef.SetTarget(remoteBranchRef.Target(), ""); err != nil {
				return err
			}
		}
	}

	return nil
}
