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

	"gopkg.in/src-d/go-git.v4"
)

// Clone the repository. Cloning will only checkout tracked branches.
// A destination path to clone to needs to be provided
func (r *Repository) Clone(path string) error {
	r.Lock()
	defer r.Unlock()

	if len(r.Config.Branches) == 0 {
		return fmt.Errorf("No tracked branches specified")
	}

	rawRepo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:  r.Config.URL,
		Auth: r.Authentication,
	})

	if err != nil {
		return err
	}

	r.Repository = rawRepo

	err = r.checkoutConfigBranches()
	if err != nil {
		return err
	}

	return nil
}
