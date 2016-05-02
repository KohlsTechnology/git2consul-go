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

	rawLocalBranchRefs, rawRemoteBranchRefs := buildRefs(r.repoConfig.Branches)
	err = origin.Fetch(rawRemoteBranchRefs, nil, "") // Should I fetch from heads or from remotes?
	if err != nil {
		return err
	}

	// Iterate on remote refs for specified branches and do a merge
	log.Debugf("All remote refs: %q", rawRemoteBranchRefs)
	for i, _ := range rawRemoteBranchRefs {
		log.Debugf("Current remote ref: %s | Index: %d", rawRemoteBranchRefs[i], i)
		remoteBranchRef, err := r.References.Lookup(rawRemoteBranchRefs[i])
		if err != nil {
			return err
		}
		log.Debugf("Remote ref: %s", remoteBranchRef.Target())

		// This also creates the branch
		localBranchRef, err := r.References.Lookup(rawLocalBranchRefs[i])
		if err != nil {
			localBranchRef, err = r.References.Create(rawLocalBranchRefs[i], remoteBranchRef.Target(), true, "")
			if err != nil {
				return err
			}
		}

		//Fast-foward changes and checkout
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

		// Actions to take depending on analysis outcome
		if analysis&git.MergeAnalysisUpToDate != 0 {
			log.Infof("Skipping pull on repository %s. Already up to date", r.repoConfig.Name)
			// return nil
		} else if analysis&git.MergeAnalysisNormal != 0 {
			// Just merge changes
			bn, _ := localBranchRef.Branch().Name()
			log.Infof("Changes detected on repository %s. Pulling commits from branch %s", r.repoConfig.Name, bn)
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

			branchRef, err := r.References.Lookup(rawLocalBranchRefs[i])
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
