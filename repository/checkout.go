package repository

import (
	"path"

	"gopkg.in/libgit2/git2go.v23"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Checkout branches specified in the config
func (r *Repository) checkoutConfigBranches() error {
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

		// Only checkout tracked branches
		// TODO: optimize this O(n^2)
		if stringInSlice(path.Base(bn), r.repoConfig.Branches) == false {
			return nil
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
