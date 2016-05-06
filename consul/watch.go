package consul

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/repository"
	"github.com/hashicorp/consul/api"
	"github.com/libgit2/git2go"
)

// Listen for changes on all registered repos
func (c *Client) WatchChanges(repos []*repository.Repository) {
	// If there changes, push to KV
	for _, r := range repos {
		go c.watchRepo(r)
	}
}

// Watch for changes on a repository
// TODO: Handle errors through channel
func (c *Client) watchRepo(repo *repository.Repository) error {
	// TODO: Initial GET on refs

	for {
		// Block until change is received
		<-repo.ChangeLock()
		repo.Lock()
		h, err := repo.Head()
		if err != nil {
			return err
		}
		b, err := h.Branch().Name()
		if err != nil {
			return err
		}

		log.Debugf("KV GET ref for %s/%s", repo.Name(), b)
		kvRef, err := c.getKVRef(repo, b)
		if err != nil {
			return err
		}

		// Local ref
		localRef := h.Target().String()

		if len(kvRef) == 0 || kvRef != localRef {
			log.Debugf("KV PUT changes for %s/%s", repo.Name(), b)
			c.pushBranch(repo, b)
			log.Debugf("KV PUT ref for %s/%s", repo.Name(), b)
			c.putKVRef(repo, b)
		}

		repo.Unlock()
	}
}

// Get local branch ref from the KV
func (c *Client) getKVRef(repo *repository.Repository, branchName string) (string, error) {
	refFile := fmt.Sprintf("%s.ref", branchName)
	key := path.Join(repo.Name(), refFile)

	kv := c.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return "", err
	}

	// If error on get, return empty value
	if pair == nil {
		return "", nil
	}

	return string(pair.Value), nil
}

// Put the local branch ref to the KV
func (c *Client) putKVRef(repo *repository.Repository, branchName string) error {
	refFile := fmt.Sprintf("%s.ref", branchName)
	key := path.Join(repo.Name(), refFile)

	rawRef, err := repo.References.Lookup("refs/heads/" + branchName)
	if err != nil {
		return err
	}
	ref := rawRef.Target().String()

	kv := c.KV()

	p := &api.KVPair{
		Key:   key,
		Value: []byte(ref),
	}

	_, err = kv.Put(p, nil)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Optimize for PUT only on changes instead of the entire repo
// TODO: Need to also push if key is absent
// Push a repository to the KV
func (c *Client) pushBranch(repo *repository.Repository, branchName string) {
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
		b, err := repo.LookupBranch(branchName, git.BranchLocal)
		if err != nil {
			return err
		}
		branch, err := b.Name()
		if err != nil {
			return err
		}

		kvPath := path.Join(repo.Name(), branch, path.Base(fullpath))

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
