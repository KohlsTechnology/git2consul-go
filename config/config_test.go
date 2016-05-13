package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	file := filepath.Join("test-fixtures", "local.json")
	content, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	config := &Config{}

	json.Unmarshal(content, config)

	jsonFromType, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	//Clean strings from trailing newlines
	want := strings.TrimSpace(string(content))
	got := strings.TrimSpace(string(jsonFromType))

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("JSON mistatch. \nfile:\n%s\ngot:\n%s\n", want, got)
	}
}
