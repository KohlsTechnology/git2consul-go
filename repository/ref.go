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
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// CheckRef checks whether a particular ref is part of the repository
func (r *Repository) CheckRef(ref string) error {
	_, err := r.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return err
	}

	return nil
}
