package main

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type config struct {
	Machines  map[string]*machine `yaml:"machines"`
	Tasks     map[string]*task    `yaml:"tasks"`
	Encrypted []string            `yaml:"encrypted"`
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

func requireConfig() (*config, error) {
	// require that the config is named 'alter.yaml'
	file := "alter.yaml"

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// require that alter.yaml exists in the cwd
	if _, err := os.Stat(path.Join(cwd, file)); os.IsNotExist(err) {
		return nil, err
	}

	cfg, err := loadConfig(file)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
