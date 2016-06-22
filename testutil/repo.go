// repo.go takes care of initializing a repository.Repository slice for testing
package testutil

// import (
// 	"os"
// 	"testing"
// 	"time"
//
// 	"github.com/Cimpress-MCP/go-git2consul/config"
// 	"github.com/Cimpress-MCP/go-git2consul/repository"
// )
//
// // Genenerates a default config
// func defaultTestConfig(t *testing.T) *config.Config {
// 	return &config.Config{
// 		LocalStore: os.TempDir(),
// 		HookSvr: &config.HookSvrConfig{
// 			Port: 9000,
// 		},
// 		Repos: []*config.Repo{
// 			&config.Repo{
// 				Name:     "test-example",
// 				Url:      testRepo.Workdir(),
// 				Branches: []string{"master"},
// 				Hooks: []*config.Hook{
// 					&config.Hook{
// 						Type:     "polling",
// 						Interval: 5 * time.Second,
// 					},
// 				},
// 			},
// 		},
// 		Consul: &config.ConsulConfig{
// 			Address: "127.0.0.1:8500",
// 		},
// 	}
// }
//
// // New Repository array object from default configuration
// func ReposFromConfig(t *testing.T) []*repository.Repository {
// 	defaultConfig := defaultTestConfig(t)
// 	repoConfig := defaultConfig.Repos[0]
//
// 	repo := &repository.Repository{
// 		Repository: testRepo,
// 		Config:     repoConfig,
// 	}
//
// 	return []*repository.Repository{repo}
// }
