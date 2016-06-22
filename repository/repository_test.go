package repository

import (
	"os"
	"testing"

	"github.com/Cimpress-MCP/go-git2consul/config/mock"
	"github.com/Cimpress-MCP/go-git2consul/testutil"
)

func TestNew(t *testing.T) {
	gitRepo, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfg := mock.Config(gitRepo.Workdir())
	repoConfig := cfg.Repos[0]

	repo, status, err := New(cfg.LocalStore, repoConfig)
	if err != nil {
		t.Fatal(err)
	}

	if status != RepositoryCloned {
		t.Fatalf("Expected clone status")
	}

	// Call New() again, this time expecting RepositoryOpened
	repo, status, err = New(cfg.LocalStore, repoConfig)
	if err != nil {
		t.Fatal(err)
	}

	if status != RepositoryOpened {
		t.Fatalf("Expected clone status")
	}

	// Cleanup cloning
	defer func() {
		err := os.RemoveAll(repo.Workdir())
		if err != nil {
			t.Fatal(err)
		}
	}()
}
