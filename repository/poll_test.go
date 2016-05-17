package repository

import (
	"os"
	"testing"
)

func TestPollBranches(t *testing.T) {

	r, cleanup := tempGitInitPath(t)
	defer cleanup()

	cfg := loadConfig(t)

	repos, err := loadRepos(cfg)
	if err != nil {
		t.Fatal(err)
	}
	repo := repos[0]

	tempCommitRepo(r, t)

	err = repo.PollBranches()
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup
	defer func() {
		repos[0].CloneCh()
		err = os.RemoveAll(repo.store)
		if err != nil {
			t.Fatal(err)
		}
	}()
}
