package repository

import (
	"testing"

	"github.com/cleung2010/go-git2consul/test"
)

func TestClone(t *testing.T) {
	repo, cleanup := test.TempGitInitPath(test.TestRepo(), t)
	defer cleanup()

	cfg := test.DefaultConfig()

}
