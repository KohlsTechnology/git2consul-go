// Package testutil takes care of initializing a local git repository for
// testing. The 'remote' should match the one specified in config/mock.
package testutil

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"gopkg.in/libgit2/git2go.v24"
)

var testRepo *git.Repository

// Return the test-fixtures path in testutil
func fixturesRepo(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Cannot find path")
	}

	testutilPath := filepath.Dir(filename)
	return filepath.Join(testutilPath, "test-fixtures", "example")
}

func copyDir(srcPath string, dstPath string) error {
	// Copy fixtures into temporary path. filepath is the full path
	var copyFn = func(path string, info os.FileInfo, err error) error {
		currentFilePath := strings.TrimPrefix(path, srcPath)
		targetPath := filepath.Join(dstPath, currentFilePath)
		if info.IsDir() {
			if targetPath != dstPath {
				err := os.Mkdir(targetPath, 0755)
				if err != nil {
					return err
				}
			}
		} else {
			src, err := os.Open(path)
			if err != nil {
				return err
			}
			dst, err := os.Create(targetPath)
			if err != nil {
				return err
			}

			_, err = io.Copy(dst, src)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := filepath.Walk(srcPath, copyFn)
	return err
}

// GitInitTestRepo coppies test-fixtures to os.TempDir() and performs a
// git-init on directory.
func GitInitTestRepo(t *testing.T) (*git.Repository, func()) {
	fixtureRepo := fixturesRepo(t)
	repoPath, err := ioutil.TempDir("", "git2consul-test-remote")
	if err != nil {
		t.Fatal(err)
	}

	err = copyDir(fixtureRepo, repoPath)
	if err != nil {
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
	testRepo = repo

	// Cleanup function that removes the repository directory
	var cleanup = func() {
		os.RemoveAll(repoPath)
	}

	return repo, cleanup
}

// GitCommitTestRepo performs a commit on the test repository, and returns
// its Oid as well as a cleanup function to revert those changes.
func GitCommitTestRepo(t *testing.T) (*git.Oid, func()) {
	// Save commmit ref for reset later
	h, err := testRepo.Head()
	if err != nil {
		t.Fatal(err)
	}

	obj, err := testRepo.Lookup(h.Target())
	if err != nil {
		t.Fatal(err)
	}

	initialCommit, err := obj.AsCommit()
	if err != nil {
		t.Fatal(err)
	}

	date := []byte(time.Now().String())
	file := path.Join(testRepo.Workdir(), "foo")
	err = ioutil.WriteFile(file, date, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Commit changes
	idx, err := testRepo.Index()
	if err != nil {
		t.Fatal(err)
	}

	treeId, err := idx.WriteTree()

	tree, err := testRepo.LookupTree(treeId)
	if err != nil {
		t.Fatal(err)
	}

	h, err = testRepo.Head()
	if err != nil {
		t.Fatal(err)
	}

	commit, err := testRepo.LookupCommit(h.Target())
	if err != nil {
		t.Fatal(err)
	}

	sig := &git.Signature{
		Name:  "Test Example",
		Email: "tes@example.com",
		When:  time.Date(2016, 01, 01, 12, 00, 00, 0, time.UTC),
	}

	oid, err := testRepo.CreateCommit("HEAD", sig, sig, "Update commit", tree, commit)
	if err != nil {
		t.Fatal(err)
	}

	// Undo commit
	var cleanup = func() {
		testRepo.ResetToCommit(initialCommit, git.ResetHard, &git.CheckoutOpts{
			Strategy: git.CheckoutForce,
		})

		testRepo.StateCleanup()
	}

	return oid, cleanup
}
