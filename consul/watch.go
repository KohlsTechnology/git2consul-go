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

func (c *Client) watchRepo(repo *repository.Repository) {
	// TODO: Initial GET on refs

	// Continuously watch changes on repo
	for {
		// Block until change is received
		<-repo.ChangeLock()

		itr, err := repo.NewBranchIterator(git.BranchLocal)
		if err != nil {
			log.Error(err)
		}

		var updateBranchFn = func(b *git.Branch, _ git.BranchType) error {
			bn, err := b.Name()
			if err != nil {
				return err
			}
			// log.Debugf("Updating for branch: %s", bn)
			log.Debugf("KV GET ref for %s/%s", repo.Name(), bn)
			ref, err := c.getKVRef(repo, bn)
			if err != nil {
				return err
			}

			// If ref doesn't exist or is not the same, push files to KV
			if ref == nil || string(ref) != b.Reference.Target().String() {
				repo.Lock()
				log.Debugf("KV PUT changes for %s/%s", repo.Name(), bn)
				c.pushBranch(repo, bn)
				log.Debugf("KV PUT ref for %s/%s", repo.Name(), bn)
				c.putKVRef(repo, bn)
				repo.Unlock()
			}

			return nil
		}

		// Update KV
		itr.ForEach(updateBranchFn)
	}
}

// Get local branch ref from the KV
func (c *Client) getKVRef(repo *repository.Repository, branchName string) ([]byte, error) {
	refFile := fmt.Sprintf("%s.ref", branchName)
	key := path.Join(repo.Name(), refFile)

	kv := c.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	// If error on get, return empty value
	if pair == nil {
		return nil, nil
	}

	return pair.Value, nil
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
