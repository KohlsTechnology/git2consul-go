package kv

import (
	"github.com/apex/log"
	"github.com/hashicorp/consul/api"
)

type KVHandler struct {
	*api.KV
	logger log.Interface
}

func New(config *api.Config) (*KVHandler, error) {
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	logger := log.Log
	logger.WithField("caller", "kv-client")

	kv := client.KV()

	handler := &KVHandler{
		KV:     kv,
		logger: logger,
	}

	return handler, nil
}
