package repository

import (
	"os"
	"path/filepath"
	"time"

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
func Poll(repos []*config.Repo) error {

	// TODO: goroutine on for interval polling
	for _, repo := range repos {
		path := defaultStore(repo.Name)
		log.Infof("Polling repository: %s from %s", repo.Name, repo.Url)

		if _, err := os.Stat(path); err != nil {
			// If there is no repo, create and clone
			if os.IsNotExist(err) {
				log.Infof("%s does not cached, cloning to %s", repo.Name, path)
				err := os.Mkdir(path, 0755)
				if err != nil {
					return err
				}

				_, err = Clone(repo)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// Pull the repository
			log.Infof("Pulling commits to repository %s", repo.Name)
			raw_repo, err := git.OpenRepository(path)
			if err != nil {
				return err
			}
			r := &Repository{
				raw_repo,
				repo,
			}

			err = r.Pull()
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (r *Repository) PollRepoByInterval(d time.Duration) {

}

func (r *Repository) PollRepoByWebhook() {

}
