package repository

import (
	"os"
	"path"
	"testing"
)

func TestClone(t *testing.T) {
	repo, cleanup := tempGitInitPath(t)
	defer cleanup()

	cfg := loadConfig(t)

	r := &Repository{
		Repository: repo,
		repoConfig: cfg.Repos[0],
		store:      path.Join(cfg.LocalStore, cfg.Repos[0].Name),
	}

	err := r.Clone()
	if err != nil {
		t.Fatal(err)
	}

	//Cleanup cloned repo
	defer func() {
		err = os.RemoveAll(r.store)
		if err != nil {
			t.Fatal(err)
		}
	}()
}
