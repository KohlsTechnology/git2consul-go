package repository

import (
	"os"
	"testing"

	"github.com/Cimpress-MCP/go-git2consul/testutil"
)

func TestDiffStatus(t *testing.T) {
	_, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfg := testutil.LoadTestConfig(t)

	repos, err := LoadRepos(cfg)
	if err != nil {
		t.Fatal(err)
	}
	repo := repos[0]

	// Push a commit to the repository
	cleanup = testutil.TempCommitTestRepo(t)

	_, err = repo.Pull("master")
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup
	defer func() {
		err = os.RemoveAll(repo.Workdir())
		if err != nil {
			t.Fatal(err)
		}
	}()

	defer cleanup()
}
