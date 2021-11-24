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
	"fmt"

	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/apex/log"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// HandleUpdate handles the update of a particular repository.
func (h *KVHandler) HandleUpdate(repo repository.Repo) error {
	w, err := repo.Worktree()
	config := repo.GetConfig()
	repo.Lock()
	defer repo.Unlock()

	if err != nil {
		return err
	}

	for _, branch := range config.Branches {
		err := w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
			Force:  true,
		})
		if err != nil {
			return err
		}
		err = h.UpdateToHead(repo)
		if err != nil {
			return err
		}
	}
	return nil
}

//UpdateToHead handles update to current HEAD comparing diffs against the KV.
func (h *KVHandler) UpdateToHead(repo repository.Repo) error {
	head, err := repo.Head()
	if err != nil {
		return err
	}
	refName := head.Name().Short()
	if err != nil {
		return err
	}

	h.logger.Infof("KV GET ref: %s/%s", repo.Name(), refName)
	kvRef, err := h.getKVRef(repo, refName)
	if err != nil {
		return err
	}

	// Local ref
	refHash := head.Hash().String()
	// log.Debugf("(consul) kvRef: %s | localRef: %s", kvRef, localRef)

	if len(kvRef) == 0 {
		log.Infof("KV PUT changes: %s/%s", repo.Name(), refName)
		err := h.putBranch(repo, plumbing.ReferenceName(head.Name().Short()))
		if err != nil {
			return err
		}

		err = h.putKVRef(repo, refName)
		if err != nil {
			return err
		}
		h.logger.Infof("KV PUT ref: %s/%s", repo.Name(), refName)
	} else if kvRef != refHash {
		// Check if the ref belongs to that repo
		err := repo.CheckRef(refName)
		if err != nil {
			return err
		}

		// Handle modified and deleted files
		deltas, err := repo.DiffStatus(kvRef)
		if err != nil {
			return err
		}
		h.handleDeltas(repo, deltas) //nolint:errcheck

		err = h.putKVRef(repo, refName)
		if err != nil {
			return err
		}
		h.logger.Infof("KV PUT ref: %s/%s", repo.Name(), refName)
	}

	return nil
}
