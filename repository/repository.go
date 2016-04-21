package repository

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/libgit2/git2go"
)

type Repository struct {
	*git.Repository
	repoConfig *config.Repo
}

// Return the default cache store, which is the OS' temp dir
func defaultStore(name string) string {
	return filepath.Join(os.TempDir(), name)
}

// Poll all repos, and either clone or pull
// TODO: Pass in local store, use channel and go to poll indefinitely
func Poll(repos *config.Repos) error {

	// TODO: goroutine on each repo whether to clone or poll
	for _, repo := range *repos {
		path := defaultStore(repo.Name)
		log.Infof("Polling repository: %s from %s", repo.Name, repo.Url)

		// Make directory
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}

		_, err = Clone(&repo)
		if err != nil {
			return err
		}

	}

	return nil
}
