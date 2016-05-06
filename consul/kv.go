package consul

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/repository"
	"github.com/hashicorp/consul/api"
	"github.com/libgit2/git2go"
)

// TODO: Optimize for PUT only on changes instead of the entire repo
// TODO: Need to also push if key is absent
// Push a repository to the KV
func (c *Client) pushBranch(repo *repository.Repository, branch *git.Branch) {
	// Checkout branch
	repo.CheckoutBranch(branch, &git.CheckoutOpts{})

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

		kvPath := path.Join(repo.Name(), branchName, path.Base(fullpath))

		kv := c.KV()
		data, err := ioutil.ReadFile(fullpath)
		if err != nil {
			return err
		}

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
}

// Push a single file
func (c *Client) pushFile(fullpath string, info os.FileInfo, err error) error {
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

	kv := c.KV()
	data, err := ioutil.ReadFile(fullpath)
	if err != nil {
		return err
	}

	// KV path, is repo/branch/file
	kvPath := path.Join("test", path.Base(fullpath))

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
