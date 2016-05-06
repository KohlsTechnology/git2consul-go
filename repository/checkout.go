package repository

import (
	"path"

	"github.com/libgit2/git2go"
)

// Checkout all remote branches
func (r *Repository) checkoutRemoteBranches() error {
	itr, err := r.NewBranchIterator(git.BranchRemote)
	if err != nil {
		return err
	}
	defer itr.Free()

	var checkoutBranchFn = func(b *git.Branch, _ git.BranchType) error {
		bn, err := b.Name()
		if err != nil {
			return err
		}
		_, err = r.References.Lookup("refs/heads/" + path.Base(bn))
		if err != nil {
			localRef, err := r.References.Create("refs/heads/"+path.Base(bn), b.Reference.Target(), true, "")
			if err != nil {
				return err
			}

			err = r.SetHead(localRef.Name())
			if err != nil {
				return err
			}

			err = r.CheckoutHead(&git.CheckoutOpts{
				Strategy: git.CheckoutForce,
			})
			if err != nil {
				return err
			}

			// r.changeCh <- struct{}{}
		}

		return nil
	}

	itr.ForEach(checkoutBranchFn)

	return nil
}

func (r *Repository) CheckoutBranch(branch *git.Branch, opts *git.CheckoutOpts) error {
	err := r.SetHead(branch.Reference.Name())
	if err != nil {
		return err
	}

	err = r.CheckoutHead(opts)
	if err != nil {
		return err
	}

	return nil
}
