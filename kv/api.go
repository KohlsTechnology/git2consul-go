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

//Handler interface for Key-Value store.
type Handler interface {
	PutKV(repository.Repo, string, []byte) error
	DeleteKV(repository.Repo, string) error
	DeleteTreeKV(repository.Repo, string) error
	HandleUpdate(repository.Repo) error
}

//API minimal Consul KV api implementation
type API interface {
	Get(string, *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error)
	Put(*api.KVPair, *api.WriteOptions) (*api.WriteMeta, error)
	Txn(api.KVTxnOps, *api.QueryOptions) (bool, *api.KVTxnResponse, *api.QueryMeta, error)
}
