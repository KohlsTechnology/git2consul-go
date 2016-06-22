package repository

import (
	"fmt"

	"gopkg.in/libgit2/git2go.v24"
)

// Clone the repository. Cloning will only checkout tracked branches.
// A destination path to clone to needs to be provided
func (r *Repository) Clone(path string) error {
	r.Lock()
	defer r.Unlock()

	// Clone the first tracked branch instead of the default branch
	if len(r.Config.Branches) == 0 {
		return fmt.Errorf("No tracked branches specified")
	}
	checkoutBranch := r.Config.Branches[0]

	raw_repo, err := git.Clone(r.Config.Url, path, &git.CloneOptions{
		CheckoutOpts: &git.CheckoutOpts{
			Strategy: git.CheckoutNone,
		},
		CheckoutBranch: checkoutBranch,
	})
	if err != nil {
		return err
	}

	r.Repository = raw_repo

	err = r.checkoutConfigBranches()
	if err != nil {
		return err
	}

	return nil
}
