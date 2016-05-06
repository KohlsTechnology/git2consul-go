package consul

import (
	"fmt"
	"path"

	"github.com/cleung2010/go-git2consul/repository"
	"github.com/hashicorp/consul/api"
)

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
