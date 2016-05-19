package repository

import "gopkg.in/libgit2/git2go.v24"

// Clone the repository
func (r *Repository) Clone() error {
	r.Lock()
	defer r.Unlock()

	// Clone the first tracked branch instead of the default branch
	checkoutBranch := r.repoConfig.Branches[0]

	raw_repo, err := git.Clone(r.repoConfig.Url, r.store, &git.CloneOptions{
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
