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
	store      config.Config.LocalStore
}

// Poll all repos, and either clone or pull
// TODO: Pass in local store, use channel and go to poll indefinitely
func Poll(cfg *config.Config) error {

	// TODO: goroutine on for interval polling
	for _, repo := range cfg.Repos {
		path := filepath.Join(cfg.LocalStore, repo.Name)
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

	// go r.PollRepoByInterval()
	// go r.PollRepoByWebhook()

	return nil
}

func PollRepos(cfg *config.Config) error {
	for _, repo := range cfg.Repos {
		// Create Repository object
		store := filepath.Join(cfg.LocalStore, repo.Name)
		r := &Repository{
			raw_repo,
			repo,
			store,
		}

		// Poll repository by interval, or webhook
		go repo.PollRepoByInterval()
		// go r.PollRepoByWebhook()
	}

	return nil
}

func (r *Repository) Poll() error {
	if _, err := os.Stat(r.store); err != nil {
		// If there is no repo, create and clone
		if os.IsNotExist(err) {
			log.Infof("%s not cached, cloning to %s", repo.Name, path)
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

	return nil
}

func (r *Repository) PollRepoByInterval() {
	hooks := r.repoConfig.Hooks
	interval := time.Second

	for _, h := range hooks {
		if h.Type == "polling" {
			interval = h.Interval
			break
		}
	}

	// If no polling found, don't poll
	if interval == 0 {
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			r.Poll()
		}
	}
}

func (r *Repository) PollRepoByWebhook() {

}
