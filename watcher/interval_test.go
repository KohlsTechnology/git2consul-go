package watch

import (
	"os"
	"testing"

	"github.com/Cimpress-MCP/go-git2consul/repository"
	"github.com/Cimpress-MCP/go-git2consul/testutil"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
)

func init() {
	log.SetHandler(discard.New())
}

func TestPollBranches(t *testing.T) {
	_, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfg := testutil.LoadTestConfig(t)

	repos, err := repository.LoadRepos(cfg)
	if err != nil {
		t.Fatal(err)
	}
	repo := repos[0]

	oid, cleanup := testutil.TempCommitTestRepo(t)
	odb, err := repo.Odb()
	if err != nil {
		t.Fatal(err)
	}

	w := &Watcher{
		Repositories: repos,
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

	// err = repo.PollBranches()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// // Verify that the file changed
	// filePath := filepath.Join("test-fixtures", "example", "foo")
	// actual, err := ioutil.ReadFile(filePath)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// if !reflect.DeepEqual(expected, actual) {
	// 	t.Fatal("Polling failed to pull files")
	// }

	// Cleanup
	defer func() {
		err = os.RemoveAll(repo.Workdir())
		if err != nil {
			t.Fatal(err)
		}
	}()

	defer cleanup()
}
