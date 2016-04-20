package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/autonomy/alterant/cache"
	"github.com/autonomy/alterant/task"
	"github.com/deckarep/golang-set"
)

// Config represents `machine.yaml`
type Config struct {
	Environment map[string]string     `yaml:"environment"`
	Tasks       map[string]*task.Task `yaml:"tasks"`
	Machine     string
	Order       []*task.Task
	Sha1        string
}

func newConfig() *Config {
	return &Config{}
}

func resolveDependencies(cfg *Config) error {
	taskDependencies := make(map[string]mapset.Set)

	for taskName, task := range cfg.Tasks {
		task.Name = taskName
		dependencySet := mapset.NewSet()

		for _, dep := range task.Dependencies {
			dependencySet.Add(dep)
		}

		taskDependencies[taskName] = dependencySet
	}

	for len(taskDependencies) != 0 {
		readySet := mapset.NewSet()

		for name, deps := range taskDependencies {
			if deps.Cardinality() == 0 {
				readySet.Add(name)
			}
		}

		if readySet.Cardinality() == 0 {
			return fmt.Errorf("Circular dependency found.")
		}

		for name := range readySet.Iter() {
			delete(taskDependencies, name.(string))
			cfg.Order = append(cfg.Order, cfg.Tasks[name.(string)])
		}

		for name, deps := range taskDependencies {
			diff := deps.Difference(readySet)
			taskDependencies[name] = diff
		}
	}

	return nil
}

func loadConfig(file string) (*Config, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	sha := cache.SHAFromBytes(bytes)

	cfg := newConfig()

	cfg.Sha1 = sha

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
	if _, err = os.Stat(path.Join(cwd, file)); os.IsNotExist(err) {
		return nil, err
	}

	cfg, err := loadConfig(file)
	if err != nil {
		return nil, err
	}

	err = resolveDependencies(cfg)
	if err != nil {
		return nil, err
	}

	// required by the custom unmarshalling of SymlinkTarget and SymlinkDestination
	os.Setenv("MACHINE", machine)

	cfg.Machine = machine

	return cfg, nil
}
