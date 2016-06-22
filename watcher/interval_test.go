package watch

import (
	"os"
	"testing"

	"github.com/Cimpress-MCP/go-git2consul/repository"
	"github.com/Cimpress-MCP/go-git2consul/repository/mock"
	"github.com/Cimpress-MCP/go-git2consul/testutil"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
)

func init() {
	log.SetHandler(discard.New())
}

func TestPollBranches(t *testing.T) {
	gitRepo, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	repo := mock.Repository(gitRepo)

	oid, _ := testutil.GitCommitTestRepo(t)
	odb, err := repo.Odb()
	if err != nil {
		t.Fatal(err)
	}

	w := &Watcher{
		Repositories: []*repository.Repository{repo},
		RepoChangeCh: make(chan *repository.Repository, 1),
		ErrCh:        make(chan error),
		RcvDoneCh:    make(chan struct{}, 1),
		SndDoneCh:    make(chan struct{}, 1),
		logger:       log.WithField("caller", "watcher"),
		hookSvr:      nil,
		once:         true,
	}

	err = w.pollBranches(repo)
	if err != nil {
		t.Fatal(err)
	}

	if !odb.Exists(oid) {
		t.Fatal("Commit not present on remote")
	}

	// Cleanup on git2consul cached repo
	defer func() {
		err = os.RemoveAll(repo.Workdir())
		if err != nil {
			t.Fatal(err)
		}
	}()
}
