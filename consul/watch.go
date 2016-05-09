package consul

import "github.com/cleung2010/go-git2consul/repository"

// Listen for changes on all registered repos
func (c *Client) WatchChanges(repos []*repository.Repository) error {
	errCh := make(chan error, 1)

	// If there changes, push to KV
	for _, r := range repos {
		go c.watchRepo(r, errCh)
	}

	select {
	case err := <-errCh:
		return err
	}
}

// Watch for changes on a repository
// TODO: Handle errors through channel
func (c *Client) watchRepo(repo *repository.Repository, errCh chan error) {
	// Initial GET on refs
	err := c.handleClone(repo)
	if err != nil {
		errCh <- err
		return
	}

	for {
		select {
		// case <-repo.CloneCh():
		// 	c.handleClone(repo)
		case <-repo.ChangeCh():
			err := c.handleChange(repo)
			if err != nil {
				errCh <- err
				return
			}
		}
	}
}
