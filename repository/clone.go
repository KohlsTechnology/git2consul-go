package repository

import (
	"path"

	"github.com/libgit2/git2go"
)

// Clone the repository
func (r *Repository) Clone() error {

	raw_repo, err := git.Clone(r.repoConfig.Url, r.store, &git.CloneOptions{})
	if err != nil {
		return err
	}

	r.Repository = raw_repo

	// Pulls all remote branches as well
	itr, err := r.NewBranchIterator(git.BranchRemote)
	if err != nil {
		return err
	}

	var checkoutBranchFn = func(b *git.Branch, _ git.BranchType) error {
		bn, err := b.Name()
		if err != nil {
			return err
		}
		_, err = r.References.Lookup("refs/heads/" + path.Base(bn))
		if err != nil {
			_, err = r.References.Create("refs/heads/"+path.Base(bn), b.Reference.Target(), true, "")
			if err != nil {
				return err
			}
		}

		return nil
	}

	itr.ForEach(checkoutBranchFn)

	r.signal <- Signal{
		Type: "clone",
	}

	return nil
}
