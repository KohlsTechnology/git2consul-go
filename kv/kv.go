package kv

import (
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/cleung2010/go-git2consul/repository"
	"github.com/hashicorp/consul/api"
)

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

func (h *KVHandler) deleteKV(repo *repository.Repository, prefix string) error {
	head, err := repo.Head()
	if err != nil {
		return err
	}

	branchName, err := head.Branch().Name()
	if err != nil {
		return err
	}

	key := path.Join(repo.Name(), branchName, prefix)

	_, err = h.Delete(key, nil)
	if err != nil {
		return err
	}

	return nil
}
