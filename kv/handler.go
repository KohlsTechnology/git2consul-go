package kv

import (
	"crypto/tls"
	"net/http"

	"github.com/apex/log"
	"github.com/Cimpress-MCP/go-git2consul/config"
	"github.com/hashicorp/consul/api"
)

type KVHandler struct {
	*api.KV
	logger *log.Entry
}

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
		KV:     kv,
		logger: logger,
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
		consulConfig.HttpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
