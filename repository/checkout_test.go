package repository

import (
	"os"
	"testing"

	"github.com/Cimpress-MCP/go-git2consul/testutil"
	"gopkg.in/libgit2/git2go.v24"
)

func TestCheckoutBranch(t *testing.T) {
	_, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfg := testutil.LoadTestConfig(t)

	repos, err := LoadRepos(cfg)
	if err != nil {
		t.Fatal(err)
	}
	repo := repos[0]

	branch, err := repo.LookupBranch("master", git.BranchLocal)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.CheckoutBranch(branch, &git.CheckoutOpts{})
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
}
