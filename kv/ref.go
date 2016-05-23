package kv

import (
	"fmt"
	"path"

	"github.com/cleung2010/go-git2consul/repository"
	"github.com/hashicorp/consul/api"
)

// Get local branch ref from the KV
func (h *KVHandler) getKVRef(repo *repository.Repository, branchName string) (string, error) {
	refFile := fmt.Sprintf("%s.ref", branchName)
	key := path.Join(repo.Name(), refFile)

	pair, _, err := h.Get(key, nil)
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
func (h *KVHandler) putKVRef(repo *repository.Repository, branchName string) error {
	refFile := fmt.Sprintf("%s.ref", branchName)
	key := path.Join(repo.Name(), refFile)

	rawRef, err := repo.References.Lookup("refs/heads/" + branchName)
	if err != nil {
		return err
	}
	ref := rawRef.Target().String()

	p := &api.KVPair{
		Key:   key,
		Value: []byte(ref),
	}

	_, err = h.Put(p, nil)
	if err != nil {
		return err
	}

	return nil
}
