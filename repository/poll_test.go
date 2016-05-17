package repository

import (
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/cleung2010/go-git2consul/test"
)

func TestPollBranches(t *testing.T) {
	_, cleanup := test.TempGitInitPath(test.TestRepo(), t)
	defer cleanup()

	file := path.Join(test.TestRepo(), "foo")
	err := ioutil.WriteFile(file, []byte(time.Now().String()), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Commit changes
	// Poll for changes
}
