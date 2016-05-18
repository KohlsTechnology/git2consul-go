package repository

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/cleung2010/go-git2consul/config"
	"gopkg.in/libgit2/git2go.v24"
)

func TestLoadRepos(t *testing.T) {
	_, cleanup := tempGitInitPath(t)
	defer cleanup()

	cfg := loadConfig(t)

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

func TestLoadRepos_invalidRepo(t *testing.T) {
	cfgPath := filepath.Join("test-fixtures", "example.json")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadRepos(cfg)
	if err == nil {
		t.Fatal("Expected failure from LoadRepos()")
	}

	// Cleanup cloning
	defer func() {
		for _, repo := range cfg.Repos {
			os.RemoveAll(filepath.Join(cfg.LocalStore, repo.Name))
		}
	}()
}

// Helper functions

// Init repository specified on test-fixtures
func tempGitInitPath(t *testing.T) (*git.Repository, func()) {
	repoPath := filepath.Join("test-fixtures", "example")
	fi, err := os.Stat(repoPath)
	if err != nil {
		t.Fatal(err)
	}
	if fi.IsDir() == false {
		t.Fatal(err)
	}

	// Init repo
	repo, err := git.InitRepository(repoPath, false)
	if err != nil {
		t.Fatal(err)
	}

	// Add files to index
	idx, err := repo.Index()
	if err != nil {
		t.Fatal(err)
	}
	err = idx.AddAll([]string{}, git.IndexAddDefault, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = idx.Write()
	if err != nil {
		t.Fatal(err)
	}

	treeId, err := idx.WriteTree()
	if err != nil {
		t.Fatal(err)
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		t.Fatal(err)
	}

	// Initial commit
	sig := &git.Signature{
		Name:  "Test Example",
		Email: "tes@example.com",
		When:  time.Date(2016, 01, 01, 12, 00, 00, 0, time.UTC),
	}

	repo.CreateCommit("HEAD", sig, sig, "Initial commit", tree)

	// Save commmit ref for reset later
	h, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}

	obj, err := repo.Lookup(h.Target())
	if err != nil {
		t.Fatal(err)
	}

	initialCommit, err := obj.AsCommit()
	if err != nil {
		t.Fatal(err)
	}

	// Reset to initial commit, and then remove .git
	var cleanup = func() {
		repo.ResetToCommit(initialCommit, git.ResetHard, &git.CheckoutOpts{
			Strategy: git.CheckoutForce,
		})

		repo.StateCleanup()

		dotgit := filepath.Join(repo.Path())
		os.RemoveAll(dotgit)
	}

	return repo, cleanup
}

// Make a commit to the repository in test-fixtures, and return
// the change for test verification
func tempCommitRepo(r *git.Repository, t *testing.T) []byte {
	// Make changes
	date := []byte(time.Now().String())
	file := path.Join("test-fixtures", "example", "foo")
	err := ioutil.WriteFile(file, date, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Commit changes
	idx, err := r.Index()
	if err != nil {
		t.Fatal(err)
	}

	treeId, err := idx.WriteTree()

	tree, err := r.LookupTree(treeId)
	if err != nil {
		t.Fatal(err)
	}

	h, err := r.Head()
	if err != nil {
		t.Fatal(err)
	}

	commit, err := r.LookupCommit(h.Target())
	if err != nil {
		log.Fatal(err)
	}

	sig := &git.Signature{
		Name:  "Test Example",
		Email: "tes@example.com",
		When:  time.Date(2016, 01, 01, 12, 00, 00, 0, time.UTC),
	}

	_, err = r.CreateCommit("HEAD", sig, sig, "Update commit", tree, commit)
	if err != nil {
		t.Fatal(err)
	}

	return date
}

func loadConfig(t *testing.T) *config.Config {
	cfgPath := filepath.Join("test-fixtures", "example.json")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}
