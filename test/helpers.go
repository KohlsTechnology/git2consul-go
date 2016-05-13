package test

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/libgit2/git2go.v24"
)

func TempGitInitPath(path string, t *testing.T) func() {
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("git init err: %s", err)
	}
	if fi.IsDir() == false {
		t.Fatalf("git init err: %s is not a directory", path)
	}

	repo, err := git.InitRepository(path, false)
	if err != nil {
		t.Fatalf("git init err: %s", err)
	}

	h, err := repo.Head()
	if err != nil {
		t.Fatalf("git init err: %s", err)
	}

	obj, err := repo.Lookup(h.Target())
	if err != nil {
		t.Fatalf("git init err: %s", err)
	}

	initialCommit := obj.AsCommit()

	// Reset to initial commit, and then remove .git
	var cleanup = func() {
		repo.ResetToCommit(initialCommit, git.ResetHard, &git.CheckoutOpts{
			Strategy: git.CheckoutForce,
		})

		dotgit := filepath.Join(repo.Path(), ".git")
		os.RemoveAll(dotgit)
	}()

	return repo, cleanup
}
