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
	"math/rand"
	"testing"

	"github.com/hashicorp/consul/api"
)

type item struct {
	value       []byte
	modifyindex uint64
}

// KV TODO write a useful documentation here
type KV struct {
	T     *testing.T
	items map[string]*item
}

// Get TODO write a useful documentation here
func (kv *KV) Get(key string, opts *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error) {
	kv.T.Logf("KV Get %s", key)
	if val, ok := kv.items[key]; ok {
		return &api.KVPair{Key: key, Value: val.value, ModifyIndex: val.modifyindex}, nil, nil
	}
	return nil, nil, nil
}

// Put TODO write a useful documentation here
func (kv *KV) Put(kvPair *api.KVPair, wOptions *api.WriteOptions) (*api.WriteMeta, error) {
	if kv.items == nil {
		kv.items = make(map[string]*item)
	}
	kv.T.Logf("KV Put %s", kvPair.Key)
	kv.items[kvPair.Key] = &item{value: kvPair.Value, modifyindex: rand.Uint64()}
	return nil, nil
}

// Delete TODO write a useful documentation here
func (kv *KV) Delete(key string, wOptions *api.WriteOptions) (*api.WriteMeta, error) {
	delete(kv.items, key)
	return nil, nil
}

// Txn TODO write a useful documentation here
func (kv *KV) Txn(txnops api.KVTxnOps, opts *api.QueryOptions) (bool, *api.KVTxnResponse, *api.QueryMeta, error) {
	var checkIndexItem *api.KVTxnOp
	if length := len(txnops); length > 1 && txnops[length-2].Verb == api.KVCheckIndex {
		checkIndexItem = txnops[length-2]
	}
	for _, item := range txnops {
		switch item.Verb {
		case api.KVSet:
			if checkIndexItem != nil && item.Key == checkIndexItem.Key {
				kvPair, _, _ := kv.Get(item.Key, nil)
				if kvPair.ModifyIndex != checkIndexItem.Index {
					return false, &api.KVTxnResponse{}, nil, nil
				}
			}
			kv.Put(&api.KVPair{Key: item.Key, Value: item.Value}, nil)
		case api.KVDelete:
			kv.Delete(item.Key, nil)
		}
	}
	return true, nil, nil, nil
}
