package runner

import (
	"github.com/cleung2010/go-git2consul/config"
	"github.com/cleung2010/go-git2consul/repository"
	"github.com/hashicorp/consul/api"
)

type Runner struct {
	ErrCh chan error

	client *api.Client

	repos repository.Repositories
}

func NewRunner(config *config.Config) (*Runner, error) {
	// TODO: Use git2consul configs for the client
	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	// Create repos from configuration
	rs, err := repository.LoadRepos(config)
	if err != nil {
		return nil, err
	}

	runner := &Runner{
		ErrCh:  make(chan error),
		client: consulClient,
		repos:  rs,
	}

	return runner, nil
}

func (r *Runner) Start() {
	// Watch for local changes to push to KV
	go r.watchKVUpdate()

	// Watch for remote changes to pull locally
	go r.watchReposUpdate()
}
