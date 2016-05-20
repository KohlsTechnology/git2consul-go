package kv

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/cleung2010/go-git2consul/repository"
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

		// log.Debugf("(consul) pushBranch(): Path: %s Key: %s", fullpath, strings.TrimPrefix(fullpath, repo.Store()))
		key := strings.TrimPrefix(fullpath, repo.Workdir())
		kvPath := path.Join(repo.Name(), branchName, key)

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

		log.Debugf("(consul)(trace): PUT KV Path: %s Key: %s", fullpath, filepath.Base(repo.Workdir()))

		return nil
	}

	err := filepath.Walk(repo.Workdir(), pushFile)
	if err != nil {
		log.WithError(err).Debug("PUT branch error")
	}

	return nil
}

func (h *KVHandler) putKV(repo *repository.Repository, prefix string) error {
	head, err := repo.Head()
	if err != nil {
		return err
	}

	branchName, err := head.Branch().Name()
	if err != nil {
		return err
	}

	key := path.Join(repo.Name(), branchName, prefix)
	filePath := filepath.Join(repo.Workdir(), prefix)
	value, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	p := &api.KVPair{
		Key:   key,
		Value: value,
	}

	_, err = h.Put(p, nil)
	if err != nil {
		return err
	}

	return nil
}
