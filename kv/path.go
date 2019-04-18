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
	"fmt"
	"path"
	"strings"

	"github.com/KohlsTechnology/git2consul-go/repository"
)

//Status codes for path formatting
const (
	SourceRootNotInPrefix = iota
	PathFormatterOK
	PathFormatterError
)

func getItemKey(repo repository.Repo, filePath string) (string, int, error) {
	return pathHandler(repo, filePath)
}
func pathHandler(repo repository.Repo, filePath string) (string, int, error) {
	filePath = strings.TrimPrefix(filePath, repository.WorkDir(repo))
	basePath, status, err := pathBaseBuilder(repo)
	if err != nil {
		return "", status, fmt.Errorf("Couldn't format the base of the path: %s", err)
	}
	corePath, status, err := pathCoreBuilder(repo, filePath)
	if err != nil {
		return "", status, err
	}
	path := path.Join(basePath, corePath)

	return path, PathFormatterOK, nil
}

func pathBaseBuilder(repo repository.Repo) (string, int, error) {
	config := repo.GetConfig()
	mountPoint := config.MountPoint
	key := ""
	repoName := ""
	if !config.SkipRepoName {
		repoName = repo.Name()
	}
	branchName, err := getBranchName(repo)
	if err != nil {
		return "", PathFormatterError, err
	}
	if len(mountPoint) > 0 {
		key = path.Join(mountPoint, repoName, branchName)
	} else {
		key = path.Join(repoName, branchName)
	}

	return key, PathFormatterOK, nil
}

func pathCoreBuilder(repo repository.Repo, filePath string) (string, int, error) {
	config := repo.GetConfig()
	sourceRoot := config.SourceRoot

	if len(sourceRoot) > 0 {
		if !strings.Contains(filePath, sourceRoot) {
			return "", SourceRootNotInPrefix, fmt.Errorf("Path: \"%s\" doesn't match the source_root: \"%s\"", filePath, sourceRoot)
		}
		filePath = strings.TrimPrefix(filePath, sourceRoot)
	}

	return filePath, PathFormatterOK, nil
}

func getBranchName(repo repository.Repo) (string, error) {
	config := repo.GetConfig()
	if config.SkipBranchName {
		return "", nil
	}
	branch, err := repo.Head()
	if err != nil {
		return "", err
	}
	return branch.Name().Short(), nil
}
