package config

import (
	"path/filepath"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
)

func init() {
	log.SetHandler(discard.New())
}

func TestLoad(t *testing.T) {
	file := filepath.Join("test-fixtures", "local.json")

	_, err := Load(file)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoad_invalidConfig(t *testing.T) {
	file := filepath.Join("test-fixtures", "invalid_config.json")

	_, err := Load(file)
	if err == nil {
		t.Fatal("Expected failure")
	}
}
