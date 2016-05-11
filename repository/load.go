package repository

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"gopkg.in/libgit2/git2go.v24"
)

func loadRepos(cfg *config.Config) (Repositories, error) {
	repos := []*Repository{}

	// Create Repository object for each repo
	for _, cRepo := range cfg.Repos {
		store := filepath.Join(cfg.LocalStore, cRepo.Name)

		r := &Repository{
			repoConfig: cRepo,
			store:      store,
			cloneCh:    make(chan struct{}, 1),
			changeCh:   make(chan struct{}, 1),
		}

		// Check if directory exists
		fi, err := os.Stat(store)
		if os.IsNotExist(err) || fi.IsDir() == false {
			err := os.Mkdir(r.store, 0755)
			if err != nil {
				return nil, err
			}

			err = r.Clone()
			if err != nil {
				return nil, err
			}
			log.Infof("(git): Repository %s not cached, cloned to %s", r.repoConfig.Name, r.store)
		} else if err != nil {
			return nil, err
		} else {
			repo, err := git.OpenRepository(store)
			if err != nil {
				return nil, err
			}

			// Check if config repo and cached repo matches
			rm, err := repo.Remotes.Lookup("origin")
			if err != nil {
				return nil, err
			}

			absPath, err := filepath.Abs(cRepo.Url)
			if err != nil {
				return nil, err
			}
			// If not equal attempt to recreate the repo
			if strings.Compare(rm.Url(), absPath) != 0 {
				log.Warnf("Diffrent %s repository exists locally, overwritting", cRepo.Name)
				os.RemoveAll(store) // Potentially dangerous?
				err := os.Mkdir(r.store, 0755)
				if err != nil {
					return nil, err
				}

				err = r.Clone()
				if err != nil {
					return nil, err
				}
			} else {
				r.Repository = repo
			}
		}

		repos = append(repos, r)
	}

	return repos, nil
}
