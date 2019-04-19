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

package kv

import (
	"os"
	"path/filepath"

	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/apex/log"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Push a repository branch to the KV
// TODO: Optimize for PUT only on changes instead of the entire repo
func (h *KVHandler) putBranch(repo repository.Repo, branch plumbing.ReferenceName) error {
	// Checkout branch
	repo.CheckoutBranch(branch)

	// h, _ := repo.Head()
	// bn, _ := h.Branch().Name()
	// log.Debugf("(consul) pushBranch(): Branch: %s Head: %s", bn, h.Target().String())
	workdir := repository.WorkDir(repo)
	sourceRoot := repo.GetConfig().SourceRoot
	var pushFile = func(fullpath string, info os.FileInfo, err error) error {
		// Walk error
		if err != nil {
			return err
		}

		// Skip the .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Do not push directories
		if info.IsDir() {
			return nil
		}

		file := Init(fullpath, repo)
		err = file.Create(h, repo)
		if err != nil {
			h.logger.Errorf("%s", err)
		}
		return nil
	}
	workdir = filepath.Join(workdir, sourceRoot)
	err := filepath.Walk(workdir, pushFile)
	if err != nil {
		log.WithError(err).Debug("PUT branch error")
		return err
	}

	return nil
}
