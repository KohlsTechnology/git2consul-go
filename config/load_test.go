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

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
)

func init() {
	log.SetHandler(discard.New())
}

func TestLoad(t *testing.T) {
	file := filepath.Join("test-fixtures", "local.json")

	_, err := Load(file)
	assert.NoError(t, err)
}

func TestLoadInvalidConfig(t *testing.T) {
	file := filepath.Join("test-fixtures", "invalid_config.json")

	_, err := Load(file)
	assert.Error(t, err)
}

func TestLoadConsulEnv(t *testing.T) {
	file := filepath.Join("test-fixtures", "local.json")

	os.Setenv("CONSUL_HTTP_ADDR", "127.0.0.1:8500")
	defer os.Unsetenv("CONSUL_HTTP_ADDR")

	os.Setenv("CONSUL_HTTP_SSL", "false")
	defer os.Unsetenv("CONSUL_HTTP_SSL")

	os.Setenv("CONSUL_HTTP_SSL_VERIFY", "false")
	defer os.Unsetenv("CONSUL_HTTP_SSL_VERIFY")

	os.Setenv("CONSUL_HTTP_TOKEN", "abcdefg123456789")
	defer os.Unsetenv("CONSUL_HTTP_TOKEN")

	_, err := Load(file)
	assert.NoError(t, err)
}
