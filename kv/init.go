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
	"path/filepath"
	"strings"

	"github.com/KohlsTechnology/git2consul-go/repository"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
)

// HandleInit handles initial fetching of the KV on start
func (h *KVHandler) HandleInit(repos []repository.Repo) error {
	for _, repo := range repos {
		err := h.handleRepoInit(repo)
		if err != nil {
			return err
		}
	}

	return nil
}

// Handles differences on all branches of a repository, comparing the ref
// of the branch against the one in the KV
func (h *KVHandler) handleRepoInit(repo repository.Repo) error {
	repo.Lock()
	defer repo.Unlock()

	storer := repo.GetStorer()
	itr, err := storer.IterReferences()
	if err != nil {
		return err
	}

	for {
		ref, err := itr.Next()
		if err != nil {
			break
		}
		if strings.Contains(ref.String(), "HEAD") {
			continue
		}

		if ref.Name().IsRemote() == false {
			h.logger.Infof("KV GET ref: %s/%s", repo.Name(), ref.Name())
			kvRef, err := h.getKVRef(repo, ref.Name().String())

			if err != nil {
				return err
			}

			localRef := ref.Hash().String()

			if len(kvRef) == 0 {
				// There is no ref in the KV, push the entire branch
				h.logger.Infof("KV PUT changes: %s/%s", repo.Name(), ref.Name())
				h.putBranch(repo, plumbing.ReferenceName(ref.Name().Short()))

				h.logger.Infof("KV PUT ref: %s/%s", repo.Name(), ref.Name())
				h.putKVRef(repo, ref.Name().String())
			} else if kvRef != localRef {
				//Check if the ref belongs to that repo
				err := repo.CheckRef(kvRef)
				if err != nil {
					return err
				}

				// Handle modified and deleted files
				deltas, err := repo.DiffStatus(kvRef)
				if err != nil {
					return err
				}
				h.handleDeltas(repo, deltas)

				err = h.putKVRef(repo, ref.Name().String())
				if err != nil {
					return err
				}
				h.logger.Debugf("KV PUT ref: %s/%s", repo.Name(), ref.Name())
			}
		}
	}
	return nil
}

// Helper function that handles deltas
func (h *KVHandler) handleDeltas(repo repository.Repo, diff object.Changes) error {
	for _, d := range diff {
		action, err := d.Action()
		if err != nil {
			return err
		}
		workDir := repository.WorkDir(repo)
		switch action {
		case merkletrie.Insert:
			filePath := filepath.Join(workDir, d.To.Name)
			h.logger.Debugf("Detected added file: %s", filePath)
			file := Init(filePath, repo)
			err := file.Create(h, repo)
			if err != nil {
				return err
			}
		case merkletrie.Modify:
			filePath := filepath.Join(workDir, d.To.Name)
			h.logger.Debugf("Detected modified file: %s", filePath)
			file := Init(filePath, repo)
			err := file.Update(h, repo)
			if err != nil {
				return err
			}
		case merkletrie.Delete:
			filePath := filepath.Join(workDir, d.From.Name)
			h.logger.Debugf("Detected deleted file: %s", filePath)
			file := Init(filePath, repo)
			err := file.Delete(h, repo)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
