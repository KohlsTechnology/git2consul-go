package repository

import (
	"os"
	"path"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/libgit2/git2go"
)

// TODO: Probably needs a lock on the repo for push/pull
type Repository struct {
	*git.Repository
	RepoConfig *config.Repo
	Store      string
}

func PollRepos(cfg *config.Config) ([]Repository, error) {
	repos := []Repository{}
	for _, repo := range cfg.Repos {
		// Create Repository object for each repo
		store := filepath.Join(cfg.LocalStore, repo.Name)
		raw_repo, err := git.OpenRepository(store)
		if err != nil {
			log.Debugf("Cannot open repository: %s", err)
		}
		r := &Repository{
			raw_repo,
			repo,
			store,
		}

		// TODO: Fix this
		repos = append(repos, *r)

		// Poll repository by interval, or webhook
		go r.pollRepoByInterval()
		// go r.PollRepoByWebhook()
	}

	return repos, nil
}

func (r *Repository) poll() error {
	if _, err := os.Stat(r.Store); err != nil {
		// If there is no repo, create and clone
		if os.IsNotExist(err) {
			log.Infof("Repository %s not cached, cloning to %s", r.RepoConfig.Name, r.Store)
			err := os.Mkdir(r.Store, 0755)
			if err != nil {
				return err
			}

			err = r.Clone()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// Pull the repository, all specified branches
		for _, branch := range r.RepoConfig.Branches {
			r.Pull(branch)
		}
	}

	return nil
}

func (r *Repository) pollRepoByInterval() {
	hooks := r.RepoConfig.Hooks
	interval := time.Second

	// Find polling hook
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

	// Initial poll
	err := r.poll()
	if err != nil {
		log.Error(err)
	}

	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-ticker.C:
			err := r.poll()
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (r *Repository) pollRepoByWebhook() {

}

func (r *Repository) GetBranchKVPath(branchName) (string, error) {
	path := path.Join(r.RepoConfig.Name, r.)

}
