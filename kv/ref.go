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
	"path"

	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/hashicorp/consul/api"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Get local branch ref from the KV
func (h *KVHandler) getKVRef(repo repository.Repo, branchName string) (string, error) {
	refFile := fmt.Sprintf("%s.ref", branchName)
	key := path.Join(repo.Name(), refFile)

	pair, _, err := h.Get(key, nil)
	if err != nil {
		return "", err
	}

	// If error on get, return empty value
	if pair == nil {
		return "", nil
	}
	//store the last modify index
	txnItem := &api.KVTxnOp{
		Verb:  api.KVCheckIndex,
		Index: pair.ModifyIndex,
		Key:   key,
	}
	h.KVTxnOps = append(h.KVTxnOps, txnItem)

	return string(pair.Value), nil
}

// Put the local branch ref to the KV
func (h *KVHandler) putKVRef(repo repository.Repo, branchName string) error {
	refFile := fmt.Sprintf("%s.ref", branchName)
	key := path.Join(repo.Name(), refFile)

	rawRef, err := repo.ResolveRevision(plumbing.Revision("refs/heads/" + branchName))
	if err != nil {
		return err
	}

	p := &api.KVPair{
		Key:   key,
		Value: []byte(rawRef.String()),
	}

	_, err = h.Put(p, nil)
	if err != nil {
		return err
	}
	err = h.Commit()
	if err != nil {
		return err
	}

	return nil
}
