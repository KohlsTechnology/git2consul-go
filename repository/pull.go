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

		// This also creates the branch
		localBranchRef, err := r.References.Lookup(rawLocalBranchRefs[i])
		if err != nil {
			localBranchRef, err = r.References.Create(rawLocalBranchRefs[i], remoteBranchRef.Target(), true, "")
			if err != nil {
				return err
			}
		}

		//Fast-foward changes and checkout
		// log.Debugf("=== %s", rawRemoteBranchRefs[i])
		err = r.SetHead(rawLocalBranchRefs[i])
		if err != nil {
			return err
		}
		err = r.CheckoutHead(&git.CheckoutOpts{})
		if err != nil {
			return err
		}

		// h, _ := r.Head()
		// log.Debugf("=== Head: %s", h.Target().String())

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
			bn, _ := localBranchRef.Branch().Name()
			log.Infof("Skipping pull on repository %s, branch %s. Already up to date", r.repoConfig.Name, bn)
		} else if analysis&git.MergeAnalysisFastForward != 0 {
			bn, _ := localBranchRef.Branch().Name()
			log.Infof("Changes detected on repository %s branch %s, Fast-forwarding", r.repoConfig.Name, bn)
			if err := r.Merge(mergeHeads, nil, nil); err != nil {
				return err
			}

			localBranchRef.SetTarget(remoteBranchRef.Target(), "")

			r.StateCleanup()

			// Point branch to the object
			// _, err := localBranchRef.SetTarget(remoteBranchRef.Target(), "")
			// if err != nil {
			// 	return err
			// }
			//
			// //Checkout head
			// err = r.SetHead(rawLocalBranchRefs[i])
			// if err != nil {
			// 	return err
			// }
			//
			// err = r.CheckoutHead(&git.CheckoutOpts{})
			// if err != nil {
			// 	return err
			// }
		} else if analysis&git.MergeAnalysisNormal != 0 {
			// Just merge changes
			bn, _ := localBranchRef.Branch().Name()
			log.Infof("Changes detected on repository %s. Pulling commits from branch %s", r.repoConfig.Name, bn)
			// if err := r.Merge(mergeHeads, nil, nil); err != nil {
			// 	return err
			// }
			//
			// // Check for conflicts
			// index, err := r.Index()
			// if err != nil {
			// 	return err
			// }
			// defer index.Free()
			//
			// if index.HasConflicts() {
			// 	iter, err := index.ConflictIterator()
			// 	if err != nil {
			// 		return fmt.Errorf("Could not create iterator for conflicts: %s", err.Error())
			// 	}
			// 	defer iter.Free()
			//
			// 	for entry, err := iter.Next(); err != nil; entry, err = iter.Next() {
			// 		fmt.Printf("CONFLICT: %s\n", entry.Our.Path)
			// 	}
			// 	return errors.New("Conflicts encountered. Please resolve them.")
			// }
			//
			// // Make the merge commit
			// sig, err := r.DefaultSignature()
			// if err != nil {
			// 	return err
			// }
			//
			// // Get Write Tree
			// treeId, err := index.WriteTree()
			// if err != nil {
			// 	return err
			// }
			//
			// tree, err := r.LookupTree(treeId)
			// if err != nil {
			// 	return err
			// }
			// defer tree.Free()
			//
			// localCommit, err := r.LookupCommit(localBranchRef.Target())
			// if err != nil {
			// 	return err
			// }
			// defer localCommit.Free()
			//
			// remoteCommit, err := r.LookupCommit(remoteBranchRef.Target())
			// if err != nil {
			// 	return err
			// }
			// defer remoteCommit.Free()
			//
			// _, err = r.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)
			// if err != nil {
			// 	return fmt.Errorf("could not create commit after merge: %s", err.Error())
			// }
			//
			// // Clean up
			// r.StateCleanup()
		}
	}

	return nil
}
