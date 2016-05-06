package repository

import (
	"os"

	"github.com/libgit2/git2go"
)

// Clone the repository
func (r *Repository) Clone() error {
	err := os.Mkdir(r.store, 0755)
	if err != nil {
		return err
	}

	raw_repo, err := git.Clone(r.repoConfig.Url, r.store, &git.CloneOptions{})
	if err != nil {
		return err
	}

	r.Repository = raw_repo

	err = r.checkoutRemoteBranches()
	if err != nil {
		return err
	}

	r.cloneCh <- struct{}{}

	return nil
}
