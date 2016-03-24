package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/autonomy/alterant/task"

	"gopkg.in/yaml.v2"
)

// Config represents `machine.yaml`
type Config struct {
	Order       []string              `yaml:"order"`
	Environment map[string]string     `yaml:"environment"`
	Tasks       map[string]*task.Task `yaml:"tasks"`
	Machine     string
}

func newConfig() *Config {
	return &Config{}
}

func loadConfig(file string) (*Config, error) {
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

// AcquireConfig unmarshalls alter.yaml and returns the representation as a Config
func AcquireConfig(machine string) (*Config, error) {
	file := machine + ".yaml"

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// require that machine + ".yaml" exists in the cwd
	if _, err := os.Stat(path.Join(cwd, file)); os.IsNotExist(err) {
		return nil, err
	}

	cfg, err := loadConfig(file)
	if err != nil {
		return nil, err
	}

	// required by the custom unmarshalling of SymlinkTarget and SymlinkDestination
	os.Setenv("MACHINE", machine)

	cfg.Machine = machine

	return cfg, nil
}
