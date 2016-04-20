package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	file := filepath.Join("test-fixtures", "default.json")
	content, err := ioutil.ReadFile(file)
	if err != nil {
		os.Exit(1)
	}

	config := &Config{}

	json.Unmarshal(content, config)

	got, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(content, got) {
		t.Fatalf("JSON mistatch. file:\n%s\ngot:\n%s\n", content, got)
	}
}
