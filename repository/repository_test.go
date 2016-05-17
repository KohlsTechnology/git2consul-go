package repository

import (
	"testing"

	"github.com/cleung2010/go-git2consul/config"
	"github.com/cleung2010/go-git2consul/test"
)

func TestLoadRepos(t *testing.T) {
	_, cleanup := test.TempGitInitPath(test.TestRepo(), t)
	defer cleanup()

	cfg, err := config.Load(test.TestConfig())
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadRepos(cfg)
	if err != nil {
		t.Fatal(err)
	}
}

// func TestLoadRepos_invalidRepo(t *testing.T) {
// 	cfg, err := config.Load(path.Join(test.ConfigPath, "local.json"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = LoadRepos(cfg)
// 	if err == nil {
// 		t.Fail()
// 	}
// }
