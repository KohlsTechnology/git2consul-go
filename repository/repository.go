/*
Copyright 2019 Kohl's Department Stores, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage"

	"github.com/KohlsTechnology/git2consul-go/config"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

const (
	refHead = "refs/heads"
)

// Repo interface represents Repository
type Repo interface {
	Name() string
	Pull(string) error
	CheckoutBranch(plumbing.ReferenceName) error
	CheckRef(string) error
	Head() (*plumbing.Reference, error)
	Lock()
	Unlock()
	DiffStatus(string) (object.Changes, error)
	Worktree() (*git.Worktree, error)
	Branch() plumbing.ReferenceName
	GetConfig() *config.Repo
	GetStorer() storage.Storer
	ResolveRevision(plumbing.Revision) (*plumbing.Hash, error)
}

// Repository is used to hold the git repository object and it's configuration
type Repository struct {
	sync.Mutex

	*git.Repository
	Config         *config.Repo
	Authentication transport.AuthMethod
}

// Status codes for Repository object creation
const (
	RepositoryError = iota // Unused, it will always get returned with an err
	RepositoryCloned
	RepositoryOpened
)

// GetConfig returns config *Repo
func (r *Repository) GetConfig() *config.Repo {
	return r.Config
}

// GetStorer returns Storer
func (r *Repository) GetStorer() storage.Storer {
	return r.Storer
}

// Name returns the repository name
func (r *Repository) Name() string {
	return r.Config.Name
}

// Branch returns the branch name
func (r *Repository) Branch() plumbing.ReferenceName {
	head, err := r.Head()
	if err != nil {
		return ""
	}
	bn := head.Name()
	if err != nil {
		return ""
	}

	return bn
}

// New is used to construct a new repository object from the configuration
func New(repoBasePath string, repoConfig *config.Repo, auth transport.AuthMethod) (*Repository, int, error) {
	repoPath := filepath.Join(repoBasePath, repoConfig.Name)

	r := &Repository{
		Repository:     &git.Repository{},
		Config:         repoConfig,
		Authentication: auth,
	}

	state, err := r.init(repoPath)
	if err != nil {
		return nil, RepositoryError, err
	}

	if r.Repository == nil {
		return nil, RepositoryError, fmt.Errorf("Could not find git repostory")
	}

	return r, state, nil
}

// Initialize git.Repository object by opening the repostiry or cloning from
// the source URL. It does not handle purging existing file or directory
// with the same path
func (r *Repository) init(repoPath string) (int, error) {
	gitRepo, err := git.PlainOpen(repoPath)
	if err != nil || gitRepo == nil {
		err := r.Clone(repoPath)
		if err != nil {
			// more explicit error handling as a workaround for the upstream issue, tracked under:
			// https://github.com/src-d/go-git/issues/741
			switch err {
			case transport.ErrAuthenticationRequired:
				os.RemoveAll(repoPath)
				return RepositoryError, err
			case transport.ErrAuthorizationFailed:
				os.RemoveAll(repoPath)
				return RepositoryError, err
			default:
				os.Remove(repoPath)
				return RepositoryError, err
			}
		}
		return RepositoryCloned, nil
	}

	r.Repository = gitRepo
	return RepositoryOpened, nil
}

//WorkDir returns working directory for a local copy of the repository.
func WorkDir(r Repo) string {
	w, _ := r.Worktree()
	return w.Filesystem.Root()
}
