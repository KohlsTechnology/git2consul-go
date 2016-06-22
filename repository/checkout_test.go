package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Cimpress-MCP/go-git2consul/config/mock"
	"github.com/Cimpress-MCP/go-git2consul/testutil"
	"gopkg.in/libgit2/git2go.v24"
)

func TestCheckoutBranch(t *testing.T) {
	gitRepo, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	repoConfig := mock.RepoConfig(gitRepo.Workdir())
	dstPath := filepath.Join(os.TempDir(), repoConfig.Name)

	localRepo, err := git.Clone(repoConfig.Url, dstPath, &git.CloneOptions{})
	if err != nil {
		t.Fatal(err)
	}

	repo := &Repository{
		Repository: localRepo,
		Config:     repoConfig,
	}

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
