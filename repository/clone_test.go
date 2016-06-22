package repository

import (
	"os"
	"path"
	"testing"

	"github.com/Cimpress-MCP/go-git2consul/config/mock"
	"github.com/Cimpress-MCP/go-git2consul/testutil"
)

func TestClone(t *testing.T) {
	gitRepo, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfg := mock.Config(gitRepo.Workdir())

	repo := &Repository{
		Config: cfg.Repos[0],
	}

	repoPath := path.Join(cfg.LocalStore, repo.Config.Name)
	err := repo.Clone(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	//Cleanup cloned repo
	defer func() {
		err = os.RemoveAll(repo.Workdir())
		if err != nil {
			t.Fatal(err)
		}
	}()
}
