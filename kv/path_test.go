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

package kv

import (
	"testing"

	"github.com/KohlsTechnology/git2consul-go/config"
	"github.com/KohlsTechnology/git2consul-go/kv/mocks"
	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/stretchr/testify/assert"
)

func TestPathHandlerWithoutMountPointAndSourceRoot(t *testing.T) {
	var repo repository.Repo = &mocks.Repo{Config: &config.Repo{}}

	filePath := "first_level/second_level/foo"
	expectedKey := "repository_mock/master/first_level/second_level/foo"

	key, status, err := pathHandler(repo, filePath)
	if err != nil {
		if status == PathFormatterError {
			t.Fatalf("Could not set the branch: %s", err)
		}
		t.Fatal(err)
	}
	assert.Equal(t, key, expectedKey)
}

func TestPathHandlerWithMountPointAndSourceRoot(t *testing.T) {
	var repo repository.Repo = &mocks.Repo{Config: &config.Repo{}}
	repo.GetConfig().MountPoint = "test_mountpoint/"
	repo.GetConfig().SourceRoot = "first_level/second_level"
	filePath := "first_level/second_level/foo"
	expectedKey := "test_mountpoint/repository_mock/master/foo"

	key, status, err := pathHandler(repo, filePath)
	if err != nil {
		if status == PathFormatterError {
			t.Fatalf("Could not set the branch: %s", err)
		}
		t.Fatal(err)
	}
	assert.Equal(t, key, expectedKey)
}

func TestPathBaseBuilderWithoutMountPoint(t *testing.T) {
	var repo repository.Repo = &mocks.Repo{Config: &config.Repo{}}
	repo.GetConfig().MountPoint = ""
	expectedKey := "repository_mock/master"
	key, status, err := pathBaseBuilder(repo)
	if err != nil {
		if status == PathFormatterError {
			t.Fatalf("Could not set the branch: %s", err)
		}
	}
	assert.Equal(t, key, expectedKey)
}

func TestPathBaseBuilderWithMountPoint(t *testing.T) {
	var repo repository.Repo = &mocks.Repo{Config: &config.Repo{}}
	repo.GetConfig().MountPoint = "test_mountpoint"
	expectedKey := "test_mountpoint/repository_mock/master"
	key, status, err := pathBaseBuilder(repo)
	if err != nil {
		if status == PathFormatterError {
			t.Fatalf("Could not set the branch: %s", err)
		}
	}
	assert.Equal(t, key, expectedKey)
}

func TestPathCoreBuilderWithSourceRoot(t *testing.T) {
	var repo repository.Repo = &mocks.Repo{Config: &config.Repo{}}
	repo.GetConfig().SourceRoot = "first_level/second_level"
	defer func() {
		repo.GetConfig().SourceRoot = ""
	}()
	expectedKey := "/foo"
	filePath := "first_level/second_level/foo"
	key, status, err := pathCoreBuilder(repo, filePath)
	assert.NoError(t, err)
	assert.Equal(t, key, expectedKey)
	assert.Equal(t, status, PathFormatterOK)
}

func TestPathCoreBuilderWithoutSourceRoot(t *testing.T) {
	var repo repository.Repo = &mocks.Repo{Config: &config.Repo{}}
	expectedKey := "first_level/second_level/foo"
	filePath := "first_level/second_level/foo"
	key, status, err := pathCoreBuilder(repo, filePath)
	assert.NoError(t, err)
	assert.Equal(t, key, expectedKey)
	assert.Equal(t, status, PathFormatterOK)
}

func TestPathBaseBuilderSkipBranch(t *testing.T) {
	var repo repository.Repo = &mocks.Repo{Config: &config.Repo{}}
	skipBranchName := repo.GetConfig().SkipBranchName
	repo.GetConfig().SkipBranchName = true
	defer func() {
		repo.GetConfig().SkipBranchName = skipBranchName
	}()
	expectedKey := "repository_mock"
	key, status, err := pathBaseBuilder(repo)
	assert.NoError(t, err)
	assert.Equal(t, key, expectedKey)
	assert.Equal(t, status, PathFormatterOK)
}

func TestPathBaseBuilderWithBranch(t *testing.T) {
	var repo repository.Repo = &mocks.Repo{Config: &config.Repo{}}
	skipBranchName := repo.GetConfig().SkipBranchName
	repo.GetConfig().SkipBranchName = false
	defer func() {
		repo.GetConfig().SkipBranchName = skipBranchName
	}()
	expectedKey := "repository_mock/master"
	key, status, err := pathBaseBuilder(repo)
	assert.NoError(t, err)
	assert.Equal(t, key, expectedKey)
	assert.Equal(t, status, PathFormatterOK)
}

func TestPathBaseBuilderSkipRepo(t *testing.T) {
	var repo repository.Repo = &mocks.Repo{Config: &config.Repo{}}
	skipRepoName := repo.GetConfig().SkipRepoName
	repo.GetConfig().SkipRepoName = true
	defer func() {
		repo.GetConfig().SkipRepoName = skipRepoName
	}()
	expectedKey := "master"
	key, status, err := pathBaseBuilder(repo)
	assert.NoError(t, err)
	assert.Equal(t, key, expectedKey)
	assert.Equal(t, status, PathFormatterOK)
}

func TestPathBaseBuilderWithRepo(t *testing.T) {
	var repo repository.Repo = &mocks.Repo{Config: &config.Repo{}}
	skipRepoName := repo.GetConfig().SkipRepoName
	repo.GetConfig().SkipRepoName = false
	defer func() {
		repo.GetConfig().SkipRepoName = skipRepoName
	}()
	expectedKey := "repository_mock/master"
	key, status, err := pathBaseBuilder(repo)
	assert.NoError(t, err)
	assert.Equal(t, key, expectedKey)
	assert.Equal(t, status, PathFormatterOK)
}
