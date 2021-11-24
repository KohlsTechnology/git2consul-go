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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KohlsTechnology/git2consul-go/repository"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

type mockHandler struct {
	t        *testing.T
	filePath string
}

var (
	yamlFile File
	textFile File
	handler  Handler
	yamlTree map[interface{}]interface{}
	keys     map[string][]byte
)

const (
	content = "---\nei_unix_cavisson::cavisson_collector_srv: 10.206.96.18\n" +
		"ei_unix_cavisson::cavisson_port: 7891\n" +
		"ei_unix_cavisson::cavisson_java_home: \"/etc/alternatives/jre_openjdk\"\n" +
		"dict:\n" +
		"  key_1: value_1\n" +
		"  key_2:\n" +
		"    - first_elem:\n" +
		"        key_3: true\n" +
		"        key_4: 2.35\n" +
		"    - second_element: value_4\n"
)

//TestFile performs tests on implemented file handlers.
// * yaml
// * text
func TestFileHandler(t *testing.T) {
	var repo repository.Repo
	yamlTree = make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(content), &yamlTree)
	if err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(os.TempDir(), "foo.yml")
	defer os.Remove(filePath)
	err = ioutil.WriteFile(filePath, []byte(content), 0700)
	if err != nil {
		t.Fatal(err)
	}
	yamlFile = &YAMLFile{filePath}
	textFile = &TextFile{filePath}
	handler = &mockHandler{
		t:        t,
		filePath: filePath,
	}
	t.Run("TestParseYAMLFile", testParseYamlEntries)
	t.Run("TestCreateYAMLFile", func(t *testing.T) { testCreateYAMLFile(t, repo) })
	t.Run("TestDeleteYAMLFile", func(t *testing.T) { testDeleteYAMLFile(t, repo) })
	t.Run("TestCreateTextFile", func(t *testing.T) { testCreateTextFile(t, repo) })
	t.Run("TestDeleteTextFile", func(t *testing.T) { testDeleteTextFile(t, repo) })
}

//testParsNodes verfies yaml file evaluation function.
func testParseYamlEntries(t *testing.T) {
	keys := entriesToKV(yamlTree)
	if string(keys["ei_unix_cavisson::cavisson_collector_srv"]) != "10.206.96.18" {
		t.Fatal("Missing key or invalid value")
	}
	if string(keys["ei_unix_cavisson::cavisson_port"]) != "7891" {
		t.Fatal("Missing key or invalid value")
	}
	if string(keys["dict/key_1"]) != "value_1" {
		t.Fatal("Missing key or invalid value")
	}
	if string(keys["dict/key_2/0/first_elem/key_3"]) != "true" {
		t.Fatal("Missing key or invalid value")
	}
	if string(keys["dict/key_2/1/second_element"]) != "value_4" {
		t.Fatal("Missing key or invalid value")
	}
	if string(keys["dict/key_2/0/first_elem/key_4"]) != "2.35" {
		t.Fatal("Missing key or invalid value")
	}
}

func testCreateYAMLFile(t *testing.T, repo repository.Repo) {
	keys = make(map[string][]byte)
	ext := filepath.Ext(yamlFile.GetPath())
	yamlPath := strings.TrimRight(yamlFile.GetPath(), ext)
	err := yamlFile.Create(handler, repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) == 0 {
		t.Fatalf("Keys empty: %+v", keys)
	}
	for k, v := range entriesToKV(yamlTree) {
		if bytes.Equal(keys[filepath.Join(yamlPath, k)], v) {
			delete(keys, filepath.Join(yamlPath, k))
		}
	}
	if len(keys) != 0 {
		t.Fatalf("Keys not empty: %+v", keys)
	}
}

func testDeleteYAMLFile(t *testing.T, repo repository.Repo) {
	err := yamlFile.Delete(handler, repo)
	assert.NoError(t, err)
}

func testCreateTextFile(t *testing.T, repo repository.Repo) {
	keys = make(map[string][]byte)
	textPath := textFile.GetPath()
	err := textFile.Create(handler, repo)
	assert.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Equal(t, keys[textPath], []byte(content))
}

func testDeleteTextFile(t *testing.T, repo repository.Repo) {
	err := textFile.Delete(handler, repo)
	assert.NoError(t, err)
}

func (a mockHandler) PutKV(repo repository.Repo, path string, content []byte) error {
	keys[path] = content
	return nil
}

func (a mockHandler) DeleteKV(repo repository.Repo, path string) error {
	if a.filePath != path {
		return fmt.Errorf("%s differs from %s", a.filePath, path)
	}
	return nil
}

func (a mockHandler) DeleteTreeKV(repo repository.Repo, path string) error {
	filePath := strings.TrimSuffix(a.filePath, filepath.Ext(a.filePath))
	if filePath != path {
		return fmt.Errorf("%s differs from %s", a.filePath, path)
	}
	return nil
}

func (a mockHandler) HandleUpdate(repo repository.Repo) error {
	return nil
}
