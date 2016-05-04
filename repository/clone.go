package repository

import (
	log "github.com/Sirupsen/logrus"
	"github.com/libgit2/git2go"
)

// Clone the repository
func (r *Repository) Clone() error {
	// Use temp dir for now

	raw_repo, err := git.Clone(r.repoConfig.Url, r.store, &git.CloneOptions{})
	if err != nil {
		return err
	}

	r.Repository = raw_repo
	// TODO: Fix consul push on a single branch
	ref, _ := raw_repo.References.Lookup("refs/heads/test")
	log.Debugf("=== References %v", ref)
	r.UpdateCh <- true

	return nil
}
