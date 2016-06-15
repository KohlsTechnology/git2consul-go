package mock

import (
	"path/filepath"
	"runtime"

	"github.com/Cimpress-MCP/go-git2consul/repository"
	"gopkg.in/libgit2/git2go.v24"
)

func fixturesPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {

	}

	testutilPath := filepath.Dir(filename)
	return filepath.Join(testutilPath, "test-fixtures")
}

func initGitRepository() *git.Repository {
	repo := &git.Repository{}

	return repo
}

// Returns a mock of a repository.Repository object
func Repository() *repository.Repository {
	repo := &repository.Repository{}

	return repo
}
