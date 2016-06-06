package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Cimpress-MCP/go-git2consul/config"
	"github.com/Cimpress-MCP/go-git2consul/testutil"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"gopkg.in/libgit2/git2go.v24"
)

func init() {
	log.SetHandler(discard.New())
}

func TestLoadRepos(t *testing.T) {
	_, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfg := testutil.LoadTestConfig(t)

	_, err := LoadRepos(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup cloning
	defer func() {
		for _, repo := range cfg.Repos {
			os.RemoveAll(filepath.Join(cfg.LocalStore, repo.Name))
		}
	}()
}

func TestLoadRepos_bareDir(t *testing.T) {
	gitRepo, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfgPath := filepath.Join(testutil.FixturesPath(t), "example.json")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	// Change some values for test runtime test repo
	cfg.Repos[0].Url = gitRepo.Workdir()

	err = os.Mkdir(filepath.Join(cfg.LocalStore, cfg.Repos[0].Name), 0755)
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadRepos(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup cloning
	defer func() {
		for _, repo := range cfg.Repos {
			os.RemoveAll(filepath.Join(cfg.LocalStore, repo.Name))
		}
	}()
}

func TestLoadRepos_invalidRepo(t *testing.T) {
	cfgPath := filepath.Join(testutil.FixturesPath(t), "example.json")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadRepos(cfg)
	if err == nil {
		t.Fatal("Expected failure for invalid repository url")
	}

	// Cleanup cloning
	defer func() {
		for _, repo := range cfg.Repos {
			os.RemoveAll(filepath.Join(cfg.LocalStore, repo.Name))
		}
	}()
}

func TestLoadRepos_existingRepo(t *testing.T) {
	_, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	cfg := testutil.LoadTestConfig(t)

	// Init a repo in the store:project
	err := os.Mkdir(filepath.Join(cfg.LocalStore, cfg.Repos[0].Name), 0755)
	if err != nil {
		t.Fatal(err)
	}
	repo, err := git.InitRepository(filepath.Join(cfg.LocalStore, cfg.Repos[0].Name), false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Remotes.Create("origin", "/foo/bar")
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadRepos(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup cloning
	defer func() {
		for _, repo := range cfg.Repos {
			os.RemoveAll(filepath.Join(cfg.LocalStore, repo.Name))
		}
	}()
}
