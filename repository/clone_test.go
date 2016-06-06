package repository

import (
	"os"
	"path"
	"testing"

	"github.com/Cimpress-MCP/go-git2consul/testutil"
)

func TestClone(t *testing.T) {
	gitRepo, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfg := testutil.LoadTestConfig(t)

	r := &Repository{
		Repository: gitRepo,
		repoConfig: cfg.Repos[0],
		store:      path.Join(cfg.LocalStore, cfg.Repos[0].Name),
	}

	err := r.Clone()
	if err != nil {
		t.Fatal(err)
	}

	//Cleanup cloned repo
	defer func() {
		err = os.RemoveAll(r.store)
		if err != nil {
			t.Fatal(err)
		}
	}()
}
