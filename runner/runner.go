package runner

import (
	"fmt"

	"github.com/cleung2010/go-git2consul/config"
	"github.com/cleung2010/go-git2consul/kv"
	"github.com/cleung2010/go-git2consul/repository"
	"github.com/cleung2010/go-git2consul/repository/watch"
	"github.com/hashicorp/consul/api"
)

type Runner struct {
	ErrCh  chan error
	DoneCh chan struct{}

	once bool

	client *api.Client

	repos []*repository.Repository
}

func NewRunner(config *config.Config, once bool) (*Runner, error) {
	// TODO: Use git2consul configs for the client
	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	// Create repos from configuration
	rs, err := repository.LoadRepos(config)
	if err != nil {
		return nil, fmt.Errorf("Cannot load repositories from configuration: %s", err)
	}

	runner := &Runner{
		ErrCh:  make(chan error),
		DoneCh: make(chan struct{}, 1),
		once:   once,
		client: consulClient,
		repos:  rs,
	}

	return runner, nil
}

func (r *Runner) Start() {
	// Watch for local changes to push to KV
	// r.watchKVUpdate()

	// Watch for remote changes to pull locally
	// r.watchReposUpdate()

	handler, err := kv.New(api.DefaultConfig())
	if err != nil {
		fmt.Println("Runner start error")
	}

	handler.HandleInit(r.repos)

	fmt.Println("Here 1")

	rw := watch.New(r.repos)
	rw.Watch()

	fmt.Println("Here 2")

	// Grab changed repos
	for {
		select {
		case err := <-rw.ErrCh:
			fmt.Println("Here ERROR")
			fmt.Println(err)
		case repo := <-rw.RepoChangeCh:
			// Handle change
			handler.HandleUpdate(repo)
			fmt.Println("Here 3")
		}
	}

	// FIXME: This doesn't work atm. Probably needs donCh on watches to block
	// until underlying goroutines are done before we can report back to r.DoneCh
	if r.once {
		r.DoneCh <- struct{}{}
	}
}
