package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Cimpress-MCP/go-git2consul/config"
	"gopkg.in/libgit2/git2go.v24"
)

// Repository is used to hold the git repository object and it's configuration
type Repository struct {
	sync.Mutex

	*git.Repository
	Config *config.Repo
}

// Status codes for Repository object creation
const (
	RepositoryError = iota // Unused, it will always get returned with an err
	RepositoryCloned
	RepositoryOpened
)

// Name returns the repository name
func (r *Repository) Name() string {
	return r.Config.Name
}

// Branch returns the branch name
func (r *Repository) Branch() string {
	head, err := r.Head()
	if err != nil {
		return ""
	}
	bn, err := head.Branch().Name()
	if err != nil {
		return ""
	}

	return bn
}

// New is used to construct a new repository object from the configuration
func New(repoBasePath string, repoConfig *config.Repo) (*Repository, int, error) {
	repoPath := filepath.Join(repoBasePath, repoConfig.Name)

	r := &Repository{
		Repository: &git.Repository{},
		Config:     repoConfig,
	}

	state, err := r.init(repoPath)
	if err != nil {
		return nil, RepositoryError, err
	}

	if r.Repository == nil {
		return nil, RepositoryError, fmt.Errorf("Could not find git repostory")
	}

	return r, state, nil
}

// Initialize git.Repository object by opening the repostory or cloning from
// the source URL. It does not handle purging existing file or directory
// with the same path
func (r *Repository) init(repoPath string) (int, error) {
	gitRepo, err := git.OpenRepository(repoPath)
	if err != nil || gitRepo == nil {
		err := r.Clone(repoPath)
		if err != nil {
			return RepositoryError, err
		}
		return RepositoryCloned, nil
	}

	// If remote URL are not the same, it will purge local copy and re-clone
	if r.mismatchRemoteUrl(gitRepo) {
		os.RemoveAll(gitRepo.Workdir())
		err := r.Clone(repoPath)
		if err != nil {
			return RepositoryError, err
		}
		return RepositoryCloned, nil
	}

	r.Repository = gitRepo

	return RepositoryOpened, nil
}

func (r *Repository) mismatchRemoteUrl(gitRepo *git.Repository) bool {
	rm, err := gitRepo.Remotes.Lookup("origin")
	if err != nil {
		return true
	}

	if strings.Compare(rm.Url(), r.Config.Url) != 0 {
		return true
	}

	return false
}
