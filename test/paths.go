package test

import "path"

const (
	TestConfigConst = "../config/test-fixtures/local.json"
	TestRepoConst   = "../repository/test-fixtures/example"
)

func TestConfig() string {
	return path.Clean(TestConfigConst)
}

func TestRepo() string {
	return path.Clean(TestRepoConst)
}
