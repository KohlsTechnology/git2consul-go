package repository

import (
	"path/filepath"
	"testing"

	"github.com/cleung2010/go-git2consul/config"
	"github.com/cleung2010/go-git2consul/test"
)

func TestLoadRepos(t *testing.T) {
	repoPath := filepath.Join("test-fixtures", "example")
	cleanup := test.TempGitInitPath(repoPath)
	defer cleanup()

	cfg := config.Load("../config/test-fixtures/local.json")
	repos, err := LoadRepo(cfg)
}
