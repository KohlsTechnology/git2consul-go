package repository

import (
	"errors"
	"fmt"

	"github.com/libgit2/git2go"
)

// Pull the repository, which is a fetch and merge
func (r *Repository) Pull() error {
	origin, err := r.Remotes.Lookup("origin")
	if err != nil {
		return err
	}
	defer origin.Free()

	// TODO: Accept branches other than master
	err = origin.Fetch([]string{"refs/heads/master"}, nil, "")
	if err != nil {
		return err
	}

	// TODO: Accept branches other than master
	remoteBranch, err := r.References.Lookup("refs/remotes/origin/master")
	if err != nil {
		return err
	}

	head, err := r.Head()
	if err != nil {
		return err
	}

	remoteBranchID := remoteBranch.Target()
	annotatedCommit, err := r.AnnotatedCommitFromRef(remoteBranch)
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
		fmt.Println("Already up to date.")
		return nil
	} else if analysis&git.MergeAnalysisNormal != 0 {
		// Just merge changes
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
				return fmt.Errorf("could not create iterator for conflicts: %s", err.Error())
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

		localCommit, err := r.LookupCommit(head.Target())
		if err != nil {
			return err
		}
		defer localCommit.Free()

		remoteCommit, err := r.LookupCommit(remoteBranchID)
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
		remoteTree, err := r.LookupTree(remoteBranchID)
		if err != nil {
			return err
		}

		// Checkout
		if err := r.CheckoutTree(remoteTree, nil); err != nil {
			return err
		}

		branchRef, err := r.References.Lookup("refs/heads/master") // TODO: not just master
		if err != nil {
			return err
		}

		// Point branch to the object
		branchRef.SetTarget(remoteBranchID, "")
		if _, err := head.SetTarget(remoteBranchID, ""); err != nil {
			return err
		}

	}

	return nil
}
