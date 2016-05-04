package repository

import "github.com/libgit2/git2go"

// Clone the repository
func (r *Repository) Clone() error {
	// Use temp dir for now

	raw_repo, err := git.Clone(r.repoConfig.Url, r.store, &git.CloneOptions{})
	if err != nil {
		return err
	}

	r.Repository = raw_repo
	// ref, _ := raw_repo.References.Lookup("refs/heads/master")
	// log.Debugf("=== References %v", ref)

	r.signal <- Signal{
		Type: "clone",
	}

	return nil
}
