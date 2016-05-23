package kv

import (
	"github.com/apex/log"
	"github.com/hashicorp/consul/api"
)

type KVHandler struct {
	*api.KV
	logger *log.Entry
}

func New(config *api.Config) (*KVHandler, error) {
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	logger := log.WithFields(log.Fields{
		"caller": "consul",
	})

	kv := client.KV()

	handler := &KVHandler{
		KV:     kv,
		logger: logger,
	}

	return handler, nil
}
