package main

import (
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestLoadConfig(t *testing.T) {
	sc := &SafeConfig{
		C: &Config{},
	}

	err := sc.reloadConfig("testdata/blackbox-good.yml")
	if err != nil {
		t.Errorf("Error loading config %v: %v", "blackbox.yml", err)
	}
}

func TestLoadBadConfigs(t *testing.T) {
	sc := &SafeConfig{
		C: &Config{},
	}
	tests := []struct {
		ConfigFile    string
		ExpectedError string
	}{
		{
			ConfigFile:    "testdata/blackbox-bad.yml",
			ExpectedError: "unknown fields in dns probe: invalid_extra_field",
		},
		{
			ConfigFile:    "testdata/invalid-dns-module.yml",
			ExpectedError: "Query name must be set for DNS module",
		},
	}
	for i, test := range tests {
		err := sc.reloadConfig(test.ConfigFile)
		if err.Error() != test.ExpectedError {
			t.Errorf("In case %v:\nExpected:\n%v\nGot:\n%v", i, test.ExpectedError, err.Error())
		}
	}
}

func TestHideConfigSecrets(t *testing.T) {
	sc := &SafeConfig{
		C: &Config{},
	}

	err := sc.reloadConfig("testdata/blackbox-good.yml")
	if err != nil {
		t.Errorf("Error loading config %v: %v", "testdata/blackbox-good.yml", err)
	}

	// String method must not reveal authentication credentials.
	sc.RLock()
	c, err := yaml.Marshal(sc.C)
	sc.RUnlock()
	if err != nil {
		t.Errorf("Error marshalling config: %v", err)
	}
	if strings.Contains(string(c), "mysecret") {
		t.Fatal("config's String method reveals authentication credentials.")
	}
}
