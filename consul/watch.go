package consul

import "github.com/cleung2010/go-git2consul/repository"

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
	// Initial GET on refs
	c.handleClone(repo)

	for {
		select {
		// case <-repo.CloneCh():
		// 	c.handleClone(repo)
		case <-repo.ChangeCh():
			c.handleChange(repo)
		}
	}
}
