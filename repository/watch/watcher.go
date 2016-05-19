package watch

import (
	"sync"

	"github.com/cleung2010/go-git2consul/repository"
)

type Watcher struct {
	sync.Mutex

	Repositories []*repository.Repository

	RepoChangeCh chan *repository.Repository
}
