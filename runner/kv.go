package runner

import (
	"fmt"
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

// Get local branch ref from the KV
func (r *Runner) getKVRef(repo *repository.Repository, branchName string) (string, error) {
	refFile := fmt.Sprintf("%s.ref", branchName)
	key := path.Join(repo.Name(), refFile)

	kv := r.client.KV()
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
func (r *Runner) putKVRef(repo *repository.Repository, branchName string) error {
	refFile := fmt.Sprintf("%s.ref", branchName)
	key := path.Join(repo.Name(), refFile)

	rawRef, err := repo.References.Lookup("refs/heads/" + branchName)
	if err != nil {
		return err
	}
	ref := rawRef.Target().String()

	kv := r.client.KV()

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

// Push a repository branch to the KV
// TODO: Optimize for PUT only on changes instead of the entire repo
func (r *Runner) putBranch(repo *repository.Repository, branch *git.Branch) error {
	// Checkout branch
	repo.CheckoutBranch(branch, &git.CheckoutOpts{
		Strategy: git.CheckoutForce,
	})

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

		log.Debugf("(consul) pushBranch(): Path: %s Key: %s", fullpath, strings.TrimPrefix(fullpath, repo.Store()))
		kvPath := path.Join(repo.Name(), branchName, strings.TrimPrefix(fullpath, repo.Store()))

		kv := r.client.KV()
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
