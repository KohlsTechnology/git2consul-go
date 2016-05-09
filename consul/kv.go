package consul

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/repository"
	"github.com/hashicorp/consul/api"
	"gopkg.in/libgit2/git2go.v23"
)

// Push a repository to the KV
// TODO: Optimize for PUT only on changes instead of the entire repo
func (c *Client) pushBranch(repo *repository.Repository, branch *git.Branch) error {
	// Checkout branch
	repo.CheckoutBranch(branch, &git.CheckoutOpts{})

	h, _ := repo.Head()
	bn, _ := h.Branch().Name()
	log.Debugf("(consul) pushBranch(): Branch: %s Head: %s", bn, h.Target().String())

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

		log.Debugf("(consul) pushBranch(): Path: %s Key: %s", fullpath, strings.TrimLeft(fullpath, repo.Store()))
		kvPath := path.Join(repo.Name(), branchName, strings.TrimPrefix(fullpath, repo.Store()))

		kv := c.KV()
		data, err := ioutil.ReadFile(fullpath)
		if err != nil {
			return err
		}

		log.Debugf("(consul) pushBranch(): Data: %s", data)

		p := &api.KVPair{
			Key:   kvPath,
			Value: data,
		}

		_, err = kv.Put(p, nil)
		if err != nil {
			return err
		}

		return nil
	}

	err := filepath.Walk(repo.Store(), pushFile)
	if err != nil {
		log.Debug(err)
	}

	return nil
}
