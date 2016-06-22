package mock

import (
	"fmt"
	"io/ioutil"

	"github.com/Cimpress-MCP/go-git2consul/config/mock"
	"github.com/Cimpress-MCP/go-git2consul/repository"

	"gopkg.in/libgit2/git2go.v24"
)

// Repository returns a mock of a repository.Repository object
func Repository(gitRepo *git.Repository) *repository.Repository {
	if gitRepo == nil {
		return nil
	}

	repoConfig := mock.RepoConfig(gitRepo.Workdir())

	dstPath, err := ioutil.TempDir("", "git2consul-test-local")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	localRepo, err := git.Clone(repoConfig.Url, dstPath, &git.CloneOptions{})
	if err != nil {
		fmt.Print(err)
		return nil
	}

	repo := &repository.Repository{
		Repository: localRepo,
		Config:     repoConfig,
	}

	return repo
}
