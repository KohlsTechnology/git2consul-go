package consul

import (
	"github.com/cleung2010/go-git2consul/config"
	"github.com/hashicorp/consul/api"
)

type Client struct {
	*api.Client
}

func NewClient(cfg *config.Config) (*Client, error) {
	// TODO: Use git2consul defaults for the client
	raw_client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	client := &Client{
		raw_client,
	}

	return client, nil
}
