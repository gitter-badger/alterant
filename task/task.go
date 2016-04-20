package task

import (
	"github.com/autonomy/alterant/cache"
	"github.com/autonomy/alterant/command"
	"github.com/autonomy/alterant/link"
	"gopkg.in/yaml.v2"
)

// Task represents a task
type Task struct {
	Dependencies []string
	Links        []*link.Link
	Commands     []*command.Command
	Name         string
	SHA1         string
}

func (t *Task) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Dependencies []string           `yaml:"dependencies"`
		Links        []*link.Link       `yaml:"links"`
		Commands     []*command.Command `yaml:"commands"`
	}

	err := unmarshal(&aux)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(&aux)
	if err != nil {
		return err
	}

	*t = Task{
		Dependencies: aux.Dependencies,
		Links:        aux.Links,
		Commands:     aux.Commands,
		SHA1:         cache.SHAFromBytes(b),
	}

	return nil
}
