package kv

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/Cimpress-MCP/go-git2consul/repository"
	"github.com/hashicorp/consul/api"
	"gopkg.in/libgit2/git2go.v24"
)

// Push a repository branch to the KV
// TODO: Optimize for PUT only on changes instead of the entire repo
func (h *KVHandler) putBranch(repo *repository.Repository, branch *git.Branch) error {
	// Checkout branch
	repo.CheckoutBranch(branch, &git.CheckoutOpts{
		Strategy: git.CheckoutForce,
	})

	// h, _ := repo.Head()
	// bn, _ := h.Branch().Name()
	// log.Debugf("(consul) pushBranch(): Branch: %s Head: %s", bn, h.Target().String())

	var pushFile = func(fullpath string, info os.FileInfo, err error) error {
		// Walk error
		if err != nil {
			return err
		}

		// Skip the .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Do not push directories
		if info.IsDir() {
			return nil
		}

		// KV path, is repo/branch/file
		branchName, err := branch.Name()
		if err != nil {
			return err
		}

		key := strings.TrimPrefix(fullpath, repo.Workdir())
		kvPath := path.Join(repo.Name(), branchName, key)
		h.logger.Debugf("KV PUT changes: %s/%s: %s", repo.Name(), branchName, kvPath)

		data, err := ioutil.ReadFile(fullpath)
		if err != nil {
			return err
		}

		// log.Debugf("(consul) pushBranch(): Data: %s", data)

		p := &api.KVPair{
			Key:   kvPath,
			Value: data,
		}

		_, err = h.Put(p, nil)
		if err != nil {
			return err
		}

		return nil
	}

	err := filepath.Walk(repo.Workdir(), pushFile)
	if err != nil {
		log.WithError(err).Debug("PUT branch error")
	}

	return nil
}
