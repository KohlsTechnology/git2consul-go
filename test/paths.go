package test

import "path"

const (
	TestConfigConst = "../config/test-fixtures/local.json"
	TestRepoConst   = "../repository/test-fixtures/example"
)

func DefaultConfigPath() string {
	return path.Clean(TestConfigConst)
}

func DefaultRepoPath() string {
	return path.Clean(TestRepoConst)
}
