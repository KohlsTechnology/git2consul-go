package git2consul

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/libgit2/git2go"
)

func CloneRepos(c *config.Config) error {
	log.Debug("Cloning repositories %+v", c.Repos)
	log.Debugf("Using %s as the local storage", os.TempDir())
	for _, repo := range c.Repos {
		log.Debugf("Cloning %s from %s", repo.Name, repo.Url)
		_, err := git.Clone(repo.Url, os.TempDir(), &git.CloneOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
