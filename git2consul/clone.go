package git2consul

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/libgit2/git2go"
)

func CloneRepos(c *config.Config) error {
	log.Debug("Cloning repositories %+v", c.Repos)
	log.Debugf("Using %s as the local storage", os.TempDir())

	for _, repo := range c.Repos {
		path := filepath.Join(os.TempDir(), repo.Name)
		log.Debugf("Cloning %s from %s", repo.Name, repo.Url)

		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}

		repo, err := git.Clone(repo.Url, path, &git.CloneOptions{})
		if err != nil {
			return err
		}
		log.Info(repo)
	}

	return nil
}
