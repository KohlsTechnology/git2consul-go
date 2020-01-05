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

package mocks

import (
	"testing"

	"github.com/KohlsTechnology/git2consul-go/config"
	"gopkg.in/src-d/go-billy.v4/osfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage"
)

// Repo TODO write a useful documentation here
type Repo struct {
	adds   []string
	Config *config.Repo
	Path   string
	branch plumbing.ReferenceName
	T      *testing.T
	hashes map[string]plumbing.Hash
}

// Name TODO write a useful documentation here
func (r *Repo) Name() string {
	return "repository_mock"
}

// GetConfig TODO write a useful documentation here
func (r *Repo) GetConfig() *config.Repo {
	return r.Config
}

// Add TODO write a useful documentation here
func (r *Repo) Add(path string) {
	r.adds = append(r.adds, path)
}

// CheckRef TODO write a useful documentation here
func (r *Repo) CheckRef(branch string) error {
	return nil
}

// CheckoutBranch TODO write a useful documentation here
func (r *Repo) CheckoutBranch(branch plumbing.ReferenceName) error {
	r.branch = branch
	return nil
}

// DiffStatus TODO write a useful documentation here
func (r *Repo) DiffStatus(commit string) (object.Changes, error) {
	var changes object.Changes
	for _, add := range r.adds {
		changes = append(changes, &object.Change{From: object.ChangeEntry{}, To: object.ChangeEntry{Name: add}})
	}
	r.adds = []string{}
	return changes, nil
}

// Head TODO write a useful documentation here
func (r *Repo) Head() (*plumbing.Reference, error) {
	if r.branch == "" {
		r.branch = plumbing.NewReferenceFromStrings("master", "").Name()
		r.Pull("master")
	}
	return plumbing.NewHashReference(r.branch, r.hashes[r.branch.Short()]), nil
}

// Pull TODO write a useful documentation here
func (r *Repo) Pull(branch string) error {
	if r.hashes == nil {
		r.hashes = make(map[string]plumbing.Hash)
	}
	if r.hashes[branch] == plumbing.ZeroHash {
		r.hashes[branch] = plumbing.ComputeHash(0, []byte(branch))
	} else {
		hash := r.hashes[branch]
		r.hashes[branch] = plumbing.ComputeHash(0, hash[:])
	}
	r.branch = plumbing.NewReferenceFromStrings(branch, "").Name()
	return nil
}

// ResolveRevision TODO write a useful documentation here
func (r *Repo) ResolveRevision(plumbing.Revision) (*plumbing.Hash, error) {
	hash := r.hashes[r.branch.Short()]
	return &hash, nil
}

// Worktree TODO write a useful documentation here
func (r *Repo) Worktree() (*git.Worktree, error) {
	return &git.Worktree{Filesystem: osfs.New(r.Path)}, nil
}

// Lock TODO write a useful documentation here
func (r *Repo) Lock() {}

// Unlock TODO write a useful documentation here
func (r *Repo) Unlock() {}

// GetStorer TODO write a useful documentation here
func (r *Repo) GetStorer() storage.Storer {
	return nil
}

// Branch TODO write a useful documentation here
func (r *Repo) Branch() plumbing.ReferenceName {
	return r.branch
}
