package consul

// import (
// 	"io/ioutil"
//
// 	"github.com/cleung2010/go-git2consul/repository"
// 	"github.com/hashicorp/consul/api"
// )
//
// // TODO: Optimize for PUT only on changes instead of the entire repo
//
// // Push a repository to the KV
// func (c *Client) Push(repo *repository.Repository) {
//
// }
//
// // Push a single file
// func (c *Client) pushFile(path string) error {
// 	kv := c.KV()
// 	data, err := ioutil.ReadFile(path)
//
// 	p := &api.KVPair{
// 		Key:   path,
// 		Value: data,
// 	}
//
// 	_, err = kv.Put(p, nil)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
