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
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/KohlsTechnology/git2consul-go/repository"
	"gopkg.in/yaml.v2"
)

//File interface to manipulate data from various types
//of files in the KV store.
type File interface {
	Update(kv Handler, repo repository.Repo) error
	Create(kv Handler, repo repository.Repo) error
	Delete(kv Handler, repo repository.Repo) error
	GetPath() string
}

//TextFile structure
type TextFile struct {
	path string
}

//YAMLFile structure
type YAMLFile struct {
	path string
}

//Init initializes new instance of File interface based on it's extension.
func Init(path string, repo repository.Repo) File {
	config := repo.GetConfig()
	expandKeys := config.ExpandKeys
	var f File
	ext := filepath.Ext(path)
	if expandKeys {
		if ext == ".yml" {
			f = &YAMLFile{path: path}
		}
	}
	if f == nil {
		f = &TextFile{path: path}
	}
	return f
}

func getContent(f File) ([]byte, error) {
	content, err := ioutil.ReadFile(f.GetPath())
	if err != nil {
		return nil, err
	}
	return content, nil
}

//GetPath returns the path to the file.
func (f *TextFile) GetPath() string {
	return f.path
}

//Create function creates the KV store entries based on the file content.
func (f *TextFile) Create(kv Handler, repo repository.Repo) error {
	content, err := getContent(f)
	if err != nil {
		return err
	}
	err = kv.PutKV(repo, f.path, content)
	if err != nil {
		return err
	}
	return nil
}

//Update functions updates the KV store based on the file content.
func (f *TextFile) Update(kv Handler, repo repository.Repo) error {
	return f.Create(kv, repo)
}

//Delete removes the key-value pair from the KV store.
func (f *TextFile) Delete(kv Handler, repo repository.Repo) error {
	err := kv.DeleteKV(repo, f.path)
	if err != nil {
		return err
	}
	return nil
}

//Create function creates the KV store entries based on the file content.
func (f *YAMLFile) Create(kv Handler, repo repository.Repo) error {
	content, err := getContent(f)
	if err != nil {
		return err
	}
	yamlTree := make(map[interface{}]interface{})
	err = yaml.Unmarshal(content, &yamlTree)
	if err != nil {
		return err
	}
	path := f.GetPath()
	extension := filepath.Ext(path)
	fileName := strings.TrimSuffix(path, extension)
	for key, value := range entriesToKV(yamlTree) {
		err = kv.PutKV(repo, filepath.Join(fileName, key), value)
		if err != nil {
			return err
		}
	}
	return nil
}

//Update functions updates the KV store based on the file content.
func (f *YAMLFile) Update(kv Handler, repo repository.Repo) error {
	f.Delete(kv, repo) //nolint:errcheck
	return f.Create(kv, repo)
}

//Delete removes the key-value pairs from the KV store under given prefix.
func (f *YAMLFile) Delete(kv Handler, repo repository.Repo) error {
	path := f.GetPath()
	extension := filepath.Ext(path)
	fileName := strings.TrimSuffix(path, extension)
	err := kv.DeleteTreeKV(repo, fileName)
	if err != nil {
		return err
	}
	return nil
}

//GetPath returns the path to the file.
func (f *YAMLFile) GetPath() string {
	return f.path
}

func entriesToKV(node map[interface{}]interface{}) map[string][]byte {
	keys := make(map[string][]byte)
	for key, value := range node {
		switch value.(type) {
		case string:
			keys[key.(string)] = []byte(value.(string))
		case int:
			keys[key.(string)] = []byte(strconv.Itoa(value.(int)))
		case bool:
			keys[key.(string)] = []byte(strconv.FormatBool(value.(bool)))
		case float64:
			keys[key.(string)] = []byte(strconv.FormatFloat(value.(float64), 'f', 2, 64))
		case map[interface{}]interface{}:
			for k, v := range entriesToKV(value.(map[interface{}]interface{})) {
				keys[filepath.Join(key.(string), k)] = v
			}
		case []interface{}:
			for index, item := range value.([]interface{}) {
				for k, v := range entriesToKV(item.(map[interface{}]interface{})) {
					keys[filepath.Join(key.(string), strconv.Itoa(index), k)] = v
				}
			}
		}
	}
	return keys
}
