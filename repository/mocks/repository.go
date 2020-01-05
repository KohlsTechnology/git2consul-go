/*
Copyright 2019 Kohl's Department Stores, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mocks

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// InitRemote TODO write useful documentation here
func InitRemote(t *testing.T) (*git.Repository, string) {
	repoPath, err := ioutil.TempDir("", "git2consul-remote")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Initializing new repo in %s", repoPath)
	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		t.Fatal(err)
	}

	Add(t, repo, "example/foo.txt", []byte("Example content foo.txt"))
	Add(t, repo, "example/boo.txt", []byte("Example content boo.txt"))
	Commit(t, repo, "Initial commit")
	return repo, repoPath
}

// Add TODO write useful documentation here
func Add(t *testing.T, repo *git.Repository, path string, content []byte) {
	w, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	root := w.Filesystem.Root()
	dir := filepath.Dir(path)
	fileName := filepath.Base(path)
	_, err = os.Stat(filepath.Join(root, dir))
	if os.IsNotExist(err) {
		err := os.Mkdir(filepath.Join(root, dir), 0700)
		if err != nil {
			t.Fatal(err)
		}
	}
	f, err := os.Create(filepath.Join(root, dir, fileName))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		t.Fatal(err)
	}
	w.Add(path)
}

func getSignature() *object.Signature {
	when := time.Now()
	sig := &object.Signature{
		Name:  "foo",
		Email: "foo@foo.foo",
		When:  when,
	}
	return sig
}

// Commit TODO write useful documentation here
func Commit(t *testing.T, repo *git.Repository, message string) {
	w, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	w.Commit(message, &git.CommitOptions{Author: getSignature()})
}
