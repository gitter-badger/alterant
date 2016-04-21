package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/autonomy/alterant/hasher"
	"github.com/autonomy/alterant/task"
	"github.com/deckarep/golang-set"
)

// Config represents `machine.yaml`
type Config struct {
	Machine     string
	Environment map[string]string
	Tasks       []*task.Task
	SHA1        string
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Environment map[string]string     `yaml:"environment"`
		Tasks       map[string]*task.Task `yaml:"tasks"`
	}

	err := unmarshal(&aux)
	if err != nil {
		return err
	}

	for name, task := range aux.Tasks {
		task.Name = name
	}

	b, err := yaml.Marshal(&aux)
	if err != nil {
		return err
	}

	tasks, err := resolveDependencies(aux.Tasks)
	if err != nil {
		return err
	}

	*c = Config{
		Environment: aux.Environment,
		Tasks:       tasks,
		SHA1:        hasher.SHA1FromBytes(b),
	}

	return nil
}
func newConfig() *Config {
	return &Config{}
}

func resolveDependencies(tasks map[string]*task.Task) ([]*task.Task, error) {
	var t []*task.Task
	taskDependencies := make(map[string]mapset.Set)

	for taskName, task := range tasks {
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
			return nil, fmt.Errorf("Circular dependency found.")
		}

		for name := range readySet.Iter() {
			delete(taskDependencies, name.(string))
			t = append(t, tasks[name.(string)])
		}

		for name, deps := range taskDependencies {
			diff := deps.Difference(readySet)
			taskDependencies[name] = diff
		}
	}

	return t, nil
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
	if _, err = os.Stat(path.Join(cwd, file)); os.IsNotExist(err) {
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
