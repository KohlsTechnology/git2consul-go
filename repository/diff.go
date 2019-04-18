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
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// DiffStatus compares the current workdir with a target ref and return the modified files
func (r *Repository) DiffStatus(ref string) (object.Changes, error) {
	sourceRoot := strings.TrimPrefix(r.GetConfig().SourceRoot, "/")
	h0, err := r.Head()
	if err != nil {
		return nil, err
	}
	c0, err := r.CommitObject(h0.Hash())
	if err != nil {
		return nil, err
	}
	c1, err := r.CommitObject(plumbing.NewHash(ref))
	if err != nil {
		return nil, err
	}

	commits := []*object.Commit{c0, c1}
	if len(commits[0].ParentHashes) != 0 {
		commits = []*object.Commit{c1, c0}
	}

	t0, err := r.TreeObject(commits[0].TreeHash)
	if err != nil {
		return nil, err
	}
	t1, err := r.TreeObject(commits[1].TreeHash)
	if err != nil {
		return nil, err
	}
	diff, err := t0.Diff(t1)
	if err != nil {
		return nil, err
	}
	return applySourceRoot(diff, sourceRoot), nil
}

func applySourceRoot(changes object.Changes, sourceRoot string) object.Changes {
	var selected object.Changes
	empty := object.ChangeEntry{}
	for _, change := range changes {
		name := ""
		if change.From != empty {
			name = change.From.Name
		} else {
			name = change.To.Name
		}
		if strings.HasPrefix(name, sourceRoot) {
			selected = append(selected, change)
		}
	}
	return selected
}
