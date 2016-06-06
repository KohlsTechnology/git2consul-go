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

	"github.com/Cimpress-MCP/go-git2consul/config"

	"gopkg.in/libgit2/git2go.v24"
)

var testRepo *git.Repository

// Get the test-fixtures directory
func FixturesPath(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Cannot find path")
	}

	testutilPath := filepath.Dir(filename)
	return filepath.Join(testutilPath, "test-fixtures")
}

// Copy test-fixtures to os.TempDir() and performs a git-init on directory.
// Returns git.Repository object and the cleanup function
func GitInitTestRepo(t *testing.T) (*git.Repository, func()) {
	fixturePath := FixturesPath(t)
	fixtureRepo := filepath.Join(fixturePath, "example")
	repoPath, err := ioutil.TempDir("", filepath.Join("test-git2consul-example"))
	if err != nil {
		t.Fatal(err)
	}

	// Copy fixtures into temporary path. filepath is the full path
	var copyFn = func(filepath string, info os.FileInfo, err error) error {
		projectFilePath := strings.TrimPrefix(fixtureRepo, filepath)
		targetPath := path.Join(repoPath, projectFilePath)
		if info.IsDir() {
			err := os.Mkdir(targetPath, 0755)
			if err != nil {
				return err
			}
		} else {
			src, err := os.Open(filepath)
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

	err = filepath.Walk(fixtureRepo, copyFn)

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

	// Reset to initial commit, and then remove .git
	var cleanup = func() {
		os.RemoveAll(repoPath)
	}

	testRepo = repo

	return repo, cleanup
}

func LoadTestConfig(t *testing.T) *config.Config {
	// Verify that testRepo is not nil
	if testRepo == nil {
		t.Fatal("testRepo not initialized")
	}

	fixturePath := FixturesPath(t)
	cfgPath := filepath.Join(fixturePath, "example.json")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	// Change defaults to use test settings
	cfg.LocalStore = os.TempDir()
	cfg.Repos[0].Url = testRepo.Workdir()

	return cfg
}

func TempCommitTestRepo(t *testing.T) func() {
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

	_, err = testRepo.CreateCommit("HEAD", sig, sig, "Update commit", tree, commit)
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

	return cleanup
}
