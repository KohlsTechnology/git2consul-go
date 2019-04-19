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
	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/hashicorp/consul/api"
)

// PutKV triggers an KV api request to put data to the Consul.
func (h *KVHandler) PutKV(repo repository.Repo, prefix string, value []byte) error {
	head, err := repo.Head()
	if err != nil {
		return err
	}

	branchName := head.Name().Short()

	key, status, err := getItemKey(repo, prefix)
	if err != nil {
		if status == SourceRootNotInPrefix {
			h.logger.Infof("%s Skipping!", err)
		}
		if status == PathFormatterError {
			return err
		}
		return nil
	}

	h.logger.Debugf("KV PUT: %s/%s: %s", repo.Name(), branchName, key)

	p := &api.KVPair{
		Key:   key,
		Value: value,
	}

	_, err = h.Put(p, nil)
	if err != nil {
		return err
	}

	return nil
}

//DeleteKV deletes provided item from the KV store.
func (h *KVHandler) DeleteKV(repo repository.Repo, prefix string) error {
	key, status, err := getItemKey(repo, prefix)
	if err != nil {
		if status == SourceRootNotInPrefix {
			h.logger.Infof("%s Skipping!", err)
		}
		if status == PathFormatterError {
			return err
		}
		return nil
	}

	h.logger.Infof("KV DEL %s/%s/%s", repo.Name(), repo.Branch(), key)
	_, err = h.Delete(key, nil)
	if err != nil {
		return err
	}

	return nil
}

//DeleteTreeKV deletes recursively all the keys with given prefix.
func (h *KVHandler) DeleteTreeKV(repo repository.Repo, prefix string) error {
	key, status, err := getItemKey(repo, prefix)
	if err != nil {
		if status == SourceRootNotInPrefix {
			h.logger.Infof("%s Skipping!", err)
		}
		if status == PathFormatterError {
			return err
		}
		return nil
	}

	h.logger.Infof("KV DEL %s/%s/%s", repo.Name(), repo.Branch(), key)
	_, err = h.DeleteTree(key, nil)
	if err != nil {
		return err
	}

	return nil
}
