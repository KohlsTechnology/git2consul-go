package repository

import (
	"github.com/cleung2010/go-git2consul/config"
	"github.com/libgit2/git2go"
)

// Clone the repository
func Clone(cr *config.Repo) (*Repository, error) {
	// Use temp dir for now

	raw_repo, err := git.Clone(cr.Url, defaultStore(cr.Name), &git.CloneOptions{})
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		raw_repo,
		cr,
	}

	return repo, nil
}
