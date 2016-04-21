package task

import (
	"github.com/autonomy/alterant/command"
	"github.com/autonomy/alterant/hasher"
	"github.com/autonomy/alterant/link"
	"gopkg.in/yaml.v2"
)

// Task represents a task
type Task struct {
	Dependencies []string
	Links        map[string]*link.Link
	Commands     map[string]*command.Command
	Name         string
	Queued       bool
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

	links := make(map[string]*link.Link)
	for _, link := range aux.Links {
		links[link.SHA1] = link
	}

	commands := make(map[string]*command.Command)
	for _, command := range aux.Commands {
		commands[command.SHA1] = command
	}

	b, err := yaml.Marshal(&aux)
	if err != nil {
		return err
	}

	*t = Task{
		Dependencies: aux.Dependencies,
		Links:        links,
		Commands:     commands,
		Queued:       true,
		SHA1:         hasher.SHA1FromBytes(b),
	}

	return nil
}
