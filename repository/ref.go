package repository

import "gopkg.in/libgit2/git2go.v24"

// CheckRef checks whether a particular ref is part of the repository
func (r *Repository) CheckRef(ref string) error {
	oid, err := git.NewOid(ref)
	if err != nil {
		return err
	}

	// This can be for a different repo
	_, err = r.Lookup(oid)
	if err != nil {
		return err
	}

	return nil
}
