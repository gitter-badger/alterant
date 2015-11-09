package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	Machines  map[string]machine `yaml:"machines"`
	Tasks     map[string]task    `yaml:"tasks"`
	Encrypted []string           `yaml:"encrypted"`
	path      string
	cwd       string
}

func newConfig() *config {
	return &config{}
}

func loadConfig(file string) (*config, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := newConfig()

	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
