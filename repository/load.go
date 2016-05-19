package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"gopkg.in/libgit2/git2go.v24"
)

// Populates Repository slice from configuration. It also
// handles cloning of the repository if not present
func LoadRepos(cfg *config.Config) (Repositories, error) {
	repos := []*Repository{}

	// Create Repository object for each repo
	for _, cRepo := range cfg.Repos {
		store := filepath.Join(cfg.LocalStore, cRepo.Name)

		r := &Repository{
			repoConfig: cRepo,
			store:      store,
			changeCh:   make(chan struct{}, 1),
		}

		fi, err := os.Stat(store)
		if os.IsNotExist(err) || fi.IsDir() == false {
			log.Infof("(git): Repository %s not cached, cloning to %s", cRepo.Name, store)
			err := os.Mkdir(store, 0755)
			if err != nil {
				return nil, err
			}

			err = r.Clone()
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			// Some other stat error
			return nil, err
		} else {
			// The directory exists

			// Not a git repository, remove directory and clone
			_, err := os.Stat(filepath.Join(store, ".git"))
			if os.IsNotExist(err) {
				log.Warnf("(git): %s exists locally, overwritting", cRepo.Name)
				err := os.RemoveAll(store)
				if err != nil {
					return nil, err
				}

				err = os.Mkdir(store, 0755)
				if err != nil {
					return nil, err
				}

				err = r.Clone()
				if err != nil {
					return nil, err
				}
			}

			// Open repository otherwise
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
				log.Warnf("(git): Diffrent %s repository exists locally, overwritting", cRepo.Name)
				os.RemoveAll(store) // Potentially dangerous?
				err := os.Mkdir(store, 0755)
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

	if len(repos) == 0 {
		return repos, fmt.Errorf("No repositories provided in the configuration")
	}

	return repos, nil
}
