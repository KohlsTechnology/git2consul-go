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

	"github.com/KohlsTechnology/git2consul-go/config"
	"github.com/apex/log"
	"github.com/hashicorp/consul/api"
)

const consulTxnSize = 64

// KVHandler is used to manipulate the KV
type KVHandler struct { //nolint:revive
	API
	api.KVTxnOps
	logger *log.Entry
}

//TransactionIntegrityError implements error to handle any violation of transaction atomicity.
type TransactionIntegrityError struct {
	msg string
}

func (e *TransactionIntegrityError) Error() string { return e.msg }

// New creates new KV handler to manipulate the Consul VK
func New(config *config.ConsulConfig) (*KVHandler, error) {
	client, err := newAPIClient(config)
	if err != nil {
		return nil, err
	}

	logger := log.WithFields(log.Fields{
		"caller": "consul",
	})

	kv := client.KV()

	handler := &KVHandler{
		API:      kv,
		KVTxnOps: nil,
		logger:   logger,
	}

	return handler, nil
}

func newAPIClient(config *config.ConsulConfig) (*api.Client, error) {
	consulConfig := api.DefaultConfig()

	if config.Address != "" {
		consulConfig.Address = config.Address
	}

	if config.Token != "" {
		consulConfig.Token = config.Token
	}

	if config.SSLEnable {
		consulConfig.Scheme = "https"
	}

	if !config.SSLVerify {
		consulConfig.TLSConfig.InsecureSkipVerify = true
	}

	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

//Put overrides Consul API Put function to add entry to KVTxnOps.
func (h *KVHandler) Put(kvPair *api.KVPair, wOptions *api.WriteOptions) (*api.WriteMeta, error) {
	txnItem := &api.KVTxnOp{
		Verb:  api.KVSet,
		Key:   kvPair.Key,
		Value: kvPair.Value,
	}
	h.KVTxnOps = append(h.KVTxnOps, txnItem)
	return nil, nil
}

//Delete overrides Consul API Delete function to add entry to KVTxnOps.
func (h *KVHandler) Delete(key string, wOptions *api.WriteOptions) (*api.WriteMeta, error) {
	txnItem := &api.KVTxnOp{
		Verb: api.KVDelete,
		Key:  key,
	}
	h.KVTxnOps = append(h.KVTxnOps, txnItem)
	return nil, nil
}

//DeleteTree overrides Consul API DeleteTree function to add entry to KVTxnOps.
func (h *KVHandler) DeleteTree(key string, wOptions *api.WriteOptions) (*api.WriteMeta, error) {
	txnItem := &api.KVTxnOp{
		Verb: api.KVDeleteTree,
		Key:  key,
	}
	h.KVTxnOps = append(h.KVTxnOps, txnItem)
	return nil, nil
}

//Commit function executes set of operations from KVTxnOps as single transaction.
func (h *KVHandler) Commit() error {
	defer func() {
		h.KVTxnOps = nil
	}()
	var kvTxnOps = h.KVTxnOps
	//move modify index check to the end
	if h.KVTxnOps[0].Verb == api.KVCheckIndex {
		length := len(h.KVTxnOps)
		kvTxnOps = append(h.KVTxnOps[1:length-1], h.KVTxnOps[0], h.KVTxnOps[length-1])
	}
	for _, slice := range h.splitIntoSlices(kvTxnOps, consulTxnSize) {
		err := h.executeTransaction(slice)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *KVHandler) executeTransaction(KVTxnOps api.KVTxnOps) error {
	status, response, _, err := h.Txn(KVTxnOps, nil)
	if err != nil {
		return err
	}
	h.logger.Debugf("Transaction with %d items was sent to the KV store", len(KVTxnOps))
	if !status {
		errMsg := ""
		for _, txError := range response.Errors {
			errMsg += fmt.Sprintf("%s\n", txError.What)
		}
		return &TransactionIntegrityError{fmt.Sprintf("Transaction has been rolled back due to: %s", errMsg)}
	}
	return nil
}

func (h *KVHandler) splitIntoSlices(kvTxnOps api.KVTxnOps, sliceLength int) []api.KVTxnOps {
	var kvTxnSlices []api.KVTxnOps
	for len(kvTxnOps) > 0 {
		index := sliceLength
		if index > len(kvTxnOps) {
			index = len(kvTxnOps)
		}
		var slice api.KVTxnOps
		slice = append(slice, kvTxnOps[:index]...)
		kvTxnOps = kvTxnOps[index:]
		kvTxnSlices = append(kvTxnSlices, slice)
	}
	return kvTxnSlices
}
