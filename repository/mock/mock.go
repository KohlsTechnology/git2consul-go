package mock

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Cimpress-MCP/go-git2consul/config/mock"
	"github.com/Cimpress-MCP/go-git2consul/repository"

	"gopkg.in/libgit2/git2go.v24"
)

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

func fixturesPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {

	}

	callerPath := filepath.Dir(filename)
	return filepath.Join(callerPath, "_fixtures", "example")
}

func initGitRepository() *git.Repository {
	// repo := &git.Repository{}
	// repoPath := repoPath()

	repoPath, err := ioutil.TempDir("", "git2consul-test-local")
	if err != nil {
		return nil
	}

	err = copyDir(fixturesPath(), repoPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	repo, err := git.InitRepository(fixturesPath(), false)
	if err != nil {
		return nil
	}

	return repo
}

// Returns a mock of a repository.Repository object
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
