package repository

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Cimpress-MCP/go-git2consul/config/mock"
	"github.com/Cimpress-MCP/go-git2consul/testutil"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"gopkg.in/libgit2/git2go.v24"
)

func init() {
	log.SetHandler(discard.New())
}

func TestLoadRepos(t *testing.T) {
	gitRepo, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfg := mock.Config(gitRepo.Workdir())

	repos, err := LoadRepos(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup cloning
	defer func() {
		for _, repo := range repos {
			os.RemoveAll(repo.Workdir())
		}
	}()
}

func TestLoadRepos_existingDir(t *testing.T) {
	bareDir, err := ioutil.TempDir("", "bare-dir")
	if err != nil {
		t.Fatal(err)
	}

	cfg := mock.Config(bareDir)

	_, err = LoadRepos(cfg)
	if err == nil {
		t.Fatal("Expected failure for existing repository")
	}

	// Cleanup
	defer func() {
		os.RemoveAll(bareDir)
	}()
}

func TestLoadRepos_invalidRepo(t *testing.T) {
	cfg := mock.Config("bogus-url")

	_, err := LoadRepos(cfg)
	if err == nil {
		t.Fatal("Expected failure for invalid repository url")
	}
}

func TestLoadRepos_existingRepo(t *testing.T) {
	gitRepo, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfg := mock.Config(gitRepo.Workdir())
	localRepoPath := filepath.Join(cfg.LocalStore, cfg.Repos[0].Name)

	// Init a repo in the local store, with same name are the "remote"
	err := os.Mkdir(localRepoPath, 0755)
	if err != nil {
		t.Fatal(err)
	}
	// defer os.RemoveAll(localRepoPath)

	repo, err := git.InitRepository(localRepoPath, false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Remotes.Create("origin", "/foo/bar")
	if err != nil {
		t.Fatal(err)
	}

	repos, err := LoadRepos(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup cloning
	defer func() {
		for _, repo := range repos {
			os.RemoveAll(repo.Workdir())
		}
	}()
}
