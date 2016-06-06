package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Cimpress-MCP/go-git2consul/testutil"
)

func TestNew(t *testing.T) {
	_, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfg := testutil.LoadTestConfig(t)
	repoConfig := cfg.Repos[0]
	repoPath := filepath.Join(cfg.LocalStore, repoConfig.Name)

	_, status, err := New(repoPath, repoConfig)
	if err != nil {
		t.Fatal(err)
	}

	if status != RepositoryCloned {
		t.Fatalf("Expected clone status")
	}

	// Call New() again, this time expecting RepositoryOpened
	_, status, err = New(repoPath, repoConfig)
	if err != nil {
		t.Fatal(err)
	}

	if status != RepositoryOpened {
		t.Fatalf("Expected clone status")
	}

	// Cleanup cloning
	defer func() {
		for _, repo := range cfg.Repos {
			os.RemoveAll(filepath.Join(cfg.LocalStore, repo.Name))
		}
	}()
}
