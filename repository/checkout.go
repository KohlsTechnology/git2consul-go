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
	"path"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

func (r *Repository) checkoutConfigBranches() error {
	err := r.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
		Auth:     r.Authentication,
	})

	w, err := r.Worktree()

	if err != nil {
		return err
	}

	refIter, err := remoteBranches(r.Storer)

	_ = refIter.ForEach(func(b *plumbing.Reference) error {
		branchOnRemote := StringInSlice(path.Base(b.Name().String()), r.Config.Branches)
		if branchOnRemote != false {
			err := w.Checkout(&git.CheckoutOptions{
				Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", b.Name())),
				Force:  true,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return nil
}

//CheckoutBranch performs a checkout on the specific branch
func (r *Repository) CheckoutBranch(branch plumbing.ReferenceName) error {
	err := r.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
		Auth:     r.Authentication,
		Force:    true,
	})

	w, err := r.Worktree()

	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: branch,
		Force:  true,
	})

	if err != nil {
		return err
	}

	return nil
}

func remoteBranches(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote()
	}, refs), nil
}

//LocalBranches returns an iterator to iterate only over local branches.
func LocalBranches(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return !ref.Name().IsRemote()
	}, refs), nil
}

//StringInSlice checks if value exists within slice.
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
