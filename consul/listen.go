package consul

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/repository"
	"github.com/hashicorp/consul/api"
)

// Listen for changes on all registered repos
func (c *Client) ListenForChanges(repos []repository.Repository) {
	// If there changes, push to KV
	for _, r := range repos {
		go c.listenOnRepo(&r)
	}
}

func (c *Client) listenOnRepo(repo *repository.Repository) {
	// If change is detected, push repository to the KV
	for {
		// Check the ref for the branch
		branch, err := repo.Head().Branch().Name()
		if err != nil {
			return err
		}
		kv := c.KV()
		kv.Get(, nil)

	}
}


// TODO: Optimize for PUT only on changes instead of the entire repo
// TODO: Need to also push if key is absent
// Push a repository to the KV
func (c *Client) pushRepo(repo *repository.Repository) {
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

		kv := c.KV()
		data, err := ioutil.ReadFile(fullpath)
		if err != nil {
			return err
		}

		// KV path, is repo/branch/file
		h, err := repo.Head()
		if err != nil {
			return err
		}
		branch, err := h.Branch().Name()
		if err != nil {
			return err
		}
		kvPath := path.Join(repo.RepoConfig.Name, branch, path.Base(fullpath))

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

	err := filepath.Walk(repo.Store, pushFile)
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
